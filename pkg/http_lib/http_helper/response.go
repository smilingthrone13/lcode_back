package http_helper

import (
	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Message string `json:"message"`
}

type errsResponse struct {
	Errors []error `json:"errors"`
}

func NewErrorResponse(c *gin.Context, statusCode int, msg string) {
	c.AbortWithStatusJSON(statusCode, errResponse{msg})
}

func NewErrorsResponse(c *gin.Context, statusCode int, errors []error) {
	c.AbortWithStatusJSON(statusCode, errsResponse{Errors: errors})
}
