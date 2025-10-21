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
