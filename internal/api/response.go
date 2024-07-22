package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Error().Msg(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
