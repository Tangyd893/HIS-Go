// Package main 认证服务入口
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"his-go/api/proto/auth"
	grpcauth "his-go/internal/auth"
	"his-go/internal/auth/handler"
	"his-go/internal/auth/repository"
	"his-go/internal/auth/service"
	"his-go/pkg/config"
	"his-go/pkg/database"
	hisgrpc "his-go/pkg/grpc"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/redis"
	secauth "his-go/pkg/security/auth"
	"his-go/pkg/security/jwt"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg.Server.Port = 8081

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Auth 服务启动中...")

	db, err := database.NewPostgres(
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		"his_auth", cfg.Database.SSLMode,
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

	var jwtSvc *jwt.JWTService
	if os.Getenv("USE_JWT_SIMPLE") == "true" || cfg.JWT.PrivateKey == "" {
		secret := cfg.JWT.PrivateKey
		if secret == "" {
			secret = "his-go-default-secret"
		}
		jwtSvc = jwt.NewSimpleJWTService(secret, cfg.JWT.ExpireHour)
		logger.Info("JWT 使用简化模式 HS256")
	} else {
		jwtSvc, err = jwt.NewJWTService(cfg.JWT.PrivateKey, cfg.JWT.PublicKey, cfg.JWT.ExpireHour)
		if err != nil {
			logger.Fatal("初始化 JWT 服务失败: " + err.Error())
		}
		logger.Info("JWT 使用非对称模式 RS256")
	}

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo, jwtSvc, rdb)
	authHandler := handler.NewAuthHandler(authSvc)

	sqlDB, _ := db.DB()
	deps := &health.Dependencies{
		DB:    sqlDB,
		Redis: rdb,
	}

	router := setupAuthRouter(cfg, authHandler, deps, jwtSvc)

	go startGrpcServer(authSvc, cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("认证服务已启动")
	log.Printf("[Auth] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("服务启动失败: " + err.Error())
	}
}

func setupAuthRouter(cfg *config.Config, authHandler *handler.AuthHandler, deps *health.Dependencies, jwtSvc *jwt.JWTService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logger(), middleware.Cors(), middleware.RequestID())
	router.Use(middleware.Metrics("his-auth"))

	router.GET("/health", health.HealthHandler("his-auth"))
	router.GET("/ready", health.ReadinessHandler("his-auth", deps))

	api := router.Group("/api/auth")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.RefreshToken)

		protected := api.Group("")
		protected.Use(secauth.UserContextFromGatewayOrJWT(jwtSvc))
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/current", authHandler.Current)
	}

	return router
}

func startGrpcServer(authSvc *service.AuthService, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Grpc.Host, 9081)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("gRPC 监听失败: " + err.Error())
		return
	}

	grpcServer := grpcauth.NewAuthGrpcServer(authSvc)
	s := hisgrpc.NewGrpcServer()
	auth.RegisterAuthServiceServer(s, grpcServer)

	log.Printf("[Auth] gRPC 服务监听地址: %s", addr)
	if err := s.Serve(lis); err != nil {
		logger.Error("gRPC 服务启动失败: " + err.Error())
	}
}
