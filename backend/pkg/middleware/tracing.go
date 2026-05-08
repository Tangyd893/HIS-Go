// Package middleware 链路追踪中间件（OpenTelemetry 占位）
package middleware

import (
	"github.com/gin-gonic/gin"
)

// Tracing 链路追踪中间件
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, exists := c.Get("requestID")
		if requestID, ok := requestID.(string); exists && ok {
			c.Header("X-Trace-ID", requestID)
		}
		c.Next()
	}
}
