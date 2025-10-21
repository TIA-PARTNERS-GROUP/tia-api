package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"     // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services" // <-- IMPORT models for response type
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10" // <-- IMPORT validator (although not used directly here, good practice)
)

type DailyActivityEnrolmentHandler struct {
	service  *services.DailyActivityEnrolmentService
	validate *validator.Validate // <-- ADDED (though not used in this specific file)
	routes   *constants.Routes   // <-- ADDED
}

// Updated constructor
func NewDailyActivityEnrolmentHandler(service *services.DailyActivityEnrolmentService, routes *constants.Routes) *DailyActivityEnrolmentHandler {
	return &DailyActivityEnrolmentHandler{
		service:  service,
		validate: validator.New(), // <-- Initialize validator
		routes:   routes,          // <-- Store routes
	}
}

// @Summary      Enrol in a Daily Activity
// @Description  Enrols the currently authenticated user in a specific daily activity.
// @Tags         DailyActivityEnrolments
// @Security     BearerAuth
// @Produce      json
// @Param        activityID path int true "Daily Activity ID (Note: Param name in route is ':id')"
// @Success      201 {object} models.DailyActivityEnrolment
// @Failure      400 {object} map[string]string "Invalid activity ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Activity or user not found"
// @Failure      409 {object} map[string]string "User already enrolled in this activity"
// @Router       /daily-activities/{id}/enrolments [post] // Corrected router path in comment
func (h *DailyActivityEnrolmentHandler) EnrolUser(c *gin.Context) {
	// --- USE CONSTANT ---
	userIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	// --- USE CONSTANT ---
	activityIDStr := c.Param(h.routes.ParamKeyID) // Use ParamKeyID as defined in the route
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
		// Distinguish between ErrAlreadyEnrolled and others if needed based on service error types
		if errors.Is(err, ports.ErrAlreadyEnrolled) { // Assuming service returns this specific error
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ports.ErrProjectOrUserNotFound) || errors.Is(err, ports.ErrDailyActivityNotFound) { // Adapt as per service errors
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity or user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enrol in activity"})
		return
	}

	// Consider mapping to a response DTO if needed
	c.JSON(http.StatusCreated, ports.MapToActivityEnrolmentResponse(enrolment)) // Map to DTO
}

// @Summary      Withdraw from a Daily Activity
// @Description  Withdraws the currently authenticated user from a specific daily activity.
// @Tags         DailyActivityEnrolments
// @Security     BearerAuth
// @Produce      json
// @Param        activityID path int true "Daily Activity ID (Note: Param name in route is ':id')"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Invalid activity ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Enrolment not found for this user and activity"
// @Router       /daily-activities/{id}/enrolments [delete] // Corrected router path in comment
func (h *DailyActivityEnrolmentHandler) WithdrawUser(c *gin.Context) {
	// --- USE CONSTANT ---
	userIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	// --- USE CONSTANT ---
	activityIDStr := c.Param(h.routes.ParamKeyID) // Use ParamKeyID
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
		if errors.Is(err, ports.ErrEnrolmentNotFound) { // Assuming service returns this
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
// @Param        activityID path int true "Daily Activity ID (Note: Param name in route is ':id')"
// @Success      200 {array} ports.ActivityEnrolmentResponse // Updated response type
// @Failure      400 {object} map[string]string "Invalid activity ID"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /daily-activities/{id}/enrolments [get] // Corrected router path in comment
func (h *DailyActivityEnrolmentHandler) GetEnrolmentsForActivity(c *gin.Context) {
	// --- USE CONSTANT ---
	activityIDStr := c.Param(h.routes.ParamKeyID) // Use ParamKeyID
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

	// Map to response DTOs
	response := make([]ports.ActivityEnrolmentResponse, len(enrolments))
	for i, e := range enrolments {
		response[i] = ports.MapToActivityEnrolmentResponse(&e)
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Get Enrolments for a User
// @Description  Retrieves a list of all daily activities a specific user is enrolled in.
// @Tags         DailyActivityEnrolments
// @Produce      json
// @Param        userID path int true "User ID (Note: Param name in route is ':id')"
// @Success      200 {array} ports.UserEnrolmentResponse // Updated response type
// @Failure      400 {object} map[string]string "Invalid user ID"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{id}/enrolments [get] // Corrected router path in comment
func (h *DailyActivityEnrolmentHandler) GetEnrolmentsForUser(c *gin.Context) {
	// --- USE CONSTANT ---
	userIDStr := c.Param(h.routes.ParamKeyID) // Use ParamKeyID
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Potential Authorization Check: Ensure the requesting user is the target user or an admin
	// authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	// authUserID, _ := authUserIDVal.(uint)
	// if authUserID != uint(userID) { /* && !isAdmin(authUserID) */
	//     c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	// 	   return
	// }

	enrolments, err := h.service.GetEnrolmentsForUser(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve enrolments"})
		return
	}

	// Map to response DTOs
	response := make([]ports.UserEnrolmentResponse, len(enrolments))
	for i, e := range enrolments {
		response[i] = ports.MapToUserEnrolmentResponse(&e)
	}

	c.JSON(http.StatusOK, response)
}
