package auth

import "errors"

var (
	ErrRefreshTokenExpired  = errors.New("refresh token was expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrCodeConflict         = errors.New("this code is in use")
	ErrInvalidCredentials   = errors.New("invalid credentials")
)
