// Package main 系统管理服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/internal/system/handler"
	"his-go/internal/system/repository"
	"his-go/internal/system/service"
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

	cfg.Server.Port = 8096

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-System 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_system", cfg.Database.SSLMode,
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

	sysRepo := repository.NewSystemRepository(db)
	sysSvc := service.NewSystemService(sysRepo)
	sysHandler := handler.NewSystemHandler(sysSvc)

	router := setupSystemRouter(cfg, sysHandler)

	go startGrpcServer(cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("系统管理服务已启动")
	log.Printf("[System] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupSystemRouter(cfg *config.Config, sysHandler *handler.SystemHandler) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "his-system"})
	})

	api := router.Group("/api/system")
	{
		api.GET("/dict-types", sysHandler.ListDictTypes)
		api.GET("/dict-items", sysHandler.ListDictItems)
		api.POST("/dict-item", sysHandler.CreateDictItem)
		api.GET("/params", sysHandler.ListParams)
		api.PUT("/param", sysHandler.UpdateParam)
		api.GET("/operation-logs", sysHandler.ListOperationLogs)
	}

	return router
}

func startGrpcServer(cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9096)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}
	log.Printf("[System] gRPC 服务监听地址: %s", addr)
	_ = lis
}
