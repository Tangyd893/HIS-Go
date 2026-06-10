package repository

import (
	"gorm.io/gorm"

	"his-go/internal/system/model"
	"his-go/pkg/errors"
)

// SystemRepository 系统管理数据仓库
type SystemRepository struct {
	db *gorm.DB
}

// NewSystemRepository 创建系统管理数据仓库
func NewSystemRepository(db *gorm.DB) *SystemRepository {
	return &SystemRepository{db: db}
}

// ListDictTypes 查询字典类型列表
func (r *SystemRepository) ListDictTypes() ([]model.DictType, error) {
	var list []model.DictType
	if err := r.db.Where("status = 1").Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, errors.WrapQueryError("字典类型", err)
	}
	return list, nil
}

// ListDictItems 根据字典类型查询字典项
func (r *SystemRepository) ListDictItems(dictType string) ([]model.DictItem, error) {
	var list []model.DictItem
	if err := r.db.Where("dict_type = ? AND status = 1", dictType).
		Order("sort_order ASC").Find(&list).Error; err != nil {
		return nil, errors.WrapQueryError("字典项", err)
	}
	return list, nil
}

// CreateDictItem 创建字典项
func (r *SystemRepository) CreateDictItem(item *model.DictItem) error {
	if err := r.db.Create(item).Error; err != nil {
		return errors.WrapCreateError("字典项", err)
	}
	return nil
}

// UpdateDictItem 更新字典项
func (r *SystemRepository) UpdateDictItem(item *model.DictItem) error {
	result := r.db.Model(&model.DictItem{}).Where("id = ?", item.ID).Updates(item)
	if result.Error != nil {
		return errors.WrapUpdateError("字典项", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "字典项不存在")
	}
	return nil
}

// ListParams 查询系统参数列表
func (r *SystemRepository) ListParams() ([]model.SystemParam, error) {
	var list []model.SystemParam
	if err := r.db.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, errors.WrapQueryError("系统参数", err)
	}
	return list, nil
}

// UpdateParam 更新系统参数
func (r *SystemRepository) UpdateParam(param *model.SystemParam) error {
	result := r.db.Model(&model.SystemParam{}).Where("id = ?", param.ID).Updates(param)
	if result.Error != nil {
		return errors.WrapUpdateError("系统参数", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "系统参数不存在")
	}
	return nil
}

// CreateOperationLog 创建操作日志
func (r *SystemRepository) CreateOperationLog(log *model.OperationLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return errors.WrapCreateError("操作日志", err)
	}
	return nil
}

// ListOperationLogs 分页查询操作日志
func (r *SystemRepository) ListOperationLogs(userID, module string, page, pageSize int) ([]model.OperationLog, int64, error) {
	var list []model.OperationLog
	var total int64

	query := r.db.Model(&model.OperationLog{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WrapCountError("操作日志", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, errors.WrapQueryError("操作日志", err)
	}

	return list, total, nil
}
