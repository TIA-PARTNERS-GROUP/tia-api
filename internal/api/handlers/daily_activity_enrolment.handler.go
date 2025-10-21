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

type DailyActivityEnrolmentHandler struct {
	service  *services.DailyActivityEnrolmentService
	validate *validator.Validate
	routes   *constants.Routes
}

func NewDailyActivityEnrolmentHandler(service *services.DailyActivityEnrolmentService, routes *constants.Routes) *DailyActivityEnrolmentHandler {
	return &DailyActivityEnrolmentHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes,
	}
}

// @Summary Enrol User in Daily Activity
// @Description Enrols the authenticated user in a specified daily activity.
// @Tags daily_activities, enrolments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Daily Activity ID"
// @Success 201 {object} ports.ActivityEnrolmentResponse "Enrolment successful"
// @Failure 400 {object} gin.H "Invalid activity ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Activity or user not found"
// @Failure 409 {object} gin.H "ErrAlreadyEnrolled"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /daily-activities/{id}/enrolments [post]
func (h *DailyActivityEnrolmentHandler) EnrolUser(c *gin.Context) {

	userIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	activityIDStr := c.Param(h.routes.ParamKeyID)
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
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
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}

		if errors.Is(err, ports.ErrAlreadyEnrolled) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ports.ErrProjectOrUserNotFound) || errors.Is(err, ports.ErrDailyActivityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity or user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enrol in activity"})
		return
	}

	c.JSON(http.StatusCreated, ports.MapToActivityEnrolmentResponse(enrolment))
}

// @Summary Withdraw User from Daily Activity
// @Description Withdraws the authenticated user from a specified daily activity.
// @Tags daily_activities, enrolments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Daily Activity ID"
// @Success 204 "Withdrawal successful (No Content)"
// @Failure 400 {object} gin.H "Invalid activity ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "ErrEnrolmentNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /daily-activities/{id}/enrolments [delete]
func (h *DailyActivityEnrolmentHandler) WithdrawUser(c *gin.Context) {

	userIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	activityIDStr := c.Param(h.routes.ParamKeyID)
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	err = h.service.WithdrawUser(c.Request.Context(), uint(activityID), userID)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		if errors.Is(err, ports.ErrEnrolmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to withdraw from activity"})
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get All Enrolments for Activity
// @Description Retrieves a list of all users currently enrolled in a specified daily activity.
// @Tags daily_activities, enrolments
// @Produce json
// @Param id path int true "Daily Activity ID"
// @Success 200 {array} ports.ActivityEnrolmentResponse "List of user enrolments"
// @Failure 400 {object} gin.H "Invalid activity ID"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /daily-activities/{id}/enrolments [get]
func (h *DailyActivityEnrolmentHandler) GetEnrolmentsForActivity(c *gin.Context) {

	activityIDStr := c.Param(h.routes.ParamKeyID)
	activityID, err := strconv.ParseUint(activityIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	enrolments, err := h.service.GetEnrolmentsForActivity(c.Request.Context(), uint(activityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrolments"})
		return
	}

	response := make([]ports.ActivityEnrolmentResponse, len(enrolments))
	for i, e := range enrolments {
		response[i] = ports.MapToActivityEnrolmentResponse(&e)
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Get All Enrolments for User
// @Description Retrieves a list of all daily activities a specific user is currently enrolled in.
// @Tags daily_activities, enrolments
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} ports.UserEnrolmentResponse "List of activity enrolments for the user"
// @Failure 400 {object} gin.H "Invalid user ID"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/{id}/enrolments [get]
func (h *DailyActivityEnrolmentHandler) GetEnrolmentsForUser(c *gin.Context) {

	userIDStr := c.Param(h.routes.ParamKeyID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	enrolments, err := h.service.GetEnrolmentsForUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrolments"})
		return
	}

	response := make([]ports.UserEnrolmentResponse, len(enrolments))
	for i, e := range enrolments {
		response[i] = ports.MapToUserEnrolmentResponse(&e)
	}
	c.JSON(http.StatusOK, response)
}
