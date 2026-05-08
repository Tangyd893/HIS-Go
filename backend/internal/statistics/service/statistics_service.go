package service

import (
	"his-go/internal/statistics/repository"
)

// StatisticsService 数据统计业务服务
type StatisticsService struct {
	repo *repository.StatisticsRepository
}

// NewStatisticsService 创建数据统计业务服务
func NewStatisticsService(repo *repository.StatisticsRepository) *StatisticsService {
	return &StatisticsService{repo: repo}
}

// GetOperationStats 获取运营统计
func (s *StatisticsService) GetOperationStats(startDate, endDate string) (*repository.OperationStatsResult, error) {
	return s.repo.GetOperationStats(startDate, endDate)
}

// GetDeptWorkload 获取科室工作负载
func (s *StatisticsService) GetDeptWorkload(startDate, endDate string) ([]repository.DeptWorkloadResult, error) {
	return s.repo.GetDeptWorkload(startDate, endDate)
}

// GetRevenueTrend 获取收入趋势
func (s *StatisticsService) GetRevenueTrend(startDate, endDate string) ([]repository.RevenueTrendResult, error) {
	return s.repo.GetRevenueTrend(startDate, endDate)
}
