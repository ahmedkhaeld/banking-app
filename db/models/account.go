package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Owner         string     `gorm:"index;not null"`
	Balance       int64      `gorm:"not null"`
	Currency      string     `gorm:"not null"`
	CreatedAt     time.Time  `gorm:"not null;autoCreateTime"`
	Entries       []Entry    `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TransfersFrom []Transfer `gorm:"foreignKey:FromAccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TransfersTo   []Transfer `gorm:"foreignKey:ToAccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Account) TableName() string { return "accounts" }
