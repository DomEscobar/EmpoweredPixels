package identity

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidRefresh     = errors.New("invalid refresh token")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidVerification = errors.New("invalid verification")
)
