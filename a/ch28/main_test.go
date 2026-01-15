package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestLRUCache tests the LRU cache implementation
func TestLRUCache(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		cache := NewLRUCache(2)
		if cache == nil {
			t.Fatal("NewLRUCache returned nil")
		}

		// Test initial state
		if cache.Size() != 0 {
			t.Errorf("Expected size 0, got %d", cache.Size())
		}
		if cache.Capacity() != 2 {
			t.Errorf("Expected capacity 2, got %d", cache.Capacity())
		}

		// Test miss
		_, found := cache.Get("a")
		if found {
			t.Error("Expected miss for non-existent key")
		}

		// Test put and get
		cache.Put("a", 1)
		value, found := cache.Get("a")
		if !found || value != 1 {
			t.Errorf("Expected (1, true), got (%v, %v)", value, found)
		}

		if cache.Size() != 1 {
			t.Errorf("Expected size 1, got %d", cache.Size())
		}
	})

	t.Run("LRU Eviction", func(t *testing.T) {
		cache := NewLRUCache(2)

		// Fill cache
		cache.Put("a", 1)
		cache.Put("b", 2)

		// Access "a" to make it most recently used
		cache.Get("a")

		// Add new item, should evict "b" (least recently used)
		cache.Put("c", 3)

		// "b" should be evicted
		_, found := cache.Get("b")
		if found {
			t.Error("Expected 'b' to be evicted")
		}

		// "a" and "c" should still be present
		value, found := cache.Get("a")
		if !found || value != 1 {
			t.Errorf("Expected 'a' to be present with value 1, got (%v, %v)", value, found)
		}

		value, found = cache.Get("c")
		if !found || value != 3 {
			t.Errorf("Expected 'c' to be present with value 3, got (%v, %v)", value, found)
		}
	})

	t.Run("Update Existing Key", func(t *testing.T) {
		cache := NewLRUCache(2)

		cache.Put("a", 1)
		cache.Put("a", 2) // Update existing key

		value, found := cache.Get("a")
		if !found || value != 2 {
			t.Errorf("Expected updated value 2, got (%v, %v)", value, found)
		}

		if cache.Size() != 1 {
			t.Errorf("Expected size 1 after update, got %d", cache.Size())
		}
	})

	t.Run("Delete", func(t *testing.T) {
		cache := NewLRUCache(2)

		cache.Put("a", 1)
		cache.Put("b", 2)

		deleted := cache.Delete("a")
		if !deleted {
			t.Error("Expected Delete to return true for existing key")
		}

		_, found := cache.Get("a")
		if found {
			t.Error("Expected 'a' to be deleted")
		}

		if cache.Size() != 1 {
			t.Errorf("Expected size 1 after delete, got %d", cache.Size())
		}

		deleted = cache.Delete("nonexistent")
		if deleted {
			t.Error("Expected Delete to return false for non-existent key")
		}
	})

	t.Run("Clear", func(t *testing.T) {
		cache := NewLRUCache(2)

		cache.Put("a", 1)
		cache.Put("b", 2)
		cache.Clear()

		if cache.Size() != 0 {
			t.Errorf("Expected size 0 after clear, got %d", cache.Size())
		}

		_, found := cache.Get("a")
		if found {
			t.Error("Expected cache to be empty after clear")
		}
	})
}

// TestLFUCache tests the LFU cache implementation
func TestLFUCache(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		cache := NewLFUCache(2)
		if cache == nil {
			t.Fatal("NewLFUCache returned nil")
		}

		// Test initial state
		if cache.Size() != 0 {
			t.Errorf("Expected size 0, got %d", cache.Size())
		}
		if cache.Capacity() != 2 {
			t.Errorf("Expected capacity 2, got %d", cache.Capacity())
		}
	})

	t.Run("LFU Eviction", func(t *testing.T) {
		cache := NewLFUCache(2)

		// Fill cache
		cache.Put("a", 1)
		cache.Put("b", 2)

		// Access "a" multiple times to increase its frequency
		cache.Get("a")
		cache.Get("a")
		// Now "a" has frequency 3, "b" has frequency 1

		// Add new item, should evict "b" (least frequently used)
		cache.Put("c", 3)

		// "b" should be evicted
		_, found := cache.Get("b")
		if found {
			t.Error("Expected 'b' to be evicted (least frequently used)")
		}

		// "a" and "c" should still be present
		value, found := cache.Get("a")
		if !found || value != 1 {
			t.Errorf("Expected 'a' to be present with value 1, got (%v, %v)", value, found)
		}

		value, found = cache.Get("c")
		if !found || value != 3 {
			t.Errorf("Expected 'c' to be present with value 3, got (%v, %v)", value, found)
		}
	})

	t.Run("Tie Breaking", func(t *testing.T) {
		cache := NewLFUCache(2)

		cache.Put("a", 1)
		cache.Put("b", 2)
		// Both have frequency 1, should evict the oldest one when tied

		cache.Put("c", 3)

		// "a" should be evicted (inserted first, so oldest among tied frequencies)
		_, found := cache.Get("a")
		if found {
			t.Error("Expected 'a' to be evicted (oldest among tied frequencies)")
		}
	})
}

// TestFIFOCache tests the FIFO cache implementation
func TestFIFOCache(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		cache := NewFIFOCache(2)
		if cache == nil {
			t.Fatal("NewFIFOCache returned nil")
		}

		// Test initial state
		if cache.Size() != 0 {
			t.Errorf("Expected size 0, got %d", cache.Size())
		}
		if cache.Capacity() != 2 {
			t.Errorf("Expected capacity 2, got %d", cache.Capacity())
		}
	})

	t.Run("FIFO Eviction", func(t *testing.T) {
		cache := NewFIFOCache(2)

		// Fill cache
		cache.Put("a", 1)
		cache.Put("b", 2)

		// Access "a" (shouldn't affect eviction order in FIFO)
		cache.Get("a")

		// Add new item, should evict "a" (first in, first out)
		cache.Put("c", 3)

		// "a" should be evicted
		_, found := cache.Get("a")
		if found {
			t.Error("Expected 'a' to be evicted (first in, first out)")
		}

		// "b" and "c" should still be present
		value, found := cache.Get("b")
		if !found || value != 2 {
			t.Errorf("Expected 'b' to be present with value 2, got (%v, %v)", value, found)
		}

		value, found = cache.Get("c")
		if !found || value != 3 {
			t.Errorf("Expected 'c' to be present with value 3, got (%v, %v)", value, found)
		}
	})
}

// TestThreadSafeCache tests the thread-safe wrapper
func TestThreadSafeCache(t *testing.T) {
	t.Run("Concurrent Access", func(t *testing.T) {
		baseCache := NewLRUCache(100)
		cache := NewThreadSafeCache(baseCache)
		if cache == nil {
			t.Fatal("NewThreadSafeCache returned nil")
		}

		const numGoroutines = 10
		const numOperations = 100

		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					cache.Put(key, j)
					cache.Get(key)
					if j%10 == 0 {
						cache.Delete(key)
					}
				}
			}(i)
		}

		wg.Wait()

		// Should not panic and should have some items
		if cache.Size() < 0 {
			t.Error("Cache size should not be negative after concurrent operations")
		}
	})

	t.Run("Interface Compatibility", func(t *testing.T) {
		baseCache := NewLRUCache(2)
		cache := NewThreadSafeCache(baseCache)

		// Test that all methods work through the interface
		cache.Put("a", 1)
		value, found := cache.Get("a")
		if !found || value != 1 {
			t.Errorf("Expected (1, true), got (%v, %v)", value, found)
		}

		if cache.Size() != 1 {
			t.Errorf("Expected size 1, got %d", cache.Size())
		}

		if cache.Capacity() != 2 {
			t.Errorf("Expected capacity 2, got %d", cache.Capacity())
		}

		deleted := cache.Delete("a")
		if !deleted {
			t.Error("Expected Delete to return true")
		}

		cache.Clear()
		if cache.Size() != 0 {
			t.Errorf("Expected size 0 after clear, got %d", cache.Size())
		}
	})
}

// TestCacheFactory tests the factory functions
func TestCacheFactory(t *testing.T) {
	t.Run("NewCache", func(t *testing.T) {
		lruCache := NewCache(LRU, 2)
		if lruCache == nil {
			t.Error("Expected non-nil LRU cache from factory")
		}

		lfuCache := NewCache(LFU, 2)
		if lfuCache == nil {
			t.Error("Expected non-nil LFU cache from factory")
		}

		fifoCache := NewCache(FIFO, 2)
		if fifoCache == nil {
			t.Error("Expected non-nil FIFO cache from factory")
		}

		// Test that they work
		lruCache.Put("test", 1)
		value, found := lruCache.Get("test")
		if !found || value != 1 {
			t.Error("LRU cache from factory not working correctly")
		}
	})

	t.Run("NewThreadSafeCacheWithPolicy", func(t *testing.T) {
		cache := NewThreadSafeCacheWithPolicy(LRU, 2)
		if cache == nil {
			t.Error("Expected non-nil thread-safe cache from factory")
		}

		// Test basic functionality
		cache.Put("test", 1)
		value, found := cache.Get("test")
		if !found || value != 1 {
			t.Error("Thread-safe cache from factory not working correctly")
		}
	})
}

// TestHitRate tests hit rate calculation
func TestHitRate(t *testing.T) {
	cache := NewLRUCache(2)

	// Initial hit rate should be 0 (no operations)
	if cache.HitRate() != 0.0 {
		t.Errorf("Expected initial hit rate 0.0, got %f", cache.HitRate())
	}

	cache.Put("a", 1)
	cache.Put("b", 2)

	// Miss
	cache.Get("c")
	hitRate := cache.HitRate()
	if hitRate != 0.0 {
		t.Errorf("Expected hit rate 0.0 after miss, got %f", hitRate)
	}

	// Hit
	cache.Get("a")
	hitRate = cache.HitRate()
	if hitRate != 0.5 {
		t.Errorf("Expected hit rate 0.5, got %f", hitRate)
	}

	// Another hit
	cache.Get("b")
	hitRate = cache.HitRate()
	expectedHitRate := 2.0 / 3.0 // 2 hits out of 3 total gets
	if hitRate < expectedHitRate-0.01 || hitRate > expectedHitRate+0.01 {
		t.Errorf("Expected hit rate ~%.3f, got %f", expectedHitRate, hitRate)
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("Zero Capacity", func(t *testing.T) {
		cache := NewLRUCache(0)
		if cache == nil {
			t.Skip("Zero capacity cache not supported")
		}

		cache.Put("a", 1)
		_, found := cache.Get("a")
		if found {
			t.Error("Zero capacity cache should not store any items")
		}

		if cache.Size() != 0 {
			t.Error("Zero capacity cache should always have size 0")
		}
	})

	t.Run("Single Capacity", func(t *testing.T) {
		cache := NewLRUCache(1)

		cache.Put("a", 1)
		cache.Put("b", 2) // Should evict "a"

		_, found := cache.Get("a")
		if found {
			t.Error("Expected 'a' to be evicted in single capacity cache")
		}

		value, found := cache.Get("b")
		if !found || value != 2 {
			t.Error("Expected 'b' to be present in single capacity cache")
		}
	})

	t.Run("Nil Values", func(t *testing.T) {
		cache := NewLRUCache(2)

		cache.Put("nil", nil)
		value, found := cache.Get("nil")
		if !found || value != nil {
			t.Errorf("Expected (nil, true), got (%v, %v)", value, found)
		}
	})

	t.Run("Empty Key", func(t *testing.T) {
		cache := NewLRUCache(2)

		cache.Put("", "empty")
		value, found := cache.Get("")
		if !found || value != "empty" {
			t.Errorf("Expected ('empty', true), got (%v, %v)", value, found)
		}
	})
}

// BenchmarkCacheOperations benchmarks cache performance
func BenchmarkCacheOperations(b *testing.B) {
	cache := NewLRUCache(1000)

	// Populate cache
	for i := 0; i < 1000; i++ {
		cache.Put(fmt.Sprintf("key-%d", i), i)
	}

	b.Run("Get", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Get(fmt.Sprintf("key-%d", i%1000))
		}
	})

	b.Run("Put", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Put(fmt.Sprintf("key-%d", i), i)
		}
	})

	b.Run("Delete", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Delete(fmt.Sprintf("key-%d", i))
		}
	})
}

// TestCacheComparison compares different cache policies
func TestCacheComparison(t *testing.T) {
	const capacity = 3

	lruCache := NewLRUCache(capacity)
	lfuCache := NewLFUCache(capacity)
	fifoCache := NewFIFOCache(capacity)

	// Sequence of operations that will show different behaviors
	operations := []struct {
		op    string
		key   string
		value int
	}{
		{"put", "a", 1},
		{"put", "b", 2},
		{"put", "c", 3},
		{"get", "a", 0}, // Access "a"
		{"get", "a", 0}, // Access "a" again (higher frequency)
		{"put", "d", 4}, // This will trigger eviction
	}

	for _, op := range operations {
		if op.op == "put" {
			lruCache.Put(op.key, op.value)
			lfuCache.Put(op.key, op.value)
			fifoCache.Put(op.key, op.value)
		} else if op.op == "get" {
			lruCache.Get(op.key)
			lfuCache.Get(op.key)
			fifoCache.Get(op.key)
		}
	}

	// Check final states
	t.Logf("After operations, cache sizes - LRU: %d, LFU: %d, FIFO: %d",
		lruCache.Size(), lfuCache.Size(), fifoCache.Size())

	// All should have the same size (capacity)
	if lruCache.Size() != capacity || lfuCache.Size() != capacity || fifoCache.Size() != capacity {
		t.Error("All caches should be at capacity after operations")
	}

	// But they might have evicted different items
	// This is mainly for observational purposes in this test
}

// TestConcurrentStress performs stress testing with multiple goroutines
func TestConcurrentStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	cache := NewThreadSafeCache(NewLRUCache(100))
	const numGoroutines = 50
	const numOperations = 1000

	start := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j%100) // Limited key space for conflicts
				cache.Put(key, j)
				cache.Get(key)
				if j%10 == 0 {
					cache.Delete(key)
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("Stress test completed in %v", duration)
	t.Logf("Final cache size: %d, hit rate: %.3f", cache.Size(), cache.HitRate())

	if cache.Size() < 0 || cache.Size() > cache.Capacity() {
		t.Errorf("Invalid cache size after stress test: %d", cache.Size())
	}
}
