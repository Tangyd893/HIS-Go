// Package main 网关服务入口 — 基于 httputil.ReverseProxy 的反向代理
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"his-go/pkg/config"
	"his-go/pkg/health"
	"his-go/pkg/logger"
	"his-go/pkg/middleware"
	"his-go/pkg/nacos"
	"his-go/pkg/security/auth"
	"his-go/pkg/security/jwt"
)

var routeManager *nacos.RouteManager

var serviceRoutes = map[string]string{
	"/api/auth":          "http://localhost:8081",
	"/api/user":          "http://localhost:8082",
	"/api/registration":  "http://localhost:8083",
	"/api/clinic":        "http://localhost:8084",
	"/api/prescription":  "http://localhost:8085",
	"/api/billing":       "http://localhost:8086",
	"/api/pharmacy":      "http://localhost:8087",
	"/api/examination":   "http://localhost:8088",
	"/api/inpatient":     "http://localhost:8089",
	"/api/schedule":      "http://localhost:8090",
	"/api/outpatient":    "http://localhost:8091",
	"/api/followup":      "http://localhost:8092",
	"/api/health-record": "http://localhost:8093",
	"/api/notification":  "http://localhost:8094",
	"/api/statistics":    "http://localhost:8095",
	"/api/system":        "http://localhost:8096",
	"/api/emr":           "http://localhost:8097",
}

var dockerServiceMapping = map[string]string{
	"/api/auth":          "his-auth",
	"/api/user":          "his-user",
	"/api/registration":  "his-registration",
	"/api/clinic":        "his-clinic",
	"/api/prescription":  "his-prescription",
	"/api/billing":       "his-billing",
	"/api/pharmacy":      "his-pharmacy",
	"/api/examination":   "his-examination",
	"/api/inpatient":     "his-inpatient",
	"/api/schedule":      "his-schedule",
	"/api/outpatient":    "his-outpatient",
	"/api/followup":      "his-followup",
	"/api/health-record": "his-health-record",
	"/api/notification":  "his-notification",
	"/api/statistics":    "his-statistics",
	"/api/system":        "his-system",
	"/api/emr":           "his-emr",
}

var jwtWhitelist = []string{
	"/health",
	"/ready",
	"/metrics",
	"/api/auth/login",
	"/api/auth/refresh",
}

var jwtSvc *jwt.JWTService

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	logger.Info("HIS-Go Gateway 服务启动中...")

	if os.Getenv("RUNNING_IN_DOCKER") == "true" {
		switchToDockerRoutes()
	}

	routeManager = nacos.NewRouteManager(cfg.Nacos.ServerAddr())
	if routeManager.IsEnabled() {
		loadRoutesFromNacos(cfg)
	}

	jwtSvc = initJWT(cfg)

	router := setupRouter(cfg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("网关服务已启动")
	log.Printf("[Gateway] 服务监听地址: %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("网关服务启动失败: " + err.Error())
	}
}

func initJWT(cfg *config.Config) *jwt.JWTService {
	if cfg.JWT.PublicKey != "" {
		svc, err := jwt.NewVerifyOnlyJWTService(cfg.JWT.PublicKey)
		if err != nil {
			logger.Fatal("初始化 JWT 验证服务失败: " + err.Error())
		}
		logger.Info("JWT 验证模式 RS256")
		return svc
	}

	secret := cfg.JWT.PrivateKey
	if secret == "" {
		secret = "his-go-default-secret"
	}
	logger.Info("JWT 验证模式 HS256")
	return jwt.NewSimpleJWTService(secret, cfg.JWT.ExpireHour)
}

func switchToDockerRoutes() {
	for prefix, serviceName := range dockerServiceMapping {
		serviceRoutes[prefix] = fmt.Sprintf("http://%s:%s", serviceName, extractPort(serviceRoutes[prefix]))
	}
	logger.Info("已切换到 Docker 容器网络路由模式")
}

func extractPort(urlStr string) string {
	if idx := strings.LastIndex(urlStr, ":"); idx != -1 {
		return urlStr[idx+1:]
	}
	return "8080"
}

func setupRouter(cfg *config.Config) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.Cors())
	router.Use(middleware.RequestID())
	router.Use(middleware.Tracing())
	router.Use(middleware.Metrics("his-gateway"))

	router.GET("/health", health.HealthHandler("his-gateway"))
	router.GET("/ready", health.ReadinessHandler("his-gateway", nil))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Any("/api/*path", gatewayJwtAuth(), proxyHandler)

	return router
}

func gatewayJwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		for _, path := range jwtWhitelist {
			if reqPath == path {
				c.Next()
				return
			}
		}
		auth.JwtAuth(jwtSvc)(c)
	}
}

func proxyHandler(c *gin.Context) {
	reqPath := c.Request.URL.Path

	target := resolveTarget(reqPath)
	if target == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "未找到对应微服务路由",
			"path":    reqPath,
		})
		return
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "目标地址解析失败",
		})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = reqPath

		req.Host = targetURL.Host

		clientIP := c.ClientIP()
		if prior := req.Header.Get("X-Forwarded-For"); prior != "" {
			req.Header.Set("X-Forwarded-For", prior+", "+clientIP)
		} else {
			req.Header.Set("X-Forwarded-For", clientIP)
		}
		req.Header.Set("X-Real-IP", clientIP)
		req.Header.Set("X-Forwarded-Proto", c.Request.URL.Scheme)
		req.Header.Set("X-Forwarded-Host", c.Request.Host)

		if requestID := c.GetString("requestID"); requestID != "" {
			req.Header.Set("X-Request-ID", requestID)
		}

		if userCtx := auth.GetUserContext(c); userCtx != nil {
			req.Header.Set("X-User-ID", userCtx.UserID)
			req.Header.Set("X-Username", userCtx.Username)
			req.Header.Set("X-User-Role", userCtx.Role)
			req.Header.Set("X-User-Dept", userCtx.DeptID)
			req.Header.Set("X-User-Perms", strings.Join(userCtx.Perms, ","))
		}
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("X-Proxy-By", "his-gateway")
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("代理请求失败: " + err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(fmt.Sprintf(
			`{"code":%d,"message":"目标服务不可用","service":"%s"}`,
			http.StatusServiceUnavailable, targetURL.Host,
		)))
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func resolveTarget(path string) string {
	for prefix, target := range serviceRoutes {
		if strings.HasPrefix(path, prefix) {
			return target
		}
	}
	return ""
}

func loadRoutesFromNacos(cfg *config.Config) {
	logger.Info("已启用 Nacos 动态服务发现（网关将优先使用 Nacos 注册中心获取微服务地址）")

	for prefix, serviceName := range nacos.ServiceMap {
		host := strings.Replace(serviceName, "his-", "", 1)
		defaultTarget := serviceRoutes[prefix]
		if defaultTarget == "" {
			defaultTarget = "http://localhost:" + extractPort(serviceRoutes[prefix])
		}
		routeManager.RegisterRoute(prefix, host, []string{defaultTarget})
	}

	if err := routeManager.InitFromNacos(); err != nil {
		logger.Error("Nacos 服务发现初始化失败: " + err.Error())
		return
	}

	// 用 Nacos 路由覆盖静态路由表
	nacosRoutes := routeManager.GetAllRoutes()
	for prefix, target := range nacosRoutes {
		serviceRoutes[prefix] = target
	}
	logger.Info("Nacos 动态路由加载完成，共 " + fmt.Sprint(len(nacosRoutes)) + " 条路由")
}
