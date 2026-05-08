package system

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/system"
	sysmodel "his-go/internal/system/model"
	syssvc "his-go/internal/system/service"
)

// SystemGrpcServer gRPC 系统管理服务实现
type SystemGrpcServer struct {
	system.UnimplementedSystemServiceServer
	svc *syssvc.SystemService
}

// NewSystemGrpcServer 创建 gRPC 系统管理服务
func NewSystemGrpcServer(svc *syssvc.SystemService) *SystemGrpcServer {
	return &SystemGrpcServer{svc: svc}
}

// ListDictTypes 查询字典类型列表
func (s *SystemGrpcServer) ListDictTypes(ctx context.Context, req *common.Empty) (*system.DictTypeListResponse, error) {
	_ = req
	types, err := s.svc.ListDictTypes()
	if err != nil {
		return nil, err
	}
	pbList := make([]*system.DictType, len(types))
	for i, t := range types {
		pbList[i] = &system.DictType{
			Id:       t.ID,
			DictName: t.DictName,
			DictType: t.DictType,
			Status:   int32(t.Status),
			Remark:   t.Remark,
		}
	}
	return &system.DictTypeListResponse{
		Base:  &common.BaseResponse{Code: 0, Message: "查询成功"},
		Types: pbList,
	}, nil
}

// ListDictItems 查询字典项列表
func (s *SystemGrpcServer) ListDictItems(ctx context.Context, req *system.DictItemQueryRequest) (*system.DictItemListResponse, error) {
	items, err := s.svc.ListDictItems(req.DictType)
	if err != nil {
		return nil, err
	}
	pbList := make([]*system.DictItem, len(items))
	for i, item := range items {
		pbList[i] = &system.DictItem{
			Id:        item.ID,
			DictType:  item.DictType,
			Label:     item.Label,
			Value:     item.Value,
			SortOrder: int32(item.SortOrder),
			Status:    int32(item.Status),
		}
	}
	return &system.DictItemListResponse{
		Base:  &common.BaseResponse{Code: 0, Message: "查询成功"},
		Items: pbList,
	}, nil
}

// CreateDictItem 创建字典项
func (s *SystemGrpcServer) CreateDictItem(ctx context.Context, req *system.DictItem) (*system.DictItem, error) {
	item := &sysmodel.DictItem{
		ID:        req.Id,
		DictType:  req.DictType,
		Label:     req.Label,
		Value:     req.Value,
		SortOrder: int(req.SortOrder),
		Status:    int8(req.Status),
	}
	if err := s.svc.CreateDictItem(item); err != nil {
		return nil, err
	}
	return req, nil
}

// ListParams 查询系统参数列表
func (s *SystemGrpcServer) ListParams(ctx context.Context, req *common.Empty) (*system.ParamListResponse, error) {
	_ = req
	params, err := s.svc.ListParams()
	if err != nil {
		return nil, err
	}
	pbList := make([]*system.SystemParam, len(params))
	for i, p := range params {
		pbList[i] = &system.SystemParam{
			Id:         p.ID,
			ParamName:  p.ParamName,
			ParamKey:   p.ParamKey,
			ParamValue: p.ParamValue,
			Remark:     p.Remark,
		}
	}
	return &system.ParamListResponse{
		Base:   &common.BaseResponse{Code: 0, Message: "查询成功"},
		Params: pbList,
	}, nil
}

// UpdateParam 更新系统参数
func (s *SystemGrpcServer) UpdateParam(ctx context.Context, req *system.SystemParam) (*system.SystemParam, error) {
	param := &sysmodel.SystemParam{
		ID:         req.Id,
		ParamName:  req.ParamName,
		ParamKey:   req.ParamKey,
		ParamValue: req.ParamValue,
		Remark:     req.Remark,
	}
	if err := s.svc.UpdateParam(param); err != nil {
		return nil, err
	}
	return req, nil
}

// ListOperationLogs 分页查询操作日志
func (s *SystemGrpcServer) ListOperationLogs(ctx context.Context, req *system.LogListRequest) (*system.LogListResponse, error) {
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
	logs, total, err := s.svc.ListOperationLogs(req.UserId, req.Module, page, pageSize)
	if err != nil {
		return nil, err
	}
	pbList := make([]*system.OperationLog, len(logs))
	for i, l := range logs {
		pbList[i] = &system.OperationLog{
			Id:        l.ID,
			UserId:    l.UserID,
			Username:  l.Username,
			Module:    l.Module,
			Action:    l.Action,
			Method:    l.Method,
			Url:       l.URL,
			Ip:        l.IP,
			Params:    l.Params,
			Result:    l.Result,
			Status:    int32(l.Status),
			CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &system.LogListResponse{
		Base: &common.BaseResponse{Code: 0, Message: "查询成功"},
		Logs: pbList,
		Page: &common.PageResponse{Total: total, Page: int32(page), PageSize: int32(pageSize)},
	}, nil
}
