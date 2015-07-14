/*
   Copyright 2015 Comcast Cable Communications Management, LLC

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
	"strings"
)

type ServerResponse struct {
	Version  string   `json:"version"`
	Response []Server `json:"response"`
}

type Server struct {
	DomainName     string `json:"domainName"`
	HostName       string `json:"hostName"`
	Id             string `json:"id"`
	IloIpAddress   string `json:"iloIpAddress"`
	IloIpGateway   string `json:"iloIpGateway"`
	IloIpNetmask   string `json:"iloIpNetmask"`
	IloPassword    string `json:"iloPassword"`
	IloUsername    string `json:"iloUsername"`
	InterfaceMtu   string `json:"interfaceMtu"`
	InterfaceName  string `json:"interfaceName"`
	Ip6Address     string `json:"ip6Address"`
	Ip6Gateway     string `json:"ip6Gateway"`
	IpAddress      string `json:"ipAddress"`
	IpGateway      string `json:"ipGateway"`
	IpNetoask      string `json:"ipNetoask"`
	LastUpdated    string `json:"lastUpdated"`
	Location       string `json:"cachegroup"`
	MgmtIpAddress  string `json:"mgmtIpAddress"`
	MgmtIpGateway  string `json:"mgmtIpGateway"`
	MgmtIpNetoask  string `json:"mgmtIpNetmask"`
	PhysLocation   string `json:"physLocation"`
	Profile        string `json:"profile"`
	Rack           string `json:"rack"`
	RouterHostName string `json:"routerHostName"`
	RouterPortName string `json:"routerPortName"`
	Status         string `json:"status"`
	TcpPort        string `json:"tcpPort"`
	Type           string `json:"type"`
	XmppId         string `json:"xmppId"`
	XmppPasswd     string `json:"xmppPasswd"`
}

// Servers
// Get an array of servers
func (to *Session) Servers() ([]Server, error) {
	body, err := to.getBytes("/api/1.1/servers.json")
	if err != nil {
		return nil, err
	}
	serverList, err := serverUnmarshall(body)
	return serverList.Response, err
}

func serverUnmarshall(body []byte) (ServerResponse, error) {
	var data ServerResponse
	err := json.Unmarshal(body, &data)
	return data, err
}

// ServersFqdn
// Returns a the full domain name for the server short name passed in.
func (to *Session) ServersFqdn(n string) (string, error) {
	var fdn string
	fdn = ""
	servers, err := to.Servers()
	if err != nil {
		return "Error", err
	}
	for _, server := range servers {
		if server.HostName == n {
			fdn = server.HostName + "." + server.DomainName
		}
	}
	if fdn == "" {
		return "Error", fmt.Errorf("No Server %s found", n)
	} else {
		return fdn, nil
	}
}

// ShortNameSearch
// Returns a slice of short server names that match a greedy match.
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
