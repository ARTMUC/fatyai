package profilepersistence

import (
	"fmt"

	"github.com/google/uuid"

	profiledomain "github.com/artmuc/fatyai/internal/domain/profile"
)

// Translator converts between the Profile domain aggregate and persistence models.
type Translator struct{}

// ToModel maps a domain Profile to a ProfileWriteModel (used for INSERT / UPDATE).
func (t Translator) ToModel(p *profiledomain.Profile) *ProfileWriteModel {
	return &ProfileWriteModel{
		ID:             p.ID().String(),
		UserID:         p.UserID(),
		Gender:         string(p.Gender()),
		BirthYear:      p.BirthYear(),
		HeightCm:       p.HeightCm(),
		WeightKg:       p.WeightKg(),
		ActivityLevel:  string(p.ActivityLevel()),
		GoalKgPerWeek:  p.GoalKgPerWeek(),
		TDEE:           p.TDEE(),
		TargetCalories: p.TargetCalories(),
		SafetyFloor:    p.SafetyFloor(),
		Onboarded:      p.Onboarded(),
	}
}

// ToDomain maps a ProfileModel (read model) back to a domain Profile.
func (t Translator) ToDomain(m *ProfileModel) (*profiledomain.Profile, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("translator: invalid profile id %q: %w", m.ID, err)
	}
	return profiledomain.Reconstitute(
		id,
		m.UserID,
		profiledomain.Gender(m.Gender),
		m.BirthYear,
		m.HeightCm,
		m.WeightKg,
		profiledomain.ActivityLevel(m.ActivityLevel),
		m.GoalKgPerWeek,
		m.TDEE,
		m.TargetCalories,
		m.SafetyFloor,
		m.Onboarded,
	), nil
}
