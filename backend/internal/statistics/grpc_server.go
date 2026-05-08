package statistics

import (
	"context"

	"his-go/api/proto/common"
	"his-go/api/proto/statistics"
	statssvc "his-go/internal/statistics/service"
)

// StatisticsGrpcServer gRPC 数据统计服务实现
type StatisticsGrpcServer struct {
	statistics.UnimplementedStatisticsServiceServer
	svc *statssvc.StatisticsService
}

// NewStatisticsGrpcServer 创建 gRPC 数据统计服务
func NewStatisticsGrpcServer(svc *statssvc.StatisticsService) *StatisticsGrpcServer {
	return &StatisticsGrpcServer{svc: svc}
}

// GetOperationStats 获取运营统计
func (s *StatisticsGrpcServer) GetOperationStats(ctx context.Context, req *statistics.StatsRequest) (*statistics.OperationStats, error) {
	result, err := s.svc.GetOperationStats(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	return &statistics.OperationStats{
		TotalRegistrations: int32(result.Registrations),
		TotalClinicVisits:  int32(result.Visits),
		TotalPrescriptions: int32(result.Prescriptions),
		TotalRevenue:       result.TotalRevenue,
	}, nil
}

// GetDeptWorkload 获取科室工作量
func (s *StatisticsGrpcServer) GetDeptWorkload(ctx context.Context, req *statistics.StatsRequest) (*statistics.DeptWorkloadResponse, error) {
	workloads, err := s.svc.GetDeptWorkload(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	pbList := make([]*statistics.DeptWorkload, len(workloads))
	for i, w := range workloads {
		pbList[i] = &statistics.DeptWorkload{
			DeptId:     w.DeptID,
			DeptName:   w.DeptName,
			VisitCount: int32(w.VisitCount),
			Revenue:    w.Revenue,
		}
	}
	return &statistics.DeptWorkloadResponse{
		Base:      &common.BaseResponse{Code: 0, Message: "查询成功"},
		Workloads: pbList,
	}, nil
}

// GetRevenueTrend 获取收入趋势
func (s *StatisticsGrpcServer) GetRevenueTrend(ctx context.Context, req *statistics.StatsRequest) (*statistics.RevenueTrendResponse, error) {
	trends, err := s.svc.GetRevenueTrend(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	pbList := make([]*statistics.RevenueTrend, len(trends))
	for i, t := range trends {
		pbList[i] = &statistics.RevenueTrend{
			Date:   t.Date,
			Amount: t.Revenue,
		}
	}
	return &statistics.RevenueTrendResponse{
		Base:   &common.BaseResponse{Code: 0, Message: "查询成功"},
		Trends: pbList,
	}, nil
}
