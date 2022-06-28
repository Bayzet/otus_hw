package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/models"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/stdlib"
)

var eventsFixture = []models.Event{
	{ID: uuid.New(), Title: "event 1", Date: time.Date(2022, 5, 30, 1, 0, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 1.1", Date: time.Date(2022, 5, 30, 1, 1, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 2", Date: time.Date(2022, 5, 31, 1, 1, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 3", Date: time.Date(2022, 6, 1, 1, 0, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 4", Date: time.Date(2022, 6, 2, 1, 0, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 5", Date: time.Date(2022, 6, 3, 1, 0, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 7", Date: time.Date(2022, 6, 4, 1, 0, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 7", Date: time.Date(2022, 6, 5, 1, 0, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 8", Date: time.Date(2022, 6, 6, 1, 0, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 9", Date: time.Date(2022, 6, 7, 1, 0, 0, 0, time.UTC), User: 1},
	{ID: uuid.New(), Title: "event 10", Date: time.Date(2022, 6, 8, 1, 0, 0, 0, time.UTC), User: 1},
}

func initDB() *sql.DB {
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/calendar")
	if err != nil {
		log.Fatal(err)
	}

	err = db.PingContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("truncate table events")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("truncate table users")
	if err != nil {
		log.Fatal(err)
	}

	stmt, _ := db.Prepare(`INSERT INTO events(id, title, date, user_id)
		values($1, $2, $3, $4)`)
	defer stmt.Close()

	for _, e := range eventsFixture {
		_, err := stmt.Exec(e.ID, e.Title, e.Date, e.User)
		if err != nil {
			log.Fatal(err) // nolint:gocritic
		}
	}

	return db
}

func TestStorage_CreateEvent(t *testing.T) {
	ctx := context.Background()

	s := New(initDB())

	event := models.Event{
		ID:    uuid.New(),
		Title: "Test event 1",
		Date:  time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
		User:  1,
	}

	err := s.CreateEvent(ctx, &event)
	require.NoError(t, err)

	row := s.db.QueryRowContext(ctx, "SELECT id FROM events WHERE id = $1", event.ID)

	var id uuid.UUID
	_ = row.Scan(&id)
	require.Equal(t, event.ID, id)
}

func TestStorage_UpdateEvent(t *testing.T) {
	ctx := context.Background()

	s := New(initDB())

	updEvent := eventsFixture[0]
	updEvent.Title = "upd"

	err := s.UpdateEvent(ctx, &updEvent)
	require.NoError(t, err)

	row := s.db.QueryRowContext(ctx, "SELECT id, title FROM events WHERE id = $1", updEvent.ID)
	var id uuid.UUID
	var title string

	_ = row.Scan(&id, &title)
	require.Equal(t, updEvent.ID, id)
	require.Equal(t, updEvent.Title, title)
}

func TestStorage_DeleteEvent(t *testing.T) {
	ctx := context.Background()

	s := New(initDB())

	err := s.DeleteEvent(ctx, &eventsFixture[0])
	require.NoError(t, err)

	row := s.db.QueryRowContext(ctx, "SELECT id FROM events WHERE id = $1", eventsFixture[0].ID)
	var id *uuid.UUID
	_ = row.Scan(id)
	require.Nil(t, id)
}

func TestStorage_FindEventById(t *testing.T) {
	tests := []struct {
		in  uuid.UUID
		exp *models.Event
	}{
		{
			eventsFixture[2].ID,
			&eventsFixture[2],
		},
		{
			uuid.New(),
			nil,
		},
	}

	ctx := context.Background()

	s := New(initDB())

	for _, tt := range tests {
		event := s.FindEventByID(ctx, tt.in)
		require.EqualValues(t, tt.exp, event)
	}
}

func TestStorage_ListEventsForDay(t *testing.T) {
	ctx := context.Background()

	s := New(initDB())

	tests := []struct {
		in  time.Time
		exp []models.Event
	}{
		{
			in:  time.Date(2022, 5, 30, 1, 0, 0, 0, time.UTC),
			exp: eventsFixture[0:2],
		},
		{
			in:  time.Date(2022, 5, 31, 1, 0, 0, 0, time.UTC),
			exp: []models.Event{eventsFixture[2]},
		},
		{
			in:  time.Date(2022, 5, 29, 1, 0, 0, 0, time.UTC),
			exp: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case: %v", i), func(t *testing.T) {
			eventList, err := s.ListEventsForDay(ctx, tt.in)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.exp, eventList)
		})
	}
}

func TestStorage_ListEventsForWeek(t *testing.T) {
	ctx := context.Background()

	s := New(initDB())

	tests := []struct {
		in     time.Time
		exp    []models.Event
		expErr error
	}{
		{
			in:     time.Date(2022, 5, 30, 1, 0, 0, 0, time.UTC),
			exp:    eventsFixture[0:8],
			expErr: nil,
		},
		{
			in:     time.Date(2022, 6, 13, 1, 0, 0, 0, time.UTC),
			exp:    nil,
			expErr: nil,
		},
		{
			in:     time.Date(2022, 5, 31, 1, 0, 0, 0, time.UTC),
			exp:    nil,
			expErr: storage.ErrDayNotMonday,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case: %v", i), func(t *testing.T) {
			eventList, err := s.ListEventsForWeek(ctx, tt.in)
			require.ErrorIs(t, err, tt.expErr)
			require.ElementsMatch(t, tt.exp, eventList)
		})
	}
}

func TestStorage_ListEventsForMonth(t *testing.T) {
	ctx := context.Background()

	s := New(initDB())

	tests := []struct {
		in  time.Time
		exp []models.Event
	}{
		{
			in:  time.Date(2022, 5, 30, 1, 0, 0, 0, time.UTC),
			exp: eventsFixture[0:3],
		},
		{
			in:  time.Date(2022, 6, 10, 1, 0, 0, 0, time.UTC),
			exp: eventsFixture[3:],
		},
		{
			in:  time.Date(2022, 7, 10, 1, 0, 0, 0, time.UTC),
			exp: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case: %v", i), func(t *testing.T) {
			eventList, err := s.ListEventsForMonth(ctx, tt.in)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.exp, eventList)
		})
	}
}
