package handlers
import (
	"errors"
	"net/http"
	"strconv"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
type DailyActivityHandler struct {
	service  *services.DailyActivityService
	validate *validator.Validate
	routes   *constants.Routes 
}
func NewDailyActivityHandler(service *services.DailyActivityService, routes *constants.Routes) *DailyActivityHandler {
	return &DailyActivityHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes, 
	}
}
func (h *DailyActivityHandler) CreateDailyActivity(c *gin.Context) {
	var input ports.CreateDailyActivityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	
	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserID, _ := authUserIDVal.(uint)
	_ = authUserID 
	activity, err := h.service.CreateDailyActivity(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ports.ErrActivityNameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create daily activity"})
		return
	}
	c.JSON(http.StatusCreated, activity) 
}
func (h *DailyActivityHandler) GetAllDailyActivities(c *gin.Context) {
	activities, err := h.service.GetAllDailyActivities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve daily activities"})
		return
	}
	c.JSON(http.StatusOK, activities) 
}
func (h *DailyActivityHandler) GetDailyActivityByID(c *gin.Context) {
	
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	activity, err := h.service.GetDailyActivityByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, ports.ErrDailyActivityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve daily activity"})
		return
	}
	c.JSON(http.StatusOK, activity) 
}
