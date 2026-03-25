package weightapplication

import "time"

// LogWeightRequest carries data for a new morning weight entry.
type LogWeightRequest struct {
	UserID     string
	WeightKg   float64
	MeasuredAt time.Time
}

// WeightEntryDTO is the read-model for a single weight measurement.
type WeightEntryDTO struct {
	ID         string
	WeightKg   float64
	MeasuredAt time.Time
}

// WeightHistoryDTO aggregates historical weight data with trend analysis.
type WeightHistoryDTO struct {
	Entries      []WeightEntryDTO
	TrendSlope   float64  // kg/day (negative = losing weight)
	ForecastKg   float64  // weight at ForecastDate
	ForecastDays int      // days until goal weight at current pace
	ForecastDate time.Time
	GoalWeightKg float64
}

// CorrectionDTO describes a potential calorie correction (F5).
type CorrectionDTO struct {
	ShouldAdjust   bool
	SuggestedDelta float64
	Reason         string
	CurrentTarget  float64
	ProposedTarget float64
}
