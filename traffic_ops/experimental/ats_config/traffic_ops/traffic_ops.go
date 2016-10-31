
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// \todo move Traffic Ops API to its own package

package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// getClient gets a http client object.
// This exists to encapsulate TLS cert verification failure skipping.
// TODO(fix to not skip cert verification, when the api cert is valid)
func getClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

type TrafficOpsParameter struct {
	LastUpdated string `json:"lastUpdated"`
	Value       string `json:"value"`
	Name        string `json:"name"`
	ConfigFile  string `json:"configFile"`
}

type TrafficOpsParametersResponse struct {
	Response []TrafficOpsParameter `json:"response"`
}

func GetParameters(uri, cookie, profile string) ([]TrafficOpsParameter, error) {
	endpointUri := uri + "/api/1.2/parameters/profile/" + profile + ".json"
	req, err := http.NewRequest("GET", endpointUri, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "mojolicious", Value: cookie})

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paramsResp TrafficOpsParametersResponse
	if err := json.Unmarshal(data, &paramsResp); err != nil {
		return nil, err
	}

	return paramsResp.Response, nil
}

// GetTrafficOpsCookie logs in to Traffic Ops and returns an auth cookie
func GetTrafficOpsCookie(cdnUri, user, pass string) (string, error) {
	uri := cdnUri + `/api/1.2/user/login`
	postdata := `{"u":"` + user + `", "p":"` + pass + `"}`
	req, err := http.NewRequest("POST", uri, strings.NewReader(postdata))
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == `mojolicious` {
			return cookie.Value, nil
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return "", errors.New("No login cookie received: " + string(data))
}

type TrafficOpsServer struct {
	Profile        string `json:"profile"`
	IloUsername    string `json:"iloUsername"`
	Status         string `json:"statusn"`
	IpAddress      string `json:"ipAddress"`
	PhysLocation   string `json:"physLocation"`
	Cachegroup     string `json:"cachegroup"`
	Ip6Gateway     string `json:"ip6Gateway"`
	InterfaceName  string `json:"interfaceName"`
	IloPassword    string `json:"iloPassword"`
	Id             string `json:"id"`
	RouterPortName string `json:"routerPortNamen"`
	LastUpdated    string `json:"lastUpdated"`
	IpNetmask      string `json:"ipNetmask"`
	TcpPort        string `json:"tcpPort"`
	IpGateway      string `json:"ipGateway"`
	MgmtIpAddress  string `json:"mgmtIpAddress"`
	Ip6Address     string `json:"ip6Address"`
	IloIpGateway   string `json:"iloIpGateway"`
	InterfaceMtu   string `json:"interfaceMtu"`
	CdnName        string `json:"cdnName"`
	HostName       string `json:"hostName"`
	IloIpAddress   string `json:"iloIpAddress"`
	MgmtIpNetmask  string `json:"mgmtIpNetmask"`
	Rack           string `json:"rack"`
	MgmtIpGateway  string `json:"mgmtIpGateway"`
	Type           string `json:"type"`
	IloIpNetmask   string `json:"iloIpNetmask"`
	DomainName     string `json:"domainName"`
	RouterHostName string `json:"routerHostName"`
}

type TrafficOpsServersResponse struct {
	Response []TrafficOpsServer `json:"response"`
}

func GetServers(trafficOpsUri, trafficOpsCookie string) ([]TrafficOpsServer, error) {
	uri := trafficOpsUri + `/api/1.2/servers.json`
	req, err := http.NewRequest("GET", uri, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "mojolicious", Value: trafficOpsCookie})

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var serversResponse TrafficOpsServersResponse
	if err := json.Unmarshal(data, &serversResponse); err != nil {
		return nil, err
	}

	return serversResponse.Response, nil
}

func GetServerProfileName(trafficOpsUri, cookie, serverHostname string) (string, error) {
	servers, err := GetServers(trafficOpsUri, cookie)
	if err != nil {
		return "", err
	}

	for _, server := range servers {
		if server.HostName == serverHostname {
			return server.Profile, nil
		}
	}

	return "", errors.New("Server not found")
}
