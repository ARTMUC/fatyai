// Package correction implements the domain service for F5 — dynamic calorie correction.
// It is a pure function with no state or infrastructure dependencies.
package correction

import (
	"fmt"

	profiledomain "github.com/artmuc/fatyai/internal/domain/profile"
	weightdomain "github.com/artmuc/fatyai/internal/domain/weight"
)

// Result describes whether and how to adjust the calorie target.
type Result struct {
	ShouldAdjust bool
	DeltaKcal    float64 // negative value (reduction)
	Reason       string
}

// Evaluate compares the user's actual weight loss over the provided entries
// against the expected loss based on their profile goal.
//
// Requires at least 14 days of entries to produce a meaningful result.
// If actual loss is less than 50% of expected → recommends cutting 150 kcal/day.
func Evaluate(p *profiledomain.Profile, entries []*weightdomain.WeightEntry) Result {
	if len(entries) < 2 {
		return Result{}
	}

	// Oldest entry first, newest last.
	first := entries[0]
	last := entries[len(entries)-1]

	days := last.MeasuredAt().Sub(first.MeasuredAt()).Hours() / 24
	if days < 7 {
		return Result{}
	}

	expectedLossKg := p.GoalKgPerWeek() * days / 7
	actualLossKg := first.WeightKg() - last.WeightKg() // positive = lost weight

	// If actual loss is less than half the expected → suggest correction.
	if actualLossKg < expectedLossKg*0.5 {
		delta := -150.0
		proposed := p.TargetCalories() + delta
		if proposed < p.SafetyFloor() {
			proposed = p.SafetyFloor()
			delta = proposed - p.TargetCalories()
		}
		return Result{
			ShouldAdjust: true,
			DeltaKcal:    delta,
			Reason: fmt.Sprintf(
				"Your progress (%.1f kg lost) is below target (%.1f kg). Adjusting target to %.0f kcal.",
				actualLossKg, expectedLossKg, proposed,
			),
		}
	}

	return Result{}
}
