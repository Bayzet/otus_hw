package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
}

type Storage interface {
	CreateEvent(context.Context, *storage.Event) error
	UpdateEvent(context.Context, *storage.Event) error
	DeleteEvent(context.Context, uuid.UUID) error
	ListEventsForDay(context.Context, time.Time) ([]storage.Event, error)
	ListEventsForWeek(context.Context, time.Time) ([]storage.Event, error)
	ListEventsForMonth(context.Context, time.Time) ([]storage.Event, error)
	FindEventByID(ctx context.Context, uuid uuid.UUID) *storage.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
