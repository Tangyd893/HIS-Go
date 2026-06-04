package model

import (
	"time"

	"gorm.io/gorm"
)

// ChronicContract 慢病签约
type ChronicContract struct {
	ID           string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PatientID    string         `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	DoctorID     string         `gorm:"column:doctor_id;not null;type:varchar(64);index" json:"doctorId,omitempty"`
	DiseaseType  string         `gorm:"column:disease_type;size:100" json:"diseaseType"` // 慢病类型：高血压/糖尿病/冠心病等
	ContractDate string         `gorm:"column:contract_date;size:10" json:"contractDate"` // 签约日期 YYYY-MM-DD
	EndDate      string         `gorm:"column:end_date;size:10" json:"endDate"`           // 到期日期 YYYY-MM-DD
	Status       int8           `gorm:"column:status;default:1" json:"status"`            // 0已解约 1签约中 2已到期
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (ChronicContract) TableName() string { return "chronic_contracts" }
