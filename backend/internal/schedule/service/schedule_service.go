package service

import (
	"his-go/internal/schedule/model"
	"his-go/internal/schedule/repository"
	"his-go/pkg/errors"
)

// ScheduleService 排班业务服务
type ScheduleService struct {
	repo    *repository.ScheduleRepository
	regSync *repository.RegistrationSyncRepository
}

// NewScheduleService 创建排班业务服务
func NewScheduleService(repo *repository.ScheduleRepository, regSync *repository.RegistrationSyncRepository) *ScheduleService {
	return &ScheduleService{repo: repo, regSync: regSync}
}

// GenerateWeekly 按科室医生生成一周排班，并同步到挂号号源库
func (s *ScheduleService) GenerateWeekly(startDate, endDate, deptID string) ([]model.ScheduleInfo, error) {
	if startDate >= endDate {
		return nil, errors.NewAppError(errors.CodeParamInvalid, "开始日期必须早于结束日期")
	}
	schedules, err := s.repo.GenerateWeekly(startDate, endDate, deptID)
	if err != nil {
		return nil, err
	}
	if s.regSync != nil {
		_ = s.regSync.SyncSchedules(schedules)
	}
	return schedules, nil
}

// ListByDept 按科室和日期查询排班
func (s *ScheduleService) ListByDept(deptID, date string) ([]model.ScheduleInfo, error) {
	return s.repo.FindByDeptAndDate(deptID, date)
}

// ListByDoctor 按医生和日期查询排班
func (s *ScheduleService) ListByDoctor(doctorID, date string) ([]model.ScheduleInfo, error) {
	return s.repo.FindByDoctor(doctorID, date)
}

// UpdateSchedule 更新排班信息
func (s *ScheduleService) UpdateSchedule(schedule *model.ScheduleInfo) error {
	return s.repo.UpdateSchedule(schedule)
}

// CancelSchedule 取消排班
func (s *ScheduleService) CancelSchedule(id string) error {
	return s.repo.CancelSchedule(id)
}
