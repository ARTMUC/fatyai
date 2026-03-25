package weightpersistence

import (
	"fmt"

	"github.com/google/uuid"

	weightdomain "github.com/artmuc/fatyai/internal/domain/weight"
)

// Translator converts between WeightEntry domain aggregate and persistence models.
type Translator struct{}

// ToModel maps a domain WeightEntry to a WeightEntryWriteModel (used for INSERT / UPDATE).
func (t Translator) ToModel(e *weightdomain.WeightEntry) *WeightEntryWriteModel {
	return &WeightEntryWriteModel{
		ID:         e.ID().String(),
		UserID:     e.UserID(),
		WeightKg:   e.WeightKg(),
		MeasuredAt: e.MeasuredAt(),
	}
}

// ToDomain maps a WeightEntryModel (read model) back to a domain WeightEntry.
func (t Translator) ToDomain(m *WeightEntryModel) (*weightdomain.WeightEntry, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("translator: invalid weight entry id %q: %w", m.ID, err)
	}
	return weightdomain.Reconstitute(id, m.UserID, m.WeightKg, m.MeasuredAt), nil
}
