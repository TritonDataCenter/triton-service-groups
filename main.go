package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/jackc/pgx"
	tsgRouter "github.com/joyent/triton-service-groups/router"
	"github.com/joyent/triton-service-groups/session"
)

func main() {
	dbPool, err := initDb()
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	// Session is a global application struct
	// This allows us to share information
	// between handlers easily
	session := &session.TsgSession{
		DbPool: dbPool,
	}

	router := tsgRouter.MakeRouter(session)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	authenticatedRouter := tsgRouter.AuthenticationHandler(session, loggedRouter)

	log.Fatal(http.ListenAndServe(":3000", authenticatedRouter))
}

func initDb() (*pgx.ConnPool, error) {
	dbPoolConfig := pgx.ConnPoolConfig{
		MaxConnections: 5,
		AfterConnect:   nil,
		AcquireTimeout: 0,
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Database: "tsg",
			Port:     26257,
			User:     "root",
		},
	}

	connPool, err := pgx.NewConnPool(dbPoolConfig)
	if err != nil {
		return nil, err
	}

	return connPool, nil
}
