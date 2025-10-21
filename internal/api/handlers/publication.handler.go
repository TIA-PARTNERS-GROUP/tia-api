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
