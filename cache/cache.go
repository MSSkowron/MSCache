package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		mu:   sync.RWMutex{},
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[string(key)] = value
	log.Printf("[Cache] SET %s to %s\n", string(key), string(value))

	// It is better to create a go routine that clears up expired keys every some period of time
	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			delete(c.data, string(key))
		}()
	}

	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keyStr := string(key)

	val, ok := c.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyStr)
	}

	log.Printf("[Cache] GET %s = %s\n", string(key), string(val))

	return val, nil
}

func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, string(key))

	return nil
}

func (c *Cache) Contains(key []byte) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}
