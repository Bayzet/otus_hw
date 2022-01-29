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
	listItem, ok := lru.items[key]
	if ok {
		item := listItem.Value.(cacheItem)
		item.value = value
		listItem.Value = item
		lru.queue.MoveToFront(listItem)
	} else {
		item := cacheItem{key, value}
		if lru.capacity == lru.queue.Len() {
			lastItem := lru.queue.Back().Value
			delete(lru.items, lastItem.(cacheItem).key)
			lru.queue.Remove(lru.queue.Back())
		}
		listItem := lru.queue.PushFront(item)
		lru.items[key] = listItem
	}

	return ok
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	if listItem, ok := lru.items[key]; ok {
		lru.queue.MoveToFront(listItem)
		item := listItem.Value.(cacheItem)

		return item.value, true
	} else {
		return nil, false
	}
}

func (lru *lruCache) Clear() {
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}
