package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/followup/model"
	"his-go/internal/followup/repository"
)

func setupFollowupService(t *testing.T) (*FollowupService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS followup_plans (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			template_id TEXT,
			plan_name TEXT,
			start_date TEXT,
			end_date TEXT,
			frequency INTEGER NOT NULL,
			status INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS followup_tasks (
			id TEXT PRIMARY KEY,
			plan_id TEXT NOT NULL,
			assignee_id TEXT,
			execute_date TEXT NOT NULL,
			type INTEGER DEFAULT 1,
			status INTEGER DEFAULT 0,
			content TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS satisfaction_surveys (
			id TEXT PRIMARY KEY,
			followup_task_id TEXT NOT NULL,
			patient_id TEXT NOT NULL,
			score INTEGER NOT NULL,
			feedback TEXT,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	repo := repository.NewFollowupRepository(db)
	svc := NewFollowupService(repo)

	return svc, db
}

// ==================== GenerateTasks 测试 ====================

func TestGenerateTasks_Weekly(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-05-01",
		EndDate:   "2026-05-31",
		Frequency: 1,
		PlanName:  "术后随访",
	}

	tasks, err := svc.GenerateTasks(plan)
	if err != nil {
		t.Fatalf("生成任务失败: %v", err)
	}

	if len(tasks) < 4 {
		t.Errorf("每周随访期望至少4条任务（4周），实际=%d", len(tasks))
	}

	if tasks[0].ExecuteDate != "2026-05-01" {
		t.Errorf("第一条任务日期期望=2026-05-01，实际=%s", tasks[0].ExecuteDate)
	}

	for _, task := range tasks {
		if task.Status != 0 {
			t.Errorf("任务状态期望=0(待执行)，实际=%d", task.Status)
		}
		if task.Type != 1 {
			t.Errorf("任务类型期望=1(电话随访)，实际=%d", task.Type)
		}
	}
}

func TestGenerateTasks_Biweekly(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-05-01",
		EndDate:   "2026-05-31",
		Frequency: 2,
		PlanName:  "慢病随访",
	}

	tasks, err := svc.GenerateTasks(plan)
	if err != nil {
		t.Fatalf("生成任务失败: %v", err)
	}

	if len(tasks) < 2 {
		t.Errorf("每两周随访期望至少2条任务，实际=%d", len(tasks))
	}

	firstDate := tasks[0].ExecuteDate
	secondDate := tasks[1].ExecuteDate
	if firstDate != "2026-05-01" {
		t.Errorf("第一条任务日期期望=2026-05-01，实际=%s", firstDate)
	}
	if secondDate != "2026-05-15" {
		t.Errorf("第二条任务日期期望=2026-05-15（14天后），实际=%s", secondDate)
	}
}

func TestGenerateTasks_Monthly(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-01-15",
		EndDate:   "2026-05-15",
		Frequency: 3,
		PlanName:  "长期随访",
	}

	tasks, err := svc.GenerateTasks(plan)
	if err != nil {
		t.Fatalf("生成任务失败: %v", err)
	}

	if len(tasks) < 4 {
		t.Errorf("每月随访4个月期望至少4条任务，实际=%d", len(tasks))
	}

	expectedDates := []string{"2026-01-15", "2026-02-15", "2026-03-15", "2026-04-15", "2026-05-15"}
	for i, expected := range expectedDates {
		if i >= len(tasks) {
			break
		}
		if tasks[i].ExecuteDate != expected {
			t.Errorf("第%d条任务日期期望=%s，实际=%s", i+1, expected, tasks[i].ExecuteDate)
		}
	}
}

func TestGenerateTasks_StartEqualsEnd(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-05-10",
		EndDate:   "2026-05-10",
		Frequency: 1,
	}

	tasks, err := svc.GenerateTasks(plan)
	if err != nil {
		t.Fatalf("开始等于结束时应生成1条任务，错误: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("期望1条任务，实际=%d", len(tasks))
	}
	if tasks[0].ExecuteDate != "2026-05-10" {
		t.Errorf("任务日期期望=2026-05-10，实际=%s", tasks[0].ExecuteDate)
	}
}

func TestGenerateTasks_StartAfterEnd(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-05-15",
		EndDate:   "2026-05-10",
		Frequency: 1,
	}

	_, err := svc.GenerateTasks(plan)
	if err == nil {
		t.Error("期望开始日期晚于结束日期时返回错误")
	}
}

func TestGenerateTasks_InvalidStartDate(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-13-01",
		EndDate:   "2026-05-10",
		Frequency: 1,
	}

	_, err := svc.GenerateTasks(plan)
	if err == nil {
		t.Error("期望无效开始日期格式时返回错误")
	}
}

func TestGenerateTasks_InvalidEndDate(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-05-01",
		EndDate:   "not-a-date",
		Frequency: 1,
	}

	_, err := svc.GenerateTasks(plan)
	if err == nil {
		t.Error("期望无效结束日期格式时返回错误")
	}
}

func TestGenerateTasks_ZeroFrequencyDefaultsToOne(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-05-01",
		EndDate:   "2026-05-15",
		Frequency: 0,
	}

	tasks, err := svc.GenerateTasks(plan)
	if err != nil {
		t.Fatalf("频率为0时应默认为1，错误: %v", err)
	}
	if len(tasks) < 2 {
		t.Errorf("频率默认1时两周期望至少2条任务，实际=%d", len(tasks))
	}
}

func TestGenerateTasks_FrequencyExceedsThree(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-05-01",
		EndDate:   "2026-05-15",
		Frequency: 5,
	}

	tasks, err := svc.GenerateTasks(plan)
	if err != nil {
		t.Fatalf("频率超过3时应默认为1，错误: %v", err)
	}
	if len(tasks) < 2 {
		t.Errorf("频率默认1时两周期望至少2条任务，实际=%d", len(tasks))
	}
}

func TestGenerateTasks_ContentSet(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		StartDate: "2026-05-01",
		EndDate:   "2026-05-01",
		Frequency: 1,
		PlanName:  "糖尿病随访",
	}

	tasks, err := svc.GenerateTasks(plan)
	if err != nil {
		t.Fatalf("生成任务失败: %v", err)
	}

	if tasks[0].Content != "糖尿病随访" {
		t.Errorf("任务内容期望='糖尿病随访'，实际=%s", tasks[0].Content)
	}
}

// ==================== SubmitSurvey 测试 ====================

func TestSubmitSurvey_ValidScore(t *testing.T) {
	svc, db := setupFollowupService(t)

	validScores := []int{1, 2, 3, 4, 5}
	for _, score := range validScores {
		survey := &model.SatisfactionSurvey{
			ID:             "survey-" + string(rune('0'+score)),
			FollowupTaskID: "task-1",
			PatientID:      "patient-1",
			Score:          score,
			Feedback:       "服务很好",
		}

		err := svc.SubmitSurvey(survey)
		if err != nil {
			t.Errorf("评分%d应通过，错误: %v", score, err)
		}

		var count int64
		db.Model(&model.SatisfactionSurvey{}).Where("id = ?", survey.ID).Count(&count)
		if count != 1 {
			t.Errorf("评分%d的调查应被保存，但未找到", score)
		}
	}
}

func TestSubmitSurvey_ScoreTooLow(t *testing.T) {
	svc, _ := setupFollowupService(t)

	survey := &model.SatisfactionSurvey{
		ID:             "survey-low",
		FollowupTaskID: "task-1",
		PatientID:      "patient-1",
		Score:          0,
	}

	err := svc.SubmitSurvey(survey)
	if err == nil {
		t.Error("期望评分0时返回错误")
	}
}

func TestSubmitSurvey_ScoreTooHigh(t *testing.T) {
	svc, _ := setupFollowupService(t)

	survey := &model.SatisfactionSurvey{
		ID:             "survey-high",
		FollowupTaskID: "task-1",
		PatientID:      "patient-1",
		Score:          6,
	}

	err := svc.SubmitSurvey(survey)
	if err == nil {
		t.Error("期望评分6时返回错误")
	}
}

func TestSubmitSurvey_NegativeScore(t *testing.T) {
	svc, _ := setupFollowupService(t)

	survey := &model.SatisfactionSurvey{
		ID:             "survey-neg",
		FollowupTaskID: "task-1",
		PatientID:      "patient-1",
		Score:          -1,
	}

	err := svc.SubmitSurvey(survey)
	if err == nil {
		t.Error("期望负分时返回错误")
	}
}

// ==================== CreatePlan 测试 ====================

func TestCreatePlan_SetsStatusToOne(t *testing.T) {
	svc, db := setupFollowupService(t)

	plan := &model.FollowupPlan{
		ID:        "plan-1",
		PatientID: "patient-1",
		PlanName:  "术后随访",
		StartDate: "2026-05-01",
		EndDate:   "2026-05-07",
		Frequency: 1,
		Status:    0,
	}

	err := svc.CreatePlan(plan)
	if err != nil {
		t.Fatalf("创建计划失败: %v", err)
	}

	if plan.Status != 1 {
		t.Errorf("创建后状态期望=1(进行中)，实际=%d", plan.Status)
	}

	var saved model.FollowupPlan
	db.Where("id = ?", "plan-1").First(&saved)
	if saved.Status != 1 {
		t.Errorf("数据库中状态期望=1，实际=%d", saved.Status)
	}
}

func TestCreatePlan_GeneratesTasks(t *testing.T) {
	svc, db := setupFollowupService(t)

	plan := &model.FollowupPlan{
		ID:        "plan-2",
		PatientID: "patient-2",
		PlanName:  "慢病随访",
		StartDate: "2026-05-01",
		EndDate:   "2026-05-21",
		Frequency: 1,
	}

	err := svc.CreatePlan(plan)
	if err != nil {
		t.Fatalf("创建计划失败: %v", err)
	}

	var taskCount int64
	db.Model(&model.FollowupTask{}).Where("plan_id = ?", "plan-2").Count(&taskCount)
	if taskCount < 3 {
		t.Errorf("3周每周1次期望至少3条任务，实际=%d", taskCount)
	}
}

func TestCreatePlan_InvalidDateFails(t *testing.T) {
	svc, _ := setupFollowupService(t)

	plan := &model.FollowupPlan{
		ID:        "plan-3",
		PatientID: "patient-3",
		PlanName:  "无效计划",
		StartDate: "2026-05-15",
		EndDate:   "2026-05-10",
		Frequency: 1,
	}

	err := svc.CreatePlan(plan)
	if err == nil {
		t.Error("期望开始日期晚于结束日期时返回错误")
	}
}
