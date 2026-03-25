package activityapplication

import "time"

// LogActivityRequest carries data for a new activity entry.
type LogActivityRequest struct {
	UserID         string
	ActivityType   string
	DurationMin    int
	Intensity      string
	CaloriesBurned float64
	LoggedAt       time.Time
}

// ActivityDTO is the read-model for a single activity entry.
type ActivityDTO struct {
	ID             string
	ActivityType   string
	DurationMin    int
	Intensity      string
	CaloriesBurned float64
	LoggedAt       time.Time
}
