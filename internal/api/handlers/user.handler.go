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

type UserHandler struct {
	userService *services.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// @Summary      Create a new user
// @Description  Creates a new user account with the provided details.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body ports.UserCreationSchema true "User Creation Data"
// @Success      201  {object}  ports.UserResponse
// @Failure      400  {object}  map[string]string "Validation error"
// @Failure      409  {object}  map[string]string "User with email already exists"
// @Failure      422  {object}  map[string]string "Password complexity error"
// @Router       /users [post]
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

// @Summary      Get a user by ID
// @Description  Retrieves the details of a single user by their unique ID.
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  ports.UserResponse
// @Failure      400  {object}  map[string]string "Invalid user ID format"
// @Failure      404  {object}  map[string]string "User not found"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.FindUserByID(c.Request.Context(), uint(id))
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

// @Summary      Get all users
// @Description  Retrieves a list of all user accounts.
// @Tags         Users
// @Produce      json
// @Success      200  {array}  ports.UserResponse
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
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

// @Summary      Update a user
// @Description  Updates a user's details by their ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        user body ports.UserUpdateSchema true "User Update Data"
// @Success      200  {object}  ports.UserResponse
// @Failure      400  {object}  map[string]string "Invalid request body or ID"
// @Failure      404  {object}  map[string]string "User not found"
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input ports.UserUpdateSchema
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), uint(id), input)
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

// @Summary      Delete a user
// @Description  Deletes a user by their ID.
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string "User not found"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), uint(id))
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
