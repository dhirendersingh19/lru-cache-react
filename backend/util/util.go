package util

import (
	"container/list"
	"sync"
	"time"
)

// Node represents a node in the doubly linked list
type Node struct {
	key   string
	value interface{}
	prev  *Node
	next  *Node
}

// LRUcache implements a Least Recently Used cache
type LRUcache struct {
	capacity int
	cache    sync.Map   // Concurrent map for key-value pairs
	list     *list.List // Doubly linked list for tracking order
	head     *Node      // Head node of the list
	tail     *Node      // Tail node of the list
	cleanup  *time.Ticker
}

// NewLRUcache creates a new LRU cache with specified capacity
func NewLRUcache(capacity int) *LRUcache {
	cache := &LRUcache{
		capacity: capacity,
		cache:    sync.Map{},
		list:     list.New(),
	}
	cache.head, cache.tail = &Node{}, &Node{}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head
	cache.cleanup = time.NewTicker(5 * time.Minute) // Set ticker for 5-minute cleanup
	go cache.cleanupLoop()
	return cache
}

// Get retrieves a value from the cache based on the key
func (c *LRUcache) Get(key string) (interface{}, bool) {
	node, ok := c.cache.Load(key)
	if !ok {
		return nil, false
	}
	c.moveToHead(node.(*Node))
	return node.(*Node).value, true
}

// Set stores a key-value pair in the cache
func (c *LRUcache) Set(key string, value interface{}) {
	node, ok := c.cache.Load(key)
	if ok {
		node.(*Node).value = value
		c.moveToHead(node.(*Node))
		return
	}
	newNode := &Node{key: key, value: value}
	c.cache.Store(key, newNode)
	c.list.PushBack(newNode)
	c.removeLeastRecentlyUsed()
}

// removeLeastRecentlyUsed removes the least recently used node from the cache
func (c *LRUcache) removeLeastRecentlyUsed() {
	if c.list.Len() > c.capacity {
		oldest := c.list.Remove(c.tail.prev).(*Node)
		c.cache.Delete(oldest.key)
	}
}

// moveToHead moves a node to the head of the list
func (c *LRUcache) moveToHead(node *Node) {
	c.list.Remove(node)
	c.list.PushFront(node)
}

// cleanupLoop continuously checks for expired entries and removes them
func (c *LRUcache) cleanupLoop() {
	for range c.cleanup.C {
		c.cache.Range(func(key, value interface{}) bool {
			// Implement logic to check if entry is expired based on your needs
			// You can store an expiry timestamp with the value
			// and compare it with the current time
			return true // Continue iterating through the map
		})
		// After iterating through the map, remove expired entries here
	}
}

// func main() {
// 	cache := NewLRUcache(10)
// 	cache.Set("key1", "value1")
// 	cache.Set("key2", "value2")
// 	// ... (use Get and Set methods)

// 	time.Sleep(6 * time.Minute) // Wait for cleanup to run

// 	val, ok := cache.Get("key1")
// 	if !ok {
// 		fmt.Println("key1 expired or not found")
// 	} else {
// 		fmt.Println("key1 value:", val)
// 	}
// }
