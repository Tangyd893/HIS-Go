package model

import (
	"time"
)

// ClinicRecord 门诊诊疗记录表
type ClinicRecord struct {
	ID             string     `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	RegistrationID string     `gorm:"column:registration_id;type:varchar(64);index" json:"registrationId,omitempty"`
	PatientID      string     `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId,omitempty"`
	PatientName    string     `gorm:"column:patient_name;size:50" json:"patientName,omitempty"`
	DoctorID       string     `gorm:"column:doctor_id;not null;type:varchar(64);index" json:"doctorId,omitempty"`
	DoctorName     string     `gorm:"-" json:"doctorName,omitempty"`
	ChiefComplaint string     `gorm:"column:chief_complaint;type:text" json:"chiefComplaint,omitempty"`
	PresentIllness string     `gorm:"column:present_illness;type:text" json:"presentIllness,omitempty"`
	Diagnosis      string     `gorm:"column:diagnosis;type:text" json:"diagnosis,omitempty"`
	IcdCode        string     `gorm:"column:icd_code;size:20" json:"icdCode,omitempty"`
	Advice         string     `gorm:"column:advice;type:text" json:"advice,omitempty"`
	Status         int8       `gorm:"column:status;default:0" json:"status"`
	VisitTime      *time.Time `gorm:"column:visit_time" json:"visitTime,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
}

func (ClinicRecord) TableName() string { return "clinic_records" }
