package model

import "time"

// BillDetail 账单明细模型，表名 bill_details
type BillDetail struct {
	ID        string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	BillID    string    `gorm:"column:bill_id;not null;type:varchar(64);index" json:"billId"`
	ItemType  int8      `gorm:"column:item_type;not null;comment:'1挂号费 2检查费 3药费 4治疗费'" json:"itemType"`
	ItemName  string    `gorm:"column:item_name;not null;size:200" json:"itemName"`
	UnitPrice float64   `gorm:"column:unit_price;not null" json:"unitPrice"`
	Quantity  int       `gorm:"column:quantity;not null;default:1" json:"quantity"`
	Amount    float64   `gorm:"column:amount;not null" json:"amount"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (BillDetail) TableName() string { return "bill_details" }
