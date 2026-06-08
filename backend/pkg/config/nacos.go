// Package config Nacos 服务注册与配置中心辅助函数
package config

import (
	"fmt"
	"os"
)

// DefaultNacosConfig 返回 Nacos 默认配置
func DefaultNacosConfig() NacosConfig {
	return NacosConfig{
		Host:      getNacosEnv("NACOS_HOST", "127.0.0.1"),
		Port:      8848,
		Namespace: getNacosEnv("NACOS_NAMESPACE", "public"),
		Group:     "DEFAULT_GROUP",
	}
}

// ServerAddr 返回 Nacos 服务端地址
func (c NacosConfig) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsConfigured 检查是否已配置 Nacos
func (c NacosConfig) IsConfigured() bool {
	return c.Host != "" && c.Port > 0
}

// UseNacos 判断是否应启用 Nacos
func UseNacos() bool {
	return os.Getenv("USE_NACOS") == "true" || os.Getenv("NACOS_HOST") != ""
}

func getNacosEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
