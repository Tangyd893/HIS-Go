package model

import (
	"time"

	"gorm.io/gorm"
)

// AuthUser 用户认证模型，表名 users
type AuthUser struct {
	ID            string         `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	Username      string         `gorm:"column:username;uniqueIndex;not null;size:50"`
	Password      string         `gorm:"column:password;not null;size:255"`
	RealName      string         `gorm:"column:real_name;size:50"`
	Phone         string         `gorm:"column:phone;size:20"`
	Email         string         `gorm:"column:email;size:100"`
	Avatar        string         `gorm:"column:avatar;size:255"`
	Role          string         `gorm:"column:role;size:20;default:doctor"`
	DeptID        string         `gorm:"column:dept_id;size:64"`
	Status        int8           `gorm:"column:status;default:1"`
	LastLoginTime *time.Time     `gorm:"column:last_login_time"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (AuthUser) TableName() string { return "users" }
