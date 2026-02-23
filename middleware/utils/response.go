package utils

import (
	"middleware/models"

	"github.com/gofiber/fiber/v2"
)

// SuccessResponse sends a standard success response
func SuccessResponse(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"data":    data,
		"message": message,
	})
}

// SuccessListResponse sends a paginated list response
func SuccessListResponse(c *fiber.Ctx, data interface{}, total int64, page, limit int) error {
	totalPages := CalculateTotalPages(total, limit)

	response := models.PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "List retrieved successfully",
	})
}

// ErrorResponse sends a standard error response
func ErrorResponse(c *fiber.Ctx, statusCode int, code string, message string) error {
	response := models.ErrorResponse{
		Code:    code,
		Message: message,
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"error":   response,
	})
}

// BadRequest sends a 400 error
func BadRequest(c *fiber.Ctx, code string, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, code, message)
}

// Unauthorized sends a 401 error
func Unauthorized(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, "auth_failed", message)
}

// Forbidden sends a 403 error
func Forbidden(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, "access_denied", message)
}

// NotFound sends a 404 error
func NotFound(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, "not_found", message)
}

// InternalError sends a 500 error
func InternalError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, "internal_error", message)
}

// Conflict sends a 409 error
func Conflict(c *fiber.Ctx, code string, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, code, message)
}
