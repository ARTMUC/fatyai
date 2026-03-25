package meal

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the persistence contract for the Meal aggregate.
type Repository interface {
	Save(ctx context.Context, m *Meal) error
	FindByID(ctx context.Context, id uuid.UUID) (*Meal, error)
	FindByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*Meal, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
