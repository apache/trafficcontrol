package auth

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
	"strings"
	"testing"
)

func ExampleCurrentUser_Can() {
	cu := CurrentUser{}
	fmt.Println(cu.Can("anything"))
	// Output: false
}

func TestCurrentUser_Can(t *testing.T) {
	cu := CurrentUser{}
	cu.perms = map[string]struct{}{"do-something": {}}
	if !cu.Can("do-something") {
		t.Error("user cannot do something they have Permission to do")
	}
	if cu.Can("do-something-else") {
		t.Error("user can do something they don't have Permission to do")
	}
}

func ExampleCurrentUser_MissingPermissions() {
	cu := CurrentUser{}
	missingPerms := cu.MissingPermissions("do something", "do anything")
	fmt.Println(strings.Join(missingPerms, ", "))
	// Output: do something, do anything
}

func TestCurrentUser_MissingPermissions(t *testing.T) {
	cu := CurrentUser{}
	cu.perms = map[string]struct{}{"do-something": {}}
	missing := cu.MissingPermissions("do-something", "do-something-else")
	if len(missing) != 1 {
		t.Fatalf("Expected checking user with one Permission for two Permissions to be missing one, actually missing: %d", len(missing))
	}
	if missing[0] != "do-something-else" {
		t.Errorf("Expected user to be missing 'do-something-else' Permission, actually missing: %s", missing[0])
	}
}
