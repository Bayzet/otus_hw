package memorystorage

import (
	"context"
	"fmt"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var events = []storage.Event{
	storage.Event{ID: uuid.New(), Title: "event 1", Date: time.Date(2022, 5, 30, 1, 0, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 1.1", Date: time.Date(2022, 5, 30, 1, 1, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 2", Date: time.Date(2022, 5, 31, 1, 1, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 3", Date: time.Date(2022, 6, 1, 1, 0, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 4", Date: time.Date(2022, 6, 2, 1, 0, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 5", Date: time.Date(2022, 6, 3, 1, 0, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 7", Date: time.Date(2022, 6, 4, 1, 0, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 7", Date: time.Date(2022, 6, 5, 1, 0, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 8", Date: time.Date(2022, 6, 6, 1, 0, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 9", Date: time.Date(2022, 6, 7, 1, 0, 0, 0, time.UTC), User: 1},
	storage.Event{ID: uuid.New(), Title: "event 10", Date: time.Date(2022, 6, 8, 1, 0, 0, 0, time.UTC), User: 1},
}

func TestStorage_CreateEvent(t *testing.T) {
	s := New()

	err := s.CreateEvent(context.Background(), &events[0])
	require.NoError(t, err)
	require.Equal(t, 1, s.countRows())

	err = s.CreateEvent(context.Background(), &events[1])
	require.NoError(t, err)
	require.Equal(t, 2, s.countRows())
}

func TestStorage_FindEventById(t *testing.T) {
	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	e := s.FindEventById(ctx, events[3].ID)
	require.NotNil(t, e)
	require.EqualValues(t, &events[3], e)
}

func TestStorage_UpdateEvent(t *testing.T) {
	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	updEvent := events[2]
	updEvent.Title = "upd"

	err := s.UpdateEvent(ctx, &updEvent)
	require.NoError(t, err)

	e := s.FindEventById(ctx, updEvent.ID)
	require.NotNil(t, e)
	require.EqualValues(t, &updEvent, e)
}

func TestStorage_DeleteEvent(t *testing.T) {
	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	err := s.DeleteEvent(ctx, events[0].ID)
	require.NoError(t, err)

	e := s.FindEventById(ctx, events[0].ID)
	require.Nil(t, e)

	require.Equal(t, 10, s.countRows())
}

func TestStorage_ListEventsForDay(t *testing.T) {
	tests := []struct {
		in  time.Time
		exp []storage.Event
	}{
		{
			in:  time.Date(2022, 5, 30, 1, 0, 0, 0, time.UTC),
			exp: events[0:2],
		},
		{
			in:  time.Date(2022, 5, 31, 1, 0, 0, 0, time.UTC),
			exp: []storage.Event{events[2]},
		},
		{
			in:  time.Date(2022, 5, 29, 1, 0, 0, 0, time.UTC),
			exp: nil,
		},
	}

	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
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
	tests := []struct {
		in     time.Time
		exp    []storage.Event
		expErr error
	}{
		{
			in:     time.Date(2022, 5, 30, 1, 0, 0, 0, time.UTC),
			exp:    events[0:8],
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

	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
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
	tests := []struct {
		in  time.Time
		exp []storage.Event
	}{
		{
			in:  time.Date(2022, 5, 30, 1, 0, 0, 0, time.UTC),
			exp: events[0:3],
		},
		{
			in:  time.Date(2022, 6, 10, 1, 0, 0, 0, time.UTC),
			exp: events[3:],
		},
		{
			in:  time.Date(2022, 7, 10, 1, 0, 0, 0, time.UTC),
			exp: nil,
		},
	}

	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case: %v", i), func(t *testing.T) {
			eventList, err := s.ListEventsForMonth(ctx, tt.in)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.exp, eventList)
		})
	}
}
