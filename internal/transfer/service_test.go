package transfer

import (
	"log"
	"os"
	"testing"

	"context"

	"github.com/ahmedkhaeld/banking-app/db"
	"github.com/ahmedkhaeld/banking-app/db/models"
	"github.com/ahmedkhaeld/banking-app/internal/account"
	"github.com/ahmedkhaeld/banking-app/internal/auth"
	"github.com/ahmedkhaeld/banking-app/internal/user"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Use the real DB from  db package
func setupTestRepository(t *testing.T) *Repository {
	repo := InitRepository()
	t.Cleanup(func() {
		repo.Repository.DB.Exec("DELETE FROM transfers")
		repo.Repository.DB.Exec("DELETE FROM entries")
		repo.Repository.DB.Exec("DELETE FROM accounts")
		repo.Repository.DB.Exec("DELETE FROM users")

	})
	return repo
}

func setupTestService(t *testing.T) *Service {
	repo := setupTestRepository(t)
	return NewService(repo)
}

func TestMain(m *testing.M) {
	// load the environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Connect to the test database
	dsn := os.Getenv("DB_SOURCE_TEST")
	if err := db.Open(dsn); err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	err := db.AddUUIDExtension()
	if err != nil {
		panic("failed to add UUID extension: " + err.Error())
	}

	// Run migrations
	if err := db.DB.AutoMigrate(&models.Account{}, &models.Transfer{}, &models.Entry{}, &models.User{}); err != nil {
		panic("failed to run migrations: " + err.Error())
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func createTestUser(t *testing.T) *models.User {
	hashed, err := auth.HashPassword("password123")
	assert.NoError(t, err)
	userModel := &models.User{
		ID:       uuid.New(),
		Username: "testuser_" + uuid.New().String()[:8],
		Password: hashed,
		FullName: "Test User",
		Email:    "test_" + uuid.New().String()[:8] + "@example.com",
	}
	repo := user.InitRepository()
	err = repo.Repository.DB.Create(userModel).Error
	assert.NoError(t, err)
	return userModel
}

func createTestAccount(t *testing.T, userID uuid.UUID, username string, balance int64, currency string) *models.Account {
	acc := &models.Account{
		ID:       uuid.New(),
		UserID:   userID,
		Owner:    username,
		Balance:  balance,
		Currency: currency,
	}
	repo := account.InitRepository()
	err := repo.Repository.DB.Create(acc).Error
	assert.NoError(t, err)
	return acc
}

func TestTransfer_Success(t *testing.T) {
	service := setupTestService(t)
	user1 := createTestUser(t)
	user2 := createTestUser(t)
	acc1 := createTestAccount(t, user1.ID, user1.Username, 1000, "USD")
	acc2 := createTestAccount(t, user2.ID, user2.Username, 100, "USD")

	req := CreateTransferRequest{
		FromAccountID: acc1.ID.String(),
		ToAccountID:   acc2.ID.String(),
		Amount:        200,
	}
	resp, err := service.Transfer(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Amount, resp.Amount)
	assert.Equal(t, req.FromAccountID, resp.FromAccountID)
	assert.Equal(t, req.ToAccountID, resp.ToAccountID)

	// Check balances updated
	repo := account.InitRepository()
	var updatedAcc1, updatedAcc2 models.Account
	err = repo.Repository.DB.First(&updatedAcc1, "id = ?", acc1.ID).Error
	assert.NoError(t, err)
	err = repo.Repository.DB.First(&updatedAcc2, "id = ?", acc2.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(800), updatedAcc1.Balance)
	assert.Equal(t, int64(300), updatedAcc2.Balance)
}

func TestTransfer_InvalidFromAccount(t *testing.T) {
	service := setupTestService(t)
	user := createTestUser(t)
	acc := createTestAccount(t, user.ID, user.Username, 100, "USD")
	req := CreateTransferRequest{
		FromAccountID: uuid.New().String(), // invalid (does not exist)
		ToAccountID:   acc.ID.String(),
		Amount:        50,
	}
	resp, err := service.Transfer(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestTransfer_InvalidToAccount(t *testing.T) {
	service := setupTestService(t)
	user := createTestUser(t)
	acc := createTestAccount(t, user.ID, user.Username, 100, "USD")
	req := CreateTransferRequest{
		FromAccountID: acc.ID.String(),
		ToAccountID:   uuid.New().String(), // invalid (does not exist)
		Amount:        50,
	}
	resp, err := service.Transfer(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestTransfer_NegativeAmount(t *testing.T) {
	service := setupTestService(t)
	user1 := createTestUser(t)
	user2 := createTestUser(t)
	acc1 := createTestAccount(t, user1.ID, user1.Username, 1000, "USD")
	acc2 := createTestAccount(t, user2.ID, user2.Username, 100, "USD")
	req := CreateTransferRequest{
		FromAccountID: acc1.ID.String(),
		ToAccountID:   acc2.ID.String(),
		Amount:        -100,
	}
	resp, err := service.Transfer(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestTransfer_ZeroAmount(t *testing.T) {
	service := setupTestService(t)
	user1 := createTestUser(t)
	user2 := createTestUser(t)
	acc1 := createTestAccount(t, user1.ID, user1.Username, 1000, "USD")
	acc2 := createTestAccount(t, user2.ID, user2.Username, 100, "USD")
	req := CreateTransferRequest{
		FromAccountID: acc1.ID.String(),
		ToAccountID:   acc2.ID.String(),
		Amount:        0,
	}
	resp, err := service.Transfer(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
