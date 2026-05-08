package repository

import (
	"time"

	"gorm.io/gorm"

	"his-go/internal/auth/model"
)

// AuthRepository 认证数据仓库
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository 创建认证数据仓库
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// FindByUsername 根据用户名查找用户
func (r *AuthRepository) FindByUsername(username string) (*model.AuthUser, error) {
	var user model.AuthUser
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *AuthRepository) FindByID(id string) (*model.AuthUser, error) {
	var user model.AuthUser
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin 更新用户最后登录时间
func (r *AuthRepository) UpdateLastLogin(userID string) error {
	now := time.Now()
	return r.db.Model(&model.AuthUser{}).Where("id = ?", userID).Update("last_login_time", now).Error
}

// FindRolesByUserID 根据用户ID查找所属角色列表
func (r *AuthRepository) FindRolesByUserID(userID string) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Table("roles").
		Select("roles.*").
		Joins("JOIN role_permissions rp ON roles.id = rp.role_id").
		Joins("JOIN users u ON u.role = roles.role_code").
		Where("u.id = ?", userID).
		Group("roles.id").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// FindPermsByRoleIDs 根据角色ID列表查找权限列表
func (r *AuthRepository) FindPermsByRoleIDs(roleIDs []string) ([]model.Permission, error) {
	if len(roleIDs) == 0 {
		return []model.Permission{}, nil
	}
	var perms []model.Permission
	err := r.db.Table("permissions").
		Select("DISTINCT permissions.*").
		Joins("JOIN role_permissions rp ON permissions.id = rp.perm_id").
		Where("rp.role_id IN ?", roleIDs).
		Find(&perms).Error
	if err != nil {
		return nil, err
	}
	return perms, nil
}
