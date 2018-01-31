package session

import "github.com/jackc/pgx"

type TsgSession struct {
	AccountId string
	DbPool    *pgx.ConnPool
}
