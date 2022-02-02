package hw04lrucache

import "sync"

var mu sync.Mutex

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
	mu.Lock()
	defer mu.Unlock()
	listItem, ok := lru.items[key]
	if ok {
		item := listItem.Value.(cacheItem)
		item.value = value
		listItem.Value = item
		lru.queue.MoveToFront(listItem)
	} else {
		item := cacheItem{key, value}

		if lru.capacity == lru.queue.Len() {
			item := lru.queue.Back().Value
			delete(lru.items, item.(cacheItem).key)
			lru.queue.Remove(lru.queue.Back())
		}

		listItem = lru.queue.PushFront(item)
		lru.items[key] = listItem
	}

	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()
	if listItem, ok := lru.items[key]; ok {
		lru.queue.MoveToFront(listItem)
		item := listItem.Value.(cacheItem)
		return item.value, true
	}

	return nil, false
}

func (lru *lruCache) Clear() {
	mu.Lock()
	defer mu.Unlock()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}
