package repository

import (
	"time"

	"his-go/internal/inpatient/model"
	"his-go/pkg/errors"

	"gorm.io/gorm"
)

// InpatientRepository 住院管理数据仓库
type InpatientRepository struct {
	db *gorm.DB
}

// NewInpatientRepository 创建住院管理数据仓库
func NewInpatientRepository(db *gorm.DB) *InpatientRepository {
	return &InpatientRepository{db: db}
}

// Create 创建入院记录并初始化床位状态
func (r *InpatientRepository) Create(record *model.InpatientRecord) error {
	now := time.Now()
	record.AdmissionDate = &now
	record.Status = 1
	return r.db.Create(record).Error
}

// FindByID 根据ID查找住院记录
func (r *InpatientRepository) FindByID(id string) (*model.InpatientRecord, error) {
	var record model.InpatientRecord
	err := r.db.Where("id = ?", id).First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "住院记录不存在")
		}
		return nil, err
	}
	return &record, nil
}

// FindByPatientID 根据患者ID查找住院记录
func (r *InpatientRepository) FindByPatientID(patientID string) ([]model.InpatientRecord, error) {
	var records []model.InpatientRecord
	err := r.db.Where("patient_id = ?", patientID).Order("created_at DESC").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// List 分页查询住院列表
func (r *InpatientRepository) List(deptID string, status int, page, pageSize int) ([]model.InpatientRecord, int64, error) {
	var records []model.InpatientRecord
	var total int64

	query := r.db.Model(&model.InpatientRecord{})
	if deptID != "" {
		query = query.Where("dept_id = ?", deptID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// Discharge 出院：状态改为2，记录出院时间
func (r *InpatientRepository) Discharge(id string) error {
	now := time.Now()
	result := r.db.Model(&model.InpatientRecord{}).
		Where("id = ? AND status = 1", id).
		Updates(map[string]interface{}{
			"status":         2,
			"discharge_date": &now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "住院记录不存在或状态不允许出院")
	}
	return nil
}

// CreateOrder 创建医嘱
func (r *InpatientRepository) CreateOrder(order *model.MedicalOrder) error {
	return r.db.Create(order).Error
}

// ListOrders 查询住院患者的医嘱列表
func (r *InpatientRepository) ListOrders(inpatientID string) ([]model.MedicalOrder, error) {
	var orders []model.MedicalOrder
	err := r.db.Where("inpatient_id = ?", inpatientID).
		Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// CreateNursingRecord 创建护理记录
func (r *InpatientRepository) CreateNursingRecord(record *model.NursingRecord) error {
	return r.db.Create(record).Error
}

// ListNursingRecords 查询住院患者的护理记录列表
func (r *InpatientRepository) ListNursingRecords(inpatientID string) ([]model.NursingRecord, error) {
	var records []model.NursingRecord
	err := r.db.Where("inpatient_id = ?", inpatientID).
		Order("created_at DESC").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
