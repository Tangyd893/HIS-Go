package model

import (
	"time"

	"gorm.io/gorm"
)

// Notification 消息通知
type Notification struct {
	ID         string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	TemplateID string         `gorm:"column:template_id;type:varchar(64)"` // 关联模板ID
	ReceiverID string         `gorm:"column:receiver_id;not null;type:varchar(64);index"`
	Title      string         `gorm:"column:title;size:200"`    // 通知标题
	Content    string         `gorm:"column:content;type:text"` // 通知内容
	Channel    int8           `gorm:"column:channel;not null"`  // 1SMS 2邮件 3站内信 4微信
	Status     int8           `gorm:"column:status;default:0"`  // 0待发送 1已发送 2发送失败
	SendTime   *time.Time     `gorm:"column:send_time"`         // 发送时间
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Notification) TableName() string { return "notifications" }
