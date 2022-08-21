package grpc

import (
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/gen/pb/calendarpb"
	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/models"
)

func convertEventsToGRPC(events []models.Event) *calendarpb.Events {
	grpcEvents := []*calendarpb.Event{}

	for _, event := range events {
		grpcEvents = append(grpcEvents, &calendarpb.Event{
			Id:    event.ID.String(),
			Title: event.Title,
			Date: &calendarpb.DateTime{
				Year:   int32(event.Date.Year()),
				Month:  int32(event.Date.Month()),
				Day:    int32(event.Date.Day()),
				Hour:   int32(event.Date.Hour()),
				Minute: int32(event.Date.Minute()),
			},
			UserId: int64(event.User),
		})
	}

	return &calendarpb.Events{Events: grpcEvents}
}
