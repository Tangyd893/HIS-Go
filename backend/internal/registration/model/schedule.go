package model

import (
	"time"

	"gorm.io/plugin/optimisticlock"
)

// Schedule 号源表
type Schedule struct {
	ID          string  `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	DeptID      string  `gorm:"column:dept_id;not null;type:varchar(64);index" json:"deptId"`
	DeptName    string  `gorm:"column:dept_name;size:100" json:"deptName"`
	DoctorID    string  `gorm:"column:doctor_id;not null;type:varchar(64);index" json:"doctorId"`
	DoctorName  string  `gorm:"column:doctor_name;size:50" json:"doctorName"`
	Date        string  `gorm:"column:date;not null;size:10;index" json:"date"`
	TimeSlot    int     `gorm:"column:time_slot;not null;default:1" json:"timeSlot"`
	TotalCount  int     `gorm:"column:total_count;not null;default:0" json:"totalCount"`
	RemainCount int     `gorm:"column:remain_count;not null;default:0" json:"remainCount"`
	Fee         float64 `gorm:"column:fee;type:decimal(10,2);default:0" json:"fee"`
	Status      int8    `gorm:"column:status;default:1" json:"status"`

	Version optimisticlock.Version `json:"-"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
}

func (Schedule) TableName() string { return "schedules" }
