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

// getAuthUserID retrieves the authenticated user's ID from the context.
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

// checkUserOwnership verifies if the auth user matches the target user in the URL.
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

// AddUserSkill handles POST /users/:id/skills
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

	// Override UserID from DTO with the one from the URL param for consistency
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

// GetUserSkills handles GET /users/:id/skills
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

// UpdateUserSkill handles PUT /users/:id/skills/:skillID
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

// RemoveUserSkill handles DELETE /users/:id/skills/:skillID
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
