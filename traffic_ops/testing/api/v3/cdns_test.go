package v3

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
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"net/http"
	"testing"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestCDNs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters}, func() {
		GetTestCDNsIMS(t)
		currentTime := time.Now().Add(-1 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		UpdateTestCDNs(t)
		GetTestCDNs(t)
		GetTestCDNsIMSAfterChange(t, header)
	})
}

func GetTestCDNsIMSAfterChange(t *testing.T, header http.Header) {
	for _, cdn := range testData.CDNs {
		_, reqInf, err := TOSession.GetCDNByNameIMS(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, cdn := range testData.CDNs {
		_, reqInf, err := TOSession.GetCDNByNameIMS(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCDNsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	for _, cdn := range testData.CDNs {
		futureTime := time.Now().AddDate(0,0,1)
		time := futureTime.Format(time.RFC1123)
		header.Set(rfc.IfModifiedSince, time)
		_, reqInf, err := TOSession.GetCDNByNameIMS(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.CreateCDN(cdn)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE cdns: %v", err)
		}
	}

}

func UpdateTestCDNs(t *testing.T) {

	firstCDN := testData.CDNs[0]
	// Retrieve the CDN by name so we can get the id for the Update
	resp, _, err := TOSession.GetCDNByName(firstCDN.Name)
	if err != nil {
		t.Errorf("cannot GET CDN by name: '%s', %v", firstCDN.Name, err)
	}
	remoteCDN := resp[0]
	expectedCDNDomain := "domain2"
	remoteCDN.DomainName = expectedCDNDomain
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateCDNByID(remoteCDN.ID, remoteCDN)
	if err != nil {
		t.Errorf("cannot UPDATE CDN by id: %v - %v", err, alert)
	}

	// Retrieve the CDN to check CDN name got updated
	resp, _, err = TOSession.GetCDNByID(remoteCDN.ID)
	if err != nil {
		t.Errorf("cannot GET CDN by name: '$%s', %v", firstCDN.Name, err)
	}
	respCDN := resp[0]
	if respCDN.DomainName != expectedCDNDomain {
		t.Errorf("results do not match actual: %s, expected: %s", respCDN.DomainName, expectedCDNDomain)
	}

}

func GetTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		resp, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v", err, resp)
		}
	}
}

func DeleteTestCDNs(t *testing.T) {

	for _, cdn := range testData.CDNs {
		// Retrieve the CDN by name so we can get the id for the Update
		resp, _, err := TOSession.GetCDNByName(cdn.Name)
		if err != nil {
			t.Errorf("cannot GET CDN by name: %v - %v", cdn.Name, err)
		}
		if len(resp) > 0 {
			respCDN := resp[0]

			_, _, err := TOSession.DeleteCDNByID(respCDN.ID)
			if err != nil {
				t.Errorf("cannot DELETE CDN by name: '%s' %v", respCDN.Name, err)
			}

			// Retrieve the CDN to see if it got deleted
			cdns, _, err := TOSession.GetCDNByName(cdn.Name)
			if err != nil {
				t.Errorf("error deleting CDN name: %s", err.Error())
			}
			if len(cdns) > 0 {
				t.Errorf("expected CDN name: %s to be deleted", cdn.Name)
			}
		}
	}
}
