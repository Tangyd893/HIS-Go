package examination

import (
	"github.com/google/wire"

	"his-go/internal/examination/handler"
	"his-go/internal/examination/repository"
	"his-go/internal/examination/service"
)

// ProviderSet 检查检验服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewExaminationRepository,
	service.NewExaminationService,
	handler.NewExaminationHandler,
)
