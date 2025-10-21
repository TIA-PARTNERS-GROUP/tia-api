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
