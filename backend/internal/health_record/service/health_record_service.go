package service

import (
	"his-go/internal/health_record/model"
	"his-go/internal/health_record/repository"
)

// HealthRecordService 健康档案业务服务
type HealthRecordService struct {
	repo *repository.HealthRecordRepository
}

// NewHealthRecordService 创建健康档案业务服务
func NewHealthRecordService(repo *repository.HealthRecordRepository) *HealthRecordService {
	return &HealthRecordService{repo: repo}
}

// GetSummary 获取档案摘要
func (s *HealthRecordService) GetSummary(patientID string) (*model.HealthRecordSummary, error) {
	return s.repo.GetSummary(patientID)
}

// GetTimeline 获取时间轴
func (s *HealthRecordService) GetTimeline(patientID string) ([]model.TimelineEvent, error) {
	return s.repo.GetTimeline(patientID)
}

// AddTimelineEvent 添加时间轴事件
func (s *HealthRecordService) AddTimelineEvent(event *model.TimelineEvent) error {
	return s.repo.AddTimelineEvent(event)
}

// GrantAuthorization 授权查看档案
func (s *HealthRecordService) GrantAuthorization(auth *model.RecordAuthorization) error {
	auth.Status = 1 // 已授权
	return s.repo.GrantAuthorization(auth)
}

// RevokeAuthorization 撤销授权
func (s *HealthRecordService) RevokeAuthorization(patientID, doctorID string) error {
	return s.repo.RevokeAuthorization(patientID, doctorID)
}

// CheckAuthorization 检查授权状态
func (s *HealthRecordService) CheckAuthorization(patientID, doctorID string) bool {
	return s.repo.CheckAuthorization(patientID, doctorID)
}
