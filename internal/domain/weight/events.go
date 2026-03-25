package weight

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the marker interface for weight domain events.
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type WeightLogged struct {
	EntryID    uuid.UUID
	UserID     string
	WeightKg   float64
	OccurredOn time.Time
}

func (e WeightLogged) EventName() string     { return "weight.logged" }
func (e WeightLogged) OccurredAt() time.Time { return e.OccurredOn }
