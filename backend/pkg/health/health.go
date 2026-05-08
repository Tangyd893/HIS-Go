package health

import (
	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"

	"his-go/pkg/logger"
	"his-go/pkg/redis"
)

const checkTimeout = 3 * time.Second

type Dependencies struct {
	DB    *sql.DB
	Redis *redis.Client
}

type Status struct {
	Status  string            `json:"status"`
	Service string            `json:"service"`
	Details map[string]string `json:"details,omitempty"`
}

func HealthHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, Status{
			Status:  "UP",
			Service: serviceName,
		})
	}
}

func ReadinessHandler(serviceName string, deps *Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthy := true
		details := make(map[string]string)

		if deps == nil {
			c.JSON(200, Status{
				Status:  "UP",
				Service: serviceName,
				Details: map[string]string{"note": "无外部依赖"},
			})
			return
		}

		if deps.DB != nil {
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
			defer cancel()
			if err := deps.DB.PingContext(ctx); err != nil {
				healthy = false
				details["database"] = "DOWN: " + err.Error()
			} else {
				details["database"] = "UP"
			}
			logger.Debug("健康检查 DB 耗时: " + time.Since(start).String())
		}

		if deps.Redis != nil {
			ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
			defer cancel()
			if err := deps.Redis.Ping(ctx).Err(); err != nil {
				healthy = false
				details["redis"] = "DOWN: " + err.Error()
			} else {
				details["redis"] = "UP"
			}
		}

		status := "UP"
		if !healthy {
			status = "DOWN"
		}

		c.JSON(200, Status{
			Status:  status,
			Service: serviceName,
			Details: details,
		})
	}
}
