package gocache

// Cache with expiration time of items, automatically clears items after it has expired
type TTLCache[K comparable, V any] struct{
	cache Cache[K,V]
}