package cache

import "errors"

// FIFOCache 基于先进先出(First-In-First-Out)策略的缓存实现
// 当缓存达到容量限制时，最早插入的元素会被优先淘汰
// K为键类型，必须支持比较操作；V为值类型，可以是任意类型
type FIFOCache[K comparable, V any] struct {
	cache    map[K]V   // 存储键值对的哈希表，提供O(1)时间复杂度的读写操作
	queue    []K       // 维护元素插入顺序的队列，用于实现FIFO淘汰策略
	capacity int       // 缓存的最大容量，超过此容量将触发淘汰机制
}

// NewFIFOCache 创建新的FIFO缓存实例
// capacity为缓存容量，必须大于0，否则返回错误
// 返回值:
//   *FIFOCache[K, V]: 成功创建的缓存实例
//   error: 当capacity <= 0时返回非nil错误
func NewFIFOCache[K comparable, V any](capacity int) (*FIFOCache[K, V], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}
	return &FIFOCache[K, V]{
		cache:    make(map[K]V),
		queue:    make([]K, 0, capacity),
		capacity: capacity,
	}, nil
}

// Get 从缓存中获取键对应的值
// 注意：FIFO策略中，Get操作不会改变元素的淘汰顺序
// 参数:
//   key: 要查找的键
// 返回值:
//   value: 键对应的值，如果键不存在则返回V类型的零值
//   exists: 布尔值，表示键是否存在于缓存中
func (f *FIFOCache[K, V]) Get(key K) (value V, exists bool) {
	value, exists = f.cache[key]
	return
}

// Set 将键值对存入缓存
// 如果键已存在，仅更新值而不改变其在队列中的位置
// 如果键不存在且缓存已满，会先移除最早插入的键（队列头部元素），再插入新键值对
// 参数:
//   key: 要存储的键
//   value: 要存储的值
func (f *FIFOCache[K, V]) Set(key K, value V) {
	// 如果键已存在，仅更新值
	if _, exists := f.cache[key]; exists {
		f.cache[key] = value
		return
	}

	// 如果缓存已满，移除队列头部元素
	if len(f.queue) >= f.capacity {
		oldestKey := f.queue[0]
		f.queue = f.queue[1:]
		delete(f.cache, oldestKey)
	}

	// 添加新元素到队列尾部和缓存
	f.queue = append(f.queue, key)
	f.cache[key] = value
}

// Delete 从缓存中删除指定键
// 如果键不存在，此操作无效果
// 参数:
//   key: 要删除的键
func (f *FIFOCache[K, V]) Delete(key K) {
	if _, exists := f.cache[key]; !exists {
		return
	}

	// 从缓存中删除键
	delete(f.cache, key)

	// 从队列中删除键
	for i, k := range f.queue {
		if k == key {
			// 保持队列连续性（不改变剩余元素顺序）
			f.queue = append(f.queue[:i], f.queue[i+1:]...)
			break
		}
	}
}

// Len 返回当前缓存中的元素数量
// 返回值:
//   int: 缓存中已存储的键值对数量
func (f *FIFOCache[K, V]) Len() int {
	return len(f.cache)
}

// Clear 清空缓存中的所有元素
// 此操作会重置缓存的内部状态，包括哈希表和队列
func (f *FIFOCache[K, V]) Clear() {
	f.cache = make(map[K]V)
	f.queue = make([]K, 0, f.capacity)
}