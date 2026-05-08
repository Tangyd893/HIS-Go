package notification

import (
	"context"
	"strings"

	"his-go/api/proto/common"
	"his-go/api/proto/notification"
	notifmodel "his-go/internal/notification/model"
	notifsvc "his-go/internal/notification/service"
)

// NotificationGrpcServer gRPC 消息通知服务实现
type NotificationGrpcServer struct {
	notification.UnimplementedNotificationServiceServer
	svc *notifsvc.NotificationService
}

// NewNotificationGrpcServer 创建 gRPC 消息通知服务
func NewNotificationGrpcServer(svc *notifsvc.NotificationService) *NotificationGrpcServer {
	return &NotificationGrpcServer{svc: svc}
}

// SendNotification 发送通知
func (s *NotificationGrpcServer) SendNotification(ctx context.Context, req *notification.Notification) (*notification.Notification, error) {
	notif := protoToNotification(req)
	if err := s.svc.Send(notif); err != nil {
		return nil, err
	}
	return notificationToProto(notif), nil
}

// BatchSend 批量发送通知
func (s *NotificationGrpcServer) BatchSend(ctx context.Context, req *notification.BatchSendRequest) (*notification.BatchSendResponse, error) {
	if err := s.svc.BatchSend(req.ReceiverIds, req.TemplateId, req.Params); err != nil {
		return nil, err
	}
	return &notification.BatchSendResponse{
		Base:         &common.BaseResponse{Code: 0, Message: "批量发送完成"},
		SuccessCount: int32(len(req.ReceiverIds)),
		FailCount:    0,
	}, nil
}

// ListTemplates 查询通知模板列表
func (s *NotificationGrpcServer) ListTemplates(ctx context.Context, req *common.Empty) (*notification.TemplateListResponse, error) {
	_ = req
	templates, err := s.svc.ListTemplates()
	if err != nil {
		return nil, err
	}
	pbList := make([]*notification.NotificationTemplate, len(templates))
	for i, t := range templates {
		pbList[i] = templateToProto(&t)
	}
	return &notification.TemplateListResponse{
		Base:      &common.BaseResponse{Code: 0, Message: "查询成功"},
		Templates: pbList,
	}, nil
}

// CreateTemplate 创建通知模板
func (s *NotificationGrpcServer) CreateTemplate(ctx context.Context, req *notification.NotificationTemplate) (*notification.NotificationTemplate, error) {
	tmpl := protoToTemplate(req)
	if err := s.svc.CreateTemplate(tmpl); err != nil {
		return nil, err
	}
	return templateToProto(tmpl), nil
}

// ---- 转换辅助函数 ----

func notificationToProto(n *notifmodel.Notification) *notification.Notification {
	pb := &notification.Notification{
		Id:         n.ID,
		TemplateId: n.TemplateID,
		ReceiverId: n.ReceiverID,
		Title:      n.Title,
		Content:    n.Content,
		Channel:    int32(n.Channel),
		Status:     int32(n.Status),
		CreatedAt:  n.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if !n.SendTime.IsZero() {
		pb.SendTime = n.SendTime.Format("2006-01-02 15:04:05")
	}
	return pb
}

func protoToNotification(pb *notification.Notification) *notifmodel.Notification {
	return &notifmodel.Notification{
		ID:         pb.Id,
		TemplateID: pb.TemplateId,
		ReceiverID: pb.ReceiverId,
		Title:      pb.Title,
		Content:    pb.Content,
		Channel:    int8(pb.Channel),
		Status:     int8(pb.Status),
	}
}

func templateToProto(t *notifmodel.NotificationTemplate) *notification.NotificationTemplate {
	return &notification.NotificationTemplate{
		Id:              t.ID,
		Name:            t.Name,
		TitleTemplate:   t.TitleTemplate,
		ContentTemplate: t.ContentTemplate,
		Channel:         int32(t.Channel),
		Params:          strings.Split(t.Params, ","),
	}
}

func protoToTemplate(pb *notification.NotificationTemplate) *notifmodel.NotificationTemplate {
	return &notifmodel.NotificationTemplate{
		ID:              pb.Id,
		Name:            pb.Name,
		TitleTemplate:   pb.TitleTemplate,
		ContentTemplate: pb.ContentTemplate,
		Channel:         int8(pb.Channel),
		Params:          strings.Join(pb.Params, ","),
	}
}
