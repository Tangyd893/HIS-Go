package model

import (
	"time"

	"gorm.io/gorm"
)

// DictType 字典类型
type DictType struct {
	ID        string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	DictName  string         `gorm:"column:dict_name;not null;size:100" json:"dictLabel"`
	DictType  string         `gorm:"column:dict_type;not null;uniqueIndex;size:50" json:"dictType"`
	Status    int8           `gorm:"column:status;default:1" json:"status"`
	Remark    string         `gorm:"column:remark;size:500" json:"remark,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (DictType) TableName() string { return "dict_types" }
