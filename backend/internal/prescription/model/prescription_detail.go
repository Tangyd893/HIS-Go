package model

// PrescriptionDetail 处方明细表
type PrescriptionDetail struct {
	ID             string  `gorm:"column:id;primaryKey;type:varchar(64)"`
	PrescriptionID string  `gorm:"column:prescription_id;not null;type:varchar(64);index"`
	DrugID         string  `gorm:"column:drug_id;not null;type:varchar(64)"`
	DrugName       string  `gorm:"column:drug_name;size:200"`
	Specification  string  `gorm:"column:specification;size:100"`                  // 规格
	Dosage         float64 `gorm:"column:dosage;default:0"`                        // 剂量
	Usage          string  `gorm:"column:usage;size:100"`                          // 用法
	Frequency      string  `gorm:"column:frequency;size:50"`                       // 频次
	Days           int     `gorm:"column:days;default:1"`                          // 天数
	Quantity       int     `gorm:"column:quantity;default:0"`                      // 数量
	UnitPrice      float64 `gorm:"column:unit_price;type:decimal(10,2);default:0"` // 单价
	Note           string  `gorm:"column:note;type:text"`                          // 备注
}

func (PrescriptionDetail) TableName() string { return "prescription_details" }
