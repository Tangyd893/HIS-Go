package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/prescription/model"
	"his-go/internal/prescription/repository"
)

func setupPrescriptionService(t *testing.T) (*PrescriptionService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS prescriptions (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			patient_name TEXT,
			doctor_id TEXT NOT NULL,
			diagnosis_id TEXT,
			prescription_type INTEGER NOT NULL,
			status INTEGER DEFAULT 0,
			note TEXT,
			version INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS prescription_details (
			id TEXT PRIMARY KEY,
			prescription_id TEXT NOT NULL,
			drug_id TEXT NOT NULL,
			drug_name TEXT,
			specification TEXT,
			dosage REAL DEFAULT 0,
			usage TEXT,
			frequency TEXT,
			days INTEGER DEFAULT 1,
			quantity INTEGER DEFAULT 0,
			unit_price REAL DEFAULT 0,
			note TEXT
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	repo := repository.NewPrescriptionRepository(db)
	svc := NewPrescriptionService(repo)

	return svc, db
}

func TestPrescriptionService_Create(t *testing.T) {
	svc, db := setupPrescriptionService(t)

	prescription := &model.Prescription{
		ID:               "pres-001",
		PatientID:        "patient-1",
		PatientName:      "张三",
		DoctorID:         "doctor-1",
		DiagnosisID:      "diag-1",
		PrescriptionType: 1,
	}
	details := []model.PrescriptionDetail{
		{ID: "det-1", DrugID: "drug-1", DrugName: "阿莫西林胶囊", Specification: "0.5g×24粒", Dosage: 0.5, Usage: "口服", Frequency: "tid", Days: 7, Quantity: 21, UnitPrice: 25.50},
		{ID: "det-2", DrugID: "drug-2", DrugName: "布洛芬缓释胶囊", Specification: "0.3g×20粒", Dosage: 0.3, Usage: "口服", Frequency: "bid", Days: 3, Quantity: 6, UnitPrice: 18.00},
	}

	err := svc.Create(prescription, details)
	if err != nil {
		t.Fatalf("创建处方失败: %v", err)
	}

	var saved model.Prescription
	if err := db.Where("id = ?", "pres-001").First(&saved).Error; err != nil {
		t.Fatalf("查询处方失败: %v", err)
	}
	if saved.Status != 0 {
		t.Errorf("期望Status=0(草稿)，实际=%d", saved.Status)
	}
	if saved.PatientName != "张三" {
		t.Errorf("期望PatientName='张三'，实际=%s", saved.PatientName)
	}

	var savedDetails []model.PrescriptionDetail
	if err := db.Where("prescription_id = ?", "pres-001").Find(&savedDetails).Error; err != nil {
		t.Fatalf("查询处方明细失败: %v", err)
	}
	if len(savedDetails) != 2 {
		t.Errorf("期望2条明细，实际=%d", len(savedDetails))
	}
}

func TestPrescriptionService_GetByID(t *testing.T) {
	svc, db := setupPrescriptionService(t)

	prescription := &model.Prescription{
		ID: "pres-002", PatientID: "patient-2", DoctorID: "doctor-2",
		PrescriptionType: 2, Status: 0,
	}
	if err := db.Create(prescription).Error; err != nil {
		t.Fatalf("创建处方失败: %v", err)
	}

	result, err := svc.GetByID("pres-002")
	if err != nil {
		t.Fatalf("获取处方失败: %v", err)
	}
	if result.PatientID != "patient-2" {
		t.Errorf("期望PatientID='patient-2'，实际=%s", result.PatientID)
	}
}

func TestPrescriptionService_ListByPatient(t *testing.T) {
	svc, db := setupPrescriptionService(t)

	prescriptions := []model.Prescription{
		{ID: "pr-1", PatientID: "patient-x", DoctorID: "doc-a", PrescriptionType: 1, Status: 0},
		{ID: "pr-2", PatientID: "patient-x", DoctorID: "doc-b", PrescriptionType: 1, Status: 2},
		{ID: "pr-3", PatientID: "patient-y", DoctorID: "doc-a", PrescriptionType: 2, Status: 0},
	}
	for _, p := range prescriptions {
		if err := db.Create(&p).Error; err != nil {
			t.Fatalf("创建处方失败: %v", err)
		}
	}

	list, total, err := svc.ListByPatient("patient-x", 1, 10)
	if err != nil {
		t.Fatalf("查询处方列表失败: %v", err)
	}
	if total != 2 {
		t.Errorf("期望total=2，实际=%d", total)
	}
	if len(list) != 2 {
		t.Errorf("期望len(list)=2，实际=%d", len(list))
	}
}

func TestPrescriptionService_SubmitForReview(t *testing.T) {
	svc, db := setupPrescriptionService(t)

	prescription := &model.Prescription{
		ID: "pres-review", PatientID: "patient-3", DoctorID: "doctor-3",
		PrescriptionType: 1, Status: 0,
	}
	if err := db.Create(prescription).Error; err != nil {
		t.Fatalf("创建处方失败: %v", err)
	}

	err := svc.SubmitForReview("pres-review")
	if err != nil {
		t.Fatalf("提交审核失败: %v", err)
	}

	var updated model.Prescription
	db.Where("id = ?", "pres-review").First(&updated)
	if updated.Status != 1 {
		t.Errorf("期望Status=1(待审核)，实际=%d", updated.Status)
	}
}

func TestPrescriptionService_Review_Approve(t *testing.T) {
	svc, db := setupPrescriptionService(t)

	prescription := &model.Prescription{
		ID: "pres-approve", PatientID: "patient-4", DoctorID: "doctor-4",
		PrescriptionType: 1, Status: 1,
	}
	if err := db.Create(prescription).Error; err != nil {
		t.Fatalf("创建处方失败: %v", err)
	}

	err := svc.Review("pres-approve", true, "审核通过，用药合理")
	if err != nil {
		t.Fatalf("审核处方失败: %v", err)
	}

	var updated model.Prescription
	db.Where("id = ?", "pres-approve").First(&updated)
	if updated.Status != 2 {
		t.Errorf("期望Status=2(已审核)，实际=%d", updated.Status)
	}
	if updated.Note != "审核通过，用药合理" {
		t.Errorf("期望Note='审核通过，用药合理'，实际=%s", updated.Note)
	}
}

func TestPrescriptionService_Review_Reject(t *testing.T) {
	svc, db := setupPrescriptionService(t)

	prescription := &model.Prescription{
		ID: "pres-reject", PatientID: "patient-5", DoctorID: "doctor-5",
		PrescriptionType: 1, Status: 1,
	}
	if err := db.Create(prescription).Error; err != nil {
		t.Fatalf("创建处方失败: %v", err)
	}

	err := svc.Review("pres-reject", false, "药品剂量过大，请修改")
	if err != nil {
		t.Fatalf("审核处方失败: %v", err)
	}

	var updated model.Prescription
	db.Where("id = ?", "pres-reject").First(&updated)
	if updated.Status != 0 {
		t.Errorf("期望Status=0(退回草稿)，实际=%d", updated.Status)
	}
	if updated.Note != "药品剂量过大，请修改" {
		t.Errorf("期望Note='药品剂量过大，请修改'，实际=%s", updated.Note)
	}
}

func TestPrescriptionService_Cancel(t *testing.T) {
	svc, db := setupPrescriptionService(t)

	prescription := &model.Prescription{
		ID: "pres-cancel", PatientID: "patient-6", DoctorID: "doctor-6",
		PrescriptionType: 1, Status: 2,
	}
	if err := db.Create(prescription).Error; err != nil {
		t.Fatalf("创建处方失败: %v", err)
	}

	err := svc.Cancel("pres-cancel")
	if err != nil {
		t.Fatalf("取消处方失败: %v", err)
	}

	var updated model.Prescription
	db.Where("id = ?", "pres-cancel").First(&updated)
	if updated.Status != 3 {
		t.Errorf("期望Status=3(已取消)，实际=%d", updated.Status)
	}
}

func TestPrescriptionService_SubmitForReview_InvalidStatus(t *testing.T) {
	svc, db := setupPrescriptionService(t)

	prescription := &model.Prescription{
		ID: "pres-bad", PatientID: "patient-7", DoctorID: "doctor-7",
		PrescriptionType: 1, Status: 2,
	}
	if err := db.Create(prescription).Error; err != nil {
		t.Fatalf("创建处方失败: %v", err)
	}

	err := svc.SubmitForReview("pres-bad")
	if err == nil {
		t.Error("期望非草稿状态提交审核时返回错误")
	}
}
