package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/gen/pb/calendarpb"
)

var ErrEventDateNotRequest = errors.New("параметр Date не передан")

func (s *Server) CreateEvent(ctx context.Context, req *calendarpb.CreateEventRequest) (*calendarpb.SuccessOrErrorResponse, error) {
	if req.Event.Date == nil {
		return &calendarpb.SuccessOrErrorResponse{
			Id:      "",
			Success: false,
			Error:   fmt.Sprintf("Ошибка создания события: %s", ErrEventDateNotRequest.Error()),
		}, ErrEventDateNotRequest
	}

	date := time.Date(
		int(req.Event.Date.Year),
		time.Month(int(req.Event.Date.Month)),
		int(req.Event.Date.Day),
		int(req.Event.Date.Hour),
		int(req.Event.Date.Minute),
		0,
		0,
		time.UTC,
	)

	eventID, err := s.app.CreateEvent(ctx, req.Event.Title, date, int(req.Event.UserId))
	if err != nil {
		return &calendarpb.SuccessOrErrorResponse{
			Id:      eventID.String(),
			Success: false,
			Error:   "Ошибка создания события",
		}, errors.Wrap(err, "Ошибка создания события")
	}

	return &calendarpb.SuccessOrErrorResponse{
		Id:      eventID.String(),
		Success: true,
		Error:   "",
	}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *calendarpb.UpdateEventRequest) (*calendarpb.SuccessOrErrorResponse, error) {
	date := time.Date(
		int(req.Event.Date.Year),
		time.Month(int(req.Event.Date.Month)),
		int(req.Event.Date.Day),
		int(req.Event.Date.Hour),
		int(req.Event.Date.Minute),
		0,
		0,
		time.UTC,
	)

	eventID, err := uuid.Parse(req.Event.Id)
	if err != nil {
		return &calendarpb.SuccessOrErrorResponse{
			Id:      eventID.String(),
			Success: false,
			Error:   "Ошибка обработки eventID",
		}, errors.Wrap(err, "Ошибка парсинга ИД")
	}

	err = s.app.UpdateEvent(ctx, eventID, req.Event.Title, date, int(req.Event.UserId))
	if err != nil {
		return &calendarpb.SuccessOrErrorResponse{
			Id:      eventID.String(),
			Success: false,
			Error:   "Ошибка обновления события",
		}, errors.Wrap(err, "Ошибка обновления события")
	}

	return &calendarpb.SuccessOrErrorResponse{
		Id:      eventID.String(),
		Success: true,
		Error:   "",
	}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *calendarpb.DeleteEventRequest) (*calendarpb.SuccessOrErrorResponse, error) {
	eventID, err := uuid.Parse(req.Id)
	if err != nil {
		return &calendarpb.SuccessOrErrorResponse{
			Success: false,
			Error:   "Ошибка обработки eventID",
		}, errors.Wrap(err, "Ошибка парсинга ИД")
	}

	err = s.app.DeleteEvent(ctx, eventID)
	if err != nil {
		return &calendarpb.SuccessOrErrorResponse{
			Success: false,
			Error:   "Ошибка удаления события",
		}, errors.Wrap(err, "Ошибка удаления события")
	}

	return &calendarpb.SuccessOrErrorResponse{
		Id:      eventID.String(),
		Success: true,
		Error:   "",
	}, nil
}

func (s *Server) ListEventsForDay(ctx context.Context, req *calendarpb.ListEventsForDayRequest) (*calendarpb.Events, error) {
	date := time.Date(
		int(req.Date.Year),
		time.Month(int(req.Date.Month)),
		int(req.Date.Day),
		int(req.Date.Hour),
		int(req.Date.Minute),
		0,
		0,
		time.UTC,
	)

	events, err := s.app.ListEventsForDay(ctx, date)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка выборки событий за день")
	}

	return convertEventsToGRPC(events), nil
}

func (s *Server) ListEventsForWeek(ctx context.Context, req *calendarpb.ListEventsForWeekRequest) (*calendarpb.Events, error) {
	date := time.Date(
		int(req.Date.Year),
		time.Month(int(req.Date.Month)),
		int(req.Date.Day),
		int(req.Date.Hour),
		int(req.Date.Minute),
		0,
		0,
		time.UTC,
	)

	events, err := s.app.ListEventsForWeek(ctx, date)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка выборки событий за неделю")
	}

	return convertEventsToGRPC(events), nil
}

func (s *Server) ListEventsForMonth(ctx context.Context, req *calendarpb.ListEventsForMonthRequest) (*calendarpb.Events, error) {
	date := time.Date(
		int(req.Date.Year),
		time.Month(int(req.Date.Month)),
		int(req.Date.Day),
		int(req.Date.Hour),
		int(req.Date.Minute),
		0,
		0,
		time.UTC,
	)

	events, err := s.app.ListEventsForMonth(ctx, date)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка выборки событий за месяц")
	}

	return convertEventsToGRPC(events), nil
}

func (s *Server) FindEventByID(ctx context.Context, req *calendarpb.FindEventByIDRequest) (event *calendarpb.Event, err error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка поиска события по ID")
	}

	evnt := s.app.FindEventByID(ctx, id)
	if evnt == nil {
		return nil, nil
	}

	return &calendarpb.Event{
		Id:    evnt.ID.String(),
		Title: evnt.Title,
		Date: &calendarpb.DateTime{
			Year:   int32(evnt.Date.Year()),
			Month:  int32(evnt.Date.Month()),
			Day:    int32(evnt.Date.Day()),
			Hour:   int32(evnt.Date.Hour()),
			Minute: int32(evnt.Date.Minute()),
		},
		UserId: int64(evnt.User),
	}, nil
}
