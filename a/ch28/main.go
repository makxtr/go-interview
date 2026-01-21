package cache

import (
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

type Node struct {
	key   string
	value interface{}
}

//
// LRU Cache Implementation
//

type LRUNode struct {
	Node
	prev *LRUNode
	next *LRUNode
}

type LRUCache struct {
	capacity int
	cache    map[string]*LRUNode
	head     *LRUNode // Most recently used
	tail     *LRUNode // Least recently used
	hits     uint64
	misses   uint64
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		return nil
	}
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*LRUNode),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	node, ok := c.cache[key]
	if !ok {
		c.misses++
		return nil, false
	}
	c.hits++
	c.moveToFront(node)
	return node.value, true
}

func (c *LRUCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}

	if node, ok := c.cache[key]; ok {
		node.value = value
		c.moveToFront(node)
		return
	}

	if len(c.cache) >= c.capacity {
		delete(c.cache, c.tail.key)
		c.removeNode(c.tail)
	}

	newNode := &LRUNode{Node: Node{key: key, value: value}}
	c.cache[key] = newNode
	c.addNode(newNode)
}

func (c *LRUCache) Delete(key string) bool {
	node, ok := c.cache[key]
	if !ok {
		return false
	}
	delete(c.cache, key)
	c.removeNode(node)
	return true
}

func (c *LRUCache) Clear() {
	c.cache = make(map[string]*LRUNode)
	c.head = nil
	c.tail = nil
	c.hits = 0
	c.misses = 0
}

func (c *LRUCache) Size() int {
	return len(c.cache)
}

func (c *LRUCache) Capacity() int {
	return c.capacity
}

func (c *LRUCache) HitRate() float64 {
	totalReq := c.hits + c.misses
	if totalReq == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(totalReq)
}

// LRU helper methods
func (c *LRUCache) addNode(node *LRUNode) {
	node.prev = nil
	node.next = c.head
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node
	if c.tail == nil {
		c.tail = node
	}
}

func (c *LRUCache) removeNode(node *LRUNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}
}

func (c *LRUCache) moveToFront(node *LRUNode) {
	if node == c.head {
		return
	}
	c.removeNode(node)
	c.addNode(node)
}

//
// LFU Cache Implementation
//

type LFUNode struct {
	Node
	freq int
	prev *LFUNode
	next *LFUNode
}

type LFUDoublyLinkedList struct {
	head *LFUNode
	tail *LFUNode
	len  int
}

func NewLFUDoublyLinkedList() *LFUDoublyLinkedList {
	return &LFUDoublyLinkedList{}
}

func (l *LFUDoublyLinkedList) Add(node *LFUNode) {
	node.prev = nil
	node.next = nil
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		node.next = l.head
		l.head.prev = node
		l.head = node
	}
	l.len++
}

func (l *LFUDoublyLinkedList) Remove(node *LFUNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
	node.prev = nil
	node.next = nil
	l.len--
}

func (l *LFUDoublyLinkedList) RemoveTail() *LFUNode {
	if l.tail == nil {
		return nil
	}
	node := l.tail
	l.Remove(node)
	return node
}

type LFUCache struct {
	capacity   int
	minFreq    int
	cache      map[string]*LFUNode
	freqGroups map[int]*LFUDoublyLinkedList
	hits       uint64
	misses     uint64
}

func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
		return nil
	}
	return &LFUCache{
		capacity:   capacity,
		cache:      make(map[string]*LFUNode),
		freqGroups: make(map[int]*LFUDoublyLinkedList),
	}
}

func (c *LFUCache) Get(key string) (interface{}, bool) {
	node, ok := c.cache[key]
	if !ok {
		c.misses++
		return nil, false
	}
	c.hits++
	c.updateNodeFreq(node)
	return node.value, true
}

func (c *LFUCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}
	if node, ok := c.cache[key]; ok {
		node.value = value
		c.updateNodeFreq(node)
		return
	}
	if len(c.cache) >= c.capacity {
		list := c.freqGroups[c.minFreq]
		nodeToEvict := list.RemoveTail()
		if nodeToEvict != nil {
			delete(c.cache, nodeToEvict.key)
		}
	}
	newNode := &LFUNode{Node: Node{key: key, value: value}, freq: 1}
	c.cache[key] = newNode
	c.minFreq = 1
	list, exists := c.freqGroups[1]
	if !exists {
		list = NewLFUDoublyLinkedList()
		c.freqGroups[1] = list
	}
	list.Add(newNode)
}

func (c *LFUCache) Delete(key string) bool {
	node, ok := c.cache[key]
	if !ok {
		return false
	}
	list := c.freqGroups[node.freq]
	list.Remove(node)
	delete(c.cache, key)
	return true
}

func (c *LFUCache) Clear() {
	c.cache = make(map[string]*LFUNode)
	c.freqGroups = make(map[int]*LFUDoublyLinkedList)
	c.minFreq = 0
	c.hits = 0
	c.misses = 0
}

func (c *LFUCache) Size() int {
	return len(c.cache)
}

func (c *LFUCache) Capacity() int {
	return c.capacity
}

func (c *LFUCache) HitRate() float64 {
	totalReq := c.hits + c.misses
	if totalReq == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(totalReq)
}

func (c *LFUCache) updateNodeFreq(node *LFUNode) {
	oldFreq := node.freq
	list := c.freqGroups[oldFreq]
	list.Remove(node)
	if oldFreq == c.minFreq && list.len == 0 {
		c.minFreq++
	}
	newFreq := oldFreq + 1
	node.freq = newFreq
	newList, exists := c.freqGroups[newFreq]
	if !exists {
		newList = NewLFUDoublyLinkedList()
		c.freqGroups[newFreq] = newList
	}
	newList.Add(node)
}

//
// FIFO Cache Implementation
//

type FIFONode struct {
	Node
	next *FIFONode
}

type FIFOCache struct {
	capacity int
	cache    map[string]*FIFONode
	head     *FIFONode
	tail     *FIFONode
	hits     uint64
	misses   uint64
}

func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		return nil
	}
	return &FIFOCache{
		capacity: capacity,
		cache:    make(map[string]*FIFONode),
	}
}

func (c *FIFOCache) evict() {
	if c.head == nil {
		return
	}
	delete(c.cache, c.head.key)
	c.head = c.head.next
	if c.head == nil {
		c.tail = nil
	}
}

func (c *FIFOCache) Get(key string) (interface{}, bool) {
	if n, ok := c.cache[key]; ok {
		c.hits++
		return n.value, true
	}
	c.misses++
	return nil, false
}

func (c *FIFOCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}
	if n, ok := c.cache[key]; ok {
		n.value = value
		return
	}
	if len(c.cache) >= c.capacity {
		c.evict()
	}
	newNode := &FIFONode{Node: Node{key: key, value: value}}
	c.cache[key] = newNode
	if c.head == nil {
		c.head = newNode
		c.tail = newNode
	} else {
		c.tail.next = newNode
		c.tail = newNode
	}
}

func (c *FIFOCache) Delete(key string) bool {
	nodeToDelete, ok := c.cache[key]
	if !ok {
		return false
	}
	delete(c.cache, key)
	if nodeToDelete == c.head {
		c.head = c.head.next
		if c.head == nil {
			c.tail = nil
		}
		return true
	}
	prev := c.head
	for prev != nil && prev.next != nodeToDelete {
		prev = prev.next
	}
	if prev != nil {
		prev.next = nodeToDelete.next
		if nodeToDelete == c.tail {
			c.tail = prev
		}
	}
	return true
}

func (c *FIFOCache) Clear() {
	c.cache = make(map[string]*FIFONode)
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
	if totalReq == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(totalReq)
}

//
// Thread-Safe Cache Wrapper
//

type ThreadSafeCache struct {
	cache Cache
	mu    sync.RWMutex
}

func NewThreadSafeCache(cache Cache) *ThreadSafeCache {
	return &ThreadSafeCache{cache: cache}
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Get(key)
}

func (c *ThreadSafeCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Put(key, value)
}

func (c *ThreadSafeCache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Delete(key)
}

func (c *ThreadSafeCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Clear()
}

func (c *ThreadSafeCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Size()
}

func (c *ThreadSafeCache) Capacity() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Capacity()
}

func (c *ThreadSafeCache) HitRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.HitRate()
}

//
// Cache Factory Functions
//

func NewCache(policy CachePolicy, capacity int) Cache {
	switch policy {
	case LRU:
		return NewLRUCache(capacity)
	case LFU:
		return NewLFUCache(capacity)
	case FIFO:
		return NewFIFOCache(capacity)
	default:
		// Return LRU as a sensible default
		return NewLRUCache(capacity)
	}
}

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache {
	return NewThreadSafeCache(NewCache(policy, capacity))
}
