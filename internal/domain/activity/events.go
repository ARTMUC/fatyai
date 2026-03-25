package activity

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the marker interface for activity domain events.
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type ActivityLogged struct {
	ActivityID     uuid.UUID
	UserID         string
	CaloriesBurned float64
	OccurredOn     time.Time
}

func (e ActivityLogged) EventName() string     { return "activity.logged" }
func (e ActivityLogged) OccurredAt() time.Time { return e.OccurredOn }

type ActivityDeleted struct {
	ActivityID uuid.UUID
	OccurredOn time.Time
}

func (e ActivityDeleted) EventName() string     { return "activity.deleted" }
func (e ActivityDeleted) OccurredAt() time.Time { return e.OccurredOn }
