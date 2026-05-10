// Package response 统一响应封装
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页响应数据
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// codeToHTTPStatus 错误码到 HTTP 状态码的映射
func codeToHTTPStatus(code int) int {
	switch code {
	case errors.CodeSuccess:
		return http.StatusOK
	case errors.CodeParamInvalid:
		return http.StatusBadRequest
	case errors.CodeUnauthorized, errors.CodeAuthExpired:
		return http.StatusUnauthorized
	case errors.CodeForbidden:
		return http.StatusForbidden
	case errors.CodeNotFound, errors.CodeUserNotFound, errors.CodeAuthNotFound:
		return http.StatusNotFound
	case errors.CodeConflict, errors.CodeRegistrationRepeat, errors.CodeAuthDuplicate:
		return http.StatusConflict
	case errors.CodeTimeout:
		return http.StatusRequestTimeout
	case errors.CodeRateLimited:
		return http.StatusTooManyRequests
	case errors.CodeServiceUnavail:
		return http.StatusServiceUnavailable
	case errors.CodeInternalError, errors.CodeUnknownError:
		return http.StatusInternalServerError
	default:
		if code >= 20000 {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: "成功",
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: "成功",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Fail 失败响应（根据错误码自动映射 HTTP 状态码）
func Fail(c *gin.Context, code int) {
	c.JSON(codeToHTTPStatus(code), Response{
		Code:    code,
		Message: errors.GetMessage(code),
	})
}

// FailWithMsg 带自定义消息的失败响应（根据错误码自动映射 HTTP 状态码）
func FailWithMsg(c *gin.Context, code int, message string) {
	c.JSON(codeToHTTPStatus(code), Response{
		Code:    code,
		Message: message,
	})
}
