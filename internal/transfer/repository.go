package transfer

import (
	"context"
	"errors"

	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/db"
	"github.com/ahmedkhaeld/banking-app/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type model = models.Transfer

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

// TransferTxParams holds the parameters for a transfer transaction
// (useful for service and controller layers)
type TransferTxParams struct {
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

// TransferTxResult holds the result of a transfer transaction
type TransferTxResult struct {
	Transfer    models.Transfer
	FromEntry   models.Entry
	ToEntry     models.Entry
	FromAccount models.Account
	ToAccount   models.Account
}

func (r *Repository) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	fromID, err := uuid.Parse(args.FromAccountID)
	if err != nil {
		return result, errors.New("invalid from_account_id")
	}
	toID, err := uuid.Parse(args.ToAccountID)
	if err != nil {
		return result, errors.New("invalid to_account_id")
	}
	if args.Amount <= 0 {
		return result, errors.New("amount must be positive")
	}
	err = r.Repository.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Create Transfer
		transfer := models.Transfer{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        args.Amount,
		}
		if err := tx.Create(&transfer).Error; err != nil {
			return err
		}
		result.Transfer = transfer

		// Step 2: Create Entries
		fromEntry := models.Entry{
			AccountID: fromID,
			Amount:    -args.Amount,
		}
		if err := tx.Create(&fromEntry).Error; err != nil {
			return err
		}
		result.FromEntry = fromEntry

		toEntry := models.Entry{
			AccountID: toID,
			Amount:    args.Amount,
		}
		if err := tx.Create(&toEntry).Error; err != nil {
			return err
		}
		result.ToEntry = toEntry

		// Step 3: Update balances (avoid deadlock by ordering by ID)
		var fromAccount, toAccount models.Account
		if fromID.String() < toID.String() {
			if err := updateBalance(tx, fromID, -args.Amount, &fromAccount); err != nil {
				return err
			}
			if err := updateBalance(tx, toID, args.Amount, &toAccount); err != nil {
				return err
			}
		} else {
			if err := updateBalance(tx, toID, args.Amount, &toAccount); err != nil {
				return err
			}
			if err := updateBalance(tx, fromID, -args.Amount, &fromAccount); err != nil {
				return err
			}
		}
		result.FromAccount = fromAccount
		result.ToAccount = toAccount
		return nil
	})
	return result, err
}

// updateBalance updates the balance of an account and returns the updated account
func updateBalance(tx *gorm.DB, accountID uuid.UUID, amount int64, account *models.Account) error {
	if err := tx.Model(&models.Account{}).Where("id = ?", accountID).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		return err
	}
	return tx.Where("id = ?", accountID).First(account).Error
}

// FindAllByAccountID returns all transfers where the account is either the sender or receiver
func (r *Repository) FindAllByAccountID(ctx context.Context, accountID string) ([]models.Transfer, error) {
	var transfers []models.Transfer
	id, err := uuid.Parse(accountID)
	if err != nil {
		return nil, errors.New("invalid account_id")
	}
	err = r.Repository.DB.WithContext(ctx).
		Where("from_account_id = ? OR to_account_id = ?", id, id).
		Order("created_at DESC").
		Find(&transfers).Error
	if err != nil {
		return nil, err
	}
	return transfers, nil
}
