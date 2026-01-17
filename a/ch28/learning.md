# Learning Materials for Cache Implementation

## Introduction to Caching
Caching is a fundamental technique in computer science used to store frequently accessed data in a fast-access location. A good cache implementation can dramatically improve application performance by reducing the time needed to access data from slower storage systems.

### Why Caching Matters
*   **Performance:** Reduces latency by storing frequently accessed data closer to the application.
*   **Resource Efficiency:** Reduces load on backend systems like databases.
*   **Scalability:** Helps applications handle more requests with the same resources.
*   **Cost Reduction:** Minimizes expensive operations like database queries or API calls.

## Cache Fundamentals
A cache is essentially a key-value store with limited capacity. When the cache reaches its capacity, it must decide which items to remove to make space for new ones. This decision is made by an eviction policy.

## Cache Eviction Policies

### 1. LRU (Least Recently Used)
LRU evicts the item that was accessed (read or written) least recently.

#### Algorithm:
*   Maintain a doubly-linked list of cache entries ordered by access time.
*   Use a hash map for O(1) key lookups.
*   On access: move item to front of list.
*   On eviction: remove item from back of list.

**Time Complexity:** O(1) for all operations.
**Space Complexity:** O(n) where n is cache capacity.

```go
// LRU Cache Implementation Concept
type LRUCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node // Most recently used
	tail     *Node // Least recently used
}

type Node struct {
	key   string
	value interface{}
	prev  *Node
	next  *Node
}
```

#### Use Cases:
*   Operating system page replacement.
*   CPU cache management.
*   Web browser cache.
*   Database buffer pools.

#### Advantages:
*   Good temporal locality performance.
*   Intuitive eviction strategy.
*   Works well for most general-purpose scenarios.

#### Disadvantages:
*   Doesn't consider access frequency.
*   Can be affected by sequential scans that destroy cache locality.

### 2. LFU (Least Frequently Used)
LFU evicts the item that has been accessed the fewest times.

#### Algorithm:
*   Maintain a frequency counter for each item.
*   Use a min-heap or frequency buckets for efficient eviction.
*   On access: increment frequency counter.
*   On eviction: remove item with lowest frequency.

**Time Complexity:** O(1) for get/put with proper implementation.
**Space Complexity:** O(n).

```go
// LFU Cache Implementation Concept
type LFUCache struct {
	capacity   int
	minFreq    int
	cache      map[string]*Node
	freqGroups map[int]*FreqGroup // frequency -> list of nodes
}

type FreqGroup struct {
	freq int
	head *Node
	tail *Node
}
```

#### Use Cases:
*   Long-running applications with stable access patterns.
*   Scientific computing with repeated data access.
*   CDN systems.

#### Advantages:
*   Excellent for workloads with clear hot data.
*   Adapts well to changing access patterns over time.
*   Good for scenarios where some data is accessed much more frequently.

#### Disadvantages:
*   More complex implementation.
*   New items are immediately evicted if cache is full.
*   Frequency counts can become stale over time.

### 3. FIFO (First In, First Out)
FIFO evicts the oldest item in the cache, regardless of access patterns.

#### Algorithm:
*   Maintain insertion order using a queue or linked list.
*   On insertion: add to front.
*   On eviction: remove from back.

**Time Complexity:** O(1) for all operations.
**Space Complexity:** O(n).

```go
// FIFO Cache Implementation Concept
type FIFOCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node // Newest item
	tail     *Node // Oldest item
}
```

#### Use Cases:
*   Simple caching scenarios.
*   When access patterns are unknown.
*   Embedded systems with memory constraints.

#### Advantages:
*   Simple to implement and understand.
*   Predictable behavior.
*   No access pattern tracking needed.

#### Disadvantages:
*   Ignores access patterns completely.
*   May evict frequently used items.
*   Generally poor cache hit rates.

## Advanced Cache Concepts

### Thread Safety
Real-world caches must handle concurrent access:

```go
type ThreadSafeCache struct {
	mu    sync.RWMutex
	cache Cache
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
```

### Cache Metrics
Important metrics to track:

*   **Hit Rate:** Percentage of requests served from cache.
*   **Miss Rate:** Percentage of requests not in cache.
*   **Eviction Count:** Number of items evicted.
*   **Average Response Time:** Performance measurement.

```go
type CacheMetrics struct {
	hits      int64
	misses    int64
	evictions int64
}

func (m *CacheMetrics) HitRate() float64 {
	total := m.hits + m.misses
	if total == 0 {
		return 0
	}
	return float64(m.hits) / float64(total)
}
```

### TTL (Time To Live)
Items can automatically expire after a certain time:

```go
type CacheEntry struct {
	value     interface{}
	timestamp time.Time
	ttl       time.Duration
}

func (e *CacheEntry) IsExpired() bool {
	if e.ttl == 0 {
		return false
	}
	return time.Since(e.timestamp) > e.ttl
}
```

## Implementation Strategies

### Memory Management
```go
// Proper cleanup to prevent memory leaks
func (c *Cache) evict(node *Node) {
	// Remove from hash map
	delete(c.cache, node.key)

	// Remove from linked list
	c.removeFromList(node)

	// Clear references to help GC
	node.prev = nil
	node.next = nil
	node.value = nil
}
```

### Interface Design
Design for flexibility and testability:

```go
type Cache interface {
	Get(key string) (value interface{}, found bool)
	Put(key string, value interface{})
	Delete(key string) bool
	Clear()
	Size() int
	Capacity() int
}

type EvictionPolicy interface {
	OnAccess(key string)
	OnInsert(key string)
	OnDelete(key string)
	SelectVictim() string
}
```

### Performance Optimization

#### Memory Layout
```go
// Use struct of arrays for better cache locality
type OptimizedLRU struct {
	keys     []string
	values   []interface{}
	prev     []int
	next     []int
	keyIndex map[string]int
	head     int
	tail     int
	size     int
	capacity int
}
```

#### Avoiding Allocations
```go
// Pre-allocate node pool to reduce GC pressure
type NodePool struct {
	nodes []Node
	free  []int
}

func (p *NodePool) Get() *Node {
	if len(p.free) == 0 {
		return &Node{}
	}
	idx := p.free[len(p.free)-1]
	p.free = p.free[:len(p.free)-1]
	return &p.nodes[idx]
}
```

## Testing Strategies

### Unit Tests
```go
func TestCacheEviction(t *testing.T) {
	cache := NewLRUCache(2)

	// Fill cache
	cache.Put("a", 1)
	cache.Put("b", 2)

	// Trigger eviction
	cache.Put("c", 3)

	// Verify oldest item was evicted
	_, found := cache.Get("a")
	assert.False(t, found)
}
```

### Concurrency Tests
```go
func TestConcurrentAccess(t *testing.T) {
	cache := NewThreadSafeCache(NewLRUCache(100)) // Assuming NewThreadSafeCache takes a Cache interface

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				cache.Put(key, j)
				cache.Get(key)
			}
		}(i)
	}
	wg.Wait()
}
```

### Benchmark Tests
```go
func BenchmarkCacheGet(b *testing.B) {
	cache := NewLRUCache(1000)

	// Populate cache
	for i := 0; i < 1000; i++ {
		cache.Put(fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("key-%d", i%1000))
	}
}
```

## Real-World Considerations

### Cache Stampede Prevention
```go
type SafeCache struct {
	cache  Cache
	groups singleflight.Group
}

func (c *SafeCache) GetOrCompute(key string, compute func() interface{}) interface{} {
	if value, found := c.cache.Get(key); found {
		return value
	}

	// Prevent cache stampede
	value, _, _ := c.groups.Do(key, func() (interface{}, error) {
		if value, found := c.cache.Get(key); found {
			return value, nil
		}

		computed := compute()
		c.cache.Put(key, computed)
		return computed, nil
	})

	return value
}
```

### Distributed Caching
```go
type DistributedCache interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{})
	Invalidate(key string)
	InvalidatePattern(pattern string)
}
```

### Memory Pressure Handling
```go
func (c *Cache) handleMemoryPressure() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if m.Alloc > c.maxMemory {
		// Aggressively evict items
		evictCount := c.size / 4
		for i := 0; i < evictCount; i++ {
			c.evictLRU() // Assuming an evictLRU method exists
		}
	}
}
```

## Comparison of Cache Policies

| Policy | Get Time | Put Time | Space | Best Use Case                       |
| :----- | :------- | :------- | :---- | :---------------------------------- |
| LRU    | O(1)     | O(1)     | O(n)  | General purpose, temporal locality  |
| LFU    | O(1)     | O(1)     | O(n)  | Stable access patterns, hot data    |
| FIFO   | O(1)     | O(1)     | O(n)  | Simple scenarios, unknown patterns  |

## Further Reading
*   [Cache Replacement Policies](https://en.wikipedia.org/wiki/Cache_replacement_policies)
*   [The LFU-DA Cache Algorithm](https://dl.acm.org/doi/10.1145/1273440.1273441)
*   [Caching at Scale with Redis](https://redis.io/topics/lru-cache)
*   [Linux Kernel Page Cache](https://www.kernel.org/doc/html/latest/admin-guide/mm/concepts.html#page-cache)
*   [Caffeine: A High Performance Java Caching Library](https://github.com/ben-manes/caffeine)
