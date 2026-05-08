package schedule

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/schedule"
	schedmodel "his-go/internal/schedule/model"
	schedsvc "his-go/internal/schedule/service"
)

// ScheduleGrpcServer gRPC 排班管理服务实现
type ScheduleGrpcServer struct {
	schedule.UnimplementedScheduleServiceServer
	svc *schedsvc.ScheduleService
}

// NewScheduleGrpcServer 创建 gRPC 排班管理服务
func NewScheduleGrpcServer(svc *schedsvc.ScheduleService) *ScheduleGrpcServer {
	return &ScheduleGrpcServer{svc: svc}
}

// GenerateWeeklySchedules 生成排班
func (s *ScheduleGrpcServer) GenerateWeeklySchedules(ctx context.Context, req *schedule.GenerateRequest) (*schedule.ScheduleListResponse, error) {
	schedules, err := s.svc.GenerateWeekly(req.StartDate, req.EndDate, "")
	if err != nil {
		return nil, err
	}
	pbList := make([]*schedule.ScheduleInfo, len(schedules))
	for i, sc := range schedules {
		pbList[i] = scheduleToProto(&sc)
	}
	return &schedule.ScheduleListResponse{
		Base:      &common.BaseResponse{Code: 0, Message: "排班生成成功"},
		Schedules: pbList,
	}, nil
}

// ListSchedules 查询排班列表
func (s *ScheduleGrpcServer) ListSchedules(ctx context.Context, req *schedule.ScheduleQueryRequest) (*schedule.ScheduleListResponse, error) {
	var schedules []schedmodel.ScheduleInfo
	var err error
	if req.DoctorId != "" {
		schedules, err = s.svc.ListByDoctor(req.DoctorId, req.StartDate)
	} else {
		schedules, err = s.svc.ListByDept(req.DeptId, req.StartDate)
	}
	if err != nil {
		return nil, err
	}
	pbList := make([]*schedule.ScheduleInfo, len(schedules))
	for i, sc := range schedules {
		pbList[i] = scheduleToProto(&sc)
	}
	return &schedule.ScheduleListResponse{
		Base:      &common.BaseResponse{Code: 0, Message: "查询成功"},
		Schedules: pbList,
	}, nil
}

// UpdateSchedule 更新排班
func (s *ScheduleGrpcServer) UpdateSchedule(ctx context.Context, req *schedule.ScheduleInfo) (*schedule.ScheduleInfo, error) {
	sc := protoToSchedule(req)
	if err := s.svc.UpdateSchedule(sc); err != nil {
		return nil, err
	}
	return scheduleToProto(sc), nil
}

// CancelSchedule 取消排班
func (s *ScheduleGrpcServer) CancelSchedule(ctx context.Context, req *common.IdRequest) (*common.BaseResponse, error) {
	if err := s.svc.CancelSchedule(req.Id); err != nil {
		return nil, err
	}
	return &common.BaseResponse{Code: 0, Message: "取消成功"}, nil
}

// ---- 转换辅助函数 ----

func scheduleToProto(sc *schedmodel.ScheduleInfo) *schedule.ScheduleInfo {
	return &schedule.ScheduleInfo{
		Id:              sc.ID,
		DoctorId:        sc.DoctorID,
		DoctorName:      sc.DoctorName,
		DeptId:          sc.DeptID,
		DeptName:        sc.DeptName,
		WorkDate:        sc.WorkDate,
		TimeSlot:        int32(sc.TimeSlot),
		MaxPatients:     int32(sc.MaxPatients),
		CurrentPatients: int32(sc.CurrentPatients),
		RoomNo:          sc.RoomNo,
		Status:          int32(sc.Status),
	}
}

func protoToSchedule(pb *schedule.ScheduleInfo) *schedmodel.ScheduleInfo {
	return &schedmodel.ScheduleInfo{
		ID:              pb.Id,
		DoctorID:        pb.DoctorId,
		DoctorName:      pb.DoctorName,
		DeptID:          pb.DeptId,
		DeptName:        pb.DeptName,
		WorkDate:        pb.WorkDate,
		TimeSlot:        int(pb.TimeSlot),
		MaxPatients:     int(pb.MaxPatients),
		CurrentPatients: int(pb.CurrentPatients),
		RoomNo:          pb.RoomNo,
		Status:          int8(pb.Status),
	}
}
