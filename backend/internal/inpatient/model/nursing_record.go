package model

import "time"

// NursingRecord 护理记录模型，表名 nursing_records
type NursingRecord struct {
	ID          string     `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	InpatientID string     `gorm:"column:inpatient_id;not null;type:varchar(64);index"`
	NurseID     string     `gorm:"column:nurse_id;not null;type:varchar(64)"`
	RecordTime  *time.Time `gorm:"column:record_time;comment:'记录时间'"`
	Content     string     `gorm:"column:content;not null;type:text"`
	VitalSigns  string     `gorm:"column:vital_signs;type:text;comment:'生命体征JSON'"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (NursingRecord) TableName() string { return "nursing_records" }
