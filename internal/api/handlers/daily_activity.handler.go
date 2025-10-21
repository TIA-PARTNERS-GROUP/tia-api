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
	routes   *constants.Routes // <-- ADDED
}

// Updated constructor
func NewDailyActivityHandler(service *services.DailyActivityService, routes *constants.Routes) *DailyActivityHandler {
	return &DailyActivityHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes, // <-- ADDED
	}
}

// @Summary      Create Daily Activity
// @Description  Creates a new daily activity.
// @Tags         DailyActivities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        activity body ports.CreateDailyActivityInput true "Daily Activity Details"
// @Success      201 {object} models.DailyActivity
// @Failure      400 {object} map[string]string "Validation error or invalid request body"
// @Failure      409 {object} map[string]string "Activity name already exists"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /daily-activities [post]
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

	// Assuming only admins/authorized users can create. Get user ID for service layer check.
	// --- USE CONSTANT ---
	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserID, _ := authUserIDVal.(uint)
	_ = authUserID // Placeholder - Pass to service if needed for permission checks

	activity, err := h.service.CreateDailyActivity(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ports.ErrActivityNameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create daily activity"})
		return
	}

	c.JSON(http.StatusCreated, activity) // Returning the model directly seems fine here
}

// @Summary      Get All Daily Activities
// @Description  Retrieves a list of all daily activities.
// @Tags         DailyActivities
// @Produce      json
// @Success      200 {array} models.DailyActivity
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /daily-activities [get]
func (h *DailyActivityHandler) GetAllDailyActivities(c *gin.Context) {
	activities, err := h.service.GetAllDailyActivities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve daily activities"})
		return
	}

	c.JSON(http.StatusOK, activities) // Returning models directly
}

// @Summary      Get Daily Activity By ID
// @Description  Retrieves a single daily activity by its ID.
// @Tags         DailyActivities
// @Produce      json
// @Param        id path int true "Daily Activity ID"
// @Success      200 {object} models.DailyActivity
// @Failure      404 {object} map[string]string "Activity not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /daily-activities/{id} [get]
func (h *DailyActivityHandler) GetDailyActivityByID(c *gin.Context) {
	// --- USE CONSTANT ---
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

	c.JSON(http.StatusOK, activity) // Returning model directly
}

// --- REMOVED EnrolUserInActivity and WithdrawUserFromActivity ---
// These functions were duplicated from daily_activity_enrolment.handler.go
// and were not used by the routes defined for DailyActivityHandler.
