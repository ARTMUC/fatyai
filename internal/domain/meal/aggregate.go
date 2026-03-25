package meal

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Meal is the aggregate root for a single eating event.
type Meal struct {
	id        uuid.UUID
	userID    string
	name      string
	calories  float64
	proteinG  float64
	fatG      float64
	carbsG    float64
	weightG   float64
	estimated bool
	source    MealSource
	eatenAt   time.Time

	events []DomainEvent
}

// NewMeal creates a new Meal and records a MealLogged event.
func NewMeal(
	userID string,
	name string,
	calories, proteinG, fatG, carbsG, weightG float64,
	estimated bool,
	source MealSource,
	eatenAt time.Time,
) (*Meal, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	if !source.IsValid() {
		return nil, errors.New("invalid meal source")
	}
	m := &Meal{
		id:        uuid.New(),
		userID:    userID,
		name:      name,
		calories:  calories,
		proteinG:  proteinG,
		fatG:      fatG,
		carbsG:    carbsG,
		weightG:   weightG,
		estimated: estimated,
		source:    source,
		eatenAt:   eatenAt,
	}
	m.record(MealLogged{
		MealID:       m.id,
		UserID:       userID,
		Name:         name,
		CaloriesKcal: calories,
		OccurredOn:   time.Now(),
	})
	return m, nil
}

// Reconstitute rebuilds a Meal from persistence without triggering events.
func Reconstitute(
	id uuid.UUID,
	userID string,
	name string,
	calories, proteinG, fatG, carbsG, weightG float64,
	estimated bool,
	source MealSource,
	eatenAt time.Time,
) *Meal {
	return &Meal{
		id:        id,
		userID:    userID,
		name:      name,
		calories:  calories,
		proteinG:  proteinG,
		fatG:      fatG,
		carbsG:    carbsG,
		weightG:   weightG,
		estimated: estimated,
		source:    source,
		eatenAt:   eatenAt,
	}
}

// Edit updates the meal's editable fields.
func (m *Meal) Edit(name string, calories, proteinG, fatG, carbsG, weightG float64, eatenAt time.Time) error {
	if name == "" {
		return ErrEmptyName
	}
	m.name = name
	m.calories = calories
	m.proteinG = proteinG
	m.fatG = fatG
	m.carbsG = carbsG
	m.weightG = weightG
	m.eatenAt = eatenAt
	m.estimated = false
	m.record(MealEdited{MealID: m.id, OccurredOn: time.Now()})
	return nil
}

// -----------------------------------------------------------------
// Getters
// -----------------------------------------------------------------

func (m *Meal) ID() uuid.UUID      { return m.id }
func (m *Meal) UserID() string     { return m.userID }
func (m *Meal) Name() string       { return m.name }
func (m *Meal) CaloriesKcal() float64 { return m.calories }
func (m *Meal) ProteinG() float64  { return m.proteinG }
func (m *Meal) FatG() float64      { return m.fatG }
func (m *Meal) CarbsG() float64    { return m.carbsG }
func (m *Meal) WeightG() float64   { return m.weightG }
func (m *Meal) Estimated() bool    { return m.estimated }
func (m *Meal) Source() MealSource { return m.source }
func (m *Meal) EatenAt() time.Time { return m.eatenAt }

// -----------------------------------------------------------------
// Events
// -----------------------------------------------------------------

func (m *Meal) PullEvents() []DomainEvent {
	events := m.events
	m.events = nil
	return events
}

func (m *Meal) record(e DomainEvent) {
	m.events = append(m.events, e)
}
