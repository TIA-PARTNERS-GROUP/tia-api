package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// @Summary      Create Event
// @Description  Creates a new event
// @Tags         Events
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input body ports.CreateEventInput true "Event input"
// @Success      201 {object} models.Event
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /events/ [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	userID, _ := userIDVal.(uint)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	var input ports.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	input.UserID = &userID

	event, err := h.service.CreateEvent(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ports.ErrDatabase) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// @Summary      Get Event by ID
// @Description  Retrieves a specific event by its ID
// @Tags         Events
// @Produce      json
// @Param        id path int true "Event ID"
// @Success      200 {object} models.Event
// @Failure      400 {object} map[string]string "Invalid event ID"
// @Failure      404 {object} map[string]string "Event not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /events/{id}/ [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 32)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// @Summary      Get Events
// @Description  Retrieves events with optional filtering
// @Tags         Events
// @Produce      json
// @Param        event_type query string false "Filter by event type"
// @Param        user_id query int false "Filter by user ID"
// @Success      200 {array} models.Event
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /events/ [get]
func (h *EventHandler) GetEvents(c *gin.Context) {
	var filters ports.EventsFilter

	if eventType := c.Query("event_type"); eventType != "" {
		filters.EventType = &eventType
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userIDUint := uint(userID)
			filters.UserID = &userIDUint
		}
	}

	events, err := h.service.GetEvents(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
