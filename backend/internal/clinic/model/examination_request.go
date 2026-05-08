package model

import (
	"time"
)

// ExaminationRequest 检查申请单
type ExaminationRequest struct {
	ID                string `gorm:"column:id;primaryKey;type:varchar(64)"`
	ClinicRecordID    string `gorm:"column:clinic_record_id;not null;type:varchar(64);index"`
	PatientID         string `gorm:"column:patient_id;not null;type:varchar(64);index"`
	ExamType          string `gorm:"column:exam_type;size:50"`            // 检查类型
	ExamItem          string `gorm:"column:exam_item;size:200"`           // 检查项目
	BodyPart          string `gorm:"column:body_part;size:100"`           // 检查部位
	ClinicalDiagnosis string `gorm:"column:clinical_diagnosis;type:text"` // 临床诊断
	Note              string `gorm:"column:note;type:text"`               // 备注
	Status            int8   `gorm:"column:status;default:0"`             // 状态

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ExaminationRequest) TableName() string { return "examination_requests" }
