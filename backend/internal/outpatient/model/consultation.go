package model

import (
	"time"

	"gorm.io/gorm"
)

// Consultation 在线问诊
type Consultation struct {
	ID          string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PatientID   string         `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	DoctorID    string         `gorm:"column:doctor_id;not null;type:varchar(64);index" json:"doctorId,omitempty"`
	Type        int8           `gorm:"column:type;not null" json:"type"` // 1图文 2视频
	Description string         `gorm:"column:description;type:text" json:"description,omitempty"`
	Images      string         `gorm:"column:images;type:text" json:"images,omitempty"` // JSON数组格式存储图片URL
	Status      int8           `gorm:"column:status;default:0" json:"status"` // 0待接诊 1接诊中 2已完成 3已取消
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Consultation) TableName() string { return "consultations" }
