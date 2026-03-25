package activity

import (
	"time"

	"github.com/google/uuid"
)

// Activity represents a single physical activity session.
type Activity struct {
	id             uuid.UUID
	userID         string
	activityType   string
	durationMin    int
	intensity      Intensity
	caloriesBurned float64
	loggedAt       time.Time

	events []DomainEvent
}

// NewActivity creates a new Activity and records an ActivityLogged event.
func NewActivity(
	userID string,
	activityType string,
	durationMin int,
	intensity Intensity,
	caloriesBurned float64,
	loggedAt time.Time,
) (*Activity, error) {
	if activityType == "" {
		return nil, ErrEmptyType
	}
	if durationMin <= 0 {
		return nil, ErrInvalidDuration
	}
	if !intensity.IsValid() {
		return nil, ErrInvalidIntensity
	}
	a := &Activity{
		id:             uuid.New(),
		userID:         userID,
		activityType:   activityType,
		durationMin:    durationMin,
		intensity:      intensity,
		caloriesBurned: caloriesBurned,
		loggedAt:       loggedAt,
	}
	a.record(ActivityLogged{
		ActivityID:     a.id,
		UserID:         userID,
		CaloriesBurned: caloriesBurned,
		OccurredOn:     time.Now(),
	})
	return a, nil
}

// Reconstitute rebuilds an Activity from persistence without triggering events.
func Reconstitute(
	id uuid.UUID,
	userID string,
	activityType string,
	durationMin int,
	intensity Intensity,
	caloriesBurned float64,
	loggedAt time.Time,
) *Activity {
	return &Activity{
		id:             id,
		userID:         userID,
		activityType:   activityType,
		durationMin:    durationMin,
		intensity:      intensity,
		caloriesBurned: caloriesBurned,
		loggedAt:       loggedAt,
	}
}

// -----------------------------------------------------------------
// Getters
// -----------------------------------------------------------------

func (a *Activity) ID() uuid.UUID           { return a.id }
func (a *Activity) UserID() string          { return a.userID }
func (a *Activity) ActivityType() string    { return a.activityType }
func (a *Activity) DurationMin() int        { return a.durationMin }
func (a *Activity) Intensity() Intensity    { return a.intensity }
func (a *Activity) CaloriesBurned() float64 { return a.caloriesBurned }
func (a *Activity) LoggedAt() time.Time     { return a.loggedAt }

// -----------------------------------------------------------------
// Events
// -----------------------------------------------------------------

func (a *Activity) PullEvents() []DomainEvent {
	events := a.events
	a.events = nil
	return events
}

func (a *Activity) record(e DomainEvent) {
	a.events = append(a.events, e)
}
