package logger

//go:generate mockgen -source logger.go -destination logger_mock_gen.go -package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"
)

type Logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
}

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

type logger struct {
	level LogLevel
	file  *os.File
}

func New(level, file string) (*logger, error) { //nolint:golint,revive
	logLevel := LogLevel(level)

	if err := logLevel.Validate(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Переданный уровень логирования %v", level))
	}

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0o666) //nolint:varnamelen
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Не удалось обработать файл %v", file))
	}

	return &logger{
		level: logLevel,
		file:  f,
	}, nil
}

func (l logger) Info(msg string) {
	if l.level != consts.LogLevelError {
		err := l.saveMsg(consts.LogLevelInfo, msg)
		if err != nil {
			fmt.Printf("Ошибка записи логов: %v", err.Error())
		}
	}
}

func (l logger) Warn(msg string) {
	if l.level != consts.LogLevelError {
		err := l.saveMsg(consts.LogLevelWarn, msg)
		if err != nil {
			fmt.Printf("Ошибка записи логов: %v", err.Error())
		}
	}
}

func (l logger) Error(msg string) {
	err := l.saveMsg(consts.LogLevelError, msg)
	if err != nil {
		fmt.Printf("Ошибка записи логов: %v", err.Error())
	}
}

func (l logger) Debug(msg string) {
	if l.level == consts.LogLevelDebug {
		err := l.saveMsg(consts.LogLevelDebug, msg)
		if err != nil {
			fmt.Printf("Ошибка записи логов: %v", err.Error())
		}
	}
}

func (l logger) saveMsg(level, msg string) error {
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

func (l logger) Close() {
	l.file.Close()
}
