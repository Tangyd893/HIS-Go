package prescription

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/prescription"
	psmodel "his-go/internal/prescription/model"
	pssvc "his-go/internal/prescription/service"
)

// PrescriptionGrpcServer gRPC 处方服务实现
type PrescriptionGrpcServer struct {
	prescription.UnimplementedPrescriptionServiceServer
	svc *pssvc.PrescriptionService
}

// NewPrescriptionGrpcServer 创建 gRPC 处方服务
func NewPrescriptionGrpcServer(svc *pssvc.PrescriptionService) *PrescriptionGrpcServer {
	return &PrescriptionGrpcServer{svc: svc}
}

// CreatePrescription 创建处方
func (s *PrescriptionGrpcServer) CreatePrescription(ctx context.Context, req *prescription.PrescriptionInfo) (*prescription.PrescriptionInfo, error) {
	ps, details := protoToPrescription(req)
	if err := s.svc.Create(ps, details); err != nil {
		return nil, err
	}
	return prescriptionToProto(ps, details), nil
}

// GetPrescription 获取处方详情
func (s *PrescriptionGrpcServer) GetPrescription(ctx context.Context, req *common.IdRequest) (*prescription.PrescriptionInfo, error) {
	ps, err := s.svc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return prescriptionToProto(ps, ps.Details), nil
}

// ListPrescriptions 分页查询处方列表
func (s *PrescriptionGrpcServer) ListPrescriptions(ctx context.Context, req *prescription.PrescriptionListRequest) (*prescription.PrescriptionListResponse, error) {
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

	var prescriptions []psmodel.Prescription
	var total int64
	var err error

	if req.DoctorId != "" {
		prescriptions, total, err = s.svc.ListByDoctor(req.DoctorId, page, pageSize)
	} else {
		prescriptions, total, err = s.svc.ListByPatient(req.PatientId, page, pageSize)
	}
	if err != nil {
		return nil, err
	}

	pbList := make([]*prescription.PrescriptionInfo, len(prescriptions))
	for i, p := range prescriptions {
		pbList[i] = prescriptionToProto(&p, p.Details)
	}
	return &prescription.PrescriptionListResponse{
		Base:          &common.BaseResponse{Code: 0, Message: "查询成功"},
		Prescriptions: pbList,
		Page:          &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// ReviewPrescription 审核处方
func (s *PrescriptionGrpcServer) ReviewPrescription(ctx context.Context, req *prescription.ReviewRequest) (*prescription.PrescriptionInfo, error) {
	if err := s.svc.Review(req.PrescriptionId, req.Approved, req.Comment); err != nil {
		return nil, err
	}
	ps, err := s.svc.GetByID(req.PrescriptionId)
	if err != nil {
		return nil, err
	}
	return prescriptionToProto(ps, ps.Details), nil
}

// CancelPrescription 取消处方
func (s *PrescriptionGrpcServer) CancelPrescription(ctx context.Context, req *common.IdRequest) (*common.BaseResponse, error) {
	if err := s.svc.Cancel(req.Id); err != nil {
		return nil, err
	}
	return &common.BaseResponse{Code: 0, Message: "取消成功"}, nil
}

// ---- 转换辅助函数 ----

func prescriptionToProto(ps *psmodel.Prescription, details []psmodel.PrescriptionDetail) *prescription.PrescriptionInfo {
	pb := &prescription.PrescriptionInfo{
		Id:               ps.ID,
		PatientId:        ps.PatientID,
		PatientName:      ps.PatientName,
		DoctorId:         ps.DoctorID,
		DiagnosisId:      ps.DiagnosisID,
		PrescriptionType: int32(ps.PrescriptionType),
		Status:           int32(ps.Status),
		Note:             ps.Note,
		CreatedAt:        ps.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	pbDetails := make([]*prescription.PrescriptionDetail, len(details))
	for i, d := range details {
		pbDetails[i] = &prescription.PrescriptionDetail{
			Id:             d.ID,
			PrescriptionId: d.PrescriptionID,
			DrugId:         d.DrugID,
			DrugName:       d.DrugName,
			Specification:  d.Specification,
			Dosage:         d.Dosage,
			Usage:          d.Usage,
			Frequency:      d.Frequency,
			Days:           int32(d.Days),
			Quantity:       int32(d.Quantity),
			UnitPrice:      d.UnitPrice,
			Note:           d.Note,
		}
	}
	pb.Details = pbDetails
	return pb
}

func protoToPrescription(pb *prescription.PrescriptionInfo) (*psmodel.Prescription, []psmodel.PrescriptionDetail) {
	ps := &psmodel.Prescription{
		ID:               pb.Id,
		PatientID:        pb.PatientId,
		PatientName:      pb.PatientName,
		DoctorID:         pb.DoctorId,
		DiagnosisID:      pb.DiagnosisId,
		PrescriptionType: int8(pb.PrescriptionType),
		Status:           int8(pb.Status),
		Note:             pb.Note,
	}
	details := make([]psmodel.PrescriptionDetail, len(pb.Details))
	for i, d := range pb.Details {
		details[i] = psmodel.PrescriptionDetail{
			ID:             d.Id,
			PrescriptionID: d.PrescriptionId,
			DrugID:         d.DrugId,
			DrugName:       d.DrugName,
			Specification:  d.Specification,
			Dosage:         d.Dosage,
			Usage:          d.Usage,
			Frequency:      d.Frequency,
			Days:           int(d.Days),
			Quantity:       int(d.Quantity),
			UnitPrice:      d.UnitPrice,
			Note:           d.Note,
		}
	}
	return ps, details
}
