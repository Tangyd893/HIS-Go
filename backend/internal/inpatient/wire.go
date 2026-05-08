package inpatient

import (
	"github.com/google/wire"

	"his-go/internal/inpatient/handler"
	"his-go/internal/inpatient/repository"
	"his-go/internal/inpatient/service"
)

// ProviderSet 住院管理服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewInpatientRepository,
	service.NewInpatientService,
	handler.NewInpatientHandler,
)
