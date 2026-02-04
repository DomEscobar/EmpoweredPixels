package identity

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"empoweredpixels/internal/domain/identity"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	users         UserRepository
	tokens        TokenRepository
	verifications VerificationRepository
	jwtSecret     []byte
	tokenDays     int
	now           func() time.Time
}

func NewService(
	users UserRepository,
	tokens TokenRepository,
	verifications VerificationRepository,
	jwtSecret string,
	tokenDays int,
	now func() time.Time,
) *Service {
	if now == nil {
		now = time.Now
	}

	return &Service{
		users:         users,
		tokens:        tokens,
		verifications: verifications,
		jwtSecret:     []byte(jwtSecret),
		tokenDays:     tokenDays,
		now:           now,
	}
}

type LoginInput struct {
	User     string
	Password string
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

type TokenOutput struct {
	UserID  int64  `json:"userId"`
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput) error {
	existing, err := s.users.FindByNameOrEmail(ctx, input.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrUserExists
	}

	existing, err = s.users.FindByNameOrEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrUserExists
	}

	salt, err := identity.GenerateSalt()
	if err != nil {
		return err
	}
	passwordHash, err := identity.HashPassword(input.Password, salt)
	if err != nil {
		return err
	}

	user := &identity.User{
		Name:       input.Username,
		Email:      input.Email,
		Password:   passwordHash,
		Salt:       salt,
		Created:    s.now(),
		LastLogin:  s.now(),
		IsVerified: false,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return err
	}

	verification := &identity.Verification{
		ID:     uuid.NewString(),
		UserID: user.ID,
	}

	return s.verifications.Create(ctx, verification)
}

func (s *Service) Verify(ctx context.Context, value string) error {
	removed, err := s.verifications.DeleteByID(ctx, value)
	if err != nil {
		return err
	}
	if !removed {
		return ErrInvalidVerification
	}

	return nil
}

func (s *Service) Token(ctx context.Context, input LoginInput) (*TokenOutput, error) {
	user, err := s.users.FindByNameOrEmail(ctx, input.User)
	if err != nil {
		return nil, err
	}
	if user == nil || !identity.PasswordMatches(input.Password, user.Salt, user.Password) {
		return nil, ErrInvalidCredentials
	}

	if err := s.users.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, err
	}

	return s.issueToken(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, userID int64, refresh string) (*TokenOutput, error) {
	token, err := s.tokens.FindByUserIDAndRefresh(ctx, userID, refresh)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, ErrInvalidRefresh
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidRefresh
	}

	return s.issueToken(ctx, user)
}

func (s *Service) issueToken(ctx context.Context, user *identity.User) (*TokenOutput, error) {
	now := s.now()
	refreshValue, err := randomToken(32)
	if err != nil {
		return nil, err
	}

	exp := jwt.NewNumericDate(now.Add(time.Duration(s.tokenDays) * 24 * time.Hour))
	iat := jwt.NewNumericDate(now)
	claims := jwt.MapClaims{
		"sub":  "auth",
		"exp":  exp,
		"iat":  iat,
		"name": user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	entity := &identity.Token{
		ID:           uuid.NewString(),
		UserID:       user.ID,
		Value:        signed,
		RefreshValue: refreshValue,
		Issued:       now,
	}

	if err := s.tokens.Upsert(ctx, entity); err != nil {
		return nil, err
	}

	return &TokenOutput{
		UserID:  user.ID,
		Token:   signed,
		Refresh: refreshValue,
	}, nil
}

func randomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

