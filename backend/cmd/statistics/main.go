// Package main 数据统计服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/statistics/handler"
	"his-go/internal/statistics/repository"
	"his-go/internal/statistics/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"

	"his-go/api/proto/statistics"
	grpcstats "his-go/internal/statistics"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8095

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Statistics 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_statistics", cfg.Database.SSLMode,
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

	statsRepo := repository.NewStatisticsRepository(db)
	statsSvc := service.NewStatisticsService(statsRepo)
	statsHandler := handler.NewStatisticsHandler(statsSvc)

	router := setupStatisticsRouter(cfg, statsHandler)

	go startGrpcServer(statsSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("数据统计服务已启动")
	log.Printf("[Statistics] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupStatisticsRouter(cfg *config.Config, statsHandler *handler.StatisticsHandler) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "his-statistics"})
	})

	api := router.Group("/api/statistics")
	{
		api.POST("/operation", statsHandler.GetOperationStats)
		api.POST("/dept-workload", statsHandler.GetDeptWorkload)
		api.POST("/revenue-trend", statsHandler.GetRevenueTrend)
	}

	return router
}

func startGrpcServer(svc *service.StatisticsService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9095)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcstats.NewStatisticsGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	statistics.RegisterStatisticsServiceServer(s, grpcSrv)

	log.Printf("[Statistics] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
