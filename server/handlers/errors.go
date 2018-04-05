package handlers

import "errors"

var (
	ErrNoConnPool    = errors.New("handlers can't access database pool")
	ErrNoNomadClient = errors.New("handlers can't access nomad client")
	ErrFailedAuth    = errors.New("failed request authentication")
	ErrFailedSession = errors.New("failed session authentication")
	ErrFailedAccount = errors.New("failed account authentication")
	ErrFailedKey     = errors.New("failed key authentication")
	ErrNoSession     = errors.New("failed to get authenticated session")
)
