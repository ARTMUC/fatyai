package weight

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the persistence contract for the WeightEntry aggregate.
type Repository interface {
	Save(ctx context.Context, e *WeightEntry) error
	FindByUserAndDate(ctx context.Context, userID string, date time.Time) (*WeightEntry, error)
	FindByUserSince(ctx context.Context, userID string, since time.Time) ([]*WeightEntry, error)
	FindByID(ctx context.Context, id uuid.UUID) (*WeightEntry, error)
}
