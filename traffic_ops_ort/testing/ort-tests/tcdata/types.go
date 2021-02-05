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
)

func (r *TCData) CreateTestTypes(t *testing.T) {
	t.Log("---- CreateTestTypes ----")

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
	dbQueryTemplate := "INSERT INTO type (name, description, use_in_table) VALUES ('%v', '%v', '%v');"

	for _, typ := range r.TestData.Types {
		foundTypes, _, err := TOSession.GetTypeByName(typ.Name)
		if err == nil && len(foundTypes) > 0 {
			t.Logf("Type %v already exists (%v match(es))", typ.Name, len(foundTypes))
			continue
		}

		if typ.UseInTable != "server" {
			err = execSQL(db, fmt.Sprintf(dbQueryTemplate, typ.Name, typ.Description, typ.UseInTable))
		} else {
			_, _, err = TOSession.CreateType(typ)
		}

		if err != nil {
			t.Fatalf("could not CREATE types: %v", err)
		}
	}
}

func (r *TCData) DeleteTestTypes(t *testing.T) {
	t.Log("---- DeleteTestTypes ----")

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
	dbDeleteTemplate := "DELETE FROM type WHERE name='%v';"

	for _, typ := range r.TestData.Types {
		// Retrieve the Type by name so we can get the id for the Update
		resp, _, err := TOSession.GetTypeByName(typ.Name)
		if err != nil || len(resp) == 0 {
			t.Fatalf("cannot GET Type by name: %v - %v", typ.Name, err)
		}
		respType := resp[0]

		if respType.UseInTable != "server" {
			err := execSQL(db, fmt.Sprintf(dbDeleteTemplate, respType.Name))
			if err != nil {
				t.Fatalf("cannot DELETE Type by name: %v", err)
			}
		} else {
			delResp, _, err := TOSession.DeleteTypeByID(respType.ID)
			if err != nil {
				t.Fatalf("cannot DELETE Type by name: %v - %v", err, delResp)
			}
		}

		// Retrieve the Type to see if it got deleted
		types, _, err := TOSession.GetTypeByName(typ.Name)
		if err != nil {
			t.Errorf("error deleting Type name: %v", err)
		}
		if len(types) > 0 {
			t.Errorf("expected Type name: %s to be deleted", typ.Name)
		}
	}
}
