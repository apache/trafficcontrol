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

package api

import (
	"encoding/json"
	"net/http"
	"testing"

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestCDNs(t *testing.T) {

	TestPostCDNs(t)
	TestPutCDNs(t)
	TestGetCDNs(t)

}

func TestPostCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		alerts, _, err := TOSession.CreateCDN(cdn)
		log.Debugln("Alerts: %v", alerts)
		if err != nil {
			t.Errorf("could not POST to cdns: %v\n", err)
		}
	}

}

func TestPutCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {

		b, err := json.Marshal(cdn)
		if err != nil {
			t.Errorf("could not marshal data %v\n", err)
		}
		resp, err := Request(*TOSession, http.MethodPost, client.API_v2_CDNs, b)
		if err != nil {
			t.Errorf("could not POST to cdns: %v\n", err)
		}
		defer resp.Body.Close()

		var alerts tc.Alerts
		if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
			t.Errorf("could not decode alert response: %v\n", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Response was not successful %v\n", err)
		}
	}

	for _, cdn := range testData.CDNs {
		alerts, _, err := TOSession.DeleteCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot DELETE CDN by name: %v - %v\n", err, alerts)
		}
	}

}

func TestGetCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		alerts, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v\n", err, alerts)
		}
	}
}

func TestDeleteCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		alerts, _, err := TOSession.DeleteCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot DELETE CDN by name: %v - %v\n", err, alerts)
		}
	}
}
