package storage

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID    uuid.UUID
	Title string
	Date  time.Time
	User  int
}
