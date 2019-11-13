package tc

// OSVersionsResponse is the JSON representation of the
// OS versions for ISO generation.
type OSVersionsResponse struct {
	Versions map[string]string `json:"response"`
}
