package hw04lrucache

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
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	l := new(list)
	return l
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := new(ListItem)
	li.Value = v

	if l.len == 0 {
		l.front = li
		l.back = li
	} else {
		li.Next = l.front
		l.front.Prev = li
	}

	l.len++

	l.front = li

	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := new(ListItem)
	li.Value = v

	if l.len == 0 {
		l.front = li
		l.back = li
	} else {
		li.Prev = l.back
		l.back.Next = li
	}

	l.len++

	l.back = li

	return li
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		i.Next.Prev = nil
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		i.Prev.Next = nil
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
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

	l.front.Prev = i
	i.Next = l.front
	i.Prev = nil
	l.front = i
}
