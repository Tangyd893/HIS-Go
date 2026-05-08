package model

import (
	"time"

	"gorm.io/gorm"
)

// OperationLog 操作日志
type OperationLog struct {
	ID        string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);index"`
	Username  string         `gorm:"column:username;size:50"`
	Module    string         `gorm:"column:module;size:50;index"` // 操作模块
	Action    string         `gorm:"column:action;size:50"`       // 操作动作
	Method    string         `gorm:"column:method;size:10"`       // HTTP方法
	URL       string         `gorm:"column:url;size:200"`         // 请求URL
	IP        string         `gorm:"column:ip;size:50"`           // 操作IP
	Params    string         `gorm:"column:params;type:text"`     // 请求参数
	Result    string         `gorm:"column:result;type:text"`     // 操作结果
	Status    int8           `gorm:"column:status;default:1"`     // 0失败 1成功
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (OperationLog) TableName() string { return "operation_logs" }
