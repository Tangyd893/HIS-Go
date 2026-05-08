// Package main 院外患者服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/outpatient/handler"
	"his-go/internal/outpatient/repository"
	"his-go/internal/outpatient/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"

	"his-go/api/proto/outpatient"
	grpcoutpat "his-go/internal/outpatient"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8091

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Outpatient 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_outpatient", cfg.Database.SSLMode,
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

	outpatientRepo := repository.NewOutpatientRepository(db)
	outpatientSvc := service.NewOutpatientService(outpatientRepo)
	outpatientHandler := handler.NewOutpatientHandler(outpatientSvc)

	router := setupOutpatientRouter(cfg, outpatientHandler)

	go startGrpcServer(outpatientSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("院外患者服务已启动")
	log.Printf("[Outpatient] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupOutpatientRouter(cfg *config.Config, outpatientHandler *handler.OutpatientHandler) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "his-outpatient"})
	})

	api := router.Group("/api/outpatient")
	{
		api.POST("/consultation", outpatientHandler.CreateConsultation)
		api.GET("/consultation/:id", outpatientHandler.GetConsultation)
		api.GET("/consultations", outpatientHandler.ListConsultations)
		api.POST("/message", outpatientHandler.SendMessage)
		api.GET("/messages", outpatientHandler.GetMessages)
		api.POST("/contract", outpatientHandler.CreateChronicContract)
		api.POST("/health-data", outpatientHandler.ReportHealthData)
		api.GET("/health-data", outpatientHandler.ListHealthData)
	}

	return router
}

func startGrpcServer(svc *service.OutpatientService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9091)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcoutpat.NewOutpatientGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	outpatient.RegisterOutpatientServiceServer(s, grpcSrv)

	log.Printf("[Outpatient] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
