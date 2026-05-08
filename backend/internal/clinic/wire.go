package clinic

import (
	"github.com/google/wire"

	"his-go/internal/clinic/handler"
	"his-go/internal/clinic/repository"
	"his-go/internal/clinic/service"
)

// ProviderSet 门诊诊疗服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewClinicRepository,
	service.NewClinicService,
	handler.NewClinicHandler,
)
