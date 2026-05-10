package emr

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/emr"
	emrmodel "his-go/internal/emr/model"
	emrsvc "his-go/internal/emr/service"
)

// EMRGrpcServer gRPC 电子病历服务实现
type EMRGrpcServer struct {
	emr.UnimplementedEMRServiceServer
	svc *emrsvc.EMRService
}

// NewEMRGrpcServer 创建 gRPC 电子病历服务
func NewEMRGrpcServer(svc *emrsvc.EMRService) *EMRGrpcServer {
	return &EMRGrpcServer{svc: svc}
}

// CreateRecord 创建病历记录
func (s *EMRGrpcServer) CreateRecord(ctx context.Context, req *emr.MedicalRecord) (*emr.MedicalRecord, error) {
	record := protoToMedicalRecord(req)
	if err := s.svc.CreateRecord(record); err != nil {
		return nil, err
	}
	return medicalRecordToProto(record), nil
}

// GetRecord 获取病历记录
func (s *EMRGrpcServer) GetRecord(ctx context.Context, req *common.IdRequest) (*emr.MedicalRecord, error) {
	record, err := s.svc.GetRecord(req.Id)
	if err != nil {
		return nil, err
	}
	return medicalRecordToProto(record), nil
}

// ListRecords 分页查询病历记录列表
func (s *EMRGrpcServer) ListRecords(ctx context.Context, req *emr.RecordListRequest) (*emr.RecordListResponse, error) {
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
	records, total, err := s.svc.ListRecords(req.PatientId, page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*emr.MedicalRecord, len(records))
	for i, r := range records {
		pbList[i] = medicalRecordToProto(&r)
	}
	return &emr.RecordListResponse{
		Base:    &common.BaseResponse{Code: 0, Message: "查询成功"},
		Records: pbList,
		Page:    &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// QualityControl 病历质控
func (s *EMRGrpcServer) QualityControl(ctx context.Context, req *emr.QualityControlRequest) (*emr.MedicalRecord, error) {
	if err := s.svc.QualityControl(req.RecordId, req.ReviewerId, int(req.Level), req.Comment); err != nil {
		return nil, err
	}
	record, err := s.svc.GetRecord(req.RecordId)
	if err != nil {
		return nil, err
	}
	return medicalRecordToProto(record), nil
}

// ListTemplates 查询病历模板列表
func (s *EMRGrpcServer) ListTemplates(ctx context.Context, req *common.Empty) (*emr.TemplateListResponse, error) {
	_ = req
	templates, err := s.svc.ListTemplates()
	if err != nil {
		return nil, err
	}
	pbList := make([]*emr.RecordTemplate, len(templates))
	for i, t := range templates {
		pbList[i] = &emr.RecordTemplate{
			Id:      t.ID,
			Name:    t.Name,
			DeptId:  t.DeptID,
			Content: t.Content,
			Type:    int32(t.Type),
		}
	}
	return &emr.TemplateListResponse{
		Base:      &common.BaseResponse{Code: 0, Message: "查询成功"},
		Templates: pbList,
	}, nil
}

// CDSSCheck CDSS 合理性校验
func (s *EMRGrpcServer) CDSSCheck(ctx context.Context, req *emr.CDSSCheckRequest) (*emr.CDSSCheckResult, error) {
	warnings, err := s.svc.CDSSCheck(req.PatientId, req.DrugId, req.Diagnosis)
	if err != nil {
		return nil, err
	}
	safe := len(warnings) == 0
	return &emr.CDSSCheckResult{
		Safe:        safe,
		Warnings:    warnings,
		Suggestions: warnings,
	}, nil
}

// ---- 转换辅助函数 ----

func medicalRecordToProto(r *emrmodel.MedicalRecord) *emr.MedicalRecord {
	return &emr.MedicalRecord{
		Id:             r.ID,
		PatientId:      r.PatientID,
		ClinicRecordId: r.ClinicRecordID,
		TemplateId:     r.TemplateID,
		ChiefComplaint: r.ChiefComplaint,
		PresentIllness: r.PresentIllness,
		PastHistory:    r.PastHistory,
		PhysicalExam:   r.PhysicalExam,
		AuxiliaryExam:  r.AuxiliaryExam,
		Diagnosis:      r.Diagnosis,
		TreatmentPlan:  r.TreatmentPlan,
		QualityLevel:   int32(r.QualityLevel),
		Status:         int32(r.Status),
		CreatedAt:      r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func protoToMedicalRecord(pb *emr.MedicalRecord) *emrmodel.MedicalRecord {
	return &emrmodel.MedicalRecord{
		ID:             pb.Id,
		PatientID:      pb.PatientId,
		ClinicRecordID: pb.ClinicRecordId,
		TemplateID:     pb.TemplateId,
		ChiefComplaint: pb.ChiefComplaint,
		PresentIllness: pb.PresentIllness,
		PastHistory:    pb.PastHistory,
		PhysicalExam:   pb.PhysicalExam,
		AuxiliaryExam:  pb.AuxiliaryExam,
		Diagnosis:      pb.Diagnosis,
		TreatmentPlan:  pb.TreatmentPlan,
		QualityLevel:   int(pb.QualityLevel),
		Status:         int8(pb.Status),
	}
}
