package account

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ahmedkhaeld/banking-app/db"
	"github.com/ahmedkhaeld/banking-app/db/models"
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
		repo.Repository.DB.Exec("DELETE FROM accounts")
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
	if err := db.DB.AutoMigrate(&models.Account{}); err != nil {
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

func TestCreateAccount_Success(t *testing.T) {
	service := InitService()
	usr := createTestUser(t)
	balance := int64(1000)
	req := CreateAccountRequest{
		Currency: "USD",
		Balance:  &balance,
	}
	resp, err := service.createAccount(req, usr.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, usr.ID.String(), resp.UserID)
	assert.Equal(t, usr.Username, resp.Owner)
	assert.Equal(t, req.Currency, resp.Currency)
	assert.Equal(t, balance, resp.Balance)
}

func TestCreateAccount_UserDoesNotExist(t *testing.T) {
	service := InitService()
	balance := int64(1000)
	req := CreateAccountRequest{
		Currency: "USD",
		Balance:  &balance,
	}
	fakeUserID := uuid.New().String()
	resp, err := service.createAccount(req, fakeUserID)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetAccountBalance_Success(t *testing.T) {
	service := InitService()
	usr := createTestUser(t)
	balance := int64(500)
	accReq := CreateAccountRequest{
		Currency: "EUR",
		Balance:  &balance,
	}
	accResp, err := service.createAccount(accReq, usr.ID.String())
	assert.NoError(t, err)
	balResp, err := service.getAccountBalance(accResp.ID, usr.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, accResp.ID, balResp.ID)
	assert.Equal(t, balance, balResp.Balance)
	assert.Equal(t, accReq.Currency, balResp.Currency)
}

func TestUpdateBalance_Success(t *testing.T) {
	service := InitService()
	usr := createTestUser(t)
	initBalance := int64(200)
	accReq := CreateAccountRequest{
		Currency: "EGP",
		Balance:  &initBalance,
	}
	accResp, err := service.createAccount(accReq, usr.ID.String())
	assert.NoError(t, err)
	addAmount := int64(300)
	updated, err := service.updateBalance(context.Background(), accResp.ID, addAmount)
	assert.NoError(t, err)
	assert.Equal(t, initBalance+addAmount, updated.Balance)
}

func TestUpdateBalance_InvalidAccountID(t *testing.T) {
	service := InitService()
	_, err := service.updateBalance(context.Background(), "not-a-uuid", 100)
	assert.Error(t, err)
}

func TestIsAccountOwnedByUser(t *testing.T) {
	service := InitService()
	usr := createTestUser(t)
	balance := int64(100)
	accReq := CreateAccountRequest{
		Currency: "CAD",
		Balance:  &balance,
	}
	accResp, err := service.createAccount(accReq, usr.ID.String())
	assert.NoError(t, err)
	// Should be true for owner
	ok := service.isAccountOwnedByUser(accResp.ID, usr.ID.String())
	assert.True(t, ok)
	// Should be false for random user
	ok = service.isAccountOwnedByUser(accResp.ID, uuid.New().String())
	assert.False(t, ok)
}
