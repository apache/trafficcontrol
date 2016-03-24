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
	"fmt"
	"io/ioutil"
	"testing"
)

func TestServer(t *testing.T) {
	fmt.Println("Running Server Tests")
	text, err := ioutil.ReadFile("testdata/servers.json")
	if err != nil {
		t.Skip("Skipping servers test, no servers.json found.")
	}
	serverList, err := serverUnmarshall(text)
	if err != nil {
		t.Fatal(err)
	}
	for _, server := range serverList.Response {
		if len(server.HostName) == 0 {
			t.Fatal("server result does not contain 'HostName'")
		}
		if len(server.DomainName) == 0 {
			t.Fatal("server result does not contain 'DomainName'")
		}
		name := fmt.Sprintf("%s.%s", server.HostName, server.DomainName)
		if len(server.Id) == 0 {
			t.Errorf("Id is null for server: %s", name)
		}
		if len(server.IloIpAddress) == 0 {
			t.Errorf("IloIpAddress is null for server: %s", name)
		}
		if len(server.IloIpGateway) == 0 {
			t.Errorf("IloIpGateway is null for server: %s", name)
		}
		if len(server.IloIpNetmask) == 0 {
			t.Errorf("IloIpNetmask is null for server: %s", name)
		}
		if len(server.IloPassword) == 0 {
			t.Errorf("IloIpPassword is null for server: %s", name)
		}
		if len(server.IloUsername) == 0 {
			t.Errorf("IloUsername is null for server: %s", name)
		}
	}
}
