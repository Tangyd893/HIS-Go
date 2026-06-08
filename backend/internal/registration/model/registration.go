package model

import (
	"time"
)

// Registration 挂号记录表
type Registration struct {
	ID               string `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	PatientID        string `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	PatientName      string `gorm:"column:patient_name;size:50" json:"patientName"`
	ScheduleID       string `gorm:"column:schedule_id;not null;type:varchar(64);index" json:"scheduleId"`
	RegistrationDate string `gorm:"column:registration_date;not null;size:10" json:"registrationDate"`
	QueueNumber      int    `gorm:"column:queue_number;default:0" json:"queueNumber"`
	Status           int8   `gorm:"column:status;default:0" json:"status"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
}

func (Registration) TableName() string { return "registrations" }
