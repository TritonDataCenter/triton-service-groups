//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
