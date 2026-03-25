package weightapplication

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/artmuc/fatyai/internal/domain/correction"
	profiledomain "github.com/artmuc/fatyai/internal/domain/profile"
	weightdomain "github.com/artmuc/fatyai/internal/domain/weight"
	"github.com/artmuc/fatyai/internal/eventbus"
)

// Service is the application service for weight tracking (F6) and correction (F5).
type Service struct {
	repo        weightdomain.Repository
	profileRepo profiledomain.Repository
	bus         *eventbus.Bus
}

// NewService creates a new weight application service.
func NewService(repo weightdomain.Repository, profileRepo profiledomain.Repository, bus *eventbus.Bus) *Service {
	return &Service{repo: repo, profileRepo: profileRepo, bus: bus}
}

// LogWeight records a morning weight entry (upserts for today).
func (s *Service) LogWeight(ctx context.Context, req LogWeightRequest) (WeightEntryDTO, error) {
	measuredAt := req.MeasuredAt
	if measuredAt.IsZero() {
		measuredAt = time.Now()
	}

	e, err := weightdomain.NewWeightEntry(req.UserID, req.WeightKg, measuredAt)
	if err != nil {
		return WeightEntryDTO{}, fmt.Errorf("create weight entry: %w", err)
	}
	if err := s.repo.Save(ctx, e); err != nil {
		return WeightEntryDTO{}, fmt.Errorf("persist weight entry: %w", err)
	}
	eventbus.PublishAll(s.bus, e.PullEvents())
	return WeightEntryDTO{ID: e.ID().String(), WeightKg: e.WeightKg(), MeasuredAt: e.MeasuredAt()}, nil
}

// GetHistory returns the last `days` weight entries with trend analysis.
func (s *Service) GetHistory(ctx context.Context, userID string, days int) (WeightHistoryDTO, error) {
	since := time.Now().AddDate(0, 0, -days)
	entries, err := s.repo.FindByUserSince(ctx, userID, since)
	if err != nil {
		return WeightHistoryDTO{}, fmt.Errorf("load weight history: %w", err)
	}

	dtos := make([]WeightEntryDTO, 0, len(entries))
	for _, e := range entries {
		dtos = append(dtos, WeightEntryDTO{ID: e.ID().String(), WeightKg: e.WeightKg(), MeasuredAt: e.MeasuredAt()})
	}

	dto := WeightHistoryDTO{Entries: dtos}

	if len(entries) >= 2 {
		slope := linearRegressionSlope(entries)
		dto.TrendSlope = slope

		// Forecast: goal weight from profile
		if p, err := s.profileRepo.FindByUserID(ctx, userID); err == nil {
			goalWeight := p.WeightKg() - float64(p.GoalKgPerWeek())*float64(days)/7
			dto.GoalWeightKg = goalWeight
			if slope < 0 {
				latestWeight := entries[len(entries)-1].WeightKg()
				daysToGoal := (latestWeight - goalWeight) / (-slope)
				dto.ForecastDays = int(math.Round(daysToGoal))
				dto.ForecastDate = time.Now().AddDate(0, 0, dto.ForecastDays)
			}
		}
	}

	return dto, nil
}

// CheckCorrection evaluates whether a calorie adjustment is needed (F5).
func (s *Service) CheckCorrection(ctx context.Context, userID string) (CorrectionDTO, error) {
	p, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return CorrectionDTO{}, nil // no profile yet — no correction
	}

	since := time.Now().AddDate(0, 0, -14)
	entries, err := s.repo.FindByUserSince(ctx, userID, since)
	if err != nil {
		return CorrectionDTO{}, fmt.Errorf("load weight entries: %w", err)
	}

	result := correction.Evaluate(p, entries)
	if !result.ShouldAdjust {
		return CorrectionDTO{}, nil
	}

	return CorrectionDTO{
		ShouldAdjust:   true,
		SuggestedDelta: result.DeltaKcal,
		Reason:         result.Reason,
		CurrentTarget:  p.TargetCalories(),
		ProposedTarget: p.TargetCalories() + result.DeltaKcal,
	}, nil
}

// ApplyCorrection persists the suggested calorie adjustment to the profile.
func (s *Service) ApplyCorrection(ctx context.Context, userID string, deltaKcal float64) error {
	p, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("load profile: %w", err)
	}
	if err := p.AdjustTargetCalories(deltaKcal); err != nil {
		return fmt.Errorf("adjust target: %w", err)
	}
	return s.profileRepo.Save(ctx, p)
}

// -----------------------------------------------------------------
// Linear regression helper
// -----------------------------------------------------------------

// linearRegressionSlope returns the slope (kg/day) of the weight trend
// using ordinary least squares.
func linearRegressionSlope(entries []*weightdomain.WeightEntry) float64 {
	n := float64(len(entries))
	if n < 2 {
		return 0
	}
	base := entries[0].MeasuredAt()

	var sumX, sumY, sumXY, sumX2 float64
	for _, e := range entries {
		x := e.MeasuredAt().Sub(base).Hours() / 24 // day index
		y := e.WeightKg()
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denom := n*sumX2 - sumX*sumX
	if denom == 0 {
		return 0
	}
	return (n*sumXY - sumX*sumY) / denom
}
