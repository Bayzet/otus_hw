package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     *sync.Mutex
	events map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]storage.Event),
		mu:     &sync.Mutex{},
	}
}

func (s Storage) countRows() int {
	return len(s.events)
}

func (s Storage) CreateEvent(ctx context.Context, e *storage.Event) error {
	s.mu.Lock()
	s.events[e.ID] = *e
	s.mu.Unlock()

	return nil
}

func (s Storage) UpdateEvent(ctx context.Context, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[e.ID]; !ok {
		return errors.Wrap(storage.ErrEventNotFound, fmt.Sprintf("Ошибка обновления события %v", e.ID))
	}

	s.events[e.ID] = *e

	return nil
}

func (s Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if _, ok := s.events[id]; !ok {
		return errors.Wrap(storage.ErrEventNotFound, fmt.Sprintf("Ошибка удаления события %v", id))
	}

	delete(s.events, id)

	return nil
}

func (s Storage) FindEventById(ctx context.Context, id uuid.UUID) *storage.Event {
	e, ok := s.events[id]
	if !ok {
		return nil
	}

	return &e
}

func (s Storage) ListEventsForDay(ctx context.Context, t time.Time) ([]storage.Event, error) {
	var events []storage.Event
	y, m, d := t.Date()

	for _, e := range s.events {
		ey, em, ed := e.Date.Date()
		if y == ey && m == em && d == ed {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s Storage) ListEventsForWeek(ctx context.Context, t time.Time) ([]storage.Event, error) {
	if t.Weekday() != time.Monday {
		return nil, errors.Wrap(storage.ErrDayNotMonday, fmt.Sprintf("Ошибка, переданный день - %v", t.Weekday()))
	}

	var events []storage.Event
	sevenDayHour, _ := time.ParseDuration("167h59m59s")
	firstDayOfWeek := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	lastDayOfWeek := firstDayOfWeek.Add(sevenDayHour)

	for _, e := range s.events {
		if e.Date == firstDayOfWeek || e.Date == lastDayOfWeek || (e.Date.After(firstDayOfWeek) && e.Date.Before(lastDayOfWeek)) {
			events = append(events, e)
		}
	}

	return events, nil
}

func (s Storage) ListEventsForMonth(ctx context.Context, t time.Time) ([]storage.Event, error) {
	var events []storage.Event

	for _, e := range s.events {
		if e.Date.Month() == t.Month() {
			events = append(events, e)
		}
	}

	return events, nil
}
