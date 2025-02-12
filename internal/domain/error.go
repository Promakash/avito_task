package domain

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrLowBalance       = errors.New("not enough coins on the balance")
	ErrMerchNotFound    = errors.New("merch not found")
	ErrUserExists       = errors.New("user already exist")
	ErrInvalidAuthToken = errors.New("invalid auth token")
	ErrUnauthorized     = errors.New("unauthorized")
)
