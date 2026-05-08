package repository

import (
	"gorm.io/gorm"

	"his-go/internal/user/model"
)

// DepartmentRepository 科室数据仓库
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository 创建科室数据仓库
func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

// ListAll 查询所有科室
func (r *DepartmentRepository) ListAll() ([]model.Department, error) {
	var depts []model.Department
	err := r.db.Order("sort_order ASC, created_at ASC").Find(&depts).Error
	if err != nil {
		return nil, err
	}
	return depts, nil
}

// FindByID 根据ID查找科室
func (r *DepartmentRepository) FindByID(id string) (*model.Department, error) {
	var dept model.Department
	err := r.db.Where("id = ?", id).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// Create 创建科室
func (r *DepartmentRepository) Create(dept *model.Department) error {
	return r.db.Create(dept).Error
}

// Update 更新科室信息
func (r *DepartmentRepository) Update(dept *model.Department) error {
	return r.db.Save(dept).Error
}
