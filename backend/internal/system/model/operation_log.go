package model

import (
	"time"

	"gorm.io/gorm"
)

// OperationLog 操作日志
type OperationLog struct {
	ID        string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);index" json:"userId"`
	Username  string         `gorm:"column:username;size:50" json:"username,omitempty"`
	Module    string         `gorm:"column:module;size:50;index" json:"module"`
	Action    string         `gorm:"column:action;size:50" json:"operation"`
	Method    string         `gorm:"column:method;size:10" json:"method,omitempty"`
	URL       string         `gorm:"column:url;size:200" json:"url,omitempty"`
	IP        string         `gorm:"column:ip;size:50" json:"ip,omitempty"`
	Params    string         `gorm:"column:params;type:text" json:"params,omitempty"`
	Result    string         `gorm:"column:result;type:text" json:"result,omitempty"`
	Status    int8           `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (OperationLog) TableName() string { return "operation_logs" }
