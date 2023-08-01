package cache

import (
	"errors"
	"time"
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
