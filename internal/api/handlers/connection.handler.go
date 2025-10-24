package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	routes *constants.Routes
}

func NewConnectionHandler(routes *constants.Routes) *ConnectionHandler {
	return &ConnectionHandler{
		routes: routes,
	}
}

// @Summary Get Complementary Partners
// @Description Retrieves complementary business partners for a user based on business type compatibility
// @Tags connections
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]interface{} "Complementary partners retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/complementary/{userId} [get]
func (h *ConnectionHandler) GetComplementaryPartners(c *gin.Context) {
	userIDStr := c.Param(h.routes.ParamKeyUserID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	if authUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	recommendations, err := h.callConnectionAnalyzer(uint(userID), "complementary")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get complementary partners", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// @Summary Get Alliance Partners
// @Description Retrieves alliance partners for a user based on skill compatibility and project collaboration potential
// @Tags connections
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]interface{} "Alliance partners retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/alliance/{userId} [get]
func (h *ConnectionHandler) GetAlliancePartners(c *gin.Context) {
	userIDStr := c.Param(h.routes.ParamKeyUserID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	if authUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	recommendations, err := h.callConnectionAnalyzer(uint(userID), "alliance")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alliance partners", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// @Summary Get Mastermind Partners
// @Description Retrieves mastermind partners for a user based on complementary skills and business phase
// @Tags connections
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]interface{} "Mastermind partners retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/mastermind/{userId} [get]
func (h *ConnectionHandler) GetMastermindPartners(c *gin.Context) {
	userIDStr := c.Param(h.routes.ParamKeyUserID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	if authUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	recommendations, err := h.callConnectionAnalyzer(uint(userID), "mastermind")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get mastermind partners", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// @Summary Get All Connection Recommendations
// @Description Retrieves all types of connection recommendations for a user
// @Tags connections
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]interface{} "All recommendations retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/recommendations/{userId} [get]
func (h *ConnectionHandler) GetAllRecommendations(c *gin.Context) {
	userIDStr := c.Param(h.routes.ParamKeyUserID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	if authUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	recommendations, err := h.callConnectionAnalyzer(uint(userID), "recommendations")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// @Summary Get Connection Analysis
// @Description Retrieves comprehensive connection analysis for a user
// @Tags connections
// @Produce json
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]interface{} "Connection analysis retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/analysis/{userId} [get]
func (h *ConnectionHandler) GetConnectionAnalysis(c *gin.Context) {
	userIDStr := c.Param(h.routes.ParamKeyUserID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	if authUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	analysis, err := h.callConnectionAnalyzer(uint(userID), "analysis")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get connection analysis"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

func (h *ConnectionHandler) callConnectionAnalyzer(userID uint, endpoint string) (map[string]interface{}, error) {
	baseURL := "http://connection-analyzer:8082/api/v1/connections"
	var url string
	
	switch endpoint {
	case "complementary":
		url = baseURL + "/complementary/" + strconv.FormatUint(uint64(userID), 10)
	case "alliance":
		url = baseURL + "/alliance/" + strconv.FormatUint(uint64(userID), 10)
	case "mastermind":
		url = baseURL + "/mastermind/" + strconv.FormatUint(uint64(userID), 10)
	case "recommendations":
		url = baseURL + "/recommendations/" + strconv.FormatUint(uint64(userID), 10)
	case "analysis":
		url = baseURL + "/analysis/" + strconv.FormatUint(uint64(userID), 10)
	default:
		return nil, errors.New("invalid endpoint")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to connection analyzer: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("connection analyzer service error (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result, nil
}
