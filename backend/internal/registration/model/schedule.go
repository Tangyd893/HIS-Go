package model

import (
	"time"

	"gorm.io/plugin/optimisticlock"
)

// Schedule 号源表
type Schedule struct {
	ID          string  `gorm:"column:id;primaryKey;type:varchar(64)"`
	DeptID      string  `gorm:"column:dept_id;not null;type:varchar(64);index"`
	DeptName    string  `gorm:"column:dept_name;size:100"`
	DoctorID    string  `gorm:"column:doctor_id;not null;type:varchar(64);index"`
	DoctorName  string  `gorm:"column:doctor_name;size:50"`
	Date        string  `gorm:"column:date;not null;size:10;index"`      // 排班日期，格式 YYYY-MM-DD
	TimeSlot    int     `gorm:"column:time_slot;not null;default:1"`     // 时段：1上午 2下午 3晚上
	TotalCount  int     `gorm:"column:total_count;not null;default:0"`   // 总号源数
	RemainCount int     `gorm:"column:remain_count;not null;default:0"`  // 剩余号源数
	Fee         float64 `gorm:"column:fee;type:decimal(10,2);default:0"` // 挂号费
	Status      int8    `gorm:"column:status;default:1"`                 // 状态：1启用 0停用

	Version optimisticlock.Version

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Schedule) TableName() string { return "schedules" }
