package middlewares

import (
	"github.com/Oxyrus/shopping/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware provides an authorization mechanism that forces
// clients to send a valid JWT in the Authorization header for all
// the requests that hit routes that use this middleware.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.ValidateToken(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
