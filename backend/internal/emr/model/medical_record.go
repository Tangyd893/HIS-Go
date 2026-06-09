package model

import "time"

// MedicalRecord 病历模型，表名 medical_records
type MedicalRecord struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PatientID      string    `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	PatientName    string    `gorm:"-" json:"patientName,omitempty"`
	ClinicRecordID string    `gorm:"column:clinic_record_id;type:varchar(64)" json:"clinicRecordId,omitempty"`
	DoctorID       string    `gorm:"-" json:"doctorId,omitempty"`
	TemplateID     string    `gorm:"column:template_id;type:varchar(64)" json:"templateId"`
	ChiefComplaint string    `gorm:"column:chief_complaint;type:text" json:"chiefComplaint,omitempty"`
	PresentIllness string    `gorm:"column:present_illness;type:text" json:"presentIllness,omitempty"`
	PastHistory    string    `gorm:"column:past_history;type:text" json:"pastHistory,omitempty"`
	PhysicalExam   string    `gorm:"column:physical_exam;type:text" json:"physicalExam,omitempty"`
	AuxiliaryExam  string    `gorm:"column:auxiliary_exam;type:text" json:"auxiliaryExam,omitempty"`
	Diagnosis      string    `gorm:"column:diagnosis;type:text" json:"diagnosis,omitempty"`
	TreatmentPlan  string    `gorm:"column:treatment_plan;type:text" json:"treatmentPlan,omitempty"`
	QualityLevel   int       `gorm:"column:quality_level;default:0" json:"qualityLevel,omitempty"`
	Status         int8      `gorm:"column:status;default:0" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
}

func (MedicalRecord) TableName() string { return "medical_records" }
