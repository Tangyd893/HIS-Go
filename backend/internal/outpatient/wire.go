package outpatient

import (
	"github.com/google/wire"

	"his-go/internal/outpatient/assistant"
	"his-go/internal/outpatient/handler"
	"his-go/internal/outpatient/repository"
	"his-go/internal/outpatient/service"
)

// ProviderSet 院外患者服务依赖注入集合
// 注：assistant.Service 需在 main.go 中手动创建（依赖环境变量），不通过 wire
var ProviderSet = wire.NewSet(
	repository.NewOutpatientRepository,
	service.NewOutpatientService,
	handler.NewOutpatientHandler,
)

// ProvideAssistantService 为 wire 提供就诊助手服务（环境变量驱动）
func ProvideAssistantService() *assistant.Service {
	cfg := assistant.LoadConfig()
	return assistant.NewService(cfg)
}
