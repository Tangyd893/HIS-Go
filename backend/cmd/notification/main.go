// Package main 消息通知服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/notification/handler"
	"his-go/internal/notification/repository"
	"his-go/internal/notification/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/mq"
	"his-go/pkg/redis"

	"his-go/api/proto/notification"
	grpcnotif "his-go/internal/notification"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8094

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Notification 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_notification", cfg.Database.SSLMode,
		cfg.Database.MaxIdleConns, cfg.Database.MaxOpenConns, cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		logger.Fatal("数据库连接失败: " + err.Error())
	}

	rdb, err := redis.NewClient(
		cfg.Redis.Host, cfg.Redis.Port,
		cfg.Redis.Password, cfg.Redis.DB, cfg.Redis.PoolSize,
	)
	if err != nil {
		logger.Fatal("Redis 连接失败: " + err.Error())
	}
	_ = rdb

	rabbitMQ, err := mq.NewRabbitMQ(
		cfg.RabbitMQ.Host, cfg.RabbitMQ.Port,
		cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.VHost,
	)
	if err != nil {
		logger.Fatal("RabbitMQ 连接失败: " + err.Error())
	}
	defer rabbitMQ.Close()
	_ = rabbitMQ

	notifRepo := repository.NewNotificationRepository(db)
	notifSvc := service.NewNotificationService(notifRepo)
	notifHandler := handler.NewNotificationHandler(notifSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupNotificationRouter(cfg, notifHandler, deps)

	go startGrpcServer(notifSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("消息通知服务已启动")
	log.Printf("[Notification] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupNotificationRouter(cfg *config.Config, notifHandler *handler.NotificationHandler, deps *health.Dependencies) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-notification"))

	router.GET("/health", health.HealthHandler("his-notification"))
	router.GET("/ready", health.ReadinessHandler("his-notification", deps))

	api := router.Group("/api/notification")
	{
		api.POST("/send", notifHandler.SendNotification)
		api.POST("/batch-send", notifHandler.BatchSend)
		api.GET("/templates", notifHandler.ListTemplates)
		api.POST("/template", notifHandler.CreateTemplate)
	}

	return router
}

func startGrpcServer(svc *service.NotificationService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9094)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcnotif.NewNotificationGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	notification.RegisterNotificationServiceServer(s, grpcSrv)

	log.Printf("[Notification] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
