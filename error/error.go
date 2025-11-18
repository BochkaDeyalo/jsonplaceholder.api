package apierror

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Predefined errors
var (
	BadRequest = &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Invalid request",
	}

	ValidationFailed = &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Validation failed",
	}

	InternalServerError = &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}

	ExternalServiceError = &ErrorResponse{
		Code:    http.StatusBadGateway,
		Message: "External service error",
	}
)

func SendError(c *gin.Context, errResp *ErrorResponse, details string) {
	response := *errResp
	if details != "" {
		response.Details = details
	}
	c.JSON(response.Code, response)
}
