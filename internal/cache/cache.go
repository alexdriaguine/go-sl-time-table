package cache

import (
	"time"
)

type Clock interface {
	Now() time.Time
}

type CacheValue[TValue any] struct {
	value   TValue
	expires time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

type Cache[TKey comparable, TValue any] struct {
	store map[TKey]CacheValue[TValue]
	clock Clock
}

func NewCache[TKey comparable, TValue any]() *Cache[TKey, TValue] {
	store := map[TKey]CacheValue[TValue]{}
	cache := &Cache[TKey, TValue]{
		store,
		SystemClock{},
	}
	return cache
}

func (c *Cache[TKey, TValue]) Get(key TKey) (TValue, bool) {
	val, found := c.store[key]
	if !found {
		var zero TValue
		return zero, true
	}

	if c.clock.Now().After(val.expires) {
		delete(c.store, key)
		var zero TValue
		return zero, false
	}

	return val.value, true
}

func (c *Cache[TKey, TValue]) Set(key TKey, value TValue, ttl time.Duration) {
	c.store[key] = CacheValue[TValue]{
		value:   value,
		expires: time.Now().Add(ttl),
	}
}
