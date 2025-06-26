package user

import (
	"log"
	"os"
	"testing"

	"github.com/ElegantSoft/go-restful-generator/crud"
	"github.com/ahmedkhaeld/banking-app/db"
	"github.com/ahmedkhaeld/banking-app/db/models"
	"github.com/ahmedkhaeld/banking-app/internal/auth"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Use the real DB from your db package
func setupTestRepository(t *testing.T) *Repository {
	repo := InitRepository()
	t.Cleanup(func() {
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
	// This assumes you have a .env file in the root of your project
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
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		panic("failed to run migrations: " + err.Error())
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestCreateUser_Success(t *testing.T) {
	service := setupTestService(t)
	request := &CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		FullName: "Test User",
		Email:    "test@example.com",
	}

	user, err := service.createUser(request)
	assert.NoError(t, err)
	assert.Equal(t, request.Username, user.Username)
	assert.Equal(t, request.Email, user.Email)
	assert.Equal(t, request.FullName, user.FullName)
	assert.NotEqual(t, request.Password, user.Password) // Password should be hashed

	// Verify the user was actually created in the database
	var dbUser models.User
	repo := InitRepository()
	err = repo.Repository.DB.Where("username = ?", request.Username).First(&dbUser).Error
	assert.NoError(t, err)
	assert.Equal(t, request.Username, dbUser.Username)
}

func TestCreateUser_UsernameExists(t *testing.T) {
	service := setupTestService(t)

	// Create a user first using the real model
	testUser := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: "hashed",
		FullName: "Test User",
		Email:    "test@example.com",
	}
	repo := InitRepository()
	err := repo.Repository.DB.Create(testUser).Error
	assert.NoError(t, err)

	request := &CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		FullName: "Test User 2",
		Email:    "test2@example.com",
	}

	user, err := service.createUser(request)
	assert.ErrorIs(t, err, ErrUsernameExists)
	assert.Nil(t, user)
}

func TestCreateUser_EmailExists(t *testing.T) {
	service := setupTestService(t)

	testUser := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: "hashed",
		FullName: "Test User",
		Email:    "test@example.com",
	}
	repo := InitRepository()
	err := repo.Repository.DB.Create(testUser).Error
	assert.NoError(t, err)

	request := &CreateUserRequest{
		Username: "testuser2",
		Password: "password123",
		FullName: "Test User 2",
		Email:    "test@example.com",
	}

	user, err := service.createUser(request)
	assert.ErrorIs(t, err, ErrEmailExists)
	assert.Nil(t, user)
}

func TestLoginUser_Success(t *testing.T) {
	service := setupTestService(t)
	hashed, err := auth.HashPassword("password123")
	assert.NoError(t, err)

	testUser := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: hashed,
		FullName: "Test User",
		Email:    "test@example.com",
	}
	repo := InitRepository()
	err = repo.Repository.DB.Create(testUser).Error
	assert.NoError(t, err)

	result, err := service.loginUser("testuser", "password123")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	service := setupTestService(t)
	hashed, err := auth.HashPassword("password123")
	assert.NoError(t, err)

	testUser := &models.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: hashed,
		FullName: "Test User",
		Email:    "test@example.com",
	}
	repo := InitRepository()
	err = repo.Repository.DB.Create(testUser).Error
	assert.NoError(t, err)

	result, err := service.loginUser("testuser", "wrongpassword")
	assert.ErrorIs(t, err, ErrInvalidUsernameOrPassword)
	assert.Nil(t, result)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	service := setupTestService(t)

	result, err := service.loginUser("nonexistent", "password123")
	assert.ErrorIs(t, err, ErrInvalidUsernameOrPassword)
	assert.Nil(t, result)
}

func TestUpdate_Success(t *testing.T) {
	service := setupTestService(t)
	repo := InitRepository()
	oldUser := &models.User{ID: uuid.New(), Username: "old", Password: "pass", FullName: "Old Name", Email: "old@example.com"}
	err := repo.Repository.DB.Create(oldUser).Error
	assert.NoError(t, err)

	newUser := &models.User{ID: oldUser.ID, Username: "new", Password: "pass", FullName: "New Name", Email: "new@example.com"}
	err = service.Update(oldUser, newUser)
	assert.NoError(t, err)

	var dbUser models.User
	err = repo.Repository.DB.First(&dbUser, "id = ?", oldUser.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "new", dbUser.Username)
	assert.Equal(t, "New Name", dbUser.FullName)
	assert.Equal(t, "new@example.com", dbUser.Email)
}

func TestFindOne_Success(t *testing.T) {
	service := setupTestService(t)
	repo := InitRepository()
	user := &models.User{ID: uuid.New(), Username: "found", Password: "pass", FullName: "Found Name", Email: "found@example.com"}
	err := repo.Repository.DB.Create(user).Error
	assert.NoError(t, err)

	var result models.User
	api := crud.GetAllRequest{
		Filter: []string{"username||eq||found"},
	}
	err = service.FindOne(api, &result)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Username, result.Username)
}

func TestFindOne_Failure(t *testing.T) {
	service := setupTestService(t)
	var result models.User
	api := crud.GetAllRequest{
		Filter: []string{"username||eq||notfound"},
	}
	err := service.FindOne(api, &result)
	assert.Error(t, err)
}
