package registration

import (
	"context"
	"fmt"
	"time"

	"his-go/api/proto/common"
	"his-go/api/proto/registration"
	regmodel "his-go/internal/registration/model"
	regsvc "his-go/internal/registration/service"
)

// RegistrationGrpcServer gRPC 挂号服务实现
type RegistrationGrpcServer struct {
	registration.UnimplementedRegistrationServiceServer
	svc *regsvc.RegistrationService
}

// NewRegistrationGrpcServer 创建 gRPC 挂号服务
func NewRegistrationGrpcServer(svc *regsvc.RegistrationService) *RegistrationGrpcServer {
	return &RegistrationGrpcServer{svc: svc}
}

// ListSchedules 查询号源列表
func (s *RegistrationGrpcServer) ListSchedules(ctx context.Context, req *registration.ScheduleQueryRequest) (*registration.ScheduleListResponse, error) {
	schedules, err := s.svc.ListSchedules(req.DeptId, req.Date)
	if err != nil {
		return nil, err
	}
	pbList := make([]*registration.ScheduleInfo, len(schedules))
	for i, sc := range schedules {
		pbList[i] = scheduleToProto(&sc)
	}
	return &registration.ScheduleListResponse{
		Base:      &common.BaseResponse{Code: 0, Message: "查询成功"},
		Schedules: pbList,
	}, nil
}

// Register 挂号
// 注意: proto RegistrationRequest 当前缺少 patient_name 和 registration_date 字段，
// 待 proto 补齐后从请求中获取
func (s *RegistrationGrpcServer) Register(ctx context.Context, req *registration.RegistrationRequest) (*registration.RegistrationInfo, error) {
	today := time.Now().Format("2006-01-02")
	reg, err := s.svc.Register(req.PatientId, "", req.ScheduleId, today)
	if err != nil {
		return nil, err
	}
	return registrationToProto(reg), nil
}

// CancelRegistration 取消挂号
func (s *RegistrationGrpcServer) CancelRegistration(ctx context.Context, req *common.IdRequest) (*common.BaseResponse, error) {
	if err := s.svc.Cancel(req.Id); err != nil {
		return nil, err
	}
	return &common.BaseResponse{Code: 0, Message: "取消成功"}, nil
}

// SignIn 签到
func (s *RegistrationGrpcServer) SignIn(ctx context.Context, req *common.IdRequest) (*common.BaseResponse, error) {
	if err := s.svc.SignIn(req.Id); err != nil {
		return nil, err
	}
	return &common.BaseResponse{Code: 0, Message: "签到成功"}, nil
}

// GetQueueStatus 查询排队状态
func (s *RegistrationGrpcServer) GetQueueStatus(ctx context.Context, req *common.IdRequest) (*registration.QueueStatus, error) {
	rank, err := s.svc.GetQueueStatus(req.Id)
	if err != nil {
		return nil, err
	}

	reg, err := s.svc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &registration.QueueStatus{
		RegistrationId: req.Id,
		QueueNumber:    int32(reg.QueueNumber),
		CurrentNumber:  0,
		WaitCount:      int32(rank),
		Status:         "waiting",
	}, nil
}

// CallNext 叫号
func (s *RegistrationGrpcServer) CallNext(ctx context.Context, req *common.IdRequest) (*registration.RegistrationInfo, error) {
	memberID, err := s.svc.CallNext(req.Id)
	if err != nil {
		return nil, err
	}

	reg, err := s.svc.GetByID(memberID)
	if err != nil {
		return nil, err
	}
	return registrationToProto(reg), nil
}

// ---- 转换辅助函数 ----

func scheduleToProto(sc *regmodel.Schedule) *registration.ScheduleInfo {
	return &registration.ScheduleInfo{
		Id:          sc.ID,
		DeptId:      sc.DeptID,
		DeptName:    sc.DeptName,
		DoctorId:    sc.DoctorID,
		DoctorName:  sc.DoctorName,
		Date:        sc.Date,
		TimeSlot:    fmt.Sprintf("%d", sc.TimeSlot),
		TotalCount:  int32(sc.TotalCount),
		RemainCount: int32(sc.RemainCount),
		Fee:         sc.Fee,
	}
}

func registrationToProto(reg *regmodel.Registration) *registration.RegistrationInfo {
	return &registration.RegistrationInfo{
		Id:               reg.ID,
		PatientId:        reg.PatientID,
		PatientName:      reg.PatientName,
		ScheduleId:       reg.ScheduleID,
		RegistrationDate: reg.RegistrationDate,
		QueueNumber:      int32(reg.QueueNumber),
		Status:           int32(reg.Status),
		CreatedAt:        reg.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
