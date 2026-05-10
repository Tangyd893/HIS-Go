package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/health_record/model"
	"his-go/internal/health_record/repository"
	"his-go/pkg/errors"
)

// newTestHealthRecordService 创建测试用的 HealthRecordService（基于 SQLite 内存数据库）
func newTestHealthRecordService(t *testing.T) (*HealthRecordService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	// SQLite 不支持 PostgreSQL 的 gen_random_uuid()，手动建表
	tables := []string{
		`CREATE TABLE IF NOT EXISTS health_record_summaries (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL UNIQUE,
			patient_name TEXT,
			total_visits INTEGER DEFAULT 0,
			total_prescriptions INTEGER DEFAULT 0,
			total_examinations INTEGER DEFAULT 0,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS timeline_events (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			date TEXT,
			event_type TEXT,
			description TEXT,
			related_id TEXT,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS record_authorizations (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			doctor_id TEXT NOT NULL,
			auth_time TEXT,
			expire_time TEXT,
			status INTEGER DEFAULT 1,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
	}

	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	// 为 timeline_events 创建索引
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_timeline_events_patient_id ON timeline_events(patient_id)`,
		`CREATE INDEX IF NOT EXISTS idx_record_authorizations_patient_id ON record_authorizations(patient_id)`,
		`CREATE INDEX IF NOT EXISTS idx_record_authorizations_doctor_id ON record_authorizations(doctor_id)`,
	}
	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			t.Fatalf("创建索引失败: %v\nSQL: %s", err, idx)
		}
	}

	repo := repository.NewHealthRecordRepository(db)
	svc := NewHealthRecordService(repo)

	return svc, db
}

// TestHealthRecordService_GetSummary_EmptyPatientID 患者ID为空时应返回参数无效错误
func TestHealthRecordService_GetSummary_EmptyPatientID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	summary, err := svc.GetSummary("")
	if err == nil {
		t.Error("期望患者ID为空时返回错误")
	}
	if summary != nil {
		t.Error("期望患者ID为空时返回 nil")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "患者ID不能为空" {
		t.Errorf("期望错误详情='患者ID不能为空'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_GetSummary_Success 正常获取档案摘要（无记录时自动创建）
func TestHealthRecordService_GetSummary_Success(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	summary, err := svc.GetSummary("P001")
	if err != nil {
		t.Fatalf("获取档案摘要失败: %v", err)
	}
	if summary.PatientID != "P001" {
		t.Errorf("期望 PatientID='P001'，实际='%s'", summary.PatientID)
	}
}

// TestHealthRecordService_GetTimeline_EmptyPatientID 患者ID为空时应返回参数无效错误
func TestHealthRecordService_GetTimeline_EmptyPatientID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	events, err := svc.GetTimeline("")
	if err == nil {
		t.Error("期望患者ID为空时返回错误")
	}
	if events != nil {
		t.Error("期望患者ID为空时返回 nil")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "患者ID不能为空" {
		t.Errorf("期望错误详情='患者ID不能为空'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_GetTimeline_Success 正常获取时间轴事件列表
func TestHealthRecordService_GetTimeline_Success(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	events, err := svc.GetTimeline("P001")
	if err != nil {
		t.Fatalf("获取时间轴失败: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("期望 0 条事件，实际=%d", len(events))
	}
}

// TestHealthRecordService_AddTimelineEvent_EmptyPatientID 患者ID为空时应返回参数无效错误
func TestHealthRecordService_AddTimelineEvent_EmptyPatientID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	event := &model.TimelineEvent{
		PatientID:   "",
		EventType:   "visit",
		Description: "测试事件",
	}
	err := svc.AddTimelineEvent(event)
	if err == nil {
		t.Error("期望患者ID为空时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "患者ID不能为空" {
		t.Errorf("期望错误详情='患者ID不能为空'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_AddTimelineEvent_EmptyEventType 事件类型为空时应返回参数无效错误
func TestHealthRecordService_AddTimelineEvent_EmptyEventType(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	event := &model.TimelineEvent{
		PatientID:   "P001",
		EventType:   "",
		Description: "测试事件",
	}
	err := svc.AddTimelineEvent(event)
	if err == nil {
		t.Error("期望事件类型为空时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "事件类型不能为空" {
		t.Errorf("期望错误详情='事件类型不能为空'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_AddTimelineEvent_InvalidEventType 无效事件类型时应返回参数无效错误
func TestHealthRecordService_AddTimelineEvent_InvalidEventType(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	event := &model.TimelineEvent{
		PatientID:   "P001",
		EventType:   "invalid_type",
		Description: "测试事件",
	}
	err := svc.AddTimelineEvent(event)
	if err == nil {
		t.Error("期望无效事件类型时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "无效的事件类型，合法值：visit/prescription/examination/followup" {
		t.Errorf("期望错误详情='无效的事件类型，合法值：visit/prescription/examination/followup'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_AddTimelineEvent_Success 正常添加时间轴事件
func TestHealthRecordService_AddTimelineEvent_Success(t *testing.T) {
	svc, db := newTestHealthRecordService(t)

	event := &model.TimelineEvent{
		PatientID:   "P001",
		EventType:   "visit",
		Description: "门诊就诊",
		Date:        "2024-01-15",
	}
	err := svc.AddTimelineEvent(event)
	if err != nil {
		t.Fatalf("添加时间轴事件失败: %v", err)
	}

	var saved model.TimelineEvent
	if err := db.Where("patient_id = ?", "P001").First(&saved).Error; err != nil {
		t.Fatalf("查询已保存事件失败: %v", err)
	}
	if saved.EventType != "visit" {
		t.Errorf("期望 EventType='visit'，实际='%s'", saved.EventType)
	}
	if saved.Description != "门诊就诊" {
		t.Errorf("期望 Description='门诊就诊'，实际='%s'", saved.Description)
	}
}

// TestHealthRecordService_GrantAuthorization_EmptyPatientID 患者ID为空时应返回参数无效错误
func TestHealthRecordService_GrantAuthorization_EmptyPatientID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	auth := &model.RecordAuthorization{
		PatientID: "",
		DoctorID:  "D001",
	}
	err := svc.GrantAuthorization(auth)
	if err == nil {
		t.Error("期望患者ID为空时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "患者ID不能为空" {
		t.Errorf("期望错误详情='患者ID不能为空'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_GrantAuthorization_EmptyDoctorID 医生ID为空时应返回参数无效错误
func TestHealthRecordService_GrantAuthorization_EmptyDoctorID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	auth := &model.RecordAuthorization{
		PatientID: "P001",
		DoctorID:  "",
	}
	err := svc.GrantAuthorization(auth)
	if err == nil {
		t.Error("期望医生ID为空时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "医生ID不能为空" {
		t.Errorf("期望错误详情='医生ID不能为空'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_GrantAuthorization_Duplicate 重复授权时应返回重复授权错误
func TestHealthRecordService_GrantAuthorization_Duplicate(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	auth := &model.RecordAuthorization{
		PatientID: "P001",
		DoctorID:  "D001",
	}
	err := svc.GrantAuthorization(auth)
	if err != nil {
		t.Fatalf("首次授权失败: %v", err)
	}

	auth2 := &model.RecordAuthorization{
		PatientID: "P001",
		DoctorID:  "D001",
	}
	err = svc.GrantAuthorization(auth2)
	if err == nil {
		t.Error("期望重复授权时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeAuthDuplicate {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeAuthDuplicate, appErr.Code)
	}
	if appErr.Detail != "已存在有效授权" {
		t.Errorf("期望错误详情='已存在有效授权'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_GrantAuthorization_Success 正常授权成功
func TestHealthRecordService_GrantAuthorization_Success(t *testing.T) {
	svc, db := newTestHealthRecordService(t)

	auth := &model.RecordAuthorization{
		PatientID: "P001",
		DoctorID:  "D001",
	}
	err := svc.GrantAuthorization(auth)
	if err != nil {
		t.Fatalf("授权失败: %v", err)
	}
	if auth.Status != 1 {
		t.Errorf("期望 Status=1，实际=%d", auth.Status)
	}
	if auth.AuthTime == "" {
		t.Error("期望 AuthTime 不为空")
	}
	if auth.ExpireTime == "" {
		t.Error("期望 ExpireTime 不为空")
	}

	var saved model.RecordAuthorization
	if err := db.Where("patient_id = ? AND doctor_id = ?", "P001", "D001").First(&saved).Error; err != nil {
		t.Fatalf("查询已保存授权失败: %v", err)
	}
	if saved.Status != 1 {
		t.Errorf("期望保存的 Status=1，实际=%d", saved.Status)
	}
}

// TestHealthRecordService_RevokeAuthorization_EmptyPatientID 患者ID为空时应返回参数无效错误
func TestHealthRecordService_RevokeAuthorization_EmptyPatientID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	err := svc.RevokeAuthorization("", "D001")
	if err == nil {
		t.Error("期望患者ID为空时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "患者ID和医生ID不能为空" {
		t.Errorf("期望错误详情='患者ID和医生ID不能为空'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_RevokeAuthorization_EmptyDoctorID 医生ID为空时应返回参数无效错误
func TestHealthRecordService_RevokeAuthorization_EmptyDoctorID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	err := svc.RevokeAuthorization("P001", "")
	if err == nil {
		t.Error("期望医生ID为空时返回错误")
	}
	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Fatalf("期望返回 *AppError，实际类型=%T", err)
	}
	if appErr.Code != errors.CodeParamInvalid {
		t.Errorf("期望错误码=%d，实际=%d", errors.CodeParamInvalid, appErr.Code)
	}
	if appErr.Detail != "患者ID和医生ID不能为空" {
		t.Errorf("期望错误详情='患者ID和医生ID不能为空'，实际='%s'", appErr.Detail)
	}
}

// TestHealthRecordService_RevokeAuthorization_Success 正常撤销授权
func TestHealthRecordService_RevokeAuthorization_Success(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	auth := &model.RecordAuthorization{
		PatientID: "P001",
		DoctorID:  "D001",
	}
	if err := svc.GrantAuthorization(auth); err != nil {
		t.Fatalf("授权失败: %v", err)
	}

	if !svc.CheckAuthorization("P001", "D001") {
		t.Fatal("撤销前期望授权有效")
	}

	err := svc.RevokeAuthorization("P001", "D001")
	if err != nil {
		t.Fatalf("撤销授权失败: %v", err)
	}

	if svc.CheckAuthorization("P001", "D001") {
		t.Error("撤销后期望授权无效")
	}
}

// TestHealthRecordService_CheckAuthorization_EmptyPatientID 患者ID为空时应返回 false
func TestHealthRecordService_CheckAuthorization_EmptyPatientID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	authorized := svc.CheckAuthorization("", "D001")
	if authorized {
		t.Error("期望患者ID为空时返回 false")
	}
}

// TestHealthRecordService_CheckAuthorization_EmptyDoctorID 医生ID为空时应返回 false
func TestHealthRecordService_CheckAuthorization_EmptyDoctorID(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	authorized := svc.CheckAuthorization("P001", "")
	if authorized {
		t.Error("期望医生ID为空时返回 false")
	}
}

// TestHealthRecordService_CheckAuthorization_Success 正常检查授权返回 true
func TestHealthRecordService_CheckAuthorization_Success(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	auth := &model.RecordAuthorization{
		PatientID: "P001",
		DoctorID:  "D001",
	}
	if err := svc.GrantAuthorization(auth); err != nil {
		t.Fatalf("授权失败: %v", err)
	}

	if !svc.CheckAuthorization("P001", "D001") {
		t.Error("期望授权有效时返回 true")
	}
}

// TestHealthRecordService_CheckAuthorization_NotAuthorized 未授权时应返回 false
func TestHealthRecordService_CheckAuthorization_NotAuthorized(t *testing.T) {
	svc, _ := newTestHealthRecordService(t)

	authorized := svc.CheckAuthorization("P999", "D999")
	if authorized {
		t.Error("期望未授权时返回 false")
	}
}
