package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BusinessConnectionHandler struct {
	service  *services.BusinessConnectionService
	validate *validator.Validate
}

func NewBusinessConnectionHandler(s *services.BusinessConnectionService) *BusinessConnectionHandler {
	return &BusinessConnectionHandler{
		service:  s,
		validate: validator.New(),
	}
}

// @Summary      Initiate a business connection
// @Description  Creates a new connection request between two businesses. Requires authentication.
// @Tags         Business Connections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        connection body ports.CreateBusinessConnectionInput true "Connection Initiation Data"
// @Success      201 {object} ports.BusinessConnectionResponse
// @Failure      400 {object} map[string]string "Validation error, bad request (e.g., connecting to self)"
// @Failure      404 {object} map[string]string "Initiating user or business not found"
// @Failure      409 {object} map[string]string "Connection already exists"
// @Router       /connections [post]
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

	userID, _ := c.Get("userID")
	input.InitiatedByUserID = userID.(uint)

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

// @Summary      Get a business connection by ID
// @Description  Retrieves details of a single business connection. Requires authentication.
// @Tags         Business Connections
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Connection ID"
// @Success      200 {object} ports.BusinessConnectionResponse
// @Failure      400 {object} map[string]string "Invalid ID format"
// @Failure      404 {object} map[string]string "Connection not found"
// @Router       /connections/{id} [get]
func (h *BusinessConnectionHandler) GetBusinessConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
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

// @Summary      Get connections for a specific business
// @Description  Retrieves all connections (incoming and outgoing) for a specific business ID. Requires authentication.
// @Tags         Business Connections
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Business ID"
// @Param        status query string false "Filter by status (pending, active, rejected, inactive)" Enums(pending, active, rejected, inactive)
// @Param        type query string false "Filter by connection type (Partnership, Supplier, Client, Referral, Collaboration)" Enums(Partnership, Supplier, Client, Referral, Collaboration)
// @Success      200 {object} ports.BusinessConnectionsResponse
// @Failure      400 {object} map[string]string "Invalid Business ID format"
// @Router       /businesses/{id}/connections [get]
func (h *BusinessConnectionHandler) GetBusinessConnections(c *gin.Context) {
	businessID, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

	connections, err := h.service.GetBusinessConnections(c.Request.Context(), uint(businessID), connType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve connections"})
		return
	}

	c.JSON(http.StatusOK, ports.MapToBusinessConnectionsResponse(connections))
}

// @Summary      Update a business connection
// @Description  Updates the status or notes of a business connection. Requires authentication.
// @Tags         Business Connections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Connection ID"
// @Param        update body ports.UpdateBusinessConnectionInput true "Fields to update"
// @Success      200 {object} ports.BusinessConnectionResponse
// @Failure      400 {object} map[string]string "Invalid ID or request body"
// @Failure      404 {object} map[string]string "Connection not found"
// @Router       /connections/{id} [put]
func (h *BusinessConnectionHandler) UpdateBusinessConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

// @Summary      Accept a business connection request
// @Description  Changes the status of a pending connection to 'active'. Requires authentication.
// @Tags         Business Connections
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Connection ID"
// @Success      200 {object} ports.BusinessConnectionResponse
// @Failure      400 {object} map[string]string "Invalid ID format"
// @Failure      404 {object} map[string]string "Connection not found or not pending"
// @Router       /connections/{id}/accept [post]
func (h *BusinessConnectionHandler) AcceptBusinessConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
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

// @Summary      Reject a business connection request
// @Description  Changes the status of a pending connection to 'rejected'. Requires authentication.
// @Tags         Business Connections
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Connection ID"
// @Success      200 {object} ports.BusinessConnectionResponse
// @Failure      400 {object} map[string]string "Invalid ID format"
// @Failure      404 {object} map[string]string "Connection not found or not pending"
// @Router       /connections/{id}/reject [post]
func (h *BusinessConnectionHandler) RejectBusinessConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
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

// @Summary      Delete a business connection
// @Description  Deletes a business connection by its ID. Requires authentication.
// @Tags         Business Connections
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Connection ID"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Invalid ID format"
// @Failure      404 {object} map[string]string "Connection not found"
// @Router       /connections/{id} [delete]
func (h *BusinessConnectionHandler) DeleteBusinessConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
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
