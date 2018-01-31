//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package router

import (
	"net/http"

	"github.com/joyent/triton-service-groups/session"
)

func isAuthenticated(session *session.TsgSession, r *http.Request) bool {
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
