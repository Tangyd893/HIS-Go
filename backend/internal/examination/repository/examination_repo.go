package repository

import (
	"fmt"

	"gorm.io/gorm"

	"his-go/internal/examination/model"
	"his-go/pkg/errors"
)

// ExaminationRepository 检查检验数据仓库
type ExaminationRepository struct {
	db *gorm.DB
}

// NewExaminationRepository 创建检查检验数据仓库
func NewExaminationRepository(db *gorm.DB) *ExaminationRepository {
	return &ExaminationRepository{db: db}
}

// CreateReport 创建检查报告
func (r *ExaminationRepository) CreateReport(report *model.ExaminationReport) error {
	if err := r.db.Create(report).Error; err != nil {
		return fmt.Errorf("创建检查报告失败: %w", err)
	}
	return nil
}

// FindByID 根据ID查询报告
func (r *ExaminationRepository) FindByID(id string) (*model.ExaminationReport, error) {
	var report model.ExaminationReport
	if err := r.db.Where("id = ?", id).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "检查报告不存在")
		}
		return nil, fmt.Errorf("查询检查报告失败: %w", err)
	}
	return &report, nil
}

// ListByPatient 分页查询患者检查报告
func (r *ExaminationRepository) ListByPatient(patientID string, status int, page, pageSize int) ([]model.ExaminationReport, int64, error) {
	var list []model.ExaminationReport
	var total int64

	query := r.db.Model(&model.ExaminationReport{}).Where("patient_id = ?", patientID)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计检查报告失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询检查报告列表失败: %w", err)
	}

	return list, total, nil
}

// UpdateReport 更新检查报告
func (r *ExaminationRepository) UpdateReport(report *model.ExaminationReport) error {
	result := r.db.Model(&model.ExaminationReport{}).Where("id = ?", report.ID).Updates(report)
	if result.Error != nil {
		return fmt.Errorf("更新检查报告失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "检查报告不存在")
	}
	return nil
}

// Review 审核检查报告
func (r *ExaminationRepository) Review(reportID, reviewerID string, approved bool, comment string) error {
	var report model.ExaminationReport
	if err := r.db.Where("id = ?", reportID).First(&report).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeNotFound, "检查报告不存在")
		}
		return fmt.Errorf("查询检查报告失败: %w", err)
	}

	if report.Status != 1 {
		return errors.NewAppError(errors.CodeConflict, "只有已检查状态的报告可进行审核")
	}

	updates := map[string]interface{}{
		"status":      int8(3), // 已发布
		"reviewer_id": reviewerID,
		"conclusion":  comment,
	}
	if approved {
		updates["status"] = int8(2) // 已审核
	}

	result := r.db.Model(&report).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("审核检查报告失败: %w", result.Error)
	}
	return nil
}
