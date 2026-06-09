package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/inpatient/model"
	"his-go/internal/inpatient/repository"
)

func setupInpatientService(t *testing.T) (*InpatientService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS inpatient_records (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			patient_name TEXT NOT NULL,
			admission_date DATETIME,
			discharge_date DATETIME,
			dept_id TEXT NOT NULL,
			room_no TEXT,
			bed_no TEXT,
			diagnosis TEXT,
			deposit REAL NOT NULL DEFAULT 0,
			total_cost REAL NOT NULL DEFAULT 0,
			status INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS medical_orders (
			id TEXT PRIMARY KEY,
			inpatient_id TEXT NOT NULL,
			doctor_id TEXT NOT NULL DEFAULT '',
			order_type INTEGER NOT NULL,
			content TEXT NOT NULL,
			start_time DATETIME,
			end_time DATETIME,
			status INTEGER DEFAULT 1,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS nursing_records (
			id TEXT PRIMARY KEY,
			inpatient_id TEXT NOT NULL,
			nurse_id TEXT NOT NULL,
			record_time DATETIME,
			content TEXT NOT NULL,
			vital_signs TEXT DEFAULT '',
			created_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v", err)
		}
	}

	repo := repository.NewInpatientRepository(db)
	svc := NewInpatientService(repo)

	return svc, db
}

func TestInpatientService_AdmitPatient(t *testing.T) {
	svc, db := setupInpatientService(t)

	record := &model.InpatientRecord{
		ID:          "ip-001",
		PatientID:   "patient-1",
		PatientName: "王小明",
		DeptID:      "dept_001",
		BedNo:       "A101",
		Deposit:     5000,
	}

	err := svc.AdmitPatient(record)
	if err != nil {
		t.Fatalf("入院登记失败: %v", err)
	}

	var saved model.InpatientRecord
	if err := db.Where("id = ?", "ip-001").First(&saved).Error; err != nil {
		t.Fatalf("查询住院记录失败: %v", err)
	}
	if saved.Status != 1 {
		t.Errorf("期望Status=1(在院)，实际=%d", saved.Status)
	}
	if saved.PatientName != "王小明" {
		t.Errorf("期望PatientName='王小明'，实际=%s", saved.PatientName)
	}
}

func TestInpatientService_AdmitPatient_Nil(t *testing.T) {
	svc, _ := setupInpatientService(t)

	err := svc.AdmitPatient(nil)
	if err == nil {
		t.Error("期望 nil 记录时返回错误")
	}
}

func TestInpatientService_AdmitPatient_EmptyPatientID(t *testing.T) {
	svc, _ := setupInpatientService(t)

	record := &model.InpatientRecord{
		ID:          "ip-002",
		PatientName: "测试",
		DeptID:      "dept_001",
	}

	err := svc.AdmitPatient(record)
	if err == nil {
		t.Error("期望空患者ID时返回错误")
	}
}

func TestInpatientService_AdmitPatient_EmptyPatientName(t *testing.T) {
	svc, _ := setupInpatientService(t)

	record := &model.InpatientRecord{
		ID:        "ip-003",
		PatientID: "patient-3",
		DeptID:    "dept_001",
	}

	err := svc.AdmitPatient(record)
	if err == nil {
		t.Error("期望空患者姓名时返回错误")
	}
}

func TestInpatientService_AdmitPatient_EmptyDeptID(t *testing.T) {
	svc, _ := setupInpatientService(t)

	record := &model.InpatientRecord{
		ID:          "ip-004",
		PatientID:   "patient-4",
		PatientName: "测试",
	}

	err := svc.AdmitPatient(record)
	if err == nil {
		t.Error("期望空科室ID时返回错误")
	}
}

func TestInpatientService_AdmitPatient_NegativeDeposit(t *testing.T) {
	svc, _ := setupInpatientService(t)

	record := &model.InpatientRecord{
		ID:          "ip-005",
		PatientID:   "patient-5",
		PatientName: "测试",
		DeptID:      "dept_001",
		Deposit:     -100,
	}

	err := svc.AdmitPatient(record)
	if err == nil {
		t.Error("期望押金为负数时返回错误")
	}
}

func TestInpatientService_AdmitPatient_DuplicateAdmit(t *testing.T) {
	svc, db := setupInpatientService(t)

	// 先创建一个在院的记录
	db.Exec(`INSERT INTO inpatient_records (id, patient_id, patient_name, dept_id, bed_no, status) 
		VALUES ('ip-existing', 'patient-dup', '重复患者', 'dept_001', 'B01', 1)`)

	record := &model.InpatientRecord{
		ID:          "ip-006",
		PatientID:   "patient-dup",
		PatientName: "重复患者",
		DeptID:      "dept_002",
	}

	err := svc.AdmitPatient(record)
	if err == nil {
		t.Error("期望重复入院时返回错误")
	}
}

func TestInpatientService_DischargePatient(t *testing.T) {
	svc, db := setupInpatientService(t)

	db.Exec(`INSERT INTO inpatient_records (id, patient_id, patient_name, dept_id, bed_no, status) 
		VALUES ('ip-dc', 'patient-dc', '出院患者', 'dept_001', 'C01', 1)`)

	err := svc.DischargePatient("ip-dc")
	if err != nil {
		t.Fatalf("出院失败: %v", err)
	}

	var saved model.InpatientRecord
	db.Where("id = ?", "ip-dc").First(&saved)
	if saved.Status != 2 {
		t.Errorf("期望Status=2(已出院)，实际=%d", saved.Status)
	}
}

func TestInpatientService_DischargePatient_EmptyID(t *testing.T) {
	svc, _ := setupInpatientService(t)

	err := svc.DischargePatient("")
	if err == nil {
		t.Error("期望空ID时返回错误")
	}
}

func TestInpatientService_DischargePatient_NotAdmitted(t *testing.T) {
	svc, db := setupInpatientService(t)

	db.Exec(`INSERT INTO inpatient_records (id, patient_id, patient_name, dept_id, bed_no, status) 
		VALUES ('ip-done', 'patient-done', '已出院', 'dept_001', 'D01', 2)`)

	err := svc.DischargePatient("ip-done")
	if err == nil {
		t.Error("期望已出院患者再次出院时返回错误")
	}
}

func TestInpatientService_CreateMedicalOrder(t *testing.T) {
	svc, db := setupInpatientService(t)

	db.Exec(`INSERT INTO inpatient_records (id, patient_id, patient_name, dept_id, bed_no, status) 
		VALUES ('ip-ord', 'patient-ord', '医嘱患者', 'dept_001', 'E01', 1)`)

	order := &model.MedicalOrder{
		ID:          "mo-001",
		InpatientID: "ip-ord",
		DoctorID:    "doctor-1",
		OrderType:   1,
		Content:     "青霉素 80万U 肌注 bid",
	}

	err := svc.CreateMedicalOrder(order)
	if err != nil {
		t.Fatalf("创建医嘱失败: %v", err)
	}
}

func TestInpatientService_CreateMedicalOrder_Nil(t *testing.T) {
	svc, _ := setupInpatientService(t)

	err := svc.CreateMedicalOrder(nil)
	if err == nil {
		t.Error("期望 nil 医嘱时返回错误")
	}
}

func TestInpatientService_CreateMedicalOrder_InvalidType(t *testing.T) {
	svc, db := setupInpatientService(t)

	db.Exec(`INSERT INTO inpatient_records (id, patient_id, patient_name, dept_id, bed_no, status) 
		VALUES ('ip-type', 'patient-type', '类型测试', 'dept_001', 'F01', 1)`)

	order := &model.MedicalOrder{
		ID:          "mo-002",
		InpatientID: "ip-type",
		OrderType:   99,
		Content:     "测试",
	}

	err := svc.CreateMedicalOrder(order)
	if err == nil {
		t.Error("期望无效医嘱类型时返回错误")
	}
}

func TestInpatientService_CreateNursingRecord(t *testing.T) {
	svc, db := setupInpatientService(t)

	db.Exec(`INSERT INTO inpatient_records (id, patient_id, patient_name, dept_id, bed_no, status) 
		VALUES ('ip-nurse', 'patient-nurse', '护理患者', 'dept_001', 'G01', 1)`)

	record := &model.NursingRecord{
		ID:          "nr-001",
		InpatientID: "ip-nurse",
		NurseID:     "nurse-1",
		Content:     "体温36.5℃，血压120/80，精神状态良好",
	}

	err := svc.CreateNursingRecord(record)
	if err != nil {
		t.Fatalf("创建护理记录失败: %v", err)
	}
}

func TestInpatientService_CreateNursingRecord_Nil(t *testing.T) {
	svc, _ := setupInpatientService(t)

	err := svc.CreateNursingRecord(nil)
	if err == nil {
		t.Error("期望 nil 护理记录时返回错误")
	}
}
