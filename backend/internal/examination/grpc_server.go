package examination

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/examination"
	exammodel "his-go/internal/examination/model"
	examsvc "his-go/internal/examination/service"
)

// ExaminationGrpcServer gRPC 检查检验服务实现
type ExaminationGrpcServer struct {
	examination.UnimplementedExaminationServiceServer
	svc *examsvc.ExaminationService
}

// NewExaminationGrpcServer 创建 gRPC 检查检验服务
func NewExaminationGrpcServer(svc *examsvc.ExaminationService) *ExaminationGrpcServer {
	return &ExaminationGrpcServer{svc: svc}
}

// CreateReport 创建检查报告
func (s *ExaminationGrpcServer) CreateReport(ctx context.Context, req *examination.ExaminationReport) (*examination.ExaminationReport, error) {
	report := protoToReport(req)
	if err := s.svc.CreateReport(report); err != nil {
		return nil, err
	}
	return reportToProto(report), nil
}

// GetReport 获取检查报告
func (s *ExaminationGrpcServer) GetReport(ctx context.Context, req *common.IdRequest) (*examination.ExaminationReport, error) {
	report, err := s.svc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return reportToProto(report), nil
}

// ReviewReport 审核检查报告
func (s *ExaminationGrpcServer) ReviewReport(ctx context.Context, req *examination.ReviewReportRequest) (*examination.ExaminationReport, error) {
	if err := s.svc.Review(req.ReportId, req.ReviewerId, req.Approved, req.Comment); err != nil {
		return nil, err
	}
	report, err := s.svc.GetByID(req.ReportId)
	if err != nil {
		return nil, err
	}
	return reportToProto(report), nil
}

// ListReports 分页查询检查报告列表
func (s *ExaminationGrpcServer) ListReports(ctx context.Context, req *examination.ReportListRequest) (*examination.ReportListResponse, error) {
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
	reports, total, err := s.svc.ListByPatient(req.PatientId, int(req.Status), page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*examination.ExaminationReport, len(reports))
	for i, r := range reports {
		pbList[i] = reportToProto(&r)
	}
	return &examination.ReportListResponse{
		Base:    &common.BaseResponse{Code: 0, Message: "查询成功"},
		Reports: pbList,
		Page:    &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// ---- 转换辅助函数 ----

func reportToProto(r *exammodel.ExaminationReport) *examination.ExaminationReport {
	return &examination.ExaminationReport{
		Id:                   r.ID,
		PatientId:            r.PatientID,
		PatientName:          r.PatientName,
		ExaminationRequestId: r.ExamRequestID,
		ExamType:             r.ExamType,
		ExamItem:             r.ExamItem,
		BodyPart:             r.BodyPart,
		Findings:             r.Findings,
		Impression:           r.Impression,
		Conclusion:           r.Conclusion,
		TechnicianId:         r.TechnicianID,
		ReviewerId:           r.ReviewerID,
		Status:               int32(r.Status),
		CreatedAt:            r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func protoToReport(pb *examination.ExaminationReport) *exammodel.ExaminationReport {
	return &exammodel.ExaminationReport{
		ID:            pb.Id,
		PatientID:     pb.PatientId,
		PatientName:   pb.PatientName,
		ExamRequestID: pb.ExaminationRequestId,
		ExamType:      pb.ExamType,
		ExamItem:      pb.ExamItem,
		BodyPart:      pb.BodyPart,
		Findings:      pb.Findings,
		Impression:    pb.Impression,
		Conclusion:    pb.Conclusion,
		TechnicianID:  pb.TechnicianId,
		ReviewerID:    pb.ReviewerId,
		Status:        int8(pb.Status),
	}
}
