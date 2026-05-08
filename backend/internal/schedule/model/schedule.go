package model

import (
	"time"

	"gorm.io/gorm"
)

// ScheduleInfo 排班信息
type ScheduleInfo struct {
	ID              string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	DoctorID        string         `gorm:"column:doctor_id;not null;type:varchar(64);index"`
	DoctorName      string         `gorm:"column:doctor_name;size:50"`
	DeptID          string         `gorm:"column:dept_id;not null;type:varchar(64);index"`
	DeptName        string         `gorm:"column:dept_name;size:50"`
	WorkDate        string         `gorm:"column:work_date;not null;size:10"` // 排班日期，格式 YYYY-MM-DD
	TimeSlot        int            `gorm:"column:time_slot;not null"`         // 1上午 2下午 3晚上
	MaxPatients     int            `gorm:"column:max_patients;not null"`      // 最大接诊数
	CurrentPatients int            `gorm:"column:current_patients;default:0"` // 当前已挂号数
	RoomNo          string         `gorm:"column:room_no;size:20"`            // 诊室号
	Status          int8           `gorm:"column:status;default:1"`           // 0停诊 1正常 2已满
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ScheduleInfo) TableName() string { return "schedules" }
