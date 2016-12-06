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
	"testing"
	"time"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestStatsSummaryAll(t *testing.T) {

	uri := fmt.Sprintf("/api/1.2/stats_summary.json")
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiStatsSummaryRes traffic_ops.StatsSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiStatsSummaryRes); err != nil {
		t.Errorf("Could not decode stats summary json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiStatsSummary := apiStatsSummaryRes.Response

	clientStatsSummary, err := to.SummaryStats("", "", "")
	if err != nil {
		t.Errorf("Could not get stats summary from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiStatsSummary) != len(clientStatsSummary) {
		t.Errorf("Stats Summary Response Length -- expected %v, got %v\n", len(apiStatsSummary), len(clientStatsSummary))
	}

	for _, apiSs := range apiStatsSummary {
		match := false
		for _, clientSs := range clientStatsSummary {
			if apiSs == clientSs {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a stats summary matching %+v\n", apiSs)
		}
	}
}

func TestStatsSummarybyCDN(t *testing.T) {
	cdn, err := GetCdn()
	if err != nil {
		t.Error("Could not get a CDN, response was %v\n", err)
		t.FailNow()
	}
	uri := fmt.Sprintf("/api/1.2/stats_summary.json?cdnName=%s", cdn.Name)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiStatsSummaryRes traffic_ops.StatsSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiStatsSummaryRes); err != nil {
		t.Errorf("Could not decode stats summary json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiStatsSummary := apiStatsSummaryRes.Response

	clientStatsSummary, err := to.SummaryStats(cdn.Name, "", "")
	if err != nil {
		t.Errorf("Could not get stats summary from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiStatsSummary) != len(clientStatsSummary) {
		t.Errorf("Stats Summary Response Length -- expected %v, got %v\n", len(apiStatsSummary), len(clientStatsSummary))
	}

	for _, apiSs := range apiStatsSummary {
		match := false
		for _, clientSs := range clientStatsSummary {
			if apiSs == clientSs {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a stats summary matching %+v\n", apiSs)
		}
	}
}

func TestStatsSummaryByDs(t *testing.T) {
	ds, err := GetDeliveryService("")
	if err != nil {
		t.Error("Could not get a DS, response was %v\n", err)
		t.FailNow()
	}
	uri := fmt.Sprintf("/api/1.2/stats_summary.json?deliveryServiceName=%s", ds.XMLID)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiStatsSummaryRes traffic_ops.StatsSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiStatsSummaryRes); err != nil {
		t.Errorf("Could not decode stats summary json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiStatsSummary := apiStatsSummaryRes.Response

	clientStatsSummary, err := to.SummaryStats("", ds.XMLID, "")
	if err != nil {
		t.Errorf("Could not get stats summary from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiStatsSummary) != len(clientStatsSummary) {
		t.Errorf("Stats Summary Response Length -- expected %v, got %v\n", len(apiStatsSummary), len(clientStatsSummary))
	}

	for _, apiSs := range apiStatsSummary {
		match := false
		for _, clientSs := range clientStatsSummary {
			if apiSs == clientSs {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a stats summary matching %+v\n", apiSs)
		}
	}
}

func TestStatsSummaryByStatName(t *testing.T) {
	uri := fmt.Sprintf("/api/1.2/stats_summary.json?statName=daily_bytesserved")
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiStatsSummaryRes traffic_ops.StatsSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiStatsSummaryRes); err != nil {
		t.Errorf("Could not decode stats summary json.  Error is: %v\n", err)
		t.FailNow()
	}
	apiStatsSummary := apiStatsSummaryRes.Response

	clientStatsSummary, err := to.SummaryStats("", "", "daily_bytesserved")
	if err != nil {
		t.Errorf("Could not get stats summary from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiStatsSummary) != len(clientStatsSummary) {
		t.Errorf("Stats Summary Response Length -- expected %v, got %v\n", len(apiStatsSummary), len(clientStatsSummary))
	}

	for _, apiSs := range apiStatsSummary {
		match := false
		for _, clientSs := range clientStatsSummary {
			if apiSs == clientSs {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a stats summary matching %+v\n", apiSs)
		}
	}
}

func TestAddSummaryStats(t *testing.T) {
	cdn, err := GetCdn()
	if err != nil {
		t.Errorf("Could not get a CDN, response was %v\n", err)
		t.FailNow()
	}
	ds, err := GetDeliveryService(cdn.Name)
	if err != nil {
		t.Errorf("Could not get a DS, response was %v\n", err)
		t.FailNow()
	}
	now := time.Now()
	summaryTime := now.Format(time.RFC3339)
	statDate := now.Format("2006-01-02")

	testStatsSummay := new(traffic_ops.StatsSummary)
	testStatsSummay.CDNName = cdn.Name
	testStatsSummay.DeliveryService = ds.XMLID
	testStatsSummay.StatDate = statDate
	testStatsSummay.StatName = "testStatName"
	testStatsSummay.StatValue = "1234"
	testStatsSummay.SummaryTime = summaryTime

	err = to.AddSummaryStats(*testStatsSummay)
	if err != nil {
		t.Errorf("Could not add Summary Stats, response was %v\n", err)
		t.FailNow()
	}

	ssRes, err := to.SummaryStats(testStatsSummay.CDNName, testStatsSummay.DeliveryService, testStatsSummay.StatName)
	if err != nil {
		t.Errorf("Could not get a SummaryStats, error was: %v\n", err)
		t.FailNow()
	}
	match := false
	for _, ss := range ssRes {
		if ss.CDNName == testStatsSummay.CDNName &&
			ss.DeliveryService == testStatsSummay.DeliveryService &&
			ss.StatDate == testStatsSummay.StatDate &&
			ss.StatName == testStatsSummay.StatName &&
			ss.StatValue == testStatsSummay.StatValue {
			match = true
		}
	}
	if !match {
		t.Errorf("Stats Summary not found in Traffic Ops after Adding.  Summary Stats Response was: %+v, expecting: %+v\n", ssRes, *testStatsSummay)
	}
}
