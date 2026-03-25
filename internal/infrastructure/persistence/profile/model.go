package profilepersistence

import (
	"time"

	"gorm.io/gorm"
)

// ProfileWriteModel contains only business/domain fields used for INSERT and UPDATE.
// Timestamp columns are intentionally absent so MySQL manages them via defaults.
type ProfileWriteModel struct {
	ID             string  `gorm:"primaryKey"`
	UserID         string  `gorm:"type:char(36);uniqueIndex;not null"`
	Gender         string  `gorm:"not null"`
	BirthYear      int     `gorm:"not null"`
	HeightCm       float64 `gorm:"not null"`
	WeightKg       float64 `gorm:"not null"`
	ActivityLevel  string  `gorm:"not null"`
	GoalKgPerWeek  float64 `gorm:"not null;default:0"`
	TDEE           float64 `gorm:"column:tdee;not null;default:0"`
	TargetCalories float64 `gorm:"not null;default:0"`
	SafetyFloor    float64 `gorm:"not null;default:1200"`
	Onboarded      bool    `gorm:"not null;default:false"`
}

func (ProfileWriteModel) TableName() string { return "profiles" }

// ProfileModel is the read model. It embeds ProfileWriteModel and adds timestamp
// fields populated by GORM when reading rows from the database.
type ProfileModel struct {
	ProfileWriteModel
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
