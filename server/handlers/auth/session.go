package auth

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// authSession a private struct which is only accessible by pulling out of the
// current request `context.Context`.
type Session struct {
	AccountID      string
	KeyFingerprint string

	*parsedRequest
}

// IsAuthenticated encapsulates whatever it means for an authSession to be
// deemed authenticated.
func (a Session) IsAuthenticated() bool {
	return a.AccountID != ""
}

func NewSession(req *http.Request) (Session, error) {
	parsedReq, err := parseRequest(req)
	if err != nil {
		return Session{}, errors.Wrap(err, "failed to parse auth request")
	}

	return Session{
		parsedRequest: parsedReq,
	}, nil
}

// EnsureKey checks Triton for an active TSG account key. If we cannot find one,
// create one and store it in Triton as well as Vault.
//
// Other edge cases will be developed later like the account having multiple TSG
// keys or no active keys but Vault stored keys, etc.
func (s Session) EnsureKey(ctx context.Context) error {
	keychain := NewKeychain(s.parsedRequest)

	if err := keychain.CheckTriton(ctx); err != nil {
		return errors.Wrap(err, "failed to check triton account keys")
	}

	if !keychain.HasKey() {
		keypair, err := NewKeyPair(1024)
		if err != nil {
			return errors.Wrap(err, "failed to generate new TSG key")
		}

		err = keychain.AddKey(ctx, keypair)
		if err != nil {
			return errors.Wrap(err, "failed to add TSG key in Triton")
		}
	}

	if keychain.HasKey() {
		s.AccountID = s.parsedRequest.accountName
		s.KeyFingerprint = keychain.AccountKey.Fingerprint
	}

	return nil
}
