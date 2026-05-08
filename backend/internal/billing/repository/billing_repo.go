package repository

import (
	"his-go/internal/billing/model"
	"his-go/pkg/errors"

	"gorm.io/gorm"
)

// BillingRepository 收费结算数据仓库
type BillingRepository struct {
	db *gorm.DB
}

// NewBillingRepository 创建收费结算数据仓库
func NewBillingRepository(db *gorm.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

// Create 事务创建账单及明细
func (r *BillingRepository) Create(bill *model.Bill, details []model.BillDetail) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(bill).Error; err != nil {
			return err
		}
		for i := range details {
			details[i].BillID = bill.ID
		}
		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByID 根据ID查找账单
func (r *BillingRepository) FindByID(id string) (*model.Bill, error) {
	var bill model.Bill
	err := r.db.Where("id = ?", id).First(&bill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "账单不存在")
		}
		return nil, err
	}
	return &bill, nil
}

// FindByBillNo 根据流水号查找账单
func (r *BillingRepository) FindByBillNo(billNo string) (*model.Bill, error) {
	var bill model.Bill
	err := r.db.Where("bill_no = ?", billNo).First(&bill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "账单不存在")
		}
		return nil, err
	}
	return &bill, nil
}

// ListByPatient 分页查询患者账单列表
func (r *BillingRepository) ListByPatient(patientID string, status int, page, pageSize int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	query := r.db.Model(&model.Bill{}).Where("patient_id = ?", patientID)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

// Pay 支付账单，使用乐观锁更新状态为已支付
func (r *BillingRepository) Pay(billID string, payMethod int8) error {
	result := r.db.Model(&model.Bill{}).
		Where("id = ? AND status = 0", billID).
		Updates(map[string]interface{}{
			"status":     1,
			"pay_method": payMethod,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodePayFailed, "账单状态不允许支付")
	}
	return nil
}

// Refund 退款，状态改为已退款
func (r *BillingRepository) Refund(billID string) error {
	result := r.db.Model(&model.Bill{}).
		Where("id = ? AND status = 1", billID).
		Update("status", 2)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRefundFailed, "账单状态不允许退款")
	}
	return nil
}

// FindDetailsByBillID 根据账单ID查找明细
func (r *BillingRepository) FindDetailsByBillID(billID string) ([]model.BillDetail, error) {
	var details []model.BillDetail
	err := r.db.Where("bill_id = ?", billID).Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}
