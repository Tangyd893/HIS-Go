package model

import "time"

// Bill 账单模型，表名 bills
type Bill struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	PatientID      string    `gorm:"column:patient_id;not null;type:varchar(64);index" json:"patientId"`
	PatientName    string    `gorm:"-" json:"patientName,omitempty"`
	RegistrationID string    `gorm:"column:registration_id;type:varchar(64)" json:"registrationId,omitempty"`
	BillNo         string    `gorm:"column:bill_no;uniqueIndex;not null;size:64" json:"billNo"`
	TotalAmount    float64   `gorm:"column:total_amount;not null" json:"totalAmount"`
	PaidAmount     float64   `gorm:"column:paid_amount;not null;default:0" json:"paidAmount"`
	PayMethod      int8      `gorm:"column:pay_method;default:0" json:"payMethod,omitempty"`
	Status         int8         `gorm:"column:status;default:0" json:"status"`
	Details        []BillDetail `gorm:"-" json:"details,omitempty"`
	CreatedAt      time.Time    `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
}

func (Bill) TableName() string { return "bills" }
