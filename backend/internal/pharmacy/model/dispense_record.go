package model

import "time"

// DispenseRecord 发药记录模型，表名 dispense_records
type DispenseRecord struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PrescriptionID string    `gorm:"column:prescription_id;not null;type:varchar(64);index"`
	PatientID      string    `gorm:"column:patient_id;not null;type:varchar(64)"`
	DrugID         string    `gorm:"column:drug_id;not null;type:varchar(64);index"`
	Quantity       int       `gorm:"column:quantity;not null"`
	DispenserID    string    `gorm:"column:dispenser_id;not null;type:varchar(64)"`
	CheckerID      string    `gorm:"column:checker_id;type:varchar(64)"`
	Status         int8      `gorm:"column:status;default:1;comment:'0已退药 1已发药'"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (DispenseRecord) TableName() string { return "dispense_records" }
