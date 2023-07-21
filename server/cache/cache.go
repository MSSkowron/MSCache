package cache

import "time"

// Key is a string that represents a key in the cache.
type Key string

// Value is a struct that represents a value in the cache.
type Value struct {
	Value []byte
	TTL   time.Duration
}

// Cache is an interface that describes the behavior of a cache.
type Cache interface {
	Set(Key, Value) error
	Get(Key) (Value, error)
	Delete(Key) error
	Contains(Key) (bool, error)
}
