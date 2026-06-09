package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/system/model"
	"his-go/internal/system/repository"
)

func setupSystemService(t *testing.T) (*SystemService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS dict_types (
			id TEXT PRIMARY KEY,
			dict_name TEXT NOT NULL,
			dict_type TEXT NOT NULL UNIQUE,
			status INTEGER DEFAULT 1,
			remark TEXT,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS dict_items (
			id TEXT PRIMARY KEY,
			dict_type TEXT NOT NULL,
			label TEXT NOT NULL,
			value TEXT NOT NULL,
			sort_order INTEGER DEFAULT 0,
			status INTEGER DEFAULT 1,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS system_params (
			id TEXT PRIMARY KEY,
			param_key TEXT NOT NULL UNIQUE,
			param_value TEXT,
			description TEXT,
			status INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS operation_logs (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			username TEXT,
			module TEXT,
			action TEXT,
			method TEXT,
			url TEXT,
			ip TEXT,
			params TEXT,
			result TEXT,
			status INTEGER DEFAULT 1,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v", err)
		}
	}

	repo := repository.NewSystemRepository(db)
	svc := NewSystemService(repo)

	return svc, db
}

func TestSystemService_CreateDictItem(t *testing.T) {
	svc, db := setupSystemService(t)

	// 先创建字典类型
	db.Exec(`INSERT INTO dict_types (id, dict_name, dict_type, status) VALUES ('dt-1', '性别', 'gender', 1)`)

	item := &model.DictItem{
		ID:        "di-001",
		DictType:  "gender",
		Label:     "男",
		Value:     "M",
		SortOrder: 1,
	}

	err := svc.CreateDictItem(item)
	if err != nil {
		t.Fatalf("创建字典项失败: %v", err)
	}

	var saved model.DictItem
	if err := db.Where("id = ?", "di-001").First(&saved).Error; err != nil {
		t.Fatalf("查询字典项失败: %v", err)
	}
	if saved.Label != "男" {
		t.Errorf("期望Label='男'，实际=%s", saved.Label)
	}
	if saved.Status != 1 {
		t.Errorf("期望Status=1，实际=%d", saved.Status)
	}
}

func TestSystemService_ListDictItems(t *testing.T) {
	svc, db := setupSystemService(t)

	db.Exec(`INSERT INTO dict_types (id, dict_name, dict_type, status) VALUES ('dt-2', '支付方式', 'pay_method', 1)`)
	db.Exec(`INSERT INTO dict_items (id, dict_type, label, value, sort_order, status) VALUES ('di-1', 'pay_method', '现金', '1', 1, 1)`)
	db.Exec(`INSERT INTO dict_items (id, dict_type, label, value, sort_order, status) VALUES ('di-2', 'pay_method', '微信', '2', 2, 1)`)
	db.Exec(`INSERT INTO dict_items (id, dict_type, label, value, sort_order, status) VALUES ('di-3', 'pay_method', '支付宝', '3', 3, 0)`)

	items, err := svc.ListDictItems("pay_method")
	if err != nil {
		t.Fatalf("查询字典项失败: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("期望2条(status=1)，实际=%d", len(items))
	}
}

func TestSystemService_CreateOperationLog(t *testing.T) {
	svc, db := setupSystemService(t)

	log := &model.OperationLog{
		ID:     "log-001",
		UserID: "admin-1",
		Module: "系统管理",
		Action: "修改参数",
		Result: "修改了JWT过期时间",
	}

	err := svc.CreateOperationLog(log)
	if err != nil {
		t.Fatalf("创建操作日志失败: %v", err)
	}

	var saved model.OperationLog
	if err := db.Where("id = ?", "log-001").First(&saved).Error; err != nil {
		t.Fatalf("查询操作日志失败: %v", err)
	}
	if saved.Module != "系统管理" {
		t.Errorf("期望Module='系统管理'，实际=%s", saved.Module)
	}
}

func TestSystemService_UpdateParam(t *testing.T) {
	svc, db := setupSystemService(t)

	db.Exec(`INSERT INTO system_params (id, param_key, param_value, description, status) VALUES ('sp-1', 'jwt_expire_hour', '24', 'JWT过期时间(小时)', 1)`)

	param := &model.SystemParam{
		ID:         "sp-1",
		ParamKey:   "jwt_expire_hour",
		ParamValue: "48",
	}

	err := svc.UpdateParam(param)
	if err != nil {
		t.Fatalf("更新参数失败: %v", err)
	}

	var saved model.SystemParam
	if err := db.Where("id = ?", "sp-1").First(&saved).Error; err != nil {
		t.Fatalf("查询参数失败: %v", err)
	}
	if saved.ParamValue != "48" {
		t.Errorf("期望ParamValue='48'，实际=%s", saved.ParamValue)
	}
}
