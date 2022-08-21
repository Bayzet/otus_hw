//nolint:varnamelen
package internalhttp

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"
	models2 "github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/server/http/models"
)

type CalendarAPIv1 interface {
	RegisterHTTPHandlers(*mux.Router)
	CreateEventHandler(http.ResponseWriter, *http.Request)
	UpdateEventHandler(http.ResponseWriter, *http.Request)
	DeleteEventHandler(http.ResponseWriter, *http.Request)
	ListEventsHandler(http.ResponseWriter, *http.Request)
	FindEventByIDHandler(http.ResponseWriter, *http.Request)
}

type CalendarApp struct {
	App app.Application
}

func (c *CalendarApp) RegisterHTTPHandlers(r *mux.Router) {
	r.
		Methods(http.MethodPost).
		Path("/event/create").
		HandlerFunc(c.CreateEventHandler)
	r.
		Methods(http.MethodPost).
		Path("/event/update").
		HandlerFunc(c.UpdateEventHandler)
	r.
		Methods(http.MethodPost).
		Path("/event/delete").
		HandlerFunc(c.DeleteEventHandler)
	r.
		Methods(http.MethodPost).
		Path("/event/list/{type}").
		HandlerFunc(c.ListEventsHandler)
	r.
		Methods(http.MethodPost).
		Path("/event/find").
		HandlerFunc(c.FindEventByIDHandler)
}

func (c *CalendarApp) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	br := decodeJSONBody(r, &event)
	if br != nil {
		errorsBadRequest(w, br)

		return
	}

	err := event.Validate([]string{"date", "title", "user"})
	if err != nil {
		errorsBadRequest(w, &BadRequest{StatusCode: http.StatusBadRequest, Msg: err.Error()})

		return
	}

	id, err := c.App.CreateEvent(r.Context(), event.Title, event.Date, event.User)
	if err != nil {
		errorResponse(w, err)

		return
	}

	goodResponse(w, models.Response{ID: id})
}

func (c *CalendarApp) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	br := decodeJSONBody(r, &event)
	if br != nil {
		errorsBadRequest(w, br)

		return
	}

	err := event.Validate([]string{"id", "title", "date", "user"})
	if err != nil {
		errorsBadRequest(w, &BadRequest{StatusCode: http.StatusBadRequest, Msg: err.Error()})

		return
	}

	err = c.App.UpdateEvent(r.Context(), event.ID, event.Title, event.Date, event.User)
	if err != nil {
		errorResponse(w, err)

		return
	}

	goodResponse(w, models.Response{ID: event.ID})
}

func (c *CalendarApp) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	br := decodeJSONBody(r, &event)
	if br != nil {
		errorsBadRequest(w, br)

		return
	}

	err := event.Validate([]string{"id"})
	if err != nil {
		errorsBadRequest(w, &BadRequest{StatusCode: http.StatusBadRequest, Msg: err.Error()})

		return
	}

	err = c.App.DeleteEvent(r.Context(), event.ID)
	if err != nil {
		errorResponse(w, err)

		return
	}

	goodResponse(w, models.Response{ID: event.ID})
}

func (c *CalendarApp) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	br := decodeJSONBody(r, &event)
	if br != nil {
		errorsBadRequest(w, br)

		return
	}

	err := event.Validate([]string{"date"})
	if err != nil {
		errorsBadRequest(w, &BadRequest{StatusCode: http.StatusBadRequest, Msg: err.Error()})

		return
	}

	var events []models2.Event

	vars := mux.Vars(r)
	switch vars["type"] {
	case consts.TypeListEventByDay:
		events, err = c.App.ListEventsForDay(r.Context(), event.Date)
	case consts.TypeListEventByWeek:
		events, err = c.App.ListEventsForWeek(r.Context(), event.Date)
	case consts.TypeListEventByMonth:
		events, err = c.App.ListEventsForMonth(r.Context(), event.Date)
	default:
		errorResponse(w, errors.New("не существующий тип"))

		return
	}

	if err != nil {
		errorResponse(w, err)

		return
	}

	goodResponse(w, events)
}

func (c *CalendarApp) FindEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	br := decodeJSONBody(r, &event)
	if br != nil {
		errorsBadRequest(w, br)

		return
	}

	e := c.App.FindEventByID(r.Context(), event.ID)

	goodResponse(w, e)
}
