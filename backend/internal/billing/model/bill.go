package model

import "time"

// Bill 账单模型，表名 bills
type Bill struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PatientID      string    `gorm:"column:patient_id;not null;type:varchar(64);index"`
	RegistrationID string    `gorm:"column:registration_id;type:varchar(64)"`
	BillNo         string    `gorm:"column:bill_no;uniqueIndex;not null;size:64"`
	TotalAmount    float64   `gorm:"column:total_amount;not null"`
	PaidAmount     float64   `gorm:"column:paid_amount;not null;default:0"`
	PayMethod      int8      `gorm:"column:pay_method;default:0;comment:'1现金 2微信 3支付宝 4医保'"`
	Status         int8      `gorm:"column:status;default:0;comment:'0待支付 1已支付 2已退款 3已取消'"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Bill) TableName() string { return "bills" }
