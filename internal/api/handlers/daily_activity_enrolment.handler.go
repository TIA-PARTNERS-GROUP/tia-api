package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
)

type DailyActivityEnrolmentHandler struct {
	service *services.DailyActivityEnrolmentService
}

func NewDailyActivityEnrolmentHandler(service *services.DailyActivityEnrolmentService) *DailyActivityEnrolmentHandler {
	return &DailyActivityEnrolmentHandler{service: service}
}

// @Summary      Enrol in a Daily Activity
// @Description  Enrols the currently authenticated user in a specific daily activity.
// @Tags         DailyActivityEnrolments
// @Security     BearerAuth
// @Produce      json
// @Param        activityID path int true "Daily Activity ID"
// @Success      201 {object} models.DailyActivityEnrolment
// @Failure      400 {object} map[string]string "Invalid activity ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Activity or user not found"
// @Failure      409 {object} map[string]string "User already enrolled in this activity"
// @Router       /daily-activities/{activityID}/enrolments/ [post]
func (h *DailyActivityEnrolmentHandler) EnrolUser(c *gin.Context) {
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

	input := ports.EnrolmentInput{
		DailyActivityID: uint(activityID),
		UserID:          userID,
	}

	enrolment, err := h.service.EnrolUser(c.Request.Context(), input)
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

// @Summary      Withdraw from a Daily Activity
// @Description  Withdraws the currently authenticated user from a specific daily activity.
// @Tags         DailyActivityEnrolments
// @Security     BearerAuth
// @Produce      json
// @Param        activityID path int true "Daily Activity ID"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Invalid activity ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Enrolment not found for this user and activity"
// @Router       /daily-activities/{activityID}/enrolments/ [delete]
func (h *DailyActivityEnrolmentHandler) WithdrawUser(c *gin.Context) {
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

	err = h.service.WithdrawUser(c.Request.Context(), uint(activityID), userID)
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

// @Summary      Get Enrolments for an Activity
// @Description  Retrieves a list of all users enrolled in a specific daily activity.
// @Tags         DailyActivityEnrolments
// @Produce      json
// @Param        activityID path int true "Daily Activity ID"
// @Success      200 {array} models.DailyActivityEnrolment
// @Failure      400 {object} map[string]string "Invalid activity ID"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /daily-activities/{activityID}/enrolments/ [get]
func (h *DailyActivityEnrolmentHandler) GetEnrolmentsForActivity(c *gin.Context) {
	activityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	enrolments, err := h.service.GetEnrolmentsForActivity(c.Request.Context(), uint(activityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrolments"})
		return
	}

	c.JSON(http.StatusOK, enrolments)
}

// @Summary      Get Enrolments for a User
// @Description  Retrieves a list of all daily activities a specific user is enrolled in.
// @Tags         DailyActivityEnrolments
// @Produce      json
// @Param        userID path int true "User ID"
// @Success      200 {array} models.DailyActivityEnrolment
// @Failure      400 {object} map[string]string "Invalid user ID"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{userID}/enrolments/ [get]
func (h *DailyActivityEnrolmentHandler) GetEnrolmentsForUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	enrolments, err := h.service.GetEnrolmentsForUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrolments"})
		return
	}

	c.JSON(http.StatusOK, enrolments)
}
