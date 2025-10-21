package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"     // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services" // <-- IMPORT models
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10" // <-- IMPORT validator
)

type EventHandler struct {
	service  *services.EventService
	validate *validator.Validate // <-- ADD
	routes   *constants.Routes   // <-- ADD
}

// Updated constructor
func NewEventHandler(service *services.EventService, routes *constants.Routes) *EventHandler {
	return &EventHandler{
		service:  service,
		validate: validator.New(), // <-- ADD
		routes:   routes,          // <-- ADD
	}
}

// @Summary      Create Event
// @Description  Creates a new event
// @Tags         Events
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input body ports.CreateEventInput true "Event input"
// @Success      201 {object} ports.EventResponse // <-- Use DTO
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /events [post] // Corrected path
func (h *EventHandler) CreateEvent(c *gin.Context) {
	// --- USE CONSTANT ---
	userIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	// No need to check for userID == 0 here, the middleware ensures it exists if auth is required.
	// If the route allows anonymous event creation, this check might differ.
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}) // Or handle anonymous case
		return
	}
	userID, _ := userIDVal.(uint) // Assume type assertion is safe after exists check

	var input ports.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	// Validate the input DTO
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Associate the event with the authenticated user
	input.UserID = &userID

	event, err := h.service.CreateEvent(c.Request.Context(), input)
	if err != nil {
		// Log the actual error for debugging
		// log.Printf("Error creating event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, ports.MapEventToResponse(event)) // <-- Use DTO Mapper
}

// @Summary      Get Event by ID
// @Description  Retrieves a specific event by its ID
// @Tags         Events
// @Produce      json
// @Param        id path int true "Event ID"
// @Success      200 {object} ports.EventResponse // <-- Use DTO
// @Failure      400 {object} map[string]string "Invalid event ID"
// @Failure      404 {object} map[string]string "Event not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /events/{id} [get] // Corrected path
func (h *EventHandler) GetEventByID(c *gin.Context) {
	// --- USE CONSTANT ---
	eventIDStr := c.Param(h.routes.ParamKeyID)
	eventID, err := strconv.ParseUint(eventIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := h.service.GetEventByID(c.Request.Context(), uint(eventID))
	if err != nil {
		if errors.Is(err, ports.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}
		// Log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	c.JSON(http.StatusOK, ports.MapEventToResponse(event)) // <-- Use DTO Mapper
}

// @Summary      Get Events
// @Description  Retrieves events with optional filtering
// @Tags         Events
// @Produce      json
// @Param        event_type query string false "Filter by event type"
// @Param        user_id query int false "Filter by user ID"
// @Success      200 {array} ports.EventResponse // <-- Use DTO
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /events [get] // Corrected path
func (h *EventHandler) GetEvents(c *gin.Context) {
	var filters ports.EventsFilter

	// Use BindQuery for cleaner query parameter handling
	if err := c.BindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	events, err := h.service.GetEvents(c.Request.Context(), filters)
	if err != nil {
		// Log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	// Map results to response DTOs
	response := make([]ports.EventResponse, len(events))
	for i, event := range events {
		response[i] = ports.MapEventToResponse(&event)
	}

	c.JSON(http.StatusOK, response) // <-- Use mapped DTOs
}
