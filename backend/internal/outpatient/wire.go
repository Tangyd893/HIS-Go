package outpatient

import (
	"github.com/google/wire"

	"his-go/internal/outpatient/handler"
	"his-go/internal/outpatient/repository"
	"his-go/internal/outpatient/service"
)

// ProviderSet 院外患者服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewOutpatientRepository,
	service.NewOutpatientService,
	handler.NewOutpatientHandler,
)
