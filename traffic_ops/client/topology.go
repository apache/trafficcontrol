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

package client

import (
	"encoding/json"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"net"
	"net/http"
	"net/url"
)

const (
	ApiTopologies = apiBase + "/topologies"
)

// CreateTopology creates a topology and returns the response.
func (to *Session) CreateTopology(top tc.Topology) (*tc.TopologyResponse, ReqInf, error, int) {
	var (
		statusCode = http.StatusNotAcceptable
		remoteAddr net.Addr
	)
	reqBody, err := json.Marshal(top)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err, statusCode
	}
	resp, remoteAddr, err := to.request(http.MethodPost, ApiTopologies, reqBody)
	if resp != nil {
		statusCode = resp.StatusCode
	}
	if err != nil {
		return nil, reqInf, err, statusCode
	}
	defer log.Close(resp.Body, "unable to close response")
	var topResp tc.TopologyResponse
	if err = json.NewDecoder(resp.Body).Decode(&topResp); err != nil {
		return nil, reqInf, err, statusCode
	}
	return &topResp, reqInf, nil, statusCode
}

// GetTopologies returns all topologies.
func (to *Session) GetTopologies() ([]tc.Topology, ReqInf, error, int) {
	var statusCode = http.StatusNotAcceptable
	resp, remoteAddr, err := to.request(http.MethodGet, ApiTopologies, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		statusCode = resp.StatusCode
	}
	if err != nil {
		return nil, reqInf, err, statusCode
	}
	defer log.Close(resp.Body, "unable to close response")

	var data tc.TopologiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err, statusCode
	}

	return data.Response, reqInf, nil, statusCode
}

// GetTopology returns the given topology by name.
func (to *Session) GetTopology(name string) (*tc.Topology, ReqInf, error, int) {
	var statusCode = http.StatusNotAcceptable
	reqUrl := fmt.Sprintf("%s?name=%s", ApiTopologies, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, reqUrl, nil)
	if resp != nil {
		statusCode = resp.StatusCode
	}
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err, statusCode
	}
	defer log.Close(resp.Body, "unable to close response")

	var data tc.TopologiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err, statusCode
	}

	if len(data.Response) == 1 {
		return &data.Response[0], reqInf, nil, statusCode
	}
	return nil, reqInf, fmt.Errorf("expected one topology in response, instead got: %+v", data.Response), statusCode
}

// UpdateTopology updates a Topology by name.
func (to *Session) UpdateTopology(name string, t tc.Topology) (*tc.TopologyResponse, ReqInf, error, int) {
	var (
		statusCode = http.StatusNotAcceptable
		remoteAddr net.Addr
	)
	reqBody, err := json.Marshal(t)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err, statusCode
	}
	route := fmt.Sprintf("%s?name=%s", ApiTopologies, name)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody)
	if resp != nil {
		statusCode = resp.StatusCode
	}
	if err != nil {
		return nil, reqInf, err, statusCode
	}
	defer log.Close(resp.Body, "unable to close response")
	var response = new(tc.TopologyResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	return response, reqInf, err, statusCode
}

// DeleteTopology deletes the given topology by name.
func (to *Session) DeleteTopology(name string) (tc.Alerts, ReqInf, error, int) {
	var statusCode = http.StatusNotAcceptable
	reqUrl := fmt.Sprintf("%s?name=%s", ApiTopologies, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodDelete, reqUrl, nil)
	if resp != nil {
		statusCode = resp.StatusCode
	}
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err, statusCode
	}
	defer log.Close(resp.Body, "unable to close response")
	var alerts tc.Alerts
	if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err, statusCode
	}
	return alerts, reqInf, nil, statusCode
}
