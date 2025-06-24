package cache

import (
	"testing"
	"time"
)

// TestTimedCache_Basic 测试基本的Set和Get操作
func TestTimedCache_Basic(t *testing.T) {
	cache, err := NewTimedCache[int, string](100, 1*time.Second)
	if err != nil {
		t.Fatalf("创建Timed缓存失败: %v", err)
	}

	// 测试Set和Get
	cache.Set(1, "a")
	val, exists := cache.Get(1)
	if !exists || val != "a" {
		t.Errorf("Get(1) = %v, %v; 期望 'a', true", val, exists)
	}

	// 测试更新值
	cache.Set(1, "a_updated")
	val, exists = cache.Get(1)
	if !exists || val != "a_updated" {
		t.Errorf("Get(1) = %v, %v; 期望 'a_updated', true", val, exists)
	}
}

// TestTimedCache_Expiration 测试缓存过期机制
func TestTimedCache_Expiration(t *testing.T) {
	cache, err := NewTimedCache[int, string](100, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("创建Timed缓存失败: %v", err)
	}

	cache.Set(1, "a")
	// 立即访问应该存在
	_, exists := cache.Get(1)
	if !exists {
		t.Error("Get(1) 应该存在")
	}

	// 等待过期
	time.Sleep(100 * time.Millisecond)

	// 过期后应该不存在
	_, exists = cache.Get(1)
	if exists {
		t.Error("Get(1) 应该过期，但存在")
	}
}

// TestTimedCache_SetWithTTL 测试自定义TTL
func TestTimedCache_SetWithTTL(t *testing.T) {
	cache, err := NewTimedCache[int, string](100, 1*time.Second)
	if err != nil {
		t.Fatalf("创建Timed缓存失败: %v", err)
	}

	// 设置不同TTL的条目
	cache.SetWithTTL(1, "a", 30*time.Millisecond)
	cache.SetWithTTL(2, "b", 150*time.Millisecond)

	// 等待30ms，检查第一个过期
	time.Sleep(50 * time.Millisecond)
	_, exists1 := cache.Get(1)
	_, exists2 := cache.Get(2)

	if exists1 {
		t.Error("Get(1) 应该过期，但存在")
	}
	if !exists2 {
		t.Error("Get(2) 不应该过期，但不存在")
	}

	// 等待剩余时间，检查第二个过期
	time.Sleep(150 * time.Millisecond)
	_, exists2 = cache.Get(2)
	if exists2 {
		t.Error("Get(2) 应该过期，但存在")
	}
}

// TestTimedCache_Delete 测试删除操作
func TestTimedCache_Delete(t *testing.T) {
	cache, err := NewTimedCache[int, string](100, 1*time.Second)
	if err != nil {
		t.Fatalf("创建Timed缓存失败: %v", err)
	}

	cache.Set(1, "a")
	cache.Delete(1)

	_, exists := cache.Get(1)
	if exists {
		t.Error("Get(1) 在删除后应该不存在")
	}
}

// TestTimedCache_Len 测试Len方法
func TestTimedCache_Len(t *testing.T) {
	cache, err := NewTimedCache[int, string](100, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("创建Timed缓存失败: %v", err)
	}

	if cache.Len() != 0 {
		t.Errorf("Len() = %d; 期望 0", cache.Len())
	}

	cache.Set(1, "a")
	if cache.Len() != 1 {
		t.Errorf("Len() = %d; 期望 1", cache.Len())
	}

	// 等待过期
	time.Sleep(100 * time.Millisecond)
	// 触发清理
	cache.Get(1)
	if cache.Len() != 0 {
		t.Errorf("过期后 Len() = %d; 期望 0", cache.Len())
	}
}

// TestTimedCache_Clear 测试Clear方法
func TestTimedCache_Clear(t *testing.T) {
	cache, err := NewTimedCache[int, string](100, 1*time.Second)
	if err != nil {
		t.Fatalf("创建Timed缓存失败: %v", err)
	}

	cache.Set(1, "a")
	cache.Set(2, "b")
	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("Clear() 后 Len() = %d; 期望 0", cache.Len())
	}
}

// BenchmarkTimedCache_SetGet 基准测试Set和Get操作性能
func BenchmarkTimedCache_SetGet(b *testing.B) {
	cache, _ := NewTimedCache[int, int](1000, 1*time.Second)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := i % 1000
		cache.Set(key, i)
		cache.Get(key)
	}
}

// BenchmarkTimedCache_Expiration 基准测试过期清理性能
func BenchmarkTimedCache_Expiration(b *testing.B) {
	cache, _ := NewTimedCache[int, int](100, 10*time.Millisecond)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.SetWithTTL(i, i, time.Duration(i%20)*time.Millisecond)
		if i%10 == 0 {
			// 触发清理
			cache.Get(0)
		}
	}
}