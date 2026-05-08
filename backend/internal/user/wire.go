package user

import (
	"github.com/google/wire"

	"his-go/internal/user/handler"
	"his-go/internal/user/repository"
	"his-go/internal/user/service"
)

// ProviderSet 用户管理模块 wire 依赖注入集合
var ProviderSet = wire.NewSet(
	repository.NewPatientRepository,
	repository.NewEmployeeRepository,
	repository.NewDepartmentRepository,
	service.NewUserService,
	handler.NewUserHandler,
)
