package request

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"testing"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TODeliveryServiceRequest{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("DeliveryServiceRequest must be Creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("DeliveryServiceRequest must be Reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("DeliveryServiceRequest must be Updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("DeliveryServiceRequest must be Deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("DeliveryServiceRequest must be Identifier")
	}
	if _, ok := i.(api.Tenantable); !ok {
		t.Errorf("DeliveryServiceRequest must be Tenantable")
	}
}

func TestGetDeliveryServiceRequest(t *testing.T) {
	s := "this is not a valid xmlid.  Bad characters and too long."
	i := 1
	b := true
	u := "UPDATE"
	st := tc.RequestStatusSubmitted
	ds := tc.DeliveryServiceNullable{}
	ds.XMLID = &s
	ds.CDNID = &i
	ds.LogsEnabled = &b
	ds.DSCP = nil
	ds.GeoLimit = &i
	ds.Active = &b
	ds.TypeID = &i
	r := &TODeliveryServiceRequest{DeliveryServiceRequestNullable: tc.DeliveryServiceRequestNullable{
		ChangeType:      &u,
		Status:          &st,
		DeliveryService: &ds,
	}}

	expectedErrors := []string{
		/*
			`'regionalGeoBlocking' is required`,
			`'xmlId' cannot contain spaces`,
			`'dscp' is required`,
			`'displayName' cannot be blank`,
			`'geoProvider' is required`,
			`'typeId' is required`,
		*/
	}

	r.SetKeys(map[string]interface{}{"id": 10})
	keys, _ := r.GetKeys()
	if keys["id"].(int) != 10 {
		t.Errorf("expected ID to be %d,  not %d", 10, keys["id"].(int))
	}
	exp := "10"
	if s != r.GetAuditName() {
		t.Errorf("expected AuditName to be '%s',  not '%s'", s, r.GetAuditName())
	}
	exp = "deliveryservice_request"
	if r.GetType() != "deliveryservice_request" {
		t.Errorf("expected Type to be %s,  not %s", exp, r.GetType())
	}

	var errs []error
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	/* TODO: this section panics when deliveryservice.Validate() tries to get the type name.
	q := `insert into type (name, description, use_in_table) values ('HTTP', 'HTTP Content Routing', 'deliveryservice') ON CONFLICT (name) DO NOTHING;`
	qe := `insert into type \(name, description, use_in_table\) values \('HTTP', 'HTTP Content Routing', 'deliveryservice'\) ON CONFLICT \(name\) DO NOTHING;`
	mock.ExpectExec(qe).WillReturnResult(sqlmock.NewResult(1, 1))
	res, err := db.Exec(q)

	// db.Exec(`insert into type (name, description, use_in_table) values ('HTTP_NO_CACHE', 'HTTP Content Routing, no caching', 'deliveryservice') ON CONFLICT (name) DO NOTHING;`)
	//db.Exec(`insert into type (name, description, use_in_table) values ('HTTP_LIVE', 'HTTP Content routing cache in RAM', 'deliveryservice') ON CONFLICT (name) DO NOTHING;`)
	if err != nil {
		t.Error(err)
	}
	mock.ExpectQuery(`SELECT name from type where id=\$1`).WillReturnRows(sqlmock.NewRows([]string{"name"}))

	errs := r.Validate(db)
	*/
	if len(errs) != len(expectedErrors) {
		for _, e := range errs {
			t.Error(e)
		}
	}

	for e := range expectedErrors {
		t.Error(e)
	}

	/*
		if r.Update(db *sqlx.Tx, ctx context.Context) {
			t.Errorf("expected ID to be %d,  not %d", 10, r.GetID())
		}
		if r.Insert(db *sqlx.Tx, ctx context.Context) {
			t.Errorf("expected ID to be %d,  not %d", 10, r.GetID())
		}
		if r.Delete(db *sqlx.Tx, ctx context.Context) {
			t.Errorf("expected ID to be %d,  not %d", 10, r.GetID())
		}
	*/
}
