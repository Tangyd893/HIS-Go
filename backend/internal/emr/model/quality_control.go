package model

import "time"

// QualityControl 质控记录模型，表名 quality_controls
type QualityControl struct {
	ID         string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	RecordID   string    `gorm:"column:record_id;not null;type:varchar(64);index"`
	ReviewerID string    `gorm:"column:reviewer_id;not null;type:varchar(64)"`
	Level      int       `gorm:"column:level;not null;default:0;comment:'质控等级'"`
	Comment    string    `gorm:"column:comment;type:text;comment:'质控评语'"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (QualityControl) TableName() string { return "quality_controls" }
