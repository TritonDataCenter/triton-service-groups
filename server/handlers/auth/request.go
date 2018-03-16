package auth

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type parsedRequest struct {
	authHeader  string
	dateHeader  string
	accountName string
	fingerprint string
}

var (
	ErrUnauthRequest = errors.New("received unauthenticated request")
	ErrBadKeyID      = errors.New("couldn't parse keyId within authorization header")
	ErrParseAuth     = errors.New("bad values parsed from keyId header")
)

func parseRequest(req *http.Request) (*parsedRequest, error) {
	dateHeader := req.Header.Get("Date")
	authHeader := req.Header.Get("Authorization")

	if dateHeader == "" || authHeader == "" {
		return nil, ErrUnauthRequest
	}

	re, err := regexp.Compile("keyId=\"(.*?)\"")
	if err != nil {
		return nil, err
	}

	matches := re.FindStringSubmatch(fmt.Sprintf("%s", authHeader))
	if len(matches) != 2 {
		return nil, ErrBadKeyID
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

	if accountName == "" || fingerprint == "" {
		return nil, ErrParseAuth
	}

	return &parsedRequest{
		dateHeader:  dateHeader,
		authHeader:  authHeader,
		accountName: accountName,
		fingerprint: fingerprint,
	}, nil
}

func (r *parsedRequest) hasValues() bool {
	return r.dateHeader != "" &&
		r.authHeader != "" &&
		r.accountName != "" &&
		r.fingerprint != ""
}

func (r *parsedRequest) getHeader() *http.Header {
	header := &http.Header{}
	header.Set("date", r.dateHeader)
	header.Set("Authorization", r.authHeader)

	return header
}
