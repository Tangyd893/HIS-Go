package model

import (
	"time"

	"gorm.io/gorm"
)

// ExaminationReport 检查检验报告
type ExaminationReport struct {
	ID            string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PatientID     string         `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	PatientName   string         `gorm:"column:patient_name;size:50" json:"patientName,omitempty"`
	ExamRequestID string         `gorm:"column:exam_request_id;type:varchar(64)" json:"examRequestId,omitempty"`
	ExamType      string         `gorm:"column:exam_type;size:50" json:"examType"`   // 检查类型：CT/DR/MR/超声/检验
	ExamItem      string         `gorm:"column:exam_item;size:100" json:"examItem"`  // 检查项目
	BodyPart      string         `gorm:"column:body_part;size:50" json:"bodyPart,omitempty"`
	Findings      string         `gorm:"column:findings;type:text" json:"findings,omitempty"`
	Impression    string         `gorm:"column:impression;type:text" json:"impression,omitempty"`
	Conclusion    string         `gorm:"column:conclusion;type:text" json:"result"` // 诊断结论 → 前端 result
	TechnicianID  string         `gorm:"column:technician_id;type:varchar(64)" json:"technicianId,omitempty"`
	ReviewerID    string         `gorm:"column:reviewer_id;type:varchar(64)" json:"reviewerId,omitempty"`
	Status        int8           `gorm:"column:status;default:0" json:"status"` // 0待检查 1已检查 2已审核 3已发布
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (ExaminationReport) TableName() string { return "examination_reports" }
