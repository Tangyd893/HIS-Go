package model

import (
	"time"

	"gorm.io/gorm"
)

// FollowupTask 随访任务
type FollowupTask struct {
	ID          string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PlanID      string         `gorm:"column:plan_id;not null;type:varchar(64);index" json:"planId"`
	AssigneeID  string         `gorm:"column:assignee_id;type:varchar(64)" json:"assigneeId,omitempty"`
	ExecuteDate string         `gorm:"column:execute_date;size:10" json:"executeDate"`
	Type        int8           `gorm:"column:type;not null" json:"type"`
	Content     string         `gorm:"column:content;type:text" json:"content,omitempty"`
	Status      int8           `gorm:"column:status;default:0" json:"status"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (FollowupTask) TableName() string { return "followup_tasks" }
