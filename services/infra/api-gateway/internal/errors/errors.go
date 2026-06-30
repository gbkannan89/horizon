package errors

import (
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error     APIError `json:"error"`
	RequestID string   `json:"request_id"`
}

func Respond(c *gin.Context, status int, errType, code, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Error:     APIError{Type: errType, Code: code, Message: message},
		RequestID: c.GetString("correlation_id"),
	})
}

func Validation(c *gin.Context, message string)    { Respond(c, 400, "VALIDATION_ERROR", "INVALID_INPUT", message) }
func NotFound(c *gin.Context, message string)       { Respond(c, 404, "NOT_FOUND", "RESOURCE_NOT_FOUND", message) }
func Conflict(c *gin.Context, message string)       { Respond(c, 409, "CONFLICT", "CONFLICT", message) }
func Internal(c *gin.Context, message string)       { Respond(c, 500, "INTERNAL_ERROR", "INTERNAL_ERROR", message) }

func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		Validation(c, "Invalid request body: "+err.Error())
		return false
	}
	return true
}
