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

type ProjectHandler struct {
	projectService *services.ProjectService
	validate       *validator.Validate
	routes         *constants.Routes
}

func NewProjectHandler(projectService *services.ProjectService, routes *constants.Routes) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		validate:       validator.New(),
		routes:         routes,
	}
}

func (h *ProjectHandler) getAuthUserID(c *gin.Context) (uint, error) {
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

func (h *ProjectHandler) checkProjectManager(c *gin.Context, projectID uint) (uint, *ports.ApiError) {
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

// @Summary Create New Project
// @Description Creates a new project record. Requires authentication.
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project body ports.CreateProjectInput true "Project creation details (Name, ManagedByUserID, ProjectStatus)"
// @Success 201 {object} ports.ProjectResponse "Project created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrManagerNotFound"
// @Failure 409 {object} map[string]interface{} "ErrProjectNameExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	_, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ports.CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.projectService.CreateProject(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapToProjectResponse(project))
}

// @Summary Get Project by ID
// @Description Retrieves a specific project record by its unique ID.
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} ports.ProjectResponse "Project retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid project ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrProjectNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	_, err = h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	project, err := h.projectService.GetProjectByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, ports.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToProjectResponse(project))
}

// @Summary Get All Projects
// @Description Retrieves a list of all project records.
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ports.ProjectResponse "List of projects"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve projects"
// @Router /projects [get]
func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	_, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projects, err := h.projectService.FindAllProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}

	projectResponses := make([]ports.ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = ports.MapToProjectResponse(&project)
	}
	c.JSON(http.StatusOK, projectResponses)
}

// @Summary Update Project Details
// @Description Updates an existing project record. Only the Project Manager can perform this action.
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param project body ports.UpdateProjectInput true "Fields to update"
// @Success 200 {object} ports.ProjectResponse "Project updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid project ID or request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the project manager)"
// @Failure 404 {object} map[string]interface{} "ErrProjectNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.UpdateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), uint(projectID), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToProjectResponse(project))
}

// @Summary Delete Project
// @Description Deletes a project record and all related data (members, regions, skills). Only the Project Manager can perform this action.
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 204 "Project deleted successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid project ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the project manager)"
// @Failure 404 {object} map[string]interface{} "ErrProjectNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	err = h.projectService.DeleteProject(c.Request.Context(), uint(projectID))
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
