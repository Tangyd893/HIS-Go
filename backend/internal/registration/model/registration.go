package model

import (
	"time"
)

// Registration 挂号记录表
type Registration struct {
	ID               string `gorm:"column:id;primaryKey;type:varchar(64)"`
	PatientID        string `gorm:"column:patient_id;not null;type:varchar(64);index"`
	PatientName      string `gorm:"column:patient_name;size:50"`
	ScheduleID       string `gorm:"column:schedule_id;not null;type:varchar(64);index"`
	RegistrationDate string `gorm:"column:registration_date;not null;size:10"` // 挂号日期 YYYY-MM-DD
	QueueNumber      int    `gorm:"column:queue_number;default:0"`             // 排队号
	Status           int8   `gorm:"column:status;default:0"`                   // 0已预约 1已签到 2已就诊 3已取消

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Registration) TableName() string { return "registrations" }
