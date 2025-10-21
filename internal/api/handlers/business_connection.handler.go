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

// @Summary Initiate Business Connection Request
// @Description Creates a new connection request between two businesses. The initiating user is taken from the auth context.
// @Tags business_connections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param connection body ports.CreateBusinessConnectionInput true "Connection request details (InitiatingBusinessID, ReceivingBusinessID, ConnectionType)"
// @Success 201 {object} ports.BusinessConnectionResponse "Connection request created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body, validation failed, or ErrCannotConnectToSelf"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "ErrBusinessConnectionAlreadyExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections [post]
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

// @Summary Get Business Connection by ID
// @Description Retrieves a specific connection record by its unique ID.
// @Tags business_connections
// @Produce json
// @Security BearerAuth
// @Param id path int true "Connection ID"
// @Success 200 {object} ports.BusinessConnectionResponse "Connection retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid connection ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrBusinessConnectionNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/{id} [get]
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

// @Summary List Connections for a Business
// @Description Retrieves a list of all connections (initiated and received) associated with a specific business ID.
// @Tags business_connections
// @Produce json
// @Security BearerAuth
// @Param id path int true "Business ID"
// @Param status query string false "Filter by connection status (pending, active, rejected, inactive)"
// @Param type query string false "Filter by connection type (Partnership, Client, Supplier, etc.)"
// @Success 200 {object} ports.BusinessConnectionsResponse "List of connections"
// @Failure 400 {object} map[string]interface{} "Invalid business ID or query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /businesses/{id}/connections [get]
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

// @Summary Update Business Connection Details
// @Description Updates modifiable fields of an existing connection (e.g., Notes, Type). This is typically restricted to the initiating user.
// @Tags business_connections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Connection ID"
// @Param connection body ports.UpdateBusinessConnectionInput true "Fields to update (e.g., ConnectionType, Notes)"
// @Success 200 {object} ports.BusinessConnectionResponse "Connection updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or connection ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not authorized to update)"
// @Failure 404 {object} map[string]interface{} "ErrBusinessConnectionNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/{id} [put]
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

// @Summary Accept Pending Connection
// @Description Updates the status of a specific connection request to 'active'. Restricted to the receiving business's operator.
// @Tags business_connections
// @Produce json
// @Security BearerAuth
// @Param id path int true "Connection ID"
// @Success 200 {object} ports.BusinessConnectionResponse "Connection successfully accepted and set to active"
// @Failure 400 {object} map[string]interface{} "Invalid connection ID or ErrConnectionNotPending"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the receiving business's operator)"
// @Failure 404 {object} map[string]interface{} "ErrBusinessConnectionNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/{id}/accept [patch]
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

// @Summary Reject Pending Connection
// @Description Updates the status of a specific connection request to 'rejected'. Restricted to the receiving business's operator.
// @Tags business_connections
// @Produce json
// @Security BearerAuth
// @Param id path int true "Connection ID"
// @Success 200 {object} ports.BusinessConnectionResponse "Connection successfully rejected"
// @Failure 400 {object} map[string]interface{} "Invalid connection ID or ErrConnectionNotPending"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the receiving business's operator)"
// @Failure 404 {object} map[string]interface{} "ErrBusinessConnectionNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/{id}/reject [patch]
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

// @Summary Delete Business Connection
// @Description Deletes a specific connection record. Restricted to the initiating user or the receiving business's operator.
// @Tags business_connections
// @Produce json
// @Security BearerAuth
// @Param id path int true "Connection ID"
// @Success 204 "Connection deleted successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid connection ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not authorized to delete)"
// @Failure 404 {object} map[string]interface{} "ErrBusinessConnectionNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /connections/{id} [delete]
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
