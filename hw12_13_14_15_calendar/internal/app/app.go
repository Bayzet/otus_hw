package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/models"
)

type App struct {
	Storage Storage
}

type Storage interface {
	CreateEvent(context.Context, *models.Event) error
	UpdateEvent(context.Context, *models.Event) error
	DeleteEvent(context.Context, uuid.UUID) error
	ListEventsForDay(context.Context, time.Time) ([]models.Event, error)
	ListEventsForWeek(context.Context, time.Time) ([]models.Event, error)
	ListEventsForMonth(context.Context, time.Time) ([]models.Event, error)
	FindEventByID(ctx context.Context, uuid uuid.UUID) *models.Event
}

func New(storage Storage) *App {
	return &App{
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
