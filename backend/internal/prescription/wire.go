package prescription

import (
	"github.com/google/wire"

	"his-go/internal/prescription/handler"
	"his-go/internal/prescription/repository"
	"his-go/internal/prescription/service"
)

// ProviderSet 处方管理服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewPrescriptionRepository,
	service.NewPrescriptionService,
	handler.NewPrescriptionHandler,
)
