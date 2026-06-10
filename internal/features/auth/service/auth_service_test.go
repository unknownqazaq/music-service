package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"music-service/internal/features/auth/service"
	users_model "music-service/internal/features/users/model"
)

// --- Mock UserRepository ---

type mockUserRepo struct {
	user *users_model.User
	err  error
}

func (m *mockUserRepo) Create(ctx context.Context, u *users_model.User) (*users_model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	u.ID = 1
	return u, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*users_model.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*users_model.User, error) {
	return m.user, m.err
}

// --- Mock RefreshTokenRepository ---

type mockRefreshRepo struct {
	tokens map[string]int64
}

func newMockRefreshRepo() *mockRefreshRepo {
	return &mockRefreshRepo{tokens: make(map[string]int64)}
}

func (m *mockRefreshRepo) Save(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	m.tokens[token] = userID
	return nil
}

func (m *mockRefreshRepo) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	if id, ok := m.tokens[token]; ok {
		return id, nil
	}
	return 0, service.ErrInvalidRefreshToken
}

func (m *mockRefreshRepo) Delete(ctx context.Context, token string) error {
	delete(m.tokens, token)
	return nil
}

func (m *mockRefreshRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	for k, v := range m.tokens {
		if v == userID {
			delete(m.tokens, k)
		}
	}
	return nil
}

func newTestAuthService(userRepo *mockUserRepo) (*service.AuthService, *mockRefreshRepo) {
	refreshRepo := newMockRefreshRepo()
	svc := service.NewAuthService(
		userRepo,
		refreshRepo,
		"test_access_secret",
		"test_refresh_secret",
		15*time.Minute,
		24*time.Hour,
	)
	return svc, refreshRepo
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)
	return string(h)
}

// --- Tests ---

func TestAuthService_Register_Success(t *testing.T) {
	repo := &mockUserRepo{}
	svc, _ := newTestAuthService(repo)

	result, err := svc.Register(context.Background(), service.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "USER", result.Role)
	assert.Equal(t, "FREE", result.SubscriptionType)
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := &mockUserRepo{
		user: &users_model.User{
			ID:               1,
			Email:            "test@example.com",
			PasswordHash:     hashPassword(t, "password123"),
			Role:             "USER",
			SubscriptionType: "FREE",
		},
	}
	svc, _ := newTestAuthService(repo)

	result, err := svc.Login(context.Background(), service.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	repo := &mockUserRepo{
		user: &users_model.User{
			ID:           1,
			Email:        "test@example.com",
			PasswordHash: hashPassword(t, "correctpassword"),
		},
	}
	svc, _ := newTestAuthService(repo)

	result, err := svc.Login(context.Background(), service.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	assert.Nil(t, result)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		user: nil,
		err:  assert.AnError,
	}
	svc, _ := newTestAuthService(repo)

	result, err := svc.Login(context.Background(), service.LoginInput{
		Email:    "notfound@example.com",
		Password: "password",
	})

	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	assert.Nil(t, result)
}

func TestAuthService_Refresh_Success(t *testing.T) {
	repo := &mockUserRepo{
		user: &users_model.User{
			ID:               1,
			Email:            "test@example.com",
			PasswordHash:     hashPassword(t, "password123"),
			Role:             "USER",
			SubscriptionType: "FREE",
		},
	}
	svc, refreshRepo := newTestAuthService(repo)

	// Логин для получения refresh token
	loginResult, err := svc.Login(context.Background(), service.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	oldRefreshToken := loginResult.RefreshToken

	// Обновление токена
	refreshResult, err := svc.Refresh(context.Background(), oldRefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, refreshResult.AccessToken)
	assert.NotEmpty(t, refreshResult.RefreshToken)
	// Старый токен должен быть удалён (ротация)
	_, exists := refreshRepo.tokens[oldRefreshToken]
	assert.False(t, exists, "old refresh token should be deleted after rotation")
}

func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	repo := &mockUserRepo{}
	svc, _ := newTestAuthService(repo)

	result, err := svc.Refresh(context.Background(), "invalid_token_xyz")
	assert.ErrorIs(t, err, service.ErrInvalidRefreshToken)
	assert.Nil(t, result)
}
