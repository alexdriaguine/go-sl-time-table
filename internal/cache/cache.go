package cache

import (
	"sync"
	"time"
)

// Clock interface for mocking out time.Now()
// in tests for ttl/expires
type Clock interface {
	Now() time.Time
}
type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

type Cacher[TKey comparable, TValue any] interface {
	Get(key TKey) (val TValue, found bool)
	Set(key TKey, val TValue, ttl time.Duration)
}

type CacheValue[TValue any] struct {
	value   TValue
	expires time.Time
}

type InMemoryCache[TKey comparable, TValue any] struct {
	store map[TKey]CacheValue[TValue]
	clock Clock
	mu    sync.RWMutex
}

func NewCache[TKey comparable, TValue any]() *InMemoryCache[TKey, TValue] {
	// maps are actually reference types, underlying points to a hash table
	// so no need to have pointers to them
	store := map[TKey]CacheValue[TValue]{}

	// return copies lock value: github.com/alexdriaguine/go-sl-time-table/internal/cache.InMemoryCache[TKey, TValue] contains sync.RWMutexcopylocksdefault
	// because of locks and that they must never be copies
	// the cache must always be a pointer
	return &InMemoryCache[TKey, TValue]{
		store,
		SystemClock{},
		sync.RWMutex{},
	}
}

// Get passes lock by value: github.com/alexdriaguine/go-sl-time-table/internal/cache.InMemoryCache[TKey, TValue] contains sync.RWMutexcopylocksdefault
// the cache MUST use pointer receivers, same reason to why we must
// use the InMemoryCache as a pointer, because we can never copy locks, the lock.
func (c *InMemoryCache[TKey, TValue]) Get(key TKey) (TValue, bool) {
	// can't defer here because in case of expired ttl, need to
	// upgrade to a read lock and delete the entry
	c.mu.RLock()
	val, found := c.store[key]
	if !found {
		c.mu.RUnlock()

		// just creates an empty zero value of type TValue
		var empty TValue
		return empty, false
	}

	if c.clock.Now().After(val.expires) {
		c.mu.RUnlock()
		c.mu.Lock()

		delete(c.store, key)

		c.mu.Unlock()

		var empty TValue
		return empty, false
	}

	c.mu.RUnlock()
	return val.value, true
}

func (c *InMemoryCache[TKey, TValue]) Set(key TKey, value TValue, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = CacheValue[TValue]{
		value:   value,
		expires: time.Now().Add(ttl),
	}
}
