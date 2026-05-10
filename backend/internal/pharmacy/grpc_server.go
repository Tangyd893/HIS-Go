package pharmacy

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/pharmacy"
	pharmmodel "his-go/internal/pharmacy/model"
	pharmsvc "his-go/internal/pharmacy/service"
)

// PharmacyGrpcServer gRPC 药房管理服务实现
type PharmacyGrpcServer struct {
	pharmacy.UnimplementedPharmacyServiceServer
	svc *pharmsvc.PharmacyService
}

// NewPharmacyGrpcServer 创建 gRPC 药房管理服务
func NewPharmacyGrpcServer(svc *pharmsvc.PharmacyService) *PharmacyGrpcServer {
	return &PharmacyGrpcServer{svc: svc}
}

// ListDrugs 分页查询药品列表
func (s *PharmacyGrpcServer) ListDrugs(ctx context.Context, req *pharmacy.DrugListRequest) (*pharmacy.DrugListResponse, error) {
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
	drugs, total, err := s.svc.ListDrugs(req.Name, page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*pharmacy.DrugInfo, len(drugs))
	for i, d := range drugs {
		pbList[i] = drugToProto(&d)
	}
	return &pharmacy.DrugListResponse{
		Base:  &common.BaseResponse{Code: 0, Message: "查询成功"},
		Drugs: pbList,
		Page:  &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}

// GetDrug 获取药品详情
func (s *PharmacyGrpcServer) GetDrug(ctx context.Context, req *common.IdRequest) (*pharmacy.DrugInfo, error) {
	drug, err := s.svc.GetDrug(req.Id)
	if err != nil {
		return nil, err
	}
	return drugToProto(drug), nil
}

// AddStock 增加库存
func (s *PharmacyGrpcServer) AddStock(ctx context.Context, req *pharmacy.StockRequest) (*pharmacy.DrugInfo, error) {
	if err := s.svc.AddStock(req.DrugId, int(req.Quantity)); err != nil {
		return nil, err
	}
	drug, err := s.svc.GetDrug(req.DrugId)
	if err != nil {
		return nil, err
	}
	return drugToProto(drug), nil
}

// ReduceStock 减少库存
func (s *PharmacyGrpcServer) ReduceStock(ctx context.Context, req *pharmacy.StockRequest) (*pharmacy.DrugInfo, error) {
	if err := s.svc.AddStock(req.DrugId, -int(req.Quantity)); err != nil {
		return nil, err
	}
	drug, err := s.svc.GetDrug(req.DrugId)
	if err != nil {
		return nil, err
	}
	return drugToProto(drug), nil
}

// DispenseDrug 发药
func (s *PharmacyGrpcServer) DispenseDrug(ctx context.Context, req *pharmacy.DispenseRequest) (*pharmacy.DispenseRecord, error) {
	if err := s.svc.DispenseDrug(req.PrescriptionId, req.DrugId, int(req.Quantity), req.DispenserId); err != nil {
		return nil, err
	}
	return &pharmacy.DispenseRecord{
		PrescriptionId: req.PrescriptionId,
		DrugId:         req.DrugId,
		Quantity:       req.Quantity,
		DispenserId:    req.DispenserId,
		Status:         1,
	}, nil
}

// CheckExpiredDrugs 校验过期药品
func (s *PharmacyGrpcServer) CheckExpiredDrugs(ctx context.Context, req *common.Empty) (*pharmacy.ExpiredDrugResponse, error) {
	_ = req
	s.svc.CheckAndAlertExpired()
	return &pharmacy.ExpiredDrugResponse{
		Base: &common.BaseResponse{Code: 0, Message: "过期检查已触发"},
	}, nil
}

// ---- 转换辅助函数 ----

func drugToProto(d *pharmmodel.Drug) *pharmacy.DrugInfo {
	return &pharmacy.DrugInfo{
		Id:            d.ID,
		Name:          d.Name,
		GenericName:   d.GenericName,
		Specification: d.Specification,
		Manufacturer:  d.Manufacturer,
		BatchNo:       d.BatchNo,
		PurchasePrice: d.PurchasePrice,
		RetailPrice:   d.RetailPrice,
		Stock:         int32(d.Stock),
		MinStock:      int32(d.MinStock),
		ExpiryDate:    d.ExpiryDate,
		Status:        int32(d.Status),
	}
}
