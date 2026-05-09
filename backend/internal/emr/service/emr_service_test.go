package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/emr/model"
	"his-go/internal/emr/repository"
)

func setupEMRService(t *testing.T) (*EMRService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS medical_records (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			clinic_record_id TEXT,
			template_id TEXT,
			chief_complaint TEXT,
			present_illness TEXT,
			past_history TEXT,
			physical_exam TEXT,
			auxiliary_exam TEXT,
			diagnosis TEXT,
			treatment_plan TEXT,
			quality_level INTEGER DEFAULT 0,
			status INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	repo := repository.NewEMRRepository(db)
	return NewEMRService(repo), db
}

// ==================== validateRecord 测试 ====================

func TestValidateRecord_AllValid(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "patient-1",
		ChiefComplaint: "头痛3天",
		PresentIllness: "患者3天前无明显诱因出现头痛",
		PastHistory:    "否认高血压、糖尿病史",
		PhysicalExam:   "T36.5℃ P72次/分 R18次/分 BP120/80mmHg",
		AuxiliaryExam:  "头颅CT未见明显异常",
		Diagnosis:      "紧张性头痛",
		TreatmentPlan:  "1.注意休息 2.布洛芬200mg tid",
	}

	err := svc.validateRecord(record)
	if err != nil {
		t.Errorf("期望所有字段有效时通过，错误: %v", err)
	}
}

func TestValidateRecord_EmptyPatientID(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "",
		ChiefComplaint: "头痛",
		PresentIllness: "有",
		PastHistory:    "有",
		PhysicalExam:   "有",
		AuxiliaryExam:  "有",
		Diagnosis:      "有",
		TreatmentPlan:  "有",
	}

	err := svc.validateRecord(record)
	if err == nil {
		t.Error("期望空PatientID返回错误")
	}
}

func TestValidateRecord_EmptyChiefComplaint(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "p1",
		ChiefComplaint: "",
		PresentIllness: "有",
		PastHistory:    "有",
		PhysicalExam:   "有",
		AuxiliaryExam:  "有",
		Diagnosis:      "有",
		TreatmentPlan:  "有",
	}

	err := svc.validateRecord(record)
	if err == nil {
		t.Error("期望空主诉返回错误")
	}
}

func TestValidateRecord_EmptyPresentIllness(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "p1",
		ChiefComplaint: "有",
		PresentIllness: "",
		PastHistory:    "有",
		PhysicalExam:   "有",
		AuxiliaryExam:  "有",
		Diagnosis:      "有",
		TreatmentPlan:  "有",
	}

	err := svc.validateRecord(record)
	if err == nil {
		t.Error("期望空现病史返回错误")
	}
}

func TestValidateRecord_EmptyPastHistory(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "p1",
		ChiefComplaint: "有",
		PresentIllness: "有",
		PastHistory:    "",
		PhysicalExam:   "有",
		AuxiliaryExam:  "有",
		Diagnosis:      "有",
		TreatmentPlan:  "有",
	}

	err := svc.validateRecord(record)
	if err == nil {
		t.Error("期望空既往史返回错误")
	}
}

func TestValidateRecord_EmptyPhysicalExam(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "p1",
		ChiefComplaint: "有",
		PresentIllness: "有",
		PastHistory:    "有",
		PhysicalExam:   "",
		AuxiliaryExam:  "有",
		Diagnosis:      "有",
		TreatmentPlan:  "有",
	}

	err := svc.validateRecord(record)
	if err == nil {
		t.Error("期望空体格检查返回错误")
	}
}

func TestValidateRecord_EmptyAuxiliaryExam(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "p1",
		ChiefComplaint: "有",
		PresentIllness: "有",
		PastHistory:    "有",
		PhysicalExam:   "有",
		AuxiliaryExam:  "",
		Diagnosis:      "有",
		TreatmentPlan:  "有",
	}

	err := svc.validateRecord(record)
	if err == nil {
		t.Error("期望空辅助检查返回错误")
	}
}

func TestValidateRecord_EmptyDiagnosis(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "p1",
		ChiefComplaint: "有",
		PresentIllness: "有",
		PastHistory:    "有",
		PhysicalExam:   "有",
		AuxiliaryExam:  "有",
		Diagnosis:      "",
		TreatmentPlan:  "有",
	}

	err := svc.validateRecord(record)
	if err == nil {
		t.Error("期望空诊断返回错误")
	}
}

func TestValidateRecord_EmptyTreatmentPlan(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "p1",
		ChiefComplaint: "有",
		PresentIllness: "有",
		PastHistory:    "有",
		PhysicalExam:   "有",
		AuxiliaryExam:  "有",
		Diagnosis:      "有",
		TreatmentPlan:  "",
	}

	err := svc.validateRecord(record)
	if err == nil {
		t.Error("期望空处理计划返回错误")
	}
}

// ==================== CreateRecord 测试 ====================

func TestCreateRecord_Valid(t *testing.T) {
	svc, db := setupEMRService(t)

	record := &model.MedicalRecord{
		ID:             "mr-1",
		PatientID:      "patient-1",
		ChiefComplaint: "头痛",
		PresentIllness: "有",
		PastHistory:    "有",
		PhysicalExam:   "有",
		AuxiliaryExam:  "有",
		Diagnosis:      "有",
		TreatmentPlan:  "有",
	}

	err := svc.CreateRecord(record)
	if err != nil {
		t.Fatalf("创建有效病历失败: %v", err)
	}

	var saved model.MedicalRecord
	db.Where("id = ?", "mr-1").First(&saved)
	if saved.ID != "mr-1" {
		t.Error("病历未保存到数据库")
	}
}

func TestCreateRecord_Invalid(t *testing.T) {
	svc, _ := setupEMRService(t)

	record := &model.MedicalRecord{
		PatientID:      "",
		ChiefComplaint: "头痛",
	}

	err := svc.CreateRecord(record)
	if err == nil {
		t.Error("期望无效病历返回错误")
	}
}

// ==================== UpdateRecord 测试 ====================

func TestUpdateRecord_Valid(t *testing.T) {
	svc, db := setupEMRService(t)

	existing := &model.MedicalRecord{
		ID:             "mr-2",
		PatientID:      "patient-1",
		ChiefComplaint: "旧主诉",
		PresentIllness: "旧现病史",
		PastHistory:    "旧既往史",
		PhysicalExam:   "旧体格检查",
		AuxiliaryExam:  "旧辅助检查",
		Diagnosis:      "旧诊断",
		TreatmentPlan:  "旧处理计划",
	}
	if err := db.Create(existing).Error; err != nil {
		t.Fatalf("创建现有记录失败: %v", err)
	}

	record := &model.MedicalRecord{
		ID:             "mr-2",
		PatientID:      "patient-1",
		ChiefComplaint: "更新后主诉",
		PresentIllness: "更新后现病史",
		PastHistory:    "更新后既往史",
		PhysicalExam:   "更新后体格检查",
		AuxiliaryExam:  "更新后辅助检查",
		Diagnosis:      "更新后诊断",
		TreatmentPlan:  "更新后处理计划",
	}

	err := svc.UpdateRecord(record)
	if err != nil {
		t.Fatalf("更新有效病历失败: %v", err)
	}

	var updated model.MedicalRecord
	db.Where("id = ?", "mr-2").First(&updated)
	if updated.ChiefComplaint != "更新后主诉" {
		t.Errorf("期望更新后主诉='更新后主诉'，实际=%s", updated.ChiefComplaint)
	}
}

func TestUpdateRecord_Invalid(t *testing.T) {
	svc, db := setupEMRService(t)

	existing := &model.MedicalRecord{
		ID:             "mr-3",
		PatientID:      "patient-1",
		ChiefComplaint: "旧主诉",
		PresentIllness: "有",
		PastHistory:    "有",
		PhysicalExam:   "有",
		AuxiliaryExam:  "有",
		Diagnosis:      "有",
		TreatmentPlan:  "有",
	}
	if err := db.Create(existing).Error; err != nil {
		t.Fatalf("创建现有记录失败: %v", err)
	}

	record := &model.MedicalRecord{
		ID:        "mr-3",
		Diagnosis: "缺少必填字段",
	}

	err := svc.UpdateRecord(record)
	if err == nil {
		t.Error("期望校验失败时返回错误")
	}
}
