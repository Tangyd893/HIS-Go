package service

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"his-go/internal/notification/model"
	"his-go/internal/notification/repository"
)

func setupNotificationService(t *testing.T) (*NotificationService, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.New(sqlite.Config{DriverName: "sqlite", DSN: ":memory:"}), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接 SQLite 失败: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			template_id TEXT,
			receiver_id TEXT NOT NULL,
			title TEXT,
			content TEXT,
			channel INTEGER NOT NULL DEFAULT 1,
			status INTEGER DEFAULT 0,
			send_time DATETIME,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS notification_templates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			channel INTEGER NOT NULL DEFAULT 1,
			title_template TEXT,
			content_template TEXT NOT NULL,
			params TEXT,
			status INTEGER DEFAULT 1,
			created_at DATETIME,
			deleted_at DATETIME
		)`,
	}
	for _, ddl := range tables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("建表失败: %v", err)
		}
	}

	repo := repository.NewNotificationRepository(db)
	svc := NewNotificationService(repo)

	return svc, db
}

func TestNotificationService_Send(t *testing.T) {
	svc, _ := setupNotificationService(t)

	notif := &model.Notification{
		ID:         "notif-001",
		ReceiverID: "user-001",
		Title:      "就诊提醒",
		Content:    "您有明天上午的门诊预约",
		Channel:    3,
	}

	err := svc.Send(notif)
	if err != nil {
		t.Fatalf("发送通知失败: %v", err)
	}
}

func TestNotificationService_Send_NilNotif(t *testing.T) {
	svc, _ := setupNotificationService(t)

	err := svc.Send(nil)
	if err == nil {
		t.Error("期望 nil 通知时返回错误")
	}
}

func TestNotificationService_Send_EmptyReceiver(t *testing.T) {
	svc, _ := setupNotificationService(t)

	notif := &model.Notification{
		ID:      "notif-002",
		Title:   "测试",
		Channel: 1,
	}

	err := svc.Send(notif)
	if err == nil {
		t.Error("期望接收人为空时返回错误")
	}
}

func TestNotificationService_Send_InvalidChannel(t *testing.T) {
	svc, _ := setupNotificationService(t)

	notif := &model.Notification{
		ID:         "notif-003",
		ReceiverID: "user-001",
		Title:      "测试",
		Channel:    99,
	}

	err := svc.Send(notif)
	if err == nil {
		t.Error("期望无效渠道时返回错误")
	}
}

func TestNotificationService_BatchSend_EmptyReceivers(t *testing.T) {
	svc, _ := setupNotificationService(t)

	err := svc.BatchSend([]string{}, "tpl-001", map[string]string{"name": "test"})
	if err == nil {
		t.Error("期望空接收人列表时返回错误")
	}
}

func TestNotificationService_BatchSend_EmptyTemplateID(t *testing.T) {
	svc, _ := setupNotificationService(t)

	err := svc.BatchSend([]string{"user-001"}, "", map[string]string{})
	if err == nil {
		t.Error("期望空模板ID时返回错误")
	}
}

func TestNotificationService_CreateTemplate(t *testing.T) {
	svc, db := setupNotificationService(t)

	tmpl := &model.NotificationTemplate{
		ID:              "tpl-001",
		Name:            "就诊提醒模板",
		Channel:         3,
		TitleTemplate:   "就诊提醒 - {patientName}",
		ContentTemplate: "尊敬的{patientName}，您预约了{date}的{dept}门诊，请按时就诊。",
	}

	err := svc.CreateTemplate(tmpl)
	if err != nil {
		t.Fatalf("创建模板失败: %v", err)
	}

	var saved model.NotificationTemplate
	if err := db.Where("id = ?", "tpl-001").First(&saved).Error; err != nil {
		t.Fatalf("查询模板失败: %v", err)
	}
	if saved.Name != "就诊提醒模板" {
		t.Errorf("期望Name='就诊提醒模板'，实际=%s", saved.Name)
	}
}

func TestNotificationService_CreateTemplate_Nil(t *testing.T) {
	svc, _ := setupNotificationService(t)

	err := svc.CreateTemplate(nil)
	if err == nil {
		t.Error("期望 nil 模板时返回错误")
	}
}

func TestNotificationService_CreateTemplate_EmptyName(t *testing.T) {
	svc, _ := setupNotificationService(t)

	tmpl := &model.NotificationTemplate{
		ID:              "tpl-002",
		Channel:         1,
		ContentTemplate: "内容",
	}

	err := svc.CreateTemplate(tmpl)
	if err == nil {
		t.Error("期望空名称时返回错误")
	}
}

func TestNotificationService_CreateTemplate_EmptyContent(t *testing.T) {
	svc, _ := setupNotificationService(t)

	tmpl := &model.NotificationTemplate{
		ID:      "tpl-003",
		Name:    "测试",
		Channel: 1,
	}

	err := svc.CreateTemplate(tmpl)
	if err == nil {
		t.Error("期望空内容模板时返回错误")
	}
}

func TestNotificationService_ConsumeNotification_EmptyTemplateID(t *testing.T) {
	svc, _ := setupNotificationService(t)

	err := svc.ConsumeNotification("", "user-001", map[string]string{})
	if err == nil {
		t.Error("期望空模板ID时返回错误")
	}
}

func TestNotificationService_ConsumeNotification_EmptyReceiver(t *testing.T) {
	svc, _ := setupNotificationService(t)

	err := svc.ConsumeNotification("tpl-001", "", map[string]string{})
	if err == nil {
		t.Error("期望空接收人时返回错误")
	}
}
