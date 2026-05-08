package model

import (
	"time"

	"gorm.io/gorm"
)

// TimelineEvent 时间轴事件
type TimelineEvent struct {
	ID          string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID   string         `gorm:"column:patient_id;not null;type:varchar(64);index"`
	Date        string         `gorm:"column:date;size:10"`                // 事件日期 YYYY-MM-DD
	EventType   string         `gorm:"column:event_type;size:30"`          // visit/prescription/examination/followup
	Description string         `gorm:"column:description;size:500"`        // 事件描述
	RelatedID   string         `gorm:"column:related_id;type:varchar(64)"` // 关联业务记录ID
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (TimelineEvent) TableName() string { return "timeline_events" }
