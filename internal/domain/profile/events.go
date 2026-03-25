package profile

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the marker interface for profile domain events.
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type ProfileCreated struct {
	ProfileID  uuid.UUID
	UserID     string
	OccurredOn time.Time
}

func (e ProfileCreated) EventName() string    { return "profile.created" }
func (e ProfileCreated) OccurredAt() time.Time { return e.OccurredOn }

type ProfileUpdated struct {
	ProfileID  uuid.UUID
	OccurredOn time.Time
}

func (e ProfileUpdated) EventName() string    { return "profile.updated" }
func (e ProfileUpdated) OccurredAt() time.Time { return e.OccurredOn }

type TargetCaloriesAdjusted struct {
	ProfileID  uuid.UUID
	OldTarget  float64
	NewTarget  float64
	OccurredOn time.Time
}

func (e TargetCaloriesAdjusted) EventName() string    { return "profile.target_calories_adjusted" }
func (e TargetCaloriesAdjusted) OccurredAt() time.Time { return e.OccurredOn }
