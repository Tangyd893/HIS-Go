package model

import "time"

// DispenseRecord 发药记录模型，表名 dispense_records
type DispenseRecord struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PrescriptionID string    `gorm:"column:prescription_id;not null;type:varchar(64);index" json:"prescriptionId"`
	PatientID      string    `gorm:"column:patient_id;not null;type:varchar(64)" json:"patientId"`
	DrugID         string    `gorm:"column:drug_id;not null;type:varchar(64);index" json:"drugId"`
	Quantity       int       `gorm:"column:quantity;not null" json:"quantity"`
	DispenserID    string    `gorm:"column:dispenser_id;not null;type:varchar(64)" json:"dispenserId"`
	CheckerID      string    `gorm:"column:checker_id;type:varchar(64)" json:"checkerId,omitempty"`
	Status         int8      `gorm:"column:status;default:1" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (DispenseRecord) TableName() string { return "dispense_records" }
