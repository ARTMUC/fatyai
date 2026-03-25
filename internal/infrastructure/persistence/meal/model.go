package mealpersistence

import (
	"time"

	"gorm.io/gorm"
)

// MealWriteModel contains only business/domain fields used for INSERT and UPDATE.
// Timestamp columns are intentionally absent so MySQL manages them via defaults.
type MealWriteModel struct {
	ID           string  `gorm:"primaryKey"`
	UserID       string  `gorm:"type:char(36);index;not null"`
	Name         string  `gorm:"not null"`
	CaloriesKcal float64 `gorm:"not null;default:0"`
	ProteinG     float64 `gorm:"not null;default:0"`
	FatG         float64 `gorm:"not null;default:0"`
	CarbsG       float64 `gorm:"not null;default:0"`
	WeightG      float64 `gorm:"not null;default:0"`
	Estimated    bool    `gorm:"not null;default:false"`
	Source       string  `gorm:"not null;default:'manual'"`
	EatenAt      time.Time `gorm:"not null"`
}

func (MealWriteModel) TableName() string { return "meals" }

// MealModel is the read model. It embeds MealWriteModel and adds timestamp fields
// populated by GORM when reading rows from the database.
type MealModel struct {
	MealWriteModel
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
