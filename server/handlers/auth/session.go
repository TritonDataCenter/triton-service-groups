package auth

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// authSession a private struct which is only accessible by pulling out of the
// current request `context.Context`.
type Session struct {
	AccountID      string
	KeyFingerprint string
}

// IsAuthenticated encapsulates whatever it means for an authSession to be
// deemed authenticated.
func (a Session) IsAuthenticated() bool {
	return a.AccountID != ""
}

func NewSession(req *http.Request) (Session, error) {
	parsed, err := parseRequest(req)
	if err != nil {
		return Session{}, err
	}

	log.Debug().Msgf("demo authentication found %q", parsed.accountName)

	return Session{
		AccountID: parsed.accountName,
		// KeyFingerprint: parsed.fingerprint,
	}, nil
}
