package repository

import (
	"time"

	"his-go/internal/pharmacy/model"
	"his-go/pkg/errors"

	"gorm.io/gorm"
)

// PharmacyRepository 药房管理数据仓库
type PharmacyRepository struct {
	db *gorm.DB
}

// NewPharmacyRepository 创建药房管理数据仓库
func NewPharmacyRepository(db *gorm.DB) *PharmacyRepository {
	return &PharmacyRepository{db: db}
}

// FindByID 根据ID查找药品
func (r *PharmacyRepository) FindByID(id string) (*model.Drug, error) {
	var drug model.Drug
	err := r.db.Where("id = ?", id).First(&drug).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "药品不存在")
		}
		return nil, err
	}
	return &drug, nil
}

// ListDrugs 分页查询药品列表
func (r *PharmacyRepository) ListDrugs(name string, page, pageSize int) ([]model.Drug, int64, error) {
	var drugs []model.Drug
	var total int64

	query := r.db.Model(&model.Drug{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&drugs).Error; err != nil {
		return nil, 0, err
	}

	return drugs, total, nil
}

// AddStock 增加药品库存
func (r *PharmacyRepository) AddStock(drugID string, qty int) error {
	result := r.db.Model(&model.Drug{}).Where("id = ?", drugID).
		Update("stock", gorm.Expr("stock + ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "药品不存在")
	}
	return nil
}

// ReduceStock 减少药品库存，检查库存是否充足
func (r *PharmacyRepository) ReduceStock(drugID string, qty int) error {
	result := r.db.Model(&model.Drug{}).
		Where("id = ? AND stock >= ?", drugID, qty).
		Update("stock", gorm.Expr("stock - ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeStockInsufficient, "库存不足")
	}
	return nil
}

// Dispense 发药：事务减库存 + 创建发药记录
func (r *PharmacyRepository) Dispense(prescriptionID string, drugID string, qty int, dispenserID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.Drug{}).
			Where("id = ? AND stock >= ?", drugID, qty).
			Update("stock", gorm.Expr("stock - ?", qty))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.NewAppError(errors.CodeStockInsufficient, "库存不足")
		}

		record := &model.DispenseRecord{
			PrescriptionID: prescriptionID,
			DrugID:         drugID,
			Quantity:       qty,
			DispenserID:    dispenserID,
			Status:         1,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindExpiredDrugs 查询过期药品
func (r *PharmacyRepository) FindExpiredDrugs() ([]model.Drug, error) {
	var drugs []model.Drug
	now := time.Now().Format("2006-01-02")
	err := r.db.Where("expiry_date < ?", now).Find(&drugs).Error
	if err != nil {
		return nil, err
	}
	return drugs, nil
}
