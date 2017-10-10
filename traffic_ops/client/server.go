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
func (to *Session) Servers() ([]tc.Server, error) {
	url := "/api/1.2/servers.json"
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data tc.ServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response, nil
}

// Server gets a server by hostname
func (to *Session) Server(name string) (*tc.Server, error) {
	url := fmt.Sprintf("/api/1.2/servers/hostname/%s/details", name)
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := tc.ServersDetailResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// ServersByType gets an array of serves of a specified type.
func (to *Session) ServersByType(qparams url.Values) ([]tc.Server, error) {
	url := fmt.Sprintf("/api/1.2/servers.json?%s", qparams.Encode())
	resp, err := to.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data tc.ServersResponse
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
