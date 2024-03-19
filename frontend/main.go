package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type entry struct {
	key   string
	value interface{}
}

type LRUCache struct {
	capacity  int
	eviction  *list.List
	cache     map[string]*list.Element
	accessMap sync.Map
	ticker    *time.Ticker
}

func NewLRUCache(capacity int) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		eviction: list.New(),
		cache:    make(map[string]*list.Element),
		ticker:   time.NewTicker(5 * time.Second),
	}

	// Goroutine for automatic eviction
	go func() {
		for range cache.ticker.C {
			cache.evictLRU()
		}
	}()

	return cache
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	if elem, ok := c.cache[key]; ok {
		c.eviction.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, value interface{}) {
	c.accessMap.Store(key, time.Now())
	if elem, ok := c.cache[key]; ok {
		c.eviction.MoveToFront(elem)
		elem.Value.(*entry).value = value
	} else {
		if len(c.cache) >= c.capacity {
			c.evictLRU()
		}
		elem := c.eviction.PushFront(&entry{key: key, value: value})
		c.cache[key] = elem
	}
}

func (c *LRUCache) evictLRU() {
	if elem := c.eviction.Back(); elem != nil {
		c.eviction.Remove(elem)
		delete(c.cache, elem.Value.(*entry).key)
	}
}

func main() {
	cache := NewLRUCache(2)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// After 5 minutes, "key1" and its value will be automatically evicted
	time.Sleep(3 * time.Second)

	value, exists := cache.Get("key1")
	fmt.Println("Value exists:", exists) // Should print false
	fmt.Println("Value:", value)         // Should print nil

	value, exists = cache.Get("key2")
	fmt.Println("Value exists:", exists) // Should print true
	fmt.Println("Value:", value)         // Should print "value2"

	time.Sleep(1 * time.Second)

	value, exists = cache.Get("key1")
	fmt.Println("Value exists:", exists) // Should print false
	fmt.Println("Value:", value)         // Should print nil

	value, exists = cache.Get("key2")
	fmt.Println("Value exists:", exists) // Should print true
	fmt.Println("Value:", value)         // Should print "value2"
}
