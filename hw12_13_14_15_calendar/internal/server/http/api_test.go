package internalhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"

	"github.com/gorilla/mux"

	models2 "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/models"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/server/http/models"
	"golang.org/x/net/context"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/golang/mock/gomock"

	app2 "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
)

var events = []models2.Event{
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

func TestCalendarAPI_CreateEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := app2.New(memorystorage.New())

	mockLogger := logger.NewMockLogger(ctrl)
	router := mux.NewRouter()
	server := NewServer(router, "127.0.0.1", "1234", mockLogger)
	calendarAPI := CalendarApp{app}
	calendarAPI.RegisterHTTPHandlers(router)
	defer server.Stop(ctx)

	go server.Start(ctx)

	// Слип, иначе тест завершается быстрее, чем поднимается HTTP сервер
	time.Sleep(1 * time.Millisecond)

	client := http.DefaultClient

	tests := []struct {
		name       string
		in         string
		expEmptyID bool
		expErr     string
	}{
		{
			name:       "Данные корректны и полны",
			in:         `{"title":"test 1","date":"2022-05-01T00:04:58Z","user":1}`,
			expEmptyID: false,
		},
		{
			name:   "user не валиден",
			in:     `{"title":"test 1","date":"2022-05-01T00:04:58Z","user":"текст"}`,
			expErr: "Request body contains an invalid value for the \"user\" field (at position 67)",
		},
		{
			name:   "user не валиден",
			in:     `{"title":"test 1","date":"2022-05-01T00:04:58Z","user":0}`,
			expErr: "1 error occurred:\n\t* поле user не может быть пустым или равно 0\n\n",
		},
		{
			name:       "title не валиден",
			in:         `{"title":0,"date":"2022-05-01T00:04:58Z","user":1}`,
			expEmptyID: true,
			expErr:     "Request body contains an invalid value for the \"title\" field (at position 10)",
		},
		{
			name:       "date не валиден",
			in:         `{"title":"test 1","user":1,"date":""}`,
			expEmptyID: true,
			expErr:     "parsing time \"\\\"\\\"\" as \"\\\"2006-01-02T15:04:05Z07:00\\\"\": cannot parse \"\\\"\" as \"2006\"",
		},
		{
			name:       "Пустое сообщение",
			in:         ``,
			expEmptyID: true,
			expErr:     "Request body must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			res, err := client.Post("http://127.0.0.1:1234/event/create", "application/json", bytes.NewReader([]byte(tt.in)))
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			if tt.expErr != "" {
				require.Equal(t, tt.expErr, string(body))
			} else {
				var actualRes models.Response
				json.Unmarshal(body, &actualRes)

				var emptyUUID uuid.UUID
				if tt.expEmptyID {
					require.Equal(t, emptyUUID, actualRes.ID)
				} else {
					require.NotEqual(t, emptyUUID, actualRes.ID)
				}
			}

		})
	}
}

func TestCalendarAPI_UpdateEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := app2.New(memorystorage.New())
	id, err := app.CreateEvent(ctx, "title", time.Now(), 1)
	require.NoError(t, err)

	mockLogger := logger.NewMockLogger(ctrl)
	router := mux.NewRouter()
	server := NewServer(router, "127.0.0.1", "1234", mockLogger)
	calendarAPI := CalendarApp{app}
	calendarAPI.RegisterHTTPHandlers(router)
	defer server.Stop(ctx)

	go server.Start(ctx)

	// Слип, иначе тест завершается быстрее, чем поднимается HTTP сервер
	time.Sleep(1 * time.Millisecond)

	client := http.DefaultClient

	tests := []struct {
		name   string
		in     string
		expID  bool
		expErr string
	}{
		{
			name:  "Данные корректны и полны. Событие существует",
			in:    fmt.Sprintf(`{"id":"%s","title":"test 2","date":"2022-05-01T00:04:58Z","user":1}`, id.String()),
			expID: true,
		},
		{
			name:   "Данные корректны, но неполны, title отсутствует. Событие существует",
			in:     fmt.Sprintf(`{"id":"%s","date":"2022-05-01T00:04:58Z","user":1}`, id.String()),
			expErr: "1 error occurred:\n\t* поле title не может быть пустым\n\n",
		},
		{
			name:   "Данные корректны, но неполны, title и date отсутствуют. Событие существует",
			in:     fmt.Sprintf(`{"id":"%s","user":1}`, id.String()),
			expErr: "2 errors occurred:\n\t* поле title не может быть пустым\n\t* поле date не может быть пустым\n\n",
		},
		{
			name:   "Данные корректны и полны. Событие отсутствует",
			in:     `{"id":"1dc565af-026b-403f-b1f5-27bee27a302a","title":"test 2","date":"2022-05-01T00:04:58Z","user":1}`,
			expErr: "Ошибка обновления события: Ошибка обновления события 1dc565af-026b-403f-b1f5-27bee27a302a: событие не найдено",
		},
		{
			name:   "Пустое сообщение",
			in:     ``,
			expErr: "Request body must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			r := bytes.NewReader([]byte(tt.in))
			res, err := client.Post("http://127.0.0.1:1234/event/update", "application/json", r)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			if tt.expErr != "" {
				require.Equal(t, tt.expErr, string(body))
			} else {
				var actualRes models.Response
				json.Unmarshal(body, &actualRes)

				require.Equal(t, id, actualRes.ID)

			}
		})
	}
}

func TestCalendarAPI_DeleteEvent(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := app2.New(memorystorage.New())
	id, err := app.CreateEvent(ctx, "title", time.Now(), 1)
	require.NoError(t, err)

	mockLogger := logger.NewMockLogger(ctrl)
	router := mux.NewRouter()
	server := NewServer(router, "127.0.0.1", "1234", mockLogger)
	calendarAPI := CalendarApp{app}
	calendarAPI.RegisterHTTPHandlers(router)
	defer server.Stop(ctx)

	go server.Start(ctx)

	// Слип, иначе тест завершается быстрее, чем поднимается HTTP сервер
	time.Sleep(1 * time.Millisecond)

	client := http.DefaultClient

	tests := []struct {
		name   string
		in     string
		expID  bool
		expErr string
	}{
		{
			name:  "Событие существует",
			in:    fmt.Sprintf(`{"id":"%s"}`, id.String()),
			expID: true,
		},
		{
			name:   "Событие отсутствует",
			in:     `{"id":"1dc565af-026b-403f-b1f5-27bee27a302a"}`,
			expErr: "Ошибка удаления события: Ошибка удаления события 1dc565af-026b-403f-b1f5-27bee27a302a: событие не найдено",
		},
		{
			name:   "ИД события не валидно",
			in:     `{"id":"not-valid-uuid"}`,
			expErr: "invalid UUID length: 14",
		},
		{
			name:   "ИД события не валидно",
			in:     `{"id":0}`,
			expErr: "Request body contains an invalid value for the \"id\" field (at position 7)",
		},
		{
			name:   "Пустое сообщение",
			in:     ``,
			expErr: "Request body must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			r := bytes.NewReader([]byte(tt.in))
			res, err := client.Post("http://127.0.0.1:1234/event/delete", "application/json", r)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			if tt.expErr != "" {
				require.Equal(t, tt.expErr, string(body))
			} else {
				var actualRes models.Response
				json.Unmarshal(body, &actualRes)

				require.Equal(t, id, actualRes.ID)
			}
		})
	}
}

func TestCalendarAPI_ListEventsForDay(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := memorystorage.New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	app := app2.New(s)

	mockLogger := logger.NewMockLogger(ctrl)
	router := mux.NewRouter()
	server := NewServer(router, "127.0.0.1", "1234", mockLogger)
	calendarAPI := CalendarApp{app}
	calendarAPI.RegisterHTTPHandlers(router)
	defer server.Stop(ctx)

	go server.Start(ctx)

	// Слип, иначе тест завершается быстрее, чем поднимается HTTP сервер
	time.Sleep(1 * time.Millisecond)

	client := http.DefaultClient

	tests := []struct {
		in  string
		exp []models.Event
	}{
		{
			in: `{"date":"2022-05-30T01:00:00Z"}`,
			exp: []models.Event{
				models.Event(events[0]),
				models.Event(events[1]),
			},
		},
		{
			in: `{"date":"2022-05-31T01:00:00Z"}`,
			exp: []models.Event{
				models.Event(events[2]),
			},
		},
		{
			in:  `{"date":"2022-05-29T01:00:00Z"}`,
			exp: []models.Event{},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			r := bytes.NewReader([]byte(tt.in))
			res, err := client.Post("http://127.0.0.1:1234/event/list/"+consts.TypeListEventByDay, "application/json", r)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			var actualRes []models.Event
			json.Unmarshal(body, &actualRes)

			require.ElementsMatch(t, tt.exp, actualRes)
		})
	}
}

func TestCalendarAPI_ListEventsForWeek(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := memorystorage.New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	app := app2.New(s)

	mockLogger := logger.NewMockLogger(ctrl)
	router := mux.NewRouter()
	server := NewServer(router, "127.0.0.1", "1234", mockLogger)
	calendarAPI := CalendarApp{app}
	calendarAPI.RegisterHTTPHandlers(router)
	defer server.Stop(ctx)

	go server.Start(ctx)

	// Слип, иначе тест завершается быстрее, чем поднимается HTTP сервер
	time.Sleep(1 * time.Millisecond)

	client := http.DefaultClient

	tests := []struct {
		in  string
		exp []models.Event
	}{
		{
			in: `{"date":"2022-05-30T01:00:00Z"}`,
			exp: []models.Event{
				models.Event(events[0]),
				models.Event(events[1]),
				models.Event(events[2]),
				models.Event(events[3]),
				models.Event(events[4]),
				models.Event(events[5]),
				models.Event(events[6]),
				models.Event(events[7]),
			},
		},
		{
			in:  `{"date":"2022-05-31T01:00:00Z"}`,
			exp: []models.Event{},
		},
		{
			in:  `{"date":"2022-05-29T01:00:00Z"}`,
			exp: []models.Event{},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			r := bytes.NewReader([]byte(tt.in))
			res, err := client.Post("http://127.0.0.1:1234/event/list/"+consts.TypeListEventByWeek, "application/json", r)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			var actualRes []models.Event
			json.Unmarshal(body, &actualRes)

			require.ElementsMatch(t, tt.exp, actualRes)
		})
	}
}

func TestCalendarAPI_ListEventsForMonth(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := memorystorage.New()
	for _, e := range events {
		_ = s.CreateEvent(ctx, &e)
	}

	app := app2.New(s)

	mockLogger := logger.NewMockLogger(ctrl)
	router := mux.NewRouter()
	server := NewServer(router, "127.0.0.1", "1234", mockLogger)
	calendarAPI := CalendarApp{app}
	calendarAPI.RegisterHTTPHandlers(router)
	defer server.Stop(ctx)

	go server.Start(ctx)

	// Слип, иначе тест завершается быстрее, чем поднимается HTTP сервер
	time.Sleep(1 * time.Millisecond)

	client := http.DefaultClient

	tests := []struct {
		in  string
		exp []models.Event
	}{
		{
			in: `{"date":"2022-05-30T01:00:00Z"}`,
			exp: []models.Event{
				models.Event(events[0]),
				models.Event(events[1]),
				models.Event(events[2]),
			},
		},
		{
			in: `{"date":"2022-06-10T01:00:00Z"}`,
			exp: []models.Event{
				models.Event(events[3]),
				models.Event(events[4]),
				models.Event(events[5]),
				models.Event(events[6]),
				models.Event(events[7]),
				models.Event(events[8]),
				models.Event(events[9]),
				models.Event(events[10]),
			},
		},
		{
			in:  `{"date":"2022-07-10T01:00:00Z"}`,
			exp: []models.Event{},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			r := bytes.NewReader([]byte(tt.in))
			res, err := client.Post("http://127.0.0.1:1234/event/list/"+consts.TypeListEventByMonth, "application/json", r)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			var actualRes []models.Event
			json.Unmarshal(body, &actualRes)

			require.ElementsMatch(t, tt.exp, actualRes)
		})
	}
}

func TestCalendarAPI_FindEventByID(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := app2.New(memorystorage.New())

	timeDate := time.Date(2022, 5, 2, 1, 0, 0, 0, time.UTC)
	id, err := app.CreateEvent(ctx, "title", timeDate, 1)
	require.NoError(t, err)

	mockLogger := logger.NewMockLogger(ctrl)
	router := mux.NewRouter()
	server := NewServer(router, "127.0.0.1", "1234", mockLogger)
	calendarAPI := CalendarApp{app}
	calendarAPI.RegisterHTTPHandlers(router)
	defer server.Stop(ctx)

	go server.Start(ctx)

	// Слип, иначе тест завершается быстрее, чем поднимается HTTP сервер
	time.Sleep(1 * time.Millisecond)

	client := http.DefaultClient

	tests := []struct {
		name string
		in   string
		exp  models.Event
	}{
		{
			name: "Событие существует",
			in:   fmt.Sprintf(`{"id":"%s"}`, id.String()),
			exp: models.Event{
				ID:    id,
				Title: "title",
				Date:  timeDate,
				User:  1,
			},
		},
		{
			name: "Событие отсутствует",
			in:   `{"id":"1dc565af-026b-403f-b1f5-27bee27a302a"}`,
			exp:  models.Event{},
		},
		{
			name: "ИД события не вылидно",
			in:   `{"id":"not-valid-uuid"}`,
			exp:  models.Event{},
		},
		{
			name: "Пустое сообщение",
			in:   ``,
			exp:  models.Event{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.EXPECT().Info(gomock.Any())
			r := bytes.NewReader([]byte(tt.in))
			res, err := client.Post("http://127.0.0.1:1234/event/find", "application/json", r)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			var actualRes models.Event
			json.Unmarshal(body, &actualRes)

			require.EqualValues(t, tt.exp, actualRes)
		})
	}
}
