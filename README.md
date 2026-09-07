# Go Cache
## Introduction
This is a caching package made using generics, maps and mutexes. Simple implementation of caching and does not contain memory management, limits, etc yet.  
**The cache is safe to use concurrently and it manages its own mutex**

## Quick start
```go
package main

import gocache "github.com/impl0x/go-cache"

type item struct {
	message string
}
func main() {
	// Instantiation of a new cache using generics
	cache := gocache.NewCache[string, item]()

	// Adding some items to the cache
	cache.Add("item1", item{"this is item 1"})
	cache.Add("item2", item{"this is item 2"})
	
	// Getting a item from the cache
	item1, ok := cache.Get("item1")
	if ok {
		println(item1.message)
	}
}
```
## Documentation
There are 2 types of caches as of now, the normal Cache and a TTLCache.
### Cache
```go
package main

import gocache "github.com/impl0x/go-cache"

type item struct {
	message string
}

func main() {
	// Instantiation of a new cache using generics
	cache := gocache.NewCache[string, item]()

	// Adding a item
	cache.Add("item", item{"this is a item"})
	
	// Getting a item
	i, ok := cache.Get("item")

	// Updating a item
	cache.Update("item",item{"this is a updated item"})

	// Deleting a item
	cache.Delete("item")
}
```
it is also possible to provide extra configuration using the constructor function `NewCacheWithConfig` where you need to provide a type of `CacheConfig` containing configuration values. 
### TTL Cache - time lived cache items
This is a wrapper on the original Cache struct and it takes a expiry duration for each item and automatically cleans it from the cache when the expiry hits.  
``` go
package main

import (
	"time"

	gocache "github.com/impl0x/go-cache"
)

type item struct {
	message string
}

func main() {
	const cleanInterval = 10 * time.Minute
	const expiryDuration = 5 * time.Minute
	// Instantiation of a new cache using generics
	cache := gocache.NewTTLCache[string, item](cleanInterval)

	// Adding a item with its expiry time
	cache.Add("item", item{"this is a item"}, time.Now().Add(expiryDuration))

	// Getting a item
	i, ok, expiresAt := cache.Get("item")

	// Updating a item
	cache.Update("item", item{"this is a updated item"}, time.Time{}) // providing empty time to leave it unchanged

	// Deleting a item
	cache.Delete("item")
}
```
The usage is the same, with one extra added parameter of expiresAt.  
it is also possible to provide extra configuration using the constructor function `NewTTLCacheWithConfig` where you need to provide a type of `TTLCacheConfig` containing configuration values. 

### Extra functions 
**`LoopFunc`**  
This takes in a function which is provided the key and value, (and expiresAt in the case of TTLCache) on each iteration of the underlying map. So this can be used to track or find values via the Value instead of the key. Although not recommended to use if the cache is big as it takes O(N) time to loop over all the items in a map. The loop can be quit early by returning 0 from the function, which stops the iteration and is a signal that you have found what you are looking for. The loop manages its own read lock so it is not safe to call any other method from inside the function as it will cause deadlocks. Returning anything other than 0 signals the loop to continue and keep calling the function. 