package model

import (
	"time"

	"gorm.io/gorm"
)

// DictType 字典类型
type DictType struct {
	ID        string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	DictName  string         `gorm:"column:dict_name;not null;size:100"`            // 字典名称
	DictType  string         `gorm:"column:dict_type;not null;uniqueIndex;size:50"` // 字典类型标识
	Status    int8           `gorm:"column:status;default:1"`                       // 0禁用 1启用
	Remark    string         `gorm:"column:remark;size:500"`                        // 备注
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (DictType) TableName() string { return "dict_types" }
