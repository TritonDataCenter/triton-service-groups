package auth

type Config struct {
	// Name of the datacenter in which this TSG service is operating. This is
	// used to create unique key names per-DC. The value is also available in
	// the HTTP request Session object.
	Datacenter string

	// URL of Triton's CloudAPI in which to scale instances. This is made
	// available within the HTTP request Session object.
	TritonURL string

	// URL of Triton's CloudAPI in which to authenticate incoming API
	// requests. This is only used by internal auth processes. It can be set to
	// the same CloudAPI used by TritonURL as well.
	AuthURL string

	// Prefix name used when creating a new key in Triton. This defaults to
	// "TSG_Management" but can be configured with whatever an end user
	// prefers. The current Datacenter is also appended to this value at
	// runtime.
	KeyNamePrefix string

	// Enable or disable whitelisting behavior. This feature only accepts
	// requests from user accounts that have previously been authenticated. If
	// this is set to true than a Triton account must be manually added to the
	// tsg_accounts table, auto account creation will be disabled.
	EnableWhitelist bool
}
