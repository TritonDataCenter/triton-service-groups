package keys

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

// Key represents the data associated with an tsg_keys row.
type Key struct {
	ID          string
	Name        string
	Fingerprint string
	Material    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Archived    bool

	store *Store
}

// New constructs a new Key with the Store for backend persistence.
func New(store *Store) *Key {
	return &Key{
		store: store,
	}
}

// Insert inserts a new key into the tsg_keys table.
func (k *Key) Insert(ctx context.Context) error {
	query := `
INSERT INTO tsg_keys (name, fingerprint, material, archived, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW());
`
	_, err := k.store.pool.ExecEx(ctx, query, nil,
		k.Name,
		k.Fingerprint,
		k.Material,
		k.Archived,
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert key")
	}

	key, err := k.store.FindByName(ctx, k.Name)
	if err != nil {
		return errors.Wrap(err, "failed to find key after insert")
	}

	k.ID = key.ID
	k.CreatedAt = key.CreatedAt
	k.UpdatedAt = key.UpdatedAt

	return nil
}

// Save saves an keys.Key object and it's field values.
func (k *Key) Save(ctx context.Context) error {
	if k.ID == "" {
		return ErrMissingID
	}

	query := `
UPDATE tsg_keys SET (name, fingerprint, material, archived, updated_at) = ($2, $3, $4, $5, $6)
WHERE id = $1;
`
	updatedAt := time.Now()

	_, err := k.store.pool.ExecEx(ctx, query, nil,
		k.ID,
		k.Name,
		k.Fingerprint,
		k.Material,
		k.Archived,
		updatedAt,
	)
	if err != nil {
		return err
	}

	k.UpdatedAt = updatedAt

	return nil
}

// Exists returns a boolean and error. True if the row exists, false if it
// doesn't, error if there was an error executing the query.
func (k *Key) Exists(ctx context.Context) (bool, error) {
	if k.Name == "" && k.ID == "" {
		return false, ErrExists
	}

	var count int

	query := `
SELECT 1 FROM tsg_keys
WHERE (id = $1 OR name = $2) AND archived = false;
`
	err := k.store.pool.QueryRowEx(ctx, query, nil,
		k.ID,
		k.Name,
	).Scan(&count)
	switch err {
	case nil:
		return true, nil
	case pgx.ErrNoRows:
		return false, nil
	default:
		return false, errors.Wrap(err, "failed to check key existence")
	}
}
