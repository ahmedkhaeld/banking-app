package user

import (
	"errors"
	"fmt"

	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/db/models"
	"github.com/ahmedkhaeld/banking-app/internal/auth"
)

// Errors
var (
	ErrUsernameExists            = errors.New("username already exists")
	ErrEmailExists               = errors.New("email already exists")
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
)

type Service struct {
	crud.Service[model]
	repo *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		Service: *crud.NewService(repository),
		repo:    repository,
	}
}

func InitService() *Service {
	return &Service{
		repo:    InitRepository(),
		Service: *crud.NewService(InitRepository()),
	}
}

func (s *Service) CreateUser(req *CreateUserRequest) (*models.User, error) {
	exists, err := s.repo.usernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// Validate email uniqueness
	emailExists, err := s.repo.emailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrEmailExists
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userModel := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		FullName: req.FullName,
		Email:    req.Email,
	}
	err = s.Service.Create(userModel)
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

func (s *Service) LoginUser(username, password string) (*models.User, error) {
	user, err := s.repo.getByUsername(username)
	if err != nil {
		return nil, ErrInvalidUsernameOrPassword
	}
	if err := auth.CheckPassword(password, user.Password); err != nil {
		return nil, ErrInvalidUsernameOrPassword
	}
	return user, nil
}

func (s *Service) FindOneByID(id string) (*models.User, error) {
	var user models.User
	api := crud.GetAllRequest{
		Filter: []string{fmt.Sprintf("id||eq||%s", id)},
	}
	err := s.FindOne(api, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
