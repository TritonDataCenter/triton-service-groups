package auth

import (
	"context"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-service-groups/accounts"
	"github.com/pkg/errors"
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

// Save saves the TSG account from the Triton Account.
func (ac *AccountCheck) SaveAccount(ctx context.Context) error {
	a := accounts.New(ac.store)
	a.AccountName = ac.TritonAccount.Login
	a.TritonUUID = ac.TritonAccount.ID

	if err := a.Save(ctx); err != nil {
		return err
	}

	return nil
}

// HasAccount returns a boolean whether or not we've authenticated with Triton.
func (ac *AccountCheck) HasAccount() bool {
	return ac.TritonAccount != nil
}
