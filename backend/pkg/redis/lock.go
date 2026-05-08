// Package redis 分布式锁实现
package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Lock 获取分布式锁
// key: 锁的键
// ttl: 锁超时时间
// 返回锁标识（用于解锁时验证）和错误
func (c *Client) Lock(ctx context.Context, key string, ttl time.Duration) (string, error) {
	lockValue := uuid.New().String()
	ok, err := c.SetNX(ctx, key, lockValue, ttl)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil // 获取锁失败
	}
	return lockValue, nil
}

// Unlock 释放分布式锁（使用 Lua 脚本保证原子性）
func (c *Client) Unlock(ctx context.Context, key, lockValue string) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	result, err := c.Eval(ctx, script, []string{key}, lockValue).Int64()
	if err != nil {
		return err
	}
	if result == 0 {
		return nil // 锁已被其他持有者释放
	}
	return nil
}
