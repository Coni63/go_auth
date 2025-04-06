package middlewares

import (
	"fmt"

	"go_auth/initializers"

	"github.com/gin-gonic/gin"
)

// Prometheus middleware for Gin
func PrometheusStatusCodeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		initializers.CountRequests.Inc() // Increment the counter for each request

		// Process request
		c.Next()

		// After request is handled, get the status code
		statusCode := c.Writer.Status()
		initializers.CountStatusCodes.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
	}
}
