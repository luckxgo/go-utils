package cache

import (
	"testing"
)

// TestFIFOCache_Basic 测试基本的Set和Get操作
func TestFIFOCache_Basic(t *testing.T) {
	fifo, err := NewFIFOCache[int, string](2)
	if err != nil {
		t.Fatalf("创建FIFO缓存失败: %v", err)
	}

	// 测试Set和Get
	fifo.Set(1, "a")
	val, exists := fifo.Get(1)
	if !exists || val != "a" {
		t.Errorf("Get(1) = %v, %v; 期望 'a', true", val, exists)
	}

	// 测试更新值
	fifo.Set(1, "a_updated")
	val, exists = fifo.Get(1)
	if !exists || val != "a_updated" {
		t.Errorf("Get(1) = %v, %v; 期望 'a_updated', true", val, exists)
	}
}

// TestFIFOCache_Eviction 测试缓存淘汰机制
func TestFIFOCache_Eviction(t *testing.T) {
	fifo, err := NewFIFOCache[int, string](2)
	if err != nil {
		t.Fatalf("创建FIFO缓存失败: %v", err)
	}

	fifo.Set(1, "a")
	fifo.Set(2, "b")
	fifo.Set(3, "c") // 触发淘汰

	// 验证最早插入的1被淘汰
	_, exists := fifo.Get(1)
	if exists {
		t.Error("Get(1) 应该被淘汰，但存在")
	}

	// 验证2和3存在
	val, exists := fifo.Get(2)
	if !exists || val != "b" {
		t.Errorf("Get(2) = %v, %v; 期望 'b', true", val, exists)
	}

	val, exists = fifo.Get(3)
	if !exists || val != "c" {
		t.Errorf("Get(3) = %v, %v; 期望 'c', true", val, exists)
	}
}

// TestFIFOCache_Delete 测试删除操作
func TestFIFOCache_Delete(t *testing.T) {
	fifo, err := NewFIFOCache[int, string](2)
	if err != nil {
		t.Fatalf("创建FIFO缓存失败: %v", err)
	}

	fifo.Set(1, "a")
	fifo.Delete(1)

	_, exists := fifo.Get(1)
	if exists {
		t.Error("Get(1) 在删除后应该不存在")
	}

	// 测试删除不存在的键
	fifo.Delete(2) // 应该无错误
}

// TestFIFOCache_Len 测试Len方法
func TestFIFOCache_Len(t *testing.T) {
	fifo, err := NewFIFOCache[int, string](2)
	if err != nil {
		t.Fatalf("创建FIFO缓存失败: %v", err)
	}

	if fifo.Len() != 0 {
		t.Errorf("Len() = %d; 期望 0", fifo.Len())
	}

	fifo.Set(1, "a")
	if fifo.Len() != 1 {
		t.Errorf("Len() = %d; 期望 1", fifo.Len())
	}

	fifo.Set(2, "b")
	if fifo.Len() != 2 {
		t.Errorf("Len() = %d; 期望 2", fifo.Len())
	}

	fifo.Set(3, "c") // 触发淘汰
	if fifo.Len() != 2 {
		t.Errorf("Len() = %d; 期望 2", fifo.Len())
	}
}

// TestFIFOCache_Clear 测试Clear方法
func TestFIFOCache_Clear(t *testing.T) {
	fifo, err := NewFIFOCache[int, string](2)
	if err != nil {
		t.Fatalf("创建FIFO缓存失败: %v", err)
	}

	fifo.Set(1, "a")
	fifo.Set(2, "b")
	fifo.Clear()

	if fifo.Len() != 0 {
		t.Errorf("Clear() 后 Len() = %d; 期望 0", fifo.Len())
	}

	_, exists := fifo.Get(1)
	if exists {
		t.Error("Clear() 后 Get(1) 应该不存在")
	}
}

// BenchmarkFIFOCache_SetGet 基准测试Set和Get操作性能
func BenchmarkFIFOCache_SetGet(b *testing.B) {
	fifo, _ := NewFIFOCache[int, int](1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := i % 1000
		fifo.Set(key, i)
		fifo.Get(key)
	}
}

// BenchmarkFIFOCache_Eviction 基准测试缓存淘汰性能
func BenchmarkFIFOCache_Eviction(b *testing.B) {
	fifo, _ := NewFIFOCache[int, int](100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fifo.Set(i, i)
	}
}