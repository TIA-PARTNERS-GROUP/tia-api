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
	projectService       *services.ProjectService
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

// @Summary Add Project Member
// @Description Adds a user to a project with a specified role. Only accessible by the **Project Manager**.
// @Tags projects, members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param member body ports.AddProjectMemberInput true "Member details (UserID, Role)"
// @Success 201 {object} ports.ProjectMemberResponse "Member added successfully"
// @Failure 400 {object} map[string]interface{} "Invalid project ID, request body, or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the project manager)"
// @Failure 404 {object} map[string]interface{} "ErrProjectNotFound or ErrUserNotFound"
// @Failure 409 {object} map[string]interface{} "ErrProjectMemberAlreadyExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects/{id}/members [post]
func (h *ProjectMemberHandler) AddProjectMember(c *gin.Context) {
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

	var input ports.AddProjectMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

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

// @Summary Get All Project Members
// @Description Retrieves a list of all members associated with a project. Accessible by any authenticated user.
// @Tags projects, members
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} ports.ProjectMembersResponse "List of project members"
// @Failure 400 {object} map[string]interface{} "Invalid project ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects/{id}/members [get]
func (h *ProjectMemberHandler) GetProjectMembers(c *gin.Context) {
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

	members, err := h.projectMemberService.GetProjectMembers(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project members"})
		return
	}
	c.JSON(http.StatusOK, ports.MapToProjectMembersResponse(members))
}

// @Summary Get Specific Project Member
// @Description Retrieves a specific project member record by Project ID and User ID.
// @Tags projects, members
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param userID path int true "User ID of the member"
// @Success 200 {object} ports.ProjectMemberResponse "Member retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrProjectMemberNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects/{id}/members/{userID} [get]
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

// @Summary Get Projects by User
// @Description Retrieves a list of all projects the specified user is a member of. Requires self-management.
// @Tags users, members
// @Produce json
// @Security BearerAuth
// @Param id path int true "Target User ID"
// @Param role query string false "Filter by member role (manager, contributor, reviewer)"
// @Success 200 {object} ports.ProjectMembersResponse "List of project memberships"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the target user)"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/project-memberships [get]
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

	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only view your own project memberships"})
		return
	}

	var role *models.ProjectMemberRole
	roleStr := c.Query("role")
	if roleStr != "" {
		pmRole := models.ProjectMemberRole(roleStr)
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

// @Summary Update Project Member Role
// @Description Updates the role of an existing project member. Only accessible by the **Project Manager**.
// @Tags projects, members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param userID path int true "User ID of the member"
// @Param role body ports.UpdateProjectMemberRoleInput true "New role for the member"
// @Success 200 {object} ports.ProjectMemberResponse "Member role updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID, request body, or ErrInvalidRole"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the project manager)"
// @Failure 404 {object} map[string]interface{} "ErrProjectMemberNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects/{id}/members/{userID} [put]
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

// @Summary Remove Project Member
// @Description Removes a user from a project. Allowed for the **Project Manager** (to remove anyone) or the **User** (to remove self).
// @Tags projects, members
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param userID path int true "User ID of the member to remove"
// @Success 204 "Member removed successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid ID or ErrCannotRemoveManager"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the manager and not the user)"
// @Failure 404 {object} map[string]interface{} "ErrProjectMemberNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /projects/{id}/members/{userID} [delete]
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

	if _, apiErr := h.checkProjectManager(c, uint(projectID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

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
