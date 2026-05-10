package service

import (
	"his-go/internal/notification/model"
	"his-go/internal/notification/repository"
	"his-go/pkg/errors"
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
	if notif == nil {
		return errors.NewAppError(errors.CodeParamInvalid, "通知对象不能为空")
	}
	if notif.ReceiverID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "接收人ID不能为空")
	}
	if notif.Channel < 1 || notif.Channel > 4 {
		return errors.NewAppError(errors.CodeParamInvalid, "通知渠道无效，有效值：1SMS、2邮件、3站内信、4微信")
	}
	if notif.Title == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "通知标题不能为空")
	}
	notif.Status = 0 // 待发送
	return s.repo.Send(notif)
}

// BatchSend 批量发送通知
func (s *NotificationService) BatchSend(receiverIDs []string, templateID string, params map[string]string) error {
	if len(receiverIDs) == 0 {
		return errors.NewAppError(errors.CodeParamInvalid, "接收人ID列表不能为空")
	}
	if templateID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "模板ID不能为空")
	}
	if params == nil {
		return errors.NewAppError(errors.CodeParamInvalid, "模板参数不能为nil")
	}
	return s.repo.BatchSend(receiverIDs, templateID, params)
}

// ListTemplates 查询模板列表
func (s *NotificationService) ListTemplates() ([]model.NotificationTemplate, error) {
	return s.repo.ListTemplates()
}

// CreateTemplate 创建通知模板
func (s *NotificationService) CreateTemplate(template *model.NotificationTemplate) error {
	if template == nil {
		return errors.NewAppError(errors.CodeParamInvalid, "模板对象不能为空")
	}
	if template.Name == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "模板名称不能为空")
	}
	if template.Channel < 1 || template.Channel > 4 {
		return errors.NewAppError(errors.CodeParamInvalid, "通知渠道无效，有效值：1SMS、2邮件、3站内信、4微信")
	}
	if template.ContentTemplate == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "模板内容不能为空")
	}
	return s.repo.CreateTemplate(template)
}

// ConsumeNotification 消费消息队列通知
func (s *NotificationService) ConsumeNotification(templateID string, receiverID string, params map[string]string) error {
	if templateID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "模板ID不能为空")
	}
	if receiverID == "" {
		return errors.NewAppError(errors.CodeParamInvalid, "接收人ID不能为空")
	}
	if params == nil {
		return errors.NewAppError(errors.CodeParamInvalid, "模板参数不能为nil")
	}
	return s.repo.ConsumeNotification(templateID, receiverID, params)
}
