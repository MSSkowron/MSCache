package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/MSSkowron/MSCache/server/logger"
)

var (
	// ErrKeyIsEmpty is returned when the key is empty.
	ErrKeyIsEmpty = errors.New("key cannot be empty")
	// ErrValueIsNil is returned when the value is nil.
	ErrValueIsNil = errors.New("value cannot be nil")
	// ErrValueIsEmpty is returned when the value is empty.
	ErrValueIsEmpty = errors.New("value cannot be empty")
	// ErrInvalidTTL is returned when the TTL is less than or equal to 0.
	ErrInvalidTTL = errors.New("ttl must be greater than 0")
	// ErrKeyNotFound is returned when the key is not found in the cache.
	ErrKeyNotFound = errors.New("key not found")
)

// InMemoryCache is a struct that represents a key-value In-Memory Cache.
type InMemoryCache struct {
	data map[Key]Value // data stores key-value pairs in the cache.
	mu   sync.RWMutex  // mu is a read-write mutex used to synchronize concurrent access to the cache.
}

// New creates a new InMemoryCache.
func New() *InMemoryCache {
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

	logger.CustomLogger.Info.Printf("SET %s to %+v\n", key, value)

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

	logger.CustomLogger.Info.Printf("GET %s", key)

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

	logger.CustomLogger.Info.Printf("DELETE %s", key)

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

	logger.CustomLogger.Info.Printf("CONTAINS %s", key)

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
