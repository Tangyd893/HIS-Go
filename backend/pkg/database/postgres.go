// Package database 数据库连接管理
package database

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewPostgres 创建 PostgreSQL 连接
func NewPostgres(host string, port int, user, password, dbname, sslMode string, maxIdleConns, maxOpenConns int, connMaxLifetime int) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		host, port, url.QueryEscape(user), url.QueryEscape(password), url.QueryEscape(dbname), sslMode,
	)

	logLevel := gormlogger.Warn
	if os.Getenv("GORM_LOG_LEVEL") == "info" {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库 Ping 失败: %w", err)
	}

	log.Println("[DB] PostgreSQL 连接成功:", dbname)
	return db, nil
}
