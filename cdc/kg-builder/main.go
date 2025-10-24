package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	_ "github.com/go-sql-driver/mysql"
)

type KnowledgeGraphBuilder struct {
	consumer sarama.ConsumerGroup
	driver   neo4j.DriverWithContext
	ctx      context.Context
}

type User struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Business struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	BusinessType     string    `json:"business_type"`
	BusinessCategory string    `json:"business_category"`
	BusinessPhase    string    `json:"business_phase"`
	Description      string    `json:"description"`
	Website          string    `json:"website"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Project struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Skill struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BusinessConnection struct {
	ID             uint      `json:"id"`
	Business1ID    uint      `json:"business1_id"`
	Business2ID    uint      `json:"business2_id"`
	ConnectionType string    `json:"connection_type"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserSkill struct {
	UserID           uint      `json:"user_id"`
	SkillID          uint      `json:"skill_id"`
	ProficiencyLevel string    `json:"proficiency_level"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ProjectSkill struct {
	ProjectID       uint      `json:"project_id"`
	SkillID         uint      `json:"skill_id"`
	ImportanceLevel string    `json:"importance_level"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DebeziumMessage struct {
	Schema struct {
		Type   string `json:"type"`
		Fields []struct {
			Type     string `json:"type"`
			Optional bool   `json:"optional"`
			Field    string `json:"field"`
			Name     string `json:"name"`
		} `json:"fields"`
		Optional bool   `json:"optional"`
		Name     string `json:"name"`
	} `json:"schema"`
	Payload struct {
		Before map[string]interface{} `json:"before"`
		After  map[string]interface{} `json:"after"`
		Source struct {
			Version   string `json:"version"`
			Connector string `json:"connector"`
			Name      string `json:"name"`
			TsMs      int64  `json:"ts_ms"`
			Snapshot  string `json:"snapshot"`
			DB        string `json:"db"`
			Sequence  string `json:"sequence"`
			Table     string `json:"table"`
			ServerID  int    `json:"server_id"`
			Gtid      string `json:"gtid"`
			File      string `json:"file"`
			Pos       int    `json:"pos"`
			Row       int    `json:"row"`
			Thread    int    `json:"thread"`
			Query     string `json:"query"`
		} `json:"source"`
		Op    string `json:"op"`
		TsMs  int64  `json:"ts_ms"`
		TsUs  int64  `json:"ts_us"`
		TsNs  int64  `json:"ts_ns"`
		TsSec int64  `json:"ts_sec"`
	} `json:"payload"`
}

func main() {
	fmt.Println("Starting Knowledge Graph Builder...")

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

	kafkaBrokers := getEnv("KAFKA_BROKERS", "kafka:29092")
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(strings.Split(kafkaBrokers, ","), "kg-builder-group", config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumer.Close()

	kgb := &KnowledgeGraphBuilder{
		consumer: consumer,
		driver:   driver,
		ctx:      context.Background(),
	}

	if err := kgb.initializeSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}
	fmt.Println("Knowledge graph schema initialized")

	fmt.Println("Loading existing data from MySQL...")
	if err := kgb.loadExistingData(); err != nil {
		log.Printf("Warning: Failed to load existing data: %v", err)
	} else {
		fmt.Println("Existing data loaded successfully")
	}

	fmt.Println("Starting to consume CDC messages...")
	for {
		err := consumer.Consume(kgb.ctx, []string{"tia-db.tia-dev.users", "tia-db.tia-dev.businesses", "tia-db.tia-dev.projects", "tia-db.tia-dev.skills", "tia-db.tia-dev.business_connections", "tia-db.tia-dev.user_skills", "tia-db.tia-dev.project_skills"}, kgb)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}
	}
}

func (kgb *KnowledgeGraphBuilder) initializeSchema() error {
	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	constraints := []string{
		"CREATE CONSTRAINT ON (u:User) ASSERT u.id IS UNIQUE",
		"CREATE CONSTRAINT ON (b:Business) ASSERT b.id IS UNIQUE",
		"CREATE CONSTRAINT ON (p:Project) ASSERT p.id IS UNIQUE",
		"CREATE CONSTRAINT ON (s:Skill) ASSERT s.id IS UNIQUE",
		"CREATE INDEX ON :User(email)",
		"CREATE INDEX ON :Business(name)",
		"CREATE INDEX ON :Skill(name)",
	}

	for _, constraint := range constraints {
		if _, err := session.Run(kgb.ctx, constraint, nil); err != nil {
			log.Printf("Warning: Failed to create constraint/index: %v", err)
		}
	}

	fmt.Println("Knowledge graph schema initialized")
	return nil
}

func (kgb *KnowledgeGraphBuilder) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (kgb *KnowledgeGraphBuilder) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (kgb *KnowledgeGraphBuilder) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := kgb.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (kgb *KnowledgeGraphBuilder) processMessage(message *sarama.ConsumerMessage) error {
	var debeziumMsg DebeziumMessage
	if err := json.Unmarshal(message.Value, &debeziumMsg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}

	tableName := debeziumMsg.Payload.Source.Table
	operation := debeziumMsg.Payload.Op

	fmt.Printf("Processing %s operation on table %s\n", operation, tableName)

	switch tableName {
	case "users":
		return kgb.processUserChange(operation, debeziumMsg.Payload.Before, debeziumMsg.Payload.After)
	case "businesses":
		return kgb.processBusinessChange(operation, debeziumMsg.Payload.Before, debeziumMsg.Payload.After)
	case "projects":
		return kgb.processProjectChange(operation, debeziumMsg.Payload.Before, debeziumMsg.Payload.After)
	case "skills":
		return kgb.processSkillChange(operation, debeziumMsg.Payload.Before, debeziumMsg.Payload.After)
	case "business_connections":
		return kgb.processBusinessConnectionChange(operation, debeziumMsg.Payload.Before, debeziumMsg.Payload.After)
	case "user_skills":
		return kgb.processUserSkillChange(operation, debeziumMsg.Payload.Before, debeziumMsg.Payload.After)
	case "project_skills":
		return kgb.processProjectSkillChange(operation, debeziumMsg.Payload.Before, debeziumMsg.Payload.After)
	default:
		fmt.Printf("Unknown table: %s\n", tableName)
		return nil
	}
}

func (kgb *KnowledgeGraphBuilder) processUserChange(op string, before, after map[string]interface{}) error {
	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	switch op {
	case "c", "u": 
		if after == nil {
			return nil
		}
		query := `
		MERGE (u:User {id: $id})
		SET u.firstName = $firstName,
		    u.lastName = $lastName,
		    u.email = $email,
		    u.active = $active,
		    u.emailVerified = $emailVerified,
		    u.createdAt = $createdAt,
		    u.updatedAt = $updatedAt
		`
		params := map[string]interface{}{
			"id":            after["id"],
			"firstName":     after["first_name"],
			"lastName":      after["last_name"],
			"email":         after["login_email"],
			"active":        after["active"],
			"emailVerified": after["email_verified"],
			"createdAt":     after["created_at"],
			"updatedAt":     after["updated_at"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err

	case "d": 
		if before == nil {
			return nil
		}
		query := `MATCH (u:User {id: $id}) DETACH DELETE u`
		params := map[string]interface{}{"id": before["id"]}
		_, err := session.Run(kgb.ctx, query, params)
		return err
	}

	return nil
}

func (kgb *KnowledgeGraphBuilder) processBusinessChange(op string, before, after map[string]interface{}) error {
	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	switch op {
	case "c", "u":
		if after == nil {
			return nil
		}
		query := `
		MERGE (b:Business {id: $id})
		SET b.name = $name,
		    b.tagline = $tagline,
		    b.website = $website,
		    b.businessType = $businessType,
		    b.businessCategory = $businessCategory,
		    b.businessPhase = $businessPhase,
		    b.active = $active,
		    b.createdAt = $createdAt,
		    b.updatedAt = $updatedAt
		WITH b
		MATCH (u:User {id: $operatorUserId})
		MERGE (u)-[:OPERATES]->(b)
		`
		params := map[string]interface{}{
			"id":               after["id"],
			"name":             after["name"],
			"tagline":          after["tagline"],
			"website":          after["website"],
			"businessType":     after["business_type"],
			"businessCategory": after["business_category"],
			"businessPhase":    after["business_phase"],
			"active":           after["active"],
			"createdAt":        after["created_at"],
			"updatedAt":        after["updated_at"],
			"operatorUserId":   after["operator_user_id"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err

	case "d":
		if before == nil {
			return nil
		}
		query := `MATCH (b:Business {id: $id}) DETACH DELETE b`
		params := map[string]interface{}{"id": before["id"]}
		_, err := session.Run(kgb.ctx, query, params)
		return err
	}

	return nil
}

func (kgb *KnowledgeGraphBuilder) processProjectChange(op string, before, after map[string]interface{}) error {
	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	switch op {
	case "c", "u":
		if after == nil {
			return nil
		}
		query := `
		MERGE (p:Project {id: $id})
		SET p.name = $name,
		    p.description = $description,
		    p.projectStatus = $projectStatus,
		    p.startDate = $startDate,
		    p.targetEndDate = $targetEndDate,
		    p.createdAt = $createdAt,
		    p.updatedAt = $updatedAt
		WITH p
		MATCH (u:User {id: $managedByUserId})
		MERGE (u)-[:MANAGES]->(p)
		`
		
		if businessID, ok := after["business_id"]; ok && businessID != nil {
			query += `
			WITH p
			MATCH (b:Business {id: $businessId})
			MERGE (b)-[:HAS_PROJECT]->(p)
			`
		}

		params := map[string]interface{}{
			"id":              after["id"],
			"name":            after["name"],
			"description":     after["description"],
			"projectStatus":   after["project_status"],
			"startDate":       after["start_date"],
			"targetEndDate":   after["target_end_date"],
			"createdAt":       after["created_at"],
			"updatedAt":       after["updated_at"],
			"managedByUserId": after["managed_by_user_id"],
		}

		if businessID, ok := after["business_id"]; ok && businessID != nil {
			params["businessId"] = businessID
		}

		_, err := session.Run(kgb.ctx, query, params)
		return err

	case "d":
		if before == nil {
			return nil
		}
		query := `MATCH (p:Project {id: $id}) DETACH DELETE p`
		params := map[string]interface{}{"id": before["id"]}
		_, err := session.Run(kgb.ctx, query, params)
		return err
	}

	return nil
}

func (kgb *KnowledgeGraphBuilder) processSkillChange(op string, before, after map[string]interface{}) error {
	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	switch op {
	case "c", "u":
		if after == nil {
			return nil
		}
		query := `
		MERGE (s:Skill {id: $id})
		SET s.category = $category,
		    s.name = $name,
		    s.description = $description,
		    s.active = $active,
		    s.createdAt = $createdAt
		`
		params := map[string]interface{}{
			"id":          after["id"],
			"category":    after["category"],
			"name":        after["name"],
			"description": after["description"],
			"active":      after["active"],
			"createdAt":   after["created_at"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err

	case "d":
		if before == nil {
			return nil
		}
		query := `MATCH (s:Skill {id: $id}) DETACH DELETE s`
		params := map[string]interface{}{"id": before["id"]}
		_, err := session.Run(kgb.ctx, query, params)
		return err
	}

	return nil
}

func (kgb *KnowledgeGraphBuilder) processBusinessConnectionChange(op string, before, after map[string]interface{}) error {
	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	switch op {
	case "c", "u":
		if after == nil {
			return nil
		}
		query := `
		MATCH (b1:Business {id: $initiatingBusinessId})
		MATCH (b2:Business {id: $receivingBusinessId})
		MERGE (b1)-[r:CONNECTS_TO {
			connectionType: $connectionType,
			status: $status,
			initiatedByUserId: $initiatedByUserId,
			notes: $notes,
			createdAt: $createdAt,
			updatedAt: $updatedAt
		}]->(b2)
		`
		params := map[string]interface{}{
			"initiatingBusinessId": after["initiating_business_id"],
			"receivingBusinessId":  after["receiving_business_id"],
			"connectionType":       after["connection_type"],
			"status":               after["status"],
			"initiatedByUserId":    after["initiated_by_user_id"],
			"notes":                after["notes"],
			"createdAt":            after["created_at"],
			"updatedAt":            after["updated_at"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err

	case "d":
		if before == nil {
			return nil
		}
		query := `
		MATCH (b1:Business {id: $initiatingBusinessId})-[r:CONNECTS_TO]->(b2:Business {id: $receivingBusinessId})
		WHERE r.connectionType = $connectionType
		DELETE r
		`
		params := map[string]interface{}{
			"initiatingBusinessId": before["initiating_business_id"],
			"receivingBusinessId":  before["receiving_business_id"],
			"connectionType":       before["connection_type"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err
	}

	return nil
}

func (kgb *KnowledgeGraphBuilder) processUserSkillChange(op string, before, after map[string]interface{}) error {
	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	switch op {
	case "c", "u":
		if after == nil {
			return nil
		}
		query := `
		MATCH (u:User {id: $userId})
		MATCH (s:Skill {id: $skillId})
		MERGE (u)-[r:HAS_SKILL {
			proficiencyLevel: $proficiencyLevel,
			createdAt: $createdAt
		}]->(s)
		`
		params := map[string]interface{}{
			"userId":           after["user_id"],
			"skillId":          after["skill_id"],
			"proficiencyLevel": after["proficiency_level"],
			"createdAt":        after["created_at"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err

	case "d":
		if before == nil {
			return nil
		}
		query := `
		MATCH (u:User {id: $userId})-[r:HAS_SKILL]->(s:Skill {id: $skillId})
		DELETE r
		`
		params := map[string]interface{}{
			"userId":  before["user_id"],
			"skillId": before["skill_id"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err
	}

	return nil
}

func (kgb *KnowledgeGraphBuilder) processProjectSkillChange(op string, before, after map[string]interface{}) error {
	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	switch op {
	case "c", "u":
		if after == nil {
			return nil
		}
		query := `
		MATCH (p:Project {id: $projectId})
		MATCH (s:Skill {id: $skillId})
		MERGE (p)-[r:REQUIRES_SKILL {
			importance: $importance
		}]->(s)
		`
		params := map[string]interface{}{
			"projectId":  after["project_id"],
			"skillId":    after["skill_id"],
			"importance": after["importance"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err

	case "d":
		if before == nil {
			return nil
		}
		query := `
		MATCH (p:Project {id: $projectId})-[r:REQUIRES_SKILL]->(s:Skill {id: $skillId})
		DELETE r
		`
		params := map[string]interface{}{
			"projectId": before["project_id"],
			"skillId":   before["skill_id"],
		}
		_, err := session.Run(kgb.ctx, query, params)
		return err
	}

	return nil
}

func (kgb *KnowledgeGraphBuilder) loadExistingData() error {
	dbURL := getEnv("DATABASE_URL", "root:tia-dev-password@tcp(database:3306)/tia-dev?parseTime=true")
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL: %v", err)
	}

	session := kgb.driver.NewSession(kgb.ctx, neo4j.SessionConfig{})
	defer session.Close(kgb.ctx)

	if err := kgb.loadUsers(db, session); err != nil {
		log.Printf("Warning: Failed to load users: %v", err)
	}

	if err := kgb.loadBusinesses(db, session); err != nil {
		log.Printf("Warning: Failed to load businesses: %v", err)
	}

	if err := kgb.loadProjects(db, session); err != nil {
		log.Printf("Warning: Failed to load projects: %v", err)
	}

	if err := kgb.loadSkills(db, session); err != nil {
		log.Printf("Warning: Failed to load skills: %v", err)
	}

	if err := kgb.loadBusinessConnections(db, session); err != nil {
		log.Printf("Warning: Failed to load business connections: %v", err)
	}

	if err := kgb.loadUserSkills(db, session); err != nil {
		log.Printf("Warning: Failed to load user skills: %v", err)
	}

	if err := kgb.loadProjectSkills(db, session); err != nil {
		log.Printf("Warning: Failed to load project skills: %v", err)
	}

	return nil
}

func (kgb *KnowledgeGraphBuilder) loadUsers(db *sql.DB, session neo4j.SessionWithContext) error {
	rows, err := db.Query("SELECT id, first_name, last_name, login_email, created_at, updated_at FROM users")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		var lastName *string
		err := rows.Scan(&user.ID, &user.FirstName, &lastName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			continue
		}
		if lastName != nil {
			user.LastName = *lastName
		}

		query := `
			MERGE (u:User {id: $id})
			SET u.firstName = $firstName,
				u.lastName = $lastName,
				u.email = $email,
				u.createdAt = $createdAt,
				u.updatedAt = $updatedAt
		`
		params := map[string]interface{}{
			"id":        user.ID,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		}
		_, err = session.Run(kgb.ctx, query, params)
		if err != nil {
			log.Printf("Warning: Failed to create user node: %v", err)
		}
	}
	return nil
}

func (kgb *KnowledgeGraphBuilder) loadBusinesses(db *sql.DB, session neo4j.SessionWithContext) error {
	rows, err := db.Query("SELECT id, operator_user_id, name, business_type, business_category, business_phase, description, website, created_at, updated_at FROM businesses")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var business Business
		var operatorUserID uint
		err := rows.Scan(&business.ID, &operatorUserID, &business.Name, &business.BusinessType, &business.BusinessCategory, &business.BusinessPhase, &business.Description, &business.Website, &business.CreatedAt, &business.UpdatedAt)
		if err != nil {
			continue
		}

		query := `
			MERGE (b:Business {id: $id})
			SET b.name = $name,
				b.businessType = $businessType,
				b.businessCategory = $businessCategory,
				b.businessPhase = $businessPhase,
				b.description = $description,
				b.website = $website,
				b.createdAt = $createdAt,
				b.updatedAt = $updatedAt
			WITH b
			MATCH (u:User {id: $operatorUserId})
			MERGE (u)-[:OPERATES]->(b)
		`
		params := map[string]interface{}{
			"id":               business.ID,
			"name":             business.Name,
			"businessType":     business.BusinessType,
			"businessCategory": business.BusinessCategory,
			"businessPhase":    business.BusinessPhase,
			"description":      business.Description,
			"website":          business.Website,
			"createdAt":        business.CreatedAt,
			"updatedAt":        business.UpdatedAt,
			"operatorUserId":   operatorUserID,
		}
		_, err = session.Run(kgb.ctx, query, params)
		if err != nil {
			log.Printf("Warning: Failed to create business node and relationship: %v", err)
		}
	}
	return nil
}

func (kgb *KnowledgeGraphBuilder) loadProjects(db *sql.DB, session neo4j.SessionWithContext) error {
	rows, err := db.Query("SELECT id, managed_by_user_id, name, description, project_status, created_at, updated_at FROM projects")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var project Project
		var managedByUserID uint
		var description *string
		err := rows.Scan(&project.ID, &managedByUserID, &project.Title, &description, &project.Status, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			continue
		}
		if description != nil {
			project.Description = *description
		}

		query := `
			MERGE (p:Project {id: $id})
			SET p.name = $name,
				p.description = $description,
				p.status = $status,
				p.createdAt = $createdAt,
				p.updatedAt = $updatedAt
			WITH p
			MATCH (u:User {id: $managedByUserId})
			MERGE (u)-[:MANAGES]->(p)
		`
		params := map[string]interface{}{
			"id":              project.ID,
			"name":            project.Title,
			"description":     project.Description,
			"status":          project.Status,
			"createdAt":       project.CreatedAt,
			"updatedAt":       project.UpdatedAt,
			"managedByUserId": managedByUserID,
		}
		_, err = session.Run(kgb.ctx, query, params)
		if err != nil {
			log.Printf("Warning: Failed to create project node and relationship: %v", err)
		}
	}
	return nil
}

func (kgb *KnowledgeGraphBuilder) loadSkills(db *sql.DB, session neo4j.SessionWithContext) error {
	rows, err := db.Query("SELECT id, name, category, created_at FROM skills")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var skill Skill
		var description *string
		err := rows.Scan(&skill.ID, &skill.Name, &skill.Category, &skill.CreatedAt)
		if err != nil {
			continue
		}
		if description != nil {
			skill.UpdatedAt = skill.CreatedAt
		} else {
			skill.UpdatedAt = skill.CreatedAt
		}

		query := `
			MERGE (s:Skill {id: $id})
			SET s.name = $name,
				s.category = $category,
				s.createdAt = $createdAt,
				s.updatedAt = $updatedAt
		`
		params := map[string]interface{}{
			"id":        skill.ID,
			"name":      skill.Name,
			"category":  skill.Category,
			"createdAt": skill.CreatedAt,
			"updatedAt": skill.UpdatedAt,
		}
		_, err = session.Run(kgb.ctx, query, params)
		if err != nil {
			log.Printf("Warning: Failed to create skill node: %v", err)
		}
	}
	return nil
}

func (kgb *KnowledgeGraphBuilder) loadBusinessConnections(db *sql.DB, session neo4j.SessionWithContext) error {
	rows, err := db.Query("SELECT id, initiating_business_id, receiving_business_id, connection_type, status, created_at, updated_at FROM business_connections")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var connection BusinessConnection
		err := rows.Scan(&connection.ID, &connection.Business1ID, &connection.Business2ID, &connection.ConnectionType, &connection.Status, &connection.CreatedAt, &connection.UpdatedAt)
		if err != nil {
			continue
		}

		query := `
			MATCH (b1:Business {id: $initiatingBusinessId}), (b2:Business {id: $receivingBusinessId})
			MERGE (b1)-[r:CONNECTED_TO]->(b2)
			SET r.connectionType = $connectionType,
				r.status = $status,
				r.createdAt = $createdAt,
				r.updatedAt = $updatedAt
		`
		params := map[string]interface{}{
			"initiatingBusinessId": connection.Business1ID,
			"receivingBusinessId":  connection.Business2ID,
			"connectionType":       connection.ConnectionType,
			"status":               connection.Status,
			"createdAt":            connection.CreatedAt,
			"updatedAt":            connection.UpdatedAt,
		}
		_, err = session.Run(kgb.ctx, query, params)
		if err != nil {
			log.Printf("Warning: Failed to create business connection: %v", err)
		}
	}
	return nil
}

func (kgb *KnowledgeGraphBuilder) loadUserSkills(db *sql.DB, session neo4j.SessionWithContext) error {
	rows, err := db.Query("SELECT user_id, skill_id, proficiency_level, created_at FROM user_skills")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userSkill UserSkill
		err := rows.Scan(&userSkill.UserID, &userSkill.SkillID, &userSkill.ProficiencyLevel, &userSkill.CreatedAt)
		if err != nil {
			continue
		}
		userSkill.UpdatedAt = userSkill.CreatedAt

		query := `
			MATCH (u:User {id: $userId}), (s:Skill {id: $skillId})
			MERGE (u)-[r:HAS_SKILL]->(s)
			SET r.proficiencyLevel = $proficiencyLevel,
				r.createdAt = $createdAt,
				r.updatedAt = $updatedAt
		`
		params := map[string]interface{}{
			"userId":           userSkill.UserID,
			"skillId":          userSkill.SkillID,
			"proficiencyLevel": userSkill.ProficiencyLevel,
			"createdAt":        userSkill.CreatedAt,
			"updatedAt":        userSkill.UpdatedAt,
		}
		_, err = session.Run(kgb.ctx, query, params)
		if err != nil {
			log.Printf("Warning: Failed to create user skill relationship: %v", err)
		}
	}
	return nil
}

func (kgb *KnowledgeGraphBuilder) loadProjectSkills(db *sql.DB, session neo4j.SessionWithContext) error {
	rows, err := db.Query("SELECT project_id, skill_id, importance FROM project_skills")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var projectSkill ProjectSkill
		err := rows.Scan(&projectSkill.ProjectID, &projectSkill.SkillID, &projectSkill.ImportanceLevel)
		if err != nil {
			continue
		}
		projectSkill.CreatedAt = time.Now()
		projectSkill.UpdatedAt = time.Now()

		query := `
			MATCH (p:Project {id: $projectId}), (s:Skill {id: $skillId})
			MERGE (p)-[r:REQUIRES_SKILL]->(s)
			SET r.importanceLevel = $importanceLevel,
				r.createdAt = $createdAt,
				r.updatedAt = $updatedAt
		`
		params := map[string]interface{}{
			"projectId":       projectSkill.ProjectID,
			"skillId":         projectSkill.SkillID,
			"importanceLevel": projectSkill.ImportanceLevel,
			"createdAt":       projectSkill.CreatedAt,
			"updatedAt":       projectSkill.UpdatedAt,
		}
		_, err = session.Run(kgb.ctx, query, params)
		if err != nil {
			log.Printf("Warning: Failed to create project skill relationship: %v", err)
		}
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
