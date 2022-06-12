package logger

import (
	"fmt"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"
)

type LogLevel string

func (ll LogLevel) Validate() error {
	switch ll {
	case consts.LogLevelInfo, consts.LogLevelError, consts.LogLevelDebug:
		return nil
	default:
		return ErrLogLevelNotValid
	}
}

type Logger struct {
	level LogLevel
	file  string
}

func New(level, file string) (*Logger, error) {
	l := LogLevel(level)
	err := l.Validate()
	if err != nil {
		return nil, fmt.Errorf("Переданный уровень логирования: %v. Ошибка: %w", level, err)
	}

	return &Logger{
		level: l,
		file:  file,
	}, nil
}

func (l Logger) Info(msg string) {
	fmt.Println(fmt.Sprintf("{\"INFO\":\"%v\"}", msg))
}

func (l Logger) Error(msg string) {
	fmt.Println(fmt.Sprintf("{\"ERROR\":\"%v\"}", msg))
}

func (l Logger) Debug(msg string) {
	if l.level == consts.LogLevelDebug {
		fmt.Println(fmt.Sprintf("{\"DEBUG\":\"%v\"}", msg))
	}
}
