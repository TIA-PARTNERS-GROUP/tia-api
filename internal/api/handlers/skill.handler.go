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

// @Summary Get All Skills with Filters
// @Description Retrieves a list of all skills, with options to filter by category, activity status, or search term.
// @Tags skills
// @Produce json
// @Security BearerAuth
// @Param category query string false "Filter by skill category"
// @Param active query bool false "Filter by active status (true/false)"
// @Param search query string false "Search by name, category, or description"
// @Success 200 {array} ports.SkillResponse "List of skills"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills [get]
func (h *SkillHandler) GetSkills(c *gin.Context) {
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

// @Summary Get Skill by ID
// @Description Retrieves a specific skill record by its unique ID.
// @Tags skills
// @Produce json
// @Security BearerAuth
// @Param id path int true "Skill ID"
// @Success 200 {object} ports.SkillResponse "Skill retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid skill ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrSkillNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills/{id} [get]
func (h *SkillHandler) GetSkillByID(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

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

// @Summary Create New Skill
// @Description Creates a new global skill record.
// @Tags skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param skill body ports.CreateSkillInput true "Skill creation details (Name, Category)"
// @Success 201 {object} ports.SkillResponse "Skill created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "ErrSkillNameExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills [post]
func (h *SkillHandler) CreateSkill(c *gin.Context) {
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

// @Summary Update Skill
// @Description Updates the details of an existing skill (e.g., category, name).
// @Tags skills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Skill ID"
// @Param update body ports.UpdateSkillInput true "Fields to update (Category, Name, Active)"
// @Success 200 {object} ports.SkillResponse "Skill updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid skill ID or request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrSkillNotFound"
// @Failure 409 {object} map[string]interface{} "ErrSkillNameExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills/{id} [put]
func (h *SkillHandler) UpdateSkill(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

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

// @Summary Delete Skill
// @Description Deletes a specific skill record. Fails if the skill is currently in use by a user or project.
// @Tags skills
// @Produce json
// @Security BearerAuth
// @Param id path int true "Skill ID"
// @Success 204 "Skill deleted successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid skill ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrSkillNotFound"
// @Failure 409 {object} map[string]interface{} "ErrSkillInUse"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills/{id} [delete]
func (h *SkillHandler) DeleteSkill(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

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

// @Summary Toggle Skill Status
// @Description Toggles the active status of a skill (Active -> Inactive, or vice versa).
// @Tags skills
// @Produce json
// @Security BearerAuth
// @Param id path int true "Skill ID"
// @Success 200 {object} ports.SkillResponse "Skill status toggled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid skill ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrSkillNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /skills/{id}/toggle-status [patch]
func (h *SkillHandler) ToggleSkillStatus(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

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
