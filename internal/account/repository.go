package account

import (
	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/db"
	"github.com/ahmedkhaeld/banking-app/db/models"
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
