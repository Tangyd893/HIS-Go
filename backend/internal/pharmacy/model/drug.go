package model

import "time"

// Drug 药品模型，表名 drugs
type Drug struct {
	ID            string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	Name          string    `gorm:"column:name;not null;size:200"`
	GenericName   string    `gorm:"column:generic_name;size:200"`
	Specification string    `gorm:"column:specification;size:100"`
	Manufacturer  string    `gorm:"column:manufacturer;size:200"`
	BatchNo       string    `gorm:"column:batch_no;size:100"`
	PurchasePrice float64   `gorm:"column:purchase_price;not null;default:0"`
	RetailPrice   float64   `gorm:"column:retail_price;not null;default:0"`
	Stock         int       `gorm:"column:stock;not null;default:0"`
	MinStock      int       `gorm:"column:min_stock;not null;default:0"`
	ExpiryDate    string    `gorm:"column:expiry_date;size:20;comment:'过期日期'"`
	Status        int8      `gorm:"column:status;default:1;comment:'0停用 1启用'"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Drug) TableName() string { return "drugs" }
