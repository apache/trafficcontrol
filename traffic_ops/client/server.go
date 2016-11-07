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
)

// ServerResponse ...
type ServerResponse struct {
	Version  string   `json:"version"`
	Response []Server `json:"response"`
}

// Server ...
type Server struct {
	DomainName    string `json:"domainName"`
	HostName      string `json:"hostName"`
	ID            string `json:"id"`
	IloIPAddress  string `json:"iloIpAddress"`
	IloIPGateway  string `json:"iloIpGateway"`
	IloIPNetmask  string `json:"iloIpNetmask"`
	IloPassword   string `json:"iloPassword"`
	IloUsername   string `json:"iloUsername"`
	InterfaceMtu  string `json:"interfaceMtu"`
	InterfaceName string `json:"interfaceName"`
	IP6Address    string `json:"ip6Address"`
	IP6Gateway    string `json:"ip6Gateway"`
	IPAddress     string `json:"ipAddress"`
	IPGateway     string `json:"ipGateway"`
	IPNetmask     string `json:"ipNetmask"`

	LastUpdated    string `json:"lastUpdated"`
	Cachegroup     string `json:"cachegroup"`
	MgmtIPAddress  string `json:"mgmtIpAddress"`
	MgmtIPGateway  string `json:"mgmtIpGateway"`
	MgmtIPNetmask  string `json:"mgmtIpNetmask"`
	PhysLocation   string `json:"physLocation"`
	Profile        string `json:"profile"`
	ProfileDesc    string `json:"profileDesc"`
	CDNName        string `json:"cdnName"`
	Rack           string `json:"rack"`
	RouterHostName string `json:"routerHostName"`
	RouterPortName string `json:"routerPortName"`
	Status         string `json:"status"`
	TCPPort        string `json:"tcpPort"`
	Type           string `json:"type"`
	XMPPID         string `json:"xmppId"`
	XMPPPasswd     string `json:"xmppPasswd"`
}

// Servers gets an array of servers
func (to *Session) Servers() ([]Server, error) {
	url := "/api/1.2/servers.json"
	resp, err := to.request(url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}

// ServersByType gets an array of serves of a specified type.
func (to *Session) ServersByType(qparams url.Values) ([]Server, error) {
	url := fmt.Sprintf("/api/1.2/servers.json?%s", qparams.Encode())
	resp, err := to.request(url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}

// ServersFqdn returns a the full domain name for the server short name passed in.
func (to *Session) ServersFqdn(n string) (string, error) {
	fdn := ""
	servers, err := to.Servers()
	if err != nil {
		return "Error", err
	}

	for _, server := range servers {
		if server.HostName == n {
			fdn = fmt.Sprintf("%s.%s", server.HostName, server.DomainName)
		}
	}
	if fdn == "" {
		return "Error", fmt.Errorf("No Server %s found", n)
	}
	return fdn, nil
}

// ServersShortNameSearch returns a slice of short server names that match a greedy match.
func (to *Session) ServersShortNameSearch(shortname string) ([]string, error) {
	var serverlst []string
	servers, err := to.Servers()
	if err != nil {
		serverlst = append(serverlst, "N/A")
		return serverlst, err
	}
	for _, server := range servers {
		if strings.Contains(server.HostName, shortname) {
			serverlst = append(serverlst, server.HostName)
		}
	}
	if len(serverlst) == 0 {
		serverlst = append(serverlst, "N/A")
		return serverlst, fmt.Errorf("No Servers Found")
	}
	return serverlst, nil
}
