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
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/lestrrat-go/jwx/jwk"
)

// Authentication methods used for signing Delivery Service SSL certificates.
const (
	SelfSignedCertAuthType           = "Self Signed"
	CertificateAuthorityCertAuthType = "Certificate Authority"
	LetsEncryptAuthType              = "Lets Encrypt"
)

// SSLKeysAddResponse is a struct to store the response of addition of ssl keys for a DS,
// along with any alert messages.
type SSLKeysAddResponse struct {
	Response string `json:"response"`
	Alerts
}

// DeliveryServiceSSLKeysResponse is the type of a response from Traffic Ops to
// GET requests made to its /deliveryservices/xmlId/{{XML ID}}/sslkeys API
// endpoint.
type DeliveryServiceSSLKeysResponse struct {
	Response DeliveryServiceSSLKeys `json:"response"`
	Alerts
}

// DeliveryServiceSSLKeys contains information about an SSL key and certificate
// used by a Delivery Service to secure its HTTP-delivered content.
type DeliveryServiceSSLKeys struct {
	AuthType        string                            `json:"authType,omitempty"`
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

// DeliveryServiceSSLKeysV4 is the representation of a DeliveryServiceSSLKeys in the latest minor version of
// version 4 of the Traffic Ops API.
type DeliveryServiceSSLKeysV4 = DeliveryServiceSSLKeysV40

// DeliveryServiceSSLKeysV40 structures contain information about an SSL key
// certificate pair used by a Delivery Service.
//
// "V40" is used because this structure was first introduced in version 4.0 of
// the Traffic Ops API.
type DeliveryServiceSSLKeysV40 struct {
	DeliveryServiceSSLKeysV15
	Sans []string `json:"sans,omitempty"`
}

// DeliveryServiceSSLKeysV15 structures contain information about an SSL key
// certificate pair used by a Delivery Service.
//
// "V15" is used because this structure was first introduced in version 1.5 of
// the Traffic Ops API.
//
// This is, ostensibly, an updated version of DeliveryServiceSSLKeys, but
// beware that this may not be completely accurate as the predecessor structure
// appears to be used in many more contexts than this structure.
type DeliveryServiceSSLKeysV15 struct {
	DeliveryServiceSSLKeys
	Expiration time.Time `json:"expiration,omitempty"`
}

// SSLKeyExpirationInformation contains information about an SSL key's expiration.
type SSLKeyExpirationInformation struct {
	DeliveryService string    `json:"deliveryservice"`
	CDN             string    `json:"cdn"`
	Provider        string    `json:"provider"`
	Expiration      time.Time `json:"expiration"`
	Federated       bool      `json:"federated"`
}

// SSLKeyRequestFields contain metadata information for generating SSL keys for
// Delivery Services through the Traffic Ops API. Specifically, they contain
// everything except the manner in which the generated certificates should be
// signed, information that can be extracted from the Delivery Service for
// which the request is being made, and any key/certificate pair being added
// rather than generated from this information.
type SSLKeyRequestFields struct {
	BusinessUnit *string `json:"businessUnit,omitempty"`
	City         *string `json:"city,omitempty"`
	Organization *string `json:"organization,omitempty"`
	HostName     *string `json:"hostname,omitempty"`
	Country      *string `json:"country,omitempty"`
	State        *string `json:"state,omitempty"`
	Version      *int    `json:"version,omitempty"`
}

// DeliveryServiceSSLKeysReq structures are requests for the generation of SSL
// key certificate pairs for a Delivery Service, and this is, in fact, the
// structure required for POST request bodies to the
// /deliveryservices/sslkeys/generate Traffic Ops API endpoint.
type DeliveryServiceSSLKeysReq struct {
	AuthType        *string `json:"authType,omitempty"`
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

// DeliveryServiceSSLKeysGenerationResponse is the type of a response from
// Traffic Ops to a request for generation of SSL Keys for a Delivery Service.
type DeliveryServiceSSLKeysGenerationResponse struct {
	Response string `json:"response"`
	Alerts
}

// DeliveryServiceSSLKeysCertificate contains an SSL key/certificate pair for a
// Delivery Service, as well as the Certificate Signing Request associated with
// them.
type DeliveryServiceSSLKeysCertificate struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
	CSR string `json:"csr"`
}

// Sanitize ensures that if either the DeliveryService or Key property of a
// DeliveryServiceSSLKeysReq is nil, it will be set to a reference to a copy of
// the value of the other property (if that is not nil).
//
// This does NOT ensure that both fields are not nil, nor does it ensure that
// they match.
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

// validateSharedRequiredRequestFields validates the request fields that are shared and required by both 'add' and 'generate' requests.
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

// DeliveryServiceAddSSLKeysReq structures are requests to import
// key/certificate pairs for a Delivery Service directly, and are the type of
// structures required for POST request bodies to the
// /deliveryservices/sslkeys/add endpoint of the Traffic Ops API.
type DeliveryServiceAddSSLKeysReq struct {
	DeliveryServiceSSLKeysReq
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
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
		if r.Certificate.CSR == "" && (r.AuthType == nil || *r.AuthType != LetsEncryptAuthType) {
			errs = append(errs, "certificate.csr required")
		}
	}
	if len(errs) > 0 {
		return errors.New("missing fields: " + strings.Join(errs, "; "))
	}
	return nil
}

// DeliveryServiceGenSSLKeysReq structures are requests to generate new
// key/certificate pairs for a Delivery Service, and are the type of structures
// required for POST request bodies to the /deliveryservices/sslkeys/generate
// endpoint of the Traffic Ops API.
type DeliveryServiceGenSSLKeysReq struct {
	DeliveryServiceSSLKeysReq
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
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

// DeliveryServiceAcmeSSLKeysReq structures are requests to generate new
// key/certificate pairs for a Delivery Service using an ACME provider, and are
// the type of structures required for POST request bodies to the
// /deliveryservices/sslkeys/generate/letsencrypt endpoint of the Traffic Ops
// API.
type DeliveryServiceAcmeSSLKeysReq struct {
	DeliveryServiceSSLKeysReq
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (r *DeliveryServiceAcmeSSLKeysReq) Validate(tx *sql.Tx) error {
	r.Sanitize()
	errs := r.validateSharedRequiredRequestFields()
	if len(errs) > 0 {
		return errors.New("missing fields: " + strings.Join(errs, "; "))
	}
	errs = r.validateAcmeSpecificFields()
	if len(errs) > 0 {
		return errors.New("missing fields: " + strings.Join(errs, "; "))
	}
	return nil
}

func (r *DeliveryServiceAcmeSSLKeysReq) validateAcmeSpecificFields() []string {
	errs := []string{}
	if checkNilOrEmpty(r.AuthType) {
		errs = append(errs, "authType required")
	}
	return errs
}

func checkNilOrEmpty(s *string) bool {
	return s == nil || *s == ""
}

// TrafficVaultPing represents the status of a given Traffic Vault server.
type TrafficVaultPing struct {
	Status string `json:"status"`
	Server string `json:"server"`
}

// TrafficVaultPingResponse represents the JSON HTTP response returned by the /vault/ping route.
type TrafficVaultPingResponse struct {
	Response TrafficVaultPing `json:"response"`
	Alerts
}

// URLSigKeys is the type of the `response` property of responses from Traffic
// Ops to GET requests made to the /deliverservices/xmlId/{{XML ID}}/urlkeys
// endpoint of its API.
type URLSigKeys map[string]string

// URLSignatureKeysResponse is the type of a response from Traffic Ops to a request
// for the URL Signing keys of a Delivery Service - in API version 4.0.
type URLSignatureKeysResponse struct {
	Response URLSigKeys `json:"response"`
	Alerts
}

// CDNSSLKeysResp is a slice of CDNSSLKeys.
//
// Deprecated: This is not used by any known ATC code and has no known purpose.
// Therefore, it's probably just technical debt subject to removal in the near
// future.
type CDNSSLKeysResp []CDNSSLKey

// CDNSSLKey structures represent an SSL key/certificate pair used by a CDN.
// This is the structure used by each entry of the `response` array property of
// responses from Traffic Ops to GET requests made to its
// /cdns/name/{{Name}}/sslkeys API endpoint.
type CDNSSLKey struct {
	DeliveryService string        `json:"deliveryservice"`
	HostName        string        `json:"hostname"`
	Certificate     CDNSSLKeyCert `json:"certificate"`
}

// A CDNSSLKeyCert represents an SSL key/certificate pair used by a CDN,
// without any other associated data.
type CDNSSLKeyCert struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}

// A CDNGenerateKSKReq is a request to generate Key-Signing Keys for CDNs for
// use in DNSSEC operations, and is the structure required for the bodies of
// POST requests made to the /cdns/{{Name}}/dnsseckeys/ksk/generate endpoint of
// the Traffic Ops API.
type CDNGenerateKSKReq struct {
	ExpirationDays *uint64    `json:"expirationDays"`
	EffectiveDate  *time.Time `json:"effectiveDate"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
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

// Sanitize ensures that the CDNGenerateKSKReq's EffectiveDate is not nil.
//
// If it was nil, this will set it to the current time.
func (r *CDNGenerateKSKReq) Sanitize() {
	if r.EffectiveDate == nil {
		now := time.Now()
		r.EffectiveDate = &now
	}
}

// GetRenewalKid extracts the value of the private "renewal_kid" field from the
// key set. If the set does not contain that private field, or if it does but
// the value is not a string, this returns nil, otherwise it will return a
// pointer to the value of that field.
func GetRenewalKid(set jwk.Set) *string {
	v, ok := set.Field(`renewal_kid`)
	if !ok {
		return nil
	}
	switch v := v.(type) {
	case string:
		return &v
	default:
		return nil
	}
}

// JWKSMap is a mapping of names of JSON Web Token Set to those sets.
type JWKSMap map[string]jwk.Set

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (ksm *JWKSMap) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	*ksm = make(map[string]jwk.Set)
	for k, v := range m {
		set, err := jwk.Parse(v)
		if err != nil {
			return err
		}
		(*ksm)[k] = set
	}
	return nil
}
