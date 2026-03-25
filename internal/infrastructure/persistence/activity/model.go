package activitypersistence

import (
	"time"

	"gorm.io/gorm"
)

// ActivityWriteModel contains only business/domain fields used for INSERT and UPDATE.
// Timestamp columns are intentionally absent so MySQL manages them via defaults.
type ActivityWriteModel struct {
	ID             string    `gorm:"primaryKey"`
	UserID         string    `gorm:"type:char(36);index;not null"`
	ActivityType   string    `gorm:"not null"`
	DurationMin    int       `gorm:"not null"`
	Intensity      string    `gorm:"not null"`
	CaloriesBurned float64   `gorm:"not null;default:0"`
	LoggedAt       time.Time `gorm:"not null"`
}

func (ActivityWriteModel) TableName() string { return "activities" }

// ActivityModel is the read model. It embeds ActivityWriteModel and adds timestamp
// fields populated by GORM when reading rows from the database.
type ActivityModel struct {
	ActivityWriteModel
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
