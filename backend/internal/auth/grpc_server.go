package auth

import (
	"context"

	"his-go/api/proto/auth"
	"his-go/api/proto/common"
	authsvc "his-go/internal/auth/service"
)

// AuthGrpcServer gRPC 认证服务实现
type AuthGrpcServer struct {
	auth.UnimplementedAuthServiceServer
	svc *authsvc.AuthService
}

// NewAuthGrpcServer 创建 gRPC 认证服务
func NewAuthGrpcServer(svc *authsvc.AuthService) *AuthGrpcServer {
	return &AuthGrpcServer{svc: svc}
}

// Login 用户登录
func (s *AuthGrpcServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	result, err := s.svc.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &auth.LoginResponse{
		Token:        result.Token,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		UserInfo: &auth.UserInfo{
			UserId:   result.UserInfo.UserID,
			Username: result.UserInfo.Username,
			RealName: result.UserInfo.RealName,
			Avatar:   result.UserInfo.Avatar,
			Role:     result.UserInfo.Role,
			DeptId:   result.UserInfo.DeptID,
			Perms:    result.UserInfo.Perms,
		},
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthGrpcServer) RefreshToken(ctx context.Context, req *common.IdRequest) (*auth.LoginResponse, error) {
	result, err := s.svc.RefreshToken(req.Id)
	if err != nil {
		return nil, err
	}
	return &auth.LoginResponse{
		Token:        result.Token,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		UserInfo: &auth.UserInfo{
			UserId:   result.UserInfo.UserID,
			Username: result.UserInfo.Username,
			RealName: result.UserInfo.RealName,
			Avatar:   result.UserInfo.Avatar,
			Role:     result.UserInfo.Role,
			DeptId:   result.UserInfo.DeptID,
			Perms:    result.UserInfo.Perms,
		},
	}, nil
}

// Logout 用户登出
func (s *AuthGrpcServer) Logout(ctx context.Context, req *common.IdRequest) (*common.BaseResponse, error) {
	if err := s.svc.Logout(req.Id); err != nil {
		return nil, err
	}
	return &common.BaseResponse{Code: 0, Message: "登出成功"}, nil
}

// ValidateToken 验证令牌
func (s *AuthGrpcServer) ValidateToken(ctx context.Context, req *common.IdRequest) (*common.BaseResponse, error) {
	_, err := s.svc.ValidateToken(req.Id)
	if err != nil {
		return &common.BaseResponse{Code: 401, Message: "令牌无效"}, nil
	}
	return &common.BaseResponse{Code: 0, Message: "令牌有效"}, nil
}
