package profile

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the persistence contract for the Profile aggregate.
type Repository interface {
	Save(ctx context.Context, p *Profile) error
	FindByID(ctx context.Context, id uuid.UUID) (*Profile, error)
	FindByUserID(ctx context.Context, userID string) (*Profile, error)
}
