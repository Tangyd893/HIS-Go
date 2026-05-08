package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"his-go/internal/auth/model"
	"his-go/internal/auth/repository"
	"his-go/pkg/redis"
	"his-go/pkg/security/jwt"
)

// LoginResult 登录结果
type LoginResult struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refreshToken"`
	ExpiresIn    int64    `json:"expiresIn"`
	UserInfo     UserInfo `json:"userInfo"`
}

// UserInfo 用户信息
type UserInfo struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	RealName string   `json:"realName"`
	Avatar   string   `json:"avatar"`
	Role     string   `json:"role"`
	DeptID   string   `json:"deptId"`
	Perms    []string `json:"perms"`
}

// AuthService 认证服务
type AuthService struct {
	repo   *repository.AuthRepository
	jwtSvc *jwt.JWTService
	rdb    *redis.Client
}

// NewAuthService 创建认证服务
func NewAuthService(repo *repository.AuthRepository, jwtSvc *jwt.JWTService, rdb *redis.Client) *AuthService {
	return &AuthService{
		repo:   repo,
		jwtSvc: jwtSvc,
		rdb:    rdb,
	}
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*LoginResult, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	perms, err := s.getUserPerms(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	permsStr := make([]string, len(perms))
	for i, p := range perms {
		permsStr[i] = p.PermCode
	}

	claims := &jwt.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		DeptID:   user.DeptID,
		Perms:    permsStr,
	}

	token, err := s.jwtSvc.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtSvc.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tokenKey := fmt.Sprintf("auth:token:%s", user.ID)
	if err := s.rdb.Set(ctx, tokenKey, token, 24*time.Hour); err != nil {
		return nil, err
	}

	_ = s.repo.UpdateLastLogin(user.ID)

	userInfo := UserInfo{
		UserID:   user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Role:     user.Role,
		DeptID:   user.DeptID,
		Perms:    permsStr,
	}

	return &LoginResult{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(24 * time.Hour.Seconds()),
		UserInfo:     userInfo,
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(userID string) error {
	ctx := context.Background()
	tokenKey := fmt.Sprintf("auth:token:%s", userID)
	return s.rdb.Del(ctx, tokenKey)
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(tokenString string) (*LoginResult, error) {
	claims, err := s.jwtSvc.ParseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("令牌无效或已过期")
	}

	ctx := context.Background()
	tokenKey := fmt.Sprintf("auth:token:%s", claims.UserID)
	_, err = s.rdb.Get(ctx, tokenKey)
	if err != nil {
		return nil, fmt.Errorf("令牌已失效，请重新登录")
	}

	newToken, err := s.jwtSvc.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(ctx, tokenKey, newToken, 24*time.Hour); err != nil {
		return nil, err
	}

	user, err := s.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	userInfo := UserInfo{
		UserID:   user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Role:     user.Role,
		DeptID:   user.DeptID,
		Perms:    claims.Perms,
	}

	return &LoginResult{
		Token:        newToken,
		RefreshToken: newToken,
		ExpiresIn:    int64(24 * time.Hour.Seconds()),
		UserInfo:     userInfo,
	}, nil
}

// ValidateToken 验证令牌
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	claims, err := s.jwtSvc.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tokenKey := fmt.Sprintf("auth:token:%s", claims.UserID)
	storedToken, err := s.rdb.Get(ctx, tokenKey)
	if err != nil || storedToken != tokenString {
		return nil, fmt.Errorf("令牌已失效")
	}

	return claims, nil
}

func (s *AuthService) getUserPerms(userID, role string) ([]model.Permission, error) {
	roles, err := s.repo.FindRolesByUserID(userID)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.ID
	}

	return s.repo.FindPermsByRoleIDs(roleIDs)
}
