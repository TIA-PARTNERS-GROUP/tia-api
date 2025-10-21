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

type FeedbackHandler struct {
	service  *services.FeedbackService
	validate *validator.Validate
	routes   *constants.Routes
}

func NewFeedbackHandler(service *services.FeedbackService, routes *constants.Routes) *FeedbackHandler {
	return &FeedbackHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes,
	}
}

// @Summary Submit New Feedback
// @Description Allows any client to submit new feedback (e.g., bug report, suggestion). Does not require authentication.
// @Tags feedback
// @Accept json
// @Produce json
// @Param feedback body ports.CreateFeedbackInput true "Feedback details (Name, Email, Content)"
// @Success 201 {object} ports.FeedbackResponse "Feedback submitted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 500 {object} map[string]interface{} "Internal server er
func (h *FeedbackHandler) CreateFeedback(c *gin.Context) {
	var input ports.CreateFeedbackInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	feedback, err := h.service.CreateFeedback(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feedback"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapFeedbackToResponse(feedback))
}

// @Summary Get All Feedback
// @Description Retrieves a list of all submitted feedback. Requires authentication.
// @Tags feedback
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ports.FeedbackResponse "List of feedback entries"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /feedback [get]
func (h *FeedbackHandler) GetAllFeedback(c *gin.Context) {

	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	feedbacks, err := h.service.GetAllFeedback(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
		return
	}

	response := make([]ports.FeedbackResponse, len(feedbacks))
	for i, fb := range feedbacks {
		response[i] = ports.MapFeedbackToResponse(&fb)
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Get Feedback by ID
// @Description Retrieves a specific feedback entry by its unique ID. Requires authentication.
// @Tags feedback
// @Produce json
// @Security BearerAuth
// @Param id path int true "Feedback ID"
// @Success 200 {object} ports.FeedbackResponse "Feedback entry retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid feedback ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrFeedbackNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /feedback/{id} [get]
func (h *FeedbackHandler) GetFeedbackByID(c *gin.Context) {

	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}
	feedback, err := h.service.GetFeedbackByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, ports.ErrFeedbackNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
		return
	}
	c.JSON(http.StatusOK, ports.MapFeedbackToResponse(feedback))
}

// @Summary Delete Feedback
// @Description Deletes a specific feedback entry by its unique ID. Requires authentication.
// @Tags feedback
// @Produce json
// @Security BearerAuth
// @Param id path int true "Feedback ID"
// @Success 204 "Feedback deleted successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid feedback ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrFeedbackNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /feedback/{id} [delete]
func (h *FeedbackHandler) DeleteFeedback(c *gin.Context) {

	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}
	err = h.service.DeleteFeedback(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, ports.ErrFeedbackNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete feedback"})
		return
	}
	c.Status(http.StatusNoContent)
}
