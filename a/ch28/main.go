package cache

import (
	"fmt"
	"sync"
)

// Cache interface defines the contract for all cache implementations
type Cache interface {
	Get(key string) (value interface{}, found bool)
	Put(key string, value interface{})
	Delete(key string) bool
	Clear()
	Size() int
	Capacity() int
	HitRate() float64
}

// CachePolicy represents the eviction policy type
type CachePolicy int

const (
	LRU CachePolicy = iota
	LFU
	FIFO
)

//
// LRU Cache Implementation
//

type LRUCache struct {
	// TODO: Add necessary fields for LRU implementation
	// Hint: Use a doubly-linked list + hash map
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	// TODO: Implement LRU cache constructor
	return nil
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	// TODO: Implement LRU get operation
	// Should move accessed item to front (most recently used position)
	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
	// TODO: Implement LRU put operation
	// Should add new item to front and evict least recently used if at capacity
}

func (c *LRUCache) Delete(key string) bool {
	// TODO: Implement delete operation
	return false
}

func (c *LRUCache) Clear() {
	// TODO: Implement clear operation
}

func (c *LRUCache) Size() int {
	// TODO: Return current cache size
	return 0
}

func (c *LRUCache) Capacity() int {
	// TODO: Return cache capacity
	return 0
}

func (c *LRUCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	return 0.0
}

//
// LFU Cache Implementation
//

type LFUCache struct {
	// TODO: Add necessary fields for LFU implementation
	// Hint: Use frequency tracking with efficient eviction
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache {
	// TODO: Implement LFU cache constructor
	return nil
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	// TODO: Implement LFU get operation
	// Should increment frequency count of accessed item
	return nil, false
}

func (c *LFUCache) Put(key string, value interface{}) {
	// TODO: Implement LFU put operation
	// Should evict least frequently used item if at capacity
}

func (c *LFUCache) Delete(key string) bool {
	// TODO: Implement delete operation
	return false
}

func (c *LFUCache) Clear() {
	// TODO: Implement clear operation
}

func (c *LFUCache) Size() int {
	// TODO: Return current cache size
	return 0
}

func (c *LFUCache) Capacity() int {
	// TODO: Return cache capacity
	return 0
}

func (c *LFUCache) HitRate() float64 {
	// TODO: Calculate and return hit rate
	return 0.0
}

//
// FIFO Cache Implementation
//
/*
type Node struct {
	key   string
	value interface{}
	next  *Node
}
*/

type Node struct {
	key   string
	value interface{}
	next  *Node
}

type FIFOCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node // Newest item
	tail     *Node // Oldest item
	hits     uint64
	misses   uint64
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		return nil
	}

	return &FIFOCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
		head:     nil,
		tail:     nil,
		hits:     0,
		misses:   0,
	}
}

func (c *FIFOCache) evict() {
	if nil == c.head {
		return
	}
	delete(c.cache, c.head.key)

	c.head = c.head.next
	if nil == c.head {
		c.tail = nil
	}
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	// Note: Get operations don't affect eviction order in FIFO
	if n, ok := c.cache[key]; ok {
		c.hits++
		return n.value, true
	}
	c.misses++

	return nil, false
}

func (c *FIFOCache) Put(key string, value interface{}) {
	// Should evict first-in item if at capacity
	if n, ok := c.cache[key]; ok {
		n.value = value
		return
	}
	if len(c.cache) == c.capacity {
		c.evict()
	}
	newNode := &Node{key: key, value: value, next: nil}
	c.cache[key] = newNode

	if nil == c.head {
		c.head = newNode
		c.tail = newNode
	} else {
		c.tail.next = newNode
		c.tail = newNode
	}
}

func (c *FIFOCache) Delete(key string) bool {
	// TODO: Implement delete operation
	return false
}

func (c *FIFOCache) Clear() {
	c.cache = make(map[string]*Node)
	c.head = nil
	c.tail = nil

	c.hits = 0
	c.misses = 0
}

func (c *FIFOCache) Size() int {
	return len(c.cache)
}

func (c *FIFOCache) Capacity() int {
	return c.capacity
}

func (c *FIFOCache) HitRate() float64 {
	totalReq := c.hits + c.misses

	if 0 == totalReq {
		return 0.0
	}

	fmt.Println(totalReq)

	return float64(c.hits) / float64(totalReq)
}

//
// Thread-Safe Cache Wrapper
//

type ThreadSafeCache struct {
	cache Cache
	mu    sync.RWMutex
	// TODO: Add any additional fields if needed
}

// NewThreadSafeCache wraps any cache implementation to make it thread-safe
func NewThreadSafeCache(cache Cache) *ThreadSafeCache {
	// TODO: Implement thread-safe wrapper constructor
	return nil
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	// TODO: Implement thread-safe get operation
	// Hint: Use read lock for better performance
	return nil, false
}

func (c *ThreadSafeCache) Put(key string, value interface{}) {
	// TODO: Implement thread-safe put operation
	// Hint: Use write lock
}

func (c *ThreadSafeCache) Delete(key string) bool {
	// TODO: Implement thread-safe delete operation
	return false
}

func (c *ThreadSafeCache) Clear() {
	// TODO: Implement thread-safe clear operation
}

func (c *ThreadSafeCache) Size() int {
	// TODO: Implement thread-safe size operation
	return 0
}

func (c *ThreadSafeCache) Capacity() int {
	// TODO: Implement thread-safe capacity operation
	return 0
}

func (c *ThreadSafeCache) HitRate() float64 {
	// TODO: Implement thread-safe hit rate operation
	return 0.0
}

//
// Cache Factory Functions
//

// NewCache creates a cache with the specified policy and capacity
func NewCache(policy CachePolicy, capacity int) Cache {
	// TODO: Implement cache factory
	// Should create appropriate cache type based on policy
	switch policy {
	case LRU:
		// TODO: Return LRU cache
	case LFU:
		// TODO: Return LFU cache
	case FIFO:
		// TODO: Return FIFO cache
	default:
		// TODO: Return default cache or handle error
	}
	return nil
}

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	// TODO: Implement thread-safe cache factory
	// Should create cache with policy and wrap it with thread safety
	return nil
}
