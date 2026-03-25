package mealpersistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	mealdomain "github.com/artmuc/fatyai/internal/domain/meal"
	"github.com/artmuc/fatyai/internal/repository"
)

// MealRepository is the GORM-backed adapter for meal.Repository.
type MealRepository struct {
	repository.BaseRepo[MealModel]
	translator Translator
}

// NewMealRepository creates a new MealRepository.
func NewMealRepository(db *gorm.DB) *MealRepository {
	return &MealRepository{
		BaseRepo:   repository.NewBaseRepo[MealModel](db),
		translator: Translator{},
	}
}

// Save persists the Meal using INSERT … ON CONFLICT DO UPDATE so that both
// new meals and edits work correctly. Using MealWriteModel avoids sending
// zero time.Time fields for created_at/updated_at to MySQL.
func (r *MealRepository) Save(ctx context.Context, m *mealdomain.Meal) error {
	wm := r.translator.ToModel(m)
	return r.DB().WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(wm).Error
}

func (r *MealRepository) FindByID(ctx context.Context, id uuid.UUID) (*mealdomain.Meal, error) {
	scope := r.DB().WithContext(ctx).Where("id = ?", id.String())
	m, err := r.FirstScoped(scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, mealdomain.ErrNotFound
		}
		return nil, fmt.Errorf("find meal by id: %w", err)
	}
	return r.translator.ToDomain(m)
}

func (r *MealRepository) FindByUserAndDate(ctx context.Context, userID string, date time.Time) ([]*mealdomain.Meal, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).UTC()
	end := start.Add(24 * time.Hour)

	models, err := r.FindAllScoped(
		r.DB().WithContext(ctx).
			Where("user_id = ? AND eaten_at >= ? AND eaten_at < ?", userID, start, end).
			Order("eaten_at ASC"),
	)
	if err != nil {
		return nil, fmt.Errorf("find meals by user and date: %w", err)
	}

	meals := make([]*mealdomain.Meal, 0, len(models))
	for i := range models {
		m, err := r.translator.ToDomain(models[i])
		if err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}
	return meals, nil
}

func (r *MealRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB().WithContext(ctx).Where("id = ?", id.String()).Delete(&MealModel{})
	if result.Error != nil {
		return fmt.Errorf("delete meal: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return mealdomain.ErrNotFound
	}
	return nil
}
