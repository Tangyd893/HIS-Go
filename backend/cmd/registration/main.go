// Package main 挂号预约服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/registration/handler"
	"his-go/internal/registration/repository"
	"his-go/internal/registration/service"
	"his-go/pkg/common/snowflake"
	"his-go/pkg/config"
	"his-go/pkg/database"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8083

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Registration 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_registration", cfg.Database.SSLMode,
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

	sf, err := snowflake.NewSnowflake(1, 1)
	if err != nil {
		logger.Fatal("雪花ID生成器初始化失败: " + err.Error())
	}

	regRepo := repository.NewRegistrationRepository(db, rdb)
	regSvc := service.NewRegistrationService(regRepo, rdb, sf)
	regHandler := handler.NewRegistrationHandler(regSvc)

	router := setupRegistrationRouter(cfg, regHandler)

	go startGrpcServer(cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("挂号预约服务已启动")
	log.Printf("[Registration] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupRegistrationRouter(cfg *config.Config, regHandler *handler.RegistrationHandler) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "his-registration"})
	})

	api := router.Group("/api/registration")
	{
		api.GET("/schedules", regHandler.ListSchedules)
		api.POST("/register", regHandler.Register)
		api.POST("/cancel/:id", regHandler.CancelRegistration)
		api.POST("/signin/:id", regHandler.SignIn)
		api.GET("/queue-status", regHandler.GetQueueStatus)
	}

	return router
}

func startGrpcServer(cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9083)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}
	log.Printf("[Registration] gRPC 服务监听地址: %s", addr)
	_ = lis
}
