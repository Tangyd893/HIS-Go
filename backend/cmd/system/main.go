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
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"
	"his-go/pkg/security/jwt"

	"his-go/api/proto/system"
	grpcsys "his-go/internal/system"
	hisgrpc "his-go/pkg/grpc"
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

	jwtSvc := middleware.InitJWT(cfg)

	sysRepo := repository.NewSystemRepository(db)
	sysSvc := service.NewSystemService(sysRepo)
	sysHandler := handler.NewSystemHandler(sysSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupSystemRouter(cfg, sysHandler, deps, jwtSvc)

	go startGrpcServer(sysSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("系统管理服务已启动")
	log.Printf("[System] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupSystemRouter(cfg *config.Config, sysHandler *handler.SystemHandler, deps *health.Dependencies, jwtSvc *jwt.JWTService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-system"))

	router.GET("/health", health.HealthHandler("his-system"))
	router.GET("/ready", health.ReadinessHandler("his-system", deps))

	api := router.Group("/api/system")
	api.Use(middleware.UserContext(jwtSvc))
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

func startGrpcServer(svc *service.SystemService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9096)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcsys.NewSystemGrpcServer(svc)
	s := hisgrpc.NewGrpcServer()
	system.RegisterSystemServiceServer(s, grpcSrv)

	log.Printf("[System] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
