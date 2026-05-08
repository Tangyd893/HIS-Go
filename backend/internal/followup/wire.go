package followup

import (
	"github.com/google/wire"

	"his-go/internal/followup/handler"
	"his-go/internal/followup/repository"
	"his-go/internal/followup/service"
)

// ProviderSet 随访管理服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewFollowupRepository,
	service.NewFollowupService,
	handler.NewFollowupHandler,
)
