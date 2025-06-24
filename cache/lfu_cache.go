package cache

import (
	ctl "container/list"
	"errors"
)

type lfuNode[K comparable, V any] struct {
	key   K
	value V
	freq  int
	elem  *ctl.Element
}

type LFUCache[K comparable, V any] struct {
	cache    map[K]*lfuNode[K, V]
	freqMap  map[int]*ctl.List
	minFreq  int
	capacity int
}

// NewLFUCache 创建新的LFU缓存实例
// capacity为缓存容量，必须大于0，否则返回错误
// 返回值:
//   *LFUCache[K, V]: 成功创建的缓存实例
//   error: 当capacity <= 0时返回非nil错误
func NewLFUCache[K comparable, V any](capacity int) (*LFUCache[K, V], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}
	return &LFUCache[K, V]{
		cache:    make(map[K]*lfuNode[K, V]),
		freqMap:  make(map[int]*ctl.List),
		capacity: capacity,
	}, nil
}

// Get 实现Cache接口的Get方法
func (l *LFUCache[K, V]) Get(key K) (value V, exists bool) {
	node, exists := l.cache[key]
	if !exists {
		return value, false
	}

	// 更新频率
	l.updateFreq(node)
	return node.value, true
}

// Set 实现Cache接口的Set方法
func (l *LFUCache[K, V]) Set(key K, value V) {
	if node, exists := l.cache[key]; exists {
		node.value = value
		l.updateFreq(node)
		return
	}

	// 如果缓存满，删除最低频率的节点
	if len(l.cache) >= l.capacity {
		l.evict()
	}

	// 创建新节点
	newNode := &lfuNode[K, V]{
		key:   key,
		value: value,
		freq:  1,
	}
	l.cache[key] = newNode

	// 添加到频率1的列表
	if _, ok := l.freqMap[1]; !ok {
		l.freqMap[1] = ctl.New()
	}
	elem := l.freqMap[1].PushBack(newNode)
	newNode.elem = elem

	// 更新最小频率为1
	l.minFreq = 1
}

// Delete 实现Cache接口的Delete方法
func (l *LFUCache[K, V]) Delete(key K) {
	node, exists := l.cache[key]
	if !exists {
		return
	}

	// 从频率列表中删除
	list := l.freqMap[node.freq]
	list.Remove(node.elem)
	if list.Len() == 0 {
		delete(l.freqMap, node.freq)
		// 如果删除的是最小频率的列表，更新minFreq
		if node.freq == l.minFreq {
			l.minFreq++
		}
	}

	// 从缓存中删除
	delete(l.cache, key)
}

// Len 实现Cache接口的Len方法
func (l *LFUCache[K, V]) Len() int {
	return len(l.cache)
}

// Clear 实现Cache接口的Clear方法
func (l *LFUCache[K, V]) Clear() {
	l.cache = make(map[K]*lfuNode[K, V])
	l.freqMap = make(map[int]*ctl.List)
	l.minFreq = 0
}

// updateFreq 更新节点的访问频率
// 实现逻辑：
// 1. 从旧频率链表中移除节点
// 2. 增加节点频率
// 3. 将节点添加到新频率链表
// 4. 如果旧频率是最小频率且链表为空，更新minFreq
func (l *LFUCache[K, V]) updateFreq(node *lfuNode[K, V]) {
	oldFreq := node.freq
	node.freq++

	// 从旧频率列表中删除
	oldList := l.freqMap[oldFreq]
	oldList.Remove(node.elem)
	if oldList.Len() == 0 {
		delete(l.freqMap, oldFreq)
		// 如果旧频率是最小频率，更新minFreq
		if oldFreq == l.minFreq {
			l.minFreq++
		}
	}

	// 添加到新频率列表
	newFreq := node.freq
	if _, ok := l.freqMap[newFreq]; !ok {
		l.freqMap[newFreq] = ctl.New()
	}
	elem := l.freqMap[newFreq].PushBack(node)
	node.elem = elem
}

// evict 淘汰访问频率最低的节点
// 实现逻辑：
// 1. 获取当前最小频率对应的链表
// 2. 删除链表头部节点（最久未使用）
// 3. 从缓存中删除该节点
// 4. 如果链表为空，删除对应的频率映射
func (l *LFUCache[K, V]) evict() {
	freqList := l.freqMap[l.minFreq]
	if freqList == nil || freqList.Len() == 0 {
		return
	}

	// 删除列表头部节点（最久未使用）
	elem := freqList.Front()
	node := elem.Value.(*lfuNode[K, V])
	freqList.Remove(elem)
	delete(l.cache, node.key)

	// 如果列表为空，删除频率映射
	if freqList.Len() == 0 {
		delete(l.freqMap, l.minFreq)
	}
}
