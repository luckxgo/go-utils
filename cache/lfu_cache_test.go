package cache

import (
	"testing"
)

// TestLFUCache_Basic 测试基本的Set和Get操作
func TestLFUCache_Basic(t *testing.T) {
	lfu, err := NewLFUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LFU缓存失败: %v", err)
	}

	// 测试Set和Get
	lfu.Set(1, "a")
	val, exists := lfu.Get(1)
	if !exists || val != "a" {
		t.Errorf("Get(1) = %v, %v; 期望 'a', true", val, exists)
	}

	// 测试更新值
	lfu.Set(1, "a_updated")
	val, exists = lfu.Get(1)
	if !exists || val != "a_updated" {
		t.Errorf("Get(1) = %v, %v; 期望 'a_updated', true", val, exists)
	}
}

// TestLFUCache_Eviction 测试缓存淘汰机制（最少使用）
func TestLFUCache_Eviction(t *testing.T) {
	lfu, err := NewLFUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LFU缓存失败: %v", err)
	}

	lfu.Set(1, "a")  // freq:1
	lfu.Set(2, "b")  // freq:1
	lfu.Get(1)        // freq:2
	lfu.Set(3, "c")  // 触发淘汰，淘汰频率最低的2

	// 验证2被淘汰
	_, exists := lfu.Get(2)
	if exists {
		t.Error("Get(2) 应该被淘汰，但存在")
	}

	// 验证1和3存在
	val, exists := lfu.Get(1)
	if !exists || val != "a" {
		t.Errorf("Get(1) = %v, %v; 期望 'a', true", val, exists)
	}

	val, exists = lfu.Get(3)
	if !exists || val != "c" {
		t.Errorf("Get(3) = %v, %v; 期望 'c', true", val, exists)
	}
}

// TestLFUCache_FreqOrder 测试相同频率下的淘汰顺序（最久未使用）
func TestLFUCache_FreqOrder(t *testing.T) {
	lfu, err := NewLFUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LFU缓存失败: %v", err)
	}

	lfu.Set(1, "a")  // freq:1
	lfu.Set(2, "b")  // freq:1
	lfu.Get(1)        // freq:2
	lfu.Get(2)        // freq:2
	lfu.Set(3, "c")  // 触发淘汰，相同频率下淘汰最久未使用的1

	_, exists := lfu.Get(1)
	if exists {
		t.Error("Get(1) 应该被淘汰，但存在")
	}
}

// TestLFUCache_Delete 测试删除操作
func TestLFUCache_Delete(t *testing.T) {
	lfu, err := NewLFUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LFU缓存失败: %v", err)
	}

	lfu.Set(1, "a")
	lfu.Delete(1)

	_, exists := lfu.Get(1)
	if exists {
		t.Error("Get(1) 在删除后应该不存在")
	}
}

// TestLFUCache_Len 测试Len方法
func TestLFUCache_Len(t *testing.T) {
	lfu, err := NewLFUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LFU缓存失败: %v", err)
	}

	if lfu.Len() != 0 {
		t.Errorf("Len() = %d; 期望 0", lfu.Len())
	}

	lfu.Set(1, "a")
	if lfu.Len() != 1 {
		t.Errorf("Len() = %d; 期望 1", lfu.Len())
	}
}

// TestLFUCache_Clear 测试Clear方法
func TestLFUCache_Clear(t *testing.T) {
	lfu, err := NewLFUCache[int, string](2)
	if err != nil {
		t.Fatalf("创建LFU缓存失败: %v", err)
	}

	lfu.Set(1, "a")
	lfu.Set(2, "b")
	lfu.Clear()

	if lfu.Len() != 0 {
		t.Errorf("Clear() 后 Len() = %d; 期望 0", lfu.Len())
	}
}

// BenchmarkLFUCache_SetGet 基准测试Set和Get操作性能
func BenchmarkLFUCache_SetGet(b *testing.B) {
	lfu, _ := NewLFUCache[int, int](1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := i % 1000
		lfu.Set(key, i)
		lfu.Get(key)
	}
}

// BenchmarkLFUCache_Eviction 基准测试缓存淘汰性能
func BenchmarkLFUCache_Eviction(b *testing.B) {
	lfu, _ := NewLFUCache[int, int](100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lfu.Set(i, i)
		if i%10 == 0 {
			lfu.Get(i % 100)
		}
	}
}