package account

import (
	"context"
	"errors"

	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/db/models"
	"github.com/ahmedkhaeld/banking-app/internal/user"
	"github.com/google/uuid"
)

type Service struct {
	crud.Service[model]
	repo        *Repository
	userService *user.Service
}

func NewService(repository *Repository) *Service {
	return &Service{
		Service: *crud.NewService(repository),
		repo:    repository,
	}
}

func InitService() *Service {
	return &Service{
		repo:        InitRepository(),
		Service:     *crud.NewService(InitRepository()),
		userService: user.InitService(),
	}
}

func (s *Service) createAccount(req CreateAccountRequest, userId string) (*CreateAccountResponse, error) {
	user, err := s.userService.FindOneByID(userId)
	if err != nil {
		return nil, errors.New("user does not exist")
	}
	userID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}
	account := &models.Account{
		UserID:   userID,
		Currency: req.Currency,
		Owner:    user.Username,
	}
	if req.Balance != nil {
		account.Balance = *req.Balance
	}
	err = s.repo.Repository.DB.Create(account).Error
	if err != nil {
		return nil, err
	}
	resp := &CreateAccountResponse{
		ID:        account.ID.String(),
		UserID:    account.UserID.String(),
		Owner:     account.Owner,
		Currency:  account.Currency,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	return resp, nil
}

func (s *Service) getAccountBalance(accountID, userID string) (*AccountBalanceResponse, error) {
	var account models.Account
	if err := s.repo.Repository.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		return nil, err
	}
	return &AccountBalanceResponse{
		ID:       account.ID.String(),
		Balance:  account.Balance,
		Currency: account.Currency,
	}, nil
}

func (s *Service) updateBalance(ctx context.Context, accountID string, amount int64) (*models.Account, error) {
	id, err := uuid.Parse(accountID)
	if err != nil {
		return nil, errors.New("invalid account_id format")
	}
	account, err := s.repo.updateBalance(ctx, id.String(), amount)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) isAccountOwnedByUser(accountID, userID string) bool {
	var account models.Account
	err := s.repo.Repository.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error
	return err == nil
}
