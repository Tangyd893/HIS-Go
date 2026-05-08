// Package middleware 限流中间件
package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	apperrors "his-go/pkg/errors"
	"his-go/pkg/response"
)

// RateLimiter 滑动窗口限流中间件
// limit: 时间窗口内最大请求数
// window: 时间窗口大小
func RateLimiter(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ratelimit:" + c.FullPath() + ":" + c.ClientIP()
		ctx := context.Background()

		now := time.Now().UnixMilli()
		windowStart := now - window.Milliseconds()

		pipe := rdb.Pipeline()
		pipe.ZRemRangeByScore(ctx, key, "0", formatTime(windowStart))
		pipe.ZCard(ctx, key)
		cmds, err := pipe.Exec(ctx)
		if err != nil {
			c.Next()
			return
		}

		count := cmds[1].(*redis.IntCmd).Val()
		if count >= int64(limit) {
			response.Fail(c, apperrors.CodeRateLimited)
			c.Abort()
			return
		}

		rdb.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
		rdb.Expire(ctx, key, window)
		c.Next()
	}
}

func formatTime(t int64) string {
	return time.UnixMilli(t).Format(time.RFC3339Nano)
}
