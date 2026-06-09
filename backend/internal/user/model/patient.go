package model

import (
	"time"

	"gorm.io/gorm"
)

// Patient 患者模型，表名 patients
type Patient struct {
	ID             string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	Name           string         `gorm:"column:name;not null;size:50" json:"name"`
	IdCard         string         `gorm:"column:id_card;uniqueIndex;not null;size:18" json:"idCard"`
	Phone          string         `gorm:"column:phone;size:20" json:"phone,omitempty"`
	Gender         string         `gorm:"column:gender;size:10" json:"gender,omitempty"`
	BirthDate      *time.Time     `gorm:"column:birth_date;type:date" json:"birthDate,omitempty"`
	Address        string         `gorm:"column:address;size:255" json:"address,omitempty"`
	AllergyHistory string         `gorm:"column:allergy_history;type:text" json:"allergyHistory,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Patient) TableName() string { return "patients" }
