package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DailyActivityHandler struct {
	service  *services.DailyActivityService
	validate *validator.Validate
}

func NewDailyActivityHandler(service *services.DailyActivityService) *DailyActivityHandler {
	return &DailyActivityHandler{
		service:  service,
		validate: validator.New(),
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

	c.JSON(http.StatusOK, activities)
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

// @Summary      Enrol in Daily Activity
// @Description  Enrols the authenticated user in a daily activity.
// @Tags         DailyActivities
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Daily Activity ID"
// @Success      201 {object} models.DailyActivityEnrolment
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Activity or user not found"
// @Failure      409 {object} map[string]string "User already enrolled"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /daily-activities/{id}/enrol [post]
func (h *DailyActivityHandler) EnrolUserInActivity(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, _ := userIDVal.(uint)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	activityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	input := ports.EnrolInActivityInput{
		DailyActivityID: uint(activityID),
		UserID:          userID,
	}

	enrolment, err := h.service.EnrolUserInActivity(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ports.ErrAlreadyEnrolled) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ports.ErrProjectOrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity or user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enrol in activity"})
		return
	}

	c.JSON(http.StatusCreated, enrolment)
}

// @Summary      Withdraw from Daily Activity
// @Description  Withdraws the authenticated user from a daily activity.
// @Tags         DailyActivities
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Daily Activity ID"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Enrolment not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /daily-activities/{id}/enrol [delete]
func (h *DailyActivityHandler) WithdrawUserFromActivity(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, _ := userIDVal.(uint)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	activityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	err = h.service.WithdrawUserFromActivity(c.Request.Context(), uint(activityID), userID)
	if err != nil {
		if errors.Is(err, ports.ErrEnrolmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to withdraw from activity"})
		return
	}

	c.Status(http.StatusNoContent)
}
