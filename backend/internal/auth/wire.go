package auth

import (
	"github.com/google/wire"

	"his-go/internal/auth/handler"
	"his-go/internal/auth/repository"
	"his-go/internal/auth/service"
)

// ProviderSet 认证模块 wire 依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewAuthRepository,
	service.NewAuthService,
	handler.NewAuthHandler,
)
