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
	"net"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const ApiTopologies = apiBase + "/topologies"

// CreateTopology creates a topology and returns the response.
func (to *Session) CreateTopology(top tc.Topology) (*tc.TopologyResponse, ReqInf, error) {
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(top)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, ApiTopologies, reqBody, nil)
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer log.Close(resp.Body, "unable to close response")
	var topResp tc.TopologyResponse
	if err = json.NewDecoder(resp.Body).Decode(&topResp); err != nil {
		return nil, reqInf, err
	}
	return &topResp, reqInf, nil
}

func (to *Session) GetTopologiesWithHdr(header http.Header) ([]tc.Topology, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, ApiTopologies, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer log.Close(resp.Body, "unable to close response")

	var data tc.TopologiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GetTopologies returns all topologies.
// Deprecated: GetTopologies will be removed in 6.0. Use GetTopologiesWithHdr.
func (to *Session) GetTopologies() ([]tc.Topology, ReqInf, error) {
	return to.GetTopologiesWithHdr(nil)
}

func (to *Session) GetTopologyWithHdr(name string, header http.Header) (*tc.Topology, ReqInf, error) {
	reqUrl := fmt.Sprintf("%s?name=%s", ApiTopologies, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodGet, reqUrl, nil, header)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer log.Close(resp.Body, "unable to close response")

	var data tc.TopologiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	if len(data.Response) == 1 {
		return &data.Response[0], reqInf, nil
	}
	return nil, reqInf, fmt.Errorf("expected one topology in response, instead got: %+v", data.Response)
}

// GetTopology returns the given topology by name.
// Deprecated: GetTopology will be removed in 6.0. Use GetTopologyWithHdr.
func (to *Session) GetTopology(name string) (*tc.Topology, ReqInf, error) {
	return to.GetTopologyWithHdr(name, nil)
}

// UpdateTopology updates a Topology by name.
func (to *Session) UpdateTopology(name string, t tc.Topology) (*tc.TopologyResponse, ReqInf, error) {
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(t)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	route := fmt.Sprintf("%s?name=%s", ApiTopologies, name)
	resp, remoteAddr, err := to.request(http.MethodPut, route, reqBody, nil)
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	if err != nil {
		return nil, reqInf, err
	}
	defer log.Close(resp.Body, "unable to close response")

	var response = new(tc.TopologyResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	return response, reqInf, err
}

// DeleteTopology deletes the given topology by name.
func (to *Session) DeleteTopology(name string) (tc.Alerts, ReqInf, error) {
	reqUrl := fmt.Sprintf("%s?name=%s", ApiTopologies, url.QueryEscape(name))
	resp, remoteAddr, err := to.request(http.MethodDelete, reqUrl, nil, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if resp != nil {
		reqInf.StatusCode = resp.StatusCode
	}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer log.Close(resp.Body, "unable to close response")

	var alerts tc.Alerts
	if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return tc.Alerts{}, reqInf, err
	}
	return alerts, reqInf, nil
}
