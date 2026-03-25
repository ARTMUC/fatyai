package weightpersistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	weightdomain "github.com/artmuc/fatyai/internal/domain/weight"
	"github.com/artmuc/fatyai/internal/repository"
)

// WeightRepository is the GORM-backed adapter for weight.Repository.
type WeightRepository struct {
	repository.BaseRepo[WeightEntryModel]
	translator Translator
}

// NewWeightRepository creates a new WeightRepository.
func NewWeightRepository(db *gorm.DB) *WeightRepository {
	return &WeightRepository{
		BaseRepo:   repository.NewBaseRepo[WeightEntryModel](db),
		translator: Translator{},
	}
}

// Save persists the WeightEntry using INSERT … ON CONFLICT DO UPDATE to handle
// the unique(user_id, measured_at) constraint. Using WeightEntryWriteModel avoids
// sending zero time.Time fields for created_at/updated_at to MySQL.
func (r *WeightRepository) Save(ctx context.Context, e *weightdomain.WeightEntry) error {
	wm := r.translator.ToModel(e)
	return r.DB().WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(wm).Error
}

func (r *WeightRepository) FindByID(ctx context.Context, id uuid.UUID) (*weightdomain.WeightEntry, error) {
	scope := r.DB().WithContext(ctx).Where("id = ?", id.String())
	m, err := r.FirstScoped(scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, weightdomain.ErrNotFound
		}
		return nil, fmt.Errorf("find weight entry by id: %w", err)
	}
	return r.translator.ToDomain(m)
}

func (r *WeightRepository) FindByUserAndDate(ctx context.Context, userID string, date time.Time) (*weightdomain.WeightEntry, error) {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	scope := r.DB().WithContext(ctx).Where("user_id = ? AND measured_at = ?", userID, d)
	m, err := r.FirstScoped(scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, weightdomain.ErrNotFound
		}
		return nil, fmt.Errorf("find weight entry by user and date: %w", err)
	}
	return r.translator.ToDomain(m)
}

func (r *WeightRepository) FindByUserSince(ctx context.Context, userID string, since time.Time) ([]*weightdomain.WeightEntry, error) {
	models, err := r.FindAllScoped(
		r.DB().WithContext(ctx).
			Where("user_id = ? AND measured_at >= ?", userID, since).
			Order("measured_at ASC"),
	)
	if err != nil {
		return nil, fmt.Errorf("find weight entries since: %w", err)
	}

	entries := make([]*weightdomain.WeightEntry, 0, len(models))
	for i := range models {
		e, err := r.translator.ToDomain(models[i])
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
