package model

import (
	"time"

	"gorm.io/gorm"
)

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	ID              string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	Name            string         `gorm:"column:name;not null;size:100"`     // 模板名称
	TitleTemplate   string         `gorm:"column:title_template;size:200"`    // 标题模板
	ContentTemplate string         `gorm:"column:content_template;type:text"` // 内容模板
	Channel         int8           `gorm:"column:channel;not null"`           // 1SMS 2邮件 3站内信 4微信
	Params          string         `gorm:"column:params;type:text"`           // 模板参数 JSON
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (NotificationTemplate) TableName() string { return "notification_templates" }
