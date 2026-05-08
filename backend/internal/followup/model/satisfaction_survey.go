package model

import (
	"time"

	"gorm.io/gorm"
)

// SatisfactionSurvey 满意度调查
type SatisfactionSurvey struct {
	ID             string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	FollowupTaskID string         `gorm:"column:followup_task_id;not null;type:varchar(64);index"`
	PatientID      string         `gorm:"column:patient_id;not null;type:varchar(64);index"`
	Score          int            `gorm:"column:score;not null"` // 评分 1-5
	Feedback       string         `gorm:"column:feedback;type:text"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SatisfactionSurvey) TableName() string { return "satisfaction_surveys" }
