package storage

import (
	"errors"
)

var (
	ErrUpdate       = errors.New("Ошибка обновления события")
	ErrDelete       = errors.New("Ошибка удаления события")
	ErrDayNotMonday = errors.New("Ошибка, ожидается понедельник")
)
