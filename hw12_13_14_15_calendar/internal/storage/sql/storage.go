package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) CreateEvent(ctx context.Context, e *storage.Event) error {
	query := "INSERT INTO events(`id`, `title`, `date`, user_id) values(?, ?, ?, ?)"

	_, err := s.db.ExecContext(ctx, query, e.ID, e.Title, e.Date, e.User)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) UpdateEvent(ctx context.Context, e *storage.Event) error {
	query := "UPDATE events e SET e.title = ?, e.date = ? WHERE e.id = ?"

	_, err := s.db.ExecContext(ctx, query, e.Title, e.Date, e.ID.String())
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) DeleteEvent(ctx context.Context, e *storage.Event) error {
	query := "DELETE FROM events WHERE id = ?"

	_, err := s.db.ExecContext(ctx, query, e.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) ListEventsForDay(ctx context.Context, t time.Time) ([]storage.Event, error) {
	dayBegin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, time.UTC)

	return s.getEventsByDate(ctx, dayBegin, dayEnd)
}

func (s Storage) ListEventsForWeek(ctx context.Context, t time.Time) ([]storage.Event, error) {
	if t.Weekday() != time.Monday {
		return nil, errors.Wrap(storage.ErrDayNotMonday, fmt.Sprintf("Номер переданного дня - %v", t.Weekday()))
	}

	sevenDayHour, _ := time.ParseDuration("167h59m59s")
	firstDayOfWeek := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	lastDayOfWeek := firstDayOfWeek.Add(sevenDayHour)

	return s.getEventsByDate(ctx, firstDayOfWeek, lastDayOfWeek)
}

func (s Storage) ListEventsForMonth(ctx context.Context, t time.Time) ([]storage.Event, error) {
	duration, _ := time.ParseDuration("23h59m59s")
	firstDayOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1).Add(duration)

	return s.getEventsByDate(ctx, firstDayOfMonth, lastDayOfMonth)
}

func (s Storage) FindEventByID(ctx context.Context, id uuid.UUID) *storage.Event {
	var e storage.Event

	query := "SELECT * FROM events e WHERE id = ?"
	row := s.db.QueryRowContext(ctx, query, id)

	var (
		eid   uuid.UUID
		title string
		date  []byte
		user  int64
	)

	err := row.Scan(&eid, &title, &date, &user)
	if err != nil {
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", string(date))
	if err != nil {
		log.Fatal(err)
	}

	e = storage.Event{
		ID:    eid,
		Title: title,
		Date:  t,
		User:  int(user),
	}

	return &e
}

func (s Storage) getEventsByDate(ctx context.Context, begin, end time.Time) ([]storage.Event, error) {
	var e []storage.Event

	query := "SELECT * FROM events e WHERE date BETWEEN ? AND ?"
	rows, err := s.db.QueryContext(ctx, query, begin, end)
	if err != nil {
		return nil, err
	}

	var (
		eid   uuid.UUID
		title string
		date  []byte
		user  int64
	)

	for rows.Next() {
		err := rows.Scan(&eid, &title, &date, &user)
		if err != nil {
			return nil, err
		}

		t, err := time.Parse("2006-01-02 15:04:05", string(date))
		if err != nil {
			return nil, err
		}

		e = append(e, storage.Event{
			ID:    eid,
			Title: title,
			Date:  t,
			User:  int(user),
		})
	}

	return e, nil
}
