package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BusinessTagHandler struct {
	service  *services.BusinessTagService
	validate *validator.Validate
}

func NewBusinessTagHandler(s *services.BusinessTagService) *BusinessTagHandler {
	return &BusinessTagHandler{
		service:  s,
		validate: validator.New(),
	}
}

// @Summary      Add a tag to a business
// @Description  Creates and associates a new tag with a specific business. Requires authentication.
// @Tags         Business Tags
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Business ID"
// @Param        tag body ports.CreateBusinessTagInput true "Tag Data (BusinessID in body is ignored)"
// @Success      201 {object} ports.BusinessTagResponse
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      404 {object} map[string]string "Business not found"
// @Failure      409 {object} map[string]string "Tag already exists"
// @Router       /businesses/{id}/tags [post]
func (h *BusinessTagHandler) CreateBusinessTag(c *gin.Context) {
	businessID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

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

// @Summary      Get tags for a business
// @Description  Retrieves all tags associated with a specific business.
// @Tags         Business Tags
// @Produce      json
// @Param        id   path      int  true  "Business ID"
// @Success      200 {object} ports.BusinessTagsResponse
// @Failure      400 {object} map[string]string "Invalid business ID"
// @Router       /businesses/{id}/tags [get]
func (h *BusinessTagHandler) GetBusinessTags(c *gin.Context) {
	businessID, err := strconv.ParseUint(c.Param("id"), 10, 32)
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

// @Summary      Delete a business tag
// @Description  Removes a tag from a business by the tag's unique ID. Requires authentication.
// @Tags         Business Tags
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Tag ID"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Invalid tag ID"
// @Failure      404 {object} map[string]string "Tag not found"
// @Router       /tags/{id} [delete]
func (h *BusinessTagHandler) DeleteBusinessTag(c *gin.Context) {
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || tagID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

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
