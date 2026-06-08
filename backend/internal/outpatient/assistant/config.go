// Package assistant 就诊助手（RAG + DeepSeek 分诊建议）
package assistant

import (
	"os"
	"strconv"
)

// Config 就诊助手配置（从环境变量加载）
type Config struct {
	DeepSeekAPIKey     string
	DeepSeekBaseURL    string
	DeepSeekModel      string
	SiliconFlowAPIKey  string
	SiliconFlowEmbedModel string
	RAGEnabled         bool
	TopK               int
}

// LoadConfig 从环境变量加载配置，缺失时返回安全默认值
func LoadConfig() *Config {
	cfg := &Config{
		DeepSeekAPIKey:        os.Getenv("DEEPSEEK_API_KEY"),
		DeepSeekBaseURL:       envOrDefault("DEEPSEEK_BASE_URL", "https://api.deepseek.com"),
		DeepSeekModel:         envOrDefault("DEEPSEEK_MODEL", "deepseek-chat"),
		SiliconFlowAPIKey:     os.Getenv("SILICONFLOW_API_KEY"),
		SiliconFlowEmbedModel: envOrDefault("SILICONFLOW_EMBED_MODEL", "BAAI/bge-m3"),
		RAGEnabled:            stringsToBool(os.Getenv("TRIAGE_RAG_ENABLED")),
		TopK:                  envIntOrDefault("TRIAGE_TOP_K", 5),
	}
	return cfg
}

// IsDeepSeekAvailable DeepSeek API 是否可用
func (c *Config) IsDeepSeekAvailable() bool {
	return c.DeepSeekAPIKey != "" && c.RAGEnabled
}

// IsSemanticSearchAvailable 语义检索是否可用
func (c *Config) IsSemanticSearchAvailable() bool {
	return c.SiliconFlowAPIKey != "" && c.RAGEnabled
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envIntOrDefault(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func stringsToBool(s string) bool {
	switch s {
	case "true", "TRUE", "1", "yes", "YES":
		return true
	default:
		return false
	}
}
