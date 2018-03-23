package accounts

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
func (s *Store) FindByID(ctx context.Context, accountID int64) (*Account, error) {
	var (
		id        pgtype.UUID
		keyID     pgtype.UUID
		name      string
		uuid      string
		createdAt pgtype.Timestamp
		updatedAt pgtype.Timestamp
	)

	query := `
SELECT id, account_name, triton_uuid, key_id, created_at, updated_at
FROM tsg_accounts
WHERE id = $1 AND archived = false;
`
	err := s.pool.QueryRowEx(ctx, query, nil, accountID).Scan(
		&id,
		&name,
		&uuid,
		&keyID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	acct := New(s)
	acct.ID = convert.BytesToUUID(id.Bytes)
	acct.AccountName = name
	acct.TritonUUID = uuid
	acct.KeyID = convert.BytesToUUID(keyID.Bytes)
	acct.CreatedAt = createdAt.Time
	acct.UpdatedAt = updatedAt.Time

	return acct, nil
}

// FindByName finds an account by a specific account_name.
func (s *Store) FindByName(ctx context.Context, accountName string) (*Account, error) {
	var (
		id        pgtype.UUID
		keyID     pgtype.UUID
		name      string
		uuid      string
		createdAt pgtype.Timestamp
		updatedAt pgtype.Timestamp
	)

	query := `
SELECT id, account_name, triton_uuid, key_id, created_at, updated_at
FROM tsg_accounts
WHERE account_name = $1 AND archived = false;
`
	err := s.pool.QueryRowEx(ctx, query, nil, accountName).Scan(
		&id,
		&name,
		&uuid,
		&keyID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	acct := New(s)
	acct.ID = convert.BytesToUUID(id.Bytes)
	acct.AccountName = name
	acct.TritonUUID = uuid
	acct.KeyID = convert.BytesToUUID(keyID.Bytes)
	acct.CreatedAt = createdAt.Time
	acct.UpdatedAt = updatedAt.Time

	return acct, nil
}
