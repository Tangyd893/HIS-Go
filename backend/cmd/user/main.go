// Package main 用户管理服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"his-go/api/proto/user"
	grpcuser "his-go/internal/user"
	"his-go/internal/user/handler"
	"his-go/internal/user/repository"
	"his-go/internal/user/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	hisgrpc "his-go/pkg/grpc"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8082

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-User 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_user", cfg.Database.SSLMode,
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

	patientRepo := repository.NewPatientRepository(db)
	empRepo := repository.NewEmployeeRepository(db)
	deptRepo := repository.NewDepartmentRepository(db)

	userSvc := service.NewUserService(patientRepo, empRepo, deptRepo, rdb)
	userHandler := handler.NewUserHandler(userSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupUserRouter(cfg, userHandler, deps)

	go startUserGrpcServer(userSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("用户服务已启动")
	log.Printf("[User] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupUserRouter(cfg *config.Config, userHandler *handler.UserHandler, deps *health.Dependencies) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-user"))

	router.GET("/health", health.HealthHandler("his-user"))
	router.GET("/ready", health.ReadinessHandler("his-user", deps))

	api := router.Group("/api/user")
	{
		api.GET("/patients", userHandler.ListPatients)
		api.GET("/departments", userHandler.ListDepartments)
		api.GET("/employees", userHandler.ListEmployees)
	}

	return router
}

func startUserGrpcServer(userSvc *service.UserService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9082)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcSrv := grpcuser.NewUserGrpcServer(userSvc)
	s := hisgrpc.NewGrpcServer()
	user.RegisterUserServiceServer(s, grpcSrv)

	log.Printf("[User] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
