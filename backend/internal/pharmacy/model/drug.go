package model

import "time"

// Drug 药品模型，表名 drugs
type Drug struct {
	ID            string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	Name          string    `gorm:"column:name;not null;size:200" json:"name"`
	GenericName   string    `gorm:"column:generic_name;size:200" json:"genericName,omitempty"`
	Specification string    `gorm:"column:specification;size:100" json:"specification,omitempty"`
	Manufacturer  string    `gorm:"column:manufacturer;size:200" json:"manufacturer,omitempty"`
	BatchNo       string    `gorm:"column:batch_no;size:100" json:"batchNo,omitempty"`
	PurchasePrice float64   `gorm:"column:purchase_price;not null;default:0" json:"purchasePrice,omitempty"`
	RetailPrice   float64   `gorm:"column:retail_price;not null;default:0" json:"price"`
	Stock         int       `gorm:"column:stock;not null;default:0" json:"stock"`
	MinStock      int       `gorm:"column:min_stock;not null;default:0" json:"minStock,omitempty"`
	ExpiryDate    string    `gorm:"column:expiry_date;size:20" json:"expiryDate,omitempty"`
	Status        int8      `gorm:"column:status;default:1" json:"status,omitempty"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
}

func (Drug) TableName() string { return "drugs" }
