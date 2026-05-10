// Package main 随访管理服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/followup/handler"
	"his-go/internal/followup/repository"
	"his-go/internal/followup/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"

	"his-go/api/proto/followup"
	grpcfup "his-go/internal/followup"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8092

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Followup 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_followup", cfg.Database.SSLMode,
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

	followupRepo := repository.NewFollowupRepository(db)
	followupSvc := service.NewFollowupService(followupRepo)
	followupHandler := handler.NewFollowupHandler(followupSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupFollowupRouter(cfg, followupHandler, deps)

	go startGrpcServer(followupSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("随访管理服务已启动")
	log.Printf("[Followup] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupFollowupRouter(cfg *config.Config, followupHandler *handler.FollowupHandler, deps *health.Dependencies) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-followup"))

	router.GET("/health", health.HealthHandler("his-followup"))
	router.GET("/ready", health.ReadinessHandler("his-followup", deps))

	api := router.Group("/api/followup")
	{
		api.POST("/plan", followupHandler.CreatePlan)
		api.GET("/plan/:id", followupHandler.GetPlan)
		api.GET("/plans", followupHandler.ListPlans)
		api.POST("/task/execute", followupHandler.ExecuteTask)
		api.POST("/survey", followupHandler.SubmitSurvey)
	}

	return router
}

func startGrpcServer(svc *service.FollowupService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9092)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcfup.NewFollowupGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	followup.RegisterFollowupServiceServer(s, grpcSrv)

	log.Printf("[Followup] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
