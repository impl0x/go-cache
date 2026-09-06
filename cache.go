package gocache

import "sync"

// Cache implementation with generic map and read write mutex
type Cache[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

// Instantiates a new cache using the generic types provided to this constructor function
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{}
}

// Adds a key value pair to the cache
func (c *Cache[K, V]) Add(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

// Gets a value with the key provided, returning bool to convey wether the key exists
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}

// Updates a value with the key provided. Does nothing if key doesn't exist
func (c *Cache[K, V]) Update(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

// Deletes a value with the key provided. Does nothing if key doesn't exist
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
}
