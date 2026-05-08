package clinic

import (
	"context"

	"his-go/api/proto/clinic"
	"his-go/api/proto/common"
	clinicmodel "his-go/internal/clinic/model"
	clinicsvc "his-go/internal/clinic/service"
)

// ClinicGrpcServer gRPC 门诊诊疗服务实现
type ClinicGrpcServer struct {
	clinic.UnimplementedClinicServiceServer
	svc *clinicsvc.ClinicService
}

// NewClinicGrpcServer 创建 gRPC 门诊诊疗服务
func NewClinicGrpcServer(svc *clinicsvc.ClinicService) *ClinicGrpcServer {
	return &ClinicGrpcServer{svc: svc}
}

// CreateClinicRecord 创建接诊记录
func (s *ClinicGrpcServer) CreateClinicRecord(ctx context.Context, req *clinic.ClinicRecord) (*clinic.ClinicRecord, error) {
	record := protoToClinicRecord(req)
	if err := s.svc.CreateRecord(record); err != nil {
		return nil, err
	}
	return clinicRecordToProto(record), nil
}

// GetClinicRecord 获取接诊记录
func (s *ClinicGrpcServer) GetClinicRecord(ctx context.Context, req *common.IdRequest) (*clinic.ClinicRecord, error) {
	record, err := s.svc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return clinicRecordToProto(record), nil
}

// ListClinicRecords 分页查询接诊记录
func (s *ClinicGrpcServer) ListClinicRecords(ctx context.Context, req *clinic.ClinicListRequest) (*clinic.ClinicListResponse, error) {
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

	var records []clinicmodel.ClinicRecord
	var total int64
	var err error

	if req.DoctorId != "" {
		records, total, err = s.svc.ListByDoctor(req.DoctorId, page, pageSize)
	} else {
		records, total, err = s.svc.ListByPatient(req.PatientId, page, pageSize)
	}
	if err != nil {
		return nil, err
	}

	pbList := make([]*clinic.ClinicRecord, len(records))
	for i, r := range records {
		pbList[i] = clinicRecordToProto(&r)
	}
	return &clinic.ClinicListResponse{
		Base:    &common.BaseResponse{Code: 0, Message: "查询成功"},
		Records: pbList,
		Page:    &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// CreateExaminationRequest 创建检查申请
func (s *ClinicGrpcServer) CreateExaminationRequest(ctx context.Context, req *clinic.ExaminationRequest) (*clinic.ExaminationRequest, error) {
	examReq := protoToExaminationRequest(req)
	if err := s.svc.CreateExamRequest(examReq); err != nil {
		return nil, err
	}
	return examRequestToProto(examReq), nil
}

// ---- 转换辅助函数 ----

func clinicRecordToProto(r *clinicmodel.ClinicRecord) *clinic.ClinicRecord {
	return &clinic.ClinicRecord{
		Id:             r.ID,
		RegistrationId: r.RegistrationID,
		PatientId:      r.PatientID,
		PatientName:    r.PatientName,
		DoctorId:       r.DoctorID,
		ChiefComplaint: r.ChiefComplaint,
		PresentIllness: r.PresentIllness,
		Diagnosis:      r.Diagnosis,
		IcdCode:        r.IcdCode,
		Advice:         r.Advice,
		Status:         int32(r.Status),
		VisitTime:      r.VisitTime.Format("2006-01-02 15:04:05"),
		CreatedAt:      r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func protoToClinicRecord(pb *clinic.ClinicRecord) *clinicmodel.ClinicRecord {
	return &clinicmodel.ClinicRecord{
		ID:             pb.Id,
		RegistrationID: pb.RegistrationId,
		PatientID:      pb.PatientId,
		PatientName:    pb.PatientName,
		DoctorID:       pb.DoctorId,
		ChiefComplaint: pb.ChiefComplaint,
		PresentIllness: pb.PresentIllness,
		Diagnosis:      pb.Diagnosis,
		IcdCode:        pb.IcdCode,
		Advice:         pb.Advice,
		Status:         int8(pb.Status),
	}
}

func examRequestToProto(r *clinicmodel.ExaminationRequest) *clinic.ExaminationRequest {
	return &clinic.ExaminationRequest{
		Id:                r.ID,
		ClinicRecordId:    r.ClinicRecordID,
		PatientId:         r.PatientID,
		ExamType:          r.ExamType,
		ExamItem:          r.ExamItem,
		BodyPart:          r.BodyPart,
		ClinicalDiagnosis: r.ClinicalDiagnosis,
		Note:              r.Note,
		Status:            int32(r.Status),
	}
}

func protoToExaminationRequest(pb *clinic.ExaminationRequest) *clinicmodel.ExaminationRequest {
	return &clinicmodel.ExaminationRequest{
		ID:                pb.Id,
		ClinicRecordID:    pb.ClinicRecordId,
		PatientID:         pb.PatientId,
		ExamType:          pb.ExamType,
		ExamItem:          pb.ExamItem,
		BodyPart:          pb.BodyPart,
		ClinicalDiagnosis: pb.ClinicalDiagnosis,
		Note:              pb.Note,
		Status:            int8(pb.Status),
	}
}
