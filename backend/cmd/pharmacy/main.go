// Package main 药房管理服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"his-go/internal/pharmacy/handler"
	"his-go/internal/pharmacy/repository"
	"his-go/internal/pharmacy/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/mq"
	"his-go/pkg/redis"

	"his-go/api/proto/pharmacy"
	grpcpharm "his-go/internal/pharmacy"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8087

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Pharmacy 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_pharmacy", cfg.Database.SSLMode,
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
		logger.Warn("RabbitMQ 不可用，过期药品告警通知功能已禁用（演示环境可忽略）: " + err.Error())
		rabbitMQ = nil
	}
	if rabbitMQ != nil {
		defer rabbitMQ.Close()
	}

	cronScheduler := cron.New(cron.WithSeconds())
	cronScheduler.Start()
	defer cronScheduler.Stop()

	pharmacyRepo := repository.NewPharmacyRepository(db)
	pharmacySvc := service.NewPharmacyService(pharmacyRepo, rabbitMQ, cronScheduler)
	pharmacyHandler := handler.NewPharmacyHandler(pharmacySvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupPharmacyRouter(cfg, pharmacyHandler, deps)

	go startGrpcServer(pharmacySvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("药房管理服务已启动")
	log.Printf("[Pharmacy] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupPharmacyRouter(cfg *config.Config, pharmacyHandler *handler.PharmacyHandler, deps *health.Dependencies) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-pharmacy"))

	router.GET("/health", health.HealthHandler("his-pharmacy"))
	router.GET("/ready", health.ReadinessHandler("his-pharmacy", deps))

	api := router.Group("/api/pharmacy")
	{
		api.GET("/drugs", pharmacyHandler.ListDrugs)
		api.GET("/drug/:id", pharmacyHandler.GetDrug)
		api.POST("/stock/add/:id", pharmacyHandler.AddStock)
		api.POST("/dispense", pharmacyHandler.DispenseDrug)
		api.GET("/expired", pharmacyHandler.CheckExpiredDrugs)
	}

	return router
}

func startGrpcServer(svc *service.PharmacyService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9087)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcpharm.NewPharmacyGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	pharmacy.RegisterPharmacyServiceServer(s, grpcSrv)

	log.Printf("[Pharmacy] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
