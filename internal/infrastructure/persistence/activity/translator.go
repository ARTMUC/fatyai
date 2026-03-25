package activitypersistence

import (
	"fmt"

	"github.com/google/uuid"

	activitydomain "github.com/artmuc/fatyai/internal/domain/activity"
)

// Translator converts between Activity domain aggregate and persistence models.
type Translator struct{}

// ToModel maps a domain Activity to an ActivityWriteModel (used for INSERT / UPDATE).
func (t Translator) ToModel(a *activitydomain.Activity) *ActivityWriteModel {
	return &ActivityWriteModel{
		ID:             a.ID().String(),
		UserID:         a.UserID(),
		ActivityType:   a.ActivityType(),
		DurationMin:    a.DurationMin(),
		Intensity:      string(a.Intensity()),
		CaloriesBurned: a.CaloriesBurned(),
		LoggedAt:       a.LoggedAt(),
	}
}

// ToDomain maps an ActivityModel (read model) back to a domain Activity.
func (t Translator) ToDomain(m *ActivityModel) (*activitydomain.Activity, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("translator: invalid activity id %q: %w", m.ID, err)
	}
	return activitydomain.Reconstitute(
		id,
		m.UserID,
		m.ActivityType,
		m.DurationMin,
		activitydomain.Intensity(m.Intensity),
		m.CaloriesBurned,
		m.LoggedAt,
	), nil
}
