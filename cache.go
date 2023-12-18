package plrucache

import (
	"sync"
	"time"
)

// LRUCache - has a buffer for maxSize elements and allow to store and get values from it.
// Whenever new element should be stored it remove least recent used element if there is no space left in the buffer.
type LRUCache[T any] struct {
	maxSize int
	rwLock  sync.RWMutex
	ttl     time.Duration
	stat    Stat
	tsQueue *tsQ
	lru     *staticQ
	store   map[string]item[T]
}

// New initialize LRUCache object and returns pointer to it.
func New[T any](maxSize int, ttl time.Duration) *LRUCache[T] {
	return &LRUCache[T]{
		maxSize: maxSize,
		lru:     newQueue(maxSize),
		store:   make(map[string]item[T], maxSize),
		stat:    Stat{},
		rwLock:  sync.RWMutex{},
		ttl:     ttl,
		tsQueue: newTSQ(maxSize),
	}
}

// Set stores value under associated key.
func (c *LRUCache[T]) Set(key string, val T) {
	c.rwLock.Lock()
    now := time.Now()
	expTs := now.Add(c.ttl)
	newItem := item[T]{val: val, expiration: expTs}
	if len(c.store)+1 > c.maxSize {
		if c.tsQueue.IsAnyExpired(time.Now()) {
			item, _ := c.tsQueue.Pop()
			val := c.store[item.val]
			c.lru.Delete(val.lruIdx)
			delete(c.store, item.val)
		} else {
			item, _ := c.lru.Pop()
			val := c.store[item.val]
			c.tsQueue.Delete(val.tsqIdx)
			delete(c.store, item.val)
		}
	}
    newItem.tsqIdx = c.tsQueue.Push(key, expTs)
	newItem.lruIdx = c.lru.Push(key, now)
	c.store[key] = newItem
	c.rwLock.Unlock()
}

// Get return value assosiated with key or nil if key not exists.
// Boolean flag indicates whether value was found or not.
func (c *LRUCache[T]) Get(key string) (any, bool) {
	c.rwLock.Lock()
	val, ok := c.store[key]
	if !ok {
		c.stat.Misses++
        c.rwLock.Unlock()
		return nil, false
	}
	c.stat.Hits++
    // rearange lru queue
    c.lru.Delete(val.lruIdx)
    val.lruIdx = c.lru.Push(key, time.Now())
    c.store[key] = val
	c.rwLock.Unlock()
	return val.val, true
}

// Len returns count of items int cache.
func (c *LRUCache[T]) Len() int {
	c.rwLock.RLock()
    v := len(c.store)
	c.rwLock.RUnlock()
    return v
}

// Delete removes item by key.
func (c *LRUCache[T]) Delete(key string) {
	c.rwLock.Lock()
	val, ok := c.store[key]
	if ok {
		c.lru.Delete(val.lruIdx)
		c.tsQueue.Delete(val.tsqIdx)
		delete(c.store, key)
	}
	c.rwLock.Unlock()
}

// Reset drop all items from cache.
func (c *LRUCache[T]) Reset() {
	c.rwLock.Lock()
	c.lru = newQueue(c.maxSize)
	c.store = make(map[string]item[T], c.maxSize)
	c.tsQueue = newTSQ(c.maxSize)
	c.stat = Stat{}
	c.rwLock.Unlock()
}

// Stat returns statistics of usage.
func (c *LRUCache[T]) Stat() Stat {
	return c.stat
}
