package health_record

import (
	"github.com/google/wire"

	"his-go/internal/health_record/handler"
	"his-go/internal/health_record/repository"
	"his-go/internal/health_record/service"
)

// ProviderSet 健康档案服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewHealthRecordRepository,
	service.NewHealthRecordService,
	handler.NewHealthRecordHandler,
)
