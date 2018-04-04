package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/jackc/pgx"
	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
	terrors "github.com/joyent/triton-go/errors"
	"github.com/joyent/triton-service-groups/accounts"
	"github.com/joyent/triton-service-groups/keys"
	"github.com/pkg/errors"
)

type KeyCheck struct {
	*ParsedRequest

	Key       *keys.Key
	TritonKey *account.Key

	config  *triton.ClientConfig
	store   *keys.Store
	account *accounts.Account
	dc      string
}

func NewKeyCheck(req *ParsedRequest, acct *accounts.Account, store *keys.Store, dc string) *KeyCheck {
	signer := &authentication.TestSigner{}
	config := &triton.ClientConfig{
		TritonURL:   tritonBaseURL,
		AccountName: req.AccountName,
		Signers:     []authentication.Signer{signer},
	}

	return &KeyCheck{
		ParsedRequest: req,
		account:       acct,
		config:        config,
		store:         store,
		dc:            dc,
	}
}

// newClient constructs our Triton AccountClient
func (k *KeyCheck) newClient() (*account.AccountClient, error) {
	return account.NewClient(k.config)
}

// CheckTriton checks Triton account keys for our TSG key
func (k *KeyCheck) OnTriton(ctx context.Context) error {
	a, err := k.newClient()
	if err != nil {
		return errors.Wrap(err, "failed to create account key client")
	}

	a.SetHeader(k.ParsedRequest.Header())

	input := &account.GetKeyInput{
		KeyName: keyNameForDC(k.dc),
	}
	key, err := a.Keys().Get(ctx, input)
	if err != nil {
		if terrors.IsSpecificStatusCode(err, http.StatusNotFound) {
			return nil
		}
		return errors.Wrap(err, "failed to get triton key")
	}

	k.TritonKey = key

	return nil
}

// InDatabase checks for and sets an account's key within the TSG database.
func (k *KeyCheck) InDatabase(ctx context.Context) error {
	if k.account.KeyID == "" {
		return nil
	}

	curKey, err := k.store.FindByID(ctx, k.account.KeyID)
	switch err {
	case nil:
		k.Key = curKey
		return nil
	case pgx.ErrNoRows:
		return nil
	default:
		return err
	}
}

// AddKey adds an account key into Triton, converting the passed in KeyPair into
// a Triton-Go account.Key for use by external consumers.
func (k *KeyCheck) AddTritonKey(ctx context.Context, keypair *KeyPair) error {
	a, err := k.newClient()
	if err != nil {
		return errors.Wrap(err, "failed to create new key client")
	}

	a.SetHeader(k.ParsedRequest.Header())

	createInput := &account.CreateKeyInput{
		Name: keyNameForDC(k.dc),
		Key:  keypair.PublicKeyBase64(),
	}
	key, err := a.Keys().Create(ctx, createInput)
	if err != nil {
		return errors.Wrap(err, "failed to create new account key")
	}

	k.TritonKey = key

	return nil
}

func (k *KeyCheck) InsertKey(ctx context.Context, keypair *KeyPair) error {
	key := keys.New(k.store)

	key.Name = keyNameForDC(k.dc)
	key.Fingerprint = keypair.FingerprintMD5
	key.Material = keypair.PrivateKeyPEM()
	key.AccountID = k.account.ID

	if err := key.Insert(ctx); err != nil {
		return errors.Wrap(err, "failed to store account key")
	}

	k.account.KeyID = key.ID
	if err := k.account.Save(ctx); err != nil {
		return errors.Wrap(err, "failed to store account key_id")
	}

	k.Key = key

	return nil
}

func (k *KeyCheck) HasTritonKey() bool {
	return k.TritonKey != nil
}

func (k *KeyCheck) HasKey() bool {
	return k.Key != nil
}

func keyNameForDC(dc string) string {
	return strings.Join([]string{defaultKeyName, dc}, "_")
}
