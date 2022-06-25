package logger

import (
	"fmt"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"
)

var (
	ErrLogLevelNotValid = fmt.Errorf(
		"ошибка валидации. Поддерживаемые уровни логирования: %v, %v, %v",
		consts.LogLevelInfo,
		consts.LogLevelError,
		consts.LogLevelDebug)
	ErrStorageTypeNotValid = fmt.Errorf(
		"ошибка валидации. Поддерживаемые типы storage: %v, %v",
		consts.StorageTypeMemory,
		consts.StorageTypeSQL)
)
