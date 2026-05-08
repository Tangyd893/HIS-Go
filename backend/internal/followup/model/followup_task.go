package model

import (
	"time"

	"gorm.io/gorm"
)

// FollowupTask 随访任务
type FollowupTask struct {
	ID          string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PlanID      string         `gorm:"column:plan_id;not null;type:varchar(64);index"`
	AssigneeID  string         `gorm:"column:assignee_id;type:varchar(64)"` // 执行人ID
	ExecuteDate string         `gorm:"column:execute_date;size:10"`         // 执行日期 YYYY-MM-DD
	Type        int8           `gorm:"column:type;not null"`                // 1电话 2问卷 3上门
	Content     string         `gorm:"column:content;type:text"`            // 随访内容
	Status      int8           `gorm:"column:status;default:0"`             // 0待执行 1已完成 2已跳过
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (FollowupTask) TableName() string { return "followup_tasks" }
