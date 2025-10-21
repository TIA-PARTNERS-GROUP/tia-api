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

// @Summary Add Tag to Business
// @Description Creates a new tag (e.g., 'client', 'service') and associates it with a specific business.
// @Tags businesses, tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Business ID to associate the tag with"
// @Param tag body ports.CreateBusinessTagInput true "Tag details (TagType, Description)"
// @Success 201 {object} ports.BusinessTagResponse "Tag created and associated successfully"
// @Failure 400 {object} gin.H "Invalid business ID, request body, or validation failed"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 409 {object} gin.H "ErrBusinessTagAlreadyExists"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /businesses/{id}/tags [post]
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

// @Summary Get All Tags for a Business
// @Description Retrieves all tags associated with a specific business.
// @Tags businesses, tags
// @Produce json
// @Param id path int true "Business ID"
// @Success 200 {object} ports.BusinessTagsResponse "List of business tags"
// @Failure 400 {object} gin.H "Invalid business ID"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /businesses/{id}/tags [get]
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

// @Summary Delete Business Tag
// @Description Deletes a specific business tag entry by its unique Tag ID.
// @Tags businesses, tags
// @Produce json
// @Security BearerAuth
// @Param id path int true "Unique Business Tag ID (NOT the Business ID)"
// @Success 204 "Tag deleted successfully (No Content)"
// @Failure 400 {object} gin.H "Invalid tag ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "ErrBusinessTagNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /tags/{id} [delete]
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
