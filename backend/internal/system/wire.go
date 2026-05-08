package system

import (
	"github.com/google/wire"

	"his-go/internal/system/handler"
	"his-go/internal/system/repository"
	"his-go/internal/system/service"
)

// ProviderSet 系统管理服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewSystemRepository,
	service.NewSystemService,
	handler.NewSystemHandler,
)
