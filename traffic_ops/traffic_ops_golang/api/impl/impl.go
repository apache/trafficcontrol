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

// impl was made to prevent a circular import between the tc and api
// packages. The TC structs are referenced by the api, so embedding
// something like api.IDImpl in the TC struct won't work.
//
// The APIInfoImpl references api.APIInfo, so it isn't included in
// this package.
//
// By making common implementations here we can define behavior by
// embedding it.

package impl

// IDImpl implements GetKeys and SetKeys for the Identifier interface
// They do not complete the interface. IDImpl makes the ID Nullable.
// Although natural keys are prefereable, the current most common case
// of the keys we use are for a single integer key.
//
type IDImpl struct {
	ID *int `json:"id" db:"id"`
}

func (val IDImpl) GetKeys() (map[string]interface{}, bool) {
	if val.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *val.ID}, true
}

func (val *IDImpl) SetKeys(keys map[string]interface{}) {
	id, _ := keys["id"].(int)
	val.ID = &id
}
