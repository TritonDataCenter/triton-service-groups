package auth

// authSession a private struct which is only accessible by pulling out of the
// current request `context.Context`.
type Session struct {
	AccountID string
}

// IsAuthenticated encapsulates whatever it means for an authSession to be
// deemed authenticated.
func (a Session) IsAuthenticated() bool {
	return a.AccountID != ""
}
