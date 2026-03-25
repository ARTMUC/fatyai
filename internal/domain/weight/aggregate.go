package weight

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// WeightEntry records a single morning weight measurement.
type WeightEntry struct {
	id         uuid.UUID
	userID     string
	weightKg   float64
	measuredAt time.Time // normalised to UTC midnight (date only)

	events []DomainEvent
}

// NewWeightEntry creates a new entry and records a WeightLogged event.
func NewWeightEntry(userID string, weightKg float64, measuredAt time.Time) (*WeightEntry, error) {
	if weightKg < 20 || weightKg > 500 {
		return nil, errors.New("weight must be between 20 and 500 kg")
	}
	// Normalise to date only (UTC midnight).
	date := time.Date(measuredAt.Year(), measuredAt.Month(), measuredAt.Day(), 0, 0, 0, 0, time.UTC)
	e := &WeightEntry{
		id:         uuid.New(),
		userID:     userID,
		weightKg:   weightKg,
		measuredAt: date,
	}
	e.record(WeightLogged{EntryID: e.id, UserID: userID, WeightKg: weightKg, OccurredOn: time.Now()})
	return e, nil
}

// Reconstitute rebuilds a WeightEntry from persistence without triggering events.
func Reconstitute(id uuid.UUID, userID string, weightKg float64, measuredAt time.Time) *WeightEntry {
	return &WeightEntry{id: id, userID: userID, weightKg: weightKg, measuredAt: measuredAt}
}

// -----------------------------------------------------------------
// Getters
// -----------------------------------------------------------------

func (e *WeightEntry) ID() uuid.UUID       { return e.id }
func (e *WeightEntry) UserID() string      { return e.userID }
func (e *WeightEntry) WeightKg() float64   { return e.weightKg }
func (e *WeightEntry) MeasuredAt() time.Time { return e.measuredAt }

// -----------------------------------------------------------------
// Events
// -----------------------------------------------------------------

func (e *WeightEntry) PullEvents() []DomainEvent {
	events := e.events
	e.events = nil
	return events
}

func (e *WeightEntry) record(ev DomainEvent) {
	e.events = append(e.events, ev)
}
