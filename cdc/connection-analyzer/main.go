package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ConnectionAnalyzer struct {
	driver neo4j.DriverWithContext
}

type ConnectionRecommendation struct {
	Type        string                 `json:"type"`
	Score       float64               `json:"score"`
	Reason      string                `json:"reason"`
	Business    *BusinessInfo         `json:"business,omitempty"`
	User        *UserInfo             `json:"user,omitempty"`
	Skills      []SkillMatch          `json:"skills,omitempty"`
	Projects    []ProjectMatch        `json:"projects,omitempty"`
	CompatibilityFactors []string    `json:"compatibilityFactors"`
}

type BusinessInfo struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	BusinessType     string `json:"businessType"`
	BusinessCategory string `json:"businessCategory"`
	BusinessPhase    string `json:"businessPhase"`
	Description      string `json:"description"`
	Website          string `json:"website"`
}

type UserInfo struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type SkillMatch struct {
	SkillID          uint    `json:"skillId"`
	SkillName        string  `json:"skillName"`
	SkillCategory    string  `json:"skillCategory"`
	UserProficiency  string  `json:"userProficiency"`
	ProjectImportance string `json:"projectImportance"`
	MatchScore       float64 `json:"matchScore"`
}

type ProjectMatch struct {
	ProjectID   uint   `json:"projectId"`
	ProjectName string `json:"projectName"`
	Status      string `json:"status"`
	MatchReason string `json:"matchReason"`
}

func main() {
	fmt.Println("Starting Connection Analysis Service...")

	memgraphHost := getEnv("MEMGRAPH_HOST", "memgraph")
	memgraphPort := getEnv("MEMGRAPH_PORT", "7687")
	
	var driver neo4j.DriverWithContext
	var err error
	
	for i := 0; i < 5; i++ {
		fmt.Printf("Attempting to connect to Memgraph at %s:%s (attempt %d/5)...\n", memgraphHost, memgraphPort, i+1)
		driver, err = neo4j.NewDriverWithContext(
			fmt.Sprintf("bolt://%s:%s", memgraphHost, memgraphPort),
			neo4j.NoAuth(),
		)
		if err == nil {
			ctx := context.Background()
			err = driver.VerifyConnectivity(ctx)
			if err == nil {
				fmt.Println("Successfully connected to Memgraph")
				break
			}
			driver.Close(ctx)
		}
		
		if i < 4 {
			fmt.Printf("Connection failed: %v. Retrying in %d seconds...\n", err, (i+1)*2)
			time.Sleep(time.Duration((i+1)*2) * time.Second)
		}
	}
	
	if err != nil {
		log.Fatalf("Failed to connect to Memgraph after 5 attempts: %v", err)
	}
	defer driver.Close(context.Background())

	analyzer := &ConnectionAnalyzer{driver: driver}

	router := gin.Default()
	
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	api := router.Group("/api/v1")
	{
		api.GET("/connections/complementary/:userId", analyzer.getComplementaryPartners)
		api.GET("/connections/alliance/:userId", analyzer.getAlliancePartners)
		api.GET("/connections/mastermind/:userId", analyzer.getMastermindPartners)
		api.GET("/connections/recommendations/:userId", analyzer.getAllRecommendations)
		api.GET("/connections/analysis/:userId", analyzer.getConnectionAnalysis)
	}

	port := getEnv("PORT", "8082")
	fmt.Printf("Connection Analysis Service running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func (ca *ConnectionAnalyzer) getComplementaryPartners(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	session := ca.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	query := `
	MATCH (u:User {id: $userId})-[:OPERATES]->(b1:Business)
	MATCH (b2:Business)<-[:OPERATES]-(u2:User)
	WHERE b1.id <> b2.id AND u.id <> u2.id
	  AND (
		(b1.businessType = b2.businessType AND b1.businessCategory <> b2.businessCategory) OR
		(b1.businessType <> b2.businessType AND b1.businessCategory = b2.businessCategory)
	  )
	WITH u, b1, b2, u2,
		CASE 
			WHEN b1.businessType = b2.businessType AND b1.businessCategory <> b2.businessCategory THEN 0.9
			WHEN b1.businessType <> b2.businessType AND b1.businessCategory = b2.businessCategory THEN 0.8
			ELSE 0.7
		END as compatibilityScore
	OPTIONAL MATCH (b1)-[conn:CONNECTED_TO]->(b2)
	RETURN b2.id as businessId, b2.name as businessName, b2.businessType as businessType,
		   b2.businessCategory as businessCategory, b2.businessPhase as businessPhase,
		   b2.description as description, b2.website as website,
		   u2.id as userId, u2.firstName as firstName, u2.lastName as lastName, u2.email as email,
		   compatibilityScore, 
		   CASE WHEN conn IS NOT NULL THEN conn.status ELSE 'not_connected' END as connectionStatus
	ORDER BY compatibilityScore DESC
	LIMIT 20
	`

	result, err := session.Run(context.Background(), query, map[string]interface{}{"userId": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query complementary partners"})
		return
	}

	var recommendations []ConnectionRecommendation
	for result.Next(context.Background()) {
		record := result.Record()
		
		businessInfo := &BusinessInfo{
			ID:               uint(record.Values[0].(int64)),
			Name:             record.Values[1].(string),
			BusinessType:     record.Values[2].(string),
			BusinessCategory: record.Values[3].(string),
			BusinessPhase:    record.Values[4].(string),
			Description:      getStringValue(record.Values[5]),
			Website:          getStringValue(record.Values[6]),
		}
		
		userInfo := &UserInfo{
			ID:        uint(record.Values[7].(int64)),
			FirstName: record.Values[8].(string),
			LastName:  getStringValue(record.Values[9]),
			Email:     record.Values[10].(string),
		}
		
		score := record.Values[11].(float64)
		connectionStatus := record.Values[12].(string)
		
		reason := fmt.Sprintf("Complementary business types: %s + %s", 
			getBusinessType(record.Values[2].(string)), 
			getBusinessType(record.Values[2].(string)))
		
		compatibilityFactors := []string{
			fmt.Sprintf("Business type compatibility: %.0f%%", score*100),
			fmt.Sprintf("Current connection status: %s", connectionStatus),
		}
		
		recommendations = append(recommendations, ConnectionRecommendation{
			Type:                "COMPLEMENTARY_PARTNER",
			Score:               score,
			Reason:              reason,
			Business:            businessInfo,
			User:                userInfo,
			CompatibilityFactors: compatibilityFactors,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"type": "COMPLEMENTARY_PARTNERS",
		"count": len(recommendations),
		"recommendations": recommendations,
	})
}

func (ca *ConnectionAnalyzer) getAlliancePartners(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	session := ca.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	query := `
	MATCH (u:User {id: $userId})-[:OPERATES]->(b1:Business)
	MATCH (u)-[:HAS_SKILL]->(s:Skill)<-[:HAS_SKILL]-(u2:User)-[:OPERATES]->(b2:Business)
	WHERE b1.id <> b2.id AND u.id <> u2.id
	WITH u, b1, b2, u2, collect(s) as sharedSkills
	WITH u, b1, b2, u2, sharedSkills,
		size(sharedSkills) * 0.4 + 
		CASE 
			WHEN b1.businessPhase = b2.businessPhase THEN 0.3
			ELSE 0.1
		END +
		CASE 
			WHEN b1.businessCategory = b2.businessCategory THEN 0.3
			ELSE 0.1
		END as allianceScore
	WHERE allianceScore > 0.5
	OPTIONAL MATCH (b1)-[conn:CONNECTED_TO]->(b2)
	RETURN b2.id as businessId, b2.name as businessName, b2.businessType as businessType,
		   b2.businessCategory as businessCategory, b2.businessPhase as businessPhase,
		   b2.description as description, b2.website as website,
		   u2.id as userId, u2.firstName as firstName, u2.lastName as lastName, u2.email as email,
		   allianceScore, sharedSkills,
		   CASE WHEN conn IS NOT NULL THEN conn.status ELSE 'not_connected' END as connectionStatus
	ORDER BY allianceScore DESC
	LIMIT 20
	`

	result, err := session.Run(context.Background(), query, map[string]interface{}{"userId": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query alliance partners"})
		return
	}

	var recommendations []ConnectionRecommendation
	for result.Next(context.Background()) {
		record := result.Record()
		
		businessInfo := &BusinessInfo{
			ID:               uint(record.Values[0].(int64)),
			Name:             record.Values[1].(string),
			BusinessType:     record.Values[2].(string),
			BusinessCategory: record.Values[3].(string),
			BusinessPhase:    record.Values[4].(string),
			Description:      getStringValue(record.Values[5]),
			Website:          getStringValue(record.Values[6]),
		}
		
		userInfo := &UserInfo{
			ID:        uint(record.Values[7].(int64)),
			FirstName: record.Values[8].(string),
			LastName:  getStringValue(record.Values[9]),
			Email:     record.Values[10].(string),
		}
		
		score := record.Values[11].(float64)
		sharedSkills := record.Values[12].([]interface{})
		connectionStatus := record.Values[13].(string)
		
		var skillMatches []SkillMatch
		for _, skill := range sharedSkills {
			if skillNode, ok := skill.(neo4j.Node); ok {
				props := skillNode.Props
				skillMatches = append(skillMatches, SkillMatch{
					SkillID:         uint(props["id"].(int64)),
					SkillName:       props["name"].(string),
					SkillCategory:   props["category"].(string),
					UserProficiency: "intermediate",
					MatchScore:      1.0,
				})
			}
		}
		
		reason := fmt.Sprintf("Shared %d skills with potential for project collaboration", len(sharedSkills))
		
		compatibilityFactors := []string{
			fmt.Sprintf("Shared skills: %d", len(sharedSkills)),
			fmt.Sprintf("Alliance potential: %.0f%%", score*100),
			fmt.Sprintf("Current connection status: %s", connectionStatus),
		}
		
		recommendations = append(recommendations, ConnectionRecommendation{
			Type:                "ALLIANCE_PARTNER",
			Score:               score,
			Reason:              reason,
			Business:            businessInfo,
			User:                userInfo,
			Skills:              skillMatches,
			CompatibilityFactors: compatibilityFactors,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"type": "ALLIANCE_PARTNERS",
		"count": len(recommendations),
		"recommendations": recommendations,
	})
}

func (ca *ConnectionAnalyzer) getMastermindPartners(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	session := ca.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	query := `
	MATCH (u:User {id: $userId})-[:OPERATES]->(b1:Business)
	MATCH (u2:User)-[:OPERATES]->(b2:Business)
	WHERE b1.id <> b2.id AND u.id <> u2.id
	WITH u, b1, b2, u2
	MATCH (u)-[:HAS_SKILL]->(s1:Skill)
	MATCH (u2)-[:HAS_SKILL]->(s2:Skill)
	WHERE s1.category <> s2.category
	WITH u, b1, b2, u2, collect(DISTINCT s1) as userSkills, collect(DISTINCT s2) as partnerSkills
	WITH u, b1, b2, u2, userSkills, partnerSkills,
		CASE 
			WHEN b1.businessPhase = 'Startup' AND b2.businessPhase IN ['Growth', 'Mature'] THEN 0.9
			WHEN b1.businessPhase = 'Growth' AND b2.businessPhase IN ['Startup', 'Mature'] THEN 0.8
			WHEN b1.businessPhase = 'Mature' AND b2.businessPhase IN ['Growth', 'Exit'] THEN 0.7
			ELSE 0.5
		END as phaseCompatibility,
		size(partnerSkills) * 0.3 + size(userSkills) * 0.2 as skillComplementarity
	WHERE phaseCompatibility > 0.6 OR skillComplementarity > 0.3
	WITH u, b1, b2, u2, userSkills, partnerSkills, phaseCompatibility, skillComplementarity,
		phaseCompatibility * 0.6 + skillComplementarity * 0.4 as mastermindScore
	WHERE mastermindScore > 0.5
	OPTIONAL MATCH (b1)-[conn:CONNECTED_TO]->(b2)
	RETURN b2.id as businessId, b2.name as businessName, b2.businessType as businessType,
		   b2.businessCategory as businessCategory, b2.businessPhase as businessPhase,
		   b2.description as description, b2.website as website,
		   u2.id as userId, u2.firstName as firstName, u2.lastName as lastName, u2.email as email,
		   mastermindScore, phaseCompatibility, skillComplementarity,
		   CASE WHEN conn IS NOT NULL THEN conn.status ELSE 'not_connected' END as connectionStatus
	ORDER BY mastermindScore DESC
	LIMIT 20
	`

	result, err := session.Run(context.Background(), query, map[string]interface{}{"userId": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query mastermind partners"})
		return
	}

	var recommendations []ConnectionRecommendation
	for result.Next(context.Background()) {
		record := result.Record()
		
		businessInfo := &BusinessInfo{
			ID:               uint(record.Values[0].(int64)),
			Name:             record.Values[1].(string),
			BusinessType:     record.Values[2].(string),
			BusinessCategory: record.Values[3].(string),
			BusinessPhase:    record.Values[4].(string),
			Description:      getStringValue(record.Values[5]),
			Website:          getStringValue(record.Values[6]),
		}
		
		userInfo := &UserInfo{
			ID:        uint(record.Values[7].(int64)),
			FirstName: record.Values[8].(string),
			LastName:  getStringValue(record.Values[9]),
			Email:     record.Values[10].(string),
		}
		
		score := record.Values[11].(float64)
		phaseCompatibility := record.Values[12].(float64)
		skillComplementarity := record.Values[13].(float64)
		connectionStatus := record.Values[14].(string)
		
		reason := fmt.Sprintf("Complementary business phase (%s → %s) and skill sets", 
			businessInfo.BusinessPhase, businessInfo.BusinessPhase)
		
		compatibilityFactors := []string{
			fmt.Sprintf("Phase compatibility: %.0f%%", phaseCompatibility*100),
			fmt.Sprintf("Skill complementarity: %.0f%%", skillComplementarity*100),
			fmt.Sprintf("Current connection status: %s", connectionStatus),
		}
		
		recommendations = append(recommendations, ConnectionRecommendation{
			Type:                "MASTERMIND_PARTNER",
			Score:               score,
			Reason:              reason,
			Business:            businessInfo,
			User:                userInfo,
			CompatibilityFactors: compatibilityFactors,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"type": "MASTERMIND_PARTNERS",
		"count": len(recommendations),
		"recommendations": recommendations,
	})
}

func (ca *ConnectionAnalyzer) getAllRecommendations(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	complementaryChan := make(chan []ConnectionRecommendation, 1)
	allianceChan := make(chan []ConnectionRecommendation, 1)
	mastermindChan := make(chan []ConnectionRecommendation, 1)

	go func() {
		recommendations := ca.getRecommendationsByType(userID, "COMPLEMENTARY")
		complementaryChan <- recommendations
	}()

	go func() {
		recommendations := ca.getRecommendationsByType(userID, "ALLIANCE")
		allianceChan <- recommendations
	}()

	go func() {
		recommendations := ca.getRecommendationsByType(userID, "MASTERMIND")
		mastermindChan <- recommendations
	}()

	complementary := <-complementaryChan
	alliance := <-allianceChan
	mastermind := <-mastermindChan

	c.JSON(http.StatusOK, gin.H{
		"userId": userID,
		"recommendations": map[string]interface{}{
			"complementary": gin.H{
				"type": "COMPLEMENTARY_PARTNERS",
				"count": len(complementary),
				"recommendations": complementary,
			},
			"alliance": gin.H{
				"type": "ALLIANCE_PARTNERS",
				"count": len(alliance),
				"recommendations": alliance,
			},
			"mastermind": gin.H{
				"type": "MASTERMIND_PARTNERS",
				"count": len(mastermind),
				"recommendations": mastermind,
			},
		},
	})
}

func (ca *ConnectionAnalyzer) getConnectionAnalysis(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	session := ca.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	query := `
	MATCH (u:User {id: $userId})-[:OPERATES]->(b:Business)
	OPTIONAL MATCH (b)-[conn:CONNECTS_TO]->(other:Business)
	OPTIONAL MATCH (u)-[:HAS_SKILL]->(s:Skill)
	OPTIONAL MATCH (p:Project)-[:MANAGES]->(u)
	RETURN b.id as businessId, b.name as businessName, b.businessType as businessType,
		   b.businessPhase as businessPhase,
		   count(DISTINCT conn) as totalConnections,
		   count(DISTINCT CASE WHEN conn.status = 'active' THEN conn END) as activeConnections,
		   count(DISTINCT s) as totalSkills,
		   count(DISTINCT CASE WHEN s.category = 'Technology' THEN s END) as techSkills,
		   count(DISTINCT CASE WHEN s.category = 'Business' THEN s END) as businessSkills,
		   count(DISTINCT p) as totalProjects,
		   count(DISTINCT CASE WHEN p.projectStatus = 'active' THEN p END) as activeProjects
	`

	result, err := session.Run(context.Background(), query, map[string]interface{}{"userId": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze connections"})
		return
	}

	if result.Next(context.Background()) {
		record := result.Record()
		
		analysis := map[string]interface{}{
			"businessId":         record.Values[0],
			"businessName":       record.Values[1],
			"businessType":       record.Values[2],
			"businessPhase":      record.Values[3],
			"totalConnections":   record.Values[4],
			"activeConnections":  record.Values[5],
			"totalSkills":        record.Values[6],
			"techSkills":         record.Values[7],
			"businessSkills":     record.Values[8],
			"totalProjects":      record.Values[9],
			"activeProjects":     record.Values[10],
			"connectionStrength": calculateConnectionStrength(record.Values[4].(int64), record.Values[5].(int64)),
			"skillDiversity":     calculateSkillDiversity(record.Values[6].(int64), record.Values[7].(int64), record.Values[8].(int64)),
		}

		c.JSON(http.StatusOK, gin.H{
			"userId": userID,
			"analysis": analysis,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}

func (ca *ConnectionAnalyzer) getRecommendationsByType(userID uint64, connectionType string) []ConnectionRecommendation {
	return []ConnectionRecommendation{}
}

func getStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

func getBusinessType(businessType string) string {
	switch businessType {
	case "Technology":
		return "Tech"
	case "Consulting":
		return "Consulting"
	case "Manufacturing":
		return "Manufacturing"
	case "Services":
		return "Services"
	case "Retail":
		return "Retail"
	default:
		return "Other"
	}
}

func calculateConnectionStrength(total, active int64) string {
	if total == 0 {
		return "None"
	}
	ratio := float64(active) / float64(total)
	if ratio > 0.8 {
		return "Strong"
	} else if ratio > 0.5 {
		return "Moderate"
	} else {
		return "Weak"
	}
}

func calculateSkillDiversity(total, tech, business int64) string {
	if total == 0 {
		return "None"
	}
	techRatio := float64(tech) / float64(total)
	businessRatio := float64(business) / float64(total)
	
	if techRatio > 0.6 || businessRatio > 0.6 {
		return "Specialized"
	} else if techRatio > 0.3 && businessRatio > 0.3 {
		return "Balanced"
	} else {
		return "Diverse"
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
