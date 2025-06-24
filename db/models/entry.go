package models

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AccountID uuid.UUID `gorm:"type:uuid;index;not null"`
	Amount    int64     `gorm:"not null;comment:can be negative or positive"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	Account   *Account  `gorm:"foreignKey:AccountID"`
}

func (Entry) TableName() string { return "entries" }
