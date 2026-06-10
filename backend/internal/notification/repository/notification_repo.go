package repository

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"his-go/internal/notification/model"
	"his-go/pkg/errors"
)

// NotificationRepository 消息通知数据仓库
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository 创建消息通知数据仓库
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Send 创建通知记录
func (r *NotificationRepository) Send(notif *model.Notification) error {
	if err := r.db.Create(notif).Error; err != nil {
		return fmt.Errorf("发送通知失败: %w", err)
	}
	return nil
}

// BatchSend 批量创建通知
func (r *NotificationRepository) BatchSend(receiverIDs []string, templateID string, params map[string]string) error {
	var template model.NotificationTemplate
	if err := r.db.Where("id = ?", templateID).First(&template).Error; err != nil {
		return errors.WrapQueryError("模板", err)
	}

	notifications := make([]model.Notification, 0, len(receiverIDs))
	for _, receiverID := range receiverIDs {
		title := template.TitleTemplate
		content := template.ContentTemplate
		for k, v := range params {
			title = strings.ReplaceAll(title, "{"+k+"}", v)
			content = strings.ReplaceAll(content, "{"+k+"}", v)
		}

		notifications = append(notifications, model.Notification{
			TemplateID: templateID,
			ReceiverID: receiverID,
			Title:      title,
			Content:    content,
			Channel:    template.Channel,
			Status:     0,
		})
	}

	if len(notifications) > 0 {
		if err := r.db.Create(&notifications).Error; err != nil {
			return fmt.Errorf("批量发送通知失败: %w", err)
		}
	}

	return nil
}

// ListTemplates 查询模板列表
func (r *NotificationRepository) ListTemplates() ([]model.NotificationTemplate, error) {
	var list []model.NotificationTemplate
	if err := r.db.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, errors.WrapQueryError("模板列表", err)
	}
	return list, nil
}

// CreateTemplate 创建通知模板
func (r *NotificationRepository) CreateTemplate(template *model.NotificationTemplate) error {
	if err := r.db.Create(template).Error; err != nil {
		return errors.WrapCreateError("模板", err)
	}
	return nil
}

// ConsumeNotification 从消息队列消费并创建通知
func (r *NotificationRepository) ConsumeNotification(templateID string, receiverID string, params map[string]string) error {
	var template model.NotificationTemplate
	if err := r.db.Where("id = ?", templateID).First(&template).Error; err != nil {
		return errors.WrapQueryError("模板", err)
	}

	title := template.TitleTemplate
	content := template.ContentTemplate
	for k, v := range params {
		title = strings.ReplaceAll(title, "{"+k+"}", v)
		content = strings.ReplaceAll(content, "{"+k+"}", v)
	}

	now := time.Now()
	notif := model.Notification{
		TemplateID: templateID,
		ReceiverID: receiverID,
		Title:      title,
		Content:    content,
		Channel:    template.Channel,
		Status:     1, // 已发送
		SendTime:   &now,
	}

	if err := r.db.Create(&notif).Error; err != nil {
		return fmt.Errorf("消费通知失败: %w", err)
	}
	return nil
}
