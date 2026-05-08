package inpatient

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/inpatient"
	inpatmodel "his-go/internal/inpatient/model"
	inpatsvc "his-go/internal/inpatient/service"
)

// InpatientGrpcServer gRPC 住院管理服务实现
type InpatientGrpcServer struct {
	inpatient.UnimplementedInpatientServiceServer
	svc *inpatsvc.InpatientService
}

// NewInpatientGrpcServer 创建 gRPC 住院管理服务
func NewInpatientGrpcServer(svc *inpatsvc.InpatientService) *InpatientGrpcServer {
	return &InpatientGrpcServer{svc: svc}
}

// AdmitPatient 入院登记
func (s *InpatientGrpcServer) AdmitPatient(ctx context.Context, req *inpatient.AdmitRequest) (*inpatient.InpatientRecord, error) {
	record := &inpatmodel.InpatientRecord{
		PatientID: req.PatientId,
		DeptID:    req.DeptId,
		RoomNo:    req.RoomNo,
		BedNo:     req.BedNo,
		Diagnosis: req.Diagnosis,
		Deposit:   req.Deposit,
		Status:    1,
	}
	if err := s.svc.AdmitPatient(record); err != nil {
		return nil, err
	}
	return inpatientRecordToProto(record), nil
}

// DischargePatient 出院
func (s *InpatientGrpcServer) DischargePatient(ctx context.Context, req *common.IdRequest) (*inpatient.InpatientRecord, error) {
	if err := s.svc.DischargePatient(req.Id); err != nil {
		return nil, err
	}
	record, err := s.svc.GetInpatient(req.Id)
	if err != nil {
		return nil, err
	}
	return inpatientRecordToProto(record), nil
}

// GetInpatient 获取住院记录
func (s *InpatientGrpcServer) GetInpatient(ctx context.Context, req *common.IdRequest) (*inpatient.InpatientRecord, error) {
	record, err := s.svc.GetInpatient(req.Id)
	if err != nil {
		return nil, err
	}
	return inpatientRecordToProto(record), nil
}

// CreateMedicalOrder 创建医嘱
func (s *InpatientGrpcServer) CreateMedicalOrder(ctx context.Context, req *inpatient.MedicalOrder) (*inpatient.MedicalOrder, error) {
	order := protoToMedicalOrder(req)
	if err := s.svc.CreateMedicalOrder(order); err != nil {
		return nil, err
	}
	return medicalOrderToProto(order), nil
}

// CreateNursingRecord 创建护理记录
func (s *InpatientGrpcServer) CreateNursingRecord(ctx context.Context, req *inpatient.NursingRecord) (*inpatient.NursingRecord, error) {
	record := protoToNursingRecord(req)
	if err := s.svc.CreateNursingRecord(record); err != nil {
		return nil, err
	}
	return nursingRecordToProto(record), nil
}

// ListInpatients 分页查询住院列表
func (s *InpatientGrpcServer) ListInpatients(ctx context.Context, req *inpatient.InpatientListRequest) (*inpatient.InpatientListResponse, error) {
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
	records, total, err := s.svc.ListInpatients(req.DeptId, int(req.Status), page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*inpatient.InpatientRecord, len(records))
	for i, r := range records {
		pbList[i] = inpatientRecordToProto(&r)
	}
	return &inpatient.InpatientListResponse{
		Base:    &common.BaseResponse{Code: 0, Message: "查询成功"},
		Records: pbList,
		Page:    &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// ---- 转换辅助函数 ----

func inpatientRecordToProto(r *inpatmodel.InpatientRecord) *inpatient.InpatientRecord {
	pb := &inpatient.InpatientRecord{
		Id:            r.ID,
		PatientId:     r.PatientID,
		PatientName:   r.PatientName,
		AdmissionDate: r.AdmissionDate.Format("2006-01-02"),
		DeptId:        r.DeptID,
		RoomNo:        r.RoomNo,
		BedNo:         r.BedNo,
		Diagnosis:     r.Diagnosis,
		Deposit:       r.Deposit,
		TotalCost:     r.TotalCost,
		Status:        int32(r.Status),
	}
	if !r.DischargeDate.IsZero() {
		pb.DischargeDate = r.DischargeDate.Format("2006-01-02")
	}
	pb.Orders = nil // MedicalOrders 需要通过独立查询加载
	return pb
}

func medicalOrderToProto(o *inpatmodel.MedicalOrder) *inpatient.MedicalOrder {
	pb := &inpatient.MedicalOrder{
		Id:          o.ID,
		InpatientId: o.InpatientID,
		DoctorId:    o.DoctorID,
		OrderType:   int32(o.OrderType),
		Content:     o.Content,
		Status:      int32(o.Status),
	}
	if !o.StartTime.IsZero() {
		pb.StartTime = o.StartTime.Format("2006-01-02 15:04:05")
	}
	if !o.EndTime.IsZero() {
		pb.EndTime = o.EndTime.Format("2006-01-02 15:04:05")
	}
	return pb
}

func protoToMedicalOrder(pb *inpatient.MedicalOrder) *inpatmodel.MedicalOrder {
	return &inpatmodel.MedicalOrder{
		ID:          pb.Id,
		InpatientID: pb.InpatientId,
		DoctorID:    pb.DoctorId,
		OrderType:   int8(pb.OrderType),
		Content:     pb.Content,
		Status:      int8(pb.Status),
	}
}

func nursingRecordToProto(r *inpatmodel.NursingRecord) *inpatient.NursingRecord {
	pb := &inpatient.NursingRecord{
		Id:          r.ID,
		InpatientId: r.InpatientID,
		NurseId:     r.NurseID,
		Content:     r.Content,
		VitalSigns:  r.VitalSigns,
	}
	if !r.RecordTime.IsZero() {
		pb.RecordTime = r.RecordTime.Format("2006-01-02 15:04:05")
	}
	return pb
}

func protoToNursingRecord(pb *inpatient.NursingRecord) *inpatmodel.NursingRecord {
	return &inpatmodel.NursingRecord{
		ID:          pb.Id,
		InpatientID: pb.InpatientId,
		NurseID:     pb.NurseId,
		Content:     pb.Content,
		VitalSigns:  pb.VitalSigns,
	}
}
