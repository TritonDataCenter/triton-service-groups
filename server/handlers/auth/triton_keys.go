package auth

import (
	"context"
	"strings"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
	"github.com/pkg/errors"
)

type foundKeys struct {
	vault   map[string]bool
	account map[string]bool
}

type TritonKeys struct {
	*parsedRequest

	found  foundKeys
	config *triton.ClientConfig
}

func NewTritonKeys(req *parsedRequest) *TritonKeys {
	signer := &authentication.TestSigner{}
	config := &triton.ClientConfig{
		TritonURL:   "https://us-sw-1.api.joyent.com/",
		AccountName: req.accountName,
		Signers:     []authentication.Signer{signer},
	}

	return &TritonKeys{
		parsedRequest: req,
		config:        config,
		found: foundKeys{
			account: make(map[string]bool, 0),
		},
	}
}

func (t *TritonKeys) newClient() (*account.AccountClient, error) {
	return account.NewClient(t.config)
}

func (t *TritonKeys) Check() error {
	a, err := t.newClient()
	if err != nil {
		return errors.Wrap(err, "failed to create account keys client")
	}

	a.SetHeader(t.parsedRequest.getHeader())

	listInput := &account.ListKeysInput{}
	keys, err := a.Keys().List(context.Background(), listInput)
	if err != nil {
		return errors.Wrap(err, "failed to list account keys")
	}
	for _, key := range keys {
		if strings.HasPrefix(key.Name, "tsg-") {
			t.found.account[key.Fingerprint] = true
		}
	}
	return nil
}

func (t *TritonKeys) HasKey() bool {
	return len(t.found.account) == 1
}

func (t *TritonKeys) HasKeys() bool {
	return len(t.found.account) >= 1
}
