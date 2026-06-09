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

// List 分页查询员工列表（关联科室名称）
func (r *EmployeeRepository) List(deptID, name string, page, pageSize int) ([]model.Employee, int64, error) {
	var total int64

	query := r.db.Table("employees").
		Joins("LEFT JOIN departments ON departments.id = employees.dept_id")
	if deptID != "" {
		query = query.Where("employees.dept_id = ?", deptID)
	}
	if name != "" {
		query = query.Where("employees.name LIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		model.Employee
		DeptName string `gorm:"column:dept_name"`
	}
	var rows []row
	offset := (page - 1) * pageSize
	if err := query.
		Select("employees.*, departments.name AS dept_name").
		Offset(offset).Limit(pageSize).
		Order("employees.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	employees := make([]model.Employee, len(rows))
	for i, item := range rows {
		employees[i] = item.Employee
		employees[i].DeptName = item.DeptName
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
