package service

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/user/model"
	"his-go/internal/user/repository"
)

func setupUserService(t *testing.T) (*UserService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS patients (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			id_card TEXT NOT NULL UNIQUE,
			phone TEXT,
			gender TEXT,
			birth_date DATETIME,
			address TEXT,
			allergy_history TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS employees (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			name TEXT NOT NULL,
			phone TEXT,
			dept_id TEXT,
			title TEXT,
			specialty TEXT,
			introduction TEXT,
			status INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS departments (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			parent_id TEXT,
			description TEXT,
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	patientRepo := repository.NewPatientRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)
	svc := NewUserService(patientRepo, empRepo, deptRepo, nil)

	return svc, db
}

func TestUserService_CreatePatient(t *testing.T) {
	svc, db := setupUserService(t)

	patient := &model.Patient{
		Name:    "张三",
		IdCard:  "110101199001011234",
		Phone:   "13800138000",
		Gender:  "男",
		Address: "北京市朝阳区",
	}

	err := svc.CreatePatient(patient)
	if err != nil {
		t.Fatalf("创建患者失败: %v", err)
	}
	if patient.ID == "" {
		t.Error("期望非空 ID")
	}

	var saved model.Patient
	if err := db.Where("id = ?", patient.ID).First(&saved).Error; err != nil {
		t.Fatalf("查询患者失败: %v", err)
	}
	if saved.Name != "张三" {
		t.Errorf("期望Name='张三'，实际=%s", saved.Name)
	}
}

func TestUserService_GetPatient(t *testing.T) {
	svc, db := setupUserService(t)

	patient := &model.Patient{
		ID:   "pat-001",
		Name: "李四", IdCard: "110101199002011234", Phone: "13900139000",
	}
	if err := db.Create(patient).Error; err != nil {
		t.Fatalf("创建患者失败: %v", err)
	}

	result, err := svc.GetPatient("pat-001")
	if err != nil {
		t.Fatalf("获取患者失败: %v", err)
	}
	if result.Name != "李四" {
		t.Errorf("期望Name='李四'，实际=%s", result.Name)
	}
}

func TestUserService_ListPatients(t *testing.T) {
	svc, db := setupUserService(t)

	patients := []model.Patient{
		{ID: "p-1", Name: "王五", IdCard: "110101199003011234", Phone: "13700137000"},
		{ID: "p-2", Name: "赵六", IdCard: "110101199004011234", Phone: "13600136000"},
		{ID: "p-3", Name: "王七", IdCard: "110101199005011234", Phone: "13500135000"},
	}
	for _, p := range patients {
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		if err := db.Create(&p).Error; err != nil {
			t.Fatalf("创建患者失败: %v", err)
		}
	}

	list, total, err := svc.ListPatients("王", "", 1, 10)
	if err != nil {
		t.Fatalf("查询患者列表失败: %v", err)
	}
	if total != 2 {
		t.Errorf("期望total=2，实际=%d", total)
	}
	if len(list) != 2 {
		t.Errorf("期望len(list)=2，实际=%d", len(list))
	}
}

func TestUserService_UpdatePatient(t *testing.T) {
	svc, db := setupUserService(t)

	patient := &model.Patient{
		ID: "pat-upd", Name: "原姓名", IdCard: "110101200001011234",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := db.Create(patient).Error; err != nil {
		t.Fatalf("创建患者失败: %v", err)
	}

	patient.Name = "新姓名"
	patient.Phone = "13300133000"
	err := svc.UpdatePatient(patient)
	if err != nil {
		t.Fatalf("更新患者失败: %v", err)
	}

	var updated model.Patient
	db.Where("id = ?", "pat-upd").First(&updated)
	if updated.Name != "新姓名" {
		t.Errorf("期望Name='新姓名'，实际=%s", updated.Name)
	}
}

func TestUserService_DeletePatient(t *testing.T) {
	svc, db := setupUserService(t)

	patient := &model.Patient{
		ID: "pat-del", Name: "待删除", IdCard: "110101200002011234",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := db.Create(patient).Error; err != nil {
		t.Fatalf("创建患者失败: %v", err)
	}

	err := svc.DeletePatient("pat-del")
	if err != nil {
		t.Fatalf("删除患者失败: %v", err)
	}

	var deleted model.Patient
	err = db.Unscoped().Where("id = ?", "pat-del").First(&deleted).Error
	if err != nil {
		t.Error("期望逻辑删除后仍可用Unscoped查到")
	}
	if deleted.DeletedAt.Time.IsZero() {
		t.Error("期望 DeletedAt 非零")
	}
}

func TestUserService_ListEmployees(t *testing.T) {
	svc, db := setupUserService(t)

	emps := []model.Employee{
		{ID: "emp-1", Name: "张医生", DeptID: "dept-1", Status: 1},
		{ID: "emp-2", Name: "李护士", DeptID: "dept-1", Status: 1},
		{ID: "emp-3", Name: "王主任", DeptID: "dept-2", Status: 1},
	}
	for _, e := range emps {
		if err := db.Create(&e).Error; err != nil {
			t.Fatalf("创建员工失败: %v", err)
		}
	}

	list, total, err := svc.ListEmployees("dept-1", "", 1, 10)
	if err != nil {
		t.Fatalf("查询员工列表失败: %v", err)
	}
	if total != 2 {
		t.Errorf("期望total=2，实际=%d", total)
	}
	if len(list) != 2 {
		t.Errorf("期望len(list)=2，实际=%d", len(list))
	}
}

func TestUserService_GetEmployee(t *testing.T) {
	svc, db := setupUserService(t)

	emp := &model.Employee{ID: "emp-get", Name: "刘专家", DeptID: "dept-x"}
	if err := db.Create(emp).Error; err != nil {
		t.Fatalf("创建员工失败: %v", err)
	}

	result, err := svc.GetEmployee("emp-get")
	if err != nil {
		t.Fatalf("获取员工失败: %v", err)
	}
	if result.Name != "刘专家" {
		t.Errorf("期望Name='刘专家'，实际=%s", result.Name)
	}
}

func TestUserService_ListDepartments(t *testing.T) {
	svc, db := setupUserService(t)

	depts := []model.Department{
		{ID: "dept-root", Name: "内科", ParentID: "", SortOrder: 1},
		{ID: "dept-child", Name: "心血管内科", ParentID: "dept-root", SortOrder: 1},
	}
	for _, d := range depts {
		if err := db.Create(&d).Error; err != nil {
			t.Fatalf("创建科室失败: %v", err)
		}
	}

	list, err := svc.ListDepartments()
	if err != nil {
		t.Fatalf("查询科室列表失败: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("期望根节点1个，实际=%d", len(list))
	}
	if len(list[0].Children) != 1 {
		t.Errorf("期望子节点1个，实际=%d", len(list[0].Children))
	}
}
