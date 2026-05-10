// Package service Nacos 服务发现集成
// 为 Gateway 提供基于 Nacos 的动态路由解析能力
package nacos

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"
)

// 注：生产环境需引入 Nacos Go SDK
// import "github.com/nacos-group/nacos-sdk-go/v2/clients"
// import "github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
// import "github.com/nacos-group/nacos-sdk-go/v2/vo"

// ServiceInstance Nacos 服务实例信息
type ServiceInstance struct {
	ServiceName string  `json:"serviceName"`
	IP          string  `json:"ip"`
	Port        int     `json:"port"`
	Healthy     bool    `json:"healthy"`
	Weight      float64 `json:"weight"`
	ClusterName string  `json:"clusterName"`
}

// ServiceRoute 网关路由条目
type ServiceRoute struct {
	PathPrefix   string   `json:"pathPrefix"`
	ServiceName  string   `json:"serviceName"`
	TargetURLs   []string `json:"targetUrls"`
	ActiveTarget string   `json:"activeTarget"`
	LoadBalance  string   `json:"loadBalance"` // round-robin, random, least-conn
}

// DiscoveryClient Nacos 服务发现客户端接口
type DiscoveryClient interface {
	// GetService 获取服务实例列表
	GetService(serviceName string) ([]ServiceInstance, error)
	// Subscribe 订阅服务变更
	Subscribe(serviceName string, callback func([]ServiceInstance)) error
	// Close 关闭客户端
	Close()
}

// RouteManager 网关路由管理器
type RouteManager struct {
	mu       sync.RWMutex
	routes   map[string]*ServiceRoute // pathPrefix → route
	client   DiscoveryClient
	nacosURL string
	enabled  bool
}

// NewRouteManager 创建路由管理器
func NewRouteManager(nacosURL string) *RouteManager {
	enabled := os.Getenv("USE_NACOS") == "true" || nacosURL != ""

	rm := &RouteManager{
		routes:   make(map[string]*ServiceRoute),
		nacosURL: nacosURL,
		enabled:  enabled,
	}

	return rm
}

// IsEnabled 检查 Nacos 服务发现是否启用
func (rm *RouteManager) IsEnabled() bool {
	return rm.enabled
}

// RegisterRoute 注册路由映射
func (rm *RouteManager) RegisterRoute(pathPrefix, serviceName string, defaultTargets []string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.routes[pathPrefix] = &ServiceRoute{
		PathPrefix:   pathPrefix,
		ServiceName:  serviceName,
		TargetURLs:   defaultTargets,
		ActiveTarget: defaultTargets[0],
		LoadBalance:  "round-robin",
	}
}

// GetAllRoutes 获取所有路由
func (rm *RouteManager) GetAllRoutes() map[string]string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	result := make(map[string]string, len(rm.routes))
	for prefix, route := range rm.routes {
		result[prefix] = route.ActiveTarget
	}
	return result
}

// GetTarget 根据路径前缀获取目标服务地址
func (rm *RouteManager) GetTarget(pathPrefix string) (string, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	route, ok := rm.routes[pathPrefix]
	if !ok {
		return "", false
	}
	return route.ActiveTarget, true
}

// InitFromNacos 从 Nacos 注册中心初始化路由表
func (rm *RouteManager) InitFromNacos() error {
	if !rm.enabled {
		return nil
	}

	log.Printf("[Nacos] 从注册中心加载服务列表: %s", rm.nacosURL)

	// 使用静态映射表作为后备
	// 生产环境替换为 Nacos Go SDK 真实调用:
	//   namingClient, err := clients.NewNamingClient(vo.NacosClientParam{ServerConfigs: ...})
	//   services, err := namingClient.GetService(vo.GetServiceParam{ServiceName: "his-gateway"})

	for _, route := range rm.routes {
		// 用服务名作为 Docker 容器名回退
		if len(route.TargetURLs) == 0 {
			target := fmt.Sprintf("http://%s:8080", route.ServiceName)
			route.TargetURLs = []string{target}
			route.ActiveTarget = target
		}
		log.Printf("[Nacos] 路由 %s → %s (%s)", route.PathPrefix, route.ServiceName, route.ActiveTarget)
	}

	return nil
}

// GetTargetURL 解析并验证目标 URL
func GetTargetURL(target string) (*url.URL, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("无效的目标 URL %s: %w", target, err)
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return u, nil
}

// DiscoverLoop 后台循环从 Nacos 拉取服务实例（30秒间隔）
func (rm *RouteManager) DiscoverLoop(interval time.Duration) {
	if !rm.enabled || rm.client == nil {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, route := range rm.routes {
			instances, err := rm.client.GetService(route.ServiceName)
			if err != nil {
				log.Printf("[Nacos] 获取服务 %s 实例失败: %v", route.ServiceName, err)
				continue
			}

			healthyCount := 0
			for _, inst := range instances {
				if inst.Healthy {
					healthyCount++
				}
			}
			log.Printf("[Nacos] 服务 %s: %d/%d 健康实例", route.ServiceName, healthyCount, len(instances))
		}
	}
}

// ServiceMap 17 个微服务的路径前缀到服务名映射
var ServiceMap = map[string]string{
	"/api/auth":          "his-auth",
	"/api/user":          "his-user",
	"/api/registration":  "his-registration",
	"/api/clinic":        "his-clinic",
	"/api/emr":           "his-emr",
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
}
