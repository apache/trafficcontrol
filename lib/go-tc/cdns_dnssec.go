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

	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	DNSSECKSKType          = "ksk"
	DNSSECZSKType          = "zsk"
	DNSSECKeyStatusNew     = "new"
	DNSSECKeyStatusExpired = "expired"
	DNSSECStatusExisting   = "existing"
)

type CDNDNSSECKeysResponse struct {
	Response DNSSECKeys `json:"response"`
	Alerts
}

type GenerateCDNDNSSECKeysResponse struct {
	Response string `json:"response"`
	Alerts
}

type DeleteCDNDNSSECKeysResponse GenerateCDNDNSSECKeysResponse

type RefreshDNSSECKeysResponse GenerateCDNDNSSECKeysResponse

// DNSSECKeys is the DNSSEC keys as stored in Riak, plus the DS record text.
type DNSSECKeys map[string]DNSSECKeySet

// Deprecated: use DNSSECKeysTrafficVault instead
type DNSSECKeysRiak DNSSECKeysV11

type DNSSECKeysTrafficVault DNSSECKeysV11

// DNSSECKeysV11 is the DNSSEC keys object stored in Riak. The map key strings are both DeliveryServiceNames and CDNNames.
type DNSSECKeysV11 map[string]DNSSECKeySetV11

type DNSSECKeySet struct {
	ZSK []DNSSECKey `json:"zsk"`
	KSK []DNSSECKey `json:"ksk"`
}

// DNSSECKeySetV11 is a DNSSEC key set (ZSK and KSK), as stored in Riak.
// This is specifically the key data, without the DS record text (which can be computed), and is also the format used in API 1.1 through 1.3.
type DNSSECKeySetV11 struct {
	ZSK []DNSSECKeyV11 `json:"zsk"`
	KSK []DNSSECKeyV11 `json:"ksk"`
}

type DNSSECKey struct {
	DNSSECKeyV11
	DSRecord *DNSSECKeyDSRecord `json:"dsRecord,omitempty"`
}

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

type DNSSECKeyDSRecord struct {
	DNSSECKeyDSRecordV11
	Text string `json:"text"`
}

type DNSSECKeyDSRecordV11 struct {
	Algorithm  int64  `json:"algorithm,string"`
	DigestType int64  `json:"digestType,string"`
	Digest     string `json:"digest"`
}

// CDNDNSSECGenerateReqDate is the date accepted by CDNDNSSECGenerateReq.
// This will unmarshal a UNIX epoch integer, a RFC3339 string, the old format string used by Perl '2018-08-21+14:26:06', and the old format string sent by the Portal '2018-08-21 14:14:42'.
// This exists to fix a critical bug, see https://github.com/apache/trafficcontrol/issues/2723 - it SHOULD NOT be used by any other endpoint.
type CDNDNSSECGenerateReqDate int64

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

type CDNDNSSECGenerateReq struct {
	// Key is the CDN name, as documented in the API documentation.
	Key               *string                   `json:"key"`
	TTL               *util.JSONIntStr          `json:"ttl"`
	KSKExpirationDays *util.JSONIntStr          `json:"kskExpirationDays"`
	ZSKExpirationDays *util.JSONIntStr          `json:"zskExpirationDays"`
	EffectiveDateUnix *CDNDNSSECGenerateReqDate `json:"effectiveDate"`
}

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
