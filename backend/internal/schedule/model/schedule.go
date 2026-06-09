package model

import (
	"time"

	"gorm.io/gorm"
)

// ScheduleInfo 排班信息
type ScheduleInfo struct {
	ID              string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	DoctorID        string         `gorm:"column:doctor_id;not null;type:varchar(64);index" json:"doctorId"`
	DoctorName      string         `gorm:"column:doctor_name;size:50" json:"doctorName"`
	DeptID          string         `gorm:"column:dept_id;not null;type:varchar(64);index" json:"deptId"`
	DeptName        string         `gorm:"column:dept_name;size:50" json:"deptName"`
	WorkDate        string         `gorm:"column:work_date;not null;size:10" json:"date"`
	TimeSlot        int            `gorm:"column:time_slot;not null" json:"timeSlot"`
	MaxPatients     int            `gorm:"column:max_patients;not null" json:"maxPatients"`
	CurrentPatients int            `gorm:"column:current_patients;default:0" json:"currentPatients"`
	RoomNo          string         `gorm:"column:room_no;size:20" json:"roomNo,omitempty"`
	Status          int8           `gorm:"column:status;default:1" json:"status"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (ScheduleInfo) TableName() string { return "schedules" }
