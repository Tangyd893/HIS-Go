package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

// Prescription 处方表
type Prescription struct {
	ID               string `gorm:"column:id;primaryKey;type:varchar(64)"`
	PatientID        string `gorm:"column:patient_id;not null;type:varchar(64);index"`
	PatientName      string `gorm:"column:patient_name;size:50"`
	DoctorID         string `gorm:"column:doctor_id;not null;type:varchar(64);index"`
	DiagnosisID      string `gorm:"column:diagnosis_id;type:varchar(64)"` // 关联诊断ID
	PrescriptionType int8   `gorm:"column:prescription_type;not null"`    // 1西药 2中成药 3中草药
	Status           int8   `gorm:"column:status;default:0"`              // 0草稿 1待审核 2已审核 3已收费 4已发药
	Note             string `gorm:"column:note;type:text"`                // 备注

	Details []PrescriptionDetail `gorm:"foreignKey:PrescriptionID"` // 处方明细

	Version optimisticlock.Version

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Prescription) TableName() string { return "prescriptions" }
