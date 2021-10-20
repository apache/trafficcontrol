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

package v3

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
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

func CreateTestDSRegexWithMissingPattern(t *testing.T) {
	var regex = testData.DeliveryServicesRegexes[3]
	ds, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr(regex.DSName, nil)
	if err != nil {
		t.Fatalf("unable to get ds %v: %v", regex.DSName, err)
	}
	if len(ds) == 0 {
		t.Fatalf("unable to get ds %v", regex.DSName)
	}

	var dsID int
	if ds[0].ID == nil {
		t.Fatal("ds has a nil id")
	} else {
		dsID = *ds[0].ID
	}

	regexPost := tc.DeliveryServiceRegexPost{Type: regex.Type, SetNumber: regex.SetNumber, Pattern: regex.Pattern}

	_, reqInfo, _ := TOSession.PostDeliveryServiceRegexesByDSID(dsID, regexPost)
	if reqInfo.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected: %v, but got: %v", http.StatusBadRequest, reqInfo.StatusCode)
	}
}

func loadDSRegexIDs(t *testing.T, test *tc.DeliveryServiceRegexesTest) {
	dsTypes, _, err := TOSession.GetTypeByName(test.TypeName)
	if err != nil {
		t.Fatalf("unable to get type by name %v: %v", test.TypeName, err)
	}
	if len(dsTypes) < 1 {
		t.Fatalf("could not find any types by name %v", test.TypeName)
	}
	test.Type = dsTypes[0].ID

	dses, _, err := TOSession.GetDeliveryServiceByXMLIDNullable(test.DSName)
	if err != nil {
		t.Fatalf("unable to ds by xmlid %v: %v", test.DSName, err)
	}
	if len(dses) != 1 {
		t.Fatalf("unable to find ds by xmlid %v", test.DSName)
	}
	test.DSID = *dses[0].ID
}

func QueryDSRegexTestIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	_, reqInf, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr("ds1", header)
	if err != nil {
		t.Fatalf("could not GET delivery services regex: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func QueryDSRegexTest(t *testing.T) {
	ds, _, err := TOSession.GetDeliveryServiceByXMLIDNullable("ds1")
	if err != nil {
		t.Fatalf("unable to get ds ds1: %v", err)
	}
	if len(ds) == 0 {
		t.Fatal("unable to get ds ds1")
	}
	var dsID int
	if ds[0].ID == nil {
		t.Fatal("ds has a nil id")
	} else {
		dsID = *ds[0].ID
	}

	dsRegexes, _, err := TOSession.GetDeliveryServiceRegexesByDSID(dsID, nil)
	if err != nil {
		t.Fatal("unable to get ds_regex by id " + strconv.Itoa(dsID))
	}
	if len(dsRegexes) != 4 {
		t.Fatal("expected to get 4 ds_regex, got " + strconv.Itoa(len(dsRegexes)))
	}

	params := make(map[string]string)
	params["id"] = strconv.Itoa(dsRegexes[0].ID)
	dsRegexes, _, err = TOSession.GetDeliveryServiceRegexesByDSID(dsID, params)
	if err != nil {
		t.Fatalf("unable to get ds_regex by id %v with query param %v", dsID, params["id"])
	}
	if len(dsRegexes) != 1 {
		t.Fatal("expected to get 1 ds_regex, got " + strconv.Itoa(len(dsRegexes)))
	}
}
