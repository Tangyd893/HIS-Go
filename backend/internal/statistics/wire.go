package statistics

import (
	"github.com/google/wire"

	"his-go/internal/statistics/handler"
	"his-go/internal/statistics/repository"
	"his-go/internal/statistics/service"
)

// ProviderSet 数据统计服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewStatisticsRepository,
	service.NewStatisticsService,
	handler.NewStatisticsHandler,
)
