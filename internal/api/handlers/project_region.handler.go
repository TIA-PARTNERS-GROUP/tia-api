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

type ProjectRegionHandler struct {
	projectRegionService *services.ProjectRegionService
	projectService       *services.ProjectService // Needed for manager checks
	validate             *validator.Validate
	routes               *constants.Routes
}

func NewProjectRegionHandler(
	projectRegionService *services.ProjectRegionService,
	projectService *services.ProjectService,
	routes *constants.Routes,
) *ProjectRegionHandler {
	return &ProjectRegionHandler{
		projectRegionService: projectRegionService,
		projectService:       projectService,
		validate:             validator.New(),
		routes:               routes,
	}
}

// getAuthUserID retrieves the authenticated user's ID from the context.
func (h *ProjectRegionHandler) getAuthUserID(c *gin.Context) (uint, error) {
	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		return 0, errors.New("invalid authentication context")
	}
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		return 0, errors.New("invalid authentication context")
	}
	return authUserID, nil
}

// checkProjectManager is a helper to verify if the auth user is the project manager.
func (h *ProjectRegionHandler) checkProjectManager(c *gin.Context, projectID uint) (uint, *ports.ApiError) {
	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		return 0, ports.ErrInvalidToken
	}

	project, err := h.projectService.GetProjectByID(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, ports.ErrProjectNotFound) {
			return 0, ports.ErrProjectNotFound
		}
		return 0, ports.ErrDatabase
	}

	if project.ManagedByUserID != authUserID {
		return 0, ports.ErrForbidden
	}
	return authUserID, nil
}

// AddRegionToProject handles POST /projects/:id/regions
func (h *ProjectRegionHandler) AddRegionToProject(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Auth: Only the project manager can add regions
	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.AddProjectRegionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Override ProjectID from DTO with the one from the URL param for consistency
	input.ProjectID = uint(projectID)

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	association, err := h.projectRegionService.AddRegionToProject(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusCreated, ports.MapProjectRegionToResponse(association))
}

// RemoveRegionFromProject handles DELETE /projects/:id/regions/:regionID
func (h *ProjectRegionHandler) RemoveRegionFromProject(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Auth: Only the project manager can remove regions
	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	regionID := c.Param(h.routes.ParamKeyRegionID)
	if regionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region ID is required"})
		return
	}

	err = h.projectRegionService.RemoveRegionFromProject(c.Request.Context(), uint(projectID), regionID)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetRegionsForProject handles GET /projects/:id/regions
func (h *ProjectRegionHandler) GetRegionsForProject(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Auth: Any authenticated user can see project regions
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	regions, err := h.projectRegionService.GetRegionsForProject(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve regions"})
		return
	}

	regionResponses := make([]ports.ProjectRegionResponse, len(regions))
	for i, pr := range regions {
		regionResponses[i] = ports.MapProjectRegionToResponse(&pr)
	}

	c.JSON(http.StatusOK, regionResponses)
}
