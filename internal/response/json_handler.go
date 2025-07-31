package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// JSONHandler handles JSON responses for HTTP API
type JSONHandler struct{}

// NewJSONHandler creates a new JSON response handler
func NewJSONHandler() *JSONHandler {
	return &JSONHandler{}
}

// Success sends a successful JSON response
func (h *JSONHandler) Success(c *gin.Context, data interface{}, message ...string) {
	response := BaseResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(http.StatusOK, response)
}

// Error sends an error JSON response
func (h *JSONHandler) Error(c *gin.Context, statusCode int, code, message string, details ...string) {
	response := BaseResponse{
		Success:   false,
		Timestamp: time.Now(),
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}

	if len(details) > 0 {
		response.Error.Details = details[0]
	}

	c.JSON(statusCode, response)
}

// BadRequest sends a 400 Bad Request response
func (h *JSONHandler) BadRequest(c *gin.Context, message string, details ...string) {
	h.Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// Unauthorized sends a 401 Unauthorized response
func (h *JSONHandler) Unauthorized(c *gin.Context, message string, details ...string) {
	h.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message, details...)
}

// Forbidden sends a 403 Forbidden response
func (h *JSONHandler) Forbidden(c *gin.Context, message string, details ...string) {
	h.Error(c, http.StatusForbidden, "FORBIDDEN", message, details...)
}

// NotFound sends a 404 Not Found response
func (h *JSONHandler) NotFound(c *gin.Context, message string, details ...string) {
	h.Error(c, http.StatusNotFound, "NOT_FOUND", message, details...)
}

// InternalServerError sends a 500 Internal Server Error response
func (h *JSONHandler) InternalServerError(c *gin.Context, message string, details ...string) {
	h.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, details...)
}

// ServiceUnavailable sends a 503 Service Unavailable response
func (h *JSONHandler) ServiceUnavailable(c *gin.Context, message string, details ...string) {
	h.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message, details...)
}

// NewsResponse sends a news response
func (h *JSONHandler) NewsResponse(c *gin.Context, resp *NewsResponse) {
	if resp.Success {
		c.JSON(http.StatusOK, resp)
	} else {
		status := http.StatusInternalServerError
		if resp.Error != nil && resp.Error.Code == "NOT_FOUND" {
			status = http.StatusNotFound
		}
		c.JSON(status, resp)
	}
}

// SearchResponse sends a search response
func (h *JSONHandler) SearchResponse(c *gin.Context, resp *SearchResponse) {
	if resp.Success {
		c.JSON(http.StatusOK, resp)
	} else {
		status := http.StatusInternalServerError
		if resp.Error != nil && resp.Error.Code == "NOT_FOUND" {
			status = http.StatusNotFound
		}
		c.JSON(status, resp)
	}
}

// HealthResponse sends a health check response
func (h *JSONHandler) HealthResponse(c *gin.Context, resp *HealthResponse) {
	status := http.StatusOK
	if resp.Status != "healthy" {
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, resp)
}

// PrettyJSON sends a formatted JSON response (for debugging)
func (h *JSONHandler) PrettyJSON(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

// CORS sets CORS headers for API responses
func (h *JSONHandler) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestLogger logs API requests
func (h *JSONHandler) RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}

// RateLimitError sends a rate limit exceeded response
func (h *JSONHandler) RateLimitError(c *gin.Context) {
	h.Error(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests, please try again later")
}

// ValidationError sends a validation error response
func (h *JSONHandler) ValidationError(c *gin.Context, errors map[string]string) {
	response := BaseResponse{
		Success:   false,
		Message:   "Validation failed",
		Timestamp: time.Now(),
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "Request validation failed",
		},
		Data: gin.H{"validation_errors": errors},
	}

	c.JSON(http.StatusBadRequest, response)
}
