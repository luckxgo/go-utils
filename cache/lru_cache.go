package cache

import (
	"container/list"
	"errors"
	"sync"
)

// LRUCache 基于最近最久未使用(Least Recently Used)策略的缓存实现
// 当缓存达到容量限制时，最久未被访问的元素会被优先淘汰
// 每次访问(包括Get和Set)会将元素标记为最近使用
// K为键类型，必须支持比较操作；V为值类型，可以是任意类型
type LRUCache[K comparable, V any] struct {
	cache          map[K]*list.Element // 键到链表元素的映射，提供O(1)时间复杂度的访问
	list           *list.List          // 维护访问顺序的双向链表，越靠近头部越是最近访问的元素
	capacity       int                 // 缓存的最大容量，超过此容量将触发淘汰机制
	concurrentSafe bool                // 是否启用并发安全模式
	mu             sync.RWMutex        // 读写锁，在并发安全模式下使用
}

// entry 链表节点存储的数据结构
// 包含键和值，用于在淘汰链表尾部元素时从map中删除对应条目
type entry[K comparable, V any] struct {
	key   K  // 缓存键
	value V  // 缓存值
}

// LRUOption 定义LRU缓存的配置选项函数类型
type LRUOption func(*lruCacheOptions)

// lruCacheOptions LRU缓存的配置选项
// 目前仅支持并发安全开关
// 后续可扩展添加其他配置项
type lruCacheOptions struct {
	concurrentSafe bool
}

// WithLRUConcurrentSafe 设置是否启用并发安全
// 参数safe为true时启用并发安全，使用读写锁保护所有操作
// 参数safe为false时禁用并发安全，性能更高但不保证线程安全
func WithLRUConcurrentSafe(safe bool) LRUOption {
	return func(opts *lruCacheOptions) {
		opts.concurrentSafe = safe
	}
}

// NewLRUCache 创建新的LRU缓存实例
// capacity为缓存容量，必须大于0，否则返回错误
// options为可选配置参数，可通过WithLRUConcurrentSafe等函数设置
// 返回值:
//   *LRUCache[K, V]: 成功创建的缓存实例
//   error: 当capacity <= 0时返回非nil错误
func NewLRUCache[K comparable, V any](capacity int, options ...LRUOption) (*LRUCache[K, V], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}

	// 默认配置
	opts := lruCacheOptions{
		concurrentSafe: true, // 默认启用并发安全
	}
	for _, opt := range options {
		opt(&opts)
	}

	return &LRUCache[K, V]{
		cache:          make(map[K]*list.Element),
		list:           list.New(),
		capacity:       capacity,
		concurrentSafe: opts.concurrentSafe,
	}, nil
}

// Get 从缓存中获取键对应的值
// 如果键存在，会将该键标记为最近使用(移到链表头部)并返回值
// 参数:
//   key: 要查找的键
// 返回值:
//   value: 键对应的值，如果键不存在则返回V类型的零值
//   exists: 布尔值，表示键是否存在于缓存中
func (l *LRUCache[K, V]) Get(key K) (value V, exists bool) {
	if l.concurrentSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	elem, exists := l.cache[key]
	if !exists {
		return value, false
	}

	// 将访问的元素移到链表头部（标记为最近使用）
	l.list.MoveToFront(elem)
	return elem.Value.(*entry[K, V]).value, true
}

// Set 将键值对存入缓存
// 如果键已存在，更新值并将该键标记为最近使用(移到链表头部)
// 如果键不存在且缓存已满，会先移除最久未使用的元素(链表尾部)，再插入新元素
// 参数:
//   key: 要存储的键
//   value: 要存储的值
func (l *LRUCache[K, V]) Set(key K, value V) {
	if l.concurrentSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	// 如果键已存在，更新值并移到头部
	if elem, exists := l.cache[key]; exists {
		elem.Value.(*entry[K, V]).value = value
		l.list.MoveToFront(elem)
		return
	}

	// 如果缓存满，移除链表尾部元素（最久未使用）
	if l.list.Len() >= l.capacity {
		backElem := l.list.Back()
		if backElem != nil {
			// 从map中删除对应的键
			delete(l.cache, backElem.Value.(*entry[K, V]).key)
			// 从链表中删除尾部元素
			l.list.Remove(backElem)
		}
	}

	// 创建新节点并添加到链表头部
	newElem := l.list.PushFront(&entry[K, V]{key: key, value: value})
	l.cache[key] = newElem
}

// Delete 从缓存中删除指定键
// 如果键不存在，此操作无效果
// 参数:
//   key: 要删除的键
func (l *LRUCache[K, V]) Delete(key K) {
	if l.concurrentSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	elem, exists := l.cache[key]
	if !exists {
		return
	}

	// 从链表中删除元素
	l.list.Remove(elem)
	// 从map中删除键
	delete(l.cache, key)
}

// Len 返回当前缓存中的元素数量
// 返回值:
//   int: 缓存中已存储的键值对数量
func (l *LRUCache[K, V]) Len() int {
	if l.concurrentSafe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	return l.list.Len()
}

// Clear 清空缓存中的所有元素
// 此操作会重置缓存的内部状态，包括哈希表和双向链表
func (l *LRUCache[K, V]) Clear() {
	if l.concurrentSafe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	l.list.Init()
	l.cache = make(map[K]*list.Element)
}