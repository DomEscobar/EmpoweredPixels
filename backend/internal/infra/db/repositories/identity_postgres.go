package repositories

import (
	"context"
	"errors"
	"time"

	"empoweredpixels/internal/domain/identity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

type TokenRepository struct {
	pool *pgxpool.Pool
}

type VerificationRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

func NewVerificationRepository(pool *pgxpool.Pool) *VerificationRepository {
	return &VerificationRepository{pool: pool}
}

func (r *UserRepository) FindByNameOrEmail(ctx context.Context, value string) (*identity.User, error) {
	const query = `
		select id, name, email, password, salt, is_verified, created, last_login, banned
		from users
		where name = $1 or email = $1
		limit 1`

	user := identity.User{}
	err := r.pool.QueryRow(ctx, query, value).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Salt,
		&user.IsVerified,
		&user.Created,
		&user.LastLogin,
		&user.Banned,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*identity.User, error) {
	const query = `
		select id, name, email, password, salt, is_verified, created, last_login, banned
		from users
		where id = $1`

	user := identity.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Salt,
		&user.IsVerified,
		&user.Created,
		&user.LastLogin,
		&user.Banned,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *identity.User) error {
	const query = `
		insert into users (name, email, password, salt, is_verified, created, last_login, banned)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning id`

	return r.pool.QueryRow(ctx, query,
		user.Name,
		user.Email,
		user.Password,
		user.Salt,
		user.IsVerified,
		user.Created,
		user.LastLogin,
		user.Banned,
	).Scan(&user.ID)
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	const query = `
		update users
		set last_login = $1
		where id = $2`

	_, err := r.pool.Exec(ctx, query, time.Now(), userID)
	return err
}

func (r *TokenRepository) FindByUserID(ctx context.Context, userID int64) (*identity.Token, error) {
	const query = `
		select id, user_id, value, refresh_value, issued
		from tokens
		where user_id = $1`

	token := identity.Token{}
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&token.ID,
		&token.UserID,
		&token.Value,
		&token.RefreshValue,
		&token.Issued,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) FindByUserIDAndRefresh(ctx context.Context, userID int64, refresh string) (*identity.Token, error) {
	const query = `
		select id, user_id, value, refresh_value, issued
		from tokens
		where user_id = $1 and refresh_value = $2`

	token := identity.Token{}
	err := r.pool.QueryRow(ctx, query, userID, refresh).Scan(
		&token.ID,
		&token.UserID,
		&token.Value,
		&token.RefreshValue,
		&token.Issued,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) Upsert(ctx context.Context, token *identity.Token) error {
	const query = `
		insert into tokens (id, user_id, value, refresh_value, issued)
		values ($1, $2, $3, $4, $5)
		on conflict (user_id)
		do update set value = excluded.value,
					  refresh_value = excluded.refresh_value,
					  issued = excluded.issued,
					  id = excluded.id`

	_, err := r.pool.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.Value,
		token.RefreshValue,
		token.Issued,
	)
	return err
}

func (r *VerificationRepository) Create(ctx context.Context, verification *identity.Verification) error {
	const query = `
		insert into verifications (id, user_id)
		values ($1, $2)`

	_, err := r.pool.Exec(ctx, query, verification.ID, verification.UserID)
	return err
}

func (r *VerificationRepository) DeleteByID(ctx context.Context, id string) (bool, error) {
	const query = `
		delete from verifications
		where id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
