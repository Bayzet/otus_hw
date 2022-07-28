package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/models"

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

func (s Storage) CreateEvent(ctx context.Context, e *models.Event) error {
	query := "INSERT INTO events(id, title, date, user_id) values($1, $2, $3, $4)"

	_, err := s.db.ExecContext(ctx, query, e.ID, e.Title, e.Date, e.User)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) UpdateEvent(ctx context.Context, e *models.Event) error {
	query := "UPDATE events SET title = $1, date = $2 WHERE id = $3"

	_, err := s.db.ExecContext(ctx, query, e.Title, e.Date, e.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) DeleteEvent(ctx context.Context, e *models.Event) error {
	query := "DELETE FROM events WHERE id = $1"

	_, err := s.db.ExecContext(ctx, query, e.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) ListEventsForDay(ctx context.Context, t time.Time) ([]models.Event, error) {
	dayBegin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, time.UTC)

	return s.getEventsByDate(ctx, dayBegin, dayEnd)
}

func (s Storage) ListEventsForWeek(ctx context.Context, t time.Time) ([]models.Event, error) {
	if t.Weekday() != time.Monday {
		return nil, errors.Wrap(storage.ErrDayNotMonday, fmt.Sprintf("Номер переданного дня - %v", t.Weekday()))
	}

	sevenDayHour, _ := time.ParseDuration("167h59m59s")
	firstDayOfWeek := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	lastDayOfWeek := firstDayOfWeek.Add(sevenDayHour)

	return s.getEventsByDate(ctx, firstDayOfWeek, lastDayOfWeek)
}

func (s Storage) ListEventsForMonth(ctx context.Context, t time.Time) ([]models.Event, error) {
	duration, _ := time.ParseDuration("23h59m59s")
	firstDayOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1).Add(duration)

	return s.getEventsByDate(ctx, firstDayOfMonth, lastDayOfMonth)
}

func (s Storage) FindEventByID(ctx context.Context, id uuid.UUID) *models.Event {
	var e models.Event

	query := "SELECT * FROM events WHERE id = $1"
	row := s.db.QueryRowContext(ctx, query, id)

	var (
		eid   uuid.UUID
		title string
		date  time.Time
		user  int64
	)

	err := row.Scan(&eid, &title, &date, &user)
	if err != nil {
		return nil
	}

	e = models.Event{
		ID:    eid,
		Title: title,
		Date:  date,
		User:  int(user),
	}

	return &e
}

func (s Storage) getEventsByDate(ctx context.Context, begin, end time.Time) ([]models.Event, error) {
	var e []models.Event

	query := "SELECT * FROM events e WHERE e.date BETWEEN $1 AND $2"
	rows, err := s.db.QueryContext(ctx, query, begin, end)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		eid   uuid.UUID
		title string
		date  time.Time
		user  int64
	)

	for rows.Next() {
		err := rows.Scan(&eid, &title, &date, &user)
		if err != nil {
			return nil, err
		}

		e = append(e, models.Event{
			ID:    eid,
			Title: title,
			Date:  date,
			User:  int(user),
		})
	}

	return e, nil
}
