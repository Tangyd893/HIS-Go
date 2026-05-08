// Package main 检查检验服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/examination/handler"
	"his-go/internal/examination/repository"
	"his-go/internal/examination/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"

	"his-go/api/proto/examination"
	grpcexam "his-go/internal/examination"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8088

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Examination 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_examination", cfg.Database.SSLMode,
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

	examRepo := repository.NewExaminationRepository(db)
	examSvc := service.NewExaminationService(examRepo)
	examHandler := handler.NewExaminationHandler(examSvc)

	router := setupExaminationRouter(cfg, examHandler)

	go startGrpcServer(examSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("检查检验服务已启动")
	log.Printf("[Examination] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupExaminationRouter(cfg *config.Config, examHandler *handler.ExaminationHandler) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "his-examination"})
	})

	api := router.Group("/api/examination")
	{
		api.POST("/report", examHandler.CreateReport)
		api.GET("/report/:id", examHandler.GetReport)
		api.GET("/reports", examHandler.ListReports)
		api.POST("/review", examHandler.ReviewReport)
	}

	return router
}

func startGrpcServer(svc *service.ExaminationService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9088)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcexam.NewExaminationGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	examination.RegisterExaminationServiceServer(s, grpcSrv)

	log.Printf("[Examination] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
