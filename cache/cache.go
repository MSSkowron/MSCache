package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Cache is a struct that represents key-value cache.
type Cache struct {
	mu   sync.RWMutex
	data map[string][]byte //key cannot be []byte while it is not comparable type which is required in the map data type.
}

// New returns a pointer to an empty Cache struct.
func New() *Cache {
	return &Cache{
		mu:   sync.RWMutex{},
		data: make(map[string][]byte),
	}
}

// Set adds a key element to the Cache with value specified by the value parameter.
// The element will be deleted after time spicifed by ttl.
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

// Get returns an value of element with key specified by the key parameter.
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

// Delete removes an element with the key specified by the key parameter.
func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, string(key))

	return nil
}

// Contains returns true if exists a value with the key specified by the parameter, otherwise returns false.
func (c *Cache) Contains(key []byte) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[string(key)]
	return ok
}
