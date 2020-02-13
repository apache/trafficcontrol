package v1

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

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestTypes(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Types}, func() {
		UpdateTestTypes(t)
		GetTestTypes(t)
	})
}

func CreateTestTypes(t *testing.T) {
	t.Log("---- CreateTestTypes ----")

	db, err := OpenConnection()
	if err != nil {
		t.Fatal("cannot open db")
	}
	defer func() {
		err := db.Close()
		if err != nil {
			t.Errorf("unable to close connection to db, error: %v", err.Error())
		}
	}()
	dbQueryTemplate := "INSERT INTO type (name, description, use_in_table) VALUES ('%v', '%v', '%v');"

	for _, typ := range testData.Types {
		foundTypes, _, err := TOSession.GetTypeByName(typ.Name)
		if err == nil && len(foundTypes) > 0 {
			t.Logf("Type %v already exists (%v match(es))", typ.Name, len(foundTypes))
			continue
		}

		if typ.UseInTable != "server" {
			err = execSQL(db, fmt.Sprintf(dbQueryTemplate, typ.Name, typ.Description, typ.UseInTable), "type")
		} else {
			_, _, err = TOSession.CreateType(typ)
		}

		if err != nil {
			t.Errorf("could not CREATE types: %v", err)
		}
	}

}

func UpdateTestTypes(t *testing.T) {
	t.Log("---- UpdateTestTypes ----")

	firstType := testData.Types[0]
	// Retrieve the Type by name so we can get the id for the Update
	resp, _, err := TOSession.GetTypeByName(firstType.Name)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v", firstType.Name, err)
	}
	remoteType := resp[0]
	expectedTypeName := "testType1"
	remoteType.Name = expectedTypeName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateTypeByID(remoteType.ID, remoteType)
	if err != nil {
		t.Errorf("cannot UPDATE Type by id: %v - %v", err, alert)
	}

	// Retrieve the Type to check Type name got updated
	resp, _, err = TOSession.GetTypeByID(remoteType.ID)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v", firstType.Name, err)
	}
	respType := resp[0]
	if respType.Name != expectedTypeName {
		t.Errorf("results do not match actual: %s, expected: %s", respType.Name, expectedTypeName)
	}

	t.Log("Response Type: ", respType)

	respType.Name = firstType.Name
	alert, _, err = TOSession.UpdateTypeByID(respType.ID, respType)
	if err != nil {
		t.Errorf("cannot restore UPDATE Type by id: %v - %v", err, alert)
	}
}

func GetTestTypes(t *testing.T) {
	t.Log("---- GetTestTypes ----")

	for _, typ := range testData.Types {
		resp, _, err := TOSession.GetTypeByName(typ.Name)
		if err != nil {
			t.Errorf("cannot GET Type by name: %v - %v", err, resp)

		}

		t.Log("Response: ", resp)
	}
}

func DeleteTestTypes(t *testing.T) {
	t.Log("---- DeleteTestTypes ----")

	for _, typ := range testData.Types {
		// Retrieve the Type by name so we can get the id for the Update
		resp, _, err := TOSession.GetTypeByName(typ.Name)
		if err != nil || len(resp) == 0 {
			t.Errorf("cannot GET Type by name: %v - %v", typ.Name, err)
		}
		respType := resp[0]

		delResp, _, err := TOSession.DeleteTypeByID(respType.ID)
		if err != nil {
			t.Errorf("cannot DELETE Type by name: %v - %v", err, delResp)
		}

		// Retrieve the Type to see if it got deleted
		types, _, err := TOSession.GetTypeByName(typ.Name)
		if err != nil {
			t.Errorf("error deleting Type name: %s", err.Error())
		}
		if len(types) > 0 {
			t.Errorf("expected Type name: %s to be deleted", typ.Name)
		}
	}
}
