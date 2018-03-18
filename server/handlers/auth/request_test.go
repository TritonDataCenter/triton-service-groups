package auth_test

import (
	"net/http"
	"testing"

	"github.com/joyent/triton-service-groups/server/handlers/auth"
	"github.com/stretchr/testify/assert"
)

const (
	authHeader    = "Signature keyId=\"/testaccount/keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	badAuthHeader = "Signature keyId=\"//keys/12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01\",algorithm=\"rsa-sha1\",headers=\"date\",signature=\"AABBCCDDEEFFGG\""
	dateHeader    = "Sat, 17 Mar 2018 16:12:06 UTC"
)

func TestParseRequest(t *testing.T) {
	parsedReq := &auth.ParsedRequest{
		AccountName: "testaccount",
		Fingerprint: "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01",
	}

	req := &http.Request{}
	h := http.Header{}
	h.Set("Authorization", authHeader)
	h.Set("Date", dateHeader)
	req.Header = h

	noDateReq := &http.Request{}
	h = http.Header{}
	h.Set("Authorization", authHeader)
	noDateReq.Header = h

	noAuthReq := &http.Request{}
	h = http.Header{}
	h.Set("Date", dateHeader)
	noAuthReq.Header = h

	badKeyReq := &http.Request{}
	h = http.Header{}
	h.Set("Authorization", "failed parse")
	h.Set("Date", dateHeader)
	badKeyReq.Header = h

	badNameReq := &http.Request{}
	h = http.Header{}
	h.Set("Authorization", badAuthHeader)
	h.Set("Date", dateHeader)
	badNameReq.Header = h

	tests := []struct {
		name   string
		input  *http.Request
		output *auth.ParsedRequest
		err    error
	}{
		{"successful", req, parsedReq, nil},
		{"missing date", noDateReq, nil, auth.ErrUnauthRequest},
		{"missing auth", noAuthReq, nil, auth.ErrUnauthRequest},
		{"bad auth header parse", badKeyReq, nil, auth.ErrBadKeyID},
		{"bad account name parse", badNameReq, nil, auth.ErrParseAuth},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := auth.ParseRequest(test.input)
			if test.err != nil {
				assert.Equal(t, test.err, err)
				return
			}
			if test.output != nil {
				assert.Equal(t, test.output.AccountName, output.AccountName)
				assert.Equal(t, test.output.Fingerprint, output.Fingerprint)
			}
		})
	}
}
