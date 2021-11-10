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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-util"
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

func TestInvalidationJobLegacy(t *testing.T) {
	startStr := `2009-11-10 23:00:00+00`
	start, err := time.Parse(JobLegacyTimeFormat, startStr)
	if err != nil {
		t.Fatalf("failed to parse constant test time '" + startStr + "': " + err.Error())
	}

	type Expecteds struct {
		TestName string
		Input    InvalidationJobV4
		Output   string
	}
	expecteds := []Expecteds{
		{
			"basic",
			InvalidationJobV4{42, "http://foo.com/bar", "user0", "ds0", 24, REFRESH, start},
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","ttlHours":24,"invalidationType":"REFRESH","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"special_url_characters",
			InvalidationJobV4{999, "http://foo.com/bar(.*@#$%^*(", "user0", "ds0", 24, REFRESH, start},
			`{"startTime":"` + startStr + `","id":999,"assetUrl":"http://foo.com/bar(.*@#$%^*(","createdBy":"user0","deliveryService":"ds0","ttlHours":24,"invalidationType":"REFRESH","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"special_url_characters",
			InvalidationJobV4{999, "http://foo.com/bar(.*@#$%^*(", "user0", "ds0", 24, REFRESH, start},
			`{"startTime":"` + startStr + `","id":999,"assetUrl":"http://foo.com/bar(.*@#$%^*(","createdBy":"user0","deliveryService":"ds0","ttlHours":24,"invalidationType":"REFRESH","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"refetch",
			InvalidationJobV4{999, "http://foo.com/bar(.*@#$%^*(", "user0", "ds0", 24, REFETCH, start},
			`{"startTime":"` + startStr + `","id":999,"assetUrl":"http://foo.com/bar(.*@#$%^*(##REFETCH##","createdBy":"user0","deliveryService":"ds0","ttlHours":24,"invalidationType":"REFETCH","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
	}

	for _, expected := range expecteds {
		t.Run(expected.TestName, func(t *testing.T) {
			bts, err := json.Marshal(expected.Input)
			if err != nil {
				t.Errorf("expected output '%v' actual error: %v", expected.Output, err)
			}
			actual := string(bts)
			if actual != expected.Output {
				t.Errorf("expected output '%v' actual '%v'", expected.Output, actual)
			}

			obj := InvalidationJobV4{}
			if err := json.Unmarshal([]byte(expected.Output), &obj); err != nil {
				t.Errorf("expected output '%v' to marshal, actual error: %v", expected.Output, err)
			}
			if !reflect.DeepEqual(obj, expected.Input) {
				t.Errorf("expected input to create '%+v' actual '%+v'", expected.Input, obj)
			}
		})
	}

	type ExpectedInputs struct {
		TestName string
		Output   InvalidationJobV4
		Input    string
	}

	expectedLegacyInputs := []ExpectedInputs{
		{
			"legacy_input_basic",
			InvalidationJobV4{42, "http://foo.com/bar", "user0", "ds0", 24, REFRESH, start},
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"legacy_and_new_input",
			InvalidationJobV4{42, "http://foo.com/bar", "user0", "ds0", 24, REFETCH, start},
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar##REFETCH##","createdBy":"user0","deliveryService":"ds0","ttlHours":24,"invalidationType":"REFETCH","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"legacy_input_refetch",
			InvalidationJobV4{999, "http://foo.com/bar", "user0", "ds0", 24, REFETCH, start},
			`{"startTime":"` + startStr + `","id":999,"assetUrl":"http://foo.com/bar##REFETCH##","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"special_url_characters",
			InvalidationJobV4{999, "http://foo.com/bar(.*@#$%^*(", "user0", "ds0", 24, REFRESH, start},
			`{"startTime":"` + startStr + `","id":999,"assetUrl":"http://foo.com/bar(.*@#$%^*(","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
	}
	for _, expected := range expectedLegacyInputs {
		t.Run(expected.TestName, func(t *testing.T) {
			obj := InvalidationJobV4{}
			if err := json.Unmarshal([]byte(expected.Input), &obj); err != nil {
				t.Errorf("expected legacy input '%v' to marshal, actual error: %v", expected.Input, err)
			}
			if !reflect.DeepEqual(obj, expected.Output) {
				t.Errorf("expected legacy input to create '%+v' actual '%+v'", expected.Output, obj)
			}
		})
	}

	expectedNewInputs := []ExpectedInputs{
		{
			"new_input_basic",
			InvalidationJobV4{42, "http://foo.com/bar", "user0", "ds0", 24, REFRESH, start},
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","ttlHours":24,"invalidationType":"REFRESH"}`,
		},
		{
			"legacy_and_new_input",
			InvalidationJobV4{42, "http://foo.com/bar", "user0", "ds0", 24, REFETCH, start},
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar##REFETCH##","createdBy":"user0","deliveryService":"ds0","ttlHours":24,"invalidationType":"REFETCH","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"new_input_refetch",
			InvalidationJobV4{999, "http://foo.com/bar", "user0", "ds0", 24, REFETCH, start},
			`{"startTime":"` + startStr + `","id":999,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","ttlHours":24,"invalidationType":"REFETCH"}`,
		},
		{
			"new_special_url_characters",
			InvalidationJobV4{6739471, "http://foo.com/bar(.*@#$%^*(", "user0", "ds0", 27528, REFRESH, start},
			`{"startTime":"` + startStr + `","id":6739471,"assetUrl":"http://foo.com/bar(.*@#$%^*(","createdBy":"user0","deliveryService":"ds0","ttlHours":27528,"invalidationType":"REFRESH"}`,
		},
	}
	for _, expected := range expectedNewInputs {
		t.Run(expected.TestName, func(t *testing.T) {
			obj := InvalidationJobV4{}
			if err := json.Unmarshal([]byte(expected.Input), &obj); err != nil {
				t.Errorf("expected new input '%v' to marshal, actual error: %v", expected.Input, err)
			}
			if !reflect.DeepEqual(obj, expected.Output) {
				t.Errorf("expected new input to create '%+v' actual '%+v'", expected.Output, obj)
			}
		})
	}

	type ExpectedInputErrs struct {
		TestName string
		Input    string
	}
	expectedInputErrs := []ExpectedInputErrs{
		{
			"input_err_malformed_json",
			`{"startTime"":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"input_err_time",
			`{"startTime":"` + `April 1, 2020` + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:24h"}`,
		},
		{
			"input_bad_legacy_parameters_format",
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:24s"}`,
		},
		{
			"input_bad_legacy_parameters_number_float",
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:2.4h"}`,
		},
		{
			"input_bad_legacy_parameters_number_str",
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:twentyfourh"}`,
		},
		{
			"input_bad_legacy_parameters_prefix",
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TYL:24s"}`,
		},
		{
			"input_bad_legacy_parameters_suffix",
			`{"startTime":"` + startStr + `","id":42,"assetUrl":"http://foo.com/bar","createdBy":"user0","deliveryService":"ds0","keyword":"PURGE","parameters":"TTL:24q"}`,
		},
	}
	for _, expected := range expectedInputErrs {
		t.Run(expected.TestName, func(t *testing.T) {
			obj := InvalidationJobV4{}
			if err := json.Unmarshal([]byte(expected.Input), &obj); err == nil {
				t.Errorf("expected input '%v' to have an error marshalling, actual: nil error, object %+v", expected.Input, obj)
			}
		})
	}

	// this has to be tested directly to achieve 100% "test coverage," because it will never error when called via json.Unmarshal
	t.Run("InvalidationJobV4_UnmarshalJSON_malformed_json", func(t *testing.T) {
		input := "asdf"
		obj := &InvalidationJobV4{}
		if err := obj.UnmarshalJSON([]byte(input)); err == nil {
			t.Errorf("expected input '%v' to have an error marshalling, actual: nil error, object %+v", input, obj)
		}
	})

	// create tests

	createTimeStr := start.Format(JobLegacyTimeFormat)
	type ExpectedCreateInputOutput struct {
		TestName string
		Input    InvalidationJobCreateV40PlusLegacy
		Output   InvalidationJobCreateV4
	}
	expectedCreateInputOutput := []ExpectedCreateInputOutput{
		{
			"create_legacy",
			InvalidationJobCreateV40PlusLegacy{
				InvalidationJobCreateV4{
					Regex: "/foo",
				},
				makeDSName("myds0"),
				&createTimeStr,
				makeTTLDuration(time.Hour * 7),
			},
			InvalidationJobCreateV4{
				Regex:           "/foo",
				TTLHours:        7,
				StartTime:       start,
				DeliveryService: "myds0",
			},
		},
	}
	for _, expected := range expectedCreateInputOutput {
		t.Run(expected.TestName, func(t *testing.T) {
			actual, err := InvalidationJobCreateV40LegacyToNew(expected.Input, nil)
			if err != nil {
				t.Fatalf("expected input '%+v' to not err, actual: %v", expected.Input, err)
			} else if !reflect.DeepEqual(actual, expected.Output) {
				t.Errorf("expected input '%+v' to produce output %+v, actual: %+v", expected.Input, expected.Output, actual)
			}
		})
	}

	type ExpectedCreateInputErr struct {
		TestName string
		Input    InvalidationJobCreateV40PlusLegacy
	}
	expectedCreateInputErrs := []ExpectedCreateInputErr{
		{
			"create_bad_time",
			InvalidationJobCreateV40PlusLegacy{
				InvalidationJobCreateV4{
					Regex: "/foo",
				},
				makeDSName("myds0"),
				util.StrPtr("bad time"),
				makeTTLDuration(time.Hour * 7),
			},
		},
		{
			"create_bad_ds",
			InvalidationJobCreateV40PlusLegacy{
				InvalidationJobCreateV4{
					Regex: "/foo",
				},
				makeInterfacePtr(false),
				&createTimeStr,
				makeTTLDuration(time.Hour * 7),
			},
		},
		{
			"create_bad_ttl",
			InvalidationJobCreateV40PlusLegacy{
				InvalidationJobCreateV4{
					Regex: "/foo",
				},
				makeDSName("myds0"),
				&createTimeStr,
				makeInterfacePtr(false),
			},
		},
	}
	for _, expected := range expectedCreateInputErrs {
		t.Run(expected.TestName, func(t *testing.T) {
			if _, err := InvalidationJobCreateV40LegacyToNew(expected.Input, nil); err == nil {
				t.Errorf("expected input '%+v' to produce error, actual: nil error", expected.Input)
			}
		})
	}
}

func TestLegacyTTLHours(t *testing.T) {
	inputErrs := []interface{}{
		float64(-1),
		float64(MaxTTL + 1),
		false,
		true,
		"0h",
		"999999999",
		struct{}{},
	}
	for _, input := range inputErrs {
		if _, err := legacyTTLHours(input); err == nil {
			t.Errorf("expected legacyTTLHours(%v) to error, actual: nil error", input)
		}
	}
}

func makeTTLDuration(dur time.Duration) *interface{} {
	ret := interface{}(dur.String())
	return &ret
}

func makeTTLNumber(hours float64) *interface{} {
	ret := interface{}(hours)
	return &ret
}

func makeDSName(name string) *interface{} {
	ret := interface{}(name)
	return &ret
}

func makeInterfacePtr(obj interface{}) *interface{} {
	return &obj
}
