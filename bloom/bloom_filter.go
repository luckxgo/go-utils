package bloom

import (
	"errors"
	"hash/fnv"
	"math"
)

// BloomFilter 实现布隆过滤器数据结构
// 用于高效判断元素是否存在于集合中，存在一定的误判率但不会漏判
type BloomFilter struct {
	bits  []uint64  // 位数组，使用uint64切片存储以提高空间效率
	k     int       // 哈希函数数量
	m     int       // 位数组总位数
	hashes []func([]byte) uint64 // 哈希函数列表
}

// NewBloomFilter 创建一个新的布隆过滤器
// n: 预期元素数量
// p: 可接受的误判率(0 < p < 1)
// 返回布隆过滤器实例和可能的错误
func NewBloomFilter(n int, p float64) (*BloomFilter, error) {
	if n <= 0 {
		return nil, errors.New("预期元素数量n必须大于0")
	}
	if p <= 0 || p >= 1 {
		return nil, errors.New("误判率p必须在(0, 1)范围内")
	}

	// 计算最优位数组大小m和哈希函数数量k
	m := int(-float64(n) * math.Log(p) / (math.Log(2) * math.Log(2)))
	k := int(math.Round(float64(m) / float64(n) * math.Log(2)))

	// 确保m和k至少为1
	if m <= 0 {
		m = 1
	}
	if k <= 0 {
		k = 1
	}

	// 初始化位数组，向上取整到uint64的倍数
	bits := make([]uint64, (m+63)/64)

	// 创建哈希函数列表 - 使用双重哈希策略确保独立性
	hashes := make([]func([]byte) uint64, k)
	for i := 0; i < k; i++ {
		// 捕获循环变量i的值，避免闭包引用问题
		seed := i
		hashes[i] = func(data []byte) uint64 {
			// 使用两种不同的哈希算法生成基础哈希值
			h1 := fnv.New64a()
			h1.Write(data)
			hash1 := h1.Sum64()

			h2 := fnv.New64()
			h2.Write(data)
			hash2 := h2.Sum64()

			// 结合种子生成独立的哈希函数
			return hash1 + uint64(seed)*hash2
		}
	}

	return &BloomFilter{
		bits:  bits,
		k:     k,
		m:     m,
		hashes: hashes,
	}, nil
}

// Add 将元素添加到布隆过滤器
// data: 要添加的元素字节表示
func (bf *BloomFilter) Add(data []byte) {
	for _, hash := range bf.hashes {
		idx := hash(data) % uint64(bf.m)
		bf.bits[idx/64] |= 1 << (idx % 64)
	}
}

// Contains 检查元素是否可能存在于布隆过滤器中
// 返回true表示可能存在(有一定误判率)，返回false表示一定不存在
func (bf *BloomFilter) Contains(data []byte) bool {
	for _, hash := range bf.hashes {
		idx := hash(data) % uint64(bf.m)
		if (bf.bits[idx/64] & (1 << (idx % 64))) == 0 {
			return false
		}
	}
	return true
}

// Reset 重置布隆过滤器，清除所有元素
func (bf *BloomFilter) Reset() {
	bf.bits = make([]uint64, len(bf.bits))
}