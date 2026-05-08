package service

import (
	"his-go/internal/system/model"
	"his-go/internal/system/repository"
)

// SystemService 系统管理业务服务
type SystemService struct {
	repo *repository.SystemRepository
}

// NewSystemService 创建系统管理业务服务
func NewSystemService(repo *repository.SystemRepository) *SystemService {
	return &SystemService{repo: repo}
}

// ListDictTypes 查询字典类型
func (s *SystemService) ListDictTypes() ([]model.DictType, error) {
	return s.repo.ListDictTypes()
}

// ListDictItems 查询字典项
func (s *SystemService) ListDictItems(dictType string) ([]model.DictItem, error) {
	return s.repo.ListDictItems(dictType)
}

// CreateDictItem 创建字典项
func (s *SystemService) CreateDictItem(item *model.DictItem) error {
	item.Status = 1
	return s.repo.CreateDictItem(item)
}

// UpdateDictItem 更新字典项
func (s *SystemService) UpdateDictItem(item *model.DictItem) error {
	return s.repo.UpdateDictItem(item)
}

// ListParams 查询系统参数
func (s *SystemService) ListParams() ([]model.SystemParam, error) {
	return s.repo.ListParams()
}

// UpdateParam 更新系统参数
func (s *SystemService) UpdateParam(param *model.SystemParam) error {
	return s.repo.UpdateParam(param)
}

// CreateOperationLog 创建操作日志
func (s *SystemService) CreateOperationLog(log *model.OperationLog) error {
	return s.repo.CreateOperationLog(log)
}

// ListOperationLogs 查询操作日志
func (s *SystemService) ListOperationLogs(userID, module string, page, pageSize int) ([]model.OperationLog, int64, error) {
	return s.repo.ListOperationLogs(userID, module, page, pageSize)
}
