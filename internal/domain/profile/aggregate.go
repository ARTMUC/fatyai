package profile

import (
	"time"

	"github.com/google/uuid"
)

// Profile is the aggregate root for user body metrics and caloric goals.
// All TDEE and target calorie calculations live here.
type Profile struct {
	id             uuid.UUID
	userID         string
	gender         Gender
	birthYear      int
	heightCm       float64
	weightKg       float64
	activityLevel  ActivityLevel
	goalKgPerWeek  float64
	tdee           float64
	targetCalories float64
	safetyFloor    float64
	onboarded      bool

	events []DomainEvent
}

// NewProfile creates a brand-new Profile and calculates initial TDEE/target.
func NewProfile(
	userID string,
	gender Gender,
	birthYear int,
	heightCm, weightKg float64,
	level ActivityLevel,
	goalKgPerWeek float64,
) (*Profile, error) {
	if err := validate(gender, birthYear, heightCm, weightKg, level, goalKgPerWeek); err != nil {
		return nil, err
	}
	p := &Profile{
		id:            uuid.New(),
		userID:        userID,
		gender:        gender,
		birthYear:     birthYear,
		heightCm:      heightCm,
		weightKg:      weightKg,
		activityLevel: level,
		goalKgPerWeek: goalKgPerWeek,
		safetyFloor:   safetyFloor(gender),
		onboarded:     true,
	}
	p.recalculate()
	p.record(ProfileCreated{ProfileID: p.id, UserID: p.userID, OccurredOn: time.Now()})
	return p, nil
}

// Reconstitute rebuilds a Profile from persistence without triggering events.
func Reconstitute(
	id uuid.UUID,
	userID string,
	gender Gender,
	birthYear int,
	heightCm, weightKg float64,
	level ActivityLevel,
	goalKgPerWeek, tdee, targetCalories, floor float64,
	onboarded bool,
) *Profile {
	return &Profile{
		id:             id,
		userID:         userID,
		gender:         gender,
		birthYear:      birthYear,
		heightCm:       heightCm,
		weightKg:       weightKg,
		activityLevel:  level,
		goalKgPerWeek:  goalKgPerWeek,
		tdee:           tdee,
		targetCalories: targetCalories,
		safetyFloor:    floor,
		onboarded:      onboarded,
	}
}

// UpdateMeasurements updates weight/height and recalculates targets.
func (p *Profile) UpdateMeasurements(weightKg, heightCm float64) error {
	if weightKg < 20 || weightKg > 500 {
		return ErrInvalidWeight
	}
	if heightCm < 50 || heightCm > 300 {
		return ErrInvalidHeight
	}
	p.weightKg = weightKg
	p.heightCm = heightCm
	p.recalculate()
	p.record(ProfileUpdated{ProfileID: p.id, OccurredOn: time.Now()})
	return nil
}

// ChangeActivityLevel updates the activity multiplier and recalculates TDEE/target.
func (p *Profile) ChangeActivityLevel(level ActivityLevel) error {
	if !level.IsValid() {
		return ErrInvalidActivityLevel
	}
	p.activityLevel = level
	p.recalculate()
	p.record(ProfileUpdated{ProfileID: p.id, OccurredOn: time.Now()})
	return nil
}

// ChangeGoal updates the target weight loss rate and recalculates.
func (p *Profile) ChangeGoal(goalKgPerWeek float64) error {
	if goalKgPerWeek < 0 || goalKgPerWeek > 1 {
		return ErrInvalidGoal
	}
	p.goalKgPerWeek = goalKgPerWeek
	p.recalculate()
	p.record(ProfileUpdated{ProfileID: p.id, OccurredOn: time.Now()})
	return nil
}

// AdjustTargetCalories applies a correction delta (F5), respecting the safety floor.
func (p *Profile) AdjustTargetCalories(deltaKcal float64) error {
	old := p.targetCalories
	proposed := p.targetCalories + deltaKcal
	if proposed < p.safetyFloor {
		proposed = p.safetyFloor
	}
	p.targetCalories = proposed
	p.record(TargetCaloriesAdjusted{
		ProfileID:  p.id,
		OldTarget:  old,
		NewTarget:  proposed,
		OccurredOn: time.Now(),
	})
	return nil
}

// -----------------------------------------------------------------
// Getters
// -----------------------------------------------------------------

func (p *Profile) ID() uuid.UUID            { return p.id }
func (p *Profile) UserID() string           { return p.userID }
func (p *Profile) Gender() Gender           { return p.gender }
func (p *Profile) BirthYear() int           { return p.birthYear }
func (p *Profile) HeightCm() float64        { return p.heightCm }
func (p *Profile) WeightKg() float64        { return p.weightKg }
func (p *Profile) ActivityLevel() ActivityLevel { return p.activityLevel }
func (p *Profile) GoalKgPerWeek() float64   { return p.goalKgPerWeek }
func (p *Profile) TDEE() float64            { return p.tdee }
func (p *Profile) TargetCalories() float64  { return p.targetCalories }
func (p *Profile) SafetyFloor() float64     { return p.safetyFloor }
func (p *Profile) Onboarded() bool          { return p.onboarded }

// -----------------------------------------------------------------
// Events
// -----------------------------------------------------------------

func (p *Profile) PullEvents() []DomainEvent {
	events := p.events
	p.events = nil
	return events
}

func (p *Profile) record(e DomainEvent) {
	p.events = append(p.events, e)
}

// -----------------------------------------------------------------
// Internal calculations
// -----------------------------------------------------------------

// recalculate computes BMR (Mifflin-St Jeor), TDEE, and target calories.
func (p *Profile) recalculate() {
	age := time.Now().Year() - p.birthYear
	var bmr float64
	if p.gender == GenderMale {
		bmr = 10*p.weightKg + 6.25*p.heightCm - 5*float64(age) + 5
	} else {
		bmr = 10*p.weightKg + 6.25*p.heightCm - 5*float64(age) - 161
	}

	p.tdee = bmr * p.activityLevel.Multiplier()

	// deficit = 7700 kcal per kg * goal kg/week / 7 days
	dailyDeficit := p.goalKgPerWeek * 7700 / 7
	target := p.tdee - dailyDeficit
	if target < p.safetyFloor {
		target = p.safetyFloor
	}
	p.targetCalories = target
}

func safetyFloor(g Gender) float64 {
	if g == GenderFemale {
		return 1200
	}
	return 1500
}

func validate(gender Gender, birthYear int, heightCm, weightKg float64, level ActivityLevel, goalKgPerWeek float64) error {
	if !gender.IsValid() {
		return ErrInvalidGender
	}
	age := time.Now().Year() - birthYear
	if age < 10 || age > 120 {
		return ErrInvalidAge
	}
	if heightCm < 50 || heightCm > 300 {
		return ErrInvalidHeight
	}
	if weightKg < 20 || weightKg > 500 {
		return ErrInvalidWeight
	}
	if !level.IsValid() {
		return ErrInvalidActivityLevel
	}
	if goalKgPerWeek < 0 || goalKgPerWeek > 1 {
		return ErrInvalidGoal
	}
	return nil
}
