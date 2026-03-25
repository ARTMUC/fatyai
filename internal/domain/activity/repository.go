package activity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the persistence contract for the Activity aggregate.
type Repository interface {
	Save(ctx context.Context, a *Activity) error
	FindByID(ctx context.Context, id uuid.UUID) (*Activity, error)
	FindByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*Activity, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
