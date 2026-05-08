package service

import (
	"his-go/internal/notification/model"
	"his-go/internal/notification/repository"
)

// NotificationService 消息通知业务服务
type NotificationService struct {
	repo *repository.NotificationRepository
}

// NewNotificationService 创建消息通知业务服务
func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// Send 发送单条通知
func (s *NotificationService) Send(notif *model.Notification) error {
	notif.Status = 0 // 待发送
	return s.repo.Send(notif)
}

// BatchSend 批量发送通知
func (s *NotificationService) BatchSend(receiverIDs []string, templateID string, params map[string]string) error {
	return s.repo.BatchSend(receiverIDs, templateID, params)
}

// ListTemplates 查询模板列表
func (s *NotificationService) ListTemplates() ([]model.NotificationTemplate, error) {
	return s.repo.ListTemplates()
}

// CreateTemplate 创建通知模板
func (s *NotificationService) CreateTemplate(template *model.NotificationTemplate) error {
	return s.repo.CreateTemplate(template)
}

// ConsumeNotification 消费消息队列通知
func (s *NotificationService) ConsumeNotification(templateID string, receiverID string, params map[string]string) error {
	return s.repo.ConsumeNotification(templateID, receiverID, params)
}
