package service

import (
	"time"

	"github.com/google/uuid"

	"his-go/internal/user/model"
	"his-go/internal/user/repository"
	"his-go/pkg/redis"
)

// UserService 用户管理服务
type UserService struct {
	patientRepo *repository.PatientRepository
	empRepo     *repository.EmployeeRepository
	deptRepo    *repository.DepartmentRepository
	rdb         *redis.Client
}

// NewUserService 创建用户管理服务
func NewUserService(
	patientRepo *repository.PatientRepository,
	empRepo *repository.EmployeeRepository,
	deptRepo *repository.DepartmentRepository,
	rdb *redis.Client,
) *UserService {
	return &UserService{
		patientRepo: patientRepo,
		empRepo:     empRepo,
		deptRepo:    deptRepo,
		rdb:         rdb,
	}
}

// ---- 患者管理 ----

// GetPatient 获取患者详情
func (s *UserService) GetPatient(id string) (*model.Patient, error) {
	return s.patientRepo.FindByID(id)
}

// ListPatients 分页查询患者列表
func (s *UserService) ListPatients(name, phone string, page, pageSize int) ([]model.Patient, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.patientRepo.List(name, phone, page, pageSize)
}

// CreatePatient 创建患者
func (s *UserService) CreatePatient(p *model.Patient) error {
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return s.patientRepo.Create(p)
}

// UpdatePatient 更新患者信息
func (s *UserService) UpdatePatient(p *model.Patient) error {
	p.UpdatedAt = time.Now()
	return s.patientRepo.Update(p)
}

// DeletePatient 逻辑删除患者
func (s *UserService) DeletePatient(id string) error {
	return s.patientRepo.Delete(id)
}

// ---- 员工管理 ----

// GetEmployee 获取员工详情
func (s *UserService) GetEmployee(id string) (*model.Employee, error) {
	return s.empRepo.FindByID(id)
}

// ListEmployees 分页查询员工列表
func (s *UserService) ListEmployees(deptID, name string, page, pageSize int) ([]model.Employee, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.empRepo.List(deptID, name, page, pageSize)
}

// ---- 科室管理 ----

// ListDepartments 查询科室列表（树形结构）
func (s *UserService) ListDepartments() ([]model.Department, error) {
	depts, err := s.deptRepo.ListAll()
	if err != nil {
		return nil, err
	}
	return buildDepartmentTree(depts, ""), nil
}

// buildDepartmentTree 构建科室树形结构
func buildDepartmentTree(depts []model.Department, parentID string) []model.Department {
	var result []model.Department
	for _, dept := range depts {
		if dept.ParentID == parentID {
			children := buildDepartmentTree(depts, dept.ID)
			if len(children) > 0 {
				dept.Children = children
			}
			result = append(result, dept)
		}
	}
	return result
}
