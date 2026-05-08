package model

import (
	"time"

	"gorm.io/gorm"
)

// Patient 患者模型，表名 patients
type Patient struct {
	ID             string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	Name           string         `gorm:"column:name;not null;size:50"`
	IdCard         string         `gorm:"column:id_card;uniqueIndex;not null;size:18"`
	Phone          string         `gorm:"column:phone;size:20"`
	Gender         string         `gorm:"column:gender;size:10"`
	BirthDate      *time.Time     `gorm:"column:birth_date;type:date"`
	Address        string         `gorm:"column:address;size:255"`
	AllergyHistory string         `gorm:"column:allergy_history;type:text"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Patient) TableName() string { return "patients" }
