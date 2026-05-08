package repository

import (
	"gorm.io/gorm"

	"his-go/internal/user/model"
)

// EmployeeRepository 员工数据仓库
type EmployeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository 创建员工数据仓库
func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// FindByID 根据ID查找员工
func (r *EmployeeRepository) FindByID(id string) (*model.Employee, error) {
	var emp model.Employee
	err := r.db.Where("id = ?", id).First(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// List 分页查询员工列表
func (r *EmployeeRepository) List(deptID, name string, page, pageSize int) ([]model.Employee, int64, error) {
	var employees []model.Employee
	var total int64

	query := r.db.Model(&model.Employee{})
	if deptID != "" {
		query = query.Where("dept_id = ?", deptID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

// Create 创建员工
func (r *EmployeeRepository) Create(emp *model.Employee) error {
	return r.db.Create(emp).Error
}

// Update 更新员工信息
func (r *EmployeeRepository) Update(emp *model.Employee) error {
	return r.db.Save(emp).Error
}
