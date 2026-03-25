package activitypersistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	activitydomain "github.com/artmuc/fatyai/internal/domain/activity"
	"github.com/artmuc/fatyai/internal/repository"
)

// ActivityRepository is the GORM-backed adapter for activity.Repository.
type ActivityRepository struct {
	repository.BaseRepo[ActivityModel]
	translator Translator
}

// NewActivityRepository creates a new ActivityRepository.
func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{
		BaseRepo:   repository.NewBaseRepo[ActivityModel](db),
		translator: Translator{},
	}
}

// Save persists the Activity using INSERT … ON CONFLICT DO UPDATE so that both
// new activities and edits work correctly. Using ActivityWriteModel avoids sending
// zero time.Time fields for created_at/updated_at to MySQL.
func (r *ActivityRepository) Save(ctx context.Context, a *activitydomain.Activity) error {
	wm := r.translator.ToModel(a)
	return r.DB().WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(wm).Error
}

func (r *ActivityRepository) FindByID(ctx context.Context, id uuid.UUID) (*activitydomain.Activity, error) {
	scope := r.DB().WithContext(ctx).Where("id = ?", id.String())
	m, err := r.FirstScoped(scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, activitydomain.ErrNotFound
		}
		return nil, fmt.Errorf("find activity by id: %w", err)
	}
	return r.translator.ToDomain(m)
}

func (r *ActivityRepository) FindByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*activitydomain.Activity, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	models, err := r.FindAllScoped(
		r.DB().WithContext(ctx).
			Where("user_id = ? AND logged_at >= ? AND logged_at < ?", userID, start, end).
			Order("logged_at ASC"),
	)
	if err != nil {
		return nil, fmt.Errorf("find activities by user and date: %w", err)
	}

	result := make([]*activitydomain.Activity, 0, len(models))
	for i := range models {
		a, err := r.translator.ToDomain(models[i])
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *ActivityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB().WithContext(ctx).Where("id = ?", id.String()).Delete(&ActivityModel{})
	if result.Error != nil {
		return fmt.Errorf("delete activity: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return activitydomain.ErrNotFound
	}
	return nil
}
