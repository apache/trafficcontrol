// TODO: This comment with the `build nullable` tag must be removed when the structs are made consistent!!
// skip this unless specifically testing for nullable vs non-nullable struct comparison

// Run as `go test -tags nullable`
//go:build nullable
// +build nullable

package tc

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
	"reflect"
	"strings"
	"testing"
)

func TestNullStructs(t *testing.T) {
	compareWithNullable(t, ASN{}, ASNNullable{})
	compareWithNullable(t, CacheGroup{}, CacheGroupNullable{})
	compareWithNullable(t, CDN{}, CDNNullable{})
	compareWithNullable(t, Coordinate{}, CoordinateNullable{})
	compareWithNullable(t, DeliveryServiceRequestComment{}, DeliveryServiceRequestCommentNullable{})
	compareWithNullable(t, DeliveryServiceRequest{}, DeliveryServiceRequestNullable{})
	compareWithNullable(t, DeliveryService{}, DeliveryServiceNullable{})
	compareWithNullable(t, DeliveryServiceV11{}, DeliveryServiceNullableV11{})
	compareWithNullable(t, Division{}, DivisionNullable{})
	compareWithNullable(t, Domain{}, DomainNullable{})
	compareWithNullable(t, Parameter{}, ParameterNullable{})
	compareWithNullable(t, PhysLocation{}, PhysLocationNullable{})
	compareWithNullable(t, ProfileParameter{}, ProfileParameterNullable{})
	compareWithNullable(t, Profile{}, ProfileNullable{})
	compareWithNullable(t, Server{}, ServerNullable{})
	compareWithNullable(t, StaticDNSEntry{}, StaticDNSEntryNullable{})
	compareWithNullable(t, Status{}, StatusNullable{})
	compareWithNullable(t, SteeringTarget{}, SteeringTargetNullable{})
	compareWithNullable(t, Tenant{}, TenantNullable{})
	compareWithNullable(t, Type{}, TypeNullable{})

	// No Nullable version of these types
	//compareWithNullable(t, Federation{}, FederationNullable{})
	//compareWithNullable(t, ProfileParameters{}, ProfileParametersNullable{})
}

// compareFields checks that non-nullable and nullable versions have same fields
func compareWithNullable(t *testing.T, obj interface{}, nullObj interface{}) {
	ot := reflect.TypeOf(obj)
	nt := reflect.TypeOf(nullObj)
	if strings.Replace(nt.String(), "Nullable", "", 1) != ot.String() {
		t.Errorf("expected type %s with nullable %s", ot, nt)
	}

	if ot.NumField() != nt.NumField() {
		t.Errorf("%T has %d fields, but %T has %d", obj, ot.NumField(), nullObj, nt.NumField())
	}

	seen := make(map[string]struct{}, ot.NumField())

	for i := 0; i < ot.NumField(); i++ {
		oField := ot.Field(i)
		if oField.Anonymous {
			// embedded struct -- skip it
			continue
		}
		nField, ok := nt.FieldByName(oField.Name)
		if !ok {
			t.Errorf("field %s found on %T but not %T", oField.Name, obj, nullObj)
			continue
		}

		seen[nField.Name] = struct{}{}
		oKind := oField.Type.Kind()
		nKind := nField.Type.Kind()
		//t.Logf("%T.%s is %s. %T.%s is %s", obj, oField.Name, oKind.String(), nullObj, nField.Name, nKind.String())
		if oKind == nKind {
			continue
		}
		if nKind == reflect.Ptr && nField.Type.Elem().Kind() != oKind {
			t.Errorf("%T.%s (%s) and %T.%s (%s) have mismatched types", obj, oField.Name, oField.Type.String(), nullObj, nField.Name, nField.Type.String())
		}
	}

	// check for fields in Nullable version not in non-Nullable
	for i := 0; i < nt.NumField(); i++ {
		nField := nt.Field(i)
		if nField.Anonymous {
			// embedded struct -- skip it
			continue
		}
		if _, ok := seen[nField.Name]; ok {
			// already accounted for
			continue
		}

		t.Errorf("field %s found on %T but not %T", nField.Name, obj, nullObj)
	}
}
