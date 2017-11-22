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
	"net/url"
	"strings"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// Servers gets an array of servers
// Deprecated: use GetServers
func (to *Session) Servers() ([]tc.Server, error) {
	s, _, err := to.GetServers()
	return s, err
}

func (to *Session) GetServers() ([]tc.Server, ReqInf, error) {
	url := "/api/1.2/servers.json"
	resp, remoteAddr, err := to.request("GET", url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// Server gets a server by hostname
// Deprecated: use GetServer
func (to *Session) Server(name string) (*tc.Server, error) {
	s, _, err := to.GetServer(name)
	return s, err
}

func (to *Session) GetServer(name string) (*tc.Server, ReqInf, error) {
	url := fmt.Sprintf("/api/1.2/servers/hostname/%s/details", name)
	resp, remoteAddr, err := to.request("GET", url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	data := tc.ServersDetailResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return &data.Response, reqInf, nil
}

// ServersByType gets an array of serves of a specified type.
func (to *Session) ServersByType(qparams url.Values) ([]tc.Server, error) {
	ss, _, err := to.GetServersByType(qparams)
	return ss, err
}

func (to *Session) GetServersByType(qparams url.Values) ([]tc.Server, ReqInf, error) {
	url := fmt.Sprintf("/api/1.2/servers.json?%s", qparams.Encode())
	resp, remoteAddr, err := to.request("GET", url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.ServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// ServersFqdn returns a the full domain name for the server short name passed in.
// Deprecated: use GetServersFQDN
func (to *Session) ServersFqdn(n string) (string, error) {
	f, _, err := to.GetServerFQDN(n)
	return f, err
}

func (to *Session) GetServerFQDN(n string) (string, ReqInf, error) {
	// TODO fix to only request one server
	fdn := ""
	servers, reqInf, err := to.GetServers()
	if err != nil {
		return "Error", reqInf, err
	}

	for _, server := range servers {
		if server.HostName == n {
			fdn = fmt.Sprintf("%s.%s", server.HostName, server.DomainName)
		}
	}
	if fdn == "" {
		return "Error", reqInf, fmt.Errorf("No Server %s found", n)
	}
	return fdn, reqInf, nil
}

// ServersShortNameSearch returns a slice of short server names that match a greedy match.
// Deprecated: use GetServersShortNameSearch
func (to *Session) ServersShortNameSearch(shortname string) ([]string, error) {
	ss, _, err := to.GetServersShortNameSearch(shortname)
	return ss, err
}

func (to *Session) GetServersShortNameSearch(shortname string) ([]string, ReqInf, error) {
	var serverlst []string
	servers, reqInf, err := to.GetServers()
	if err != nil {
		serverlst = append(serverlst, "N/A")
		return serverlst, reqInf, err
	}
	for _, server := range servers {
		if strings.Contains(server.HostName, shortname) {
			serverlst = append(serverlst, server.HostName)
		}
	}
	if len(serverlst) == 0 {
		serverlst = append(serverlst, "N/A")
		return serverlst, reqInf, fmt.Errorf("No Servers Found")
	}
	return serverlst, reqInf, nil
}
