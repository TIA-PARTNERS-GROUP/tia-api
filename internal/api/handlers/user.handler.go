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

type UserHandler struct {
	userService *services.UserService
	validate    *validator.Validate
	routes      *constants.Routes
}

func NewUserHandler(userService *services.UserService, routes *constants.Routes) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
		routes:      routes,
	}
}

// @Summary Register New User
// @Description Registers a new user account. Does not require prior authentication.
// @Tags users, auth
// @Accept json
// @Produce json
// @Param user body ports.UserCreationSchema true "User registration details (Name, Email, Password)"
// @Success 201 {object} ports.UserResponse "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 409 {object} map[string]interface{} "ErrUserAlreadyExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var input ports.UserCreationSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userService.CreateUser(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapUserToResponse(user))
}

// @Summary Get User by ID
// @Description Retrieves a user's profile by their unique ID. Requires authentication.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} ports.UserResponse "User retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrUserNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := h.userService.FindUserByID(c.Request.Context(), uint(id))
	if err != nil {

		if errors.Is(err, ports.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapUserToResponse(user))
}

// @Summary List All Users
// @Description Retrieves a list of all user profiles. Requires authentication.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ports.UserResponse "List of users"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {

	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	users, err := h.userService.FindAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	userResponses := make([]ports.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ports.MapUserToResponse(&user)
	}
	c.JSON(http.StatusOK, userResponses)
}

// @Summary Update User Profile
// @Description Updates the authenticated user's profile information. Requires self-management.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID (must match authenticated user)"
// @Param update body ports.UserUpdateSchema true "Fields to update (e.g., FirstName, ContactEmail)"
// @Success 200 {object} ports.UserResponse "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID or request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Cannot update another user's profile"
// @Failure 404 {object} map[string]interface{} "ErrUserNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {

	idStr := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	if authUserID != uint(targetUserID) {

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only update your own profile"})
		return
	}

	var input ports.UserUpdateSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), uint(targetUserID), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapUserToResponse(user))
}

// @Summary Delete User Account
// @Description Deletes the authenticated user's account. Requires self-management.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID (must match authenticated user)"
// @Success 204 "Account deleted successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Cannot delete another user's profile"
// @Failure 404 {object} map[string]interface{} "ErrUserNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {

	idStr := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	if authUserID != uint(targetUserID) {

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only delete your own profile"})
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), uint(targetUserID))
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
