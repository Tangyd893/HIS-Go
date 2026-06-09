// Package main 健康档案服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/health_record/handler"
	"his-go/internal/health_record/repository"
	"his-go/internal/health_record/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"
	"his-go/pkg/security/jwt"

	"his-go/api/proto/health_record"
	grpchr "his-go/internal/health_record"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8093

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-HealthRecord 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_health_record", cfg.Database.SSLMode,
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

	hrRepo := repository.NewHealthRecordRepository(db)
	hrSvc := service.NewHealthRecordService(hrRepo)
	hrHandler := handler.NewHealthRecordHandler(hrSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupHealthRecordRouter(cfg, hrHandler, deps, jwtSvc)

	go startGrpcServer(hrSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("健康档案服务已启动")
	log.Printf("[HealthRecord] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupHealthRecordRouter(cfg *config.Config, hrHandler *handler.HealthRecordHandler, deps *health.Dependencies, jwtSvc *jwt.JWTService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-health-record"))

	router.GET("/health", health.HealthHandler("his-health-record"))
	router.GET("/ready", health.ReadinessHandler("his-health-record", deps))

	api := router.Group("/api/health-record")
	api.Use(middleware.UserContext(jwtSvc))
	{
		api.GET("/summary/:patientId", hrHandler.GetSummary)
		api.GET("/timeline/:patientId", hrHandler.GetTimeline)
		api.POST("/authorization/grant", hrHandler.GrantAuthorization)
		api.POST("/authorization/revoke", hrHandler.RevokeAuthorization)
	}

	return router
}

func startGrpcServer(svc *service.HealthRecordService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9093)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpchr.NewHealthRecordGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	health_record.RegisterHealthRecordServiceServer(s, grpcSrv)

	log.Printf("[HealthRecord] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
