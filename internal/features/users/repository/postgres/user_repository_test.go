package postgres_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"music-service/internal/features/users/model"
	repo "music-service/internal/features/users/repository/postgres"
)

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	ctx := context.Background()

	home, err := os.UserHomeDir()
	if err == nil {
		colimaSocket := filepath.Join(home, ".colima/default/docker.sock")
		if _, err := os.Stat(colimaSocket); err == nil && os.Getenv("DOCKER_HOST") == "" {
			os.Setenv("DOCKER_HOST", "unix://"+colimaSocket)
			os.Setenv("TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE", "/var/run/docker.sock")
		}
	}

	pwd, err := os.Getwd()
	require.NoError(t, err)
	initScript := filepath.Join(pwd, "../../../../../migrations/001_init.up.sql")

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithInitScripts(initScript),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sqlx.Open("pgx", connStr)
	require.NoError(t, err)

	require.NoError(t, db.Ping())

	cleanup := func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return db, cleanup
}

func TestUserRepository_CreateAndGetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repository := repo.NewUserRepository(db)

	userToCreate := &model.User{
		Email:            "test@music.com",
		PasswordHash:     "hashedpassword123",
		Username:         "testmusicuser",
		Role:             model.RoleUser,
		SubscriptionType: model.SubscriptionFree,
	}

	// Test Create
	createdUser, err := repository.Create(context.Background(), userToCreate)
	require.NoError(t, err)
	require.NotZero(t, createdUser.ID)
	assert.Equal(t, userToCreate.Email, createdUser.Email)
	assert.Equal(t, userToCreate.Username, createdUser.Username)

	// Test GetByID
	foundUser, err := repository.GetByID(context.Background(), createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, createdUser.Email, foundUser.Email)

	// Test GetByEmail
	foundByEmail, err := repository.GetByEmail(context.Background(), createdUser.Email)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, foundByEmail.ID)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repository := repo.NewUserRepository(db)

	user, err := repository.GetByID(context.Background(), 99999)
	require.ErrorIs(t, err, repo.ErrUserNotFound)
	assert.Nil(t, user)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repository := repo.NewUserRepository(db)

	u1 := &model.User{
		Email:            "duplicate@music.com",
		PasswordHash:     "hash",
		Username:         "user1",
		Role:             model.RoleUser,
		SubscriptionType: model.SubscriptionFree,
	}
	_, err := repository.Create(context.Background(), u1)
	require.NoError(t, err)

	u2 := &model.User{
		Email:            "duplicate@music.com",
		PasswordHash:     "hash",
		Username:         "user2",
		Role:             model.RoleUser,
		SubscriptionType: model.SubscriptionFree,
	}
	_, err = repository.Create(context.Background(), u2)
	require.ErrorIs(t, err, repo.ErrEmailAlreadyTaken)
}

func TestUserRepository_UpdateProfile_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repository := repo.NewUserRepository(db)

	u := &model.User{
		Email:            "original@music.com",
		PasswordHash:     "hash",
		Username:         "originaluser",
		Role:             model.RoleUser,
		SubscriptionType: model.SubscriptionFree,
	}
	created, err := repository.Create(context.Background(), u)
	require.NoError(t, err)

	// Update username only
	newUsername := "updateduser"
	updated, err := repository.UpdateProfile(context.Background(), created.ID, nil, &newUsername)
	require.NoError(t, err)
	assert.Equal(t, newUsername, updated.Username)
	assert.Equal(t, created.Email, updated.Email)

	// Update email only
	newEmail := "updated@music.com"
	updated, err = repository.UpdateProfile(context.Background(), created.ID, &newEmail, nil)
	require.NoError(t, err)
	assert.Equal(t, newEmail, updated.Email)
	assert.Equal(t, newUsername, updated.Username)

	// Update both
	newEmail2 := "both@music.com"
	newUsername2 := "bothuser"
	updated, err = repository.UpdateProfile(context.Background(), created.ID, &newEmail2, &newUsername2)
	require.NoError(t, err)
	assert.Equal(t, newEmail2, updated.Email)
	assert.Equal(t, newUsername2, updated.Username)
}

func TestUserRepository_UpdateProfile_Duplicates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repository := repo.NewUserRepository(db)

	u1, err := repository.Create(context.Background(), &model.User{
		Email:            "user1@music.com",
		PasswordHash:     "hash",
		Username:         "user1",
		Role:             model.RoleUser,
		SubscriptionType: model.SubscriptionFree,
	})
	require.NoError(t, err)

	u2, err := repository.Create(context.Background(), &model.User{
		Email:            "user2@music.com",
		PasswordHash:     "hash",
		Username:         "user2",
		Role:             model.RoleUser,
		SubscriptionType: model.SubscriptionFree,
	})
	require.NoError(t, err)

	// Update u1's email to u2's email
	_, err = repository.UpdateProfile(context.Background(), u1.ID, &u2.Email, nil)
	assert.ErrorIs(t, err, repo.ErrEmailAlreadyTaken)

	// Update u1's username to u2's username
	_, err = repository.UpdateProfile(context.Background(), u1.ID, nil, &u2.Username)
	assert.ErrorIs(t, err, repo.ErrUsernameTaken)
}

