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
type EventHandler struct {
	service  *services.EventService
	validate *validator.Validate 
	routes   *constants.Routes   
}
func NewEventHandler(service *services.EventService, routes *constants.Routes) *EventHandler {
	return &EventHandler{
		service:  service,
		validate: validator.New(), 
		routes:   routes,          
	}
}
func (h *EventHandler) CreateEvent(c *gin.Context) {
	
	userIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	
	
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}) 
		return
	}
	userID, _ := userIDVal.(uint) 
	var input ports.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}
	
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	input.UserID = &userID
	event, err := h.service.CreateEvent(c.Request.Context(), input)
	if err != nil {
		
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapEventToResponse(event)) 
}
func (h *EventHandler) GetEventByID(c *gin.Context) {
	
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
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	c.JSON(http.StatusOK, ports.MapEventToResponse(event)) 
}
func (h *EventHandler) GetEvents(c *gin.Context) {
	var filters ports.EventsFilter
	
	if err := c.BindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}
	events, err := h.service.GetEvents(c.Request.Context(), filters)
	if err != nil {
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	
	response := make([]ports.EventResponse, len(events))
	for i, event := range events {
		response[i] = ports.MapEventToResponse(&event)
	}
	c.JSON(http.StatusOK, response) 
}
