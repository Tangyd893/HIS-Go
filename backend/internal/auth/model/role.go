package model

import "time"

// Role 角色模型，表名 roles
type Role struct {
	ID          string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	RoleName    string    `gorm:"column:role_name;not null;size:50"`
	RoleCode    string    `gorm:"column:role_code;uniqueIndex;not null;size:50"`
	Description string    `gorm:"column:description;size:255"`
	Status      int8      `gorm:"column:status;default:1"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Role) TableName() string { return "roles" }

// Permission 权限模型，表名 permissions
type Permission struct {
	ID        string    `gorm:"column:id;primaryKey;type:varchar(64);default:gen_random_uuid()"`
	PermName  string    `gorm:"column:perm_name;not null;size:50"`
	PermCode  string    `gorm:"column:perm_code;uniqueIndex;not null;size:50"`
	ParentID  string    `gorm:"column:parent_id;size:64"`
	SortOrder int       `gorm:"column:sort_order;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Permission) TableName() string { return "permissions" }

// RolePermission 角色权限关联模型，表名 role_permissions
type RolePermission struct {
	RoleID string `gorm:"column:role_id;primaryKey;type:varchar(64)"`
	PermID string `gorm:"column:perm_id;primaryKey;type:varchar(64)"`
}

func (RolePermission) TableName() string { return "role_permissions" }
