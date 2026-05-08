package service

import (
	"encoding/json"

	"his-go/internal/pharmacy/model"
	"his-go/internal/pharmacy/repository"
	"his-go/pkg/logger"
	"his-go/pkg/mq"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// PharmacyService 药房管理业务服务
type PharmacyService struct {
	repo *repository.PharmacyRepository
	mq   *mq.RabbitMQ
	cron *cron.Cron
}

// NewPharmacyService 创建药房管理业务服务
func NewPharmacyService(repo *repository.PharmacyRepository, rabbitMQ *mq.RabbitMQ, cronScheduler *cron.Cron) *PharmacyService {
	svc := &PharmacyService{
		repo: repo,
		mq:   rabbitMQ,
		cron: cronScheduler,
	}

	// 每30分钟执行一次过期药品扫描
	_, _ = cronScheduler.AddFunc("*/30 * * * *", svc.CheckAndAlertExpired)

	return svc
}

// ListDrugs 分页查询药品列表
func (s *PharmacyService) ListDrugs(name string, page, pageSize int) ([]model.Drug, int64, error) {
	return s.repo.ListDrugs(name, page, pageSize)
}

// GetDrug 获取药品详情
func (s *PharmacyService) GetDrug(id string) (*model.Drug, error) {
	return s.repo.FindByID(id)
}

// AddStock 增加药品库存
func (s *PharmacyService) AddStock(drugID string, qty int) error {
	return s.repo.AddStock(drugID, qty)
}

// DispenseDrug 发药
func (s *PharmacyService) DispenseDrug(prescriptionID string, drugID string, qty int, dispenserID string) error {
	return s.repo.Dispense(prescriptionID, drugID, qty, dispenserID)
}

// CheckAndAlertExpired 扫描过期药品并发送告警
func (s *PharmacyService) CheckAndAlertExpired() {
	drugs, err := s.repo.FindExpiredDrugs()
	if err != nil {
		logger.Error("扫描过期药品失败", zap.Error(err))
		return
	}

	for _, drug := range drugs {
		msg, _ := json.Marshal(map[string]interface{}{
			"drug_id":   drug.ID,
			"drug_name": drug.Name,
			"batch_no":  drug.BatchNo,
			"expiry":    drug.ExpiryDate,
			"stock":     drug.Stock,
			"event":     "drug_expired",
		})
		if err := s.mq.Publish("his.pharmacy.stock", "pharmacy.expired", msg); err != nil {
			logger.Error("发送过期药品告警失败", zap.Error(err))
		}
	}

	if len(drugs) > 0 {
		logger.Warn("发现过期药品", zap.Int("count", len(drugs)))
	}
}
