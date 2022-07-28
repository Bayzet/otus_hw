package storage

import (
	"errors"
)

var (
	ErrEventNotFound = errors.New("событие не найдено")
	ErrDayNotMonday  = errors.New("ожидается понедельник")
)
