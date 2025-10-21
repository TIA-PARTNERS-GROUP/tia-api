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
type BusinessTagHandler struct {
	service  *services.BusinessTagService
	validate *validator.Validate
	routes   *constants.Routes 
}
func NewBusinessTagHandler(s *services.BusinessTagService, routes *constants.Routes) *BusinessTagHandler {
	return &BusinessTagHandler{
		service:  s,
		validate: validator.New(),
		routes:   routes, 
	}
}
func (h *BusinessTagHandler) CreateBusinessTag(c *gin.Context) {
	
	businessIDStr := c.Param(h.routes.ParamKeyID)
	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil || businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}
	
	
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	
	_ = authUserID 
	var input ports.CreateBusinessTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	input.BusinessID = uint(businessID) 
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	tag, err := h.service.CreateBusinessTag(c.Request.Context(), input) 
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapToBusinessTagResponse(tag))
}
func (h *BusinessTagHandler) GetBusinessTags(c *gin.Context) {
	
	businessIDStr := c.Param(h.routes.ParamKeyID)
	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil || businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}
	tags, err := h.service.GetBusinessTags(c.Request.Context(), uint(businessID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tags"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToBusinessTagsResponse(tags))
}
func (h *BusinessTagHandler) DeleteBusinessTag(c *gin.Context) {
	
	tagIDStr := c.Param(h.routes.ParamKeyID)
	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil || tagID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}
	
	
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	
	_ = authUserID 
	
	if err := h.service.DeleteBusinessTag(c.Request.Context(), uint(tagID)); err != nil { 
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}
	c.Status(http.StatusNoContent)
}
