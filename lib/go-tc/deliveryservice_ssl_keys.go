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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-util"
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
	Version         string                            `json:"version"`
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
	Version     *string                            `json:"version"`
	Certificate *DeliveryServiceSSLKeysCertificate `json:"certificate,omitempty"`
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
	if r.Version == nil {
		r.Version = util.StrPtr("")
	}
}

func (r *DeliveryServiceSSLKeysReq) Validate(tx *sql.Tx) error {
	r.Sanitize()
	errs := []string{}
	if r.CDN == nil {
		errs = append(errs, "cdn required")
	}
	if r.Key == nil {
		errs = append(errs, "key required")
	}
	if r.DeliveryService == nil {
		errs = append(errs, "deliveryservice required")
	}
	if r.Key != nil && r.DeliveryService != nil && *r.Key != *r.DeliveryService {
		errs = append(errs, "deliveryservice and key must match")
	}
	if r.BusinessUnit == nil {
		errs = append(errs, "businessUnit required")
	}
	if r.City == nil {
		errs = append(errs, "city required")
	}
	if r.Organization == nil {
		errs = append(errs, "organization required")
	}
	if r.HostName == nil {
		errs = append(errs, "hostname required")
	}
	if r.Country == nil {
		errs = append(errs, "country required")
	}
	if r.State == nil {
		errs = append(errs, "state required")
	}
	// version is optional
	// certificate is optional
	if len(errs) > 0 {
		return errors.New("missing fields: " + strings.Join(errs, "; "))
	}
	return nil
}

type RiakPingResp struct {
	Status string `json:"status"`
	Server string `json:"server"`
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

type URLSigKeys map[string]string
