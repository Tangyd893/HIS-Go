// Package main 住院管理服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/inpatient/handler"
	"his-go/internal/inpatient/repository"
	"his-go/internal/inpatient/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"
	"his-go/pkg/security/jwt"

	"his-go/api/proto/inpatient"
	grpcinpat "his-go/internal/inpatient"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8089

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Inpatient 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_inpatient", cfg.Database.SSLMode,
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

	inpatientRepo := repository.NewInpatientRepository(db)
	inpatientSvc := service.NewInpatientService(inpatientRepo)
	inpatientHandler := handler.NewInpatientHandler(inpatientSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupInpatientRouter(cfg, inpatientHandler, deps, jwtSvc)

	go startGrpcServer(inpatientSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("住院管理服务已启动")
	log.Printf("[Inpatient] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupInpatientRouter(cfg *config.Config, inpatientHandler *handler.InpatientHandler, deps *health.Dependencies, jwtSvc *jwt.JWTService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-inpatient"))

	router.GET("/health", health.HealthHandler("his-inpatient"))
	router.GET("/ready", health.ReadinessHandler("his-inpatient", deps))

	api := router.Group("/api/inpatient")
	api.Use(middleware.UserContext(jwtSvc))
	{
		api.POST("/admit", inpatientHandler.AdmitPatient)
		api.POST("/discharge/:id", inpatientHandler.DischargePatient)
		api.GET("/:id", inpatientHandler.GetInpatient)
		api.GET("/list", inpatientHandler.ListInpatients)
		api.POST("/order", inpatientHandler.CreateMedicalOrder)
		api.POST("/nursing", inpatientHandler.CreateNursingRecord)
	}

	return router
}

func startGrpcServer(svc *service.InpatientService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9089)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcinpat.NewInpatientGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	inpatient.RegisterInpatientServiceServer(s, grpcSrv)

	log.Printf("[Inpatient] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
