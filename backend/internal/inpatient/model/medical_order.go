package model

import "time"

// MedicalOrder 医嘱模型，表名 medical_orders
type MedicalOrder struct {
	ID          string     `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	InpatientID string     `gorm:"column:inpatient_id;not null;type:varchar(64);index"`
	DoctorID    string     `gorm:"column:doctor_id;not null;type:varchar(64)"`
	OrderType   int8       `gorm:"column:order_type;not null;comment:'1长期 2临时'"`
	Content     string     `gorm:"column:content;not null;type:text"`
	StartTime   *time.Time `gorm:"column:start_time;comment:'开始执行时间'"`
	EndTime     *time.Time `gorm:"column:end_time;comment:'停止执行时间'"`
	Status      int8       `gorm:"column:status;default:1;comment:'0停止 1执行中'"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (MedicalOrder) TableName() string { return "medical_orders" }
