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
type L2EHandler struct {
	service  *services.L2EResponseService
	validate *validator.Validate
	routes   *constants.Routes 
}
func NewL2EHandler(service *services.L2EResponseService, routes *constants.Routes) *L2EHandler {
	return &L2EHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes, 
	}
}
func (h *L2EHandler) CreateL2EResponse(c *gin.Context) {
	var input ports.CreateL2EResponseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	
	userIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in context"})
		return
	}
	input.UserID = userID 
	
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.CreateL2EResponse(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create L2E response"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapL2EResponseToResponse(response))
}
func (h *L2EHandler) GetL2EResponsesForUser(c *gin.Context) {
	
	userIDStr := c.Param(h.routes.ParamKeyID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	
	
	
	
	
	
	responses, err := h.service.GetL2EResponsesForUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve L2E responses"})
		return
	}
	responseDTOs := make([]ports.L2EResponseResponse, len(responses))
	for i, resp := range responses {
		responseDTOs[i] = ports.MapL2EResponseToResponse(&resp)
	}
	c.JSON(http.StatusOK, responseDTOs)
}
