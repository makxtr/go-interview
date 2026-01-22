package cache

import "go-interview/a/ch28/list"

//
// LFU Cache Implementation
//

type lfuPayload struct {
	key   string
	value interface{}
	freq  int
}

type LFUCache struct {
	capacity   int
	minFreq    int
	cache      map[string]*list.DoublyNode[lfuPayload]
	freqGroups map[int]*list.DoublyLinkedList[lfuPayload]
	hits       uint64
	misses     uint64
}

func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
		return nil
	}
	return &LFUCache{
		capacity:   capacity,
		cache:      make(map[string]*list.DoublyNode[lfuPayload]),
		freqGroups: make(map[int]*list.DoublyLinkedList[lfuPayload]),
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
	return node.Value.value, true
}

func (c *LFUCache) Put(key string, value interface{}) {
	if c.capacity <= 0 {
		return
	}
	if node, ok := c.cache[key]; ok {
		node.Value.value = value
		c.updateNodeFreq(node)
		return
	}
	if len(c.cache) >= c.capacity {
		oldestFreqList := c.freqGroups[c.minFreq] // Renamed local variable
		if oldestFreqList != nil {
			nodeToEvict := oldestFreqList.Back()
			if nodeToEvict != nil {
				oldestFreqList.Remove(nodeToEvict)
				delete(c.cache, nodeToEvict.Value.key)
			}
		}
	}
	c.minFreq = 1
	payload := lfuPayload{key: key, value: value, freq: 1}
	newList, exists := c.freqGroups[1]
	if !exists {
		newList = list.NewDoubly[lfuPayload]()
		c.freqGroups[1] = newList
	}
	node := newList.PushFront(payload)
	c.cache[key] = node
}

func (c *LFUCache) updateNodeFreq(node *list.DoublyNode[lfuPayload]) {
	oldFreq := node.Value.freq
	oldFreqList := c.freqGroups[oldFreq] // Renamed local variable
	oldFreqList.Remove(node)

	if oldFreq == c.minFreq && oldFreqList.Len == 0 {
		delete(c.freqGroups, oldFreq)
		c.minFreq++
	}

	newFreq := oldFreq + 1
	node.Value.freq = newFreq

	newList, exists := c.freqGroups[newFreq]
	if !exists {
		newList = list.NewDoubly[lfuPayload]()
		c.freqGroups[newFreq] = newList
	}
	newList.PushFrontNode(node)
}

func (c *LFUCache) Delete(key string) bool {
	node, ok := c.cache[key]
	if !ok {
		return false
	}
	delete(c.cache, key)
	freqList := c.freqGroups[node.Value.freq] // Renamed local variable
	freqList.Remove(node)
	if freqList.Len == 0 {
		delete(c.freqGroups, node.Value.freq)
	}
	return true
}

func (c *LFUCache) Clear() {
	c.cache = make(map[string]*list.DoublyNode[lfuPayload])
	c.freqGroups = make(map[int]*list.DoublyLinkedList[lfuPayload])
	c.minFreq = 0
	c.hits = 0
	c.misses = 0
}

func (c *LFUCache) Size() int { return len(c.cache) }

func (c *LFUCache) Capacity() int { return c.capacity }

func (c *LFUCache) HitRate() float64 {
	total := c.hits + c.misses
	if total == 0 {
		return 0.0
	}
	return float64(c.hits) / float64(total)
}
