package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"     // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services" // <-- IMPORT models for response
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FeedbackHandler struct {
	service  *services.FeedbackService
	validate *validator.Validate
	routes   *constants.Routes // <-- ADDED
}

// Updated constructor
func NewFeedbackHandler(service *services.FeedbackService, routes *constants.Routes) *FeedbackHandler {
	return &FeedbackHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes, // <-- ADDED
	}
}

// @Summary      Create Feedback
// @Description  Submits a new piece of feedback. This is a public endpoint.
// @Tags         Feedback
// @Accept       json
// @Produce      json
// @Param        feedback body ports.CreateFeedbackInput true "Feedback Details"
// @Success      201 {object} ports.FeedbackResponse // <-- Use DTO
// @Failure      400 {object} map[string]string "Validation error or invalid request body"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /feedback [post]
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

	c.JSON(http.StatusCreated, ports.MapFeedbackToResponse(feedback)) // <-- Use DTO Mapper
}

// @Summary      Get All Feedback
// @Description  Retrieves a list of all feedback submissions. (Admin/Protected)
// @Tags         Feedback
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} ports.FeedbackResponse // <-- Use DTO
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /feedback [get]
func (h *FeedbackHandler) GetAllFeedback(c *gin.Context) {
	// Optional: Check auth user ID if needed for specific logic,
	// but access control is primarily handled by middleware here.
	// --- USE CONSTANT ---
	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		// Should not happen if middleware ran successfully
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	feedbacks, err := h.service.GetAllFeedback(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
		return
	}

	// Map to response DTOs
	response := make([]ports.FeedbackResponse, len(feedbacks))
	for i, fb := range feedbacks {
		response[i] = ports.MapFeedbackToResponse(&fb)
	}

	c.JSON(http.StatusOK, response) // <-- Use mapped DTOs
}

// @Summary      Get Feedback By ID
// @Description  Retrieves a single piece of feedback by its ID. (Admin/Protected)
// @Tags         Feedback
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Feedback ID"
// @Success      200 {object} ports.FeedbackResponse // <-- Use DTO
// @Failure      400 {object} map[string]string "Invalid Feedback ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Feedback not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /feedback/{id} [get]
func (h *FeedbackHandler) GetFeedbackByID(c *gin.Context) {
	// Optional: Check auth user ID if needed
	// --- USE CONSTANT ---
	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// --- USE CONSTANT ---
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

	c.JSON(http.StatusOK, ports.MapFeedbackToResponse(feedback)) // <-- Use DTO Mapper
}

// @Summary      Delete Feedback
// @Description  Deletes a piece of feedback by its ID. (Admin/Protected)
// @Tags         Feedback
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Feedback ID"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Invalid Feedback ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Feedback not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /feedback/{id} [delete]
func (h *FeedbackHandler) DeleteFeedback(c *gin.Context) {
	// Optional: Check auth user ID if needed
	// --- USE CONSTANT ---
	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// --- USE CONSTANT ---
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
