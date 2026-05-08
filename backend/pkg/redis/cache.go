// Package redis Redis 缓存基础封装
package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	*goredis.Client
}

func NewClient(host string, port int, password string, db int, poolSize int) (*Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	return &Client{Client: client}, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}

func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.Client.Exists(ctx, keys...).Result()
}

func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.Client.SetNX(ctx, key, value, expiration).Result()
}

func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.Client.HSet(ctx, key, values...).Err()
}

func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return c.Client.HGet(ctx, key, field).Result()
}

func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.Client.HGetAll(ctx, key).Result()
}

func (c *Client) ZAdd(ctx context.Context, key string, members ...goredis.Z) error {
	return c.Client.ZAdd(ctx, key, members...).Err()
}

func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.Client.ZRange(ctx, key, start, stop).Result()
}

func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.Client.ZRem(ctx, key, members...).Err()
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.Client.Incr(ctx, key).Result()
}

func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.Client.Decr(ctx, key).Result()
}
