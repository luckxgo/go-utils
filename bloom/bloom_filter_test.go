package bloom

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// TestNewBloomFilter 测试布隆过滤器的创建
func TestNewBloomFilter(t *testing.T) {
	// 测试正常创建
	bf, err := NewBloomFilter(1000, 0.01)
	if err != nil {
		t.Fatalf("创建布隆过滤器失败: %v", err)
	}
	if bf == nil {
		t.Fatal("返回的布隆过滤器为nil")
	}
	if bf.k <= 0 || bf.m <= 0 {
		t.Errorf("布隆过滤器参数异常: k=%d, m=%d", bf.k, bf.m)
	}

	// 测试无效参数
	_, err = NewBloomFilter(0, 0.01)
	if err == nil {
		t.Error("预期n=0时返回错误，但未返回")
	}

	_, err = NewBloomFilter(1000, 0)
	if err == nil {
		t.Error("预期p=0时返回错误，但未返回")
	}

	_, err = NewBloomFilter(1000, 1)
	if err == nil {
		t.Error("预期p=1时返回错误，但未返回")
	}
}

// TestBloomFilter_Add_Contains 测试添加和查询功能
func TestBloomFilter_Add_Contains(t *testing.T) {
	bf, err := NewBloomFilter(100, 0.01)
	if err != nil {
		t.Fatalf("创建布隆过滤器失败: %v", err)
	}

	// 添加元素
	elements := [][]byte{[]byte("test1"), []byte("test2"), []byte("test3")}
	for _, e := range elements {
		bf.Add(e)
	}

	// 检查已添加的元素
	for _, e := range elements {
		if !bf.Contains(e) {
			t.Errorf("元素 %s 应该存在，但未检测到", e)
		}
	}

	// 检查未添加的元素
	if bf.Contains([]byte("nonexistent")) {
		t.Error("检测到不存在的元素，可能是误判")
	}
}

// TestBloomFilter_FalsePositive 测试误判率
func TestBloomFilter_FalsePositive(t *testing.T) {
	// 设置较小的规模以便测试
	n := 1000
	p := 0.01
	bf, err := NewBloomFilter(n, p)
	if err != nil {
		t.Fatalf("创建布隆过滤器失败: %v", err)
	}

	// 添加n个元素
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	added := make([][]byte, n)
	for i := 0; i < n; i++ {
		data := make([]byte, 16)
		r.Read(data)
		added[i] = data
		bf.Add(data)
	}

	// 测试10000个未添加的元素，计算误判率
	testCount := 10000
	falsePositives := 0

	for i := 0; i < testCount; i++ {
		data := make([]byte, 16)
		r.Read(data)
		// 确保测试数据不在已添加列表中
		found := false
		for _, a := range added {
			if string(a) == string(data) {
				found = true
				break
			}
		}
		if found {
			continue
		}

		if bf.Contains(data) {
			falsePositives++
		}
	}

	// 计算实际误判率
	actualP := float64(falsePositives) / float64(testCount)
	// 允许一定的误差范围
	if actualP > p*2 {
		t.Errorf("误判率超出预期: 预期%.4f, 实际%.4f", p, actualP)
	}
}

// TestBloomFilter_Reset 测试重置功能
func TestBloomFilter_Reset(t *testing.T) {
	bf, err := NewBloomFilter(100, 0.01)
	if err != nil {
		t.Fatalf("创建布隆过滤器失败: %v", err)
	}

	elem := []byte("test_reset")
	bf.Add(elem)

	if !bf.Contains(elem) {
		t.Error("添加元素后检测失败")
	}

	bf.Reset()

	if bf.Contains(elem) {
		t.Error("重置后仍能检测到元素")
	}
}

// BenchmarkBloomFilter_Add 基准测试添加元素性能
func BenchmarkBloomFilter_Add(b *testing.B) {
	bf, err := NewBloomFilter(1000000, 0.01)
	if err != nil {
		b.Fatalf("创建布隆过滤器失败: %v", err)
	}
	data := []byte("benchmark_test_data")

	b.ResetTimer() // 重置计时器，排除初始化时间
	for i := 0; i < b.N; i++ {
		bf.Add(data)
	}
}

// BenchmarkBloomFilter_Contains 基准测试查询元素性能
func BenchmarkBloomFilter_Contains(b *testing.B) {
	bf, err := NewBloomFilter(1000000, 0.01)
	if err != nil {
		b.Fatalf("创建布隆过滤器失败: %v", err)
	}
	data := []byte("benchmark_test_data")
	bf.Add(data)

	b.ResetTimer() // 重置计时器，排除初始化和添加元素时间
	for i := 0; i < b.N; i++ {
		bf.Contains(data)
	}
}

// BenchmarkBloomFilter_HighLoad 添加大量元素后的性能测试
func BenchmarkBloomFilter_HighLoad(b *testing.B) {
	// 创建一个可容纳100万元素的过滤器
	bf, err := NewBloomFilter(1000000, 0.01)
	if err != nil {
		b.Fatalf("创建布隆过滤器失败: %v", err)
	}

	// 预先添加50万元素
	for i := 0; i < 500000; i++ {
		bf.Add([]byte(fmt.Sprintf("test_data_%d", i)))
	}

	testData := []byte("high_load_test_data")
	bf.Add(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Contains(testData)
	}
}
