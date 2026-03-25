package mealapplication

import (
	"time"

	activityapplication "github.com/artmuc/fatyai/internal/application/activity"
)

// ScanMealRequest carries the base64 image for AI analysis (F1).
type ScanMealRequest struct {
	UserID      string
	ImageBase64 string
	MimeType    string
}

// LogMealRequest carries the confirmed (possibly edited) meal data.
type LogMealRequest struct {
	UserID       string
	Name         string
	CaloriesKcal float64
	ProteinG     float64
	FatG         float64
	CarbsG       float64
	WeightG      float64
	Estimated    bool
	Source       string
	EatenAt      time.Time
}

// EditMealRequest carries updated fields for an existing meal.
type EditMealRequest struct {
	Name         string
	CaloriesKcal float64
	ProteinG     float64
	FatG         float64
	CarbsG       float64
	WeightG      float64
	EatenAt      time.Time
}

// MealDTO is the read-model for a single meal.
type MealDTO struct {
	ID           string
	Name         string
	CaloriesKcal float64
	ProteinG     float64
	FatG         float64
	CarbsG       float64
	WeightG      float64
	Estimated    bool
	Source       string
	EatenAt      time.Time
}

// DailyJournalDTO aggregates all meals and activities for one day.
type DailyJournalDTO struct {
	Date           time.Time
	Meals          []MealDTO
	Activities     []activityapplication.ActivityDTO
	TotalCalories  float64
	TotalProteinG  float64
	TotalFatG      float64
	TotalCarbsG    float64
	CaloriesBurned float64
	NetCalories    float64
	TargetCalories float64
}
