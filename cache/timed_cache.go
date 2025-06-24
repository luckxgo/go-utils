package cache

import (
	"container/heap"
	"errors"
	"sync"
	"time"
)

// heapEntry 用于最小堆中的元素，存储键和过期时间
type heapEntry[K comparable] struct {
	key        K          // 缓存键
	expiration int64      // 过期时间戳（纳秒）
	index      int        // 在堆中的索引，用于更新堆结构
}

// expirationHeap 实现最小堆接口，按过期时间戳升序排序
// 堆顶元素始终是最早过期的条目，用于高效获取和删除过期缓存
type expirationHeap[K comparable] []*heapEntry[K]

// Len 返回堆中元素的数量，实现heap.Interface
func (h expirationHeap[K]) Len() int { return len(h) }

// Less 比较i和j位置元素的过期时间，实现heap.Interface
// 过期时间较小的元素优先级更高（靠近堆顶）
func (h expirationHeap[K]) Less(i, j int) bool {
	return h[i].expiration < h[j].expiration
}

// Swap 交换i和j位置的元素，实现heap.Interface
// 同时更新元素的索引值以保持堆结构正确性
func (h expirationHeap[K]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

// Push 向堆中添加元素，实现heap.Interface
// 将元素添加到堆尾部并更新其索引
func (h *expirationHeap[K]) Push(x interface{}) {
	n := len(*h)
	entry := x.(*heapEntry[K])
	entry.index = n
	*h = append(*h, entry)
}

// Pop 从堆中移除并返回最小元素（堆顶），实现heap.Interface
// 将堆尾部元素移至堆顶并调整堆结构
// 返回值:
//   interface{}: 堆中最早过期的元素
func (h *expirationHeap[K]) Pop() interface{} {
	old := *h
	n := len(old)
	entry := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	entry.index = -1 // 标记为已移除
	*h = old[0 : n-1]
	return entry
}

// timedEntry 缓存中的条目，包含值和过期时间
type timedEntry[V any] struct {
	value      V          // 缓存值
	expiration int64      // 过期时间戳（纳秒）
}

// timedCacheOptions 用于配置TimedCache的选项
type timedCacheOptions struct {
	concurrentSafe bool // 是否启用并发安全
}

// TimedOption 定义配置TimedCache的函数类型
type TimedOption func(*timedCacheOptions)

// WithTimedConcurrentSafe 设置是否启用并发安全
// 参数:
//   enabled: true表示启用并发安全，false表示禁用
// 返回值:
//   TimedOption: 用于配置缓存的选项函数
func WithTimedConcurrentSafe(enabled bool) TimedOption {
	return func(o *timedCacheOptions) {
		o.concurrentSafe = enabled
	}
}

// TimedCache 基于过期时间的缓存实现
// 支持设置默认TTL(Time-To-Live)，条目过期后自动失效
// 当缓存达到容量限制时，会优先淘汰最早过期的条目
// K为键类型（必须可比较），V为值类型

type TimedCache[K comparable, V any] struct {
	cache          map[K]*timedEntry[V]   // 存储键值对的哈希表，提供O(1)时间复杂度的读写
	heap           *expirationHeap[K]     // 最小堆，用于跟踪过期时间，支持高效获取最早过期条目
	heapEntries    map[K]*heapEntry[K]    // 键到堆条目的映射，用于快速更新堆
	capacity       int                    // 最大容量，防止内存溢出
	defaultTTL     time.Duration          // 默认过期时间，当使用Set方法时应用
	concurrentSafe bool                   // 是否启用并发安全
	mu             sync.RWMutex           // 读写锁，用于并发控制
}

// NewTimedCache 创建新的超时缓存实例
// 参数:
//   capacity: 最大缓存条目数，必须大于0
//   defaultTTL: 默认过期时间，必须大于0
// 返回值:
//   *TimedCache[K, V]: 成功创建的缓存实例
//   error: 当capacity <= 0或defaultTTL <= 0时返回非nil错误
func NewTimedCache[K comparable, V any](capacity int, defaultTTL time.Duration, options ...TimedOption) (*TimedCache[K, V], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}
	if defaultTTL <= 0 {
		return nil, errors.New("default TTL must be positive")
	}
	
	opts := timedCacheOptions{
		concurrentSafe: true, // 默认启用并发安全
	}
	for _, option := range options {
		option(&opts)
	}
	
	return &TimedCache[K, V]{
		cache:          make(map[K]*timedEntry[V]),
		heap:           &expirationHeap[K]{},
		heapEntries:    make(map[K]*heapEntry[K]),
		capacity:       capacity,
		defaultTTL:     defaultTTL,
		concurrentSafe: opts.concurrentSafe,
		mu:             sync.RWMutex{},
	}, nil
}

// Get 获取缓存中键对应的值
// 调用此方法会先清理所有过期条目，然后检查指定键是否存在且有效
// 参数:
//   key: 要查找的键
// 返回值:
//   value: 键对应的值，如果键不存在或已过期则返回V类型的零值
//   exists: 布尔值，表示键是否存在且未过期
func (t *TimedCache[K, V]) Get(key K) (value V, exists bool) {
	if t.concurrentSafe {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	
	t.cleanupExpired()

	entry, exists := t.cache[key]
	if !exists {
		return value, false
	}

	now := time.Now().UnixNano()
	if entry.expiration < now {
		delete(t.cache, key)
		return value, false
	}

	return entry.value, true
}

// Set 使用默认TTL存储键值对
// 等效于调用SetWithTTL(key, value, t.defaultTTL)
// 参数:
//   key: 要存储的键
//   value: 要存储的值
func (t *TimedCache[K, V]) Set(key K, value V) {
	t.SetWithTTL(key, value, t.defaultTTL)
}

// SetWithTTL 存储带有自定义过期时间的键值对
// 如果键已存在，更新其值和过期时间
// 如果缓存满，会先淘汰最早过期的条目
// 参数:
//   key: 要存储的键
//   value: 要存储的值
//   ttl: 该条目的生存时间，必须为正数
func (t *TimedCache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	if t.concurrentSafe {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	
	t.cleanupExpired()

	expiration := time.Now().Add(ttl).UnixNano()

	// 如果键已存在，更新值和过期时间
	if entry, exists := t.cache[key]; exists {
		entry.value = value
		oldExpiration := entry.expiration
		entry.expiration = expiration
		// 查找并移除堆中该键的旧条目
		for i, e := range *t.heap {
			if e.key == key && e.expiration == oldExpiration {
				heap.Remove(t.heap, i)
				break
			}
		}
		heap.Push(t.heap, &heapEntry[K]{
			key:        key,
			expiration: expiration,
		})
		return
	}

	// 如果缓存满了，驱逐最早过期的条目
	for len(t.cache) >= t.capacity {
		if t.heap.Len() == 0 {
			break // 理论上不会发生，防止死循环
		}
		oldest := heap.Pop(t.heap).(*heapEntry[K])
		// 检查堆条目是否仍然有效（缓存中存在且过期时间匹配）
		if entry, exists := t.cache[oldest.key]; exists && entry.expiration == oldest.expiration {
			delete(t.cache, oldest.key)
		}
	}

	// 创建新条目并添加到缓存
	newEntry := &timedEntry[V]{
		value:      value,
		expiration: expiration,
	}
	t.cache[key] = newEntry

	// 添加到堆
	newHeapEntry := &heapEntry[K]{
		key:        key,
		expiration: expiration,
	}
	heap.Push(t.heap, newHeapEntry)
	t.heapEntries[key] = newHeapEntry
}

// Delete 从缓存中删除指定键
// 如果键不存在，此操作无效果
// 参数:
//   key: 要删除的键
func (t *TimedCache[K, V]) Delete(key K) {
	if t.concurrentSafe {
		t.mu.Lock()
		defer t.mu.Unlock()
	}

	// 从堆和映射中删除
	if heapEntry, exists := t.heapEntries[key]; exists {
		heap.Remove(t.heap, heapEntry.index)
		delete(t.heapEntries, key)
	}
	// 从缓存中删除
	delete(t.cache, key)
}

// Len 返回当前有效缓存条目数量
// 调用此方法会先清理所有过期条目
// 返回值:
//   int: 缓存中未过期的键值对数量
func (t *TimedCache[K, V]) Len() int {
	if t.concurrentSafe {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	t.cleanupExpired()
	return len(t.cache)
}

// Clear 清空所有缓存条目
// 此操作会重置缓存的内部状态，包括哈希表和堆
func (t *TimedCache[K, V]) Clear() {
	if t.concurrentSafe {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	t.cache = make(map[K]*timedEntry[V])
	*t.heap = (*t.heap)[:0] // 清空堆
}

// cleanupExpired 清理所有过期的缓存条目
// 此方法应在持有锁的情况下调用
func (t *TimedCache[K, V]) cleanupExpired() {
	now := time.Now().UnixNano()

	// 循环检查并移除所有过期元素
	for t.heap.Len() > 0 {
		// 获取并弹出堆顶元素（最早过期）
		entry := heap.Pop(t.heap).(*heapEntry[K])
		if entry.expiration > now {
			// 未过期，推回堆中并停止清理
		heap.Push(t.heap, entry)
			break
		}

		// 从缓存和堆条目映射中删除过期条目
		if cacheEntry, exists := t.cache[entry.key]; exists && cacheEntry.expiration == entry.expiration {
			delete(t.cache, entry.key)
		}
		delete(t.heapEntries, entry.key)
	}
}