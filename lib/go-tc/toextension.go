package tc

import "encoding/json"

// TOExtensionNullable ...
type TOExtensionNullable struct {
	ID                   *int            `json:"id"`
	Name                 *string         `json:"name"`
	Version              *string         `json:"version"`
	InfoURL              *string         `json:"info_url"`
	ScriptFile           *string         `json:"script_file"`
	IsActive             *bool           `json:"isactive"`
	AdditionConfigJSON   json.RawMessage `json:"additional_config_json"`
	Description          *string         `json:"description"`
	ServercheckShortName *string         `json:"servercheck_short_name"`
	Type                 *string         `json:"type"`
}

type TOExtensionResponse struct {
	Response []TOExtensionNullable `json:"response"`
}
