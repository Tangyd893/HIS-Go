package model

import "time"

// RecordTemplate 病历模板模型，表名 record_templates
type RecordTemplate struct {
	ID        string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	Name      string    `gorm:"column:name;not null;size:200"`
	DeptID    string    `gorm:"column:dept_id;not null;type:varchar(64);index"`
	Content   string    `gorm:"column:content;type:text"`
	Type      int       `gorm:"column:type;default:0;comment:'模板类型'"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (RecordTemplate) TableName() string { return "record_templates" }
