package logger

import (
	"errors"
	"fmt"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"
)

var (
	ErrLogLevelNotValid = errors.New(
		fmt.Sprintf(
			"Ошибка валидации. Поддерживаемые уровни логирования: %v, %v, %v",
			consts.LogLevelInfo,
			consts.LogLevelError,
			consts.LogLevelDebug))
)
