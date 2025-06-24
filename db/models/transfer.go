package models

import (
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FromAccountID uuid.UUID `gorm:"type:uuid;index;not null"`
	ToAccountID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Amount        int64     `gorm:"not null;comment:must be positive"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	FromAccount   *Account  `gorm:"foreignKey:FromAccountID"`
	ToAccount     *Account  `gorm:"foreignKey:ToAccountID"`
}

func (Transfer) TableName() string { return "transfers" }
