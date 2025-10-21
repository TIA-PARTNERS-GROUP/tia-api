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

type BusinessHandler struct {
	businessService *services.BusinessService
	validate        *validator.Validate
	routes          *constants.Routes
}

func NewBusinessHandler(businessService *services.BusinessService, routes *constants.Routes) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
		validate:        validator.New(),
		routes:          routes,
	}
}

// CreateBusiness (Authorization check remains here as it prevents unnecessary service calls)
func (h *BusinessHandler) CreateBusiness(c *gin.Context) {
	var input ports.CreateBusinessInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Authorization Check
	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	if authUserID != input.OperatorUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only create a business for yourself"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business, err := h.businessService.CreateBusiness(c.Request.Context(), input) // No authUserID passed
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create business"})
		return
	}

	c.JSON(http.StatusCreated, ports.MapBusinessToResponse(business))
}

// GetBusinessByID (No changes needed)
func (h *BusinessHandler) GetBusinessByID(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	business, err := h.businessService.GetBusinessByID(c.Request.Context(), uint(id))
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve business"})
		return
	}

	c.JSON(http.StatusOK, ports.MapBusinessToResponse(business))
}

// GetBusinesses (No changes needed)
func (h *BusinessHandler) GetBusinesses(c *gin.Context) {
	var filters ports.BusinessesFilter
	if err := c.BindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	businesses, err := h.businessService.GetBusinesses(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve businesses"})
		return
	}

	businessResponses := make([]ports.BusinessResponse, len(businesses))
	for i, biz := range businesses {
		businessResponses[i] = ports.MapBusinessToResponse(&biz)
	}

	c.JSON(http.StatusOK, businessResponses)
}

// In business.handler.go

// UpdateBusiness
func (h *BusinessHandler) UpdateBusiness(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// --- FIX: Store the authUserID ---
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 { // Also check if authUserID is zero
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	var input ports.UpdateBusinessInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Assuming the service `UpdateBusiness` handles the auth check internally.
	// --- FIX: Add authUserID argument BACK ---
	business, err := h.businessService.UpdateBusiness(c.Request.Context(), uint(id), authUserID, input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			// Check specifically for Forbidden error potentially returned by the service
			if apiErr.StatusCode == http.StatusForbidden {
				c.JSON(http.StatusForbidden, gin.H{"error": apiErr.Message})
				return
			}
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update business"})
		return
	}

	c.JSON(http.StatusOK, ports.MapBusinessToResponse(business))
}

// DeleteBusiness
func (h *BusinessHandler) DeleteBusiness(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// --- FIX: Store the authUserID ---
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 { // Also check if authUserID is zero
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	// Assuming the service `DeleteBusiness` handles the auth check internally.
	// --- FIX: Add authUserID argument BACK ---
	err = h.businessService.DeleteBusiness(c.Request.Context(), uint(id), authUserID)
	// --- FIX: Remove incorrect UpdateBusiness call from DeleteBusiness handler ---
	// business, err := h.businessService.UpdateBusiness(c.Request.Context(), uint(id), authUserID, input) // This line was incorrect
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			// Check specifically for Forbidden error potentially returned by the service
			if apiErr.StatusCode == http.StatusForbidden {
				c.JSON(http.StatusForbidden, gin.H{"error": apiErr.Message})
				return
			}
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete business"})
		return
	}

	c.Status(http.StatusNoContent)
}
