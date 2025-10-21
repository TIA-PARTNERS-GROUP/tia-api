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

type UserSkillHandler struct {
	userSkillService *services.UserSkillService
	validate         *validator.Validate
	routes           *constants.Routes
}

func NewUserSkillHandler(userSkillService *services.UserSkillService, routes *constants.Routes) *UserSkillHandler {
	return &UserSkillHandler{
		userSkillService: userSkillService,
		validate:         validator.New(),
		routes:           routes,
	}
}

func (h *UserSkillHandler) getAuthUserID(c *gin.Context) (uint, error) {
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

func (h *UserSkillHandler) checkUserOwnership(c *gin.Context) (uint, uint, *ports.ApiError) {
	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		return 0, 0, ports.ErrInvalidToken
	}

	targetUserIDStr := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return 0, 0, ports.ErrDatabase
	}

	if authUserID != uint(targetUserID) {
		return 0, 0, ports.ErrForbidden
	}

	return authUserID, uint(targetUserID), nil
}

// @Summary Add User Skill
// @Description Adds a new skill and its proficiency level to the authenticated user's profile. Enforces self-management.
// @Tags users, skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param skill body ports.CreateUserSkillInput true "Skill details (SkillID, ProficiencyLevel)"
// @Success 201 {object} ports.UserSkillResponse "Skill added successfully"
// @Failure 400 {object} gin.H "Invalid request body or validation error"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the target user)"
// @Failure 404 {object} gin.H "ErrUserNotFound or ErrSkillNotFound"
// @Failure 409 {object} gin.H "ErrUserSkillAlreadyExists"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/{id}/skills [post]
func (h *UserSkillHandler) AddUserSkill(c *gin.Context) {
	_, targetUserID, apiErr := h.checkUserOwnership(c)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.CreateUserSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	input.UserID = targetUserID

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userSkill, err := h.userSkillService.AddUserSkill(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapToUserSkillResponse(userSkill))
}

// @Summary Get User Skills
// @Description Retrieves all skills and proficiency levels associated with the user. Enforces self-management.
// @Tags users, skills
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} ports.UserSkillsResponse "List of user skills"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the target user)"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/{id}/skills [get]
func (h *UserSkillHandler) GetUserSkills(c *gin.Context) {
	_, targetUserID, apiErr := h.checkUserOwnership(c)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	skills, err := h.userSkillService.GetUserSkills(c.Request.Context(), targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user skills"})
		return
	}

	c.JSON(http.StatusOK, ports.MapToUserSkillsResponse(skills))
}

// @Summary Update User Skill Proficiency
// @Description Updates the proficiency level for an existing skill associated with the user. Enforces self-management.
// @Tags users, skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param skillID path int true "Skill ID"
// @Param update body ports.UpdateUserSkillInput true "New proficiency level"
// @Success 200 {object} ports.UserSkillResponse "Proficiency updated successfully"
// @Failure 400 {object} gin.H "Invalid ID, request body, or validation error"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the target user)"
// @Failure 404 {object} gin.H "ErrUserSkillNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/{id}/skills/{skillID} [put]
func (h *UserSkillHandler) UpdateUserSkill(c *gin.Context) {
	_, targetUserID, apiErr := h.checkUserOwnership(c)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	skillIDStr := c.Param(h.routes.ParamKeySkillID)
	skillID, err := strconv.ParseUint(skillIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	var input ports.UpdateUserSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userSkill, err := h.userSkillService.UpdateUserSkill(c.Request.Context(), targetUserID, uint(skillID), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusOK, ports.MapToUserSkillResponse(userSkill))
}

// @Summary Remove User Skill
// @Description Removes a skill association from the user's profile. Enforces self-management.
// @Tags users, skills
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param skillID path int true "Skill ID"
// @Success 204 "Skill removed successfully (No Content)"
// @Failure 400 {object} gin.H "Invalid ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden (Not the target user)"
// @Failure 404 {object} gin.H "ErrUserSkillNotFound"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /users/{id}/skills/{skillID} [delete]
func (h *UserSkillHandler) RemoveUserSkill(c *gin.Context) {
	_, targetUserID, apiErr := h.checkUserOwnership(c)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	skillIDStr := c.Param(h.routes.ParamKeySkillID)
	skillID, err := strconv.ParseUint(skillIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	err = h.userSkillService.RemoveUserSkill(c.Request.Context(), targetUserID, uint(skillID))
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
