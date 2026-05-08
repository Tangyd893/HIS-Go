package model

import (
	"time"

	"gorm.io/gorm"
)

// ExaminationReport 检查检验报告
type ExaminationReport struct {
	ID            string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID     string         `gorm:"column:patient_id;not null;type:varchar(64);index"`
	PatientName   string         `gorm:"column:patient_name;size:50"`
	ExamRequestID string         `gorm:"column:exam_request_id;type:varchar(64)"`
	ExamType      string         `gorm:"column:exam_type;size:50"`              // 检查类型：CT/DR/MR/超声/检验
	ExamItem      string         `gorm:"column:exam_item;size:100"`             // 检查项目
	BodyPart      string         `gorm:"column:body_part;size:50"`              // 检查部位
	Findings      string         `gorm:"column:findings;type:text"`             // 检查所见
	Impression    string         `gorm:"column:impression;type:text"`           // 影像印象
	Conclusion    string         `gorm:"column:conclusion;type:text"`           // 诊断结论
	TechnicianID  string         `gorm:"column:technician_id;type:varchar(64)"` // 技师ID
	ReviewerID    string         `gorm:"column:reviewer_id;type:varchar(64)"`   // 审核医师ID
	Status        int8           `gorm:"column:status;default:0"`               // 0待检查 1已检查 2已审核 3已发布
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ExaminationReport) TableName() string { return "examination_reports" }
