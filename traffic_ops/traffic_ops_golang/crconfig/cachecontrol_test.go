package crconfig

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
)

func TestETag(t *testing.T) {
	now := time.Now()
	e := ETag(now)

	eTime, err := TimeFromETag(e)
	if err != nil {
		t.Errorf("TimeFromETag(now) expected: no error, actual: %+v", err)
	}
	if !now.Equal(eTime) {
		t.Errorf("TimeFromETag(now) expected: %+v, actual: %+v %+v", now, e, eTime)
	}

	later := now.Add(time.Hour)
	eLater := ETag(later)
	eLaterTime, err := TimeFromETag(eLater)
	if err != nil {
		t.Errorf("TimeFromETag(later) expected: no error, actual: %+v", err)
	}
	if !later.Equal(eLaterTime) {
		t.Errorf("TimeFromETag(later) expected: %+v, actual: %+v", later, eLaterTime)
	}

	badTags := []string{
		``,
		`"`,
		`""`,
		`"v0-"`,
		`"v0"`,
		`"v1"`,
		`"v0-"`,
		`"v1-"`,
		`"v0-a"`,
		`"v0-asdf"`,
		ETag(time.Now().Add(time.Hour * 24 * 365 * 200)),
		ETag(time.Now().Add(time.Hour * 24 * 365 * 200 * -1)),
		`v0-asdf`, // this is a valid date, but it's 1970 - so the sanity check must catch it
		`v0-`,
	}
	for _, tag := range badTags {
		if badTime, err := TimeFromETag(tag); err == nil {
			t.Errorf("TimeFromETag("+tag+") expected: error, actual: nil time %+v", badTime)
		}
	}
}
