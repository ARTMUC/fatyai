package activity

import "errors"

// Intensity represents the effort level of a physical activity.
type Intensity string

const (
	IntensityLow    Intensity = "low"
	IntensityMedium Intensity = "medium"
	IntensityHigh   Intensity = "high"
)

func (i Intensity) IsValid() bool {
	return i == IntensityLow || i == IntensityMedium || i == IntensityHigh
}

var (
	ErrEmptyType      = errors.New("activity type cannot be empty")
	ErrInvalidDuration = errors.New("duration must be greater than 0")
	ErrInvalidIntensity = errors.New("invalid intensity: must be low, medium, or high")
)
