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
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-util"
)

const DNSSECKSKType = "ksk"
const DNSSECZSKType = "zsk"
const DNSSECKeyStatusNew = "new"
const DNSSECKeyStatusExpired = "expired"
const DNSSECStatusExisting = "existing"

// DeliveryServiceSSLKeysResponse ...
type DeliveryServiceSSLKeysResponse struct {
	Response DeliveryServiceSSLKeys `json:"response"`
}

// DeliveryServiceSSLKeys ...
type DeliveryServiceSSLKeys struct {
	CDN             string                            `json:"cdn,omitempty"`
	DeliveryService string                            `json:"deliveryservice,omitempty"`
	BusinessUnit    string                            `json:"businessUnit,omitempty"`
	City            string                            `json:"city,omitempty"`
	Organization    string                            `json:"organization,omitempty"`
	Hostname        string                            `json:"hostname,omitempty"`
	Country         string                            `json:"country,omitempty"`
	State           string                            `json:"state,omitempty"`
	Key             string                            `json:"key"`
	Version         util.JSONIntStr                   `json:"version"`
	Certificate     DeliveryServiceSSLKeysCertificate `json:"certificate,omitempty"`
}

type DeliveryServiceSSLKeysReq struct {
	CDN             *string `json:"cdn,omitempty"`
	DeliveryService *string `json:"deliveryservice,omitempty"`
	BusinessUnit    *string `json:"businessUnit,omitempty"`
	City            *string `json:"city,omitempty"`
	Organization    *string `json:"organization,omitempty"`
	HostName        *string `json:"hostname,omitempty"`
	Country         *string `json:"country,omitempty"`
	State           *string `json:"state,omitempty"`
	// Key is the XMLID of the delivery service
	Key         *string                            `json:"key"`
	Version     *util.JSONIntStr                   `json:"version"`
	Certificate *DeliveryServiceSSLKeysCertificate `json:"certificate,omitempty"`
}

// DeliveryServiceSSLKeysCertificate ...
type DeliveryServiceSSLKeysCertificate struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
	CSR string `json:"csr"`
}

func (r *DeliveryServiceSSLKeysReq) Sanitize() {
	// DeliveryService and Key are the same value, so if the user sent one but not the other, set the missing one, in the principle of "be liberal in what you accept."
	if r.DeliveryService == nil && r.Key != nil {
		k := *r.Key // sqlx fails with aliased pointers, so make a new one
		r.DeliveryService = &k
	} else if r.Key == nil && r.DeliveryService != nil {
		k := *r.DeliveryService // sqlx fails with aliased pointers, so make a new one
		r.Key = &k
	}
}

// validateSharedRequiredRequestFields validates the request fields that are shared and required by both 'add' and 'generate' requests
func (r *DeliveryServiceSSLKeysReq) validateSharedRequiredRequestFields() []string {
	errs := []string{}
	if checkNilOrEmpty(r.CDN) {
		errs = append(errs, "cdn required")
	}
	if r.Version == nil {
		errs = append(errs, "version required")
	}
	if checkNilOrEmpty(r.Key) {
		errs = append(errs, "key required")
	}
	if checkNilOrEmpty(r.DeliveryService) {
		errs = append(errs, "deliveryservice required")
	}
	if r.Key != nil && r.DeliveryService != nil && *r.Key != *r.DeliveryService {
		errs = append(errs, "deliveryservice and key must match")
	}
	if checkNilOrEmpty(r.HostName) {
		errs = append(errs, "hostname required")
	}
	return errs
}

type DeliveryServiceAddSSLKeysReq struct {
	DeliveryServiceSSLKeysReq
}

func (r *DeliveryServiceAddSSLKeysReq) Validate(tx *sql.Tx) error {
	r.Sanitize()
	errs := r.validateSharedRequiredRequestFields()
	if r.Certificate == nil {
		errs = append(errs, "certificate required")
	} else {
		if r.Certificate.Key == "" {
			errs = append(errs, "certificate.key required")
		}
		if r.Certificate.Crt == "" {
			errs = append(errs, "certificate.crt required")
		}
		if r.Certificate.CSR == "" {
			errs = append(errs, "certificate.csr required")
		}
	}
	if len(errs) > 0 {
		return errors.New("missing fields: " + strings.Join(errs, "; "))
	}
	return nil
}

type DeliveryServiceGenSSLKeysReq struct {
	DeliveryServiceSSLKeysReq
}

func (r *DeliveryServiceGenSSLKeysReq) Validate(tx *sql.Tx) error {
	r.Sanitize()
	errs := r.validateSharedRequiredRequestFields()
	if checkNilOrEmpty(r.BusinessUnit) {
		errs = append(errs, "businessUnit required")
	}
	if checkNilOrEmpty(r.City) {
		errs = append(errs, "city required")
	}
	if checkNilOrEmpty(r.Organization) {
		errs = append(errs, "organization required")
	}
	if checkNilOrEmpty(r.Country) {
		errs = append(errs, "country required")
	}
	if checkNilOrEmpty(r.State) {
		errs = append(errs, "state required")
	}
	if len(errs) > 0 {
		return errors.New("missing fields: " + strings.Join(errs, "; "))
	}
	return nil
}

func checkNilOrEmpty(s *string) bool {
	return s == nil || *s == ""
}

type RiakPingResp struct {
	Status string `json:"status"`
	Server string `json:"server"`
}

// DNSSECKeys is the DNSSEC keys as stored in Riak, plus the DS record text.
type DNSSECKeys map[string]DNSSECKeySet

// DNSSECKeysV11 is the DNSSEC keys object stored in Riak. The map key strings are both DeliveryServiceNames and CDNNames.

type DNSSECKeysRiak DNSSECKeysV11

type DNSSECKeysV11 map[string]DNSSECKeySetV11

type DNSSECKeySet struct {
	ZSK []DNSSECKey `json:"zsk"`
	KSK []DNSSECKey `json:"ksk"`
}

// DNSSECKeyDSRecordRiak is a DNSSEC key set (ZSK and KSK), as stored in Riak.
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
	Key *string `json:"key"`
	// Name is the CDN domain, as documented in the API documentation.
	Name              *string                   `json:"name"`
	TTL               *util.JSONIntStr          `json:"ttl"`
	KSKExpirationDays *util.JSONIntStr          `json:"kskExpirationDays"`
	ZSKExpirationDays *util.JSONIntStr          `json:"zskExpirationDays"`
	EffectiveDateUnix *CDNDNSSECGenerateReqDate `json:"effectiveDate"`
}

func (r CDNDNSSECGenerateReq) Validate(tx *sql.Tx) error {
	errs := []string{}
	if r.Key == nil {
		errs = append(errs, "key (cdn name) must be set")
	}
	if r.Name == nil {
		errs = append(errs, "name (cdn domain name) must be set")
	}
	if r.TTL == nil {
		errs = append(errs, "ttl must be set")
	}
	if r.KSKExpirationDays == nil {
		errs = append(errs, "kskExpirationDays must be set")
	}
	if r.ZSKExpirationDays == nil {
		errs = append(errs, "zskExpirationDays must be set")
	}
	// effective date is optional
	if len(errs) > 0 {
		return errors.New("missing fields: " + strings.Join(errs, "; "))
	}
	return nil
}

type URLSigKeys map[string]string

type CDNSSLKeysResp []CDNSSLKey

type CDNSSLKey struct {
	DeliveryService string        `json:"deliveryservice"`
	HostName        string        `json:"hostname"`
	Certificate     CDNSSLKeyCert `json:"certificate"`
}

type CDNSSLKeyCert struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}

type CDNGenerateKSKReq struct {
	ExpirationDays *uint64    `json:"expirationDays"`
	EffectiveDate  *time.Time `json:"effectiveDate"`
}

func (r *CDNGenerateKSKReq) Validate(tx *sql.Tx) error {
	r.Sanitize()
	errs := []string{}
	if r.ExpirationDays == nil || *r.ExpirationDays == 0 {
		errs = append(errs, "expiration missing")
	}
	// effective date is optional
	if len(errs) > 0 {
		return errors.New("missing fields: " + strings.Join(errs, "; "))
	}
	return nil
}

func (r *CDNGenerateKSKReq) Sanitize() {
	if r.EffectiveDate == nil {
		now := time.Now()
		r.EffectiveDate = &now
	}
}
