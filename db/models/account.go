package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Balance   int64     `json:"balance" gorm:"type:bigint;default:0"`
	Owner     string    `json:"owner" gorm:"index;not null"`
	Currency  string    `json:"currency" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime"`
	// Relationships
	Entries       []Entry    `json:"entries,omitempty" gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TransfersFrom []Transfer `json:"transfers_from,omitempty" gorm:"foreignKey:FromAccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TransfersTo   []Transfer `json:"transfers_to,omitempty" gorm:"foreignKey:ToAccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Account) TableName() string {
	return "accounts"
}
