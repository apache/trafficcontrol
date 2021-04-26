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

package v4

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestDeliveryServicesRegexes(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, DeliveryServicesRegexes}, func() {
		QueryDSRegexTest(t)
		QueryDSRegexTestIMS(t)
		CreateTestDSRegexWithMissingPattern(t)
	})
}

func CreateTestDeliveryServicesRegexes(t *testing.T) {
	db, err := OpenConnection()
	if err != nil {
		t.Fatal("cannot open db")
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Errorf("unable to close connection to db, error: %v", err)
		}
	}()

	dbRegexInsertTemplate := "INSERT INTO regex (pattern, type) VALUES ('%v', '%v');"
	dbRegexQueryTemplate := "SELECT id FROM regex order by id desc limit 1;"
	dbDSRegexInsertTemplate := "INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES ('%v', '%v', '%v');"

	for i, regex := range testData.DeliveryServicesRegexes {
		loadDSRegexIDs(t, &regex)

		err = execSQL(db, fmt.Sprintf(dbRegexInsertTemplate, regex.Pattern, regex.Type))
		if err != nil {
			t.Fatalf("unable to create regex: %v", err)
		}

		row := db.QueryRow(dbRegexQueryTemplate)
		err = row.Scan(&regex.ID)
		if err != nil {
			t.Fatalf("unable to query regex: %v", err)
		}

		err = execSQL(db, fmt.Sprintf(dbDSRegexInsertTemplate, regex.DSID, regex.ID, regex.SetNumber))
		if err != nil {
			t.Fatalf("unable to create ds regex %v", err)
		}

		testData.DeliveryServicesRegexes[i] = regex
	}
}

func DeleteTestDeliveryServicesRegexes(t *testing.T) {
	db, err := OpenConnection()
	if err != nil {
		t.Fatal("cannot open db")
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Errorf("unable to close connection to db, error: %v", err)
		}
	}()

	for _, regex := range testData.DeliveryServicesRegexes {
		err = execSQL(db, fmt.Sprintf("DELETE FROM deliveryservice_regex WHERE deliveryservice = '%v' and regex ='%v';", regex.DSID, regex.ID))
		if err != nil {
			t.Fatalf("unable to delete deliveryservice_regex by regex %v and ds %v: %v", regex.ID, regex.DSID, err)
		}

		err := execSQL(db, fmt.Sprintf("DELETE FROM regex WHERE Id = '%v';", regex.ID))
		if err != nil {
			t.Fatalf("unable to delete regex %v: %v", regex.ID, err)
		}
	}
}

// Note that this test is liable to break if the structure of the Delivery
// Service Regular Expressions in the test data is changed at all.
func CreateTestDSRegexWithMissingPattern(t *testing.T) {
	if len(testData.DeliveryServicesRegexes) < 4 {
		t.Fatal("Need at least 4 Delivery Service Regular Expressions to test creating a Delivery Service Regular Expression with a missing pattern")
	}
	var regex = testData.DeliveryServicesRegexes[3]
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", regex.DSName)
	ds, _, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Fatalf("unable to get ds '%s': %v - alerts: %+v", regex.DSName, err, ds.Alerts)
	}
	if len(ds.Response) == 0 {
		t.Fatalf("unable to get ds %v", regex.DSName)
	}

	var dsID int
	if ds.Response[0].ID == nil {
		t.Fatal("ds has a nil id")
	} else {
		dsID = *ds.Response[0].ID
	}

	regexPost := tc.DeliveryServiceRegexPost{Type: regex.Type, SetNumber: regex.SetNumber, Pattern: regex.Pattern}

	_, reqInfo, _ := TOSession.PostDeliveryServiceRegexesByDSID(dsID, regexPost, client.RequestOptions{})
	if reqInfo.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected: %v, but got: %v", http.StatusBadRequest, reqInfo.StatusCode)
	}
}

func loadDSRegexIDs(t *testing.T, test *tc.DeliveryServiceRegexesTest) {
	if test == nil {
		t.Error("loadDSRegexIDs called with nil test")
		return
	}
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", test.TypeName)
	dsTypes, _, err := TOSession.GetTypes(opts)
	if err != nil {
		t.Errorf("unable to get Types filtered by name '%s': %v - alerts: %+v", test.TypeName, err, dsTypes.Alerts)
		return
	}
	if len(dsTypes.Response) < 1 {
		t.Errorf("could not find any types by name '%s'", test.TypeName)
		return
	}
	test.Type = dsTypes.Response[0].ID

	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", test.DSName)
	dses, _, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Errorf("unable to ds by xmlid %v: %v - alerts: %+v", test.DSName, err, dses.Alerts)
		return
	}
	if len(dses.Response) != 1 {
		t.Errorf("unable to find ds by xmlid %v", test.DSName)
		return
	}
	if dses.Response[0].ID == nil {
		t.Error("Delivery Service had a null or undefined ID")
		return
	}
	test.DSID = *dses.Response[0].ID
}

func QueryDSRegexTestIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	opts.QueryParameters.Set("xmlId", "ds1")
	resp, reqInf, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Delivery Services filtered by XMLID 'ds1': %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

// Note that this test will break if the Delivery Service in the test data with
// the XMLID 'ds1' is removed or altered such that its regular expressions are
// different than at the time of this writing.
func QueryDSRegexTest(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", "ds1")
	ds, _, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Fatalf("unable to get ds 'ds1': %v - alerts: %+v", err, ds.Alerts)
	}
	if len(ds.Response) == 0 {
		t.Fatal("unable to get ds ds1")
	}
	var dsID int
	if ds.Response[0].ID == nil {
		t.Fatal("ds has a nil id")
	} else {
		dsID = *ds.Response[0].ID
	}

	dsRegexes, _, err := TOSession.GetDeliveryServiceRegexesByDSID(dsID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error fetching Regular Expressions for Delivery Service 'ds1' (#%d): %v - alerts: %+v", dsID, err, dsRegexes.Alerts)
	}
	if len(dsRegexes.Response) != 4 {
		t.Fatalf("expected to get 4 Regular Expressions for Delivery Service 'ds1' (#%d), got: %d", dsID, len(dsRegexes.Response))
	}
	regExpID := dsRegexes.Response[0].ID
	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("id", strconv.Itoa(regExpID))
	dsRegexes, _, err = TOSession.GetDeliveryServiceRegexesByDSID(dsID, opts)
	if err != nil {
		t.Errorf("Unexpected error getting Regular Expressions for Delivery Service 'ds1' (#%d) filtered by Regular Expression ID %d: %v - alerts: %+v", dsID, regExpID, err, dsRegexes.Alerts)
	}
	if len(dsRegexes.Response) != 1 {
		t.Fatalf("expected to get 1 Regular Expression for Delivery Service 'ds1' (#%d) that has ID %d, got: %d", dsID, regExpID, len(dsRegexes.Response))
	}
}
