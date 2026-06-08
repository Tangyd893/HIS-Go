// Package redis 布隆过滤器扩展，用于防止缓存穿透
package redis

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"hash"

	"github.com/redis/go-redis/v9"
)

const (
	// bloomFilterScript 布隆过滤器 Lua 脚本
	bloomAddScript = `
local key = KEYS[1]
local bits = ARGV
for i = 1, #bits do
    redis.call('SETBIT', key, bits[i], 1)
end
return 1
`

	bloomExistsScript = `
local key = KEYS[1]
local bits = ARGV
for i = 1, #bits do
    if redis.call('GETBIT', key, bits[i]) == 0 then
        return 0
    end
end
return 1
`
)

// BloomFilter 布隆过滤器（基于 Redis Bitmap）
type BloomFilter struct {
	client    *Client
	key       string
	hashFuncs int
	size      uint32
}

// NewBloomFilter 创建布隆过滤器
// key: Redis key 前缀
// expectedItems: 预估元素数量
// falsePositiveRate: 目标误判率 (0 < rate < 1)
func (c *Client) NewBloomFilter(key string, expectedItems int, falsePositiveRate float64) *BloomFilter {
	size, hashFuncs := estimateBloomParams(uint32(expectedItems), falsePositiveRate)
	return &BloomFilter{
		client:    c,
		key:       "bf:" + key,
		hashFuncs: hashFuncs,
		size:      size,
	}
}

// Add 添加元素到布隆过滤器
func (bf *BloomFilter) Add(ctx context.Context, item string) error {
	positions := bf.hashPositions(item)
	args := make([]interface{}, len(positions))
	for i, pos := range positions {
		args[i] = pos
	}

	_, err := bf.client.Eval(ctx, bloomAddScript, []string{bf.key}, args...).Result()
	return err
}

// Exists 检查元素是否可能存在
// 返回 true 表示可能存在（可能有误判），false 表示一定不存在
func (bf *BloomFilter) Exists(ctx context.Context, item string) (bool, error) {
	positions := bf.hashPositions(item)
	args := make([]interface{}, len(positions))
	for i, pos := range positions {
		args[i] = pos
	}

	result, err := bf.client.Eval(ctx, bloomExistsScript, []string{bf.key}, args...).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// Clear 清空布隆过滤器
func (bf *BloomFilter) Clear(ctx context.Context) error {
	return bf.client.Del(ctx, bf.key)
}

// hashPositions 计算元素在布隆过滤器中的多个位位置
func (bf *BloomFilter) hashPositions(item string) []uint32 {
	positions := make([]uint32, bf.hashFuncs)

	hasher := md5.New()
	hasher.Write([]byte(item))
	hashBytes := hasher.Sum(nil)

	h1 := binary.BigEndian.Uint32(hashBytes[0:4])
	h2 := binary.BigEndian.Uint32(hashBytes[4:8])

	for i := 0; i < bf.hashFuncs; i++ {
		positions[i] = (h1 + uint32(i)*h2) % bf.size
	}

	return positions
}

// NewBloomFilterClient 为独立的 go-redis Client 创建布隆过滤器
func NewBloomFilterClient(client *redis.Client, key string, expectedItems int, falsePositiveRate float64) *BloomFilter {
	size, hashFuncs := estimateBloomParams(uint32(expectedItems), falsePositiveRate)
	wrappedClient := &Client{Client: client}
	return &BloomFilter{
		client:    wrappedClient,
		key:       "bf:" + key,
		hashFuncs: hashFuncs,
		size:      size,
	}
}

// estimateBloomParams 根据期望元素数量和误判率估算位数组大小和哈希函数数量
func estimateBloomParams(n uint32, p float64) (uint32, int) {
	if p <= 0 {
		p = 0.01
	}
	if p >= 1 {
		p = 0.99
	}

	m := uint32(float64(-int64(n)) * p / (0.693 * 0.693))
	if m < 1000 {
		m = 1000
	}

	k := int(float64(m) / float64(n) * 0.693)
	if k < 1 {
		k = 1
	}
	if k > 30 {
		k = 30
	}

	return m, k
}

// DoubleHash 双重哈希辅助函数
func doubleHash(h1, h2 uint32, i int, m uint32) uint32 {
	return (h1 + uint32(i)*h2) % m
}

// HashFnv 使用 FNV-1a 哈希算法
func hashFnv(data []byte) (uint32, uint32) {
	var h1, h2 uint32 = 2166136261, 2166136261

	for i, b := range data {
		h1 ^= uint32(b)
		h1 *= 16777619
		h2 ^= uint32(data[len(data)-1-i])
		h2 *= 16777619
	}

	return h1, h2
}

// 确保 hash 包被引用（用于 crypto/md5 的 Hash 接口）
var _ hash.Hash = md5.New()
