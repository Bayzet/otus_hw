package hw04lrucache

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
	lruItem, ok := lru.items[key]
	if ok {
		lruItem.Value = value
		lru.queue.MoveToFront(lruItem)
	} else {
		if lru.capacity == lru.queue.Len() {
			delete(lru.items, lru.queue.Back().key)
			lru.queue.Remove(lru.queue.Back())
		}
		lruItem := lru.queue.PushFront(value)
		lruItem.key = key
		lru.items[key] = lruItem
	}

	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {	
	if lruItem, ok := lru.items[key]; ok {
		lru.queue.MoveToFront(lruItem)

		return lruItem.Value, true
	} else {
		return nil, false
	}
}

func (lru *lruCache) Clear() {
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}
