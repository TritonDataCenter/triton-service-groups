package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joyent/triton-service-groups/accounts"
	"github.com/joyent/triton-service-groups/keys"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/y0ssar1an/q"
)

// authSession a private struct which is only accessible by pulling out of the
// current request `context.Context`.
type Session struct {
	*ParsedRequest

	AccountID   int64
	Fingerprint string

	devMode bool
}

// NewSession constructs and returns a new Session by parsing the HTTP request,
// validating and pulling out authentication headers.
func NewSession(req *http.Request) (*Session, error) {
	if devMode := os.Getenv("TSG_DEV_MODE"); devMode == "true" {
		return &Session{
			AccountID: testAccountID,
			devMode:   true,
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

// IsAuthenticated represents whatever it means for an authSession to be deemed
// authenticated.
func (a *Session) IsAuthenticated() bool {

	return a.AccountID != 0 && a.Fingerprint != ""
}

// EnsureAccount ensures that a Triton account is authentic and an account has
// been created for it within the TSG database. Returns the TSG account that was
// either created or found.
func (s *Session) EnsureAccount(ctx context.Context, store *accounts.Store) (*accounts.Account, error) {
	check := NewAccountCheck(s.ParsedRequest, store)

	if err := check.OnTriton(ctx); err != nil {
		return nil, err
	}

	if !check.HasTritonAccount() {
		return nil, errors.New("could not authenticate account with triton")
	}

	if err := check.SaveAccount(ctx); err != nil {
		return nil, err
	}

	s.AccountID = check.Account.ID

	log.Debug().
		Str("account_id", fmt.Sprintf("%d", s.AccountID)).
		Str("account_name", check.Account.AccountName).
		Msg("auth: session account has been authenticated")

	return check.Account, nil
}

// EnsureKey checks Triton for an active TSG account key. If one cannot be found
// than a new key is created and stored it into the TSG database.
func (s *Session) EnsureKeys(ctx context.Context, acct *accounts.Account, store *keys.Store) error {
	check := NewKeyCheck(s.ParsedRequest, acct, store)

	if err := check.OnTriton(ctx); err != nil {
		err = errors.Wrap(err, "failed to check triton for key")
		log.Debug().Err(err)
		return err
	}

	if err := check.InDatabase(ctx); err != nil {
		err = errors.Wrap(err, "failed to check database for key")
		log.Debug().Err(err)
		return err
	}

	q.Q(check.TritonKey, check.Key)

	if check.HasKey() {
		if check.HasTritonKey() {
			if check.Key.Fingerprint == check.TritonKey.Fingerprint {
				log.Debug().
					Str("account_name", acct.AccountName).
					Str("fingerprint", check.Key.Fingerprint).
					Msg("auth: found existing key with matching fingerprint")

				return nil
			}

			err := errors.New("auth: found conflicting key state")
			log.Error().
				Str("account_name", acct.AccountName).
				Err(err)
			return err
		} else {
			keypair, err := DecodeKeyPair(check.Key.Material)
			if err != nil {
				err = errors.Wrap(err, "failed to generate new keypair")
				log.Error().Err(err)
				return err
			}

			err = check.AddTritonKey(ctx, keypair)
			if err != nil {
				err = errors.Wrap(err, "failed to add new key")
				log.Error().Err(err)
				return err
			}
		}
	} else {
		if check.HasTritonKey() {
			err := errors.New("auth: found key in triton not in tsg")
			log.Error().
				Str("account_name", acct.AccountName).
				Str("fingerprint", check.TritonKey.Fingerprint).
				Err(err)
			return err
		}

		// create a new key
		keypair, err := NewKeyPair(1024)
		if err != nil {
			err = errors.Wrap(err, "failed to generate new keypair")
			log.Error().Err(err)
			return err
		}

		if !check.HasTritonKey() {
			if err := check.AddTritonKey(ctx, keypair); err != nil {
				err = errors.Wrap(err, "failed to add new key")
				log.Error().Err(err)
				return err
			}
		}

		if err := check.InsertKey(ctx, keypair); err != nil {
			err = errors.Wrap(err, "failed to save new key")
			log.Error().Err(err)
			return err
		}
	}

	if check.Key.Fingerprint == check.TritonKey.Fingerprint {
		log.Debug().
			Str("account_name", acct.AccountName).
			Str("fingerprint", check.Key.Fingerprint).
			Msg("auth: found existing key with matching fingerprint")

		s.Fingerprint = check.Key.Fingerprint

		return nil
	}

	return nil
}
