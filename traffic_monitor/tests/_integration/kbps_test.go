package _integration

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
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestKBPS(t *testing.T) {
	crc, err := TMClient.CRConfig()
	if err != nil {
		t.Fatalf("client CRConfig error expected nil, actual %v\n", err)
	}

	if len(crc.ContentServers) == 0 {
		t.Fatalf("Monitor CRConfig has no servers, cannot test KBPS")
	}

	serverName := ""
	server := tc.CRConfigTrafficOpsServer{}
	for crcServerName, crcServer := range crc.ContentServers {
		server = crcServer
		serverName = crcServerName
		break
	}
	if server.Ip == nil {
		t.Fatalf("Monitor CRConfig server '" + serverName + "' has no Ip, cannot test KBPS")
	}
	if server.Port == nil {
		t.Fatalf("Monitor CRConfig server '" + serverName + "' has no Port, cannot test KBPS")
	}

	const bytesPerKilobit = 125

	expectedKbps := 10000

	httpClient := http.Client{Timeout: time.Duration(Config.Default.Session.TimeoutInSecs) * time.Second}

	kbps10 := bytesPerKilobit * expectedKbps
	uri := fmt.Sprintf(`http://%v:%v/cmd/setstat?remap=num1.example.net&stat=out_bytes&min=%v&max=%v`, *server.Ip, *server.Port, kbps10, kbps10)
	resp, err := httpClient.Get(uri)
	if err != nil {
		t.Fatalf("Error posting fake cache command '" + uri + "': " + err.Error())
	}
	defer log.Close(resp.Body, "Unable to close http client "+uri)

	time.Sleep(time.Second * 5) // TODO determine if there's a faster or more precise way to wait for polled data?

	kbps, err := TMClient.BandwidthKBPS()
	if err != nil {
		t.Fatalf("getting monitor bandwidth kbps: %v\n", err)
	}

	if kbps < float64(expectedKbps/2) || kbps > float64(expectedKbps*2) {
		t.Errorf("monitor bandwidth kbps expected %v-%v actual %v\n", expectedKbps/2, expectedKbps*2, kbps)
	}
}
