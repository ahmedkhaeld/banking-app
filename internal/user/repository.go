package user

import (
	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/db"
	"github.com/ahmedkhaeld/banking-app/db/models"
)

type model = models.User

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

func (r *Repository) usernameExists(username string) (bool, error) {
	var count int64
	err := r.Repository.DB.Model(&model{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) emailExists(email string) (bool, error) {
	var count int64
	err := r.Repository.DB.Model(&model{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) getByUsername(username string) (*model, error) {
	var user model
	err := r.Repository.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
