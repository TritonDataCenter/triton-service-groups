package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/joyent/triton-service-groups/accounts"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// authSession a private struct which is only accessible by pulling out of the
// current request `context.Context`.
type Session struct {
	*ParsedRequest

	AccountID      int
	KeyFingerprint string

	devMode bool
}

const (
	testAccountID      = 332378521158418433
	testKeyFingerprint = "12:34:56:78:90:12:34:56:78:90:12:34:56:78:90:AB"
)

// NewSession constructs and returns a new Session by parsing the HTTP request,
// validating and pulling out authentication headers.
func NewSession(req *http.Request) (*Session, error) {
	if devMode := os.Getenv("TSG_DEV_MODE"); devMode == "true" {
		return &Session{
			AccountID:      testAccountID,
			KeyFingerprint: testKeyFingerprint,
			devMode:        true,
		}, nil
	}

	parsedReq, err := ParseRequest(req)
	if err != nil {
		return &Session{}, errors.Wrap(err, "failed to parse auth request")
	}

	return &Session{
		ParsedRequest: parsedReq,
	}, nil
}

// IsAuthenticated encapsulates whatever it means for an authSession to be
// deemed authenticated.
func (a *Session) IsAuthenticated() bool {
	return a.AccountID != 0 && a.KeyFingerprint != ""
}

func (s *Session) EnsureAccount(ctx context.Context, store *accounts.Store) error {
	if s.devMode {
		log.Debug().
			Int("account_id", s.AccountID).
			Msg("auth: ignoring account via TSG_DEV_MODE")

		return nil
	}

	check := NewAccountCheck(s.ParsedRequest, store)

	if err := check.OnTriton(ctx); err != nil {
		return err
	}

	if !check.HasAccount() {
		return errors.New("could not authenticate account with triton")
	}

	if err := check.SaveAccount(ctx); err != nil {
		return err
	}

	return nil
}

// EnsureKey checks Triton for an active TSG account key. If we cannot find one,
// create one and store it in Triton as well as Vault.
//
// Other edge cases will be developed later like the account having multiple TSG
// keys or no active keys but Vault stored keys, etc.
func (s *Session) EnsureKey(ctx context.Context) error {
	if s.devMode {
		log.Debug().
			Int("account_id", s.AccountID).
			Str("fingerprint", s.KeyFingerprint).
			Msg("auth: ignoring authentication via TSG_DEV_MODE")

		return nil
	}

	keychain := NewKeychain(s.ParsedRequest)

	if err := keychain.CheckTriton(ctx); err != nil {
		err = errors.Wrap(err, "failed to check triton keys")
		log.Error().Err(err)
		return err
	}

	if keychain.HasKey() {
		log.Debug().
			Int("account_id", testAccountID).
			Str("fingerprint", keychain.AccountKey.Fingerprint).
			Msg("auth: found existing key in Triton")

		s.AccountID = testAccountID
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
			Int("account_id", testAccountID).
			Str("fingerprint", keychain.AccountKey.Fingerprint).
			Msg("auth: successfully created and stored new Triton key")

		s.AccountID = testAccountID
		s.KeyFingerprint = keychain.AccountKey.Fingerprint
	}

	return nil
}
