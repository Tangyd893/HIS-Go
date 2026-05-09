package service

import (
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/registration/model"
	"his-go/internal/registration/repository"
	"his-go/pkg/common/snowflake"
	"his-go/pkg/errors"
	rediscli "his-go/pkg/redis"
)

func setupRegistrationService(t *testing.T) (*RegistrationService, *gorm.DB, *rediscli.Client, *miniredis.Miniredis, func()) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}
	if err := db.AutoMigrate(&model.Registration{}, &model.Schedule{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	mr := miniredis.RunT(t)
	rdb := &rediscli.Client{
		Client: goredis.NewClient(&goredis.Options{Addr: mr.Addr()}),
	}

	sf, err := snowflake.NewSnowflake(1, 1)
	if err != nil {
		t.Fatalf("创建雪化算法失败: %v", err)
	}

	repo := repository.NewRegistrationRepository(db, rdb)
	svc := NewRegistrationService(repo, rdb, sf)

	cleanup := func() {
		mr.Close()
	}

	return svc, db, rdb, mr, cleanup
}

func createTestSchedule(t *testing.T, db *gorm.DB, scheduleID string, deptID string, date string, total int, remain int) {
	t.Helper()
	schedule := model.Schedule{
		ID:          scheduleID,
		DeptID:      deptID,
		DeptName:    "测试科室",
		DoctorID:    "doc-1",
		DoctorName:  "测试医生",
		Date:        date,
		TimeSlot:    1,
		TotalCount:  total,
		RemainCount: remain,
		Status:      1,
	}
	if err := db.Create(&schedule).Error; err != nil {
		t.Fatalf("创建测试号源失败: %v", err)
	}
}

func TestRegistrationService_Register_Success(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	scheduleID := "schedule-001"
	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, scheduleID, "dept-1", date, 10, 10)

	reg, err := svc.Register("patient-1", "张三", scheduleID, date)
	if err != nil {
		t.Fatalf("挂号失败: %v", err)
	}
	if reg.ID == "" {
		t.Error("期望非空 ID")
	}
	if reg.PatientID != "patient-1" {
		t.Errorf("期望 PatientID='patient-1'，实际=%s", reg.PatientID)
	}
	if reg.PatientName != "张三" {
		t.Errorf("期望 PatientName='张三'，实际=%s", reg.PatientName)
	}
	if reg.QueueNumber != 1 {
		t.Errorf("期望 QueueNumber=1，实际=%d", reg.QueueNumber)
	}
	if reg.Status != 0 {
		t.Errorf("期望 Status=0(已预约)，实际=%d", reg.Status)
	}

	// 验证号源已扣减
	result, err := svc.ListSchedules("dept-1", date)
	if err != nil {
		t.Fatalf("查询号源失败: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("期望 1 个号源，实际=%d", len(result))
	}
	if result[0].RemainCount != 9 {
		t.Errorf("期望 RemainCount=9，实际=%d", result[0].RemainCount)
	}
}

func TestRegistrationService_Register_ScheduleFull(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	scheduleID := "schedule-full"
	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, scheduleID, "dept-1", date, 10, 0)

	_, err := svc.Register("patient-1", "张三", scheduleID, date)
	if err == nil {
		t.Error("期望号源已满时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Errorf("期望 AppError 类型，实际=%T", err)
	} else if appErr.Code != errors.CodeScheduleFull {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeScheduleFull, appErr.Code)
	}
}

func TestRegistrationService_Register_NonexistentSchedule(t *testing.T) {
	svc, _, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	_, err := svc.Register("patient-1", "张三", "nonexistent", "2026-01-01")
	if err == nil {
		t.Error("期望号源不存在时返回错误")
	}
}

func TestRegistrationService_Cancel_Success(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	scheduleID := "schedule-cancel"
	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, scheduleID, "dept-1", date, 10, 10)

	reg, err := svc.Register("patient-1", "张三", scheduleID, date)
	if err != nil {
		t.Fatalf("挂号失败: %v", err)
	}

	err = svc.Cancel(reg.ID)
	if err != nil {
		t.Fatalf("取消挂号失败: %v", err)
	}

	// 验证号源已恢复
	schedules, _ := svc.ListSchedules("dept-1", date)
	if schedules[0].RemainCount != 10 {
		t.Errorf("期望 RemainCount=10（恢复后），实际=%d", schedules[0].RemainCount)
	}

	// 验证挂号状态已取消
	cancelled, _ := svc.GetByID(reg.ID)
	if cancelled.Status != 3 {
		t.Errorf("期望 Status=3(已取消)，实际=%d", cancelled.Status)
	}
}

func TestRegistrationService_Cancel_NotFound(t *testing.T) {
	svc, _, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	err := svc.Cancel("nonexistent")
	if err == nil {
		t.Error("期望不存在的挂号记录时返回错误")
	}
}

func TestRegistrationService_Cancel_AlreadyCancelled(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	scheduleID := "schedule-double-cancel"
	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, scheduleID, "dept-1", date, 10, 10)

	reg, err := svc.Register("patient-1", "张三", scheduleID, date)
	if err != nil {
		t.Fatalf("挂号失败: %v", err)
	}

	err = svc.Cancel(reg.ID)
	if err != nil {
		t.Fatalf("首次取消失败: %v", err)
	}

	err = svc.Cancel(reg.ID)
	if err == nil {
		t.Error("期望重复取消时返回错误")
	}
}

func TestRegistrationService_SignIn_Success(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	scheduleID := "schedule-signin"
	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, scheduleID, "dept-1", date, 10, 10)

	reg, err := svc.Register("patient-1", "张三", scheduleID, date)
	if err != nil {
		t.Fatalf("挂号失败: %v", err)
	}

	err = svc.SignIn(reg.ID)
	if err != nil {
		t.Fatalf("签到失败: %v", err)
	}

	updated, _ := svc.GetByID(reg.ID)
	if updated.Status != 1 {
		t.Errorf("期望 Status=1(已签到)，实际=%d", updated.Status)
	}
}

func TestRegistrationService_SignIn_NotFound(t *testing.T) {
	svc, _, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	err := svc.SignIn("nonexistent")
	if err == nil {
		t.Error("期望不存在的挂号记录时返回错误")
	}
}

func TestRegistrationService_GetQueueStatus(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	scheduleID := "schedule-queue"
	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, scheduleID, "dept-1", date, 10, 10)

	reg1, _ := svc.Register("patient-1", "张三", scheduleID, date)
	reg2, _ := svc.Register("patient-2", "李四", scheduleID, date)

	if reg1.QueueNumber != 1 {
		t.Errorf("期望 reg1.QueueNumber=1，实际=%d", reg1.QueueNumber)
	}
	if reg2.QueueNumber != 2 {
		t.Errorf("期望 reg2.QueueNumber=2，实际=%d", reg2.QueueNumber)
	}

	// 检查排队位置
	rank, err := svc.GetQueueStatus(reg2.ID)
	if err != nil {
		t.Fatalf("查询排队状态失败: %v", err)
	}
	if rank < 0 {
		t.Errorf("期望 rank>=0，实际=%d", rank)
	}
}

func TestRegistrationService_CallNext_Success(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	scheduleID := "schedule-callnext"
	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, scheduleID, "dept-1", date, 10, 10)

	reg, err := svc.Register("patient-1", "张三", scheduleID, date)
	if err != nil {
		t.Fatalf("挂号失败: %v", err)
	}

	err = svc.SignIn(reg.ID)
	if err != nil {
		t.Fatalf("签到失败: %v", err)
	}

	nextID, err := svc.CallNext(scheduleID)
	if err != nil {
		t.Fatalf("叫号失败: %v", err)
	}
	if nextID != reg.ID {
		t.Errorf("期望叫号 ID=%s，实际=%s", reg.ID, nextID)
	}
}

func TestRegistrationService_ListByPatient(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, "schedule-a", "dept-1", date, 10, 10)
	createTestSchedule(t, db, "schedule-b", "dept-1", date, 10, 10)

	_, err := svc.Register("patient-x", "王五", "schedule-a", date)
	if err != nil {
		t.Fatalf("挂号失败: %v", err)
	}
	_, err = svc.Register("patient-x", "王五", "schedule-b", date)
	if err != nil {
		t.Fatalf("挂号失败: %v", err)
	}

	list, total, err := svc.ListByPatient("patient-x", 1, 10)
	if err != nil {
		t.Fatalf("查询患者挂号失败: %v", err)
	}
	if total != 2 {
		t.Errorf("期望 total=2，实际=%d", total)
	}
	if len(list) != 2 {
		t.Errorf("期望 len(list)=2，实际=%d", len(list))
	}
}

func TestRegistrationService_Register_ConcurrentLock(t *testing.T) {
	svc, db, _, _, cleanup := setupRegistrationService(t)
	defer cleanup()

	scheduleID := "schedule-concurrent"
	date := time.Now().Format("2006-01-02")
	createTestSchedule(t, db, scheduleID, "dept-1", date, 10, 10)

	var wg sync.WaitGroup
	var successCount int32
	var mu sync.Mutex
	successPatients := make(map[string]bool)

	// 串行执行挂号（每个协程带重试等待锁释放）
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			patientID := "patient-" + string(rune('a'+idx))
			// 重试最多 5 次，每次等 50ms
			for retry := 0; retry < 5; retry++ {
				reg, err := svc.Register(patientID, "测试", scheduleID, date)
				if err == nil && reg != nil {
					mu.Lock()
					successCount++
					successPatients[patientID] = true
					mu.Unlock()
					return
				}
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}
	wg.Wait()

	if successCount < 3 {
		t.Errorf("期望至少 3 次成功，实际=%d", successCount)
	}

	schedules, _ := svc.ListSchedules("dept-1", date)
	expectedRemain := int32(10 - successCount)
	if schedules[0].RemainCount != int(expectedRemain) {
		t.Errorf("期望 RemainCount=%d，实际=%d", expectedRemain, schedules[0].RemainCount)
	}
}
