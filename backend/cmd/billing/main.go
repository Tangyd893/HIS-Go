// Package main 收费结算服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/billing/handler"
	"his-go/internal/billing/repository"
	"his-go/internal/billing/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/mq"
	"his-go/pkg/redis"
	"his-go/pkg/security/jwt"

	"his-go/api/proto/billing"
	grpcbill "his-go/internal/billing"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8086

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Billing 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_billing", cfg.Database.SSLMode,
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

	jwtSvc := middleware.InitJWT(cfg)

	rabbitMQ, err := mq.NewRabbitMQ(
		cfg.RabbitMQ.Host, cfg.RabbitMQ.Port,
		cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.VHost,
	)
	if err != nil {
		logger.Warn("RabbitMQ 不可用，缴费成功/退款事件通知功能已禁用（演示环境可忽略）: " + err.Error())
		rabbitMQ = nil
	}
	if rabbitMQ != nil {
		defer rabbitMQ.Close()
	}

	billingRepo := repository.NewBillingRepository(db)
	billingSvc := service.NewBillingService(billingRepo, rabbitMQ)
	billingHandler := handler.NewBillingHandler(billingSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupBillingRouter(cfg, billingHandler, deps, jwtSvc)

	go startGrpcServer(billingSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("收费结算服务已启动")
	log.Printf("[Billing] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupBillingRouter(cfg *config.Config, billingHandler *handler.BillingHandler, deps *health.Dependencies, jwtSvc *jwt.JWTService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-billing"))

	router.GET("/health", health.HealthHandler("his-billing"))
	router.GET("/ready", health.ReadinessHandler("his-billing", deps))

	api := router.Group("/api/billing")
	api.Use(middleware.UserContext(jwtSvc))
	{
		api.POST("/create", billingHandler.CreateBill)
		api.GET("/:id", billingHandler.GetBill)
		api.POST("/pay/:id", billingHandler.PayBill)
		api.POST("/refund/:id", billingHandler.RefundBill)
		api.GET("/list", billingHandler.ListBills)
	}

	return router
}

func startGrpcServer(svc *service.BillingService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9086)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcbill.NewBillingGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	billing.RegisterBillingServiceServer(s, grpcSrv)

	log.Printf("[Billing] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
