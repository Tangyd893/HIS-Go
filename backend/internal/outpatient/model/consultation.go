package model

import (
	"time"

	"gorm.io/gorm"
)

// Consultation 在线问诊
type Consultation struct {
	ID          string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID   string         `gorm:"column:patient_id;not null;type:varchar(64);index"`
	DoctorID    string         `gorm:"column:doctor_id;not null;type:varchar(64);index"`
	Type        int8           `gorm:"column:type;not null"` // 1图文 2视频
	Description string         `gorm:"column:description;type:text"`
	Images      string         `gorm:"column:images;type:text"` // JSON数组格式存储图片URL
	Status      int8           `gorm:"column:status;default:0"` // 0待接诊 1接诊中 2已完成 3已取消
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Consultation) TableName() string { return "consultations" }
