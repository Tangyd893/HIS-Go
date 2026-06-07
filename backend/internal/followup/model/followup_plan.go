package model

import (
	"time"

	"gorm.io/gorm"
)

// FollowupPlan 随访计划
type FollowupPlan struct {
	ID         string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PatientID  string         `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	TemplateID string         `gorm:"column:template_id;type:varchar(64)" json:"templateId,omitempty"`
	PlanName   string         `gorm:"column:plan_name;size:100" json:"content"`  // 计划名称 → 前端 content
	StartDate  string         `gorm:"column:start_date;size:10" json:"planDate"` // 开始日期 → 前端 planDate
	EndDate    string         `gorm:"column:end_date;size:10" json:"actualDate"` // 结束日期 → 前端 actualDate
	Frequency  int            `gorm:"column:frequency;not null" json:"frequency"`
	Status     int8           `gorm:"column:status;default:1" json:"status"` // 0已终止 1进行中 2已完成
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (FollowupPlan) TableName() string { return "followup_plans" }
