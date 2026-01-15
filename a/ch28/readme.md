# Challenge 28: Cache Implementation with Multiple Eviction Policies

## Problem Statement
In this challenge, you will implement a high-performance, thread-safe cache system with multiple eviction policies. This is a common interview question that tests your understanding of data structures, algorithms, concurrency, and system design.

Your task is to implement three different cache eviction policies:

*   **LRU (Least Recently Used):** Evicts the item that was accessed least recently.
*   **LFU (Least Frequently Used):** Evicts the item that has been accessed the fewest times.
*   **FIFO (First In, First Out):** Evicts the oldest item regardless of access patterns.

Each implementation must provide O(1) time complexity for `Get` and `Put` operations and must be thread-safe.

## Requirements

### Core Interface
All cache implementations must satisfy this interface:

```go
type Cache interface {
	// Get retrieves a value by key. Returns the value and true if found, or nil and false if not found.
	Get(key string) (value interface{}, found bool)

	// Put stores a key-value pair. If the cache is at capacity, it should evict according to its policy.
	Put(key string, value interface{})

	// Delete removes a key-value pair. Returns true if the key existed, false otherwise.
	Delete(key string) bool

	// Clear removes all entries from the cache.
	Clear()

	// Size returns the current number of items in the cache.
	Size() int

	// Capacity returns the maximum number of items the cache can hold.
	Capacity() int

	// HitRate returns the cache hit rate as a float between 0 and 1.
	HitRate() float64
}
```

### Performance Requirements

*   **Time Complexity:** O(1) for `Get`, `Put`, and `Delete` operations.
*   **Space Complexity:** O(n) where n is the cache capacity.
*   **Thread Safety:** All operations must be safe for concurrent use.
*   **Memory Efficiency:** Minimize memory overhead and prevent memory leaks.

### Implementation Requirements
You must implement:

*   **LRUCache:** Uses doubly-linked list + hash map.
*   **LFUCache:** Uses frequency tracking with efficient eviction.
*   **FIFOCache:** Uses queue-based eviction.
*   **ThreadSafeWrapper:** Makes any cache implementation thread-safe.
*   **CacheFactory:** Creates cache instances based on policy type.

## Function Signatures

### 1. LRU Cache

```go
type LRUCache struct {
	// Private fields - design your own implementation
}

// NewLRUCache creates a new LRU cache with the specified capacity
func NewLRUCache(capacity int) *LRUCache

// Implement Cache interface methods
func (c *LRUCache) Get(key string) (interface{}, bool)
func (c *LRUCache) Put(key string, value interface{})
func (c *LRUCache) Delete(key string) bool
func (c *LRUCache) Clear()
func (c *LRUCache) Size() int
func (c *LRUCache) Capacity() int
func (c *LRUCache) HitRate() float64
```

### 2. LFU Cache

```go
type LFUCache struct {
	// Private fields - design your own implementation
}

// NewLFUCache creates a new LFU cache with the specified capacity
func NewLFUCache(capacity int) *LFUCache

// Implement Cache interface methods
func (c *LFUCache) Get(key string) (interface{}, bool)
func (c *LFUCache) Put(key string, value interface{})
func (c *LFUCache) Delete(key string) bool
func (c *LFUCache) Clear()
func (c *LFUCache) Size() int
func (c *LFUCache) Capacity() int
func (c *LFUCache) HitRate() float64
```

### 3. FIFO Cache

```go
type FIFOCache struct {
	// Private fields - design your own implementation
}

// NewFIFOCache creates a new FIFO cache with the specified capacity
func NewFIFOCache(capacity int) *FIFOCache

// Implement Cache interface methods
func (c *FIFOCache) Get(key string) (interface{}, bool)
func (c *FIFOCache) Put(key string, value interface{})
func (c *FIFOCache) Delete(key string) bool
func (c *FIFOCache) Clear()
func (c *FIFOCache) Size() int
func (c *FIFOCache) Capacity() int
func (c *FIFOCache) HitRate() float64
```

### 4. Thread-Safe Wrapper

```go
type ThreadSafeCache struct {
	// Private fields - design your own implementation
}

// NewThreadSafeCache wraps any cache implementation to make it thread-safe
func NewThreadSafeCache(cache Cache) *ThreadSafeCache

// Implement Cache interface methods with proper locking
func (c *ThreadSafeCache) Get(key string) (interface{}, bool)
func (c *ThreadSafeCache) Put(key string, value interface{})
func (c *ThreadSafeCache) Delete(key string) bool
func (c *ThreadSafeCache) Clear()
func (c *ThreadSafeCache) Size() int
func (c *ThreadSafeCache) Capacity() int
func (c *ThreadSafeCache) HitRate() float64
```

### 5. Cache Factory

```go
type CachePolicy int

const (
	LRU CachePolicy = iota
	LFU
	FIFO
)

// NewCache creates a cache with the specified policy and capacity
func NewCache(policy CachePolicy, capacity int) Cache

// NewThreadSafeCacheWithPolicy creates a thread-safe cache with the specified policy
func NewThreadSafeCacheWithPolicy(policy CachePolicy, capacity int) Cache
```

## Input/Output Examples

### LRU Cache Example

```go
cache := NewLRUCache(2)

cache.Put("a", 1)
cache.Put("b", 2)
fmt.Println(cache.Get("a"))  // Output: 1, true

cache.Put("c", 3)            // Evicts "b" (least recently used)
fmt.Println(cache.Get("b"))  // Output: nil, false
fmt.Println(cache.Get("a"))  // Output: 1, true
fmt.Println(cache.Get("c"))  // Output: 3, true
```

### LFU Cache Example

```go
cache := NewLFUCache(2)

cache.Put("a", 1)
cache.Put("b", 2)
cache.Get("a")               // "a" now has frequency 2
cache.Get("a")               // "a" now has frequency 3

cache.Put("c", 3)            // Evicts "b" (frequency 1, least frequent)
fmt.Println(cache.Get("b"))  // Output: nil, false
fmt.Println(cache.Get("a"))  // Output: 1, true
fmt.Println(cache.Get("c"))  // Output: 3, true
```

### FIFO Cache Example

```go
cache := NewFIFOCache(2)

cache.Put("a", 1)
cache.Put("b", 2)
cache.Get("a")               // Doesn't affect eviction order

cache.Put("c", 3)            // Evicts "a" (first in, first out)
fmt.Println(cache.Get("a"))  // Output: nil, false
fmt.Println(cache.Get("b"))  // Output: 2, true
fmt.Println(cache.Get("c"))  // Output: 3, true
```

### Thread-Safe Example

```go
cache := NewThreadSafeCache(NewLRUCache(100))

// Safe for concurrent use
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
	wg.Add(1)
	go func(id int) {
		defer wg.Done()
		for j := 0; j < 100; j++ {
			key := fmt.Sprintf("key-%d-%d", id, j)
			cache.Put(key, j)
			cache.Get(key)
		}
	}(i)
}
wg.Wait()

fmt.Printf("Hit rate: %.2f\n", cache.HitRate())
```

## Evaluation Criteria

*   **Correctness (40 points)**
    *   All cache policies work correctly.
    *   Proper eviction behavior.
    *   Thread safety when wrapped.
    *   Edge cases handled properly.
*   **Performance (25 points)**
    *   O(1) time complexity for all operations.
    *   Efficient memory usage.
    *   Minimal lock contention in thread-safe version.
*   **Code Quality (20 points)**
    *   Clean, readable, and well-organized code.
    *   Proper abstraction and interfaces.
    *   Good variable and function naming.
    *   Appropriate comments.
*   **Algorithm Understanding (15 points)**
    *   Efficient implementation of each eviction policy.
    *   Understanding of data structure trade-offs.
    *   Proper handling of concurrent access patterns.

## Advanced Requirements (Bonus)
Implement these for extra credit:

*   **TTL Support:** Add time-to-live functionality.
*   **Cache Metrics:** Detailed statistics (hit rate, eviction count, etc.).
*   **Benchmark Tests:** Performance comparison between policies.
*   **Memory Pressure Handling:** Automatic eviction under memory pressure.

## Constraints

*   Cache capacity must be at least 1.
*   Keys are always non-empty strings.
*   Values can be any `interface{}` type.
*   Must handle `nil` values correctly.
*   Zero-capacity cache should always miss.
*   Thread-safe operations should use minimal locking.

## Hints

*   **LRU:** Use a doubly-linked list with a hash map pointing to nodes.
*   **LFU:** Consider using frequency buckets or a min-heap for efficient victim selection.
*   **FIFO:** A simple queue or circular buffer works well.
*   **Thread Safety:** Use `sync.RWMutex` for better read performance.
*   **Memory Management:** Be careful about memory leaks when removing nodes.
*   **Testing:** Test edge cases like capacity 1, concurrent access, and large datasets.

## Time Limit
This challenge should be completed within 60-90 minutes for an interview setting.

## Success Criteria
A successful implementation should:

*   Pass all provided test cases.
*   Demonstrate O(1) performance for basic operations.
*   Handle concurrent access safely.
*   Show understanding of different eviction policies.
*   Include proper error handling and edge case management.
