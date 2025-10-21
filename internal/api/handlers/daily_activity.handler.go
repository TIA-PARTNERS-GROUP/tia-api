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

// @Summary Create Daily Activity
// @Description Creates a new daily activity definition (e.g., "30-minute meditation").
// @Tags daily_activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param activity body ports.CreateDailyActivityInput true "Activity details (Name, Description)"
// @Success 201 {object} models.DailyActivity "Activity created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "ErrActivityNameExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /daily-activities [post]
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

// @Summary Get All Daily Activities
// @Description Retrieves a list of all defined daily activities.
// @Tags daily_activities
// @Produce json
// @Success 200 {array} models.DailyActivity "List of daily activities"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /daily-activities [get]
func (h *DailyActivityHandler) GetAllDailyActivities(c *gin.Context) {
	activities, err := h.service.GetAllDailyActivities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve daily activities"})
		return
	}
	c.JSON(http.StatusOK, activities)
}

// @Summary Get Daily Activity by ID
// @Description Retrieves a specific daily activity definition by its unique ID.
// @Tags daily_activities
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} models.DailyActivity "Activity retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid activity ID"
// @Failure 404 {object} map[string]interface{} "ErrDailyActivityNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /daily-activities/{id} [get]
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
