package profilepersistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	profiledomain "github.com/artmuc/fatyai/internal/domain/profile"
	"github.com/artmuc/fatyai/internal/repository"
)

// ProfileRepository is the GORM-backed adapter for profile.Repository.
type ProfileRepository struct {
	repository.BaseRepo[ProfileModel]
	translator Translator
}

// NewProfileRepository creates a new ProfileRepository.
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		BaseRepo:   repository.NewBaseRepo[ProfileModel](db),
		translator: Translator{},
	}
}

// Save persists the Profile using INSERT … ON CONFLICT DO UPDATE to handle
// the unique(user_id) constraint. Using ProfileWriteModel avoids sending
// zero time.Time fields for created_at/updated_at to MySQL.
func (r *ProfileRepository) Save(ctx context.Context, p *profiledomain.Profile) error {
	wm := r.translator.ToModel(p)
	return r.DB().WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(wm).Error
}

func (r *ProfileRepository) FindByID(ctx context.Context, id uuid.UUID) (*profiledomain.Profile, error) {
	scope := r.DB().WithContext(ctx).Where("id = ?", id.String())
	m, err := r.FirstScoped(scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, profiledomain.ErrNotFound
		}
		return nil, fmt.Errorf("find profile by id: %w", err)
	}
	return r.translator.ToDomain(m)
}

func (r *ProfileRepository) FindByUserID(ctx context.Context, userID string) (*profiledomain.Profile, error) {
	scope := r.DB().WithContext(ctx).Where("user_id = ?", userID)
	m, err := r.FirstScoped(scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, profiledomain.ErrNotFound
		}
		return nil, fmt.Errorf("find profile by user_id: %w", err)
	}
	return r.translator.ToDomain(m)
}
