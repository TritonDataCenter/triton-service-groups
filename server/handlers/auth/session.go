package auth

import (
	"net/http"

	"github.com/pkg/errors"
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
	emptySession := Session{}

	parsedReq, err := parseRequest(req)
	if err != nil {
		return emptySession, errors.Wrap(err, "could not parse auth request")
	}

	tritonKeys := NewTritonKeys(parsedReq)
	if err := tritonKeys.Check(); err != nil {
		return emptySession, errors.Wrap(err, "could not check triton account keys")
	}

	if !tritonKeys.HasTSG() {
		log.Debug().Msg("--- couldn't find any TSG keys for Triton account")
	}

	return Session{
		AccountID: parsedReq.accountName,
		// KeyFingerprint: parsed.fingerprint,
	}, nil
}
