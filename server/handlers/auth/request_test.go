package auth_test

import (
	"net/http"
	"testing"

	"github.com/joyent/triton-service-groups/server/handlers/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	authHeader           = "Signature keyId=\"/testaccount/keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	authUserHeader       = "Signature keyId=\"/testaccount/users/demouser/keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	shortAccountHeader   = "Signature keyId=\"/te/keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	shortUserHeader      = "Signature keyId=\"/testaccount/users/te/keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	invalidAccountHeader = "Signature keyId=\"/test+account/keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	invalidUserHeader    = "Signature keyId=\"/testaccount/users/demo?user/keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	badAuthHeader        = "Signature keyId=\"//keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	badFPrintHeader      = "Signature keyId=\"/testaccount/keys/\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	dateHeader           = "Sat, 17 Mar 2018 16:12:06 UTC"
)

func newAuthRequest(authHeader string, dateHeader string) *http.Request {
	req := &http.Request{}
	h := http.Header{}
	h.Set("Authorization", authHeader)
	h.Set("Date", dateHeader)
	req.Header = h

	return req
}

func TestParseRequest(t *testing.T) {
	parsedReq := &auth.ParsedRequest{
		AccountName: "testaccount",
		Fingerprint: "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01",
	}

	parsedUserReq := &auth.ParsedRequest{
		AccountName: "testaccount",
		UserName:    "demouser",
		Fingerprint: "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01",
	}

	tests := []struct {
		name   string
		auth   string
		date   string
		output *auth.ParsedRequest
		err    error
	}{
		{"successful", authHeader, dateHeader, parsedReq, nil},
		{"with subuser", authUserHeader, dateHeader, parsedUserReq, nil},
		{"missing date", authHeader, "", nil, auth.ErrUnauthRequest},
		{"missing auth", "", dateHeader, nil, auth.ErrUnauthRequest},
		{"bad auth header parse", "failed parse", dateHeader, nil, auth.ErrBadKeyID},
		{"bad account name parse", badAuthHeader, dateHeader, nil, auth.ErrParseValue},
		{"bad fingerprint parse", badFPrintHeader, dateHeader, nil, auth.ErrParseValue},
		{"short account name", shortAccountHeader, dateHeader, nil, auth.ErrNameLen},
		{"short user name", shortUserHeader, dateHeader, nil, auth.ErrNameLen},
		{"invalid account name", invalidAccountHeader, dateHeader, nil, auth.ErrNameFormat},
		{"invalid user name", invalidUserHeader, dateHeader, nil, auth.ErrNameFormat},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := newAuthRequest(test.auth, test.date)
			output, err := auth.ParseRequest(req)

			if test.err != nil {
				require.Error(t, err)

				assert.Equal(t, test.err, err)
				return
			}

			require.NotNil(t, output)

			if test.output != nil {
				assert.Equal(t, test.output.AccountName, output.AccountName)
				assert.Equal(t, test.output.Fingerprint, output.Fingerprint)

				if test.output.UserName != "" {
					assert.Equal(t, test.output.UserName, output.UserName)
				}
			}
		})
	}
}
