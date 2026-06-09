package model

import "time"

// RecordTemplate 病历模板模型，表名 record_templates
type RecordTemplate struct {
	ID        string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"column:name;not null;size:200" json:"name"`
	DeptID    string    `gorm:"column:dept_id;not null;type:varchar(64);index" json:"deptId"`
	Content   string    `gorm:"column:content;type:text" json:"content,omitempty"`
	Type      int       `gorm:"column:type;default:0;comment:'模板类型'" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (RecordTemplate) TableName() string { return "record_templates" }
