package notification

import (
	"github.com/google/wire"

	"his-go/internal/notification/handler"
	"his-go/internal/notification/repository"
	"his-go/internal/notification/service"
)

// ProviderSet 消息通知服务依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewNotificationRepository,
	service.NewNotificationService,
	handler.NewNotificationHandler,
)
