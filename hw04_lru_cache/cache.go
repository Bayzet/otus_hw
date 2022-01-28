package hw04lrucache

import "fmt"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	fmt.Println("=== Set START")
	fmt.Println("=== === key, value ", key, value)
	fmt.Println("=== === lru.items", lru.items)
	fmt.Println("=== === lru.queue", debug(lru.queue))
	fmt.Println("=== === lru.queue.front", lru.queue.Front())
	fmt.Println("=== === lru.queue.back", lru.queue.Back())
	
	lruItem, ok := lru.items[key]
	if ok {
		lruItem.Value = value
		lru.queue.MoveToFront(lruItem)
	} else {
		fmt.Println("=== === lru",lru)
		if lru.capacity == lru.queue.Len() {
			lru.queue.Remove(lru.queue.Back())
			delete(lru.items, key)
			fmt.Println("=== === delete key lru.items",key, lru.items)
		}
		lruItem := lru.queue.PushFront(value)
		lru.items[key] = lruItem
	}

	fmt.Println("=== === lru.items", lru.items)
	fmt.Println("=== === lru.queue", debug(lru.queue))
	fmt.Println("=== === lru.queue.front", lru.queue.Front())
	fmt.Println("=== === lru.queue.back", lru.queue.Back())
	fmt.Println("Set END")

	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	fmt.Println("=== Get START")
	fmt.Println("=== === key ", key)
	fmt.Println("=== === lru.items", lru.items)
	fmt.Println("=== === lru.queue", debug(lru.queue))
	fmt.Println("=== === lru.queue.front", lru.queue.Front())
	fmt.Println("=== === lru.queue.back", lru.queue.Back())
	
	if lruItem, ok := lru.items[key]; ok {
		fmt.Println("=== === lruItem", lruItem)
		lru.queue.MoveToFront(lruItem)

		fmt.Println("=== === lru.items", lru.items)
		fmt.Println("=== === lru.queue", debug(lru.queue))
		fmt.Println("=== === lru.queue.front", lru.queue.Front())
		fmt.Println("=== === lru.queue.back", lru.queue.Back())
		fmt.Println("=== Get END")

		return lruItem.Value, true
	} else {
		return nil, false
	}
}

func (lru *lruCache) Clear() {
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}
