package model

import (
	"time"

	"gorm.io/gorm"
)

// RecordAuthorization 档案授权
type RecordAuthorization struct {
	ID         string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID  string         `gorm:"column:patient_id;not null;type:varchar(64);index"`
	DoctorID   string         `gorm:"column:doctor_id;not null;type:varchar(64);index"`
	AuthTime   string         `gorm:"column:auth_time;size:20"`   // 授权时间 YYYY-MM-DD HH:mm:ss
	ExpireTime string         `gorm:"column:expire_time;size:20"` // 过期时间 YYYY-MM-DD HH:mm:ss
	Status     int8           `gorm:"column:status;default:1"`    // 0已撤销 1已授权 2已过期
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (RecordAuthorization) TableName() string { return "record_authorizations" }
