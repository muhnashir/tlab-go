package handler

import (
	"wallet-api/internal/domain"
	"wallet-api/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Service domain.UserService
}

func NewAuthHandler(s domain.UserService) *AuthHandler {
	return &AuthHandler{Service: s}
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} utils.ApiResponse
// @Failure 409 {object} utils.ApiResponse
// @Failure 500 {object} utils.ApiResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validation could be added here (e.g. using go-playground/validator)

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.Service.Register(c.Context(), &user); err != nil {
		// Differentiate error types (e.g. email exists vs internal server error)
		if err.Error() == "email already registered" {
			return utils.Error(c, fiber.StatusConflict, err.Error(), nil)
		}
		// Log the actual error for internal server errors via utils helper
		return utils.InternalServerError(c, "Failed to register user", err.Error())
	}

	return utils.Created(c, "User registered successfully", nil)
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to receive JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.ApiResponse
// @Failure 401 {object} utils.ApiResponse
// @Failure 500 {object} utils.ApiResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", err.Error())
	}

	token, err := h.Service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return utils.Unauthorized(c, "Invalid credentials")
		}
		return utils.InternalServerError(c, "Failed to login", err.Error())
	}

	return utils.Success(c, fiber.StatusOK, "Login successful", fiber.Map{
		"token": token,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile of the logged-in user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.ApiResponse{data=domain.User}
// @Failure 401 {object} utils.ApiResponse
// @Failure 404 {object} utils.ApiResponse
// @Failure 500 {object} utils.ApiResponse
// @Router /users/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Parse user_id from middleware
	userIDf, ok := c.Locals("user_id").(int64)
	// Wait, in middleware we cast to int64. Let's verify middleware.
	// middleware: c.Locals("user_id", int64(userIDFloat))
	// yes.

	if !ok {
		return utils.Unauthorized(c, "Invalid user session")
	}

	user, err := h.Service.GetProfile(c.Context(), userIDf)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch profile", err.Error())
	}
	if user == nil {
		return utils.NotFound(c, "User not found")
	}

	// Make sure password is not returned
	user.Password = ""

	return utils.Success(c, fiber.StatusOK, "User profile retrieved", user)
}
