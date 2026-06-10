package service

import (
	"time"

	"his-go/internal/health_record/model"
	"his-go/internal/health_record/repository"
	"his-go/pkg/errors"
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
	if patientID == "" {
		return nil, errors.NewAppError(errors.CodeParamInvalid, errors.MsgPatientIDRequired)
	}
	return s.repo.GetSummary(patientID)
}

// GetTimeline 获取时间轴
func (s *HealthRecordService) GetTimeline(patientID string) ([]model.TimelineEvent, error) {
	if patientID == "" {
		return nil, errors.NewAppError(errors.CodeParamInvalid, errors.MsgPatientIDRequired)
	}
	return s.repo.GetTimeline(patientID)
}

// AddTimelineEvent 添加时间轴事件
func (s *HealthRecordService) AddTimelineEvent(event *model.TimelineEvent) error {
	if event.PatientID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, errors.MsgPatientIDRequired)
	}
	if event.EventType == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "事件类型不能为空")
	}
	validTypes := map[string]bool{
		"visit": true, "prescription": true, "examination": true, "followup": true,
	}
	if !validTypes[event.EventType] {
		return errors.NewAppError(errors.CodeParamInvalid, "无效的事件类型，合法值：visit/prescription/examination/followup")
	}
	return s.repo.AddTimelineEvent(event)
}

// GrantAuthorization 授权查看档案（设置30天有效期，防止重复授权）
func (s *HealthRecordService) GrantAuthorization(auth *model.RecordAuthorization) error {
	if auth.PatientID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, errors.MsgPatientIDRequired)
	}
	if auth.DoctorID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, errors.MsgDoctorIDRequired)
	}

	if s.repo.CheckAuthorization(auth.PatientID, auth.DoctorID) {
		return errors.NewAppError(errors.CodeAuthDuplicate, "已存在有效授权")
	}

	now := time.Now()
	auth.AuthTime = now.Format("2006-01-02 15:04:05")
	auth.ExpireTime = now.AddDate(0, 0, 30).Format("2006-01-02 15:04:05")
	auth.Status = 1

	return s.repo.GrantAuthorization(auth)
}

// RevokeAuthorization 撤销授权（幂等：已撤销的授权不报错）
func (s *HealthRecordService) RevokeAuthorization(patientID, doctorID string) error {
	if patientID == "" || doctorID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "患者ID和医生ID不能为空")
	}
	return s.repo.RevokeAuthorization(patientID, doctorID)
}

// CheckAuthorization 检查授权状态（含有效期校验）
func (s *HealthRecordService) CheckAuthorization(patientID, doctorID string) bool {
	if patientID == "" || doctorID == "" {
		return false
	}
	return s.repo.CheckAuthorization(patientID, doctorID)
}
