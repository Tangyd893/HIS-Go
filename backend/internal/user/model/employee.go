package model

import "time"

// Employee 员工模型，表名 employees
type Employee struct {
	ID           string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()" json:"id"`
	UserID       string    `gorm:"column:user_id;size:64" json:"userId,omitempty"`
	Name         string    `gorm:"column:name;not null;size:50" json:"name"`
	Phone        string    `gorm:"column:phone;size:20" json:"phone,omitempty"`
	DeptID       string    `gorm:"column:dept_id;size:64" json:"deptId,omitempty"`
	DeptName     string    `gorm:"-" json:"deptName,omitempty"`
	Title        string    `gorm:"column:title;size:50" json:"title,omitempty"`
	Specialty    string    `gorm:"column:specialty;size:255" json:"specialty,omitempty"`
	Introduction string    `gorm:"column:introduction;type:text" json:"introduction,omitempty"`
	Status       int8      `gorm:"column:status;default:1" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
}

func (Employee) TableName() string { return "employees" }
