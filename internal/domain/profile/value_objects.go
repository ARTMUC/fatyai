package profile

import "errors"

// Gender represents biological sex used in TDEE calculations.
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

func (g Gender) IsValid() bool {
	return g == GenderMale || g == GenderFemale
}

// ActivityLevel represents the user's baseline physical activity multiplier.
type ActivityLevel string

const (
	ActivitySedentary  ActivityLevel = "sedentary"   // x1.2
	ActivityLight      ActivityLevel = "light"        // x1.375
	ActivityModerate   ActivityLevel = "moderate"     // x1.55
	ActivityActive     ActivityLevel = "active"       // x1.725
	ActivityVeryActive ActivityLevel = "very_active"  // x1.9
)

func (a ActivityLevel) IsValid() bool {
	switch a {
	case ActivitySedentary, ActivityLight, ActivityModerate, ActivityActive, ActivityVeryActive:
		return true
	}
	return false
}

func (a ActivityLevel) Multiplier() float64 {
	switch a {
	case ActivitySedentary:
		return 1.2
	case ActivityLight:
		return 1.375
	case ActivityModerate:
		return 1.55
	case ActivityActive:
		return 1.725
	case ActivityVeryActive:
		return 1.9
	}
	return 1.2
}

var (
	ErrInvalidGender        = errors.New("invalid gender: must be 'male' or 'female'")
	ErrInvalidActivityLevel = errors.New("invalid activity level")
	ErrInvalidGoal          = errors.New("goal must be between 0 and 1 kg/week")
	ErrInvalidAge           = errors.New("birth year must result in age between 10 and 120")
	ErrInvalidHeight        = errors.New("height must be between 50 and 300 cm")
	ErrInvalidWeight        = errors.New("weight must be between 20 and 500 kg")
)
