package mealapplication

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	activityapplication "github.com/artmuc/fatyai/internal/application/activity"
	activitydomain "github.com/artmuc/fatyai/internal/domain/activity"
	mealdomain "github.com/artmuc/fatyai/internal/domain/meal"
	profiledomain "github.com/artmuc/fatyai/internal/domain/profile"
	"github.com/artmuc/fatyai/internal/eventbus"
)

// Service is the application service for the Meal bounded context.
type Service struct {
	repo        mealdomain.Repository
	activityRepo activitydomain.Repository
	profileRepo  profiledomain.Repository
	vision      mealdomain.VisionAnalyzer
	bus         *eventbus.Bus
}

// NewService creates a new meal application service.
func NewService(
	repo mealdomain.Repository,
	activityRepo activitydomain.Repository,
	profileRepo profiledomain.Repository,
	vision mealdomain.VisionAnalyzer,
	bus *eventbus.Bus,
) *Service {
	return &Service{
		repo:         repo,
		activityRepo: activityRepo,
		profileRepo:  profileRepo,
		vision:       vision,
		bus:          bus,
	}
}

// ScanMeal calls the Vision AI and returns a draft MealDTO (not persisted).
func (s *Service) ScanMeal(ctx context.Context, req ScanMealRequest) (MealDTO, error) {
	result, err := s.vision.AnalyzeMeal(ctx, req.ImageBase64, req.MimeType)
	if err != nil {
		return MealDTO{}, fmt.Errorf("analyze meal: %w", err)
	}
	return MealDTO{
		Name:         result.Name,
		CaloriesKcal: result.CaloriesKcal,
		ProteinG:     result.ProteinG,
		FatG:         result.FatG,
		CarbsG:       result.CarbsG,
		WeightG:      result.WeightG,
		Estimated:    result.Estimated,
		Source:       "scan",
		EatenAt:      time.Now(),
	}, nil
}

// LogMeal saves a confirmed meal and dispatches domain events.
func (s *Service) LogMeal(ctx context.Context, req LogMealRequest) (MealDTO, error) {
	source := mealdomain.MealSource(req.Source)
	if !source.IsValid() {
		source = mealdomain.MealSourceManual
	}
	eatenAt := req.EatenAt
	if eatenAt.IsZero() {
		eatenAt = time.Now()
	}

	m, err := mealdomain.NewMeal(
		req.UserID, req.Name,
		req.CaloriesKcal, req.ProteinG, req.FatG, req.CarbsG, req.WeightG,
		req.Estimated, source, eatenAt,
	)
	if err != nil {
		return MealDTO{}, fmt.Errorf("create meal: %w", err)
	}
	if err := s.repo.Save(ctx, m); err != nil {
		return MealDTO{}, fmt.Errorf("persist meal: %w", err)
	}
	eventbus.PublishAll(s.bus, m.PullEvents())
	return toDTO(m), nil
}

// GetDailyJournal returns all meals and activities for a given date.
func (s *Service) GetDailyJournal(ctx context.Context, userID string, date time.Time) (DailyJournalDTO, error) {
	meals, err := s.repo.FindByUserAndDate(ctx, userID, date)
	if err != nil {
		return DailyJournalDTO{}, fmt.Errorf("load meals: %w", err)
	}
	activities, err := s.activityRepo.FindByUserAndDate(ctx, userID, date)
	if err != nil {
		return DailyJournalDTO{}, fmt.Errorf("load activities: %w", err)
	}

	mealDTOs := make([]MealDTO, 0, len(meals))
	var totalCal, totalProt, totalFat, totalCarbs float64
	for _, m := range meals {
		mealDTOs = append(mealDTOs, toDTO(m))
		totalCal += m.CaloriesKcal()
		totalProt += m.ProteinG()
		totalFat += m.FatG()
		totalCarbs += m.CarbsG()
	}

	actDTOs := make([]activityapplication.ActivityDTO, 0, len(activities))
	var burned float64
	for _, a := range activities {
		actDTOs = append(actDTOs, activityapplication.ActivityDTO{
			ID:             a.ID().String(),
			ActivityType:   a.ActivityType(),
			DurationMin:    a.DurationMin(),
			Intensity:      string(a.Intensity()),
			CaloriesBurned: a.CaloriesBurned(),
			LoggedAt:       a.LoggedAt(),
		})
		burned += a.CaloriesBurned()
	}

	var targetCal float64
	if p, err := s.profileRepo.FindByUserID(ctx, userID); err == nil {
		targetCal = p.TargetCalories()
	}

	return DailyJournalDTO{
		Date:           date,
		Meals:          mealDTOs,
		Activities:     actDTOs,
		TotalCalories:  totalCal,
		TotalProteinG:  totalProt,
		TotalFatG:      totalFat,
		TotalCarbsG:    totalCarbs,
		CaloriesBurned: burned,
		NetCalories:    totalCal - burned,
		TargetCalories: targetCal,
	}, nil
}

// EditMeal updates a meal's editable fields.
func (s *Service) EditMeal(ctx context.Context, mealID string, req EditMealRequest) (MealDTO, error) {
	id, err := uuid.Parse(mealID)
	if err != nil {
		return MealDTO{}, fmt.Errorf("invalid meal id: %w", err)
	}
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return MealDTO{}, err
	}
	if err := m.Edit(req.Name, req.CaloriesKcal, req.ProteinG, req.FatG, req.CarbsG, req.WeightG, req.EatenAt); err != nil {
		return MealDTO{}, fmt.Errorf("edit meal: %w", err)
	}
	if err := s.repo.Save(ctx, m); err != nil {
		return MealDTO{}, fmt.Errorf("persist meal: %w", err)
	}
	eventbus.PublishAll(s.bus, m.PullEvents())
	return toDTO(m), nil
}

// DeleteMeal removes a meal by ID.
func (s *Service) DeleteMeal(ctx context.Context, mealID string) error {
	id, err := uuid.Parse(mealID)
	if err != nil {
		return fmt.Errorf("invalid meal id: %w", err)
	}
	return s.repo.Delete(ctx, id)
}

func toDTO(m *mealdomain.Meal) MealDTO {
	return MealDTO{
		ID:           m.ID().String(),
		Name:         m.Name(),
		CaloriesKcal: m.CaloriesKcal(),
		ProteinG:     m.ProteinG(),
		FatG:         m.FatG(),
		CarbsG:       m.CarbsG(),
		WeightG:      m.WeightG(),
		Estimated:    m.Estimated(),
		Source:       string(m.Source()),
		EatenAt:      m.EatenAt(),
	}
}
