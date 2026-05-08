package model

import (
	"time"
)

// ClinicRecord 门诊诊疗记录表
type ClinicRecord struct {
	ID             string     `gorm:"column:id;primaryKey;type:varchar(64)"`
	RegistrationID string     `gorm:"column:registration_id;type:varchar(64);index"`
	PatientID      string     `gorm:"column:patient_id;not null;type:varchar(64);index"`
	PatientName    string     `gorm:"column:patient_name;size:50"`
	DoctorID       string     `gorm:"column:doctor_id;not null;type:varchar(64);index"`
	ChiefComplaint string     `gorm:"column:chief_complaint;type:text"` // 主诉
	PresentIllness string     `gorm:"column:present_illness;type:text"` // 现病史
	Diagnosis      string     `gorm:"column:diagnosis;type:text"`       // 诊断
	IcdCode        string     `gorm:"column:icd_code;size:20"`          // ICD编码
	Advice         string     `gorm:"column:advice;type:text"`          // 医嘱
	Status         int8       `gorm:"column:status;default:0"`          // 状态
	VisitTime      *time.Time `gorm:"column:visit_time"`                // 就诊时间

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ClinicRecord) TableName() string { return "clinic_records" }
