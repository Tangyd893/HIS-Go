package redis

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
)

func newTestClient(t *testing.T) (*Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := &Client{
		Client: goredis.NewClient(&goredis.Options{
			Addr: mr.Addr(),
		}),
	}
	return client, mr
}

func TestLock_Success(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	lockValue, err := client.Lock(ctx, "test:lock:1", 10*time.Second)
	if err != nil {
		t.Fatalf("获取锁失败: %v", err)
	}
	if lockValue == "" {
		t.Error("期望非空 lockValue")
	}
}

func TestLock_Conflict(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	lockValue1, err := client.Lock(ctx, "test:lock:conflict", 10*time.Second)
	if err != nil {
		t.Fatalf("首次获取锁失败: %v", err)
	}
	if lockValue1 == "" {
		t.Fatal("首次获取锁返回空值")
	}

	lockValue2, err := client.Lock(ctx, "test:lock:conflict", 10*time.Second)
	if err != nil {
		t.Fatalf("第二次获取锁不应报错: %v", err)
	}
	if lockValue2 != "" {
		t.Errorf("期望第二次获取锁返回空值，实际=%s", lockValue2)
	}
}

func TestUnlock_Success(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	lockValue, err := client.Lock(ctx, "test:lock:unlock1", 10*time.Second)
	if err != nil {
		t.Fatalf("获取锁失败: %v", err)
	}

	err = client.Unlock(ctx, "test:lock:unlock1", lockValue)
	if err != nil {
		t.Fatalf("释放锁失败: %v", err)
	}

	// 释放后应可重新获取
	newValue, err := client.Lock(ctx, "test:lock:unlock1", 1*time.Second)
	if err != nil {
		t.Fatalf("重新获取锁失败: %v", err)
	}
	if newValue == "" {
		t.Error("释放后应可重新获取锁")
	}
}

func TestUnlock_WrongToken(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	lockValue, err := client.Lock(ctx, "test:lock:wrong-token", 10*time.Second)
	if err != nil {
		t.Fatalf("获取锁失败: %v", err)
	}

	// 使用错误的 token 释放锁（不应删除 key）
	err = client.Unlock(ctx, "test:lock:wrong-token", "wrong-token-value")
	if err != nil {
		t.Fatalf("错误 token 释放锁不应报错: %v", err)
	}

	// 验证锁仍存在
	val, err := client.Get(ctx, "test:lock:wrong-token")
	if err != nil {
		t.Fatalf("获取锁值失败: %v", err)
	}
	if val != lockValue {
		t.Errorf("期望锁值=%s，实际=%s", lockValue, val)
	}
}

func TestLock_TTL(t *testing.T) {
	client, mr := newTestClient(t)
	ctx := context.Background()

	lockValue, err := client.Lock(ctx, "test:lock:ttl", 1*time.Second)
	if err != nil {
		t.Fatalf("获取锁失败: %v", err)
	}
	if lockValue == "" {
		t.Fatal("获取锁返回空值")
	}

	// 快进 2 秒，锁应过期
	mr.FastForward(2 * time.Second)

	_, err = client.Get(ctx, "test:lock:ttl")
	if err == nil {
		t.Error("期望锁已过期，key 不存在")
	}
}

func TestUnlock_AlreadyExpired(t *testing.T) {
	client, mr := newTestClient(t)
	ctx := context.Background()

	lockValue, err := client.Lock(ctx, "test:lock:expire-unlock", 1*time.Second)
	if err != nil {
		t.Fatalf("获取锁失败: %v", err)
	}

	mr.FastForward(2 * time.Second)

	// 锁已过期，释放应不报错
	err = client.Unlock(ctx, "test:lock:expire-unlock", lockValue)
	if err != nil {
		t.Fatalf("已过期锁释放不应报错: %v", err)
	}
}

func TestLock_DifferentKeys(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	v1, err := client.Lock(ctx, "test:lock:key-a", 10*time.Second)
	if err != nil || v1 == "" {
		t.Fatalf("获取锁 key-a 失败: %v", err)
	}

	v2, err := client.Lock(ctx, "test:lock:key-b", 10*time.Second)
	if err != nil || v2 == "" {
		t.Fatalf("获取锁 key-b 失败: %v", err)
	}

	if v1 == v2 {
		t.Error("期望不同键的锁值不同")
	}
}

func TestUnlock_DoubleUnlock(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	lockValue, err := client.Lock(ctx, "test:lock:double", 10*time.Second)
	if err != nil {
		t.Fatalf("获取锁失败: %v", err)
	}

	err = client.Unlock(ctx, "test:lock:double", lockValue)
	if err != nil {
		t.Fatalf("首次释放锁失败: %v", err)
	}

	err = client.Unlock(ctx, "test:lock:double", lockValue)
	if err != nil {
		t.Fatalf("重复释放锁不应报错: %v", err)
	}
}

func TestLock_Concurrency(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()
	key := "test:lock:concurrent"

	var successCount int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := client.Lock(ctx, key, 10*time.Second)
			if err == nil && v != "" {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if successCount != 1 {
		t.Errorf("期望仅 1 个协程获取锁，实际=%d", successCount)
	}
}
