package model

import "time"

// Employee 员工模型，表名 employees
type Employee struct {
	ID           string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	UserID       string    `gorm:"column:user_id;size:64"`
	Name         string    `gorm:"column:name;not null;size:50"`
	Phone        string    `gorm:"column:phone;size:20"`
	DeptID       string    `gorm:"column:dept_id;size:64"`
	Title        string    `gorm:"column:title;size:50"`
	Specialty    string    `gorm:"column:specialty;size:255"`
	Introduction string    `gorm:"column:introduction;type:text"`
	Status       int8      `gorm:"column:status;default:1"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Employee) TableName() string { return "employees" }
