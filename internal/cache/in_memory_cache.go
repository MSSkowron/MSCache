package cache

import (
	"sync"
	"time"
)

// InMemoryCache is a struct that represents a key-value In-Memory Cache.
type InMemoryCache struct {
	data map[Key]Value // data stores key-value pairs in the cache.
	mu   sync.RWMutex  // mu is a read-write mutex used to synchronize concurrent access to the cache.
}

// NewInMemoryCache creates a new InMemoryCache.
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		mu:   sync.RWMutex{},
		data: make(map[Key]Value),
	}
}

// Set adds a key-value pair to the InMemoryCache with a specified time-to-live (TTL).
// The element will be deleted after the TTL has passed.
func (c *InMemoryCache) Set(key Key, value Value) error {
	if err := c.validateKey(key); err != nil {
		return err
	}

	if err := c.validateValue(value); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value

	time.AfterFunc(value.TTL, func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		delete(c.data, key)
	})

	return nil
}

// Get returns the value of the element with the specified key.
func (c *InMemoryCache) Get(key Key) (Value, error) {
	var value Value

	if err := c.validateKey(key); err != nil {
		return value, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if !ok {
		return value, ErrKeyNotFound
	}

	return value, nil
}

// Delete removes the element with the specified key from the cache.
func (c *InMemoryCache) Delete(key Key) error {
	if err := c.validateKey(key); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)

	return nil
}

// Contains checks if a value with the specified key exists in the cache.
func (c *InMemoryCache) Contains(key Key) (bool, error) {
	if err := c.validateKey(key); err != nil {
		return false, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[key]
	return ok, nil
}

func (c *InMemoryCache) validateKey(key Key) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	return nil
}

func (c *InMemoryCache) validateValue(value Value) error {
	if value.Value == nil {
		return ErrValueIsNil
	}

	if len(value.Value) == 0 {
		return ErrValueIsEmpty
	}

	if value.TTL <= 0 {
		return ErrInvalidTTL
	}

	return nil
}
