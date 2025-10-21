package handlers

import (
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InferredConnectionHandler struct {
	service  *services.InferredConnectionService
	validate *validator.Validate
	routes   *constants.Routes
}

func NewInferredConnectionHandler(service *services.InferredConnectionService, routes *constants.Routes) *InferredConnectionHandler {
	return &InferredConnectionHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes,
	}
}

// @Summary Create Inferred Connection Record
// @Description Creates a new record for a potential connection inferred by a model. Intended for internal/system use.
// @Tags inferred_connections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param connection body ports.CreateInferredConnectionInput true "Inferred connection details"
// @Success 201 {object} ports.InferredConnectionResponse "Connection record created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /inferred-connections [post]
func (h *InferredConnectionHandler) CreateInferredConnection(c *gin.Context) {

	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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

// @Summary Get Inferred Connections by Source Entity
// @Description Retrieves all potential connections inferred from a specific source entity (e.g., a Project or Business).
// @Tags inferred_connections
// @Produce json
// @Security BearerAuth
// @Param entityType path string true "Type of the source entity (e.g., business, project)"
// @Param entityID path int true "ID of the source entity"
// @Success 200 {array} ports.InferredConnectionResponse "List of inferred connections"
// @Failure 400 {object} map[string]interface{} "Invalid entity ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /inferred-connections/source/{entityType}/{entityID} [get]
func (h *InferredConnectionHandler) GetConnectionsForSource(c *gin.Context) {

	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	entityType := c.Param(h.routes.ParamKeyEntityType)
	entityIDStr := c.Param(h.routes.ParamKeyEntityID)
	entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
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
