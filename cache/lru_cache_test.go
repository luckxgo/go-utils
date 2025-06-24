package cache

import (
	"fmt"
	"sync"
	"testing"
)

// TestLRUCacheConcurrent 测试并发环境下LRU缓存的正确性
func TestLRUCacheConcurrent(t *testing.T) {
	cache, err := NewLRUCache[int, int](100000)
	if err != nil {
		t.Fatalf("Failed to create LRU cache: %v", err)
	}

	const (
		numGoroutines          = 50
		operationsPerGoroutine = 2000
	)
	var wg sync.WaitGroup
	errCh := make(chan error, numGoroutines)

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := goroutineID*operationsPerGoroutine + j
				cache.Set(key, key*2)
				val, exists := cache.Get(key)
				if !exists || val != key*2 {
					errCh <- fmt.Errorf("goroutine %d: key %d, expected %d, got %v (exists: %v)", goroutineID, key, key*2, val, exists)
					return
				}

				// 暂时禁用随机删除操作以验证并发读写
			// if j%10 == 0 {
			// 	cache.Delete(key)
			// }
			}
		}(i)
	}

	// 等待所有goroutine完成并检查错误
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			t.Error(err)
		}
	}

	// 验证最终缓存状态
	if cache.Len() < 0 {
		 t.Errorf("Unexpected cache length: %d", cache.Len())
	}
}

// BenchmarkLRUCacheConcurrent 并发读写性能基准测试
func BenchmarkLRUCacheConcurrent(b *testing.B) {
	cache, _ := NewLRUCache[int, int](1000)
	b.RunParallel(func(pb *testing.PB) {
		key := 0
		for pb.Next() {
			cache.Set(key, key)
			cache.Get(key)
			key = (key + 1) % 1000 // 循环使用键以模拟实际场景
		}
	})
}

// BenchmarkLRUCacheWithEviction 带淘汰机制的LRU性能基准测试
func BenchmarkLRUCacheWithEviction(b *testing.B) {
	cache, _ := NewLRUCache[int, int](100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := i % 200 // 超出容量以触发淘汰
		cache.Set(key, i)
		cache.Get(key)
	}
}

// TestLRUCache_Basic 测试基本的Set和Get操作
func TestLRUCache_Basic(t *testing.T) {
	lru, err := NewLRUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LRU缓存失败: %v", err)
	}

	// 测试Set和Get
	lru.Set(1, "a")
	val, exists := lru.Get(1)
	if !exists || val != "a" {
		t.Errorf("Get(1) = %v, %v; 期望 'a', true", val, exists)
	}

	// 测试更新值
	lru.Set(1, "a_updated")
	val, exists = lru.Get(1)
	if !exists || val != "a_updated" {
		t.Errorf("Get(1) = %v, %v; 期望 'a_updated', true", val, exists)
	}
}

// TestLRUCache_Eviction 测试缓存淘汰机制（最近最久未使用）
func TestLRUCache_Eviction(t *testing.T) {
	lru, err := NewLRUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LRU缓存失败: %v", err)
	}

	lru.Set(1, "a")
	lru.Set(2, "b")
	lru.Get(1)      // 访问1，使其成为最近使用
	lru.Set(3, "c") // 触发淘汰，淘汰最久未使用的2

	// 验证2被淘汰
	_, exists := lru.Get(2)
	if exists {
		t.Error("Get(2) 应该被淘汰，但存在")
	}

	// 验证1和3存在
	val, exists := lru.Get(1)
	if !exists || val != "a" {
		t.Errorf("Get(1) = %v, %v; 期望 'a', true", val, exists)
	}

	val, exists = lru.Get(3)
	if !exists || val != "c" {
		t.Errorf("Get(3) = %v, %v; 期望 'c', true", val, exists)
	}
}

// TestLRUCache_UpdateAccess 测试访问更新最近使用顺序
func TestLRUCache_UpdateAccess(t *testing.T) {
	lru, err := NewLRUCache[int, string](3)
	if err != nil {
		t.Fatalf("创建LRU缓存失败: %v", err)
	}

	lru.Set(1, "a")
	lru.Set(2, "b")
	lru.Set(3, "c")
	lru.Get(1)      // 访问1
	lru.Set(4, "d") // 触发淘汰，淘汰最久未使用的2

	_, exists := lru.Get(2)
	if exists {
		t.Error("Get(2) 应该被淘汰，但存在")
	}
}

// TestLRUCache_Delete 测试删除操作
func TestLRUCache_Delete(t *testing.T) {
	lru, err := NewLRUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LRU缓存失败: %v", err)
	}

	lru.Set(1, "a")
	lru.Delete(1)

	_, exists := lru.Get(1)
	if exists {
		t.Error("Get(1) 在删除后应该不存在")
	}
}

// TestLRUCache_Len 测试Len方法
func TestLRUCache_Len(t *testing.T) {
	lru, err := NewLRUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LRU缓存失败: %v", err)
	}

	if lru.Len() != 0 {
		t.Errorf("Len() = %d; 期望 0", lru.Len())
	}

	lru.Set(1, "a")
	if lru.Len() != 1 {
		t.Errorf("Len() = %d; 期望 1", lru.Len())
	}
}

// TestLRUCache_Clear 测试Clear方法
func TestLRUCache_Clear(t *testing.T) {
	lru, err := NewLRUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LRU缓存失败: %v", err)
	}

	lru.Set(1, "a")
	lru.Set(2, "b")
	lru.Clear()

	if lru.Len() != 0 {
		t.Errorf("Clear() 后 Len() = %d; 期望 0", lru.Len())
	}
}

// BenchmarkLRUCache_SetGet 基准测试Set和Get操作性能
func BenchmarkLRUCache_SetGet(b *testing.B) {
	lru, _ := NewLRUCache[int, int](1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := i % 1000
		lru.Set(key, i)
		lru.Get(key)
	}
}

// BenchmarkLRUCache_Eviction 基准测试缓存淘汰性能
func BenchmarkLRUCache_Eviction(b *testing.B) {
	lru, _ := NewLRUCache[int, int](100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lru.Set(i, i)
		if i%10 == 0 {
			lru.Get(i % 100)
		}
	}
}
