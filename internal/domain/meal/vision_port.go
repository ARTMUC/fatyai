package meal

import "context"

// VisionAnalyzer is the port that the application layer uses to call
// the AI food recognition service. The OpenAI adapter implements this.
type VisionAnalyzer interface {
	AnalyzeMeal(ctx context.Context, imageBase64, mimeType string) (*VisionResult, error)
}

// VisionResult holds the AI-extracted nutritional data for a meal photo.
type VisionResult struct {
	Name         string
	CaloriesKcal float64
	ProteinG     float64
	FatG         float64
	CarbsG       float64
	WeightG      float64
	Estimated    bool
}
