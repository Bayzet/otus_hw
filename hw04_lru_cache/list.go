package hw04lrucache

import (
	"github.com/google/uuid"
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	key   uuid.UUID
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	items map[uuid.UUID]*ListItem
	front *ListItem
	back  *ListItem
}

func NewList() List {
	l := new(list)
	l.items = make(map[uuid.UUID]*ListItem)
	return l
}

func (l *list) Len() int {
	return len(l.items)
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) push(v interface{}) *ListItem {
	li := ListItem{
		Value: v,
		key:   uuid.New(),
	}

	if l.Len() == 0 {
		l.front = &li
		l.back = &li
	}

	l.items[li.key] = &li

	return &li
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := l.push(v)

	if l.Front() != li {
		l.Front().Prev = li
		li.Next = l.Front()
	}

	l.front = li

	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := l.push(v)

	if l.Back() != li {
		l.Back().Next = li
		li.Prev = l.Back()
	}

	l.back = li

	return li
}

func (l *list) Remove(i *ListItem) {
	if l.Front() != i {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if l.Back() != i {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	delete(l.items, i.key)
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Front() != i {	
		if l.Back() != i {
			i.Next.Prev = i.Prev
		} else {
			l.back = i.Prev
		}

		i.Prev.Next = i.Next
		i.Next = l.Front()
		l.front = i
	}
}
