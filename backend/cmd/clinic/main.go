// Package main 门诊诊疗服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/clinic/handler"
	"his-go/internal/clinic/repository"
	"his-go/internal/clinic/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"

	"his-go/api/proto/clinic"
	grpcclinic "his-go/internal/clinic"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8084

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Clinic 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_clinic", cfg.Database.SSLMode,
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

	clinicRepo := repository.NewClinicRepository(db)
	clinicSvc := service.NewClinicService(clinicRepo)
	clinicHandler := handler.NewClinicHandler(clinicSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupClinicRouter(cfg, clinicHandler, deps)

	go startGrpcServer(clinicSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("门诊诊疗服务已启动")
	log.Printf("[Clinic] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupClinicRouter(cfg *config.Config, clinicHandler *handler.ClinicHandler, deps *health.Dependencies) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-clinic"))

	router.GET("/health", health.HealthHandler("his-clinic"))
	router.GET("/ready", health.ReadinessHandler("his-clinic", deps))

	api := router.Group("/api/clinic")
	{
		api.POST("/record", clinicHandler.CreateClinicRecord)
		api.GET("/record/:id", clinicHandler.GetClinicRecord)
		api.GET("/records", clinicHandler.ListClinicRecords)
		api.POST("/examination-request", clinicHandler.CreateExaminationRequest)
	}

	return router
}

func startGrpcServer(svc *service.ClinicService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9084)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcclinic.NewClinicGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	clinic.RegisterClinicServiceServer(s, grpcSrv)

	log.Printf("[Clinic] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
