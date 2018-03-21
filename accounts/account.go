package accounts

import (
	"context"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

var (
	ErrExists    = errors.New("can't check existence without id or name")
	ErrMissingID = errors.New("missing identifer for save")
)

// Account represents the data associated with an tsg_accounts row.
type Account struct {
	ID          int64
	AccountName string
	TritonUUID  string
	KeyID       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time

	store *Store
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
	_, err := a.store.pool.ExecEx(ctx, query, nil,
		a.AccountName,
		a.TritonUUID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert account")
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
	if a.ID == 0 {
		return ErrMissingID
	}

	query := `
UPDATE tsg_accounts SET (account_name, triton_uuid, updated_at) = ($2, $3, $4)
WHERE id = $1;
`
	updatedAt := time.Now()

	_, err := a.store.pool.ExecEx(ctx, query, nil,
		a.ID,
		a.AccountName,
		a.TritonUUID,
		updatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to save account")
	}

	a.UpdatedAt = updatedAt

	return nil
}

// Exists returns a boolean and error. True if the row exists, false if it
// doesn't, error if there was an error executing the query.
func (a *Account) Exists(ctx context.Context) (bool, error) {
	if a.AccountName == "" && a.ID == 0 {
		return false, ErrExists
	}

	var count int

	query := `
SELECT 1 FROM tsg_accounts
WHERE (id = $1 OR account_name = $2) AND archived = false;
`
	err := a.store.pool.QueryRowEx(ctx, query, nil,
		a.ID,
		a.AccountName,
	).Scan(&count)
	switch err {
	case nil:
		return true, nil
	case pgx.ErrNoRows:
		return false, nil
	default:
		return false, errors.Wrap(err, "failed to check account existence")
	}

	return true, nil
}
