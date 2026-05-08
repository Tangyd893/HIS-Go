package registration

import (
	"github.com/google/wire"

	"his-go/internal/registration/handler"
	"his-go/internal/registration/repository"
	"his-go/internal/registration/service"
)

// ProviderSet 挂号服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewRegistrationRepository,
	service.NewRegistrationService,
	handler.NewRegistrationHandler,
)
