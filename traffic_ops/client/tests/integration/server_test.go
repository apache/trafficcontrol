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

package integration

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestServers(t *testing.T) {

	uri := fmt.Sprintf("/api/1.2/servers.json")
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiServerRes traffic_ops.ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiServerRes); err != nil {
		t.Errorf("Could not decode server json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiServers := apiServerRes.Response

	clientServers, err := to.Servers()
	if err != nil {
		t.Errorf("Could not get servers from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiServers) != len(clientServers) {
		t.Errorf("Server Response Length -- expected %v, got %v\n", len(apiServers), len(clientServers))
	}

	for _, apiServer := range apiServers {
		match := false
		for _, clientServer := range clientServers {
			if apiServer.ID == clientServer.ID {
				match = true
				compareServer(apiServer, clientServer, t)
			}
		}
		if !match {
			t.Errorf("Did not get a server matching %v\n", apiServer.ID)
		}
	}
}

func TestServersByType(t *testing.T) {
	serverType, err := GetType("server")
	if err != nil {
		t.Error("Could not get Type from client for TestServersByType")
		t.FailNow()
	}
	uri := fmt.Sprintf("/api/1.2/servers.json?type=%s", serverType.Name)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiServerRes traffic_ops.ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiServerRes); err != nil {
		t.Errorf("Could not decode server json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiServers := apiServerRes.Response

	params := make(url.Values)
	params.Add("type", serverType.Name)
	clientServers, err := to.ServersByType(params)
	if err != nil {
		t.Errorf("Could not get servers from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiServers) != len(clientServers) {
		t.Errorf("Server Response Length -- expected %v, got %v\n", len(apiServers), len(clientServers))
	}

	for _, apiServer := range apiServers {
		match := false
		for _, clientServer := range clientServers {
			if apiServer.ID == clientServer.ID {
				match = true
				compareServer(apiServer, clientServer, t)
			}
		}
		if !match {
			t.Errorf("Did not get a server matching %v\n", apiServer.ID)
		}
	}

}

func TestServersFQDN(t *testing.T) {
	serverType, err := GetType("server")
	if err != nil {
		t.Error("Could not get Type from client")
		t.FailNow()
	}

	params := make(url.Values)
	params.Add("type", serverType.Name)
	servers, err := to.ServersByType(params)
	if err != nil {
		t.Errorf("Could not get servers from client.  Error is: %v\n", err)
		t.FailNow()
	}

	serverFQDN := fmt.Sprintf("%s.%s", servers[0].HostName, servers[0].DomainName)

	clientFQDN, err := to.ServersFqdn(servers[0].HostName)
	if err != nil {
		t.Errorf("Servers FQDN failed...err: %v\n", err)
	}
	if clientFQDN != serverFQDN {
		t.Errorf("Client FQDN -- expected %v, got %v", serverFQDN, clientFQDN)
	}
}

func compareServer(server1 traffic_ops.Server, server2 traffic_ops.Server, t *testing.T) {
	if server1.CDNName != server2.CDNName {
		t.Errorf("CDNName -- Expected %v, got %v", server1.CDNName, server2.CDNName)
	}
	if server1.Cachegroup != server2.Cachegroup {
		t.Errorf("Cachegroup -- Expected %v, got %v", server1.Cachegroup, server2.Cachegroup)
	}
	if server1.DomainName != server2.DomainName {
		t.Errorf("DomainName -- Expected %v, got %v", server1.DomainName, server2.DomainName)
	}
	if server1.HostName != server2.HostName {
		t.Errorf("HostName -- Expected %v, got %v", server1.HostName, server2.HostName)
	}
	if server1.IP6Address != server2.IP6Address {
		t.Errorf("IP6Address -- Expected %v, got %v", server1.IPAddress, server2.IP6Address)
	}
	if server1.IP6Gateway != server2.IP6Gateway {
		t.Errorf("IP6Gateway -- Expected %v, got %v", server1.IP6Gateway, server2.IP6Gateway)
	}
	if server1.IPAddress != server2.IPAddress {
		t.Errorf("IPAddress -- Expected %v, got %v", server1.IPAddress, server2.IPAddress)
	}
	if server1.IPGateway != server2.IPGateway {
		t.Errorf("IPGateway -- Expected %v, got %v", server1.IPGateway, server2.IPGateway)
	}
	if server1.IPNetmask != server2.IPNetmask {
		t.Errorf("IPNetmask -- Expected %v, got %v", server1.IPNetmask, server2.IPNetmask)
	}
	if server1.IloIPAddress != server2.IloIPAddress {
		t.Errorf("IloIPAddress -- Expected %v, got %v", server1.IloIPAddress, server2.IloIPAddress)
	}
	if server1.IloIPGateway != server2.IloIPGateway {
		t.Errorf("IloIPGateway -- Expected %v, got %v", server1.IloIPGateway, server2.IloIPGateway)
	}
	if server1.IloIPNetmask != server2.IloIPNetmask {
		t.Errorf("IloIPNetmast -- Expected %v, got %v", server1.IloIPNetmask, server2.IloIPNetmask)
	}
	if server1.IloPassword != server2.IloPassword {
		t.Errorf("IloPassword -- Expected %v, got %v", server1.IloPassword, server2.IloPassword)
	}
	if server1.IloUsername != server2.IloUsername {
		t.Errorf("IloUsername -- Expected %v, got %v", server1.IloUsername, server2.IloUsername)
	}
	if server1.InterfaceMtu != server2.InterfaceMtu {
		t.Errorf("InterfaceMty -- Expected %v, got %v", server1.InterfaceMtu, server2.InterfaceMtu)
	}
	if server1.InterfaceName != server2.InterfaceName {
		t.Errorf("InterfaceName -- Expected %v, got %v", server1.InterfaceName, server2.InterfaceName)
	}
	if server1.LastUpdated != server2.LastUpdated {
		t.Errorf("LastUpdated -- Expected %v, got %v", server1.LastUpdated, server2.LastUpdated)
	}
	if server1.MgmtIPAddress != server2.MgmtIPAddress {
		t.Errorf("MgmtIPAddress -- Expected %v, got %v", server1.MgmtIPAddress, server2.MgmtIPAddress)
	}
	if server1.MgmtIPGateway != server2.MgmtIPGateway {
		t.Errorf("MgmtIPGateway -- Expected %v, got %v", server1.MgmtIPGateway, server2.MgmtIPGateway)
	}
	if server1.MgmtIPNetmask != server2.MgmtIPNetmask {
		t.Errorf("MgmtIPNetmask -- Expected %v, got %v", server1.MgmtIPNetmask, server2.MgmtIPNetmask)
	}
	if server1.PhysLocation != server2.PhysLocation {
		t.Errorf("PhysLocation -- Expected %v, got %v", server1.PhysLocation, server2.PhysLocation)
	}
	if server1.Profile != server2.Profile {
		t.Errorf("Profile -- Expected %v, got %v", server1.Profile, server2.Profile)
	}
	if server1.ProfileDesc != server2.ProfileDesc {
		t.Errorf("ProfileDesc -- Expected %v, got %v", server1.ProfileDesc, server2.ProfileDesc)
	}
	if server1.Rack != server2.Rack {
		t.Errorf("Rack -- Expected %v, got %v", server1.Rack, server2.Rack)
	}
	if server1.RouterHostName != server2.RouterHostName {
		t.Errorf("RouterHostName -- Expected %v, got %v", server1.RouterHostName, server2.RouterHostName)
	}
	if server1.RouterPortName != server2.RouterPortName {
		t.Errorf("RouterPortName -- Expected %v, got %v", server1.RouterPortName, server2.RouterPortName)
	}
	if server1.Status != server2.Status {
		t.Errorf("Status -- Expected %v, got %v", server1.Status, server2.Status)
	}
	if server1.TCPPort != server2.TCPPort {
		t.Errorf("TCPPort -- Expected %v, got %v", server1.TCPPort, server2.TCPPort)
	}
	if server1.Type != server2.Type {
		t.Errorf("Type -- Expected %v, got %v", server1.Type, server2.Type)
	}
	if server1.XMPPID != server2.XMPPID {
		t.Errorf("XMPPID -- Expected %v, got %v", server1.XMPPID, server2.XMPPID)
	}
	if server1.XMPPPasswd != server2.XMPPPasswd {
		t.Errorf("XMPPasswd -- Expected %v, got %v", server1.XMPPPasswd, server2.XMPPPasswd)
	}
}
