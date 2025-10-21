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

type PublicationHandler struct {
	publicationService *services.PublicationService
	validate           *validator.Validate
	routes             *constants.Routes
}

func NewPublicationHandler(publicationService *services.PublicationService, routes *constants.Routes) *PublicationHandler {
	return &PublicationHandler{
		publicationService: publicationService,
		validate:           validator.New(),
		routes:             routes,
	}
}

func (h *PublicationHandler) getAuthUserID(c *gin.Context) (uint, error) {
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

func (h *PublicationHandler) checkPublicationOwnership(c *gin.Context, pubID uint) (uint, *ports.ApiError) {
	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		return 0, ports.ErrInvalidToken
	}

	publication, err := h.publicationService.GetPublicationByID(c.Request.Context(), pubID)
	if err != nil {
		if errors.Is(err, ports.ErrPublicationNotFound) {
			return 0, ports.ErrPublicationNotFound
		}
		return 0, ports.ErrDatabase
	}

	if publication.UserID != authUserID {
		return 0, ports.ErrForbidden
	}
	return authUserID, nil
}

// @Summary Create New Publication
// @Description Creates a new publication (post, article, case study, etc.). The UserID in the body must match the authenticated user.
// @Tags publications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param publication body ports.CreatePublicationInput true "Publication creation details (Title, UserID, Content, Type)"
// @Success 201 {object} ports.PublicationResponse "Publication created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Cannot create publication for another user"
// @Failure 409 {object} map[string]interface{} "ErrPublicationSlugExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /publications [post]
func (h *PublicationHandler) CreatePublication(c *gin.Context) {
	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ports.CreatePublicationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if input.UserID != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Cannot create publication for another user"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publication, err := h.publicationService.CreatePublication(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapPublicationToResponse(publication))
}

// @Summary Get Publication by ID
// @Description Retrieves a specific publication record by its unique ID.
// @Tags publications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Publication ID"
// @Success 200 {object} ports.PublicationResponse "Publication retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid publication ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrPublicationNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /publications/id/{id} [get]
func (h *PublicationHandler) GetPublicationByID(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication ID"})
		return
	}

	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	publication, err := h.publicationService.GetPublicationByID(c.Request.Context(), uint(id))
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapPublicationToResponse(publication))
}

// @Summary Get Publication by Slug
// @Description Retrieves a specific publication record by its unique URL slug.
// @Tags publications
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Publication URL Slug"
// @Success 200 {object} ports.PublicationResponse "Publication retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Missing slug"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrPublicationNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /publications/slug/{slug} [get]
func (h *PublicationHandler) GetPublicationBySlug(c *gin.Context) {
	slug := c.Param(h.routes.ParamKeySlug)
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug is required"})
		return
	}

	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	publication, err := h.publicationService.GetPublicationBySlug(c.Request.Context(), slug)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapPublicationToResponse(publication))
}

// @Summary Get All Publications
// @Description Retrieves a list of all publication records.
// @Tags publications
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ports.PublicationResponse "List of publications"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve publications"
// @Router /publications [get]
func (h *PublicationHandler) GetAllPublications(c *gin.Context) {
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	publications, err := h.publicationService.FindAllPublications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve publications"})
		return
	}

	responses := make([]ports.PublicationResponse, len(publications))
	for i, pub := range publications {
		responses[i] = ports.MapPublicationToResponse(&pub)
	}
	c.JSON(http.StatusOK, responses)
}

// @Summary Update Publication
// @Description Updates an existing publication record. Only the Author (UserID) can perform this action.
// @Tags publications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Publication ID"
// @Param update body ports.UpdatePublicationInput true "Fields to update (Title, Content, Published, etc.)"
// @Success 200 {object} ports.PublicationResponse "Publication updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid publication ID or request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the author)"
// @Failure 404 {object} map[string]interface{} "ErrPublicationNotFound"
// @Failure 409 {object} map[string]interface{} "ErrPublicationSlugExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /publications/{id} [put]
func (h *PublicationHandler) UpdatePublication(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	pubID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication ID"})
		return
	}

	if _, apiErr := h.checkPublicationOwnership(c, uint(pubID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.UpdatePublicationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	publication, err := h.publicationService.UpdatePublication(c.Request.Context(), uint(pubID), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapPublicationToResponse(publication))
}

// @Summary Delete Publication
// @Description Deletes a publication record. Only the Author (UserID) can perform this action.
// @Tags publications
// @Produce json
// @Security BearerAuth
// @Param id path int true "Publication ID"
// @Success 204 "Publication deleted successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid publication ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the author)"
// @Failure 404 {object} map[string]interface{} "ErrPublicationNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /publications/{id} [delete]
func (h *PublicationHandler) DeletePublication(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	pubID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publication ID"})
		return
	}

	if _, apiErr := h.checkPublicationOwnership(c, uint(pubID)); apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	err = h.publicationService.DeletePublication(c.Request.Context(), uint(pubID))
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
