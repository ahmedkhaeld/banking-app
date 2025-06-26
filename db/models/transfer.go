package models

import (
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	FromAccountID uuid.UUID `gorm:"type:uuid;index;not null" json:"from_account_id"`
	ToAccountID   uuid.UUID `gorm:"type:uuid;index;not null" json:"to_account_id"`
	Amount        int64     `gorm:"not null;comment:must be positive" json:"amount"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	FromAccount   *Account  `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
	ToAccount     *Account  `gorm:"foreignKey:ToAccountID" json:"to_account,omitempty"`
}

func (Transfer) TableName() string { return "transfers" }
