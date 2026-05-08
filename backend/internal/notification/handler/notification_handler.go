package handler

import (
	"github.com/gin-gonic/gin"

	"his-go/internal/notification/model"
	"his-go/internal/notification/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// NotificationHandler 消息通知接口处理器
type NotificationHandler struct {
	svc *service.NotificationService
}

// NewNotificationHandler 创建消息通知接口处理器
func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// BatchSendRequest 批量发送请求
type BatchSendRequest struct {
	ReceiverIDs []string          `json:"receiver_ids" binding:"required"`
	TemplateID  string            `json:"template_id" binding:"required"`
	Params      map[string]string `json:"params"`
}

// SendNotification 发送通知
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var notif model.Notification
	if err := c.ShouldBindJSON(&notif); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.Send(&notif); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, notif)
}

// BatchSend 批量发送通知
func (h *NotificationHandler) BatchSend(c *gin.Context) {
	var req BatchSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.BatchSend(req.ReceiverIDs, req.TemplateID, req.Params); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListTemplates 查询模板列表
func (h *NotificationHandler) ListTemplates(c *gin.Context) {
	list, err := h.svc.ListTemplates()
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateTemplate 创建通知模板
func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	var template model.NotificationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateTemplate(&template); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, template)
}
