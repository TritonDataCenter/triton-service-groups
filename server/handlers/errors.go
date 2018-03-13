package handlers

import "errors"

var (
	ErrNoConnPool = errors.New("handlers can't access database pool")
	ErrFailedAuth = errors.New("failed request authentication")
	ErrNoSession  = errors.New("failed to get authenticated session")
)
