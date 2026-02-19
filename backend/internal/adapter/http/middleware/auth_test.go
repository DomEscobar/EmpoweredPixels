package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"empoweredpixels/internal/domain/identity"
	identityusecase "empoweredpixels/internal/usecase/identity"
	"github.com/golang-jwt/jwt/v5"
)

type mockTokenRepository struct {
	findByValue func(ctx context.Context, value string) (*identity.Token, error)
}

func (m *mockTokenRepository) FindByUserID(ctx context.Context, userID int64) (*identity.Token, error) {
	return nil, nil
}

func (m *mockTokenRepository) FindByUserIDAndRefresh(ctx context.Context, userID int64, refresh string) (*identity.Token, error) {
	return nil, nil
}

func (m *mockTokenRepository) FindByValue(ctx context.Context, value string) (*identity.Token, error) {
	if m.findByValue != nil {
		return m.findByValue(ctx, value)
	}
	return nil, nil
}

func (m *mockTokenRepository) Upsert(ctx context.Context, token *identity.Token) error {
	return nil
}

var _ identityusecase.TokenRepository = (*mockTokenRepository)(nil)

func createValidToken(t *testing.T, secret []byte, userID int64) string {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  "auth",
		"exp":  jwt.NewNumericDate(now.Add(24 * time.Hour)),
		"iat":  jwt.NewNumericDate(now),
		"name": userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		t.Fatal(err)
	}
	return signed
}

func TestWithUserID_TokenExistsInDB(t *testing.T) {
	secret := []byte("test-secret-key-that-is-at-least-32-characters-long")
	userID := int64(123)
	token := createValidToken(t, secret, userID)

	mockRepo := &mockTokenRepository{
		findByValue: func(ctx context.Context, value string) (*identity.Token, error) {
			if value == token {
				return &identity.Token{UserID: userID, Value: token}, nil
			}
			return nil, nil
		},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserID(r.Context())
		if !ok {
			t.Error("userID not found in context")
			return
		}
		if id != userID {
			t.Errorf("expected userID %d, got %d", userID, id)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := WithUserID(next, secret, mockRepo)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestWithUserID_TokenNotInDB(t *testing.T) {
	secret := []byte("test-secret-key-that-is-at-least-32-characters-long")
	userID := int64(123)
	token := createValidToken(t, secret, userID)

	mockRepo := &mockTokenRepository{
		findByValue: func(ctx context.Context, value string) (*identity.Token, error) {
			return nil, nil
		},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserID(r.Context())
		if ok {
			t.Errorf("expected no userID in context, got %d", id)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := WithUserID(next, secret, mockRepo)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestWithUserID_NilTokenRepository(t *testing.T) {
	secret := []byte("test-secret-key-that-is-at-least-32-characters-long")
	userID := int64(123)
	token := createValidToken(t, secret, userID)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserID(r.Context())
		if !ok {
			t.Error("userID not found in context")
			return
		}
		if id != userID {
			t.Errorf("expected userID %d, got %d", userID, id)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := WithUserID(next, secret, nil)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}
