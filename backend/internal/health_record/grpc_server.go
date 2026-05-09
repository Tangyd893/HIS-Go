package health_record

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/health_record"
	hrmodel "his-go/internal/health_record/model"
	hrsvc "his-go/internal/health_record/service"
)

// HealthRecordGrpcServer gRPC 健康档案服务实现
type HealthRecordGrpcServer struct {
	health_record.UnimplementedHealthRecordServiceServer
	svc *hrsvc.HealthRecordService
}

// NewHealthRecordGrpcServer 创建 gRPC 健康档案服务
func NewHealthRecordGrpcServer(svc *hrsvc.HealthRecordService) *HealthRecordGrpcServer {
	return &HealthRecordGrpcServer{svc: svc}
}

// GetSummary 获取健康档案总览
func (s *HealthRecordGrpcServer) GetSummary(ctx context.Context, req *common.IdRequest) (*health_record.HealthRecordSummary, error) {
	summary, err := s.svc.GetSummary(req.Id)
	if err != nil {
		return nil, err
	}
	timeline, _ := s.svc.GetTimeline(req.Id)
	events := make([]*health_record.TimelineEvent, len(timeline))
	for i, e := range timeline {
		events[i] = &health_record.TimelineEvent{
			Date:        e.Date,
			EventType:   e.EventType,
			Description: e.Description,
			RelatedId:   e.RelatedID,
		}
	}
	return &health_record.HealthRecordSummary{
		PatientId:          summary.PatientID,
		PatientName:        summary.PatientName,
		TotalVisits:        int32(summary.TotalVisits),
		TotalPrescriptions: int32(summary.TotalPrescriptions),
		TotalExaminations:  int32(summary.TotalExaminations),
		Timeline:           events,
	}, nil
}

// GetTimeline 获取健康档案时间轴
func (s *HealthRecordGrpcServer) GetTimeline(ctx context.Context, req *common.IdRequest) (*health_record.TimelineResponse, error) {
	events, err := s.svc.GetTimeline(req.Id)
	if err != nil {
		return nil, err
	}
	pbEvents := make([]*health_record.TimelineEvent, len(events))
	for i, e := range events {
		pbEvents[i] = &health_record.TimelineEvent{
			Date:        e.Date,
			EventType:   e.EventType,
			Description: e.Description,
			RelatedId:   e.RelatedID,
		}
	}
	return &health_record.TimelineResponse{
		Base:   &common.BaseResponse{Code: 0, Message: "查询成功"},
		Events: pbEvents,
	}, nil
}

// GrantAuthorization 授权档案访问
func (s *HealthRecordGrpcServer) GrantAuthorization(ctx context.Context, req *health_record.RecordAuthorization) (*health_record.RecordAuthorization, error) {
	auth := protoToAuthorization(req)
	if err := s.svc.GrantAuthorization(auth); err != nil {
		return nil, err
	}
	return authorizationToProto(auth), nil
}

// RevokeAuthorization 撤销档案授权
// 注意: proto RevokeAuthorization RPC 当前使用 common.IdRequest（仅含 id），
// 但 service 层需要 patient_id + doctor_id，待 proto 增加 RevokeAuthorizationRequest 消息后补全
func (s *HealthRecordGrpcServer) RevokeAuthorization(ctx context.Context, req *common.IdRequest) (*common.BaseResponse, error) {
	if err := s.svc.RevokeAuthorization("", req.Id); err != nil {
		return nil, err
	}
	return &common.BaseResponse{Code: 0, Message: "撤销成功"}, nil
}

// ---- 转换辅助函数 ----

func authorizationToProto(a *hrmodel.RecordAuthorization) *health_record.RecordAuthorization {
	pb := &health_record.RecordAuthorization{
		Id:        a.ID,
		PatientId: a.PatientID,
		DoctorId:  a.DoctorID,
		Status:    int32(a.Status),
	}
	if a.AuthTime != "" {
		pb.AuthTime = a.AuthTime
	}
	if a.ExpireTime != "" {
		pb.ExpireTime = a.ExpireTime
	}
	return pb
}

func protoToAuthorization(pb *health_record.RecordAuthorization) *hrmodel.RecordAuthorization {
	return &hrmodel.RecordAuthorization{
		ID:        pb.Id,
		PatientID: pb.PatientId,
		DoctorID:  pb.DoctorId,
		Status:    int8(pb.Status),
	}
}
