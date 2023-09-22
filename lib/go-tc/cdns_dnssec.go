package tc

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Types of keys used in DNSSEC and stored in Traffic Vault.
const (
	// Key-signing key.
	DNSSECKSKType = "ksk"
	// Zone-signing key.
	DNSSECZSKType = "zsk"
)

// The various allowable statuses for DNSSEC keys being inserted into or
// retrieved from Traffic Vault.
const (
	DNSSECKeyStatusNew     = "new"
	DNSSECKeyStatusExpired = "expired"
	DNSSECStatusExisting   = "existing"
)

// CDNDNSSECKeysResponse is the type of a response from Traffic Ops to GET
// requests made to its /cdns/name/{{name}}/dnsseckeys API endpoint.
type CDNDNSSECKeysResponse struct {
	Response DNSSECKeys `json:"response"`
	Alerts
}

// GenerateCDNDNSSECKeysResponse is the type of a response from Traffic Ops to
// requests made to its /cdns/dnsseckeys/generate and
// /cdns/dnsseckeys/ksk/generate API endpoints.
type GenerateCDNDNSSECKeysResponse struct {
	Response string `json:"response"`
	Alerts
}

// DeleteCDNDNSSECKeysResponse is the type of a response from Traffic Ops to
// DELETE requests made to its /cdns/name/{{name}}/dnsseckeys API endpoint.
type DeleteCDNDNSSECKeysResponse GenerateCDNDNSSECKeysResponse

// RefreshDNSSECKeysResponse is the type of a response from Traffic Ops to
// requests made to its /cdns/dnsseckeys/refresh API endpoint.
type RefreshDNSSECKeysResponse GenerateCDNDNSSECKeysResponse

// DNSSECKeys is the DNSSEC keys as stored in Traffic Vault, plus the DS record text.
type DNSSECKeys map[string]DNSSECKeySet

// DNSSECKeysRiak is the structure in which DNSSECKeys are stored in the Riak
// backend for Traffic Vault.
//
// Deprecated: use DNSSECKeysTrafficVault instead.
type DNSSECKeysRiak DNSSECKeysV11

// A DNSSECKeysTrafficVault is a mapping of CDN Names and/or Delivery Service
// XMLIDs to sets of keys used for DNSSEC with that CDN or Delivery Service.
type DNSSECKeysTrafficVault DNSSECKeysV11

// DNSSECKeysV11 is the DNSSEC keys object stored in Traffic Vault. The map key
// strings are both DeliveryServiceNames and CDNNames.
type DNSSECKeysV11 map[string]DNSSECKeySetV11

// A DNSSECKeySet is a set of keys used for DNSSEC zone and key signing.
type DNSSECKeySet struct {
	ZSK []DNSSECKey `json:"zsk"`
	KSK []DNSSECKey `json:"ksk"`
}

// DNSSECKeySetV11 is a DNSSEC key set (ZSK and KSK), as stored in Traffic
// Vault.
// This is specifically the key data, without the DS record text (which can be computed), and was also the format used in API 1.1 through 1.3.
type DNSSECKeySetV11 struct {
	ZSK []DNSSECKeyV11 `json:"zsk"`
	KSK []DNSSECKeyV11 `json:"ksk"`
}

// A DNSSECKey is a DNSSEC Key (Key-Signing or Zone-Signing) and all associated
// data - computed data as well as data that is actually stored by Traffic
// Vault.
type DNSSECKey struct {
	DNSSECKeyV11
	DSRecord *DNSSECKeyDSRecord `json:"dsRecord,omitempty"`
}

// A DNSSECKeyV11 represents a DNSSEC Key (Key-Signing or Zone-Signing) as it
// appeared in Traffic Ops API version 1.1. This structure still exists because
// it is used by modern structures, but in general should not be used on its
// own, and github.com/apache/trafficcontrol/v8/lib/go-tc.DNSSECKey should usually
// be used instead.
type DNSSECKeyV11 struct {
	InceptionDateUnix  int64                 `json:"inceptionDate"`
	ExpirationDateUnix int64                 `json:"expirationDate"`
	Name               string                `json:"name"`
	TTLSeconds         uint64                `json:"ttl,string"`
	Status             string                `json:"status"`
	EffectiveDateUnix  int64                 `json:"effectiveDate"`
	Public             string                `json:"public"`
	Private            string                `json:"private"`
	DSRecord           *DNSSECKeyDSRecordV11 `json:"dsRecord,omitempty"`
}

// DNSSECKeyDSRecordRiak is a DNSSEC key DS record, as stored in Riak.
// This is specifically the key data, without the DS record text (which can be computed), and is also the format used in API 1.1 through 1.3.
type DNSSECKeyDSRecordRiak DNSSECKeyDSRecordV11

// DNSSECKeyDSRecord structures are used by the Traffic Ops API after version
// 1.1, and they contain all the same information as well as an additional
// "Text" property of unknown purpose.
type DNSSECKeyDSRecord struct {
	DNSSECKeyDSRecordV11
	Text string `json:"text"`
}

// DNSSECKeyDSRecordV11 structures contain meta information for Key-Signing
// DNSSEC Keys (KSKs) used by Delivery Services.
type DNSSECKeyDSRecordV11 struct {
	Algorithm  int64  `json:"algorithm,string"`
	DigestType int64  `json:"digestType,string"`
	Digest     string `json:"digest"`
}

// CDNDNSSECGenerateReqDate is the date accepted by CDNDNSSECGenerateReq.
// This will unmarshal a UNIX epoch integer, a RFC3339 string, the old format string used by Perl '2018-08-21+14:26:06', and the old format string sent by the Portal '2018-08-21 14:14:42'.
// This exists to fix a critical bug, see https://github.com/apache/trafficcontrol/issues/2723 - it SHOULD NOT be used by any other endpoint.
type CDNDNSSECGenerateReqDate int64

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (i *CDNDNSSECGenerateReqDate) UnmarshalJSON(d []byte) error {
	const oldPortalDateFormat = `2006-01-02 15:04:05`
	const oldPerlUIDateFormat = `2006-01-02+15:04:05`
	if len(d) == 0 {
		return errors.New("empty object")
	}
	if d[0] == '"' {
		d = d[1 : len(d)-1] // strip JSON quotes, to accept the UNIX epoch as a string or number
	}
	if di, err := strconv.ParseInt(string(d), 10, 64); err == nil {
		*i = CDNDNSSECGenerateReqDate(di)
		return nil
	}
	if t, err := time.Parse(time.RFC3339, string(d)); err == nil {
		*i = CDNDNSSECGenerateReqDate(t.Unix())
		return nil
	}
	if t, err := time.Parse(oldPortalDateFormat, string(d)); err == nil {
		*i = CDNDNSSECGenerateReqDate(t.Unix())
		return nil
	}
	if t, err := time.Parse(oldPerlUIDateFormat, string(d)); err == nil {
		*i = CDNDNSSECGenerateReqDate(t.Unix())
		return nil
	}
	return errors.New("invalid date")
}

// A CDNDNSSECGenerateReq is the structure used to request generation of CDN
// DNSSEC Keys through the Traffic Ops API's /cdns/dnsseckeys/generate
// endpoint.
type CDNDNSSECGenerateReq struct {
	// Key is the CDN name, as documented in the API documentation.
	Key               *string                   `json:"key"`
	TTL               *util.JSONIntStr          `json:"ttl"`
	KSKExpirationDays *util.JSONIntStr          `json:"kskExpirationDays"`
	ZSKExpirationDays *util.JSONIntStr          `json:"zskExpirationDays"`
	EffectiveDateUnix *CDNDNSSECGenerateReqDate `json:"effectiveDate"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (r CDNDNSSECGenerateReq) Validate(tx *sql.Tx) error {
	validateErrs := validation.Errors{
		"key (CDN name)":    validation.Validate(r.Key, validation.NotNil),
		"ttl":               validation.Validate(r.TTL, validation.NotNil),
		"kskExpirationDays": validation.Validate(r.KSKExpirationDays, validation.NotNil),
		"zskExpirationDays": validation.Validate(r.ZSKExpirationDays, validation.NotNil),
		// effective date is optional
	}
	return util.JoinErrs(tovalidate.ToErrors(validateErrs))
}
