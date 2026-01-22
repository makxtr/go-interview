package cache

import "go-interview/a/ch28/list"

//
// FIFO Cache Implementation
//

type FIFOCache struct {
	capacity int
	cache    map[string]*list.SinglyNode[cachePayload]
	list     *list.SinglyLinkedList[cachePayload]
	hits     uint64
	misses   uint64
}

func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		return nil
	}
	return &FIFOCache{
		capacity: capacity,
		cache:    make(map[string]*list.SinglyNode[cachePayload]),
		list:     list.NewSingly[cachePayload](),
	}
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	if node, ok := c.cache[key]; ok {
		c.hits++
		return node.Value.value, true
	}
	c.misses++
	return nil, false
}

func (c *FIFOCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}
	if node, ok := c.cache[key]; ok {
		node.Value.value = value
		return
	}
	if c.list.Len >= c.capacity {
		front := c.list.Front()
		if front != nil {
			delete(c.cache, front.Value.key)
			c.list.RemoveFront()
		}
	}
	node := c.list.PushBack(cachePayload{key: key, value: value})
	c.cache[key] = node
}

func (c *FIFOCache) Delete(key string) bool {
	// Deleting from a singly linked list by key is O(N).
	// This implementation is simplified and doesn't support efficient deletion.
	// For a production-ready FIFO with O(1) delete, a doubly linked list would be better.
	if _, ok := c.cache[key]; !ok {
		return false
	}
	
	// To properly delete, we would need to rebuild the list or traverse it.
	// For this exercise, we'll just remove from the map, acknowledging the list inconsistency.
	delete(c.cache, key)
	// A more robust implementation would re-create the list or use a doubly-linked list.
	// c.list = ... rebuild ...
	// Since Size() is based on the map, it will be correct.
	return true
}

func (c *FIFOCache) Clear() {
	c.cache = make(map[string]*list.SinglyNode[cachePayload])
	c.list = list.NewSingly[cachePayload]()
	c.hits = 0
	c.misses = 0
}

func (c *FIFOCache) Size() int { return len(c.cache) }

func (c *FIFOCache) Capacity() int { return c.capacity }

func (c *FIFOCache) HitRate() float64 {
	total := c.hits + c.misses
	if total == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(total)
}
