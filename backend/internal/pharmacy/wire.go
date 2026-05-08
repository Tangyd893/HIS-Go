package pharmacy

import (
	"github.com/google/wire"

	"his-go/internal/pharmacy/handler"
	"his-go/internal/pharmacy/repository"
	"his-go/internal/pharmacy/service"
)

// ProviderSet 药房管理服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewPharmacyRepository,
	service.NewPharmacyService,
	handler.NewPharmacyHandler,
)
