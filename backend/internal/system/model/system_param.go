package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemParam 系统参数
type SystemParam struct {
	ID         string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	ParamName  string         `gorm:"column:param_name;not null;size:100"`           // 参数名称
	ParamKey   string         `gorm:"column:param_key;not null;uniqueIndex;size:50"` // 参数键名
	ParamValue string         `gorm:"column:param_value;type:text"`                  // 参数值
	Remark     string         `gorm:"column:remark;size:500"`                        // 备注
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SystemParam) TableName() string { return "system_params" }
