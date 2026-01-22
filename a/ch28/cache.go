package cache

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

// --- Common Payload for Nodes ---

type cachePayload struct {
	key   string
	value interface{}
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
