package main
import (
	"fmt"
	"sync"
	"time"
)
type CacheEntry struct {
	value      interface{}
	exp time.Time
}

type CacheV1 struct {
	data map[string]CacheEntry
	mu   sync.RWMutex
	ttl  time.Duration
}

func NewCacheV1(ttl time.Duration) *CacheV1 {
	c := &CacheV1{
		data: make(map[string]CacheEntry),
		ttl:  ttl,
	}
	go c.cleanup()
	return c
}

func (c *CacheV1) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheEntry{
		value:      value,
		exp: time.Now().Add(c.ttl),
	}
}

func (c *CacheV1) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.exp) {
		return nil, false
	}
	return entry.value, true
}

func (c *CacheV1) cleanup() {
	for {
		time.Sleep(c.ttl / 2)
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.exp) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}
type CacheV2 struct {
	data map[string]CacheEntry
	mu   sync.RWMutex
	ttl  time.Duration
}

func NewCacheV2(ttl time.Duration) *CacheV2 {
	c := &CacheV2{
		data: make(map[string]CacheEntry),
		ttl:  ttl,
	}
	go c.cleanup()
	return c
}

func (c *CacheV2) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheEntry{
		value:      value,
		exp: time.Now().Add(c.ttl),
	}
}

func (c *CacheV2) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.exp) {
		return nil, false
	}
	// Обновляем TTL при чтении
	entry.exp = time.Now().Add(c.ttl)
	c.data[key] = entry
	return entry.value, true
}

func (c *CacheV2) cleanup() {
	for {
		time.Sleep(c.ttl / 2)
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.data {
			if now.After(entry.exp) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

func main() {
	fmt.Println("TTL не обновляется при чтении:")
	cache1 := NewCacheV1(2 * time.Second)
	cache1.Set("key1", "value1")
	
	time.Sleep(1 * time.Second)
	if val, ok := cache1.Get("key1"); ok {
		fmt.Println("После 1s:", val)
	}
	
	time.Sleep(2 * time.Second)
	if val, ok := cache1.Get("key1"); ok {
		fmt.Println("После 3s:", val)
	} else {
		fmt.Println("После 3s: истек")
	}

	fmt.Println("TTL обновляется при чтении:")
	cache2 := NewCacheV2(2 * time.Second)
	cache2.Set("key2", "value2")
	
	time.Sleep(1 * time.Second)
	if val, ok := cache2.Get("key2"); ok {
		fmt.Println("После 1s:", val)
	}
	
	time.Sleep(2 * time.Second)
	if val, ok := cache2.Get("key2"); ok {
		fmt.Println("После 3s:", val, "(TTL перезаписан)")
	} else {
		fmt.Println("После 3s: истек")
	}
}
