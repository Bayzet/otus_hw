package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"
)

type LogLevel string

func (ll LogLevel) Validate() error {
	switch ll {
	case consts.LogLevelInfo,
		consts.LogLevelWarn,
		consts.LogLevelError,
		consts.LogLevelDebug:
		return nil
	default:
		return ErrLogLevelNotValid
	}
}

type Logger struct {
	level LogLevel
	file  *os.File
}

func New(level, file string) (*Logger, error) {
	l := LogLevel(level)

	if err := l.Validate(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Переданный уровень логирования %v", level))
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Не удалось обработать файл %v", file))
	}

	return &Logger{
		level: l,
		file:  f,
	}, nil
}

func (l Logger) Info(msg string) {
	if l.level != consts.LogLevelError {
		err := l.saveMsg(consts.LogLevelInfo, msg)
		if err != nil {
			fmt.Printf("Ошибка записи логов: %v", err.Error())
		}
	}
}

func (l Logger) Warn(msg string) {
	if l.level != consts.LogLevelError {
		err := l.saveMsg(consts.LogLevelWarn, msg)
		if err != nil {
			fmt.Printf("Ошибка записи логов: %v", err.Error())
		}
	}
}

func (l Logger) Error(msg string) {
	err := l.saveMsg(consts.LogLevelError, msg)
	if err != nil {
		fmt.Printf("Ошибка записи логов: %v", err.Error())
	}
}

func (l Logger) Debug(msg string) {
	if l.level == consts.LogLevelDebug {
		err := l.saveMsg(consts.LogLevelDebug, msg)
		if err != nil {
			fmt.Printf("Ошибка записи логов: %v", err.Error())
		}
	}
}

func (l Logger) saveMsg(level, msg string) error {
	_, err := l.file.Seek(0, io.SeekEnd)
	if err != nil {
		return errors.Wrap(err, "Попытка передвинуть указатель файла")
	}

	_, err = l.file.WriteString(fmt.Sprintf(
		"{\"time\":\"%v\",\"level\":\"%v\",\"message\":\"%v\"}\n",
		time.Now().Format(time.RFC3339),
		level,
		msg,
	))
	if err != nil {
		return errors.Wrap(err, "Попытка записать в файл")
	}

	return nil
}

func (l Logger) Close() {
	l.file.Close()
}
