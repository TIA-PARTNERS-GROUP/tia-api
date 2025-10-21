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
	projectService      *services.ProjectService
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

// @Summary Add Skill Requirement to Project
// @Description Associates a specific skill (by Skill ID) with a project and sets its importance level. Only accessible by the **Project Manager**.
// @Tags projects, skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param skill body ports.CreateProjectSkillInput true "Skill details (SkillID, Importance)"
// @Success 201 {object} ports.ProjectSkillResponse "Skill requirement added successfully"
// @Failure 400 {object} gin.H "Invalid project ID, skill ID, request body, or validation error"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the project manager)"
// @Failure 404 {object} gin.H "ErrProjectNotFound or ErrSkillNotFound"
// @Failure 409 {object} gin.H "ErrProjectSkillAlreadyExists"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/skills [post]
func (h *ProjectSkillHandler) AddProjectSkill(c *gin.Context) {
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

	var input ports.CreateProjectSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

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

// @Summary Get Skills Required by Project
// @Description Retrieves all skill requirements for a specific project.
// @Tags projects, skills
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} ports.ProjectSkillsResponse "List of required skills"
// @Failure 400 {object} gin.H "Invalid project ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/skills [get]
func (h *ProjectSkillHandler) GetProjectSkills(c *gin.Context) {
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

	skills, err := h.projectSkillService.GetProjectSkills(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project skills"})
		return
	}

	c.JSON(http.StatusOK, ports.MapToProjectSkillsResponse(skills))
}

// @Summary Update Project Skill Importance
// @Description Updates the importance level (required, preferred, optional) for an existing project skill. Only accessible by the **Project Manager**.
// @Tags projects, skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param skillID path int true "Skill ID"
// @Param update body ports.UpdateProjectSkillInput true "New importance level"
// @Success 200 {object} ports.ProjectSkillResponse "Skill importance updated successfully"
// @Failure 400 {object} gin.H "Invalid ID, request body, or validation error"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the project manager)"
// @Failure 404 {object} gin.H "ErrProjectSkillNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/skills/{skillID} [put]
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

// @Summary Remove Skill Requirement from Project
// @Description Removes a skill requirement association from a project. Only accessible by the **Project Manager**.
// @Tags projects, skills
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param skillID path int true "Skill ID"
// @Success 204 "Skill requirement removed successfully (No Content)"
// @Failure 400 {object} gin.H "Invalid ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the project manager)"
// @Failure 404 {object} gin.H "ErrProjectSkillNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/skills/{skillID} [delete]
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
