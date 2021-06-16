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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
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
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, typ := range testData.Types {
		opts.QueryParameters.Set("name", typ.Name)
		resp, reqInf, err := TOSession.GetTypes(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, timeStr)

	for _, typ := range testData.Types {
		opts.QueryParameters.Set("name", typ.Name)
		resp, reqInf, err := TOSession.GetTypes(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestTypesIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)

	for _, typ := range testData.Types {
		opts.QueryParameters.Set("name", typ.Name)
		resp, reqInf, err := TOSession.GetTypes(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestTypes(t *testing.T) {
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
	dbQueryTemplate := "INSERT INTO type (name, description, use_in_table) VALUES ('%s', '%s', '%s');"

	opts := client.NewRequestOptions()
	for _, typ := range testData.Types {
		opts.QueryParameters.Set("name", typ.Name)
		foundTypes, _, err := TOSession.GetTypes(opts)
		if err == nil && len(foundTypes.Response) > 0 {
			t.Logf("Type %v already exists (%v match(es))", typ.Name, len(foundTypes.Response))
			continue
		}

		var alerts tc.Alerts
		if typ.UseInTable != "server" {
			err = execSQL(db, fmt.Sprintf(dbQueryTemplate, typ.Name, typ.Description, typ.UseInTable))
		} else {
			alerts, _, err = TOSession.CreateType(typ, client.RequestOptions{})
		}

		if err != nil {
			t.Fatalf("could not create Type: %v - alerts: %+v", err, alerts.Alerts)
		}
	}
}

func SortTestTypes(t *testing.T) {
	resp, _, err := TOSession.GetTypes(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}

	sortedList := make([]string, 0, len(resp.Response))
	for _, typ := range resp.Response {
		sortedList = append(sortedList, typ.Name)
	}

	if !sort.StringsAreSorted(sortedList) {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestTypes(t *testing.T) {

	for i, typ := range testData.Types {
		expectedTypeName := fmt.Sprintf("testType%v", i)
		originalType := typ

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("name", originalType.Name)
		resp, _, err := TOSession.GetTypes(opts)
		if err != nil {
			t.Fatalf("cannot get Types filtered by name '%s': %v - alerts: %+v", originalType.Name, err, resp.Alerts)
		}
		if len(resp.Response) < 1 {
			t.Fatalf("no Types exist by name '%s'", originalType.Name)
		}

		remoteType := resp.Response[0]
		remoteType.Name = expectedTypeName
		// Ensure TO checks DB for UseInTable value
		remoteType.UseInTable = "server"

		alert, _, err := TOSession.UpdateType(remoteType.ID, remoteType, client.RequestOptions{})
		if originalType.UseInTable != "server" {
			if err == nil {
				t.Fatalf("expected update on Type #%d to fail", remoteType.ID)
			}
			continue
		} else if err != nil {
			t.Fatalf("cannot update Type: %v - alerts: %+v", err, alert.Alerts)
		}

		// Retrieve the Type to check Type name got updated
		opts.QueryParameters.Del("name")
		opts.QueryParameters.Set("id", strconv.Itoa(remoteType.ID))
		resp, _, err = TOSession.GetTypes(opts)
		opts.QueryParameters.Del("id")
		if err != nil {
			t.Fatalf("cannot get Type by ID %d: %v - alerts: %+v", originalType.ID, err, resp.Alerts)
		}
		respType := resp.Response[0]
		if respType.Name != expectedTypeName {
			t.Fatalf("results do not match actual: %s, expected: %s", respType.Name, expectedTypeName)
		}
		if respType.UseInTable != originalType.UseInTable {
			t.Fatalf("use in table should never be updated, got: %v, expected %v", respType.UseInTable, originalType.UseInTable)
		}

		// Revert name change
		respType.Name = originalType.Name
		alert, _, err = TOSession.UpdateType(respType.ID, respType, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot restore/update Type: %v - %+v", err, alert.Alerts)
		}
	}
}

func GetTestTypes(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, typ := range testData.Types {
		opts.QueryParameters.Set("name", typ.Name)
		resp, _, err := TOSession.GetTypes(opts)
		if err != nil {
			t.Errorf("cannot get Type: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func DeleteTestTypes(t *testing.T) {
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
	dbDeleteTemplate := "DELETE FROM type WHERE name='%s';"

	opts := client.NewRequestOptions()
	for _, typ := range testData.Types {
		// Retrieve the Type by name so we can get the id for the Update
		opts.QueryParameters.Set("name", typ.Name)
		resp, _, err := TOSession.GetTypes(opts)
		if err != nil || len(resp.Response) == 0 {
			t.Fatalf("cannot get Types filtered by name '%s': %v - alerts: %+v", typ.Name, err, resp.Alerts)
		}
		respType := resp.Response[0]

		if respType.UseInTable != "server" {
			err := execSQL(db, fmt.Sprintf(dbDeleteTemplate, respType.Name))
			if err != nil {
				t.Fatalf("cannot delete Type using database operations: %v", err)
			}
		} else {
			delResp, _, err := TOSession.DeleteType(respType.ID, client.RequestOptions{})
			if err != nil {
				t.Fatalf("cannot delete Type using the API: %v - alerts: %+v", err, delResp.Alerts)
			}
		}

		// Retrieve the Type to see if it got deleted
		types, _, err := TOSession.GetTypes(opts)
		if err != nil {
			t.Errorf("error fetching Types filtered by presumably deleted name: %v - alerts: %+v", err, types.Alerts)
		}
		if len(types.Response) > 0 {
			t.Errorf("expected Type '%s' to be deleted", typ.Name)
		}
	}
}
