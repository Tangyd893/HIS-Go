package model

import (
	"time"

	"gorm.io/gorm"
)

// DictItem 字典项
type DictItem struct {
	ID        string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	DictType  string         `gorm:"column:dict_type;not null;size:50;index"` // 关联字典类型
	Label     string         `gorm:"column:label;not null;size:100"`          // 字典标签
	Value     string         `gorm:"column:value;not null;size:100"`          // 字典值
	SortOrder int            `gorm:"column:sort_order;default:0"`             // 排序号
	Status    int8           `gorm:"column:status;default:1"`                 // 0禁用 1启用
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (DictItem) TableName() string { return "dict_items" }
