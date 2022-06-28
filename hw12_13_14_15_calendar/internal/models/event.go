package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID    uuid.UUID
	Title string
	Date  time.Time
	User  int
}
