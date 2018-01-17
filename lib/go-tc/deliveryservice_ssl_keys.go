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
	"encoding/json"
	"fmt"
	"strconv"
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
	DeliveryService string                            `json:"deliveryservice,omitempty"`
	BusinessUnit    string                            `json:"unit,omitempty"`
	City            string                            `json:"city,omitempty"`
	Organization    string                            `json:"org,omitempty"`
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
