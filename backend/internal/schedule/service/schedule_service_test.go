package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/schedule/model"
	"his-go/internal/schedule/repository"
)

func setupScheduleService(t *testing.T) (*ScheduleService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS schedules (
			id TEXT PRIMARY KEY,
			doctor_id TEXT NOT NULL,
			doctor_name TEXT,
			dept_id TEXT NOT NULL,
			dept_name TEXT,
			work_date TEXT NOT NULL,
			time_slot INTEGER NOT NULL,
			max_patients INTEGER NOT NULL,
			current_patients INTEGER DEFAULT 0,
			room_no TEXT,
			status INTEGER DEFAULT 1,
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

	repo := repository.NewScheduleRepository(db)
	svc := NewScheduleService(repo, nil)

	return svc, db
}

func TestScheduleService_GenerateWeekly(t *testing.T) {
	svc, db := setupScheduleService(t)

	schedules, err := svc.GenerateWeekly("2026-05-11", "2026-05-13", "dept-1")
	if err != nil {
		t.Fatalf("生成排班失败: %v", err)
	}

	if len(schedules) == 0 {
		t.Fatal("期望生成排班记录")
	}

	expectedCount := 9 // 3天 × 3时段
	if len(schedules) != expectedCount {
		t.Errorf("期望%d条排班记录，实际=%d", expectedCount, len(schedules))
	}

	var count int64
	db.Model(&model.ScheduleInfo{}).Count(&count)
	if count != int64(expectedCount) {
		t.Errorf("期望数据库中有%d条，实际=%d", expectedCount, count)
	}
}

func TestScheduleService_GenerateWeekly_InvalidDate(t *testing.T) {
	svc, _ := setupScheduleService(t)

	_, err := svc.GenerateWeekly("2026-05-15", "2026-05-10", "dept-1")
	if err == nil {
		t.Error("期望开始日期晚于结束日期时返回错误")
	}
}

func TestScheduleService_ListByDept(t *testing.T) {
	svc, db := setupScheduleService(t)

	entries := []model.ScheduleInfo{
		{ID: "s-1", DoctorID: "doc-1", DeptID: "dept-A", WorkDate: "2026-05-11", TimeSlot: 1, MaxPatients: 30, Status: 1},
		{ID: "s-2", DoctorID: "doc-2", DeptID: "dept-A", WorkDate: "2026-05-11", TimeSlot: 2, MaxPatients: 25, Status: 1},
		{ID: "s-3", DoctorID: "doc-3", DeptID: "dept-B", WorkDate: "2026-05-11", TimeSlot: 1, MaxPatients: 20, Status: 1},
	}
	for _, e := range entries {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("创建排班失败: %v", err)
		}
	}

	list, err := svc.ListByDept("dept-A", "2026-05-11")
	if err != nil {
		t.Fatalf("查询排班失败: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("期望2条排班记录，实际=%d", len(list))
	}
}

func TestScheduleService_ListByDoctor(t *testing.T) {
	svc, db := setupScheduleService(t)

	entries := []model.ScheduleInfo{
		{ID: "s-4", DoctorID: "doc-x", DeptID: "dept-C", WorkDate: "2026-05-12", TimeSlot: 1, MaxPatients: 20, Status: 1},
		{ID: "s-5", DoctorID: "doc-x", DeptID: "dept-C", WorkDate: "2026-05-12", TimeSlot: 3, MaxPatients: 15, Status: 1},
	}
	for _, e := range entries {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("创建排班失败: %v", err)
		}
	}

	list, err := svc.ListByDoctor("doc-x", "2026-05-12")
	if err != nil {
		t.Fatalf("查询排班失败: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("期望2条排班记录，实际=%d", len(list))
	}
}

func TestScheduleService_UpdateSchedule(t *testing.T) {
	svc, db := setupScheduleService(t)

	entry := &model.ScheduleInfo{
		ID: "s-upd", DoctorID: "doc-1", DeptID: "dept-1", WorkDate: "2026-05-13", TimeSlot: 1, MaxPatients: 50, Status: 1,
	}
	if err := db.Create(entry).Error; err != nil {
		t.Fatalf("创建排班失败: %v", err)
	}

	entry.MaxPatients = 80
	entry.RoomNo = "302室"
	err := svc.UpdateSchedule(entry)
	if err != nil {
		t.Fatalf("更新排班失败: %v", err)
	}

	var updated model.ScheduleInfo
	db.Where("id = ?", "s-upd").First(&updated)
	if updated.MaxPatients != 80 {
		t.Errorf("期望MaxPatients=80，实际=%d", updated.MaxPatients)
	}
}

func TestScheduleService_CancelSchedule(t *testing.T) {
	svc, db := setupScheduleService(t)

	entry := &model.ScheduleInfo{
		ID: "s-cancel", DoctorID: "doc-1", DeptID: "dept-1", WorkDate: "2026-05-14", TimeSlot: 1, MaxPatients: 30, Status: 1,
	}
	if err := db.Create(entry).Error; err != nil {
		t.Fatalf("创建排班失败: %v", err)
	}

	err := svc.CancelSchedule("s-cancel")
	if err != nil {
		t.Fatalf("取消排班失败: %v", err)
	}

	var cancelled model.ScheduleInfo
	db.Where("id = ?", "s-cancel").First(&cancelled)
	if cancelled.Status != 0 {
		t.Errorf("期望Status=0(已停诊)，实际=%d", cancelled.Status)
	}
}
