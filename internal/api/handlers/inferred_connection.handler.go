package handlers

import (
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InferredConnectionHandler struct {
	service  *services.InferredConnectionService
	validate *validator.Validate
}

func NewInferredConnectionHandler(service *services.InferredConnectionService) *InferredConnectionHandler {
	return &InferredConnectionHandler{
		service:  service,
		validate: validator.New(),
	}
}

// @Summary      Create Inferred Connection
// @Description  Creates a new inferred connection between two entities. (Protected)
// @Tags         InferredConnections
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        connection body ports.CreateInferredConnectionInput true "Inferred Connection Details"
// @Success      201 {object} ports.InferredConnectionResponse
// @Failure      400 {object} map[string]string "Validation error or invalid request body"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /inferred-connections [post]
func (h *InferredConnectionHandler) CreateInferredConnection(c *gin.Context) {
	var input ports.CreateInferredConnectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connection, err := h.service.CreateInferredConnection(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inferred connection"})
		return
	}

	c.JSON(http.StatusCreated, ports.MapInferredConnectionToResponse(connection))
}

// @Summary      Get Connections For Source
// @Description  Retrieves all inferred connections originating from a specific source entity. (Protected)
// @Tags         InferredConnections
// @Security     BearerAuth
// @Produce      json
// @Param        entityType path string true "Type of the source entity (e.g., user, business)"
// @Param        entityID path int true "ID of the source entity"
// @Success      200 {array} ports.InferredConnectionResponse
// @Failure      400 {object} map[string]string "Invalid entity ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /inferred-connections/source/{entityType}/{entityID} [get]
func (h *InferredConnectionHandler) GetConnectionsForSource(c *gin.Context) {
	entityType := c.Param("entityType")
	entityID, err := strconv.ParseUint(c.Param("entityID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	connections, err := h.service.GetConnectionsForSource(c.Request.Context(), entityType, uint(entityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inferred connections"})
		return
	}

	response := make([]ports.InferredConnectionResponse, len(connections))
	for i, conn := range connections {
		response[i] = ports.MapInferredConnectionToResponse(&conn)
	}

	c.JSON(http.StatusOK, response)
}
