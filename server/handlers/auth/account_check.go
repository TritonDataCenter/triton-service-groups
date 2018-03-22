package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-service-groups/accounts"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type AccountCheck struct {
	*ParsedRequest
	*accounts.Account

	TritonAccount *account.Account

	config *triton.ClientConfig
	store  *accounts.Store
}

func NewAccountCheck(req *ParsedRequest, store *accounts.Store) *AccountCheck {
	signer := &authentication.TestSigner{}
	config := &triton.ClientConfig{
		TritonURL:   tritonBaseURL,
		AccountName: req.AccountName,
		Signers:     []authentication.Signer{signer},
	}

	return &AccountCheck{
		ParsedRequest: req,
		config:        config,
		store:         store,
	}
}

// newClient constructs our Triton AccountClient
func (ac *AccountCheck) newClient() (*account.AccountClient, error) {
	return account.NewClient(ac.config)
}

func (ac *AccountCheck) OnTriton(ctx context.Context) error {
	a, err := ac.newClient()
	if err != nil {
		return errors.Wrap(err, "failed to create account client")
	}

	a.SetHeader(ac.ParsedRequest.Header())

	acct, err := a.Get(ctx, &account.GetInput{})
	if err != nil {
		return errors.Wrap(err, "failed to get account")
	}

	ac.TritonAccount = acct

	return nil
}

func (ac *AccountCheck) createAccount(ctx context.Context) error {
	newAccount := accounts.New(ac.store)
	newAccount.AccountName = ac.TritonAccount.Login
	newAccount.TritonUUID = ac.TritonAccount.ID

	if err := newAccount.Insert(ctx); err != nil {
		return errors.Wrapf(err, "failed to insert %q account", newAccount.AccountName)
	}

	ac.Account = newAccount

	log.Debug().
		Str("id", fmt.Sprintf("%d", ac.Account.ID)).
		Str("name", ac.Account.AccountName).
		Str("uuid", ac.Account.TritonUUID).
		Msg("auth: inserted new account into database")

	return nil
}

// Save saves the TSG account from the Triton Account.
func (ac *AccountCheck) SaveAccount(ctx context.Context) error {
	var exists bool

	curAccount, err := ac.store.FindByName(ctx, ac.TritonAccount.Login)
	switch err {
	case nil:
		exists = true
	case pgx.ErrNoRows:
		exists = false
	default:
		return err
	}

	if !exists && isWhitelistOnly {
		log.Debug().
			Str("name", ac.TritonAccount.Login).
			Str("uuid", ac.TritonAccount.ID).
			Str("module", "whitelist").
			Msg("auth: access denied to new service users")

		return ErrWhitelist
	}

	if !exists {
		if err := ac.createAccount(ctx); err != nil {
			return err
		}

		return nil
	}

	if curAccount.TritonUUID != ac.TritonAccount.ID {
		curAccount.TritonUUID = ac.TritonAccount.ID

		if err := curAccount.Save(ctx); err != nil {
			return errors.Wrapf(err, "failed to save %q account", curAccount.AccountName)
		}
	}

	ac.Account = curAccount

	log.Debug().
		Str("id", fmt.Sprintf("%d", ac.Account.ID)).
		Str("name", ac.Account.AccountName).
		Str("uuid", ac.Account.TritonUUID).
		Msg("auth: found existing account in database")

	return nil
}

// HasTritonAccount returns a boolean whether or not we've authenticated with Triton.
func (ac *AccountCheck) HasTritonAccount() bool {
	return ac.TritonAccount != nil
}

// HasAccount returns a boolean whether or not the database has a valid Account.
func (ac *AccountCheck) HasAccount() bool {
	return ac.Account != nil
}

func (ac *AccountCheck) IsAuthentic() bool {
	return ac.HasTritonAccount() && ac.HasAccount()
}
