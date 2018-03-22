package auth

import "errors"

var (
	ErrUnauthRequest = errors.New("received unauthenticated request")
	ErrMissingSig    = errors.New("missing signature within auth header")
	ErrBadKeyID      = errors.New("couldn't parse keyId within header")
	ErrParseAuth     = errors.New("failed to parse values from keyId")
	ErrParseValue    = errors.New("incorrect values parsed from keyId")
	ErrNameLen       = errors.New("parsed name is too short")
	ErrNameFormat    = errors.New("parsed name is not formatted properly")

	ErrWhitelist = errors.New("service only accessible by whitelist")
)
