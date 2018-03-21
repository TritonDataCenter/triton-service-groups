package keys

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
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
func (s *Store) FindByID(ctx context.Context, keyID int64) (*Key, error) {
	var (
		id          int64
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
	key.ID = id
	key.Name = name
	key.Fingerprint = fingerprint
	key.Material = material
	key.CreatedAt = createdAt.Time
	key.UpdatedAt = updatedAt.Time

	return key, nil
}

// FindByName finds an account by a specific account_name.
func (s *Store) FindByName(ctx context.Context, keyName string) (*Key, error) {
	var (
		id          int64
		name        string
		fingerprint string
		material    string
		createdAt   pgtype.Timestamp
		updatedAt   pgtype.Timestamp
	)

	query := `
SELECT id, name, fingerprint, material, created_at, updated_at
FROM tsg_keys
WHERE name = $1 AND archived = false;
`
	err := s.pool.QueryRowEx(ctx, query, nil, keyName).Scan(
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
	key.ID = id
	key.Name = name
	key.Fingerprint = fingerprint
	key.Material = material
	key.CreatedAt = createdAt.Time
	key.UpdatedAt = updatedAt.Time

	return key, nil
}
