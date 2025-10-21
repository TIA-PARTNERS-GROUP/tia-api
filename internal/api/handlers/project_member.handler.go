package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProjectMemberHandler struct {
	projectMemberService *services.ProjectMemberService
	projectService       *services.ProjectService // Needed for manager checks
	validate             *validator.Validate
	routes               *constants.Routes
}

func NewProjectMemberHandler(
	projectMemberService *services.ProjectMemberService,
	projectService *services.ProjectService,
	routes *constants.Routes,
) *ProjectMemberHandler {
	return &ProjectMemberHandler{
		projectMemberService: projectMemberService,
		projectService:       projectService,
		validate:             validator.New(),
		routes:               routes,
	}
}

// getAuthUserID retrieves the authenticated user's ID from the context.
func (h *ProjectMemberHandler) getAuthUserID(c *gin.Context) (uint, error) {
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
func (h *ProjectMemberHandler) checkProjectManager(c *gin.Context, projectID uint) (uint, *ports.ApiError) {
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

// AddProjectMember handles POST /projects/:id/members
func (h *ProjectMemberHandler) AddProjectMember(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Auth: Only the project manager can add members
	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.AddProjectMemberInput
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

	member, err := h.projectMemberService.AddProjectMember(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapToProjectMemberResponse(member))
}

// GetProjectMembers handles GET /projects/:id/members
func (h *ProjectMemberHandler) GetProjectMembers(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Auth: Any authenticated user can see project members
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	members, err := h.projectMemberService.GetProjectMembers(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project members"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToProjectMembersResponse(members))
}

// GetProjectMember handles GET /projects/:id/members/:userID
func (h *ProjectMemberHandler) GetProjectMember(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	userIDStr := c.Param(h.routes.ParamKeyUserID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Auth: Any authenticated user can get a specific member
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	member, err := h.projectMemberService.GetProjectMember(c.Request.Context(), uint(projectID), uint(userID))
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToProjectMemberResponse(member))
}

// GetProjectsByUser handles GET /users/:id/project-memberships
func (h *ProjectMemberHandler) GetProjectsByUser(c *gin.Context) {
	targetUserIDStr := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Auth: Users can only see their own project memberships
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only view your own project memberships"})
		return
	}

	// Check for optional "role" query param
	var role *models.ProjectMemberRole
	roleStr := c.Query("role")
	if roleStr != "" {
		pmRole := models.ProjectMemberRole(roleStr)
		// Basic validation for the role enum
		if pmRole == models.ProjectMemberRoleManager || pmRole == models.ProjectMemberRoleContributor || pmRole == models.ProjectMemberRoleReviewer {
			role = &pmRole
		}
	}

	memberships, err := h.projectMemberService.GetProjectsByUser(c.Request.Context(), authUserID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project memberships"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToProjectMembersResponse(memberships))
}

// UpdateProjectMemberRole handles PUT /projects/:id/members/:userID
func (h *ProjectMemberHandler) UpdateProjectMemberRole(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	userIDStr := c.Param(h.routes.ParamKeyUserID)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Auth: Only the project manager can update roles
	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.UpdateProjectMemberRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.projectMemberService.UpdateProjectMemberRole(c.Request.Context(), uint(projectID), uint(userID), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToProjectMemberResponse(member))
}

// RemoveProjectMember handles DELETE /projects/:id/members/:userID
func (h *ProjectMemberHandler) RemoveProjectMember(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	userIDStr := c.Param(h.routes.ParamKeyUserID)
	userIDToRemove, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Auth: Allow if user is removing themselves
	if authUserID == uint(userIDToRemove) {
		err = h.projectMemberService.RemoveProjectMember(c.Request.Context(), uint(projectID), uint(userIDToRemove))
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
		return
	}

	// Auth: Otherwise, check if the user is the project manager
	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	// Manager is allowed to remove anyone (service will block manager self-removal)
	err = h.projectMemberService.RemoveProjectMember(c.Request.Context(), uint(projectID), uint(userIDToRemove))
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
