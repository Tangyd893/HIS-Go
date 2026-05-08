package model

import "time"

// Department 科室模型，表名 departments
type Department struct {
	ID          string       `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	Name        string       `gorm:"column:name;not null;size:100"`
	ParentID    string       `gorm:"column:parent_id;size:64"`
	Description string       `gorm:"column:description;size:255"`
	SortOrder   int          `gorm:"column:sort_order;default:0"`
	CreatedAt   time.Time    `gorm:"column:created_at;autoCreateTime"`
	Children    []Department `gorm:"-" json:"children,omitempty"`
}

func (Department) TableName() string { return "departments" }
