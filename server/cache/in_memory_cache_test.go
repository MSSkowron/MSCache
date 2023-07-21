package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := New()

	assert.NotNil(t, c)
	assert.NotNil(t, c.data)

	assert.Equal(t, 0, len(c.data))
}

func TestSet(t *testing.T) {
	data := []struct {
		name          string
		key           Key
		value         Value
		expectedError error
	}{
		{
			name: "valid",
			key:  Key("key"),
			value: Value{
				Value: []byte("value"),
				TTL:   1 * time.Second,
			},
			expectedError: nil,
		},
		{
			name: "empty key",
			key:  Key(""),
			value: Value{
				Value: []byte("value"),
				TTL:   1 * time.Second,
			},
			expectedError: ErrKeyIsEmpty,
		},
		{
			name: "nil value",
			key:  Key("key"),
			value: Value{
				Value: nil,
				TTL:   1 * time.Second,
			},
			expectedError: ErrValueIsNil,
		},
		{
			name: "empty value",
			key:  Key("key"),
			value: Value{
				Value: []byte(""),
				TTL:   1 * time.Second,
			},
			expectedError: ErrValueIsEmpty,
		},
		{
			name: "zero TTL",
			key:  Key("key"),
			value: Value{
				Value: []byte("value"),
				TTL:   0,
			},
			expectedError: ErrInvalidTTL,
		},
		{
			name: "negative TTL",
			key:  Key("key"),
			value: Value{
				Value: []byte("value"),
				TTL:   -10 * time.Second,
			},
			expectedError: ErrInvalidTTL,
		},
	}

	c := New()

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := c.Set(d.key, d.value)
			assert.Equal(t, d.expectedError, err)
		})
	}

	time.Sleep(2 * time.Second)

	assert.Equal(t, 0, len(c.data))
}

func TestGet(t *testing.T) {
	data := []struct {
		name          string
		key           Key
		value         Value
		expectedError error
	}{
		{
			name:          "empty key",
			key:           Key(""),
			value:         Value{},
			expectedError: ErrKeyIsEmpty,
		},
		{
			name:          "key not found",
			key:           Key("key"),
			value:         Value{},
			expectedError: ErrKeyNotFound,
		},
	}

	c := New()

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			value, err := c.Get(d.key)
			assert.Equal(t, d.expectedError, err)
			assert.Equal(t, d.value, value)
		})
	}

	key := Key("key")
	value := Value{
		Value: []byte("value"),
		TTL:   1 * time.Second,
	}

	err := c.Set(key, value)
	assert.Nil(t, err)

	v, err := c.Get(key)
	assert.Nil(t, err)

	assert.Equal(t, value, v)
}

func TestDelete(t *testing.T) {
	data := []struct {
		name          string
		key           Key
		value         Value
		expectedError error
	}{
		{
			name:          "empty key",
			key:           Key(""),
			value:         Value{},
			expectedError: ErrKeyIsEmpty,
		},
	}

	c := New()

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := c.Delete(d.key)
			assert.Equal(t, d.expectedError, err)
		})
	}

	key := Key("key")
	value := Value{
		Value: []byte("value"),
		TTL:   5 * time.Second,
	}

	err := c.Set(key, value)
	assert.Nil(t, err)

	err = c.Delete(key)
	assert.Nil(t, err)

	_, err = c.Get(key)
	assert.Equal(t, ErrKeyNotFound, err)
}

func TestContaints(t *testing.T) {
	data := []struct {
		name          string
		key           Key
		value         Value
		expectedError error
	}{
		{
			name:          "empty key",
			key:           Key(""),
			value:         Value{},
			expectedError: ErrKeyIsEmpty,
		},
	}

	c := New()

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_, err := c.Contains(d.key)
			assert.Equal(t, d.expectedError, err)
		})
	}

	key := Key("key")
	value := Value{
		Value: []byte("value"),
		TTL:   5 * time.Second,
	}

	err := c.Set(key, value)
	assert.Nil(t, err)

	ok, err := c.Contains(key)
	assert.Nil(t, err)

	assert.True(t, ok)

	err = c.Delete(key)
	assert.Nil(t, err)

	ok, err = c.Contains(key)
	assert.Nil(t, err)

	assert.False(t, ok)
}
