package accounts

import (
	"context"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

var (
	ErrCredExists = errors.New("can't check credentials without key id and name")
	ErrExists     = errors.New("can't check existence without id or name")
	ErrMissingID  = errors.New("missing identifer for save")
)

// Account represents the data associated with an tsg_accounts row.
type Account struct {
	ID          string
	AccountName string
	TritonUUID  string
	KeyID       string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	store *Store
}

type TritonCredential struct {
	AccountName string
	KeyID       string
	KeyMaterial string
}

// New constructs a new Account with the Store for backend persistence.
func New(store *Store) *Account {
	return &Account{
		store: store,
	}
}

// Insert inserts a new account into the tsg_accounts table.
func (a *Account) Insert(ctx context.Context) error {
	query := `
INSERT INTO tsg_accounts (account_name, triton_uuid, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW());
`
	pool := a.store.pool

	tx, err := pool.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback() // nolint: errcheck

	_, err = pool.ExecEx(ctx, query, nil,
		a.AccountName,
		a.TritonUUID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert account")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	acct, err := a.store.FindByName(ctx, a.AccountName)
	if err != nil {
		return errors.Wrap(err, "failed to find account after insert")
	}

	a.ID = acct.ID
	a.CreatedAt = acct.CreatedAt
	a.UpdatedAt = acct.UpdatedAt

	return nil
}

// Save saves an accounts.Account object and it's field values.
func (a *Account) Save(ctx context.Context) error {
	if a.ID == "" {
		return ErrMissingID
	}

	updatedAt := time.Now()

	pool := a.store.pool

	tx, err := pool.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback() // nolint: errcheck

	if a.KeyID == "" {
		query := `
UPDATE tsg_accounts SET (account_name, triton_uuid, updated_at) = ($2, $3, $4)
WHERE id = $1;
`
		_, err := pool.ExecEx(ctx, query, nil,
			a.ID,
			a.AccountName,
			a.TritonUUID,
			updatedAt,
		)
		if err != nil {
			return errors.Wrap(err, "failed to save account with key")
		}
	} else {

		query := `
UPDATE tsg_accounts SET (account_name, triton_uuid, key_id, updated_at) = ($2, $3, $4, $5)
WHERE id = $1;
`
		_, err := pool.ExecEx(ctx, query, nil,
			a.ID,
			a.AccountName,
			a.TritonUUID,
			a.KeyID,
			updatedAt,
		)
		if err != nil {
			return errors.Wrap(err, "failed to save account with key")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	a.UpdatedAt = updatedAt

	return nil
}

// Exists returns a boolean and error. True if the row exists, false if it
// doesn't, error if there was an error executing the query.
func (a *Account) Exists(ctx context.Context) (bool, error) {
	if a.AccountName == "" && a.ID == "" {
		return false, ErrExists
	}

	var count int

	query := `
SELECT 1 FROM tsg_accounts
WHERE (id = $1 OR account_name = $2) AND archived = false;
`

	// NOTE(justinwr): seriously...
	accountID := "00000000-0000-0000-0000-000000000000"
	if a.ID != "" {
		accountID = a.ID
	}

	err := a.store.pool.QueryRowEx(ctx, query, nil,
		accountID,
		a.AccountName,
	).Scan(
		&count,
	)
	switch err {
	case nil:
		return true, nil
	case pgx.ErrNoRows:
		return false, nil
	default:
		return false, errors.Wrap(err, "failed to check account existence")
	}
}

// GetTritonCredential gets Triton credentials from an existing Account. If the
// account is found, then we will get the KeyID and KeyMaterial for the TSG
// Management key of that account. If we do not find any credentials, we return
// an error.
func (a *Account) GetTritonCredential(ctx context.Context) (*TritonCredential, error) {
	if a.AccountName == "" && a.KeyID == "" {
		return nil, ErrCredExists
	}

	var (
		fingerprint string
		material    string
	)

	query := `SELECT fingerprint, material FROM tsg_keys WHERE id = $1 AND archived = false;`

	err := a.store.pool.QueryRowEx(ctx, query, nil, a.KeyID).Scan(
		&fingerprint,
		&material,
	)
	if err != nil {
		return nil, err
	}

	return &TritonCredential{
		AccountName: a.AccountName,
		KeyID:       fingerprint,
		KeyMaterial: material,
	}, nil
}
