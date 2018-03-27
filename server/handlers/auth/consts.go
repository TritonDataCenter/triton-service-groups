package auth

const (
	matchName  = `^[a-zA-Z][a-zA-Z0-9_\.@]+$`
	matchKeyId = `keyId=\"(.*?)\"`

	defaultKeyName = "TSG_Management"
	tritonBaseURL  = "https://us-west-1.api.joyent.com/"

	// NOTE: if this is set to true than a triton account must be manually added
	// to the tsg_accounts table, auto account creation will be disabled
	isWhitelistOnly = true

	testAccountID = "6f873d02-172c-418f-8416-4da2b50d5c53"
)