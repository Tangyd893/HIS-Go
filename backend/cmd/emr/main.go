// Package main 电子病历服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/emr/handler"
	"his-go/internal/emr/repository"
	"his-go/internal/emr/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"
	"his-go/pkg/security/jwt"

	"his-go/api/proto/emr"
	grpcemr "his-go/internal/emr"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8097

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-EMR 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_emr", cfg.Database.SSLMode,
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

	emrRepo := repository.NewEMRRepository(db)
	emrSvc := service.NewEMRService(emrRepo)
	emrHandler := handler.NewEMRHandler(emrSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupEMRRouter(cfg, emrHandler, deps, jwtSvc)

	go startGrpcServer(emrSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("电子病历服务已启动")
	log.Printf("[EMR] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupEMRRouter(cfg *config.Config, emrHandler *handler.EMRHandler, deps *health.Dependencies, jwtSvc *jwt.JWTService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-emr"))

	router.GET("/health", health.HealthHandler("his-emr"))
	router.GET("/ready", health.ReadinessHandler("his-emr", deps))

	api := router.Group("/api/emr")
	api.Use(middleware.UserContext(jwtSvc))
	{
		api.POST("/record", emrHandler.CreateRecord)
		api.GET("/record/:id", emrHandler.GetRecord)
		api.GET("/records", emrHandler.ListRecords)
		api.POST("/quality-control/:id", emrHandler.QualityControl)
		api.GET("/templates", emrHandler.ListTemplates)
		api.POST("/cdss-check", emrHandler.CDSSCheck)
	}

	return router
}

func startGrpcServer(svc *service.EMRService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9097)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcemr.NewEMRGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	emr.RegisterEMRServiceServer(s, grpcSrv)

	log.Printf("[EMR] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
