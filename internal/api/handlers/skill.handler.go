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

type SkillHandler struct {
	skillService *services.SkillService
	validate     *validator.Validate
	routes       *constants.Routes
}

func NewSkillHandler(skillService *services.SkillService, routes *constants.Routes) *SkillHandler {
	return &SkillHandler{
		skillService: skillService,
		validate:     validator.New(),
		routes:       routes,
	}
}

// getAuthUserID retrieves the authenticated user's ID from the context.
func (h *SkillHandler) getAuthUserID(c *gin.Context) (uint, error) {
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

// GetSkills handles GET /skills (with optional filters)
func (h *SkillHandler) GetSkills(c *gin.Context) {
	// Auth: Any authenticated user can view skills
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var filters ports.SkillsFilter
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	skills, err := h.skillService.GetSkills(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve skills"})
		return
	}

	responses := make([]ports.SkillResponse, len(skills))
	for i, skill := range skills {
		responses[i] = ports.MapSkillToResponse(&skill)
	}
	c.JSON(http.StatusOK, responses)
}

// GetSkillByID handles GET /skills/:id
func (h *SkillHandler) GetSkillByID(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Auth: Any authenticated user can view a skill
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	skill, err := h.skillService.GetSkillByID(c.Request.Context(), uint(id))
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapSkillToResponse(skill))
}

// CreateSkill handles POST /skills
func (h *SkillHandler) CreateSkill(c *gin.Context) {
	// Auth: Check authentication
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ports.CreateSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skill, err := h.skillService.CreateSkill(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapSkillToResponse(skill))
}

// UpdateSkill handles PUT /skills/:id
func (h *SkillHandler) UpdateSkill(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Auth: Check authentication
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ports.UpdateSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skill, err := h.skillService.UpdateSkill(c.Request.Context(), uint(id), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapSkillToResponse(skill))
}

// DeleteSkill handles DELETE /skills/:id
func (h *SkillHandler) DeleteSkill(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Auth: Check authentication
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = h.skillService.DeleteSkill(c.Request.Context(), uint(id))
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

// ToggleSkillStatus handles PATCH /skills/:id/toggle-status
func (h *SkillHandler) ToggleSkillStatus(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Auth: Check authentication
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	skill, err := h.skillService.ToggleSkillStatus(c.Request.Context(), uint(id))
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapSkillToResponse(skill))
}
