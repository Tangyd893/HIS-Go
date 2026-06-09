// Package main 处方管理服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/prescription/handler"
	"his-go/internal/prescription/repository"
	"his-go/internal/prescription/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"
	"his-go/pkg/security/jwt"

	"his-go/api/proto/prescription"
	grpcps "his-go/internal/prescription"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8085

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Prescription 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_prescription", cfg.Database.SSLMode,
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

	prescRepo := repository.NewPrescriptionRepository(db)
	prescSvc := service.NewPrescriptionService(prescRepo)
	prescHandler := handler.NewPrescriptionHandler(prescSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupPrescriptionRouter(cfg, prescHandler, deps, jwtSvc)

	go startGrpcServer(prescSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("处方管理服务已启动")
	log.Printf("[Prescription] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupPrescriptionRouter(cfg *config.Config, prescHandler *handler.PrescriptionHandler, deps *health.Dependencies, jwtSvc *jwt.JWTService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-prescription"))

	router.GET("/health", health.HealthHandler("his-prescription"))
	router.GET("/ready", health.ReadinessHandler("his-prescription", deps))

	api := router.Group("/api/prescription")
	api.Use(middleware.UserContext(jwtSvc))
	{
		api.POST("/create", prescHandler.CreatePrescription)
		api.GET("/:id", prescHandler.GetPrescription)
		api.GET("/list", prescHandler.ListPrescriptions)
		api.POST("/review", prescHandler.ReviewPrescription)
		api.POST("/cancel/:id", prescHandler.CancelPrescription)
	}

	return router
}

func startGrpcServer(svc *service.PrescriptionService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9085)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcps.NewPrescriptionGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	prescription.RegisterPrescriptionServiceServer(s, grpcSrv)

	log.Printf("[Prescription] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
