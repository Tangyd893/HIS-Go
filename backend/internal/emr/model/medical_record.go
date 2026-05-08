package model

import "time"

// MedicalRecord 病历模型，表名 medical_records
type MedicalRecord struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID      string    `gorm:"column:patient_id;not null;type:varchar(64);index"`
	ClinicRecordID string    `gorm:"column:clinic_record_id;type:varchar(64)"`
	TemplateID     string    `gorm:"column:template_id;type:varchar(64)"`
	ChiefComplaint string    `gorm:"column:chief_complaint;type:text;comment:'主诉'"`
	PresentIllness string    `gorm:"column:present_illness;type:text;comment:'现病史'"`
	PastHistory    string    `gorm:"column:past_history;type:text;comment:'既往史'"`
	PhysicalExam   string    `gorm:"column:physical_exam;type:text;comment:'体格检查'"`
	AuxiliaryExam  string    `gorm:"column:auxiliary_exam;type:text;comment:'辅助检查'"`
	Diagnosis      string    `gorm:"column:diagnosis;type:text;comment:'诊断'"`
	TreatmentPlan  string    `gorm:"column:treatment_plan;type:text;comment:'处理计划'"`
	QualityLevel   int       `gorm:"column:quality_level;default:0;comment:'质控等级'"`
	Status         int8      `gorm:"column:status;default:0;comment:'0草稿 1已完成 2已质控 3已归档'"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (MedicalRecord) TableName() string { return "medical_records" }
