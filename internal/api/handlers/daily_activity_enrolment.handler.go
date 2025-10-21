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
