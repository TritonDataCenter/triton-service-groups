package keys

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/joyent/triton-service-groups/convert"
)

type Store struct {
	pool *pgx.ConnPool
}

// NewStore returns a new store object.
func NewStore(pool *pgx.ConnPool) *Store {
	return &Store{
		pool: pool,
	}
}

// FindByID finds an account by a specific ID.
func (s *Store) FindByID(ctx context.Context, keyID string) (*Key, error) {
	var (
		id          pgtype.UUID
		name        string
		fingerprint string
		material    string
		createdAt   pgtype.Timestamp
		updatedAt   pgtype.Timestamp
	)

	query := `
SELECT id, name, fingerprint, material, created_at, updated_at
FROM tsg_keys
WHERE id = $1 AND archived = false;
`
	err := s.pool.QueryRowEx(ctx, query, nil, keyID).Scan(
		&id,
		&name,
		&fingerprint,
		&material,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	key := New(s)
	key.ID = convert.BytesToUUID(id.Bytes)
	key.Name = name
	key.Fingerprint = fingerprint
	key.Material = material
	key.CreatedAt = createdAt.Time
	key.UpdatedAt = updatedAt.Time

	return key, nil
}

// FindByName finds an account by a specific account_name.
func (s *Store) FindByName(ctx context.Context, keyName string, accountID string) (*Key, error) {
	var (
		id          pgtype.UUID
		name        string
		fingerprint string
		material    string
		createdAt   pgtype.Timestamp
		updatedAt   pgtype.Timestamp
	)

	query := `
SELECT id, name, fingerprint, material, created_at, updated_at
FROM tsg_keys
WHERE name = $1 AND account_id = $2 AND archived = false;
`
	err := s.pool.QueryRowEx(ctx, query, nil, keyName, accountID).Scan(
		&id,
		&name,
		&fingerprint,
		&material,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	key := New(s)
	key.ID = convert.BytesToUUID(id.Bytes)
	key.AccountID = accountID
	key.Name = name
	key.Fingerprint = fingerprint
	key.Material = material
	key.CreatedAt = createdAt.Time
	key.UpdatedAt = updatedAt.Time

	return key, nil
}
