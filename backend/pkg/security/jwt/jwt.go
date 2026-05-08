// Package jwt JWT 令牌签发与解析
package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"his-go/pkg/logger"
)

var (
	ErrTokenInvalid = errors.New("token 无效")
	ErrTokenExpired = errors.New("token 已过期")
)

// Claims 自定义 JWT 声明
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Role     string   `json:"role"`
	DeptID   string   `json:"dept_id"`
	Perms    []string `json:"perms"`
	jwtlib.RegisteredClaims
}

// JWTService JWT 服务
type JWTService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	hmacSecret []byte
	expireHour int
	useHMAC    bool
}

// NewJWTService 创建 JWT 服务（RS256 非对称加密）
func NewJWTService(privateKeyPEM, publicKeyPEM string, expireHour int) (*JWTService, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}
	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return nil, err
	}
	return &JWTService{
		privateKey: privateKey,
		publicKey:  publicKey,
		expireHour: expireHour,
		useHMAC:    false,
	}, nil
}

// NewSimpleJWTService 创建简化版 JWT 服务（HS256 对称加密，仅限开发环境）
func NewSimpleJWTService(secret string, expireHour int) *JWTService {
	if secret == "" {
		secret = "his-go-default-secret"
		logger.Warn("JWT 使用默认 HS256 密钥，严禁用于生产环境")
	}
	return &JWTService{
		hmacSecret: []byte(secret),
		expireHour: expireHour,
		useHMAC:    true,
	}
}

// NewVerifyOnlyJWTService 创建仅用于验证的 JWT 服务（RS256 — 网关专用）
func NewVerifyOnlyJWTService(publicKeyPEM string) (*JWTService, error) {
	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return nil, err
	}
	return &JWTService{
		publicKey:  publicKey,
		expireHour: 2,
		useHMAC:    false,
	}, nil
}

// GenerateToken 生成 JWT Token
func (s *JWTService) GenerateToken(claims *Claims) (string, error) {
	claims.RegisteredClaims = jwtlib.RegisteredClaims{
		Issuer:    "his-go",
		Subject:   claims.UserID,
		ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Duration(s.expireHour) * time.Hour)),
		IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		NotBefore: jwtlib.NewNumericDate(time.Now()),
	}

	if s.useHMAC {
		token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
		return token.SignedString(s.hmacSecret)
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

// ParseToken 解析 JWT Token
func (s *JWTService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(token *jwtlib.Token) (interface{}, error) {
		if s.useHMAC {
			if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, ErrTokenInvalid
			}
			return s.hmacSecret, nil
		}
		if _, ok := token.Method.(*jwtlib.SigningMethodRSA); !ok {
			return nil, ErrTokenInvalid
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新 Token
func (s *JWTService) RefreshToken(claims *Claims) (string, error) {
	claims.RegisteredClaims.ExpiresAt = jwtlib.NewNumericDate(time.Now().Add(time.Duration(s.expireHour) * time.Hour))
	return s.GenerateToken(claims)
}

func parsePrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("解析私钥失败")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("私钥类型不是 RSA")
	}
	return rsaKey, nil
}

func parsePublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("解析公钥失败")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("公钥类型不是 RSA")
	}
	return rsaKey, nil
}
