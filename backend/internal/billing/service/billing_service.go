package service

import (
	"encoding/json"

	"his-go/internal/billing/model"
	"his-go/internal/billing/repository"
	"his-go/pkg/mq"
)

// BillingService 收费结算业务服务
type BillingService struct {
	repo *repository.BillingRepository
	mq   *mq.RabbitMQ
}

// NewBillingService 创建收费结算业务服务
func NewBillingService(repo *repository.BillingRepository, rabbitMQ *mq.RabbitMQ) *BillingService {
	return &BillingService{repo: repo, mq: rabbitMQ}
}

// CreateBill 创建账单
func (s *BillingService) CreateBill(bill *model.Bill, details []model.BillDetail) error {
	return s.repo.Create(bill, details)
}

// GetBill 获取账单详情
func (s *BillingService) GetBill(id string) (*model.Bill, error) {
	return s.repo.FindByID(id)
}

// ListBills 分页查询患者账单
func (s *BillingService) ListBills(patientID string, status int, page, pageSize int) ([]model.Bill, int64, error) {
	return s.repo.ListByPatient(patientID, status, page, pageSize)
}

// Pay 支付并发送缴费成功消息到 RabbitMQ
func (s *BillingService) Pay(billID string, payMethod int8) error {
	if err := s.repo.Pay(billID, payMethod); err != nil {
		return err
	}

	bill, err := s.repo.FindByID(billID)
	if err != nil {
		return err
	}

	msg, _ := json.Marshal(map[string]interface{}{
		"bill_id":    bill.ID,
		"bill_no":    bill.BillNo,
		"patient_id": bill.PatientID,
		"amount":     bill.TotalAmount,
		"pay_method": bill.PayMethod,
		"event":      "pay_success",
	})
	_ = s.mq.Publish("his.billing.pay", "billing.pay.success", msg)

	return nil
}

// Refund 退款并发送退费消息到 RabbitMQ
func (s *BillingService) Refund(billID string) error {
	if err := s.repo.Refund(billID); err != nil {
		return err
	}

	bill, err := s.repo.FindByID(billID)
	if err != nil {
		return err
	}

	msg, _ := json.Marshal(map[string]interface{}{
		"bill_id":    bill.ID,
		"bill_no":    bill.BillNo,
		"patient_id": bill.PatientID,
		"amount":     bill.TotalAmount,
		"event":      "refund_success",
	})
	_ = s.mq.Publish("his.billing.refund", "billing.refund.success", msg)

	return nil
}

// GetBillDetails 获取账单明细
func (s *BillingService) GetBillDetails(billID string) ([]model.BillDetail, error) {
	return s.repo.FindDetailsByBillID(billID)
}
