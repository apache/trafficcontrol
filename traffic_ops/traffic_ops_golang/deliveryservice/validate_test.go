package deliveryservice

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
	"fmt"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// Validate ...
func TestValidate(t *testing.T) {
	displayName := "this is gonna be a name that's a little longer than the limit of 48 chars.."

	ds := tc.DeliveryService{
		DisplayName: "this is gonna be a name that's a little longer than the limit of 48 chars..",
		MissLat:     -91.0,
		MissLong:    102.1,
	}
	expected := []string{
		`display name '` + displayName + `' can be no longer than 48 characters`,
		`missLat value -91 must not exceed +/- 90.0`,
		`missLong value 102 must not exceed +/- 90.0`,
		`xmlId is required`,
	}

	errors := Validate(ds)
	for _, e := range errors {
		fmt.Printf("%s\n", e)
	}

	if len(expected) != len(errors) {
		t.Errorf("Expected %d errors, got %d", len(expected), len(errors))
	}
}
