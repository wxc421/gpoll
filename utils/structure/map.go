package structure

import (
	"errors"
	"sync"
)

type KeyType interface {
	int | int16 | int32 | int64 | string
}

// ConcurrentMap represents safe concurrent map container
type ConcurrentMap[Key KeyType, Val any] struct {
	mu    sync.RWMutex
	items map[Key]Val
}

// Get return item and exist
func (container *ConcurrentMap[Key, T]) Get(key Key) (item T, exist bool) {
	container.mu.RLock()
	defer container.mu.RUnlock()
	item, exist = container.items[key]
	return
}

// Set sets map item
func (container *ConcurrentMap[Key, T]) Set(key Key, value T) {
	container.mu.Lock()
	container.items[key] = value
	container.mu.Unlock()
}

// Find returns item
// if key is not exist, returns not found error
func (container *ConcurrentMap[Key, T]) Find(key Key) (item T, err error) {
	var exist bool
	item, exist = container.Get(key)
	if !exist {
		err = errors.New("not found")
	}
	return
}

// Del remove keys from container
func (container *ConcurrentMap[Key, T]) Del(key Key) {
	container.mu.Lock()
	delete(container.items, key)
	container.mu.Unlock()
}

func (container *ConcurrentMap[Key, T]) Iterator(fn func(key Key, val T)) {
	container.mu.RLock()
	for key, val := range container.items {
		fn(key, val)
	}
	container.mu.RUnlock()
}

func (container *ConcurrentMap[Key, T]) Len() int {
	return len(container.items)
}
func NewConcurrentMap[Key KeyType, T any]() *ConcurrentMap[Key, T] {
	return &ConcurrentMap[Key, T]{items: make(map[Key]T)}
}
