package weightpersistence

import (
	"time"

	"gorm.io/gorm"
)

// WeightEntryWriteModel contains only business/domain fields used for INSERT and UPDATE.
// Timestamp columns are intentionally absent so MySQL manages them via defaults.
type WeightEntryWriteModel struct {
	ID         string    `gorm:"primaryKey"`
	UserID     string    `gorm:"type:char(36);index;not null"`
	WeightKg   float64   `gorm:"not null"`
	MeasuredAt time.Time `gorm:"not null;type:date"`
}

func (WeightEntryWriteModel) TableName() string { return "weight_entries" }

// WeightEntryModel is the read model. It embeds WeightEntryWriteModel and adds
// timestamp fields populated by GORM when reading rows from the database.
type WeightEntryModel struct {
	WeightEntryWriteModel
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
