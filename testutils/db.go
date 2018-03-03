package testutils

import (
	"testing"
	"time"

	"github.com/jackc/pgx"
)

type TestDB struct {
	conn *pgx.ConnPool
}

func NewTestDB(conn *pgx.ConnPool) *TestDB {
	return &TestDB{conn}
}

func (db *TestDB) Clear(t *testing.T) {
	rows, err := db.conn.Query("TRUNCATE tsg_groups")
	if err != nil {
		t.Fatalf("conn.Query failed: %v", err)
	}
	defer rows.Close()

	rows2, err2 := db.conn.Query("TRUNCATE tsg_templates CASCADE")
	if err2 != nil {
		t.Fatalf("conn.Query failed: %v", err2)
	}
	defer rows2.Close()

	time.Sleep(4 * time.Second)
}

// TODO: We should refactor how/where our database initializes so we can half
// bootstrap the application from our tests with a simple one-liner.
func InitDB() (*pgx.ConnPool, error) {
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

	return connPool, nil
}
