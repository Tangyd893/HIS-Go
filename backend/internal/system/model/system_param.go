package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemParam 系统参数
type SystemParam struct {
	ID         string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	ParamName  string         `gorm:"column:param_name;not null;size:100" json:"paramName"`
	ParamKey   string         `gorm:"column:param_key;not null;uniqueIndex;size:50" json:"paramKey"`
	ParamValue string         `gorm:"column:param_value;type:text" json:"paramValue"`
	Remark     string         `gorm:"column:remark;size:500" json:"description"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (SystemParam) TableName() string { return "system_params" }
