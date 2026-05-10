package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"his-go/internal/health_record/model"
)

// HealthRecordRepository 健康档案数据仓库
type HealthRecordRepository struct {
	db *gorm.DB
}

// NewHealthRecordRepository 创建健康档案数据仓库
func NewHealthRecordRepository(db *gorm.DB) *HealthRecordRepository {
	return &HealthRecordRepository{db: db}
}

// GetSummary 获取或创建档案摘要
func (r *HealthRecordRepository) GetSummary(patientID string) (*model.HealthRecordSummary, error) {
	var summary model.HealthRecordSummary
	err := r.db.Where("patient_id = ?", patientID).First(&summary).Error
	if err == nil {
		return &summary, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询档案摘要失败: %w", err)
	}

	summary = model.HealthRecordSummary{PatientID: patientID}
	if err := r.db.Create(&summary).Error; err != nil {
		return nil, fmt.Errorf("创建档案摘要失败: %w", err)
	}
	return &summary, nil
}

// GetTimeline 获取时间轴事件
func (r *HealthRecordRepository) GetTimeline(patientID string) ([]model.TimelineEvent, error) {
	var events []model.TimelineEvent
	if err := r.db.Where("patient_id = ?", patientID).
		Order("date DESC, created_at DESC").Find(&events).Error; err != nil {
		return nil, fmt.Errorf("查询时间轴事件失败: %w", err)
	}
	return events, nil
}

// AddTimelineEvent 添加时间轴事件
func (r *HealthRecordRepository) AddTimelineEvent(event *model.TimelineEvent) error {
	if err := r.db.Create(event).Error; err != nil {
		return fmt.Errorf("添加时间轴事件失败: %w", err)
	}
	return nil
}

// GrantAuthorization 授权查看档案
func (r *HealthRecordRepository) GrantAuthorization(auth *model.RecordAuthorization) error {
	if err := r.db.Create(auth).Error; err != nil {
		return fmt.Errorf("授权失败: %w", err)
	}
	return nil
}

// RevokeAuthorization 撤销授权
func (r *HealthRecordRepository) RevokeAuthorization(patientID, doctorID string) error {
	result := r.db.Model(&model.RecordAuthorization{}).
		Where("patient_id = ? AND doctor_id = ? AND status = 1", patientID, doctorID).
		Update("status", 0)
	if result.Error != nil {
		return fmt.Errorf("撤销授权失败: %w", result.Error)
	}
	return nil
}

// CheckAuthorization 检查授权状态（含过期校验）
func (r *HealthRecordRepository) CheckAuthorization(patientID, doctorID string) bool {
	var count int64
	now := time.Now().Format("2006-01-02 15:04:05")
	r.db.Model(&model.RecordAuthorization{}).
		Where("patient_id = ? AND doctor_id = ? AND status = 1", patientID, doctorID).
		Where("expire_time = '' OR expire_time > ?", now).
		Count(&count)
	return count > 0
}
