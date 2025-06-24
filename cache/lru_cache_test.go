package cache

import (
	"testing"
)

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
	lru.Get(1)        // 访问1，使其成为最近使用
	lru.Set(3, "c")  // 触发淘汰，淘汰最久未使用的2

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
	lru.Get(1)        // 访问1
	lru.Set(4, "d")  // 触发淘汰，淘汰最久未使用的2

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