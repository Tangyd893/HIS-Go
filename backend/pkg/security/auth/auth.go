// Package auth JWT 鉴权中间件
package auth

import (
	"strings"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	apperrors "his-go/pkg/errors"
	"his-go/pkg/logger"
	"his-go/pkg/response"
	"his-go/pkg/security/jwt"
)

// UserContext 用户上下文（放在 Gin Context 中）
type UserContext struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Role     string   `json:"role"`
	DeptID   string   `json:"dept_id"`
	Perms    []string `json:"perms"`
}

const UserContextKey = "user_context"

// JwtAuth JWT 认证中间件
func JwtAuth(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, apperrors.CodeUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Fail(c, apperrors.CodeUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwtService.ParseToken(parts[1])
		if err != nil {
			logger.Warn("JWT 解析失败", zap.Error(err))
			response.Fail(c, apperrors.CodeUnauthorized)
			c.Abort()
			return
		}

		userCtx := &UserContext{
			UserID:   claims.UserID,
			Username: claims.Username,
			Role:     claims.Role,
			DeptID:   claims.DeptID,
			Perms:    claims.Perms,
		}
		c.Set(UserContextKey, userCtx)
		c.Next()
	}
}

// GatewayUserContext 从网关转发的 X-User-* 请求头解析用户上下文
func GatewayUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			response.Fail(c, apperrors.CodeUnauthorized)
			c.Abort()
			return
		}

		permsHeader := c.GetHeader("X-User-Perms")
		var perms []string
		if permsHeader != "" {
			perms = strings.Split(permsHeader, ",")
		}

		c.Set(UserContextKey, &UserContext{
			UserID:   userID,
			Username: c.GetHeader("X-Username"),
			Role:     c.GetHeader("X-User-Role"),
			DeptID:   c.GetHeader("X-User-Dept"),
			Perms:    perms,
		})
		c.Next()
	}
}

// UserContextFromGatewayOrJWT 优先使用网关头，否则回退 JWT 解析（直连调试）
func UserContextFromGatewayOrJWT(jwtService *jwt.JWTService) gin.HandlerFunc {
	gateway := GatewayUserContext()
	jwtAuth := JwtAuth(jwtService)
	return func(c *gin.Context) {
		if c.GetHeader("X-User-ID") != "" {
			gateway(c)
			return
		}
		jwtAuth(c)
	}
}

// ResolveUserID 将 "current" 解析为当前登录用户 ID
func ResolveUserID(c *gin.Context, id string) string {
	if id != "" && id != "current" {
		return id
	}
	if uc := GetUserContext(c); uc != nil && uc.UserID != "" {
		return uc.UserID
	}
	return "demo-admin"
}

// GetUserContext 从 Gin Context 中获取用户上下文
func GetUserContext(c *gin.Context) *UserContext {
	if val, ok := c.Get(UserContextKey); ok {
		if userCtx, ok := val.(*UserContext); ok {
			return userCtx
		}
	}
	return nil
}

// RequireRole 角色权限检查中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx := GetUserContext(c)
		if userCtx == nil {
			response.Fail(c, apperrors.CodeUnauthorized)
			c.Abort()
			return
		}

		for _, role := range roles {
			if userCtx.Role == role {
				c.Next()
				return
			}
		}

		response.Fail(c, apperrors.CodeForbidden)
		c.Abort()
	}
}

// RequirePerm 权限检查中间件
func RequirePerm(perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx := GetUserContext(c)
		if userCtx == nil {
			response.Fail(c, apperrors.CodeUnauthorized)
			c.Abort()
			return
		}

		permSet := make(map[string]bool)
		for _, p := range userCtx.Perms {
			permSet[p] = true
		}

		for _, p := range perms {
			if permSet[p] {
				c.Next()
				return
			}
		}

		response.Fail(c, apperrors.CodeForbidden)
		c.Abort()
	}
}
