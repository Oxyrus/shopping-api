package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthController provides an endpoint to verify the service status.
type HealthController struct{}

// GetHealthStatus provides a handler that will return a status code 200
// if the service is operational.
func (c *HealthController) GetHealthStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Healthy"})
}
