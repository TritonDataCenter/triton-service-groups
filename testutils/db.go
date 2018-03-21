package testutils

import (
	"testing"

	"github.com/jackc/pgx"
)

type TestDB struct {
	Conn *pgx.ConnPool
}

// NewTestDB creates a new object which is used to act upon the database within
// tests.
func NewTestDB() (*TestDB, error) {
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		MaxConnections: 5,
		AfterConnect:   nil,
		AcquireTimeout: 0,
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Database: "triton_test",
			Port:     26257,
			User:     "root",
		},
	})
	if err != nil {
		return nil, err
	}

	return &TestDB{connPool}, nil
}

// Clear clears out all active tables used during automated testing.
func (db *TestDB) Clear(t *testing.T) {
	_, err := db.Conn.Exec(`DELETE FROM tsg_groups`)
	if err != nil {
		t.Fatalf("conn.Exec failed: %v", err)
	}

	_, err2 := db.Conn.Exec(`DELETE FROM tsg_templates`)
	if err2 != nil {
		t.Fatalf("conn.Exec failed: %v", err)
	}

	_, err3 := db.Conn.Exec(`DELETE FROM tsg_users`)
	if err3 != nil {
		t.Fatalf("conn.Exec failed: %v", err2)
	}

	_, err4 := db.Conn.Exec(`DELETE FROM tsg_accounts`)
	if err4 != nil {
		t.Fatalf("conn.Exec failed: %v", err2)
	}

	_, err5 := db.Conn.Exec(`DELETE FROM tsg_keys`)
	if err5 != nil {
		t.Fatalf("conn.Exec failed: %v", err2)
	}
}
