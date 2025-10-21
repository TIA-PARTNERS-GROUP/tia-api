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
	projectService       *services.ProjectService
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

// @Summary Add Region to Project
// @Description Associates a geographical region (identified by its short code/ID) with a project. Only accessible by the **Project Manager**.
// @Tags projects, regions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param region body ports.AddProjectRegionInput true "Region details (RegionID)"
// @Success 201 {object} ports.ProjectRegionResponse "Region associated successfully"
// @Failure 400 {object} gin.H "Invalid project ID, request body, or validation error"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the project manager)"
// @Failure 409 {object} gin.H "ErrRegionAlreadyAdded"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/regions [post]
func (h *ProjectRegionHandler) AddRegionToProject(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.AddProjectRegionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

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

// @Summary Remove Region from Project
// @Description Dissociates a geographical region from a project. Only accessible by the **Project Manager**.
// @Tags projects, regions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param regionID path string true "Region ID (e.g., USA, AUS)"
// @Success 204 "Region removed successfully (No Content)"
// @Failure 400 {object} gin.H "Invalid project ID or missing regionID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the project manager)"
// @Failure 404 {object} gin.H "ErrProjectRegionNotFound or ErrProjectNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/regions/{regionID} [delete]
func (h *ProjectRegionHandler) RemoveRegionFromProject(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

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

// @Summary Get Regions for Project
// @Description Retrieves all geographical regions associated with a specific project.
// @Tags projects, regions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {array} ports.ProjectRegionResponse "List of regions"
// @Failure 400 {object} gin.H "Invalid project ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/regions [get]
func (h *ProjectRegionHandler) GetRegionsForProject(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

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
