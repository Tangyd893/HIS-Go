package billing

import (
	"context"

	"his-go/api/proto/billing"
	"his-go/api/proto/common"
	billmodel "his-go/internal/billing/model"
	billsvc "his-go/internal/billing/service"
)

// BillingGrpcServer gRPC 收费结算服务实现
type BillingGrpcServer struct {
	billing.UnimplementedBillingServiceServer
	svc *billsvc.BillingService
}

// NewBillingGrpcServer 创建 gRPC 收费结算服务
func NewBillingGrpcServer(svc *billsvc.BillingService) *BillingGrpcServer {
	return &BillingGrpcServer{svc: svc}
}

// CreateBill 创建账单
func (s *BillingGrpcServer) CreateBill(ctx context.Context, req *billing.BillInfo) (*billing.BillInfo, error) {
	bill, details := protoToBill(req)
	if err := s.svc.CreateBill(bill, details); err != nil {
		return nil, err
	}
	return billToProto(bill, details), nil
}

// GetBill 获取账单详情
func (s *BillingGrpcServer) GetBill(ctx context.Context, req *common.IdRequest) (*billing.BillInfo, error) {
	bill, err := s.svc.GetBill(req.Id)
	if err != nil {
		return nil, err
	}
	details, _ := s.svc.GetBillDetails(req.Id)
	return billToProto(bill, details), nil
}

// PayBill 支付账单
func (s *BillingGrpcServer) PayBill(ctx context.Context, req *billing.PayRequest) (*billing.BillInfo, error) {
	if err := s.svc.Pay(req.BillId, int8(req.PayMethod)); err != nil {
		return nil, err
	}
	bill, err := s.svc.GetBill(req.BillId)
	if err != nil {
		return nil, err
	}
	details, _ := s.svc.GetBillDetails(req.BillId)
	return billToProto(bill, details), nil
}

// RefundBill 退款
func (s *BillingGrpcServer) RefundBill(ctx context.Context, req *common.IdRequest) (*billing.BillInfo, error) {
	if err := s.svc.Refund(req.Id); err != nil {
		return nil, err
	}
	bill, err := s.svc.GetBill(req.Id)
	if err != nil {
		return nil, err
	}
	details, _ := s.svc.GetBillDetails(req.Id)
	return billToProto(bill, details), nil
}

// ListBills 分页查询账单列表
func (s *BillingGrpcServer) ListBills(ctx context.Context, req *billing.BillListRequest) (*billing.BillListResponse, error) {
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
	bills, total, err := s.svc.ListBills(req.PatientId, int(req.Status), page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*billing.BillInfo, len(bills))
	for i, b := range bills {
		details, _ := s.svc.GetBillDetails(b.ID)
		pbList[i] = billToProto(&b, details)
	}
	return &billing.BillListResponse{
		Base:  &common.BaseResponse{Code: 0, Message: "查询成功"},
		Bills: pbList,
		Page:  &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// ---- 转换辅助函数 ----

func billToProto(bill *billmodel.Bill, details []billmodel.BillDetail) *billing.BillInfo {
	pb := &billing.BillInfo{
		Id:             bill.ID,
		PatientId:      bill.PatientID,
		RegistrationId: bill.RegistrationID,
		BillNo:         bill.BillNo,
		TotalAmount:    bill.TotalAmount,
		PaidAmount:     bill.PaidAmount,
		PayMethod:      int32(bill.PayMethod),
		Status:         int32(bill.Status),
		CreatedAt:      bill.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	pbDetails := make([]*billing.BillDetail, len(details))
	for i, d := range details {
		pbDetails[i] = &billing.BillDetail{
			Id:        d.ID,
			BillId:    d.BillID,
			ItemType:  int32(d.ItemType),
			ItemName:  d.ItemName,
			UnitPrice: d.UnitPrice,
			Quantity:  int32(d.Quantity),
			Amount:    d.Amount,
		}
	}
	pb.Details = pbDetails
	return pb
}

func protoToBill(pb *billing.BillInfo) (*billmodel.Bill, []billmodel.BillDetail) {
	bill := &billmodel.Bill{
		ID:             pb.Id,
		PatientID:      pb.PatientId,
		RegistrationID: pb.RegistrationId,
		BillNo:         pb.BillNo,
		TotalAmount:    pb.TotalAmount,
		PaidAmount:     pb.PaidAmount,
		PayMethod:      int8(pb.PayMethod),
		Status:         int8(pb.Status),
	}
	details := make([]billmodel.BillDetail, len(pb.Details))
	for i, d := range pb.Details {
		details[i] = billmodel.BillDetail{
			ID:        d.Id,
			BillID:    d.BillId,
			ItemType:  int8(d.ItemType),
			ItemName:  d.ItemName,
			UnitPrice: d.UnitPrice,
			Quantity:  int(d.Quantity),
			Amount:    d.Amount,
		}
	}
	return bill, details
}
