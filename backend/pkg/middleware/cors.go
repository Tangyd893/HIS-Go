// Package middleware 通用 HTTP 中间件
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Cors 跨域处理中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RequestID 请求ID中间件（链路追踪用）
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if query != "" {
			path = path + "?" + query
		}

		gin.DefaultWriter.Write([]byte(
			"[GIN] " + c.Request.Method + " | " +
				formatStatus(statusCode) + " | " +
				latency.String() + " | " +
				c.ClientIP() + " | " +
				path + "\n",
		))
	}
}

func formatStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green(code)
	case code >= 300 && code < 400:
		return blue(code)
	case code >= 400 && code < 500:
		return yellow(code)
	default:
		return red(code)
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// 终端颜色（简化）
func green(code int) string  { return "\033[32m" + string(rune(code)) + "\033[0m" }
func blue(code int) string   { return "\033[34m" + string(rune(code)) + "\033[0m" }
func yellow(code int) string { return "\033[33m" + string(rune(code)) + "\033[0m" }
func red(code int) string    { return "\033[31m" + string(rune(code)) + "\033[0m" }
