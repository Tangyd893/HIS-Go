package handler

import (
	"github.com/gin-gonic/gin"

	"his-go/internal/auth/service"
	apperrors "his-go/pkg/errors"
	"his-go/pkg/response"
	"his-go/pkg/security/auth"
)

// AuthHandler 认证接口处理器
type AuthHandler struct {
	svc *service.AuthService
}

// NewAuthHandler 创建认证接口处理器
func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录接口
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, apperrors.CodeParamInvalid, "请输入用户名和密码")
		return
	}

	result, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		response.FailWithMsg(c, apperrors.CodeUserPasswordWrong, err.Error())
		return
	}

	response.Success(c, result)
}

// Logout 用户登出接口
func (h *AuthHandler) Logout(c *gin.Context) {
	userCtx := auth.GetUserContext(c)
	if userCtx == nil {
		response.Fail(c, apperrors.CodeUnauthorized)
		return
	}

	if err := h.svc.Logout(userCtx.UserID); err != nil {
		response.Fail(c, apperrors.CodeInternalError)
		return
	}

	response.Success(c, nil)
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	Token string `json:"token" binding:"required"`
}

// RefreshToken 刷新令牌接口
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, apperrors.CodeParamInvalid, "请提供令牌")
		return
	}

	result, err := h.svc.RefreshToken(req.Token)
	if err != nil {
		response.FailWithMsg(c, apperrors.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, result)
}
