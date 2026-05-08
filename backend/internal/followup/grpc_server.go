package followup

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/followup"
	fupmodel "his-go/internal/followup/model"
	fupsvc "his-go/internal/followup/service"
)

// FollowupGrpcServer gRPC 随访管理服务实现
type FollowupGrpcServer struct {
	followup.UnimplementedFollowupServiceServer
	svc *fupsvc.FollowupService
}

// NewFollowupGrpcServer 创建 gRPC 随访管理服务
func NewFollowupGrpcServer(svc *fupsvc.FollowupService) *FollowupGrpcServer {
	return &FollowupGrpcServer{svc: svc}
}

// CreatePlan 创建随访计划
func (s *FollowupGrpcServer) CreatePlan(ctx context.Context, req *followup.FollowupPlan) (*followup.FollowupPlan, error) {
	plan := protoToFollowupPlan(req)
	if err := s.svc.CreatePlan(plan); err != nil {
		return nil, err
	}
	return followupPlanToProto(plan), nil
}

// GetPlan 获取随访计划
func (s *FollowupGrpcServer) GetPlan(ctx context.Context, req *common.IdRequest) (*followup.FollowupPlan, error) {
	plan, err := s.svc.GetPlan(req.Id)
	if err != nil {
		return nil, err
	}
	return followupPlanToProto(plan), nil
}

// ExecuteTask 执行随访任务
func (s *FollowupGrpcServer) ExecuteTask(ctx context.Context, req *followup.ExecuteTaskRequest) (*followup.FollowupTask, error) {
	if err := s.svc.ExecuteTask(req.TaskId, req.Result); err != nil {
		return nil, err
	}
	return &followup.FollowupTask{
		Id:     req.TaskId,
		Status: 2,
	}, nil
}

// SubmitSurvey 提交满意度调查
func (s *FollowupGrpcServer) SubmitSurvey(ctx context.Context, req *followup.SatisfactionSurvey) (*followup.SatisfactionSurvey, error) {
	survey := protoToSurvey(req)
	if err := s.svc.SubmitSurvey(survey); err != nil {
		return nil, err
	}
	return surveyToProto(survey), nil
}

// ListPlans 分页查询随访计划列表
func (s *FollowupGrpcServer) ListPlans(ctx context.Context, req *followup.PlanListRequest) (*followup.PlanListResponse, error) {
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
	plans, total, err := s.svc.ListPlans(req.PatientId, int(req.Status), page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*followup.FollowupPlan, len(plans))
	for i, p := range plans {
		pbList[i] = followupPlanToProto(&p)
	}
	return &followup.PlanListResponse{
		Base:  &common.BaseResponse{Code: 0, Message: "查询成功"},
		Plans: pbList,
		Page:  &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// ---- 转换辅助函数 ----

func followupPlanToProto(p *fupmodel.FollowupPlan) *followup.FollowupPlan {
	pb := &followup.FollowupPlan{
		Id:         p.ID,
		PatientId:  p.PatientID,
		TemplateId: p.TemplateID,
		PlanName:   p.PlanName,
		StartDate:  p.StartDate,
		EndDate:    p.EndDate,
		Frequency:  int32(p.Frequency),
		Status:     int32(p.Status),
	}
	return pb
}

func protoToFollowupPlan(pb *followup.FollowupPlan) *fupmodel.FollowupPlan {
	plan := &fupmodel.FollowupPlan{
		ID:         pb.Id,
		PatientID:  pb.PatientId,
		TemplateID: pb.TemplateId,
		PlanName:   pb.PlanName,
		Frequency:  int(pb.Frequency),
		Status:     int8(pb.Status),
	}
	return plan
}

func surveyToProto(s *fupmodel.SatisfactionSurvey) *followup.SatisfactionSurvey {
	pb := &followup.SatisfactionSurvey{
		Id:             s.ID,
		FollowupTaskId: s.FollowupTaskID,
		PatientId:      s.PatientID,
		Score:          int32(s.Score),
		Feedback:       s.Feedback,
	}
	if !s.CreatedAt.IsZero() {
		pb.CreatedAt = s.CreatedAt.Format("2006-01-02 15:04:05")
	}
	return pb
}

func protoToSurvey(pb *followup.SatisfactionSurvey) *fupmodel.SatisfactionSurvey {
	return &fupmodel.SatisfactionSurvey{
		ID:             pb.Id,
		FollowupTaskID: pb.FollowupTaskId,
		PatientID:      pb.PatientId,
		Score:          int(pb.Score),
		Feedback:       pb.Feedback,
	}
}
