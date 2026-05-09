package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/billing/model"
	"his-go/internal/billing/repository"
)

func setupBillingService(t *testing.T) (*BillingService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS bills (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			registration_id TEXT,
			bill_no TEXT NOT NULL UNIQUE,
			total_amount REAL NOT NULL,
			paid_amount REAL NOT NULL DEFAULT 0,
			pay_method INTEGER DEFAULT 0,
			status INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS bill_details (
			id TEXT PRIMARY KEY,
			bill_id TEXT NOT NULL,
			item_type INTEGER NOT NULL,
			item_name TEXT NOT NULL,
			unit_price REAL NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			amount REAL NOT NULL,
			created_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	repo := repository.NewBillingRepository(db)
	svc := NewBillingService(repo, nil)

	return svc, db
}

func TestBillingService_CreateBill(t *testing.T) {
	svc, db := setupBillingService(t)

	bill := &model.Bill{
		ID:          "bill-001",
		PatientID:   "patient-1",
		BillNo:      "B20260101001",
		TotalAmount: 350.00,
		Status:      0,
	}
	details := []model.BillDetail{
		{ID: "detail-1", ItemType: 1, ItemName: "挂号费", UnitPrice: 20.00, Quantity: 1, Amount: 20.00},
		{ID: "detail-2", ItemType: 3, ItemName: "阿莫西林", UnitPrice: 30.00, Quantity: 2, Amount: 60.00},
		{ID: "detail-3", ItemType: 2, ItemName: "血常规", UnitPrice: 100.00, Quantity: 1, Amount: 100.00},
		{ID: "detail-4", ItemType: 4, ItemName: "输液费", UnitPrice: 170.00, Quantity: 1, Amount: 170.00},
	}

	err := svc.CreateBill(bill, details)
	if err != nil {
		t.Fatalf("创建账单失败: %v", err)
	}

	var savedBill model.Bill
	if err := db.Where("id = ?", "bill-001").First(&savedBill).Error; err != nil {
		t.Fatalf("查询账单失败: %v", err)
	}
	if savedBill.BillNo != "B20260101001" {
		t.Errorf("期望BillNo='B20260101001'，实际=%s", savedBill.BillNo)
	}
	if savedBill.TotalAmount != 350.00 {
		t.Errorf("期望TotalAmount=350，实际=%f", savedBill.TotalAmount)
	}

	var savedDetails []model.BillDetail
	if err := db.Where("bill_id = ?", "bill-001").Find(&savedDetails).Error; err != nil {
		t.Fatalf("查询账单明细失败: %v", err)
	}
	if len(savedDetails) != 4 {
		t.Errorf("期望4条明细，实际=%d", len(savedDetails))
	}
}

func TestBillingService_GetBill(t *testing.T) {
	svc, db := setupBillingService(t)

	bill := &model.Bill{
		ID:          "bill-002",
		PatientID:   "patient-2",
		BillNo:      "B20260101002",
		TotalAmount: 120.00,
		Status:      0,
	}
	if err := db.Create(bill).Error; err != nil {
		t.Fatalf("创建账单失败: %v", err)
	}

	result, err := svc.GetBill("bill-002")
	if err != nil {
		t.Fatalf("获取账单失败: %v", err)
	}
	if result.BillNo != "B20260101002" {
		t.Errorf("期望BillNo='B20260101002'，实际=%s", result.BillNo)
	}
	if result.PatientID != "patient-2" {
		t.Errorf("期望PatientID='patient-2'，实际=%s", result.PatientID)
	}
}

func TestBillingService_GetBill_NotFound(t *testing.T) {
	svc, _ := setupBillingService(t)

	_, err := svc.GetBill("nonexistent")
	if err == nil {
		t.Error("期望不存在的账单时返回错误")
	}
}

func TestBillingService_ListBills(t *testing.T) {
	svc, db := setupBillingService(t)

	bills := []model.Bill{
		{ID: "bill-010", PatientID: "patient-x", BillNo: "B20260101010", TotalAmount: 100.00, Status: 0},
		{ID: "bill-011", PatientID: "patient-x", BillNo: "B20260101011", TotalAmount: 200.00, Status: 1},
		{ID: "bill-012", PatientID: "patient-x", BillNo: "B20260101012", TotalAmount: 300.00, Status: 2},
		{ID: "bill-013", PatientID: "patient-y", BillNo: "B20260101013", TotalAmount: 400.00, Status: 0},
	}
	for _, b := range bills {
		if err := db.Create(&b).Error; err != nil {
			t.Fatalf("创建测试账单失败: %v", err)
		}
	}

	list, total, err := svc.ListBills("patient-x", -1, 1, 10)
	if err != nil {
		t.Fatalf("查询账单列表失败: %v", err)
	}
	if total != 3 {
		t.Errorf("期望total=3，实际=%d", total)
	}
	if len(list) != 3 {
		t.Errorf("期望len(list)=3，实际=%d", len(list))
	}
}

func TestBillingService_ListBills_FilterByStatus(t *testing.T) {
	svc, db := setupBillingService(t)

	bills := []model.Bill{
		{ID: "bill-020", PatientID: "patient-z", BillNo: "B20260101020", TotalAmount: 100.00, Status: 0},
		{ID: "bill-021", PatientID: "patient-z", BillNo: "B20260101021", TotalAmount: 200.00, Status: 1},
		{ID: "bill-022", PatientID: "patient-z", BillNo: "B20260101022", TotalAmount: 300.00, Status: 0},
	}
	for _, b := range bills {
		if err := db.Create(&b).Error; err != nil {
			t.Fatalf("创建测试账单失败: %v", err)
		}
	}

	list, total, err := svc.ListBills("patient-z", 1, 1, 10)
	if err != nil {
		t.Fatalf("查询账单列表失败: %v", err)
	}
	if total != 1 {
		t.Errorf("期望total=1（仅已支付），实际=%d", total)
	}
	if len(list) != 1 {
		t.Errorf("期望len(list)=1，实际=%d", len(list))
	}
}

func TestBillingService_GetBillDetails(t *testing.T) {
	svc, db := setupBillingService(t)

	bill := &model.Bill{
		ID:          "bill-003",
		PatientID:   "patient-3",
		BillNo:      "B20260101003",
		TotalAmount: 250.00,
		Status:      0,
	}
	if err := db.Create(bill).Error; err != nil {
		t.Fatalf("创建账单失败: %v", err)
	}

	details := []model.BillDetail{
		{ID: "det-1", BillID: "bill-003", ItemType: 1, ItemName: "挂号费", UnitPrice: 20.00, Quantity: 1, Amount: 20.00},
		{ID: "det-2", BillID: "bill-003", ItemType: 3, ItemName: "头孢克肟", UnitPrice: 50.00, Quantity: 2, Amount: 100.00},
	}
	for _, d := range details {
		if err := db.Create(&d).Error; err != nil {
			t.Fatalf("创建账单明细失败: %v", err)
		}
	}

	result, err := svc.GetBillDetails("bill-003")
	if err != nil {
		t.Fatalf("获取账单明细失败: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("期望2条明细，实际=%d", len(result))
	}
}
