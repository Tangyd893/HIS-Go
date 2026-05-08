// Package main 排班管理服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/schedule/handler"
	"his-go/internal/schedule/repository"
	"his-go/internal/schedule/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"

	"his-go/api/proto/schedule"
	grpcsched "his-go/internal/schedule"
	hisgrpc "his-go/pkg/grpc"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8090

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Schedule 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_schedule", cfg.Database.SSLMode,
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

	scheduleRepo := repository.NewScheduleRepository(db)
	scheduleSvc := service.NewScheduleService(scheduleRepo)
	scheduleHandler := handler.NewScheduleHandler(scheduleSvc)

	router := setupScheduleRouter(cfg, scheduleHandler)

	go startGrpcServer(scheduleSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("排班管理服务已启动")
	log.Printf("[Schedule] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupScheduleRouter(cfg *config.Config, scheduleHandler *handler.ScheduleHandler) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "his-schedule"})
	})

	api := router.Group("/api/schedule")
	{
		api.POST("/generate", scheduleHandler.GenerateWeeklySchedules)
		api.GET("/list", scheduleHandler.ListSchedules)
		api.PUT("/update", scheduleHandler.UpdateSchedule)
		api.POST("/cancel/:id", scheduleHandler.CancelSchedule)
	}

	return router
}

func startGrpcServer(svc *service.ScheduleService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9090)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcsched.NewScheduleGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	schedule.RegisterScheduleServiceServer(s, grpcSrv)

	log.Printf("[Schedule] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
