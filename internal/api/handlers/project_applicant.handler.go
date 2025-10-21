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

type ProjectApplicantHandler struct {
	applicantService *services.ProjectApplicantService
	projectService   *services.ProjectService
	validate         *validator.Validate
	routes           *constants.Routes
}

func NewProjectApplicantHandler(
	applicantService *services.ProjectApplicantService,
	projectService *services.ProjectService,
	routes *constants.Routes,
) *ProjectApplicantHandler {
	return &ProjectApplicantHandler{
		applicantService: applicantService,
		projectService:   projectService,
		validate:         validator.New(),
		routes:           routes,
	}
}

func (h *ProjectApplicantHandler) getAuthUserID(c *gin.Context) (uint, error) {
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

func (h *ProjectApplicantHandler) checkProjectManager(c *gin.Context, projectID uint) (uint, *ports.ApiError) {
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

// @Summary Apply to Project
// @Description Submits an application for the authenticated user to join a project. The UserID is taken from the auth context.
// @Tags projects, applicants
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 201 "Application submitted successfully (Created)"
// @Failure 400 {object} gin.H "Invalid project ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "ErrProjectNotFound or user not found"
// @Failure 409 {object} gin.H "ErrAlreadyApplied"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/apply [post]
func (h *ProjectApplicantHandler) ApplyToProject(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	input := ports.ApplyToProjectInput{
		ProjectID: uint(projectID),
		UserID:    authUserID,
	}

	_, err = h.applicantService.ApplyToProject(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Withdraw Application
// @Description Withdraws the authenticated user's application from a project.
// @Tags projects, applicants
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 204 "Application withdrawn successfully (No Content)"
// @Failure 400 {object} gin.H "Invalid project ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "ErrApplicationNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/apply [delete]
func (h *ProjectApplicantHandler) WithdrawApplication(c *gin.Context) {
	projectIDStr := c.Param(h.routes.ParamKeyID)
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = h.applicantService.WithdrawApplication(c.Request.Context(), uint(projectID), authUserID)
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

// @Summary Get Applicants for Project
// @Description Retrieves a list of all users who have applied to a specific project. Only accessible by the **Project Manager**.
// @Tags projects, applicants
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {array} ports.ProjectApplicantResponse "List of applicants"
// @Failure 400 {object} gin.H "Invalid project ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the project manager)"
// @Failure 404 {object} gin.H "ErrProjectNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /projects/{id}/applicants [get]
func (h *ProjectApplicantHandler) GetApplicantsForProject(c *gin.Context) {
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

	applicants, err := h.applicantService.GetApplicantsForProject(c.Request.Context(), uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applicants"})
		return
	}

	applicantResponses := make([]ports.ProjectApplicantResponse, len(applicants))
	for i, app := range applicants {
		applicantResponses[i] = ports.MapProjectApplicantToResponse(&app)
	}

	c.JSON(http.StatusOK, applicantResponses)
}

// @Summary Get Applications for User
// @Description Retrieves a list of all projects the specified user has applied to. Requires authentication and self-management.
// @Tags users, applicants
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {array} ports.UserApplicationResponse "List of user applications"
// @Failure 400 {object} gin.H "Invalid user ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the target user)"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/{id}/applications [get]
func (h *ProjectApplicantHandler) GetApplicationsForUser(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only view your own applications"})
		return
	}

	applications, err := h.applicantService.GetApplicationsForUser(c.Request.Context(), authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	applicationResponses := make([]ports.UserApplicationResponse, len(applications))
	for i, app := range applications {
		applicationResponses[i] = ports.MapUserApplicationToResponse(&app)
	}

	c.JSON(http.StatusOK, applicationResponses)
}
