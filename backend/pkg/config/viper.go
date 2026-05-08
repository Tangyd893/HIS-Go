// Package config 配置管理（基于 Viper）
package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Nacos    NacosConfig    `mapstructure:"nacos"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
	Log      LogConfig      `mapstructure:"log"`
	Grpc     GrpcConfig     `mapstructure:"grpc"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

type JWTConfig struct {
	PrivateKey string `mapstructure:"private_key"`
	PublicKey  string `mapstructure:"public_key"`
	ExpireHour int    `mapstructure:"expire_hour"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Group     string `mapstructure:"group"`
}

type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json 或 console
}

type GrpcConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

var GlobalConfig *Config

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 支持环境变量覆盖
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 环境变量覆盖敏感配置
	overrideFromEnv(&cfg)

	GlobalConfig = &cfg
	return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		cfg.Database.Port = parseIntEnv(v)
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		cfg.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PORT"); v != "" {
		cfg.Redis.Port = parseIntEnv(v)
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		cfg.Redis.Password = v
	}
	if v := os.Getenv("MQ_HOST"); v != "" {
		cfg.RabbitMQ.Host = v
	}
	if v := os.Getenv("MQ_PORT"); v != "" {
		cfg.RabbitMQ.Port = parseIntEnv(v)
	}
	if v := os.Getenv("MQ_PASSWORD"); v != "" {
		cfg.RabbitMQ.Password = v
	}
	if v := os.Getenv("NACOS_HOST"); v != "" {
		cfg.Nacos.Host = v
	}
	if v := os.Getenv("NACOS_PORT"); v != "" {
		cfg.Nacos.Port = parseIntEnv(v)
	}
	if v := os.Getenv("MINIO_ENDPOINT"); v != "" {
		cfg.MinIO.Endpoint = v
	}
	if v := os.Getenv("MINIO_ACCESS_KEY"); v != "" {
		cfg.MinIO.AccessKey = v
	}
	if v := os.Getenv("MINIO_SECRET_KEY"); v != "" {
		cfg.MinIO.SecretKey = v
	}
	if v := os.Getenv("JWT_PRIVATE_KEY"); v != "" {
		cfg.JWT.PrivateKey = v
	}
	if v := os.Getenv("JWT_PUBLIC_KEY"); v != "" {
		cfg.JWT.PublicKey = v
	}
}

func parseIntEnv(v string) int {
	var result int
	fmt.Sscanf(v, "%d", &result)
	return result
}
