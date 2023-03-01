package tcdata

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
	"testing"

	"github.com/apache/trafficcontrol/v7/lib/go-tc"
)

func (r *TCData) CreateTestDeliveryServicesRegexes(t *testing.T) {
	db, err := r.OpenConnection()
	if err != nil {
		t.Fatal("cannot open db")
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Errorf("unable to close connection to db, error: %v", err)
		}
	}()

	dbRegexInsertTemplate := "INSERT INTO regex (pattern, type) VALUES ('%s', %d);"
	dbRegexQueryTemplate := "SELECT id FROM regex order by id desc limit 1;"
	dbDSRegexInsertTemplate := "INSERT INTO deliveryservice_regex (deliveryservice, regex, set_number) VALUES (%d, %d, %d);"

	for i, regex := range r.TestData.DeliveryServicesRegexes {
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

		r.TestData.DeliveryServicesRegexes[i] = regex
	}
}

func loadDSRegexIDs(t *testing.T, test *tc.DeliveryServiceRegexesTest) {
	dsTypes, _, err := TOSession.GetTypeByName(test.TypeName)
	if err != nil {
		t.Fatalf("unable to get type by name %s: %v", test.TypeName, err)
	}
	if len(dsTypes) < 1 {
		t.Fatalf("could not find any types by name %s", test.TypeName)
	}
	test.Type = dsTypes[0].ID

	dses, _, err := TOSession.GetDeliveryServiceByXMLIDNullable(test.DSName)
	if err != nil {
		t.Fatalf("unable to ds by xmlid %s: %v", test.DSName, err)
	}
	if len(dses) != 1 {
		t.Fatalf("unable to find ds by xmlid %s", test.DSName)
	}
	test.DSID = *dses[0].ID
}

func (r *TCData) DeleteTestDeliveryServicesRegexes(t *testing.T) {
	db, err := r.OpenConnection()
	if err != nil {
		t.Fatal("cannot open db")
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Errorf("unable to close connection to db, error: %v", err)
		}
	}()

	for _, regex := range r.TestData.DeliveryServicesRegexes {
		err = execSQL(db, fmt.Sprintf("DELETE FROM deliveryservice_regex WHERE deliveryservice = '%v' and regex ='%v';", regex.DSID, regex.ID))
		if err != nil {
			t.Fatalf("unable to delete deliveryservice_regex by regex %d and ds %d: %v", regex.ID, regex.DSID, err)
		}

		err := execSQL(db, fmt.Sprintf("DELETE FROM regex WHERE id = %d;", regex.ID))
		if err != nil {
			t.Fatalf("unable to delete regex %d: %v", regex.ID, err)
		}
	}
}
