package middleware

import (
	"os"

	"github.com/gin-gonic/gin"

	"his-go/pkg/config"
	secauth "his-go/pkg/security/auth"
	"his-go/pkg/security/jwt"
)

// InitJWT 从配置初始化 JWT 服务（与各微服务网关鉴权共用密钥）
func InitJWT(cfg *config.Config) *jwt.JWTService {
	if os.Getenv("USE_JWT_SIMPLE") == "true" || cfg.JWT.PrivateKey == "" {
		secret := cfg.JWT.PrivateKey
		if secret == "" {
			secret = "his-go-default-secret"
		}
		return jwt.NewSimpleJWTService(secret, cfg.JWT.ExpireHour)
	}
	svc, err := jwt.NewJWTService(cfg.JWT.PrivateKey, cfg.JWT.PublicKey, cfg.JWT.ExpireHour)
	if err != nil {
		secret := cfg.JWT.PrivateKey
		if secret == "" {
			secret = "his-go-default-secret"
		}
		return jwt.NewSimpleJWTService(secret, cfg.JWT.ExpireHour)
	}
	return svc
}

// UserContext 解析网关 X-User-* 头或 JWT（供各微服务 API 组使用）
func UserContext(jwtSvc *jwt.JWTService) gin.HandlerFunc {
	return secauth.UserContextFromGatewayOrJWT(jwtSvc)
}
