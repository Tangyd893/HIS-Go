package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/statistics/repository"
)

func TestNewStatisticsService(t *testing.T) {
	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	repo := repository.NewStatisticsRepository(db)
	svc := NewStatisticsService(repo)

	if svc == nil {
		t.Fatal("NewStatisticsService 返回 nil")
	}
	if svc.repo == nil {
		t.Fatal("StatisticsService.repo 为 nil")
	}
}

func setupStatsDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}
	// 统计数据跨多个业务表，建最小表结构以通过查询
	tables := []string{
		`CREATE TABLE IF NOT EXISTS registrations (id TEXT, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS clinic_records (id TEXT, dept_id TEXT, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS prescriptions (id TEXT, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS bills (id TEXT, related_id TEXT, total_amount REAL, status INTEGER, created_at DATETIME)`,
		`CREATE TABLE IF NOT EXISTS departments (id TEXT PRIMARY KEY, name TEXT)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v", err)
		}
	}
	return db
}

func TestStatisticsService_GetOperationStats_EmptyDB(t *testing.T) {
	db := setupStatsDB(t)
	repo := repository.NewStatisticsRepository(db)
	svc := NewStatisticsService(repo)

	stats, err := svc.GetOperationStats("2026-01-01", "2026-06-01")
	if err != nil {
		t.Fatalf("空库统计不应报错: %v", err)
	}
	if stats.Registrations != 0 {
		t.Errorf("空库期望 Registrations=0，实际=%d", stats.Registrations)
	}
	if stats.Visits != 0 {
		t.Errorf("空库期望 Visits=0，实际=%d", stats.Visits)
	}
	if stats.TotalRevenue != 0 {
		t.Errorf("空库期望 TotalRevenue=0，实际=%f", stats.TotalRevenue)
	}
}

func TestStatisticsService_GetDeptWorkload_EmptyDB(t *testing.T) {
	db := setupStatsDB(t)
	repo := repository.NewStatisticsRepository(db)
	svc := NewStatisticsService(repo)

	results, err := svc.GetDeptWorkload("2026-01-01", "2026-06-01")
	if err != nil {
		t.Fatalf("空库科室统计不应报错: %v", err)
	}
	// 空 departments 表 → 0 条结果
	if len(results) != 0 {
		t.Errorf("空库期望 0 条结果，实际=%d", len(results))
	}
}

func TestStatisticsService_GetRevenueTrend_EmptyDB(t *testing.T) {
	db := setupStatsDB(t)
	repo := repository.NewStatisticsRepository(db)
	svc := NewStatisticsService(repo)

	results, err := svc.GetRevenueTrend("2026-01-01", "2026-06-01")
	if err != nil {
		t.Fatalf("空库收入趋势不应报错: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("空库期望 0 条趋势，实际=%d", len(results))
	}
}
