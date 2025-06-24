package cache

type Cache[K comparable, V any] interface {
	// Get 获取缓存中key对应的值，如果不存在返回false
	Get(key K) (value V, exists bool)
	// Set 将key-value存入缓存，如果缓存满则根据策略驱逐旧元素
	Set(key K, value V)
	// Delete 从缓存中删除key对应的元素
	Delete(key K)
	// Len 返回当前缓存中的元素数量
	Len() int
	// Clear 清空缓存中的所有元素
	Clear()
}