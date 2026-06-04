package model

// PrescriptionDetail 处方明细表
type PrescriptionDetail struct {
	ID             string  `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	PrescriptionID string  `gorm:"column:prescription_id;not null;type:varchar(64);index" json:"prescriptionId,omitempty"`
	DrugID         string  `gorm:"column:drug_id;not null;type:varchar(64)" json:"drugId"`
	DrugName       string  `gorm:"column:drug_name;size:200" json:"drugName"`
	Specification  string  `gorm:"column:specification;size:100" json:"specification,omitempty"`        // 规格
	Dosage         float64 `gorm:"column:dosage;default:0" json:"dosage,omitempty"`                     // 剂量
	Usage          string  `gorm:"column:usage;size:100" json:"usage,omitempty"`                        // 用法
	Frequency      string  `gorm:"column:frequency;size:50" json:"frequency,omitempty"`                 // 频次
	Days           int     `gorm:"column:days;default:1" json:"days"`                                  // 天数
	Quantity       int     `gorm:"column:quantity;default:0" json:"quantity"`                           // 数量
	UnitPrice      float64 `gorm:"column:unit_price;type:decimal(10,2);default:0" json:"unitPrice,omitempty"` // 单价
	Note           string  `gorm:"column:note;type:text" json:"note,omitempty"`                         // 备注
}

func (PrescriptionDetail) TableName() string { return "prescription_details" }
