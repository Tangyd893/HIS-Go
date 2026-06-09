package model

import (
	"time"

	"gorm.io/gorm"
)

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	ID              string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	Name            string         `gorm:"column:name;not null;size:100" json:"name"`
	TitleTemplate   string         `gorm:"column:title_template;size:200" json:"titleTemplate,omitempty"`
	ContentTemplate string         `gorm:"column:content_template;type:text" json:"contentTemplate,omitempty"`
	Channel         int8           `gorm:"column:channel;not null" json:"channel"`
	Params          string         `gorm:"column:params;type:text" json:"params,omitempty"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (NotificationTemplate) TableName() string { return "notification_templates" }
