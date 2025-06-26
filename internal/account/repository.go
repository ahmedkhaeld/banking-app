package account

import (
	"context"

	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/db"
	"github.com/ahmedkhaeld/banking-app/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type model = models.Account

type Repository struct {
	crud.Repository[model]
}

func InitRepository() *Repository {
	return &Repository{
		Repository: crud.Repository[model]{
			DB:    db.DB,
			Model: model{},
		},
	}
}

// UpdateBalance updates the balance of an account by a given amount and returns the updated account
func (r *Repository) updateBalance(ctx context.Context, accountID string, amount int64) (*model, error) {
	id, err := uuid.Parse(accountID)
	if err != nil {
		return nil, err
	}
	db := r.Repository.DB.WithContext(ctx)
	if err := db.Model(&model{}).Where("id = ?", id).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		return nil, err
	}
	var account model
	if err := db.Where("id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
