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
	"fmt"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestInvalidationJobGetTTL(t *testing.T) {
	job := InvalidationJob{
		Parameters: nil,
	}
	ttl := job.TTLHours()
	if ttl != 0 {
		t.Error("expected 0 when no parameters")
	}
	job.Parameters = util.StrPtr("TTL:24h,x:asdf")
	ttl = job.TTLHours()
	if ttl != 0 {
		t.Error("expected 0 when invalid parameters")
	}

	job.Parameters = util.StrPtr("TTL:24h")
	ttl = job.TTLHours()
	if ttl != 24 {
		t.Errorf("expected ttl to be 24, got %v", ttl)
	}
}

func ExampleInvalidationJobInput_TTLHours_duration() {
	j := InvalidationJobInput{nil, nil, nil, util.InterfacePtr("121m"), nil, nil}
	ttl, e := j.TTLHours()
	if e != nil {
		fmt.Printf("Error: %v\n", e)
	}
	fmt.Println(ttl)
	// Output: 2
}

func ExampleInvalidationJobInput_TTLHours_number() {
	j := InvalidationJobInput{nil, nil, nil, util.InterfacePtr(2.1), nil, nil}
	ttl, e := j.TTLHours()
	if e != nil {
		fmt.Printf("Error: %v\n", e)
	}
	fmt.Println(ttl)
	// Output: 2
}

func ExampleInvalidationJobV4_String() {
	t, _ := time.Parse(time.RFC3339, "2021-11-08T01:02:03Z")
	j := InvalidationJobV4{
		AssetURL:         "https://example.com/.*",
		CreatedBy:        "noone",
		DeliveryService:  "demo1",
		ID:               5,
		InvalidationType: REFETCH,
		StartTime:        t,
		TTLHours:         72,
	}

	fmt.Println(j)
	// Output: InvalidationJobV4{ID: 5, AssetURL: "https://example.com/.*", CreatedBy: "noone", DeliveryService: "demo1", TTLHours: 72, InvalidationType: "REFETCH", StartTime: "2021-11-08T01:02:03Z"}
}
