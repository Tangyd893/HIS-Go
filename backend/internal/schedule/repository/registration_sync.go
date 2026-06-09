package repository

import (
	"fmt"

	"gorm.io/gorm"

	regmodel "his-go/internal/registration/model"
	"his-go/internal/schedule/model"
	"his-go/pkg/demo"
)

// RegistrationSyncRepository 将排班同步到挂号号源库 (his_registration)
type RegistrationSyncRepository struct {
	db *gorm.DB
}

// NewRegistrationSyncRepository 创建挂号号源同步仓库
func NewRegistrationSyncRepository(db *gorm.DB) *RegistrationSyncRepository {
	return &RegistrationSyncRepository{db: db}
}

// SyncSchedules 幂等同步排班到 registration.schedules
func (r *RegistrationSyncRepository) SyncSchedules(schedules []model.ScheduleInfo) error {
	for _, s := range schedules {
		doctorID := s.DoctorID
		if doctorID == "" {
			doctorID = "doctor-wang"
		}
		remain := s.MaxPatients - s.CurrentPatients
		if remain < 0 {
			remain = 0
		}
		reg := regmodel.Schedule{
			ID:          s.ID,
			DeptID:      s.DeptID,
			DeptName:    demo.DeptName(s.DeptID),
			DoctorID:    doctorID,
			DoctorName:  demo.DoctorName(doctorID),
			Date:        s.WorkDate,
			TimeSlot:    s.TimeSlot,
			TotalCount:  s.MaxPatients,
			RemainCount: remain,
			Fee:         10,
			Status:      s.Status,
		}
		if err := r.db.Save(&reg).Error; err != nil {
			return fmt.Errorf("同步挂号号源失败: %w", err)
		}
	}
	return nil
}
