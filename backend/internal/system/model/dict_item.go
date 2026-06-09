package model

import (
	"time"

	"gorm.io/gorm"
)

// DictItem 字典项
type DictItem struct {
	ID        string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	DictType  string         `gorm:"column:dict_type;not null;size:50;index" json:"dictType"`
	Label     string         `gorm:"column:label;not null;size:100" json:"dictLabel"`
	Value     string         `gorm:"column:value;not null;size:100" json:"dictValue"`
	SortOrder int            `gorm:"column:sort_order;default:0" json:"sort"`
	Status    int8           `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (DictItem) TableName() string { return "dict_items" }
