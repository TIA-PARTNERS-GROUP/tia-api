package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BusinessHandler struct {
	businessService *services.BusinessService
	validate        *validator.Validate
}

func NewBusinessHandler(businessService *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
		validate:        validator.New(),
	}
}

// @Summary      Create a new business
// @Description  Creates a new business record for the authenticated user.
// @Tags         Businesses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        business body ports.CreateBusinessInput true "Business Creation Data"
// @Success      201 {object} ports.BusinessResponse
// @Failure      400 {object} map[string]string "Validation error"
// @Failure      404 {object} map[string]string "Operator user not found"
// @Router       /businesses [post]
func (h *BusinessHandler) CreateBusiness(c *gin.Context) {
	var input ports.CreateBusinessInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	business, err := h.businessService.CreateBusiness(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusCreated, ports.MapBusinessToResponse(business))
}

// @Summary      Get a business by ID
// @Description  Retrieves the details of a single business by its unique ID.
// @Tags         Businesses
// @Produce      json
// @Param        id   path      int  true  "Business ID"
// @Success      200  {object}  ports.BusinessResponse
// @Failure      404  {object}  map[string]string "Business not found"
// @Router       /businesses/{id} [get]
func (h *BusinessHandler) GetBusinessByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusOK, ports.MapBusinessToResponse(business))
}

// @Summary      Get all businesses
// @Description  Retrieves a list of all businesses, with optional filtering.
// @Tags         Businesses
// @Produce      json
// @Param        active query bool false "Filter by active status"
// @Param        search query string false "Search term for name, tagline, or description"
// @Success      200 {array} ports.BusinessResponse
// @Router       /businesses [get]
func (h *BusinessHandler) GetBusinesses(c *gin.Context) {
	var filters ports.BusinessesFilter
	// .BindQuery() maps query parameters (e.g., ?active=true) to the struct
	if err := c.BindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
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

// @Summary      Update a business
// @Description  Updates a business's details by its ID.
// @Tags         Businesses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Business ID"
// @Param        business body ports.UpdateBusinessInput true "Business Update Data"
// @Success      200  {object}  ports.BusinessResponse
// @Failure      404  {object}  map[string]string "Business not found"
// @Router       /businesses/{id} [put]
func (h *BusinessHandler) UpdateBusiness(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	var input ports.UpdateBusinessInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	business, err := h.businessService.UpdateBusiness(c.Request.Context(), uint(id), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusOK, ports.MapBusinessToResponse(business))
}

// @Summary      Delete a business
// @Description  Deletes a business by its ID.
// @Tags         Businesses
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Business ID"
// @Success      204 "No Content"
// @Failure      404  {object}  map[string]string "Business not found"
// @Failure      409  {object}  map[string]string "Business is in use and cannot be deleted"
// @Router       /businesses/{id} [delete]
func (h *BusinessHandler) DeleteBusiness(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business ID"})
		return
	}

	if err := h.businessService.DeleteBusiness(c.Request.Context(), uint(id)); err != nil {
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
