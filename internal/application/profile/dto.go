package profileapplication

// OnboardRequest carries the user's initial profile data.
type OnboardRequest struct {
	UserID        string
	Gender        string
	BirthYear     int
	HeightCm      float64
	WeightKg      float64
	ActivityLevel string
	GoalKgPerWeek float64
}

// UpdateMeasurementsRequest carries updated body measurements.
type UpdateMeasurementsRequest struct {
	WeightKg float64
	HeightCm float64
}

// ChangeGoalRequest carries a new weight loss goal.
type ChangeGoalRequest struct {
	GoalKgPerWeek float64
}

// ProfileDTO is the read-model returned to handlers and views.
type ProfileDTO struct {
	ID             string
	UserID         string
	Gender         string
	BirthYear      int
	HeightCm       float64
	WeightKg       float64
	ActivityLevel  string
	GoalKgPerWeek  float64
	TDEE           float64
	TargetCalories float64
	Onboarded      bool
}
