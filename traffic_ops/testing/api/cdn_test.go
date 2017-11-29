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
	"testing"

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func TestCDNs(t *testing.T) {

	TestCreateCDNs(t)
	TestUpdateCDNs(t)
	TestGetCDNs(t)
	TestDeleteCDNs(t)

}

func TestCreateCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.CreateCDN(cdn)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATe cdns: %v\n", err)
		}
	}

}

func TestUpdateCDNs(t *testing.T) {

	firstCDN := testData.CDNs[0]
	// Retrieve the CDN by name so we can get the id for the Update
	resp, _, err := TOSession.GetCDNByName(firstCDN.Name)
	if err != nil {
		t.Errorf("cannot GET CDN by name: %v - %v\n", err)
	}
	remoteCDN := resp[0]
	expectedCDNName := "testCdn1"
	remoteCDN.Name = expectedCDNName
	alerts, _, err := TOSession.UpdateCDNByID(remoteCDN.ID, remoteCDN)
	if err != nil {
		t.Errorf("cannot UPDATE CDN by id: %v - %v\n", err, alerts)
	}

	// Retrieve the CDN to check CDN name got updated
	resp, _, err = TOSession.GetCDNByID(remoteCDN.ID)
	if err != nil {
		t.Errorf("cannot GET CDN by name: %v - %v\n", err)
	}
	respCDN := resp[0]
	if respCDN.Name != expectedCDNName {
		t.Errorf("results do not match actual: %s, expected: %s\n", respCDN.Name, expectedCDNName)
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

	secondCDN := testData.CDNs[1]
	resp, _, err := TOSession.DeleteCDNByName(secondCDN.Name)
	if err != nil {
		t.Errorf("cannot DELETE CDN by name: %v - %v\n", err, resp)
	}

	// Retrieve the CDN to see if it got deleted
	cdns, _, err := TOSession.GetCDNByName(secondCDN.Name)
	if len(cdns) > 0 {
		t.Errorf("expected CDN name: %s to be deleted\n", secondCDN.Name)
	}
}
