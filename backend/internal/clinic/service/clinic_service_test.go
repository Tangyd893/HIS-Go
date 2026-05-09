package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/clinic/model"
	"his-go/internal/clinic/repository"
)

func setupClinicService(t *testing.T) (*ClinicService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS clinic_records (
			id TEXT PRIMARY KEY,
			registration_id TEXT,
			patient_id TEXT NOT NULL,
			patient_name TEXT,
			doctor_id TEXT NOT NULL,
			chief_complaint TEXT,
			present_illness TEXT,
			diagnosis TEXT,
			icd_code TEXT,
			advice TEXT,
			status INTEGER DEFAULT 0,
			visit_time DATETIME,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS examination_requests (
			id TEXT PRIMARY KEY,
			clinic_record_id TEXT NOT NULL,
			patient_id TEXT NOT NULL,
			exam_type TEXT,
			exam_item TEXT,
			body_part TEXT,
			clinical_diagnosis TEXT,
			note TEXT,
			status INTEGER DEFAULT 0,
			created_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	repo := repository.NewClinicRepository(db)
	svc := NewClinicService(repo)

	return svc, db
}

func TestClinicService_CreateRecord(t *testing.T) {
	svc, db := setupClinicService(t)

	record := &model.ClinicRecord{
		ID:             "clinic-001",
		RegistrationID: "reg-1",
		PatientID:      "patient-1",
		PatientName:    "张三",
		DoctorID:       "doctor-1",
		ChiefComplaint: "发热3天",
		PresentIllness: "患者3天前无明显诱因发热",
		Diagnosis:      "上呼吸道感染",
		IcdCode:        "J06.9",
		Advice:         "多饮水，注意休息",
		Status:         0,
	}

	err := svc.CreateRecord(record)
	if err != nil {
		t.Fatalf("创建诊疗记录失败: %v", err)
	}

	var saved model.ClinicRecord
	if err := db.Where("id = ?", "clinic-001").First(&saved).Error; err != nil {
		t.Fatalf("查询诊疗记录失败: %v", err)
	}
	if saved.ChiefComplaint != "发热3天" {
		t.Errorf("期望主诉='发热3天'，实际=%s", saved.ChiefComplaint)
	}
	if saved.Diagnosis != "上呼吸道感染" {
		t.Errorf("期望诊断='上呼吸道感染'，实际=%s", saved.Diagnosis)
	}
}

func TestClinicService_GetByID(t *testing.T) {
	svc, db := setupClinicService(t)

	record := &model.ClinicRecord{
		ID: "clinic-002", PatientID: "patient-2", DoctorID: "doctor-2",
		ChiefComplaint: "腹痛", Diagnosis: "急性胃炎", IcdCode: "K29.1",
	}
	if err := db.Create(record).Error; err != nil {
		t.Fatalf("创建诊疗记录失败: %v", err)
	}

	result, err := svc.GetByID("clinic-002")
	if err != nil {
		t.Fatalf("获取诊疗记录失败: %v", err)
	}
	if result.Diagnosis != "急性胃炎" {
		t.Errorf("期望诊断='急性胃炎'，实际=%s", result.Diagnosis)
	}
}

func TestClinicService_ListByPatient(t *testing.T) {
	svc, db := setupClinicService(t)

	records := []model.ClinicRecord{
		{ID: "c-10", PatientID: "patient-x", DoctorID: "doc-a", ChiefComplaint: "头痛"},
		{ID: "c-11", PatientID: "patient-x", DoctorID: "doc-b", ChiefComplaint: "咳嗽"},
		{ID: "c-12", PatientID: "patient-y", DoctorID: "doc-a", ChiefComplaint: "发热"},
	}
	for _, r := range records {
		if err := db.Create(&r).Error; err != nil {
			t.Fatalf("创建诊疗记录失败: %v", err)
		}
	}

	list, total, err := svc.ListByPatient("patient-x", 1, 10)
	if err != nil {
		t.Fatalf("查询诊疗记录失败: %v", err)
	}
	if total != 2 {
		t.Errorf("期望total=2，实际=%d", total)
	}
	if len(list) != 2 {
		t.Errorf("期望len(list)=2，实际=%d", len(list))
	}
}

func TestClinicService_ListByDoctor(t *testing.T) {
	svc, db := setupClinicService(t)

	records := []model.ClinicRecord{
		{ID: "c-20", PatientID: "p-a", DoctorID: "doc-10", ChiefComplaint: "鼻炎"},
		{ID: "c-21", PatientID: "p-b", DoctorID: "doc-10", ChiefComplaint: "咽炎"},
	}
	for _, r := range records {
		if err := db.Create(&r).Error; err != nil {
			t.Fatalf("创建诊疗记录失败: %v", err)
		}
	}

	list, total, err := svc.ListByDoctor("doc-10", 1, 10)
	if err != nil {
		t.Fatalf("查询诊疗记录失败: %v", err)
	}
	if total != 2 {
		t.Errorf("期望total=2，实际=%d", total)
	}
	if len(list) != 2 {
		t.Errorf("期望len(list)=2，实际=%d", len(list))
	}
}

func TestClinicService_CreateExamRequest(t *testing.T) {
	svc, db := setupClinicService(t)

	req := &model.ExaminationRequest{
		ID:                "exam-001",
		ClinicRecordID:    "clinic-001",
		PatientID:         "patient-1",
		ExamType:          "影像学检查",
		ExamItem:          "胸部X光",
		BodyPart:          "胸部",
		ClinicalDiagnosis: "疑似肺炎",
		Status:            0,
	}

	err := svc.CreateExamRequest(req)
	if err != nil {
		t.Fatalf("创建检查申请失败: %v", err)
	}

	var saved model.ExaminationRequest
	if err := db.Where("id = ?", "exam-001").First(&saved).Error; err != nil {
		t.Fatalf("查询检查申请失败: %v", err)
	}
	if saved.ExamItem != "胸部X光" {
		t.Errorf("期望检查项目='胸部X光'，实际=%s", saved.ExamItem)
	}
}

func TestClinicService_ListExamRequests(t *testing.T) {
	svc, db := setupClinicService(t)

	reqs := []model.ExaminationRequest{
		{ID: "exam-a", ClinicRecordID: "c-1", PatientID: "patient-z", ExamItem: "血常规"},
		{ID: "exam-b", ClinicRecordID: "c-2", PatientID: "patient-z", ExamItem: "尿常规"},
	}
	for _, r := range reqs {
		if err := db.Create(&r).Error; err != nil {
			t.Fatalf("创建检查申请失败: %v", err)
		}
	}

	list, err := svc.ListExamRequests("patient-z")
	if err != nil {
		t.Fatalf("查询检查申请失败: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("期望2条检查申请，实际=%d", len(list))
	}
}
