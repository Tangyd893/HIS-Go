package user

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/user"
	usermodel "his-go/internal/user/model"
	usersvc "his-go/internal/user/service"
)

// UserGrpcServer gRPC 用户管理服务实现
type UserGrpcServer struct {
	user.UnimplementedUserServiceServer
	svc *usersvc.UserService
}

// NewUserGrpcServer 创建 gRPC 用户管理服务
func NewUserGrpcServer(svc *usersvc.UserService) *UserGrpcServer {
	return &UserGrpcServer{svc: svc}
}

// GetPatient 获取患者详情
func (s *UserGrpcServer) GetPatient(ctx context.Context, req *common.IdRequest) (*user.PatientInfo, error) {
	p, err := s.svc.GetPatient(req.Id)
	if err != nil {
		return nil, err
	}
	return patientToProto(p), nil
}

// ListPatients 分页查询患者列表
func (s *UserGrpcServer) ListPatients(ctx context.Context, req *user.PatientListRequest) (*user.PatientListResponse, error) {
	page := 1
	pageSize := 10
	if req.Page != nil {
		if req.Page.Page > 0 {
			page = int(req.Page.Page)
		}
		if req.Page.PageSize > 0 {
			pageSize = int(req.Page.PageSize)
		}
	}
	patients, total, err := s.svc.ListPatients(req.Name, req.Phone, page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*user.PatientInfo, len(patients))
	for i, p := range patients {
		pbList[i] = patientToProto(&p)
	}
	return &user.PatientListResponse{
		Base:     &common.BaseResponse{Code: 0, Message: "查询成功"},
		Patients: pbList,
		Page:     &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// CreatePatient 创建患者
func (s *UserGrpcServer) CreatePatient(ctx context.Context, req *user.PatientInfo) (*user.PatientInfo, error) {
	p := protoToPatient(req)
	if err := s.svc.CreatePatient(p); err != nil {
		return nil, err
	}
	return patientToProto(p), nil
}

// UpdatePatient 更新患者信息
func (s *UserGrpcServer) UpdatePatient(ctx context.Context, req *user.PatientInfo) (*user.PatientInfo, error) {
	p := protoToPatient(req)
	if err := s.svc.UpdatePatient(p); err != nil {
		return nil, err
	}
	return patientToProto(p), nil
}

// ListDepartments 查询科室列表
func (s *UserGrpcServer) ListDepartments(ctx context.Context, req *common.Empty) (*user.DepartmentListResponse, error) {
	depts, err := s.svc.ListDepartments()
	if err != nil {
		return nil, err
	}
	pbList := make([]*user.DepartmentInfo, len(depts))
	for i, d := range depts {
		pbList[i] = &user.DepartmentInfo{
			Id:          d.ID,
			Name:        d.Name,
			ParentId:    d.ParentID,
			Description: d.Description,
			SortOrder:   int32(d.SortOrder),
		}
	}
	return &user.DepartmentListResponse{
		Base:        &common.BaseResponse{Code: 0, Message: "查询成功"},
		Departments: pbList,
	}, nil
}

// GetEmployee 获取员工详情
func (s *UserGrpcServer) GetEmployee(ctx context.Context, req *common.IdRequest) (*user.EmployeeInfo, error) {
	e, err := s.svc.GetEmployee(req.Id)
	if err != nil {
		return nil, err
	}
	return employeeToProto(e), nil
}

// ListEmployees 分页查询员工列表
func (s *UserGrpcServer) ListEmployees(ctx context.Context, req *user.EmployeeListRequest) (*user.EmployeeListResponse, error) {
	page := 1
	pageSize := 10
	if req.Page != nil {
		if req.Page.Page > 0 {
			page = int(req.Page.Page)
		}
		if req.Page.PageSize > 0 {
			pageSize = int(req.Page.PageSize)
		}
	}
	employees, total, err := s.svc.ListEmployees(req.DeptId, req.Name, page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*user.EmployeeInfo, len(employees))
	for i, e := range employees {
		pbList[i] = employeeToProto(&e)
	}
	return &user.EmployeeListResponse{
		Base:      &common.BaseResponse{Code: 0, Message: "查询成功"},
		Employees: pbList,
		Page:      &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// ---- 转换辅助函数 ----

func patientToProto(p *usermodel.Patient) *user.PatientInfo {
	return &user.PatientInfo{
		Id:             p.ID,
		Name:           p.Name,
		IdCard:         p.IdCard,
		Phone:          p.Phone,
		Gender:         p.Gender,
		Address:        p.Address,
		AllergyHistory: p.AllergyHistory,
		CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func protoToPatient(pb *user.PatientInfo) *usermodel.Patient {
	return &usermodel.Patient{
		ID:             pb.Id,
		Name:           pb.Name,
		IdCard:         pb.IdCard,
		Phone:          pb.Phone,
		Gender:         pb.Gender,
		Address:        pb.Address,
		AllergyHistory: pb.AllergyHistory,
	}
}

func employeeToProto(e *usermodel.Employee) *user.EmployeeInfo {
	return &user.EmployeeInfo{
		Id:           e.ID,
		Name:         e.Name,
		Phone:        e.Phone,
		DeptId:       e.DeptID,
		Title:        e.Title,
		Specialty:    e.Specialty,
		Introduction: e.Introduction,
	}
}
