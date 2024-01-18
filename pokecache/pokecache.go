package pokecache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type CacheEntry struct {
	value     []byte
	createdAt time.Time
}

type Cache struct {
	entries  map[string]*CacheEntry
	mutex    sync.Mutex
	interval time.Duration
	shutdown chan struct{} // New channel for shutdown signal
}

// function to create a new cache
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]*CacheEntry),
		mutex:    sync.Mutex{},
		interval: interval,
		shutdown: make(chan struct{}), // Initialize the shutdown channel
	}

	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = &CacheEntry{
		createdAt: time.Now(),
		value:     value,
	}
}

var defaultCache *Cache

func init() {
	defaultCache = NewCache(7)
}

func GetDefaultCache() *Cache {
	return defaultCache
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// check if the key exists in the cach
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	// check if the entry is expired
	if time.Since(entry.createdAt) > c.interval {
		delete(c.entries, key)
		return nil, false
	}
	return entry.value, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mutex.Lock()
			for key, entry := range c.entries {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.entries, key)
				}
			}
			c.mutex.Unlock()
		}
	}
}

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
