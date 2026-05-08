// Package errors 全局错误处理
package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError 创建应用错误
func NewAppError(code int, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: GetMessage(code),
		Detail:  detail,
	}
}

// GinRecoveryHandler 全局 Panic 恢复处理
func GinRecoveryHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    CodeInternalError,
			"message": GetMessage(CodeInternalError),
			"detail":  "服务器内部错误",
		})
	})
}

// GlobalErrorMiddleware 全局错误处理中间件
func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			if appErr, ok := err.Err.(*AppError); ok {
				c.JSON(http.StatusOK, appErr)
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    CodeInternalError,
				"message": GetMessage(CodeInternalError),
			})
		}
	}
}
