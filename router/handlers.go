//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/account"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-service-groups/session"
)

func isAuthenticated(session *session.TsgSession, req *http.Request) bool {
	dateHeader := req.Header.Get("Date")
	authHeader := req.Header.Get("Authorization")

	if dateHeader == "" || authHeader == "" {
		session.AccountId = ""
		return false
	}

	re, err := regexp.Compile("keyId=\"(.*?)\"")
	if err != nil {
		return false
	}

	matches := re.FindStringSubmatch(fmt.Sprintf("%s", authHeader))
	if len(matches) != 2 {
		// fmt.Error("couldn't find keyId within authorization header")
		return false
	}

	authParts := strings.Split(matches[1], "/")
	parts := []string{}
	for _, part := range authParts {
		if part != "" && part != "keys" {
			parts = append(parts, part)
		}
	}

	accountName := parts[0]
	fingerprint := parts[1]
	signer := &authentication.TestSigner{}

	config := &triton.ClientConfig{
		TritonURL:   "https://us-east-1.api.joyent.com/",
		AccountName: accountName,
		Signers:     []authentication.Signer{signer},
	}

	a, err := account.NewClient(config)
	if err != nil {
		log.Println("failed to create account client: %v", err)
	}

	header := &http.Header{}
	header.Set("date", dateHeader)
	header.Set("Authorization", authHeader)
	a.SetHeader(header)

	input := &account.ListKeysInput{}
	keys, err := a.Keys().List(context.Background(), input)
	if err != nil {
		log.Println("failed to list account keys: %v", err)
		return false
	}
	for _, key := range keys {
		fmt.Println("Key Name", key.Name)
	}

	fmt.Println(fingerprint)

	session.AccountId = "joyent"
	return true
}

func AuthenticationHandler(session *session.TsgSession, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(session, r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
