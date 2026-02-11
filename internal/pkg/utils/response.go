package utils

import (
	"github.com/gofiber/fiber/v2"
)

// ApiResponse defines the standard response format for the API
type ApiResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success sends a successful JSON response with a standard format.
// The default status code is 200 (OK), but can be customized.
func Success(c *fiber.Ctx, code int, message string, data interface{}) error {
	if code == 0 {
		code = fiber.StatusOK
	}

	LogInfof("[SUCCESS] %s - Code: %d", message, code)

	return c.Status(code).JSON(ApiResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Error sends an error JSON response with a standard format.
// The default status code is 500 (Internal Server Error), but can be customized.
func Error(c *fiber.Ctx, code int, message string, errors interface{}) error {
	if code == 0 {
		code = fiber.StatusInternalServerError
	}

	if errors != nil {
		LogErrorf("[ERROR] %s - Code: %d, Details: %v", message, code, errors)
	} else {
		LogErrorf("[ERROR] %s - Code: %d", message, code)
	}

	return c.Status(code).JSON(ApiResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error:   errors,
	})
}

// Created sends a standardized 201 Created response
func Created(c *fiber.Ctx, message string, data interface{}) error {
	return Success(c, fiber.StatusCreated, message, data)
}

// BadRequest sends a standardized 400 Bad Request response
func BadRequest(c *fiber.Ctx, message string, errors interface{}) error {
	return Error(c, fiber.StatusBadRequest, message, errors)
}

// NotFound sends a standardized 404 Not Found response
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

// InternalServerError sends a standardized 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, message string, errors interface{}) error {
	return Error(c, fiber.StatusInternalServerError, message, errors)
}

// Forbidden sends a standardized 403 Forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message, nil)
}

// Unauthorized sends a standardized 401 Unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, nil)
}
