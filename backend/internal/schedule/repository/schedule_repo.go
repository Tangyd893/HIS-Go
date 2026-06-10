package repository

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"his-go/internal/schedule/model"
	"his-go/pkg/demo"
	"his-go/pkg/errors"
)

// ScheduleRepository 排班数据仓库
type ScheduleRepository struct {
	db *gorm.DB
}

// NewScheduleRepository 创建排班数据仓库
func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// GenerateWeekly 按科室医生生成一周排班
func (r *ScheduleRepository) GenerateWeekly(startDate, endDate string, deptID string) ([]model.ScheduleInfo, error) {
	var schedules []model.ScheduleInfo

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeParamInvalid, "开始日期格式错误")
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeParamInvalid, "结束日期格式错误")
	}

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		workDate := d.Format("2006-01-02")
		doctorID := "doctor-wang"
		for slot := 1; slot <= 3; slot++ {
			schedule := model.ScheduleInfo{
				DoctorID:        doctorID,
				DoctorName:      demo.DoctorName(doctorID),
				DeptID:          deptID,
				DeptName:        demo.DeptName(deptID),
				WorkDate:        workDate,
				TimeSlot:        slot,
				MaxPatients:     50,
				CurrentPatients: 0,
				Status:          1,
				RoomNo:          strconv.Itoa(slot) + "室",
			}
			schedules = append(schedules, schedule)
		}
	}

	if err := r.db.Create(&schedules).Error; err != nil {
		return nil, fmt.Errorf("生成排班失败: %w", err)
	}

	return schedules, nil
}

// FindByDeptAndDate 按科室和日期查询排班（deptID 为空时返回当日全部排班）
func (r *ScheduleRepository) FindByDeptAndDate(deptID, date string) ([]model.ScheduleInfo, error) {
	var list []model.ScheduleInfo
	q := r.db.Model(&model.ScheduleInfo{})
	if date != "" {
		q = q.Where("work_date = ?", date)
	}
	if deptID != "" {
		q = q.Where("dept_id = ?", deptID)
	}
	if err := q.Order("dept_id ASC, time_slot ASC").Find(&list).Error; err != nil {
		return nil, errors.WrapQueryError("排班", err)
	}
	return list, nil
}

// FindByDoctor 按医生和日期查询排班
func (r *ScheduleRepository) FindByDoctor(doctorID, date string) ([]model.ScheduleInfo, error) {
	var list []model.ScheduleInfo
	if err := r.db.Where("doctor_id = ? AND work_date = ?", doctorID, date).
		Order("time_slot ASC").Find(&list).Error; err != nil {
		return nil, errors.WrapQueryError("排班", err)
	}
	return list, nil
}

// UpdateSchedule 更新排班信息
func (r *ScheduleRepository) UpdateSchedule(schedule *model.ScheduleInfo) error {
	result := r.db.Model(&model.ScheduleInfo{}).Where("id = ?", schedule.ID).Updates(schedule)
	if result.Error != nil {
		return errors.WrapUpdateError("排班", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "排班不存在")
	}
	return nil
}

// CancelSchedule 取消排班
func (r *ScheduleRepository) CancelSchedule(id string) error {
	result := r.db.Model(&model.ScheduleInfo{}).Where("id = ?", id).Update("status", 0)
	if result.Error != nil {
		return fmt.Errorf("取消排班失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "排班不存在")
	}
	return nil
}
