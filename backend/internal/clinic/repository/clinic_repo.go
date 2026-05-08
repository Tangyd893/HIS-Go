package repository

import (
	"fmt"

	"gorm.io/gorm"

	"his-go/internal/clinic/model"
	"his-go/pkg/errors"
)

// ClinicRepository 门诊诊疗数据仓库
type ClinicRepository struct {
	db *gorm.DB
}

// NewClinicRepository 创建门诊诊疗数据仓库
func NewClinicRepository(db *gorm.DB) *ClinicRepository {
	return &ClinicRepository{db: db}
}

// CreateRecord 创建门诊诊疗记录
func (r *ClinicRepository) CreateRecord(record *model.ClinicRecord) error {
	if err := r.db.Create(record).Error; err != nil {
		return fmt.Errorf("创建诊疗记录失败: %w", err)
	}
	return nil
}

// FindByID 根据ID查询诊疗记录
func (r *ClinicRepository) FindByID(id string) (*model.ClinicRecord, error) {
	var record model.ClinicRecord
	if err := r.db.Where("id = ?", id).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "诊疗记录不存在")
		}
		return nil, fmt.Errorf("查询诊疗记录失败: %w", err)
	}
	return &record, nil
}

// ListByPatient 分页查询患者诊疗记录
func (r *ClinicRepository) ListByPatient(patientID string, page, pageSize int) ([]model.ClinicRecord, int64, error) {
	var list []model.ClinicRecord
	var total int64

	query := r.db.Model(&model.ClinicRecord{}).Where("patient_id = ?", patientID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计诊疗记录失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询诊疗记录列表失败: %w", err)
	}

	return list, total, nil
}

// ListByDoctor 分页查询医生的诊疗记录
func (r *ClinicRepository) ListByDoctor(doctorID string, page, pageSize int) ([]model.ClinicRecord, int64, error) {
	var list []model.ClinicRecord
	var total int64

	query := r.db.Model(&model.ClinicRecord{}).Where("doctor_id = ?", doctorID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计诊疗记录失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询诊疗记录列表失败: %w", err)
	}

	return list, total, nil
}

// UpdateRecord 更新诊疗记录
func (r *ClinicRepository) UpdateRecord(record *model.ClinicRecord) error {
	result := r.db.Model(&model.ClinicRecord{}).Where("id = ?", record.ID).Updates(record)
	if result.Error != nil {
		return fmt.Errorf("更新诊疗记录失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "诊疗记录不存在")
	}
	return nil
}

// CreateExamRequest 创建检查申请单
func (r *ClinicRepository) CreateExamRequest(request *model.ExaminationRequest) error {
	if err := r.db.Create(request).Error; err != nil {
		return fmt.Errorf("创建检查申请失败: %w", err)
	}
	return nil
}

// ListExamRequests 查询患者的检查申请单列表
func (r *ClinicRepository) ListExamRequests(patientID string) ([]model.ExaminationRequest, error) {
	var list []model.ExaminationRequest
	if err := r.db.Where("patient_id = ?", patientID).
		Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("查询检查申请列表失败: %w", err)
	}
	return list, nil
}
