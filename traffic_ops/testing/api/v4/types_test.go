package v4

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
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestTypes(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Types}, func() {
		GetTestTypesIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		SortTestTypes(t)
		UpdateTestTypes(t)
		GetTestTypes(t)
		GetTestTypesIMSAfterChange(t, header)
	})
}

func GetTestTypesIMSAfterChange(t *testing.T, header http.Header) {
	for _, typ := range testData.Types {
		_, reqInf, err := TOSession.GetTypeByNameWithHdr(typ.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, typ := range testData.Types {
		_, reqInf, err := TOSession.GetTypeByNameWithHdr(typ.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestTypesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	t.Log("---- GetTestTypes ----")

	for _, typ := range testData.Types {
		_, reqInf, err := TOSession.GetTypeByNameWithHdr(typ.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
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
			t.Errorf("unable to close connection to db, error: %v", err)
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
			err = execSQL(db, fmt.Sprintf(dbQueryTemplate, typ.Name, typ.Description, typ.UseInTable))
		} else {
			_, _, err = TOSession.CreateType(typ)
		}

		if err != nil {
			t.Fatalf("could not CREATE types: %v", err)
		}
	}
}

func SortTestTypes(t *testing.T) {
	var header http.Header
	var sortedList []string
	resp, _, err := TOSession.GetTypesWithHdr(header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
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

func UpdateTestTypes(t *testing.T) {
	t.Log("---- UpdateTestTypes ----")

	for i, typ := range testData.Types {
		expectedTypeName := fmt.Sprintf("testType%v", i)
		originalType := typ
		resp, _, err := TOSession.GetTypeByName(originalType.Name)
		if err != nil {
			t.Fatalf("cannot GET Type by name: %v - %v", originalType.Name, err)
		}
		if len(resp) < 1 {
			t.Fatalf("no types by name: %v", originalType.Name)
		}

		remoteType := resp[0]
		remoteType.Name = expectedTypeName
		// Ensure TO checks DB for UseInTable value
		remoteType.UseInTable = "server"

		var alert tc.Alerts
		alert, _, err = TOSession.UpdateTypeByID(remoteType.ID, remoteType)
		if originalType.UseInTable != "server" {
			if err == nil {
				t.Fatalf("expected UPDATE on type %v to fail", remoteType.ID)
			}
			continue
		} else if err != nil {
			t.Fatalf("cannot UPDATE Type by id: %v - %v", err, alert)
		}

		// Retrieve the Type to check Type name got updated
		resp, _, err = TOSession.GetTypeByID(remoteType.ID)
		if err != nil {
			t.Fatalf("cannot GET Type by ID: %v - %v", originalType.ID, err)
		}
		respType := resp[0]
		if respType.Name != expectedTypeName {
			t.Fatalf("results do not match actual: %s, expected: %s", respType.Name, expectedTypeName)
		}
		if respType.UseInTable != originalType.UseInTable {
			t.Fatalf("use in table should never be updated, got: %v, expected %v", respType.UseInTable, originalType.UseInTable)
		}

		// Revert name change
		respType.Name = originalType.Name
		alert, _, err = TOSession.UpdateTypeByID(respType.ID, respType)
		if err != nil {
			t.Fatalf("cannot restore UPDATE Type by id: %v - %v", err, alert)
		}
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
	dbDeleteTemplate := "DELETE FROM type WHERE name='%v';"

	for _, typ := range testData.Types {
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
