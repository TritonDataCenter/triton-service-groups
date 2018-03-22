package auth

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

type ParsedRequest struct {
	authHeader string
	dateHeader string

	AccountName string
	UserName    string
	Fingerprint string
}

func parseKeyId(keyId string) (string, string, string, error) {
	matches := strings.Split(keyId, `/`)

	if strings.Contains(keyId, `/users/`) {
		if len(matches) != 6 {
			return "", "", "", ErrParseAuth
		}

		return matches[1], matches[3], matches[5], nil
	}

	if len(matches) != 4 {
		return "", "", "", ErrParseAuth
	}

	return matches[1], "", matches[3], nil
}

func validateName(name string) error {
	if len(name) < 3 {
		return ErrNameLen
	}

	matched, err := regexp.MatchString(matchName, name)
	if err != nil {
		log.Error().Err(err)
		return ErrNameFormat
	}
	if !matched {
		return ErrNameFormat
	}

	return nil
}

func ParseRequest(req *http.Request) (*ParsedRequest, error) {
	dateHeader := req.Header.Get("Date")
	authHeader := req.Header.Get("Authorization")

	if dateHeader == "" || authHeader == "" {
		return nil, ErrUnauthRequest
	}

	re, err := regexp.Compile(matchKeyId)
	if err != nil {
		return nil, err
	}

	matches := re.FindStringSubmatch(fmt.Sprintf("%s", authHeader))
	if len(matches) != 2 {
		return nil, ErrBadKeyID
	}

	accountName, userName, fingerprint, err := parseKeyId(matches[1])
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if accountName == "" || fingerprint == "" {
		return nil, ErrParseValue
	}

	if err := validateName(accountName); err != nil {
		return nil, err
	}
	if userName != "" {
		if err := validateName(userName); err != nil {
			return nil, err
		}
	}

	return &ParsedRequest{
		dateHeader:  dateHeader,
		authHeader:  authHeader,
		AccountName: accountName,
		UserName:    userName,
		Fingerprint: fingerprint,
	}, nil
}

func (r *ParsedRequest) Header() *http.Header {
	header := &http.Header{}
	header.Set("date", r.dateHeader)
	header.Set("Authorization", r.authHeader)

	return header
}
