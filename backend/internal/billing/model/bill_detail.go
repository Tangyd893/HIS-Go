package model

import "time"

// BillDetail 账单明细模型，表名 bill_details
type BillDetail struct {
	ID        string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	BillID    string    `gorm:"column:bill_id;not null;type:varchar(64);index"`
	ItemType  int8      `gorm:"column:item_type;not null;comment:'1挂号费 2检查费 3药费 4治疗费'"`
	ItemName  string    `gorm:"column:item_name;not null;size:200"`
	UnitPrice float64   `gorm:"column:unit_price;not null"`
	Quantity  int       `gorm:"column:quantity;not null;default:1"`
	Amount    float64   `gorm:"column:amount;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (BillDetail) TableName() string { return "bill_details" }
