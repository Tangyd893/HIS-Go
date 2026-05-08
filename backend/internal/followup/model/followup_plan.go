package model

import (
	"time"

	"gorm.io/gorm"
)

// FollowupPlan 随访计划
type FollowupPlan struct {
	ID         string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID  string         `gorm:"column:patient_id;not null;type:varchar(64);index"`
	TemplateID string         `gorm:"column:template_id;type:varchar(64)"`
	PlanName   string         `gorm:"column:plan_name;size:100"` // 计划名称
	StartDate  string         `gorm:"column:start_date;size:10"` // 开始日期 YYYY-MM-DD
	EndDate    string         `gorm:"column:end_date;size:10"`   // 结束日期 YYYY-MM-DD
	Frequency  int            `gorm:"column:frequency;not null"` // 1每周 2每两周 3每月
	Status     int8           `gorm:"column:status;default:1"`   // 0已终止 1进行中 2已完成
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (FollowupPlan) TableName() string { return "followup_plans" }
