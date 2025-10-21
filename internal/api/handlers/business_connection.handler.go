package handlers
import (
	"errors"
	"net/http"
	"strconv"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
type BusinessConnectionHandler struct {
	service  *services.BusinessConnectionService
	validate *validator.Validate
	routes   *constants.Routes 
}
func NewBusinessConnectionHandler(service *services.BusinessConnectionService, routes *constants.Routes) *BusinessConnectionHandler {
	return &BusinessConnectionHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes, 
	}
}
func (h *BusinessConnectionHandler) CreateBusinessConnection(c *gin.Context) {
	var input ports.CreateBusinessConnectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	userIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	input.InitiatedByUserID = userID 
	connection, err := h.service.CreateBusinessConnection(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create connection"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapToBusinessConnectionResponse(connection))
}
func (h *BusinessConnectionHandler) GetBusinessConnection(c *gin.Context) {
	
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}
	
	
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	
	connection, err := h.service.GetBusinessConnection(c.Request.Context(), uint(id)) 
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve connection"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToBusinessConnectionResponse(connection))
}
func (h *BusinessConnectionHandler) GetBusinessConnections(c *gin.Context) {
	
	businessIDStr := c.Param(h.routes.ParamKeyID)
	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil || businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}
	var status *models.BusinessConnectionStatus
	if s := c.Query("status"); s != "" {
		tempStatus := models.BusinessConnectionStatus(s)
		
		status = &tempStatus
	}
	var connType *models.BusinessConnectionType
	if ct := c.Query("type"); ct != "" {
		tempType := models.BusinessConnectionType(ct)
		
		connType = &tempType
	}
	
	
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	connections, err := h.service.GetBusinessConnections(c.Request.Context(), uint(businessID), connType, status) 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve connections"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToBusinessConnectionsResponse(connections))
}
func (h *BusinessConnectionHandler) UpdateBusinessConnection(c *gin.Context) {
	
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}
	var input ports.UpdateBusinessConnectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	connection, err := h.service.UpdateBusinessConnection(c.Request.Context(), uint(id), input) 
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update connection"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToBusinessConnectionResponse(connection))
}
func (h *BusinessConnectionHandler) AcceptBusinessConnection(c *gin.Context) {
	
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}
	
	
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	connection, err := h.service.AcceptBusinessConnection(c.Request.Context(), uint(id)) 
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept connection"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToBusinessConnectionResponse(connection))
}
func (h *BusinessConnectionHandler) RejectBusinessConnection(c *gin.Context) {
	
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}
	
	
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	connection, err := h.service.RejectBusinessConnection(c.Request.Context(), uint(id)) 
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject connection"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToBusinessConnectionResponse(connection))
}
func (h *BusinessConnectionHandler) DeleteBusinessConnection(c *gin.Context) {
	
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}
	
	
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	if err := h.service.DeleteBusinessConnection(c.Request.Context(), uint(id)); err != nil { 
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete connection"})
		return
	}
	c.Status(http.StatusNoContent)
}
