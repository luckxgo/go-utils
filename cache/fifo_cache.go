package cache

import (
	"container/list"
	"errors"
	"sync"
)

// FIFOCache 基于先进先出(First-In-First-Out)策略的缓存实现
// 当缓存达到容量限制时，最早插入的元素会被优先淘汰
// K为键类型，必须支持比较操作；V为值类型，可以是任意类型
type FIFOCache[K comparable, V any] struct {
	cache          map[K]cacheEntry[K, V] // 存储键值对及链表节点的哈希表
	queue          *list.List             // 维护元素插入顺序的双向链表
	capacity       int                    // 缓存的最大容量，超过此容量将触发淘汰机制
	concurrentSafe bool                   // 是否启用并发安全模式
	mu             sync.RWMutex           // 读写锁，在并发安全模式下使用
}

// cacheEntry 缓存条目，存储值和对应的链表节点
type cacheEntry[K comparable, V any] struct {
	value V
	node  *list.Element
}

// Option 定义FIFO缓存的配置选项函数类型
type Option func(*fifoCacheOptions)

// fifoCacheOptions FIFO缓存的配置选项
type fifoCacheOptions struct {
	concurrentSafe bool // 是否启用并发安全
}

// WithConcurrentSafe 设置是否启用并发安全模式
// concurrentSafe为true时启用并发安全，false时禁用
func WithConcurrentSafe(concurrentSafe bool) Option {
	return func(o *fifoCacheOptions) {
		o.concurrentSafe = concurrentSafe
	}
}

// NewFIFOCache 创建新的FIFO缓存实例
// capacity为缓存容量，必须大于0，否则返回错误
// options为可选配置参数，可通过WithConcurrentSafe等函数设置
// 返回值:
//
//	*FIFOCache[K, V]: 成功创建的缓存实例
//	error: 当capacity <= 0时返回非nil错误
func NewFIFOCache[K comparable, V any](capacity int, options ...Option) (*FIFOCache[K, V], error) {
	if capacity <= 0 {
		return nil, errors.New("容量必须大于0")
	}

	// 默认配置
	opts := fifoCacheOptions{
		concurrentSafe: true,
	}

	// 应用用户提供的配置选项
	for _, opt := range options {
		opt(&opts)
	}

	return &FIFOCache[K, V]{
		cache:          make(map[K]cacheEntry[K, V], capacity),
		queue:          list.New(),
		capacity:       capacity,
		concurrentSafe: opts.concurrentSafe,
	}, nil
}

// Get 从缓存中获取键对应的值
// 注意：FIFO策略中，Get操作不会改变元素的淘汰顺序
// 参数:
//
//	key: 要查找的键
//
// 返回值:
//
//	value: 键对应的值，如果键不存在则返回V类型的零值
//	exists: 布尔值，表示键是否存在于缓存中
func (f *FIFOCache[K, V]) Get(key K) (V, bool) {
	// 如果启用并发安全，加读锁
	if f.concurrentSafe {
		f.mu.RLock()
		defer f.mu.RUnlock()
	}

	entry, ok := f.cache[key]
	if !ok {
		var zero V
		return zero, false
	}

	return entry.value, true
}

// Set 将键值对存入缓存
// 如果键已存在，仅更新值而不改变其在队列中的位置
// 如果键不存在且缓存已满，会先移除最早插入的键（队列头部元素），再插入新键值对
// 参数:
//
//	key: 要存储的键
//	value: 要存储的值
func (f *FIFOCache[K, V]) Set(key K, value V) {
	// 如果启用并发安全，加写锁
	if f.concurrentSafe {
		f.mu.Lock()
		defer f.mu.Unlock()
	}

	// 检查键是否已存在
	if entry, ok := f.cache[key]; ok {
		// 更新值
		entry.value = value
		f.cache[key] = entry
		return
	}

	// 如果缓存已满，淘汰最早的元素
	if f.queue.Len() >= f.capacity {
		// 移除链表头部元素
		front := f.queue.Front()
		if front != nil {
			oldKey := front.Value.(K)
			// 从哈希表中删除
			delete(f.cache, oldKey)
			// 从链表中删除
			f.queue.Remove(front)
		}
	}

	// 添加新元素到链表尾部
	node := f.queue.PushBack(key)
	// 添加到哈希表
	f.cache[key] = cacheEntry[K, V]{
		value: value,
		node:  node,
	}
}

// Delete 从缓存中删除指定键
// 如果键不存在，此操作无效果
// 参数:
//
//	key: 要删除的键
func (f *FIFOCache[K, V]) Delete(key K) bool {
	// 如果启用并发安全，加写锁
	if f.concurrentSafe {
		f.mu.Lock()
		defer f.mu.Unlock()
	}

	entry, ok := f.cache[key]
	if !ok {
		return false
	}

	// 从链表中删除节点
	f.queue.Remove(entry.node)
	// 从哈希表中删除
	delete(f.cache, key)

	return true
}

// Len 返回当前缓存中的元素数量
// 返回值:
//
//	int: 缓存中已存储的键值对数量
func (f *FIFOCache[K, V]) Len() int {
	if f.concurrentSafe {
		f.mu.RLock()
		defer f.mu.RUnlock()
	}

	return f.queue.Len()
}

// Clear 清空缓存中的所有元素
// 此操作会重置缓存的内部状态，包括哈希表和队列
func (f *FIFOCache[K, V]) Clear() {
	if f.concurrentSafe {
		f.mu.Lock()
		defer f.mu.Unlock()
	}

	// 重置哈希表
	f.cache = make(map[K]cacheEntry[K, V], f.capacity)
	// 重置链表
	f.queue.Init()
}
