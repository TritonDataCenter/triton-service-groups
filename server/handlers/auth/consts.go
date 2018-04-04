package auth

const (
	matchName  = `^[a-zA-Z][a-zA-Z0-9_\.@]+$`
	matchKeyId = `keyId=\"(.*?)\"`

	defaultKeyName = "TSG_Management"

	// NOTE: if this is set to true than a triton account must be manually added
	// to the tsg_accounts table, auto account creation will be disabled
	isWhitelistOnly = true

	testAccountID   = "6f873d02-172c-418f-8416-4da2b50d5c53"
	testFingerprint = "5a:ce:1e:1d:b0:96:78:c6:7a:f2:f8:26:e1:b3:55:79"
)
