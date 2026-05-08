package repository

import (
	"fmt"

	"gorm.io/gorm"

	"his-go/internal/prescription/model"
	"his-go/pkg/errors"
)

// PrescriptionRepository 处方数据仓库
type PrescriptionRepository struct {
	db *gorm.DB
}

// NewPrescriptionRepository 创建处方数据仓库
func NewPrescriptionRepository(db *gorm.DB) *PrescriptionRepository {
	return &PrescriptionRepository{db: db}
}

// Create 事务中同时创建处方和明细
func (r *PrescriptionRepository) Create(p *model.Prescription, details []model.PrescriptionDetail) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return fmt.Errorf("创建处方失败: %w", err)
		}

		for i := range details {
			details[i].PrescriptionID = p.ID
		}
		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return fmt.Errorf("创建处方明细失败: %w", err)
			}
		}

		return nil
	})
}

// FindByID 根据ID查询处方（预加载明细）
func (r *PrescriptionRepository) FindByID(id string) (*model.Prescription, error) {
	var p model.Prescription
	if err := r.db.Preload("Details").Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "处方不存在")
		}
		return nil, fmt.Errorf("查询处方失败: %w", err)
	}
	return &p, nil
}

// ListByPatient 分页查询患者处方
func (r *PrescriptionRepository) ListByPatient(patientID string, page, pageSize int) ([]model.Prescription, int64, error) {
	var list []model.Prescription
	var total int64

	query := r.db.Model(&model.Prescription{}).Where("patient_id = ?", patientID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计处方失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).
		Preload("Details").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询处方列表失败: %w", err)
	}

	return list, total, nil
}

// ListByDoctor 分页查询医生的处方
func (r *PrescriptionRepository) ListByDoctor(doctorID string, page, pageSize int) ([]model.Prescription, int64, error) {
	var list []model.Prescription
	var total int64

	query := r.db.Model(&model.Prescription{}).Where("doctor_id = ?", doctorID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计处方失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).
		Preload("Details").Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询处方列表失败: %w", err)
	}

	return list, total, nil
}

// UpdateStatus 使用乐观锁更新处方状态
func (r *PrescriptionRepository) UpdateStatus(id string, status int8) error {
	var p model.Prescription
	if err := r.db.Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "处方不存在")
		}
		return fmt.Errorf("查询处方失败: %w", err)
	}

	p.Status = status
	result := r.db.Model(&p).Select("status", "version").Updates(&p)
	if result.Error != nil {
		return fmt.Errorf("更新处方状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeConflict, "处方信息已变更，请重新尝试")
	}

	return nil
}

// Review 审核处方（通过/驳回）
func (r *PrescriptionRepository) Review(id string, approved bool, comment string) error {
	var p model.Prescription
	if err := r.db.Where("id = ?", id).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "处方不存在")
		}
		return fmt.Errorf("查询处方失败: %w", err)
	}

	if p.Status != 1 {
		return errors.NewAppError(errors.CodeConflict, "只有待审核状态的处方可进行审核")
	}

	if approved {
		p.Status = 2 // 已审核
	} else {
		p.Status = 0 // 退回草稿
	}
	p.Note = comment

	result := r.db.Model(&p).Select("status", "note", "version").Updates(&p)
	if result.Error != nil {
		return fmt.Errorf("审核处方失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeConflict, "处方信息已变更，请重新尝试")
	}

	return nil
}
