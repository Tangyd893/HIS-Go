package model

import (
	"time"

	"gorm.io/gorm"
)

// Notification 消息通知
type Notification struct {
	ID         string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	TemplateID string         `gorm:"column:template_id;type:varchar(64)" json:"templateId,omitempty"`
	ReceiverID string         `gorm:"column:receiver_id;not null;type:varchar(64);index" json:"receiverId"`
	Title      string         `gorm:"column:title;size:200" json:"title"`
	Content    string         `gorm:"column:content;type:text" json:"content"`
	Channel    int8           `gorm:"column:channel;not null" json:"channel"`
	Status     int8           `gorm:"column:status;default:0" json:"status"`
	SendTime   *time.Time     `gorm:"column:send_time" json:"sendTime,omitempty"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Notification) TableName() string { return "notifications" }
