package repository

import (
	"fmt"

	"gorm.io/gorm"
)

// StatisticsRepository 数据统计仓库（纯聚合查询，无持久化表）
type StatisticsRepository struct {
	db *gorm.DB
}

// NewStatisticsRepository 创建数据统计仓库
func NewStatisticsRepository(db *gorm.DB) *StatisticsRepository {
	return &StatisticsRepository{db: db}
}

// OperationStatsResult 运营统计结果
type OperationStatsResult struct {
	Registrations int64   `json:"registrations"`
	Visits        int64   `json:"visits"`
	Prescriptions int64   `json:"prescriptions"`
	TotalRevenue  float64 `json:"total_revenue"`
}

// DeptWorkloadResult 科室工作负载结果
type DeptWorkloadResult struct {
	DeptID     string  `json:"dept_id"`
	DeptName   string  `json:"dept_name"`
	VisitCount int64   `json:"visit_count"`
	Revenue    float64 `json:"revenue"`
}

// RevenueTrendResult 收入趋势结果
type RevenueTrendResult struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
}

// GetOperationStats 获取运营统计数据
func (r *StatisticsRepository) GetOperationStats(startDate, endDate string) (*OperationStatsResult, error) {
	var stats OperationStatsResult

	var regCount int64
	if err := r.db.Table("registrations").
		Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59").
		Count(&regCount).Error; err != nil {
		regCount = 0
	}
	stats.Registrations = regCount

	var visitCount int64
	if err := r.db.Table("clinic_records").
		Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59").
		Count(&visitCount).Error; err != nil {
		visitCount = 0
	}
	stats.Visits = visitCount

	var prescCount int64
	if err := r.db.Table("prescriptions").
		Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59").
		Count(&prescCount).Error; err != nil {
		prescCount = 0
	}
	stats.Prescriptions = prescCount

	var revenue float64
	if err := r.db.Table("bills").
		Where("created_at BETWEEN ? AND ? AND status = 1", startDate+" 00:00:00", endDate+" 23:59:59").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&revenue).Error; err != nil {
		revenue = 0
	}
	stats.TotalRevenue = revenue

	return &stats, nil
}

// GetDeptWorkload 按科室统计接诊量和收入
func (r *StatisticsRepository) GetDeptWorkload(startDate, endDate string) ([]DeptWorkloadResult, error) {
	var results []DeptWorkloadResult

	query := `
		SELECT d.id AS dept_id, d.name AS dept_name,
			COALESCE(COUNT(cr.id), 0) AS visit_count,
			COALESCE(SUM(b.total_amount), 0) AS revenue
		FROM departments d
		LEFT JOIN clinic_records cr ON d.id = cr.dept_id
			AND cr.created_at BETWEEN ? AND ?
		LEFT JOIN bills b ON cr.id = b.related_id
			AND b.created_at BETWEEN ? AND ?
		GROUP BY d.id, d.name
		ORDER BY revenue DESC`

	if err := r.db.Raw(query, startDate+" 00:00:00", endDate+" 23:59:59",
		startDate+" 00:00:00", endDate+" 23:59:59").Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("统计科室工作负载失败: %w", err)
	}

	return results, nil
}

// GetRevenueTrend 按日期统计收入趋势
func (r *StatisticsRepository) GetRevenueTrend(startDate, endDate string) ([]RevenueTrendResult, error) {
	var results []RevenueTrendResult

	query := `
		SELECT DATE(created_at) AS date,
			COALESCE(SUM(total_amount), 0) AS revenue
		FROM bills
		WHERE created_at BETWEEN ? AND ? AND status = 1
		GROUP BY DATE(created_at)
		ORDER BY date ASC`

	if err := r.db.Raw(query, startDate+" 00:00:00", endDate+" 23:59:59").Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("统计收入趋势失败: %w", err)
	}

	return results, nil
}
