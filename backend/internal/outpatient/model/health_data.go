package model

import (
	"time"

	"gorm.io/gorm"
)

// HealthData 健康数据
type HealthData struct {
	ID          string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PatientID   string         `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	DataType    string         `gorm:"column:data_type;not null;size:30" json:"dataType"` // blood_pressure/blood_sugar/weight
	Value       string         `gorm:"column:value;size:20" json:"value"`                  // 测量值
	Unit        string         `gorm:"column:unit;size:10" json:"unit,omitempty"`           // 单位：mmHg/mmol/L/kg
	MeasureTime string         `gorm:"column:measure_time;size:20" json:"measureTime"`      // 测量时间 YYYY-MM-DD HH:mm:ss
	Abnormal    bool           `gorm:"column:abnormal;default:false" json:"abnormal"`       // 是否异常
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (HealthData) TableName() string { return "health_data" }
