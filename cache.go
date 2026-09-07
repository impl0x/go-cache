package gocache

import "sync"

type CacheConfig struct {
	InitialMapCapacity int
}

var DefaultCacheConfig = CacheConfig{} // Default config where the InitialMapCapacity is set to 0 same as the default map initiation in go

// Cache implementation with generic map and read write mutex
type Cache[K k, V v] struct {
	m  map[K]V
	mu sync.RWMutex
}

func newCache[K k, V v](config CacheConfig) Cache[K, V] {
	return Cache[K, V]{
		m: make(map[K]V, config.InitialMapCapacity),
	}
}

// Instantiates a new cache using the generic types provided to this constructor function and default [DefaultCacheConfig]
func NewCache[K comparable, V any]() *Cache[K, V] {
	c := newCache[K, V](DefaultCacheConfig)
	return &c
}

// Instantiates a new cache using the config provided and generics
func NewCacheWithConfig[K comparable, V any](config CacheConfig) *Cache[K, V] {
	c := newCache[K, V](config)
	return &c
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

// Loops over all the items in the list and passes the key and value to the function provided
//
// the function must return a uint8, which tells when to quit the loop
//   - return 0 to exit the loop
//   - return >= 1 to continue the loop
//
// the read lock unlocks itself only after this function call has ended, do not run any other methods on this instance of [Cache] inside the provided function fn
func (c *Cache[K, V]) LoopFunc(fn func(key K, value V) uint8) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for k, v := range c.m {
		if r := fn(k, v); r == 0 {
			return
		}
	}
}