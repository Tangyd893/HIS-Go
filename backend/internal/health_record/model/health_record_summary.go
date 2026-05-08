package model

import (
	"time"

	"gorm.io/gorm"
)

// HealthRecordSummary 健康档案摘要
type HealthRecordSummary struct {
	ID                 string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID          string         `gorm:"column:patient_id;not null;type:varchar(64);uniqueIndex"`
	PatientName        string         `gorm:"column:patient_name;size:50"`
	TotalVisits        int            `gorm:"column:total_visits;default:0"`        // 就诊总次数
	TotalPrescriptions int            `gorm:"column:total_prescriptions;default:0"` // 处方总数
	TotalExaminations  int            `gorm:"column:total_examinations;default:0"`  // 检查总数
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (HealthRecordSummary) TableName() string { return "health_record_summaries" }
