package crconfig

import (
	"context"
	"encoding/json"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"strings"
	"testing"
	"time"
)

func randTopology() tc.CRConfigTopology {
	return tc.CRConfigTopology{
		Nodes: randStrArray(),
	}
}

func ExpectedMakeTops() map[string]tc.CRConfigTopology {
	return map[string]tc.CRConfigTopology{
		"top1": randTopology(),
		"top2": randTopology(),
	}
}

func MockMakeTops(mock sqlmock.Sqlmock, expected map[string]tc.CRConfigTopology) {
	rows := sqlmock.NewRows([]string{
		"name",
		"nodes"})

	for topName, top := range expected {
		nodes := "{" + strings.Join(top.Nodes, ",") + "}"
		rows = rows.AddRow(
			topName,
			nodes)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
}

func TestMakeTops(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	expected := ExpectedMakeTops()
	mock.ExpectBegin()
	MockMakeTops(mock, expected)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatal("creating transaction: ", err)
	}

	actual, err := makeTopologies(tx)
	if err != nil {
		t.Fatal("makeTopologies expected: nil error, actual: ", err)
	}

	if err = db.Close(); err != nil {
		t.Fatal("closing db: ", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("makeTopologies len expected: %v, actual: %v", len(expected), len(actual))
	}

	for topName, top := range expected {
		actualTop, ok := actual[topName]
		if !ok {
			t.Errorf("makeTopologies expected: %v, actual: missing", topName)
			continue
		}
		expectedBts, _ := json.MarshalIndent(top, " ", " ")
		actualBts, _ := json.MarshalIndent(actualTop, " ", " ")
		if !reflect.DeepEqual(expectedBts, actualBts) {
			t.Errorf("makeDSes ds %+v expected: %+v\n\nactual: %+v\n\n\n", topName, string(expectedBts), string(actualBts))
		}
	}
}
