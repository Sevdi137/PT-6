package main
import (
    "fmt"
    "sync"
    "time"
)
type Cache struct {
    data map[string]string
    mu   sync.RWMutex
}

func (c *Cache) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

func main() {
    cache := Cache{data: make(map[string]string)}
    
    for i := 0; i < 5; i++ {
        go func() {
            fmt.Println("Read:", cache.Get("test"))
        }()
    }
    
    cache.Set("test", "value")
    time.Sleep(time.Millisecond * 100)
}
