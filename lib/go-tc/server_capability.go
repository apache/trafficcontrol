package tc

// ServerServerCapability represents an association between a server capability and a server
type ServerServerCapability struct {
	LastUpdated      *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Server           *string    `json:"serverHostName,omitempty" db:"host_name"`
	ServerID         *int       `json:"serverId" db:"server"`
	ServerCapability *string    `json:"serverCapability" db:"server_capability"`
}
