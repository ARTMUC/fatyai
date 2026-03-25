package activityapplication

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	activitydomain "github.com/artmuc/fatyai/internal/domain/activity"
	"github.com/artmuc/fatyai/internal/eventbus"
)

// Service is the application service for the Activity bounded context.
type Service struct {
	repo activitydomain.Repository
	bus  *eventbus.Bus
}

// NewService creates a new activity application service.
func NewService(repo activitydomain.Repository, bus *eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// LogActivity creates and persists a new activity entry.
func (s *Service) LogActivity(ctx context.Context, req LogActivityRequest) (ActivityDTO, error) {
	loggedAt := req.LoggedAt
	if loggedAt.IsZero() {
		loggedAt = time.Now()
	}

	a, err := activitydomain.NewActivity(
		req.UserID,
		req.ActivityType,
		req.DurationMin,
		activitydomain.Intensity(req.Intensity),
		req.CaloriesBurned,
		loggedAt,
	)
	if err != nil {
		return ActivityDTO{}, fmt.Errorf("create activity: %w", err)
	}
	if err := s.repo.Save(ctx, a); err != nil {
		return ActivityDTO{}, fmt.Errorf("persist activity: %w", err)
	}
	eventbus.PublishAll(s.bus, a.PullEvents())
	return toDTO(a), nil
}

// GetDailyBurnedCalories returns total calories burned on the given day.
func (s *Service) GetDailyBurnedCalories(ctx context.Context, userID string, date time.Time) (float64, error) {
	activities, err := s.repo.FindByUserAndDate(ctx, userID, date)
	if err != nil {
		return 0, fmt.Errorf("find activities: %w", err)
	}
	var total float64
	for _, a := range activities {
		total += a.CaloriesBurned()
	}
	return total, nil
}

// DeleteActivity removes an activity by ID.
func (s *Service) DeleteActivity(ctx context.Context, activityID string) error {
	id, err := uuid.Parse(activityID)
	if err != nil {
		return fmt.Errorf("invalid activity id: %w", err)
	}
	return s.repo.Delete(ctx, id)
}

func toDTO(a *activitydomain.Activity) ActivityDTO {
	return ActivityDTO{
		ID:             a.ID().String(),
		ActivityType:   a.ActivityType(),
		DurationMin:    a.DurationMin(),
		Intensity:      string(a.Intensity()),
		CaloriesBurned: a.CaloriesBurned(),
		LoggedAt:       a.LoggedAt(),
	}
}
