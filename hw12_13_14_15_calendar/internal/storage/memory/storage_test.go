package memorystorage

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

var events = []storage.Event{
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

func TestStorage_CreateEvent(t *testing.T) {
	s := New()

	err := s.CreateEvent(context.Background(), &events[0])
	require.NoError(t, err)
	require.Equal(t, 1, s.countRows())

	err = s.CreateEvent(context.Background(), &events[1])
	require.NoError(t, err)
	require.Equal(t, 2, s.countRows())
}

func TestStorage_CreateEvent_race(t *testing.T) {
	s := New()

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			_ = s.CreateEvent(context.Background(), &events[0])
		}(wg)
	}
	wg.Wait()
}

func TestStorage_FindEventById(t *testing.T) {
	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	e := s.FindEventByID(ctx, events[3].ID)
	require.NotNil(t, e)
	require.EqualValues(t, &events[3], e)
}

func TestStorage_UpdateEvent_success(t *testing.T) {
	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	updEvent := events[2]
	updEvent.Title = "update"

	err := s.UpdateEvent(ctx, &updEvent)
	require.NoError(t, err)

	e := s.FindEventByID(ctx, updEvent.ID)
	require.NotNil(t, e)
	require.EqualValues(t, &updEvent, e)
}

func TestStorage_UpdateEvent_race(t *testing.T) {
	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	updEvent := events[2]
	updEvent.Title = "upd"

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			err := s.UpdateEvent(ctx, &updEvent)
			require.NoError(t, err)
		}(wg)
	}

	wg.Wait()
}

func TestStorage_UpdateEvent_error(t *testing.T) {
	ctx := context.Background()
	s := New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	updEvent := events[2]
	updEvent.Title = "upd"

	err := s.UpdateEvent(ctx, &updEvent)
	require.NoError(t, err)

	e := s.FindEventByID(ctx, updEvent.ID)
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

	e := s.FindEventByID(ctx, events[0].ID)
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
