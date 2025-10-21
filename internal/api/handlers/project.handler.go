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

// getAuthUserID retrieves the authenticated user's ID from the context.
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

// checkProjectManager is a helper to verify if the auth user is the project manager.
func (h *ProjectHandler) checkProjectManager(c *gin.Context, projectID uint) (uint, *ports.ApiError) {
	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		// Use a predefined error for auth issues
		return 0, ports.ErrInvalidToken
	}

	project, err := h.projectService.GetProjectByID(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, ports.ErrProjectNotFound) {
			// Use the predefined error
			return 0, ports.ErrProjectNotFound
		}
		// Use the generic database error
		return 0, ports.ErrDatabase
	}

	if project.ManagedByUserID != authUserID {
		// Use the predefined forbidden error
		return 0, ports.ErrForbidden
	}
	return authUserID, nil
}

// CreateProject handles POST /projects
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

// GetProjectByID handles GET /projects/:id
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

// GetAllProjects handles GET /projects
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

// UpdateProject handles PUT /projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check if the user is the project manager
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

// DeleteProject handles DELETE /projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check if the user is the project manager
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
