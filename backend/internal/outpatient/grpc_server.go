package outpatient

import (
	"context"
	"strings"
	"time"

	"his-go/api/proto/common"
	"his-go/api/proto/outpatient"
	outpatmodel "his-go/internal/outpatient/model"
	outpatsvc "his-go/internal/outpatient/service"
)

// OutpatientGrpcServer gRPC 院外患者服务实现
type OutpatientGrpcServer struct {
	outpatient.UnimplementedOutpatientServiceServer
	svc *outpatsvc.OutpatientService
}

// NewOutpatientGrpcServer 创建 gRPC 院外患者服务
func NewOutpatientGrpcServer(svc *outpatsvc.OutpatientService) *OutpatientGrpcServer {
	return &OutpatientGrpcServer{svc: svc}
}

// CreateConsultation 创建在线问诊
func (s *OutpatientGrpcServer) CreateConsultation(ctx context.Context, req *outpatient.Consultation) (*outpatient.Consultation, error) {
	c := protoToConsultation(req)
	if err := s.svc.CreateConsultation(c); err != nil {
		return nil, err
	}
	return consultationToProto(c), nil
}

// GetConsultation 获取在线问诊详情
func (s *OutpatientGrpcServer) GetConsultation(ctx context.Context, req *common.IdRequest) (*outpatient.Consultation, error) {
	c, err := s.svc.GetConsultation(req.Id)
	if err != nil {
		return nil, err
	}
	return consultationToProto(c), nil
}

// SendMessage 发送问诊消息
func (s *OutpatientGrpcServer) SendMessage(ctx context.Context, req *outpatient.ConsultationMessage) (*outpatient.ConsultationMessage, error) {
	msg := protoToConsultationMessage(req)
	if err := s.svc.SendMessage(msg); err != nil {
		return nil, err
	}
	return consultationMessageToProto(msg), nil
}

// CreateChronicContract 创建慢病签约
func (s *OutpatientGrpcServer) CreateChronicContract(ctx context.Context, req *outpatient.ChronicContract) (*outpatient.ChronicContract, error) {
	contract := protoToChronicContract(req)
	if err := s.svc.CreateChronicContract(contract); err != nil {
		return nil, err
	}
	return chronicContractToProto(contract), nil
}

// ReportHealthData 上报健康数据
func (s *OutpatientGrpcServer) ReportHealthData(ctx context.Context, req *outpatient.HealthData) (*outpatient.HealthData, error) {
	data := protoToHealthData(req)
	if err := s.svc.ReportHealthData(data); err != nil {
		return nil, err
	}
	return healthDataToProto(data), nil
}

// ListConsultations 分页查询在线问诊列表
func (s *OutpatientGrpcServer) ListConsultations(ctx context.Context, req *outpatient.ConsultationListRequest) (*outpatient.ConsultationListResponse, error) {
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
	consultations, total, err := s.svc.ListConsultations(req.PatientId, req.DoctorId, int(req.Status), page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*outpatient.Consultation, len(consultations))
	for i, c := range consultations {
		pbList[i] = consultationToProto(&c)
	}
	return &outpatient.ConsultationListResponse{
		Base:          &common.BaseResponse{Code: 0, Message: "查询成功"},
		Consultations: pbList,
		Page:          &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// ---- 转换辅助函数 ----

func consultationToProto(c *outpatmodel.Consultation) *outpatient.Consultation {
	pb := &outpatient.Consultation{
		Id:          c.ID,
		PatientId:   c.PatientID,
		DoctorId:    c.DoctorID,
		Type:        int32(c.Type),
		Description: c.Description,
		Images:      strings.Fields(c.Images),
		Status:      int32(c.Status),
		CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return pb
}

func protoToConsultation(pb *outpatient.Consultation) *outpatmodel.Consultation {
	now := time.Now()
	return &outpatmodel.Consultation{
		ID:          pb.Id,
		PatientID:   pb.PatientId,
		DoctorID:    pb.DoctorId,
		Type:        int8(pb.Type),
		Description: pb.Description,
		Images:      strings.Join(pb.Images, " "),
		Status:      int8(pb.Status),
		CreatedAt:   now,
	}
}

func consultationMessageToProto(m *outpatmodel.ConsultationMessage) *outpatient.ConsultationMessage {
	pb := &outpatient.ConsultationMessage{
		Id:             m.ID,
		ConsultationId: m.ConsultationID,
		SenderId:       m.SenderID,
		SenderName:     m.SenderName,
		Content:        m.Content,
		MsgType:        m.MsgType,
	}
	if !m.CreatedAt.IsZero() {
		pb.CreatedAt = m.CreatedAt.Format("2006-01-02 15:04:05")
	}
	return pb
}

func protoToConsultationMessage(pb *outpatient.ConsultationMessage) *outpatmodel.ConsultationMessage {
	return &outpatmodel.ConsultationMessage{
		ID:             pb.Id,
		ConsultationID: pb.ConsultationId,
		SenderID:       pb.SenderId,
		SenderName:     pb.SenderName,
		Content:        pb.Content,
		MsgType:        pb.MsgType,
		CreatedAt:      time.Now(),
	}
}

func chronicContractToProto(c *outpatmodel.ChronicContract) *outpatient.ChronicContract {
	pb := &outpatient.ChronicContract{
		Id:           c.ID,
		PatientId:    c.PatientID,
		DoctorId:     c.DoctorID,
		DiseaseType:  c.DiseaseType,
		ContractDate: c.ContractDate,
		EndDate:      c.EndDate,
		Status:       int32(c.Status),
	}
	return pb
}

func protoToChronicContract(pb *outpatient.ChronicContract) *outpatmodel.ChronicContract {
	return &outpatmodel.ChronicContract{
		ID:           pb.Id,
		PatientID:    pb.PatientId,
		DoctorID:     pb.DoctorId,
		DiseaseType:  pb.DiseaseType,
		ContractDate: pb.ContractDate,
		EndDate:      pb.EndDate,
		Status:       int8(pb.Status),
	}
}

func healthDataToProto(d *outpatmodel.HealthData) *outpatient.HealthData {
	pb := &outpatient.HealthData{
		Id:        d.ID,
		PatientId: d.PatientID,
		DataType:  d.DataType,
		Value:     d.Value,
		Unit:      d.Unit,
		Abnormal:  d.Abnormal,
	}
	if d.MeasureTime != "" {
		pb.MeasureTime = d.MeasureTime
	}
	return pb
}

func protoToHealthData(pb *outpatient.HealthData) *outpatmodel.HealthData {
	now := time.Now()
	return &outpatmodel.HealthData{
		ID:          pb.Id,
		PatientID:   pb.PatientId,
		DataType:    pb.DataType,
		Value:       pb.Value,
		Unit:        pb.Unit,
		MeasureTime: now.Format("2006-01-02 15:04:05"),
		Abnormal:    pb.Abnormal,
	}
}
