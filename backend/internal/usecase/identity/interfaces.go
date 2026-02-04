package identity

import (
	"context"

	"empoweredpixels/internal/domain/identity"
)

type UserRepository interface {
	FindByNameOrEmail(ctx context.Context, value string) (*identity.User, error)
	FindByID(ctx context.Context, id int64) (*identity.User, error)
	Create(ctx context.Context, user *identity.User) error
	UpdateLastLogin(ctx context.Context, userID int64) error
}

type TokenRepository interface {
	FindByUserID(ctx context.Context, userID int64) (*identity.Token, error)
	FindByUserIDAndRefresh(ctx context.Context, userID int64, refresh string) (*identity.Token, error)
	Upsert(ctx context.Context, token *identity.Token) error
}

type VerificationRepository interface {
	Create(ctx context.Context, verification *identity.Verification) error
	DeleteByID(ctx context.Context, id string) (bool, error)
}
