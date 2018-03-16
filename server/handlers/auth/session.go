package auth

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// authSession a private struct which is only accessible by pulling out of the
// current request `context.Context`.
type Session struct {
	AccountID      string
	KeyFingerprint string

	*parsedRequest
}

func NewSession(req *http.Request) (*Session, error) {
	parsedReq, err := parseRequest(req)
	if err != nil {
		return &Session{}, errors.Wrap(err, "failed to parse auth request")
	}

	return &Session{
		parsedRequest: parsedReq,
	}, nil
}

// IsAuthenticated encapsulates whatever it means for an authSession to be
// deemed authenticated.
func (a *Session) IsAuthenticated() bool {
	return a.AccountID != "" && a.KeyFingerprint != ""
}

// EnsureKey checks Triton for an active TSG account key. If we cannot find one,
// create one and store it in Triton as well as Vault.
//
// Other edge cases will be developed later like the account having multiple TSG
// keys or no active keys but Vault stored keys, etc.
func (s *Session) EnsureKey(ctx context.Context) error {
	keychain := NewKeychain(s.parsedRequest)

	if err := keychain.CheckTriton(ctx); err != nil {
		err = errors.Wrap(err, "failed to check triton keys")
		log.Error().Err(err)
		return err
	}

	// NOTE(justinwr): this is duplicate logic from below but I wanted
	// differentiating debug logs between creating/adding and existing
	if keychain.HasKey() {
		log.Debug().
			Str("account", s.parsedRequest.accountName).
			Str("fingerprint", keychain.AccountKey.Fingerprint).
			Msg("auth: found existing key in Triton")

		s.AccountID = s.parsedRequest.accountName
		s.KeyFingerprint = keychain.AccountKey.Fingerprint

		return nil
	}

	keypair, err := NewKeyPair(1024)
	if err != nil {
		err = errors.Wrap(err, "failed to generate new keypair")
		log.Error().Err(err)
		return err
	}

	err = keychain.AddKey(ctx, keypair)
	if err != nil {
		err = errors.Wrap(err, "failed to add new key")
		log.Error().Err(err)
		return err
	}

	if keychain.HasKey() {
		log.Debug().
			Str("account", s.parsedRequest.accountName).
			Str("fingerprint", keychain.AccountKey.Fingerprint).
			Msg("auth: successfully created and stored new Triton key")

		s.AccountID = s.parsedRequest.accountName
		s.KeyFingerprint = keychain.AccountKey.Fingerprint
	}

	return nil
}
