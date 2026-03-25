package mealpersistence

import (
	"fmt"

	"github.com/google/uuid"

	mealdomain "github.com/artmuc/fatyai/internal/domain/meal"
)

// Translator converts between the Meal domain aggregate and persistence models.
type Translator struct{}

// ToModel maps a domain Meal to a MealWriteModel (used for INSERT / UPDATE).
func (t Translator) ToModel(m *mealdomain.Meal) *MealWriteModel {
	return &MealWriteModel{
		ID:           m.ID().String(),
		UserID:       m.UserID(),
		Name:         m.Name(),
		CaloriesKcal: m.CaloriesKcal(),
		ProteinG:     m.ProteinG(),
		FatG:         m.FatG(),
		CarbsG:       m.CarbsG(),
		WeightG:      m.WeightG(),
		Estimated:    m.Estimated(),
		Source:       string(m.Source()),
		EatenAt:      m.EatenAt(),
	}
}

// ToDomain maps a MealModel (read model) back to a domain Meal.
func (t Translator) ToDomain(m *MealModel) (*mealdomain.Meal, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("translator: invalid meal id %q: %w", m.ID, err)
	}
	return mealdomain.Reconstitute(
		id,
		m.UserID,
		m.Name,
		m.CaloriesKcal,
		m.ProteinG,
		m.FatG,
		m.CarbsG,
		m.WeightG,
		m.Estimated,
		mealdomain.MealSource(m.Source),
		m.EatenAt,
	), nil
}
