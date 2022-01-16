package hw04lrucache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
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
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	items map[string]*ListItem
	front *ListItem
	back  *ListItem
}

func NewList() List {
	l := new(list)
	l.items = make(map[string]*ListItem)
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
	li := new(ListItem)
	li.Value = v
	hashKey := makeMD5(li.Value)

	l.items[hashKey] = li

	if l.itAddToEmptyList() {
		l.front = li
		l.back = li
	}

	return li
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := l.push(v)
	lFront := l.Front()

	if !l.itAddToEmptyList() {
		lFront.Prev = li
		li.Next = lFront
	}

	l.front = li

	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := l.push(v)
	lBack := l.Back()

	if !l.itAddToEmptyList() {
		li.Prev = lBack
		lBack.Next = li
	}

	l.back = li

	return li
}

func (l *list) Remove(i *ListItem) {
	hashKey := makeMD5(i.Value)

	prev := l.items[hashKey].Prev
	next := l.items[hashKey].Next

	if prev != nil {
		prev.Next = next
	} else {
		l.front = next
	}

	if next != nil {
		next.Prev = prev
	} else {
		l.back = prev
	}

	delete(l.items, hashKey)
}

func (l *list) MoveToFront(i *ListItem) {
	hashKey := makeMD5(i.Value)
	tmp := *l.items[hashKey]

	// Перепривязка соседних записей
	lFront := *l.front
	if tmp.Prev != nil {
		tmp.Next = &lFront
		tmp.Prev = nil
		l.front = &tmp
		l.Remove(&tmp)
	}
	l.front.Prev = &tmp
}

// https://gist.github.com/sergiotapia/8263278
func makeMD5(in interface{}) string {
	inByte, _ := getBytes(in)
	binHash := md5.Sum(inByte)

	return hex.EncodeToString(binHash[:])
}

// https://clck.ru/ahzvt
func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (l *list) itAddToEmptyList() bool {
	// Такая проверка продиктована выбраным способом добавления новых элементов.
	// Проверка идёт после добавления элемента
	// т.е. Len == 1 означает что только начали работать со списком
	// Выходит в момент первых проверок Len никогда не бывает == 0
	return l.Len() == 1
}
