package v2

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
	"sort"
	"testing"

	tc "github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestCDNs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters}, func() {
		SortTestCDNs(t)
		UpdateTestCDNs(t)
		GetTestCDNs(t)
	})
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

func SortTestCDNs(t *testing.T) {
	var sortedList []string
	resp, _, _ := TOSession.GetCDNs()
	for i, _ := range resp {
		sortedList = append(sortedList, resp[i].Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
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
