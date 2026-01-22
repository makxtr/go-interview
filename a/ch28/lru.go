package cache

import "go-interview/a/ch28/list"

//
// LRU Cache Implementation
//

type LRUCache struct {
	capacity int
	cache    map[string]*list.DoublyNode[cachePayload]
	list     *list.DoublyLinkedList[cachePayload]
	hits     uint64
	misses   uint64
}

func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		return nil
	}
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.DoublyNode[cachePayload]),
		list:     list.NewDoubly[cachePayload](),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	node, ok := c.cache[key]
	if !ok {
		c.misses++
		return nil, false
	}
	c.hits++
	c.list.MoveToFront(node)
	return node.Value.value, true
}

func (c *LRUCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}
	if node, ok := c.cache[key]; ok {
		node.Value.value = value
		c.list.MoveToFront(node)
		return
	}
	if c.list.Len >= c.capacity {
		tail := c.list.Back()
		if tail != nil {
			delete(c.cache, tail.Value.key)
			c.list.Remove(tail)
		}
	}
	node := c.list.PushFront(cachePayload{key: key, value: value})
	c.cache[key] = node
}

func (c *LRUCache) Delete(key string) bool {
	node, ok := c.cache[key]
	if !ok {
		return false
	}
	delete(c.cache, key)
	c.list.Remove(node)
	return true
}

func (c *LRUCache) Clear() {
	c.cache = make(map[string]*list.DoublyNode[cachePayload])
	c.list = list.NewDoubly[cachePayload]()
	c.hits = 0
	c.misses = 0
}

func (c *LRUCache) Size() int { return c.list.Len }

func (c *LRUCache) Capacity() int { return c.capacity }

func (c *LRUCache) HitRate() float64 {
	total := c.hits + c.misses
	if total == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(total)
}
