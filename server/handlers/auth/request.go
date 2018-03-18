package auth

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type ParsedRequest struct {
	authHeader string
	dateHeader string

	AccountName string
	Fingerprint string
}

var (
	ErrUnauthRequest = errors.New("received unauthenticated request")
	ErrBadKeyID      = errors.New("couldn't parse keyId within header")
	ErrParseAuth     = errors.New("failed to parse values from keyId")
	ErrParseValue    = errors.New("incorrect values parsed from keyId")
)

func ParseRequest(req *http.Request) (*ParsedRequest, error) {
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

	if len(parts) < 2 {
		return nil, ErrParseAuth
	}

	accountName := parts[0]
	fingerprint := parts[1]

	if accountName == "" || fingerprint == "" {
		return nil, ErrParseValue
	}

	return &ParsedRequest{
		dateHeader:  dateHeader,
		authHeader:  authHeader,
		AccountName: accountName,
		Fingerprint: fingerprint,
	}, nil
}

func (r *ParsedRequest) HasValues() bool {
	return r.dateHeader != "" &&
		r.authHeader != "" &&
		r.AccountName != "" &&
		r.Fingerprint != ""
}

func (r *ParsedRequest) Header() *http.Header {
	header := &http.Header{}
	header.Set("date", r.dateHeader)
	header.Set("Authorization", r.authHeader)

	return header
}
