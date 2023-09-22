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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestDowngrade(t *testing.T) {
	roleV4 := RoleV4{
		Name:        "rolev40_name",
		Permissions: []string{"perm1", "perm2"},
		Description: "rolev40_desc",
		LastUpdated: &time.Time{},
	}

	role := roleV4.Downgrade()
	if role.Name == nil {
		t.Errorf("role name became nil after downgrade")
	} else if *role.Name != roleV4.Name {
		t.Errorf("expected role names to be the same after downgrade, but the new role name %s doesn't match the old role name %s", *role.Name, roleV4.Name)
	}

	if role.Description == nil {
		t.Errorf("role description became nil after downgrade")
	} else if *role.Description != roleV4.Description {
		t.Errorf("expected role descriptions to be the same after downgrade, but the new role description %s doesn't match the old role description %s", *role.Description, roleV4.Description)
	}

	if role.Capabilities == nil {
		t.Errorf("role capabilities became nil after downgrade")
	} else if len(*role.Capabilities) != len(roleV4.Permissions) {
		t.Errorf("new role capabilities length %d after downgrade doesn't match old role permissions length %d", len(*role.Capabilities), len(roleV4.Permissions))
	} else {
		oldPermissions := make(map[string]struct{})
		for _, perm := range roleV4.Permissions {
			oldPermissions[perm] = struct{}{}
		}
		for _, cap := range *role.Capabilities {
			if _, ok := oldPermissions[cap]; !ok {
				t.Errorf("permission %s did not exist earlier, but is present after downgrade", cap)
			}
		}
	}
}

func TestUpgrade(t *testing.T) {
	role := Role{
		RoleV11: RoleV11{
			ID:          util.IntPtr(100),
			Name:        util.StrPtr("role_name"),
			Description: util.StrPtr("role_desc"),
			PrivLevel:   util.IntPtr(10),
		},
		Capabilities: &[]string{"cap1", "cap2"},
	}

	roleV4 := role.Upgrade()
	if roleV4.Name != *role.Name {
		t.Errorf("expected role names to be the same after upgrade, but the new role name %s doesn't match the old role name %s", roleV4.Name, *role.Name)
	}

	if *role.Description != roleV4.Description {
		t.Errorf("expected role descriptions to be the same after upgrade, but the new role description %s doesn't match the old role description %s", roleV4.Description, *role.Description)
	}

	if &roleV4 == nil {
		t.Errorf("role permissions became nil after upgrade")
	} else if len(roleV4.Permissions) != len(*role.Capabilities) {
		t.Errorf("new role permissions length %d after upgrade doesn't match old role capabilities length %d", len(roleV4.Permissions), len(*role.Capabilities))
	} else {
		oldCapabilities := make(map[string]struct{})
		for _, perm := range *role.Capabilities {
			oldCapabilities[perm] = struct{}{}
		}
		for _, perm := range roleV4.Permissions {
			if _, ok := oldCapabilities[perm]; !ok {
				t.Errorf("capability %s did not exist earlier, but is present after upgrade", perm)
			}
		}
	}
}
