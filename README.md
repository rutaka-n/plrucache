[![Go Reference](https://pkg.go.dev/badge/github.com/rutaka-n/plrucache.svg)](https://pkg.go.dev/github.com/rutaka-n/plrucache)

# pLRUCache
(pseudo) LRU cache is a library that implements in-memory (p)LRU cache with focus on efficient memory utilization.
It uses static queues to track expiration time and least recently used items in preallocated memory, so it has minimal
memory footprint. This approach helps to avoid high memory consumtion on the peak load.
It is pseudo LRU, since it relies on `time.Time` to identify least recently used items, so in theory it might drop
not actually least recently used item, but almost least recently used one.
It uses `sync.Mutex` to deal with concurrent read/write operations.

## Install
```sh
go get github.com/rutaka-n/plrucache
```
## Usage
```sh
package main

import (
	"fmt"
	lru "github.com/rutaka-n/plrucache"
	"time"
)

type item struct {
	id  int64
	val string
}

func main() {
	cacheSize := 128
	expirationTime := 300 * time.Second
	cache := lru.New[string, item](cacheSize, expirationTime)

	key := "k1"
	value := item{1, "hello, world"}
	cache.Set(key, value)

	res, ok := cache.Get(key)
    if !ok {
        panic("item is not in cache")
    }
	fmt.Printf("%+v", res)
}
```
