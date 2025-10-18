package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type StubClock struct {
	now time.Time
}

func NewStubClock() *StubClock {
	return &StubClock{
		now: time.Now(),
	}
}

func (s *StubClock) Now() time.Time {
	return s.now
}

func (s *StubClock) advanceBy(d time.Duration) {
	s.now = s.now.Add(d)
}

func TestCache(t *testing.T) {
	t.Run("add and get value", func(t *testing.T) {
		t.Skip()
		cache := NewCache[string, int]()
		key := "hello"
		val := 12

		cache.Set(key, val, 5*time.Second)

		got, ok := cache.Get(key)
		assert.True(t, ok)
		assert.Equal(t, got, val)
	})

	t.Run("ttl works", func(t *testing.T) {
		clock := NewStubClock()
		cache := &Cache[string, int]{map[string]CacheValue[int]{}, clock}
		key := "hello"
		val := 12

		cache.Set(key, val, 5*time.Minute)

		got, found := cache.Get(key)
		assert.True(t, found)
		assert.Equal(t, got, val)

		clock.advanceBy(4 * time.Minute)

		got, found = cache.Get(key)
		assert.True(t, found)
		assert.Equal(t, val, got)

		clock.advanceBy(61 * time.Second)

		_, found = cache.Get(key)
		assert.False(t, found)
		assert.Equal(t, len(cache.store), 0)
	})

}
