package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/google/uuid"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/gen/pb/calendarpb"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	app2 "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/golang/mock/gomock"
)

var events = []models.Event{
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

func convertEvent(e models.Event) *calendarpb.Event {
	return &calendarpb.Event{
		Id:    e.ID.String(),
		Title: e.Title,
		Date: &calendarpb.DateTime{
			Year:   int32(e.Date.Year()),
			Month:  int32(e.Date.Month()),
			Day:    int32(e.Date.Day()),
			Hour:   int32(e.Date.Hour()),
			Minute: int32(e.Date.Minute()),
		},
		UserId: int64(e.User),
	}
}

func TestServer_CreateEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := app2.New(memorystorage.New())

	mockLogger := logger.NewMockLogger(ctrl)

	server := NewServer(app, "127.0.0.1", "50051", mockLogger)
	defer server.Stop(ctx)

	go server.Run(ctx)

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	require.NoError(t, err)

	client := calendarpb.NewCalendarClient(conn)

	tests := []struct {
		name   string
		in     calendarpb.Event
		expRes bool
		expErr bool
	}{
		{
			name: "Успешное создание события",
			in: calendarpb.Event{
				Title: "test event",
				Date: &calendarpb.DateTime{
					Year:   2022,
					Month:  5,
					Day:    5,
					Hour:   10,
					Minute: 0,
				},
				UserId: 0,
			},
			expRes: true,
			expErr: false,
		},
		{
			name:   "Данные не полные",
			in:     calendarpb.Event{},
			expRes: false,
			expErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			e := calendarpb.CreateEventRequest{Event: &tt.in}

			res, err := client.CreateEvent(ctx, &e)
			if tt.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expRes, res.Success)
		})
	}
}

func TestServer_UpdateEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := app2.New(memorystorage.New())

	id, err := app.CreateEvent(ctx, "title 1", time.Now(), 1)
	require.NoError(t, err)

	mockLogger := logger.NewMockLogger(ctrl)

	server := NewServer(app, "127.0.0.1", "50051", mockLogger)
	defer server.Stop(ctx)

	go server.Run(ctx)

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	require.NoError(t, err)

	client := calendarpb.NewCalendarClient(conn)

	tests := []struct {
		name   string
		in     calendarpb.Event
		expRes bool
		expErr bool
	}{
		{
			name: "ИД существует",
			in: calendarpb.Event{
				Id:    id.String(),
				Title: "test event",
				Date: &calendarpb.DateTime{
					Year:   2022,
					Month:  5,
					Day:    5,
					Hour:   10,
					Minute: 0,
				},
				UserId: 1,
			},
			expRes: true,
			expErr: false,
		},
		{
			name: "ИД не существует",
			in: calendarpb.Event{
				Id:    uuid.New().String(),
				Title: "test event",
				Date: &calendarpb.DateTime{
					Year:   2022,
					Month:  5,
					Day:    5,
					Hour:   10,
					Minute: 0,
				},
				UserId: 1,
			},
			expRes: false,
			expErr: false,
		},
		{
			name: "ИД не валидный",
			in: calendarpb.Event{
				Id:    "not-valid-uuid",
				Title: "test event",
				Date: &calendarpb.DateTime{
					Year:   2022,
					Month:  5,
					Day:    5,
					Hour:   10,
					Minute: 0,
				},
				UserId: 1,
			},
			expRes: false,
			expErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			e := calendarpb.UpdateEventRequest{Event: &tt.in}

			res, err := client.UpdateEvent(ctx, &e)
			if tt.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expRes, res.Success)
		})
	}
}

func TestServer_DeleteEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := app2.New(memorystorage.New())

	id, err := app.CreateEvent(ctx, "title 1", time.Now(), 1)
	require.NoError(t, err)

	mockLogger := logger.NewMockLogger(ctrl)

	server := NewServer(app, "127.0.0.1", "50051", mockLogger)
	defer server.Stop(ctx)

	go server.Run(ctx)

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	require.NoError(t, err)

	client := calendarpb.NewCalendarClient(conn)

	tests := []struct {
		name   string
		in     string
		expRes bool
		expErr bool
	}{
		{
			name:   "ИД существует",
			in:     id.String(),
			expRes: true,
			expErr: false,
		},
		{
			name:   "ИД не существует",
			in:     uuid.New().String(),
			expRes: false,
			expErr: false,
		},
		{
			name:   "ИД не валидный",
			in:     "not-valid-uuid",
			expRes: false,
			expErr: false,
		},
	}

	for _, tt := range tests {
		mockLogger.EXPECT().Info(gomock.Any())
		e := calendarpb.DeleteEventRequest{Id: tt.in}

		res, err := client.DeleteEvent(ctx, &e)
		if tt.expErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, tt.expRes, res.Success)
	}
}

func TestServer_ListEventsForDay(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := memorystorage.New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	app := app2.New(s)

	mockLogger := logger.NewMockLogger(ctrl)

	server := NewServer(app, "127.0.0.1", "50051", mockLogger)
	defer server.Stop(ctx)

	go server.Run(ctx)

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	require.NoError(t, err)

	client := calendarpb.NewCalendarClient(conn)

	tests := []struct {
		in  calendarpb.DateTime
		exp []*calendarpb.Event
	}{
		{
			in: calendarpb.DateTime{Year: 2022, Month: 5, Day: 30, Hour: 1, Minute: 0},
			exp: []*calendarpb.Event{
				convertEvent(events[0]),
				convertEvent(events[1]),
			},
		},
		{
			in: calendarpb.DateTime{Year: 2022, Month: 5, Day: 31, Hour: 1, Minute: 0},
			exp: []*calendarpb.Event{
				convertEvent(events[2]),
			},
		},
		{
			in:  calendarpb.DateTime{Year: 2022, Month: 5, Day: 29, Hour: 1, Minute: 0},
			exp: nil,
		},
	}

	for _, tt := range tests {
		mockLogger.EXPECT().Info(gomock.Any())
		e := calendarpb.ListEventsForDayRequest{Date: &tt.in}

		res, err := client.ListEventsForDay(ctx, &e)
		require.NoError(t, err)
		require.ElementsMatch(t, res.GetEvents(), tt.exp)
	}
}

func TestServer_ListEventsForWeek(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := memorystorage.New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	app := app2.New(s)

	mockLogger := logger.NewMockLogger(ctrl)

	server := NewServer(app, "127.0.0.1", "50051", mockLogger)
	defer server.Stop(ctx)

	go server.Run(ctx)

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	require.NoError(t, err)

	client := calendarpb.NewCalendarClient(conn)

	tests := []struct {
		in     calendarpb.DateTime
		exp    []*calendarpb.Event
		expErr bool
	}{
		{
			in: calendarpb.DateTime{Year: 2022, Month: 5, Day: 30, Hour: 1, Minute: 0},
			exp: []*calendarpb.Event{
				convertEvent(events[0]),
				convertEvent(events[1]),
				convertEvent(events[2]),
				convertEvent(events[3]),
				convertEvent(events[4]),
				convertEvent(events[5]),
				convertEvent(events[6]),
				convertEvent(events[7]),
			},
		},
		{
			in:  calendarpb.DateTime{Year: 2022, Month: 6, Day: 13, Hour: 1, Minute: 0},
			exp: nil,
		},
		{
			in:     calendarpb.DateTime{Year: 2022, Month: 5, Day: 31, Hour: 1, Minute: 0},
			exp:    nil,
			expErr: true,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			e := calendarpb.ListEventsForWeekRequest{Date: &tt.in}

			res, err := client.ListEventsForWeek(ctx, &e)
			if tt.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.ElementsMatch(t, res.GetEvents(), tt.exp)
		})
	}
}

func TestServer_ListEventsForMonth(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := memorystorage.New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	app := app2.New(s)

	mockLogger := logger.NewMockLogger(ctrl)

	server := NewServer(app, "127.0.0.1", "50051", mockLogger)
	defer server.Stop(ctx)

	go server.Run(ctx)

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	require.NoError(t, err)

	client := calendarpb.NewCalendarClient(conn)

	tests := []struct {
		in  calendarpb.DateTime
		exp []*calendarpb.Event
	}{
		{
			in: calendarpb.DateTime{Year: 2022, Month: 5, Day: 30, Hour: 1, Minute: 0},
			exp: []*calendarpb.Event{
				convertEvent(events[0]),
				convertEvent(events[1]),
				convertEvent(events[2]),
			},
		},
		{
			in: calendarpb.DateTime{Year: 2022, Month: 6, Day: 10, Hour: 1, Minute: 0},
			exp: []*calendarpb.Event{
				convertEvent(events[3]),
				convertEvent(events[4]),
				convertEvent(events[5]),
				convertEvent(events[6]),
				convertEvent(events[7]),
				convertEvent(events[8]),
				convertEvent(events[9]),
				convertEvent(events[10]),
			},
		},
		{
			in:  calendarpb.DateTime{Year: 2022, Month: 7, Day: 10, Hour: 1, Minute: 0},
			exp: nil,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			e := calendarpb.ListEventsForMonthRequest{Date: &tt.in}

			res, err := client.ListEventsForMonth(ctx, &e)
			require.NoError(t, err)
			require.ElementsMatch(t, res.GetEvents(), tt.exp)
		})
	}
}

func TestServer_FindEventByID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := app2.New(memorystorage.New())

	timeNow := time.Now()
	id, err := app.CreateEvent(ctx, "title 1", timeNow, 1)
	require.NoError(t, err)

	mockLogger := logger.NewMockLogger(ctrl)

	server := NewServer(app, "127.0.0.1", "50051", mockLogger)
	defer server.Stop(ctx)

	go server.Run(ctx)

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	require.NoError(t, err)

	client := calendarpb.NewCalendarClient(conn)

	tests := []struct {
		name   string
		in     string
		exp    *calendarpb.Event
		expErr bool
	}{
		{
			name: "ИД существует",
			in:   id.String(),
			exp: &calendarpb.Event{
				Id:    id.String(),
				Title: "title 1",
				Date: &calendarpb.DateTime{
					Year:   int32(timeNow.Year()),
					Month:  int32(timeNow.Month()),
					Day:    int32(timeNow.Day()),
					Hour:   int32(timeNow.Hour()),
					Minute: int32(timeNow.Minute()),
				},
				UserId: 1,
			},
			expErr: false,
		},
		{
			name:   "ИД не существует",
			in:     uuid.New().String(),
			exp:    &calendarpb.Event{},
			expErr: true,
		},
		{
			name:   "ИД не валидно",
			in:     "not-valid-uuid",
			exp:    &calendarpb.Event{},
			expErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			e := calendarpb.FindEventByIDRequest{Id: tt.in}

			res, err := client.FindEventByID(ctx, &e)
			if tt.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.exp.GetId(), res.GetId())
			require.Equal(t, tt.exp.GetTitle(), res.GetTitle())
			require.Equal(t, tt.exp.GetDate(), res.GetDate())
			require.Equal(t, tt.exp.GetUserId(), res.GetUserId())
		})
	}
}
