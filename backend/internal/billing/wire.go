package billing

import (
	"github.com/google/wire"

	"his-go/internal/billing/handler"
	"his-go/internal/billing/repository"
	"his-go/internal/billing/service"
)

// ProviderSet 收费结算服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewBillingRepository,
	service.NewBillingService,
	handler.NewBillingHandler,
)
