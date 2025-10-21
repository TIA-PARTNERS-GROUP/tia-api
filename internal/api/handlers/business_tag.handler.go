package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BusinessTagHandler struct {
	service  *services.BusinessTagService
	validate *validator.Validate
	routes   *constants.Routes // <-- ADDED
}

// Updated constructor
func NewBusinessTagHandler(s *services.BusinessTagService, routes *constants.Routes) *BusinessTagHandler {
	return &BusinessTagHandler{
		service:  s,
		validate: validator.New(),
		routes:   routes, // <-- ADDED
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
	// --- USE CONSTANT ---
	businessIDStr := c.Param(h.routes.ParamKeyID)
	businessID, err := strconv.ParseUint(businessIDStr, 10, 32)
	if err != nil || businessID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	// Retrieve auth user ID for potential service layer checks (e.g., is user operator of businessID?)
	// --- USE CONSTANT ---
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	// Note: You might want to pass authUserID to the service to verify permissions.
	_ = authUserID // Placeholder if not used directly here

	var input ports.CreateBusinessTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	input.BusinessID = uint(businessID) // Set business ID from path param

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pass authUserID if service needs to check permissions
	tag, err := h.service.CreateBusinessTag(c.Request.Context(), input) // Modify service if needed
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
	// --- USE CONSTANT ---
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
	// --- USE CONSTANT ---
	tagIDStr := c.Param(h.routes.ParamKeyID)
	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil || tagID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	// Retrieve auth user ID for potential service layer checks (e.g., does user own the business associated with tagID?)
	// --- USE CONSTANT ---
	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, _ := authUserIDVal.(uint)
	if authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	// Note: Pass authUserID to service if it needs to verify permissions.
	_ = authUserID // Placeholder

	// Pass authUserID if service needs to check permissions
	if err := h.service.DeleteBusinessTag(c.Request.Context(), uint(tagID)); err != nil { // Modify service if needed
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
