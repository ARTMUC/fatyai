package meal

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the marker interface for meal domain events.
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type MealLogged struct {
	MealID       uuid.UUID
	UserID       string
	Name         string
	CaloriesKcal float64
	OccurredOn   time.Time
}

func (e MealLogged) EventName() string     { return "meal.logged" }
func (e MealLogged) OccurredAt() time.Time { return e.OccurredOn }

type MealEdited struct {
	MealID     uuid.UUID
	OccurredOn time.Time
}

func (e MealEdited) EventName() string     { return "meal.edited" }
func (e MealEdited) OccurredAt() time.Time { return e.OccurredOn }

type MealDeleted struct {
	MealID     uuid.UUID
	UserID     string
	OccurredOn time.Time
}

func (e MealDeleted) EventName() string     { return "meal.deleted" }
func (e MealDeleted) OccurredAt() time.Time { return e.OccurredOn }
