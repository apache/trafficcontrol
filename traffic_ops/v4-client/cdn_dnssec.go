package client

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
	"fmt"
	"net/url"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	apiCDNsDNSSECKeysGenerate    = "/cdns/dnsseckeys/generate"
	apiCDNsNameDNSSECKeys        = "/cdns/name/%s/dnsseckeys"
	apiCDNsDNSSECRefresh         = "/cdns/dnsseckeys/refresh"
	apiCDNsDNSSECKeysKSKGenerate = "/cdns/%s/dnsseckeys/ksk/generate"
)

// GenerateCDNDNSSECKeys generates DNSSEC keys for the given CDN.
func (to *Session) GenerateCDNDNSSECKeys(req tc.CDNDNSSECGenerateReq, opts RequestOptions) (tc.GenerateCDNDNSSECKeysResponse, toclientlib.ReqInf, error) {
	var resp tc.GenerateCDNDNSSECKeysResponse
	reqInf, err := to.post(apiCDNsDNSSECKeysGenerate, opts, req, &resp)
	return resp, reqInf, err
}

// GetCDNDNSSECKeys gets the DNSSEC keys for the given CDN.
func (to *Session) GetCDNDNSSECKeys(name string, opts RequestOptions) (tc.CDNDNSSECKeysResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(apiCDNsNameDNSSECKeys, url.PathEscape(name))
	var resp tc.CDNDNSSECKeysResponse
	reqInf, err := to.get(route, opts, &resp)
	return resp, reqInf, err
}

// DeleteCDNDNSSECKeys deletes all the DNSSEC keys for the given CDN.
func (to *Session) DeleteCDNDNSSECKeys(name string, opts RequestOptions) (tc.DeleteCDNDNSSECKeysResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(apiCDNsNameDNSSECKeys, url.PathEscape(name))
	var resp tc.DeleteCDNDNSSECKeysResponse
	reqInf, err := to.del(route, opts, &resp)
	return resp, reqInf, err
}

// RefreshDNSSECKeys asynchronously regenerates any expired DNSSEC keys in all CDNs.
func (to *Session) RefreshDNSSECKeys(opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	var resp tc.Alerts
	reqInf, err := to.put(apiCDNsDNSSECRefresh, opts, nil, &resp)
	return resp, reqInf, err
}

// GenerateCDNDNSSECKSK generates the DNSSEC KSKs (key-signing key) for the given CDN.
func (to *Session) GenerateCDNDNSSECKSK(name string, req tc.CDNGenerateKSKReq, opts RequestOptions) (tc.GenerateCDNDNSSECKeysResponse, toclientlib.ReqInf, error) {
	route := fmt.Sprintf(apiCDNsDNSSECKeysKSKGenerate, url.PathEscape(name))
	var resp tc.GenerateCDNDNSSECKeysResponse
	reqInf, err := to.post(route, opts, req, &resp)
	return resp, reqInf, err
}
