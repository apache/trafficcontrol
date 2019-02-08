package config

// Should mappings from error code to description replace the need for issue?
// No. An 'issue' can describe the error with more granularity.
type NegativeTest struct {
	Issue    string
	Config   string
	Expected int
}

// coverage should allow to to verify we could correctly make a config without an error
type PositiveTest struct {
	Coverage string // can be map correspond "dest_host" and integer? sounds like an enum
	Config   string
}
