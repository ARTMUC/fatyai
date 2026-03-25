package profileapplication

import (
	"context"
	"fmt"

	"github.com/artmuc/fatyai/internal/domain/profile"
	"github.com/artmuc/fatyai/internal/eventbus"
)

// Service is the application service for the Profile bounded context.
type Service struct {
	repo profile.Repository
	bus  *eventbus.Bus
}

// NewService creates a new profile application service.
func NewService(repo profile.Repository, bus *eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// Onboard creates (or updates) a user profile from the onboarding form.
func (s *Service) Onboard(ctx context.Context, req OnboardRequest) (ProfileDTO, error) {
	// Check for existing profile.
	existing, err := s.repo.FindByUserID(ctx, req.UserID)
	if err == nil {
		// Profile already exists — update measurements, activity level and goal.
		if err2 := existing.UpdateMeasurements(req.WeightKg, req.HeightCm); err2 != nil {
			return ProfileDTO{}, fmt.Errorf("update measurements: %w", err2)
		}
		if req.ActivityLevel != "" {
			if err2 := existing.ChangeActivityLevel(profile.ActivityLevel(req.ActivityLevel)); err2 != nil {
				return ProfileDTO{}, fmt.Errorf("change activity level: %w", err2)
			}
		}
		if err2 := existing.ChangeGoal(req.GoalKgPerWeek); err2 != nil {
			return ProfileDTO{}, fmt.Errorf("change goal: %w", err2)
		}
		if err2 := s.repo.Save(ctx, existing); err2 != nil {
			return ProfileDTO{}, fmt.Errorf("persist profile: %w", err2)
		}
		eventbus.PublishAll(s.bus, existing.PullEvents())
		return toDTO(existing), nil
	}

	p, err := profile.NewProfile(
		req.UserID,
		profile.Gender(req.Gender),
		req.BirthYear,
		req.HeightCm,
		req.WeightKg,
		profile.ActivityLevel(req.ActivityLevel),
		req.GoalKgPerWeek,
	)
	if err != nil {
		return ProfileDTO{}, fmt.Errorf("create profile: %w", err)
	}

	if err := s.repo.Save(ctx, p); err != nil {
		return ProfileDTO{}, fmt.Errorf("persist profile: %w", err)
	}

	eventbus.PublishAll(s.bus, p.PullEvents())
	return toDTO(p), nil
}

// GetByUserID returns the profile for a given user, or an error if not found.
func (s *Service) GetByUserID(ctx context.Context, userID string) (ProfileDTO, error) {
	p, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return ProfileDTO{}, err
	}
	return toDTO(p), nil
}

// AdjustTarget applies a calorie correction delta (F5).
func (s *Service) AdjustTarget(ctx context.Context, userID string, deltaKcal float64) (ProfileDTO, error) {
	p, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return ProfileDTO{}, fmt.Errorf("load profile: %w", err)
	}
	if err := p.AdjustTargetCalories(deltaKcal); err != nil {
		return ProfileDTO{}, fmt.Errorf("adjust target: %w", err)
	}
	if err := s.repo.Save(ctx, p); err != nil {
		return ProfileDTO{}, fmt.Errorf("persist profile: %w", err)
	}
	eventbus.PublishAll(s.bus, p.PullEvents())
	return toDTO(p), nil
}

func toDTO(p *profile.Profile) ProfileDTO {
	return ProfileDTO{
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
		Onboarded:      p.Onboarded(),
	}
}
