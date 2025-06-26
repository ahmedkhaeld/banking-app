package transfer

import (
	"context"

	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/db/models"
	"github.com/google/uuid"
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

// Transfer performs a money transfer between accounts using a transaction
func (s *Service) Transfer(ctx context.Context, req CreateTransferRequest) (*CreateTransferResponse, error) {
	params := TransferTxParams(req)
	result, err := s.repo.TransferTx(ctx, params)
	if err != nil {
		return nil, err
	}
	resp := &CreateTransferResponse{
		ID:            result.Transfer.ID.String(),
		FromAccountID: result.Transfer.FromAccountID.String(),
		ToAccountID:   result.Transfer.ToAccountID.String(),
		Amount:        result.Transfer.Amount,
		CreatedAt:     result.Transfer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	return resp, nil
}

// FindAllByAccountID returns all transfers for a given account as sender or receiver
func (s *Service) FindAllByAccountID(ctx context.Context, accountID string) ([]CreateTransferResponse, error) {
	transfers, err := s.repo.FindAllByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	responses := make([]CreateTransferResponse, 0, len(transfers))
	for _, t := range transfers {
		responses = append(responses, CreateTransferResponse{
			ID:            t.ID.String(),
			FromAccountID: t.FromAccountID.String(),
			ToAccountID:   t.ToAccountID.String(),
			Amount:        t.Amount,
			CreatedAt:     t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return responses, nil
}

// isAccountBelongsToUser checks if the account belongs to the user
func (s *Service) isAccountBelongsToUser(ctx context.Context, accountID, userID string) bool {
	db := s.repo.Repository.DB
	var account models.Account
	aid, err := uuid.Parse(accountID)
	if err != nil {
		return false
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return false
	}
	err = db.WithContext(ctx).Where("id = ? AND user_id = ?", aid, uid).First(&account).Error
	return err == nil
}
