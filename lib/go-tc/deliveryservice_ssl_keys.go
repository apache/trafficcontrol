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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// DeliveryServiceSSLKeysResponse ...
type DeliveryServiceSSLKeysResponse struct {
	Response DeliveryServiceSSLKeys `json:"response"`
}

// DeliveryServiceSSLKeysCertificate ...
type DeliveryServiceSSLKeysCertificate struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
	CSR string `json:"csr"`
}

// DeliveryServiceSSLKeys ...
type DeliveryServiceSSLKeys struct {
	CDN             string                            `json:"cdn,omitempty"`
	DeliveryService string                            `json:"DeliveryService,omitempty"`
	BusinessUnit    string                            `json:"businessUnit,omitempty"`
	City            string                            `json:"city,omitempty"`
	Organization    string                            `json:"organization,omitempty"`
	Hostname        string                            `json:"hostname,omitempty"`
	Country         string                            `json:"country,omitempty"`
	State           string                            `json:"state,omitempty"`
	Key             string                            `json:"key"`
	Version         int                               `json:"version"`
	Certificate     DeliveryServiceSSLKeysCertificate `json:"certificate,omitempty"`
}

/*
 * The DeliveryServicesSSLKeys are stored in RIAK as JSON.
 * It was found that the "Version" field has been written to
 * RIAK as both a string numeral enclosed in quotes ie,
 *	"version: "1"
 * and sometimes as an integer ie,
 *	"version: 1
 * In order to deal with this problem, a custom Unmarshal() workaround
 * is used, see below.
 *
 */
func (v *DeliveryServiceSSLKeys) UnmarshalJSON(b []byte) (err error) {
	type Alias DeliveryServiceSSLKeys
	o := &struct {
		Version interface{} `json:"version"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	if err = json.Unmarshal(b, &o); err == nil {
		switch t := o.Version.(type) {
		case float64:
			v.Version = int(t)
			break
		case int:
			v.Version = t
			break
		case string:
			v.Version, err = strconv.Atoi(t)
			break
		default:
			err = fmt.Errorf("Version field is an unandled type: %T", t)
		}
	}
	return err
}

// DNSSECKeys is the DNSSEC keys object stored in Riak. The map key strings are both DeliveryServiceNames and CDNNames.
type DNSSECKeys map[string]DNSSECKeySet

type DNSSECKeySet struct {
	ZSK []DNSSECKey `json:"zsk"`
	KSK []DNSSECKey `json:"ksk"`
}

type DNSSECKey struct {
	InceptionDateUnix  int64              `json:"inceptionDate"`
	ExpirationDateUnix int64              `json:"expirationDate"`
	Name               string             `json:"name"`
	TTLSeconds         uint64             `json:"ttl,string"`
	Status             string             `json:"status"`
	EffectiveDateUnix  int64              `json:"effectiveDate"`
	Public             string             `json:"public"`
	Private            string             `json:"private"`
	DSRecord           *DNSSECKeyDSRecord `json:"dsRecord,omitempty"`
}

type DNSSECKeyDSRecord struct {
	Algorithm  int64  `json:"algorithm,string"`
	DigestType int64  `json:"digestType,string"`
	Digest     string `json:"digest"`
}

type CDNDNSSECGenerateReq struct {
	// Key is the CDN name, as documented in the API documentation.
	Key *string `json:"key"`
	// Name is the CDN domain, as documented in the API documentation.
	Name              *string `json:"name"`
	TTL               *uint64 `json:"ttl,string"`
	KSKExpirationDays *uint64 `json:"kskExpirationDays,string"`
	ZSKExpirationDays *uint64 `json:"zskExpirationDays,string"`
	EffectiveDateUnix *int64  `json:"effectiveDate"`
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
