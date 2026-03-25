package meal

// MealSource indicates whether the meal was logged via photo scan or manually.
type MealSource string

const (
	MealSourceScan   MealSource = "scan"
	MealSourceManual MealSource = "manual"
)

func (s MealSource) IsValid() bool {
	return s == MealSourceScan || s == MealSourceManual
}
