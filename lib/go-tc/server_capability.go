package tc

// ServerServerCapabilityNullable
type ServerServerCapabilityNullable struct {
	LastUpdated      *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Server           *string    `json:"serverHostName,omitempty" db:"host_name"`
	ServerID         *int       `json:"serverId" db:"server"`
	ServerCapability *string    `json:"serverCapability" db:"server_capability"`
}
