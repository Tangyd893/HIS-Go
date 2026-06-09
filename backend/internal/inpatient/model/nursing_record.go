package model

import "time"

// NursingRecord 护理记录模型，表名 nursing_records
type NursingRecord struct {
	ID          string     `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	InpatientID string     `gorm:"column:inpatient_id;not null;type:varchar(64);index" json:"inpatientId"`
	NurseID     string     `gorm:"column:nurse_id;not null;type:varchar(64)" json:"nurseId"`
	RecordTime  *time.Time `gorm:"column:record_time" json:"recordTime,omitempty"`
	Content     string     `gorm:"column:content;not null;type:text" json:"content"`
	VitalSigns  string     `gorm:"column:vital_signs;type:text" json:"vitalSigns,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (NursingRecord) TableName() string { return "nursing_records" }
