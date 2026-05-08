package schedule

import (
	"github.com/google/wire"

	"his-go/internal/schedule/handler"
	"his-go/internal/schedule/repository"
	"his-go/internal/schedule/service"
)

// ProviderSet 排班管理服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewScheduleRepository,
	service.NewScheduleService,
	handler.NewScheduleHandler,
)
