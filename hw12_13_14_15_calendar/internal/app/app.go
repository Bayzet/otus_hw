package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/models"
)

type Application interface {
	CreateEvent(context.Context, string, time.Time, int) (uuid.UUID, error)
	UpdateEvent(context.Context, uuid.UUID, string, time.Time, int) error
	DeleteEvent(context.Context, uuid.UUID) error
	ListEventsForDay(context.Context, time.Time) ([]models.Event, error)
	ListEventsForWeek(context.Context, time.Time) ([]models.Event, error)
	ListEventsForMonth(context.Context, time.Time) ([]models.Event, error)
	FindEventByID(context.Context, uuid.UUID) *models.Event
}

type App struct {
	Repository Repository
}

type Repository interface {
	CreateEvent(context.Context, *models.Event) error
	UpdateEvent(context.Context, *models.Event) error
	DeleteEvent(context.Context, uuid.UUID) error
	ListEventsForDay(context.Context, time.Time) ([]models.Event, error)
	ListEventsForWeek(context.Context, time.Time) ([]models.Event, error)
	ListEventsForMonth(context.Context, time.Time) ([]models.Event, error)
	FindEventByID(ctx context.Context, uuid uuid.UUID) *models.Event
}

func New(repo Repository) *App {
	return &App{
		Repository: repo,
	}
}

func (a *App) CreateEvent(ctx context.Context, title string, date time.Time, userID int) (uuid.UUID, error) {
	event := models.Event{
		ID:    uuid.New(),
		Title: title,
		Date:  date,
		User:  userID,
	}

	return event.ID, errors.Wrap(a.Repository.CreateEvent(ctx, &event), "Ошибка создания события")
}

func (a *App) UpdateEvent(ctx context.Context, id uuid.UUID, title string, date time.Time, userID int) error {
	event := models.Event{
		ID:    id,
		Title: title,
		Date:  date,
		User:  userID,
	}

	return errors.Wrap(a.Repository.UpdateEvent(ctx, &event), "Ошибка обновления события")
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return errors.Wrap(a.Repository.DeleteEvent(ctx, id), "Ошибка удаления события")
}

func (a *App) ListEventsForDay(ctx context.Context, date time.Time) ([]models.Event, error) {
	events, err := a.Repository.ListEventsForDay(ctx, date)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка ListEventsForDay")
	}

	return events, nil
}

func (a *App) ListEventsForWeek(ctx context.Context, date time.Time) ([]models.Event, error) {
	events, err := a.Repository.ListEventsForWeek(ctx, date)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка ListEventsForWeek")
	}

	return events, nil
}

func (a *App) ListEventsForMonth(ctx context.Context, date time.Time) ([]models.Event, error) {
	events, err := a.Repository.ListEventsForMonth(ctx, date)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка ListEventsForMonth")
	}

	return events, nil
}

func (a *App) FindEventByID(ctx context.Context, id uuid.UUID) *models.Event {
	return a.Repository.FindEventByID(ctx, id)
}
