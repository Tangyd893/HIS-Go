package service

import (
	"testing"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/pharmacy/model"
	"his-go/internal/pharmacy/repository"
)

func setupPharmacyService(t *testing.T) (*PharmacyService, *gorm.DB, *cron.Cron) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS drugs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			generic_name TEXT,
			specification TEXT,
			manufacturer TEXT,
			batch_no TEXT,
			purchase_price REAL NOT NULL DEFAULT 0,
			retail_price REAL NOT NULL DEFAULT 0,
			stock INTEGER NOT NULL DEFAULT 0,
			min_stock INTEGER NOT NULL DEFAULT 0,
			expiry_date TEXT,
			status INTEGER DEFAULT 1,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS dispense_records (
			id TEXT PRIMARY KEY,
			prescription_id TEXT NOT NULL,
			patient_id TEXT NOT NULL,
			drug_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			dispenser_id TEXT NOT NULL,
			checker_id TEXT,
			status INTEGER DEFAULT 1,
			created_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v\nSQL: %s", err, ddl)
		}
	}

	repo := repository.NewPharmacyRepository(db)
	cr := cron.New(cron.WithSeconds())
	svc := NewPharmacyService(repo, nil, cr)

	return svc, db, cr
}

func createTestDrug(t *testing.T, db *gorm.DB, id, name, batchNo string, stock int, expiryDate string, retailPrice float64) *model.Drug {
	t.Helper()
	drug := &model.Drug{
		ID:          id,
		Name:        name,
		GenericName: name + "通用名",
		BatchNo:     batchNo,
		Stock:       stock,
		MinStock:    10,
		ExpiryDate:  expiryDate,
		RetailPrice: retailPrice,
		Status:      1,
	}
	if err := db.Create(drug).Error; err != nil {
		t.Fatalf("创建测试药品失败: %v", err)
	}
	return drug
}

func TestPharmacyService_ListDrugs(t *testing.T) {
	svc, db, cr := setupPharmacyService(t)
	defer cr.Stop()

	createTestDrug(t, db, "drug-1", "阿莫西林胶囊", "BATCH001", 100, "2027-12-31", 25.50)
	createTestDrug(t, db, "drug-2", "头孢克肟片", "BATCH002", 50, "2027-06-30", 55.00)
	createTestDrug(t, db, "drug-3", "布洛芬缓释胶囊", "BATCH003", 200, "2028-03-15", 18.00)

	list, total, err := svc.ListDrugs("", 1, 10)
	if err != nil {
		t.Fatalf("查询药品列表失败: %v", err)
	}
	if total != 3 {
		t.Errorf("期望total=3，实际=%d", total)
	}
	if len(list) != 3 {
		t.Errorf("期望len(list)=3，实际=%d", len(list))
	}
}

func TestPharmacyService_ListDrugs_FilterByName(t *testing.T) {
	svc, db, cr := setupPharmacyService(t)
	defer cr.Stop()

	createTestDrug(t, db, "drug-4", "阿莫西林胶囊", "BATCH004", 100, "2027-12-31", 25.50)
	createTestDrug(t, db, "drug-5", "阿奇霉素片", "BATCH005", 60, "2027-08-20", 32.00)
	createTestDrug(t, db, "drug-6", "头孢克肟片", "BATCH006", 50, "2027-06-30", 55.00)

	list, total, err := svc.ListDrugs("阿", 1, 10)
	if err != nil {
		t.Fatalf("查询药品列表失败: %v", err)
	}
	if total != 2 {
		t.Errorf("期望total=2，实际=%d", total)
	}
	if len(list) != 2 {
		t.Errorf("期望len(list)=2，实际=%d", len(list))
	}
}

func TestPharmacyService_GetDrug(t *testing.T) {
	svc, db, cr := setupPharmacyService(t)
	defer cr.Stop()

	createTestDrug(t, db, "drug-get", "阿司匹林", "BATCH-GET", 80, "2027-11-01", 12.00)

	drug, err := svc.GetDrug("drug-get")
	if err != nil {
		t.Fatalf("获取药品失败: %v", err)
	}
	if drug.Name != "阿司匹林" {
		t.Errorf("期望Name='阿司匹林'，实际=%s", drug.Name)
	}
	if drug.Stock != 80 {
		t.Errorf("期望Stock=80，实际=%d", drug.Stock)
	}
}

func TestPharmacyService_GetDrug_NotFound(t *testing.T) {
	svc, _, cr := setupPharmacyService(t)
	defer cr.Stop()

	_, err := svc.GetDrug("nonexistent")
	if err == nil {
		t.Error("期望不存在的药品时返回错误")
	}
}

func TestPharmacyService_AddStock(t *testing.T) {
	svc, db, cr := setupPharmacyService(t)
	defer cr.Stop()

	createTestDrug(t, db, "drug-stock", "维生素C片", "BATCH-STOCK", 50, "2027-09-30", 5.00)

	err := svc.AddStock("drug-stock", 30)
	if err != nil {
		t.Fatalf("增加库存失败: %v", err)
	}

	var updated model.Drug
	if err := db.Where("id = ?", "drug-stock").First(&updated).Error; err != nil {
		t.Fatalf("查询药品失败: %v", err)
	}
	if updated.Stock != 80 {
		t.Errorf("期望Stock=80，实际=%d", updated.Stock)
	}
}

func TestPharmacyService_AddStock_NotFound(t *testing.T) {
	svc, _, cr := setupPharmacyService(t)
	defer cr.Stop()

	err := svc.AddStock("nonexistent", 10)
	if err == nil {
		t.Error("期望不存在的药品时返回错误")
	}
}

func TestPharmacyService_DispenseDrug_Success(t *testing.T) {
	svc, db, cr := setupPharmacyService(t)
	defer cr.Stop()

	createTestDrug(t, db, "drug-dispense", "对乙酰氨基酚片", "BATCH-DISP", 100, "2027-12-31", 15.00)

	err := svc.DispenseDrug("prescription-1", "drug-dispense", 5, "dispenser-zhang")
	if err != nil {
		t.Fatalf("发药失败: %v", err)
	}

	var updated model.Drug
	if err := db.Where("id = ?", "drug-dispense").First(&updated).Error; err != nil {
		t.Fatalf("查询药品失败: %v", err)
	}
	if updated.Stock != 95 {
		t.Errorf("期望Stock=95，实际=%d", updated.Stock)
	}

	var records []model.DispenseRecord
	if err := db.Where("drug_id = ?", "drug-dispense").Find(&records).Error; err != nil {
		t.Fatalf("查询发药记录失败: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("期望1条发药记录，实际=%d", len(records))
	}
	if records[0].Quantity != 5 {
		t.Errorf("期望Quantity=5，实际=%d", records[0].Quantity)
	}
	if records[0].DispenserID != "dispenser-zhang" {
		t.Errorf("期望DispenserID='dispenser-zhang'，实际=%s", records[0].DispenserID)
	}
	if records[0].Status != 1 {
		t.Errorf("期望Status=1(已发药)，实际=%d", records[0].Status)
	}
}

func TestPharmacyService_DispenseDrug_StockInsufficient(t *testing.T) {
	svc, db, cr := setupPharmacyService(t)
	defer cr.Stop()

	createTestDrug(t, db, "drug-low", "胰岛素注射液", "BATCH-LOW", 3, "2027-10-15", 180.00)

	err := svc.DispenseDrug("prescription-2", "drug-low", 10, "dispenser-li")
	if err == nil {
		t.Error("期望库存不足时返回错误")
	}

	var unchanged model.Drug
	if err := db.Where("id = ?", "drug-low").First(&unchanged).Error; err != nil {
		t.Fatalf("查询药品失败: %v", err)
	}
	if unchanged.Stock != 3 {
		t.Errorf("期望库存未变Stock=3，实际=%d", unchanged.Stock)
	}
}

func TestPharmacyService_DispenseDrug_NonexistentDrug(t *testing.T) {
	svc, _, cr := setupPharmacyService(t)
	defer cr.Stop()

	err := svc.DispenseDrug("prescription-3", "nonexistent", 1, "dispenser-wang")
	if err == nil {
		t.Error("期望不存在的药品时返回错误")
	}
}
