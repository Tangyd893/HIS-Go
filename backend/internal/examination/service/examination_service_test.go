package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/examination/model"
	"his-go/internal/examination/repository"
)

func setupExaminationService(t *testing.T) (*ExaminationService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS examination_reports (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			patient_name TEXT,
			exam_request_id TEXT,
			exam_type TEXT,
			exam_item TEXT,
			body_part TEXT,
			findings TEXT,
			impression TEXT,
			conclusion TEXT,
			technician_id TEXT,
			reviewer_id TEXT,
			status INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	repo := repository.NewExaminationRepository(db)
	return NewExaminationService(repo), db
}

// ==================== CreateReport 测试 ====================

func TestCreateReport_SetsStatusZero(t *testing.T) {
	svc, db := setupExaminationService(t)

	report := &model.ExaminationReport{
		ID:        "er-1",
		PatientID: "patient-1",
		ExamType:  "CT",
		ExamItem:  "头颅CT",
		Status:    99,
	}

	err := svc.CreateReport(report)
	if err != nil {
		t.Fatalf("创建报告失败: %v", err)
	}

	if report.Status != 0 {
		t.Errorf("创建后状态期望=0(待检查)，实际=%d", report.Status)
	}

	var saved model.ExaminationReport
	db.Where("id = ?", "er-1").First(&saved)
	if saved.Status != 0 {
		t.Errorf("数据库中状态期望=0，实际=%d", saved.Status)
	}
}

// ==================== UpdateReport 测试 ====================

func TestUpdateReport_StatusZeroAllowsUpdate(t *testing.T) {
	svc, db := setupExaminationService(t)

	existing := &model.ExaminationReport{
		ID:        "er-upd-1",
		PatientID: "patient-1",
		ExamType:  "CT",
		Status:    0,
	}
	if err := db.Create(existing).Error; err != nil {
		t.Fatalf("创建现有报告失败: %v", err)
	}

	update := &model.ExaminationReport{
		ID:         "er-upd-1",
		PatientID:  "patient-1",
		Findings:   "右肺上叶结节",
		Impression: "考虑良性结节",
		Status:     0,
	}

	err := svc.UpdateReport(update)
	if err != nil {
		t.Fatalf("期望待检查状态允许更新，错误: %v", err)
	}

	if update.Status != 1 {
		t.Errorf("更新后状态期望=1(已检查)，实际=%d", update.Status)
	}

	var saved model.ExaminationReport
	db.Where("id = ?", "er-upd-1").First(&saved)
	if saved.Status != 1 {
		t.Errorf("数据库中状态期望=1，实际=%d", saved.Status)
	}
}

func TestUpdateReport_NonZeroStatusRejectsUpdate(t *testing.T) {
	svc, db := setupExaminationService(t)

	existing := &model.ExaminationReport{
		ID:        "er-upd-2",
		PatientID: "patient-2",
		ExamType:  "MR",
		Status:    1,
	}
	if err := db.Create(existing).Error; err != nil {
		t.Fatalf("创建现有报告失败: %v", err)
	}

	update := &model.ExaminationReport{
		ID:       "er-upd-2",
		Findings: "修改已检查报告",
	}

	err := svc.UpdateReport(update)
	if err == nil {
		t.Error("期望非待检查状态(Status=1)拒绝更新")
	}
}

func TestUpdateReport_StatusTwoRejectsUpdate(t *testing.T) {
	svc, db := setupExaminationService(t)

	existing := &model.ExaminationReport{
		ID:        "er-upd-3",
		PatientID: "patient-3",
		ExamType:  "DR",
		Status:    2,
	}
	if err := db.Create(existing).Error; err != nil {
		t.Fatalf("创建现有报告失败: %v", err)
	}

	update := &model.ExaminationReport{
		ID:       "er-upd-3",
		Findings: "修改已审核报告",
	}

	err := svc.UpdateReport(update)
	if err == nil {
		t.Error("期望已审核状态(Status=2)拒绝更新")
	}
}

func TestUpdateReport_StatusThreeRejectsUpdate(t *testing.T) {
	svc, db := setupExaminationService(t)

	existing := &model.ExaminationReport{
		ID:        "er-upd-4",
		PatientID: "patient-4",
		ExamType:  "超声",
		Status:    3,
	}
	if err := db.Create(existing).Error; err != nil {
		t.Fatalf("创建现有报告失败: %v", err)
	}

	update := &model.ExaminationReport{
		ID:       "er-upd-4",
		Findings: "修改已发布报告",
	}

	err := svc.UpdateReport(update)
	if err == nil {
		t.Error("期望已发布状态(Status=3)拒绝更新")
	}
}

func TestUpdateReport_NotFound(t *testing.T) {
	svc, _ := setupExaminationService(t)

	update := &model.ExaminationReport{
		ID:       "non-existent",
		Findings: "不存在的报告",
	}

	err := svc.UpdateReport(update)
	if err == nil {
		t.Error("期望更新不存在的报告返回错误")
	}
}
