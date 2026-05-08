// Package middleware Panic 恢复中间件
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	apperrors "his-go/pkg/errors"
	"his-go/pkg/logger"
	"his-go/pkg/response"
)

// Recovery Panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic 恢复",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    apperrors.CodeInternalError,
					"message": apperrors.GetMessage(apperrors.CodeInternalError),
				})
			}
		}()
		c.Next()
	}
}

// AbortWithError 中断请求并返回错误
func AbortWithError(c *gin.Context, code int, detail string) {
	response.Fail(c, code)
	c.Abort()
}
