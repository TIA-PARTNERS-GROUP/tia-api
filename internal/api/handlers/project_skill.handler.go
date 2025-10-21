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

type ProjectSkillHandler struct {
	projectSkillService *services.ProjectSkillService
	projectService      *services.ProjectService // Needed for manager checks
	validate            *validator.Validate
	routes              *constants.Routes
}

func NewProjectSkillHandler(
	projectSkillService *services.ProjectSkillService,
	projectService *services.ProjectService,
	routes *constants.Routes,
) *ProjectSkillHandler {
	return &ProjectSkillHandler{
		projectSkillService: projectSkillService,
		projectService:      projectService,
		validate:            validator.New(),
		routes:              routes,
	}
}

// getAuthUserID retrieves the authenticated user's ID from the context.
func (h *ProjectSkillHandler) getAuthUserID(c *gin.Context) (uint, error) {
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
func (h *ProjectSkillHandler) checkProjectManager(c *gin.Context, projectID uint) (uint, *ports.ApiError) {
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

// AddProjectSkill handles POST /projects/:id/skills
func (h *ProjectSkillHandler) AddProjectSkill(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Auth: Only the project manager can add skills
	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.CreateProjectSkillInput
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

	projectSkill, err := h.projectSkillService.AddProjectSkill(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusCreated, ports.MapToProjectSkillResponse(projectSkill))
}

// GetProjectSkills handles GET /projects/:id/skills
func (h *ProjectSkillHandler) GetProjectSkills(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Auth: Any authenticated user can view project skills
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	skills, err := h.projectSkillService.GetProjectSkills(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project skills"})
		return
	}

	c.JSON(http.StatusOK, ports.MapToProjectSkillsResponse(skills))
}

// UpdateProjectSkill handles PUT /projects/:id/skills/:skillID
func (h *ProjectSkillHandler) UpdateProjectSkill(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	skillIDStr := c.Param(h.routes.ParamKeySkillID)
	skillID, err := strconv.ParseUint(skillIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Auth: Only the project manager can update skills
	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.UpdateProjectSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectSkill, err := h.projectSkillService.UpdateProjectSkill(c.Request.Context(), uint(projectID), uint(skillID), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusOK, ports.MapToProjectSkillResponse(projectSkill))
}

// RemoveProjectSkill handles DELETE /projects/:id/skills/:skillID
func (h *ProjectSkillHandler) RemoveProjectSkill(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	skillIDStr := c.Param(h.routes.ParamKeySkillID)
	skillID, err := strconv.ParseUint(skillIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Auth: Only the project manager can remove skills
	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	err = h.projectSkillService.RemoveProjectSkill(c.Request.Context(), uint(projectID), uint(skillID))
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
