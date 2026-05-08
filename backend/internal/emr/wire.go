package emr

import (
	"github.com/google/wire"

	"his-go/internal/emr/handler"
	"his-go/internal/emr/repository"
	"his-go/internal/emr/service"
)

// ProviderSet 电子病历服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewEMRRepository,
	service.NewEMRService,
	handler.NewEMRHandler,
)
