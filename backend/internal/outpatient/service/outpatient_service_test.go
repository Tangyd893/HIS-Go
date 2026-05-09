package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/outpatient/model"
	"his-go/internal/outpatient/repository"
)

func setupOutpatientService(t *testing.T) *OutpatientService {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS health_data (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			data_type TEXT NOT NULL,
			value TEXT,
			unit TEXT,
			measure_time TEXT,
			abnormal INTEGER DEFAULT 0,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS consultations (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			doctor_id TEXT NOT NULL,
			type INTEGER NOT NULL,
			description TEXT,
			images TEXT,
			status INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS chronic_contracts (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			doctor_id TEXT NOT NULL,
			disease_type TEXT,
			contract_date TEXT,
			end_date TEXT,
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

	repo := repository.NewOutpatientRepository(db)
	return NewOutpatientService(repo)
}

// ==================== isBloodPressureAbnormal 测试 ====================

func TestBloodPressureNormal(t *testing.T) {
	svc := setupOutpatientService(t)

	normalValues := []string{"120/80", "130/85", "110/70", "140/90", "90/60"}
	for _, v := range normalValues {
		if svc.isBloodPressureAbnormal(v) {
			t.Errorf("期望%s为正常血压，但判定为异常", v)
		}
	}
}

func TestBloodPressureHigh(t *testing.T) {
	svc := setupOutpatientService(t)

	highValues := []string{"141/80", "160/95", "150/100", "180/120", "140/91"}
	for _, v := range highValues {
		if !svc.isBloodPressureAbnormal(v) {
			t.Errorf("期望%s为异常高血压，但判定为正常", v)
		}
	}
}

func TestBloodPressureLow(t *testing.T) {
	svc := setupOutpatientService(t)

	lowValues := []string{"89/60", "80/50", "70/50", "89/59", "85/55"}
	for _, v := range lowValues {
		if !svc.isBloodPressureAbnormal(v) {
			t.Errorf("期望%s为异常低血压，但判定为正常", v)
		}
	}
}

func TestBloodPressureInvalidFormat(t *testing.T) {
	svc := setupOutpatientService(t)

	invalidValues := []string{"120", "120/80/90", "abc/80", "120/abc", "", "/", "120/"}
	for _, v := range invalidValues {
		if svc.isBloodPressureAbnormal(v) {
			t.Errorf("期望格式无效的'%s'判定为正常，但判为异常", v)
		}
	}
}

func TestBloodPressureWhitespace(t *testing.T) {
	svc := setupOutpatientService(t)

	if !svc.isBloodPressureAbnormal(" 150 / 100 ") {
		t.Error("期望含空格的150/100正常判定为异常")
	}

	if svc.isBloodPressureAbnormal(" 120 / 80 ") {
		t.Error("期望含空格的120/80正常判定为正常")
	}
}

// ==================== isBloodSugarAbnormal 测试 ====================

func TestBloodSugarNormal(t *testing.T) {
	svc := setupOutpatientService(t)

	normalValues := []string{"4.5", "5.5", "6.0", "3.9", "7.0"}
	for _, v := range normalValues {
		if svc.isBloodSugarAbnormal(v) {
			t.Errorf("期望血糖%s为正常，但判定为异常", v)
		}
	}
}

func TestBloodSugarHigh(t *testing.T) {
	svc := setupOutpatientService(t)

	highValues := []string{"7.1", "8.0", "12.5", "20.0", "7.01"}
	for _, v := range highValues {
		if !svc.isBloodSugarAbnormal(v) {
			t.Errorf("期望血糖%s为异常偏高，但判定为正常", v)
		}
	}
}

func TestBloodSugarLow(t *testing.T) {
	svc := setupOutpatientService(t)

	lowValues := []string{"3.8", "3.0", "2.5", "1.0", "3.89"}
	for _, v := range lowValues {
		if !svc.isBloodSugarAbnormal(v) {
			t.Errorf("期望血糖%s为异常偏低，但判定为正常", v)
		}
	}
}

func TestBloodSugarInvalidFormat(t *testing.T) {
	svc := setupOutpatientService(t)

	invalidValues := []string{"abc", "", "12.3.4", "nan"}
	for _, v := range invalidValues {
		if svc.isBloodSugarAbnormal(v) {
			t.Errorf("期望无效血糖值'%s'判定为正常，但判为异常", v)
		}
	}
}

func TestBloodSugarNegative(t *testing.T) {
	svc := setupOutpatientService(t)

	if !svc.isBloodSugarAbnormal("-1.0") {
		t.Error("期望负血糖值判定为异常")
	}
}

// ==================== CheckHealthDataAbnormal 路由测试 ====================

func TestCheckHealthDataAbnormal_BloodPressure(t *testing.T) {
	svc := setupOutpatientService(t)

	data := &model.HealthData{
		DataType: "blood_pressure",
		Value:    "150/100",
	}

	if !svc.CheckHealthDataAbnormal(data) {
		t.Error("期望高血压数据判定为异常")
	}

	data.Value = "120/80"
	if svc.CheckHealthDataAbnormal(data) {
		t.Error("期望正常血压数据判定为非异常")
	}
}

func TestCheckHealthDataAbnormal_BloodSugar(t *testing.T) {
	svc := setupOutpatientService(t)

	data := &model.HealthData{
		DataType: "blood_sugar",
		Value:    "8.0",
	}

	if !svc.CheckHealthDataAbnormal(data) {
		t.Error("期望高血糖数据判定为异常")
	}

	data.Value = "5.5"
	if svc.CheckHealthDataAbnormal(data) {
		t.Error("期望正常血糖数据判定为非异常")
	}
}

func TestCheckHealthDataAbnormal_UnknownType(t *testing.T) {
	svc := setupOutpatientService(t)

	unknownTypes := []string{"weight", "height", "", "heart_rate"}
	for _, dt := range unknownTypes {
		data := &model.HealthData{
			DataType: dt,
			Value:    "999",
		}

		if svc.CheckHealthDataAbnormal(data) {
			t.Errorf("期望未知类型%s判定为非异常", dt)
		}
	}
}

// ==================== ReportHealthData 测试 ====================

func TestReportHealthData_MarksAbnormal(t *testing.T) {
	svc := setupOutpatientService(t)

	data := &model.HealthData{
		ID:        "hd-1",
		PatientID: "patient-1",
		DataType:  "blood_pressure",
		Value:     "160/100",
	}

	err := svc.ReportHealthData(data)
	if err != nil {
		t.Fatalf("上报健康数据失败: %v", err)
	}

	if !data.Abnormal {
		t.Error("期望异常血压被标记为Abnormal=true")
	}
}

func TestReportHealthData_MarksNormal(t *testing.T) {
	svc := setupOutpatientService(t)

	data := &model.HealthData{
		ID:        "hd-2",
		PatientID: "patient-1",
		DataType:  "blood_pressure",
		Value:     "120/80",
	}

	err := svc.ReportHealthData(data)
	if err != nil {
		t.Fatalf("上报健康数据失败: %v", err)
	}

	if data.Abnormal {
		t.Error("期望正常血压被标记为Abnormal=false")
	}
}

// ==================== CreateConsultation 测试 ====================

func TestCreateConsultation_SetsStatusZero(t *testing.T) {
	svc := setupOutpatientService(t)

	c := &model.Consultation{
		ID:        "consult-1",
		PatientID: "patient-1",
		Status:    99,
	}

	err := svc.CreateConsultation(c)
	if err != nil {
		t.Fatalf("创建问诊失败: %v", err)
	}

	if c.Status != 0 {
		t.Errorf("创建后状态期望=0(待接诊)，实际=%d", c.Status)
	}
}

// ==================== CreateChronicContract 测试 ====================

func TestCreateChronicContract_SetsStatusOne(t *testing.T) {
	svc := setupOutpatientService(t)

	contract := &model.ChronicContract{
		ID:        "contract-1",
		PatientID: "patient-1",
		Status:    0,
	}

	err := svc.CreateChronicContract(contract)
	if err != nil {
		t.Fatalf("创建慢病签约失败: %v", err)
	}

	if contract.Status != 1 {
		t.Errorf("创建后状态期望=1(签约中)，实际=%d", contract.Status)
	}
}
