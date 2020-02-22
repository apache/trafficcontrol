package toreq

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
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestWriteCacheJSON(t *testing.T) {
	type Obj struct {
		S string `json:"s"`
	}

	fileName := "TestWriteCacheJSON.json"
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, fileName)
	defer os.Remove(filePath)

	{
		largeObj := Obj{
			S: `
    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

    Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst.
`,
		}

		WriteCache(CacheFormatJSON.Encoder, tmpDir, fileName, largeObj)
		loadedLargeObj := Obj{}
		if err := ReadCache(CacheFormatJSON.Decoder, tmpDir, fileName, time.Hour, &loadedLargeObj); err != nil {
			t.Fatalf("ReadCache large error expected nil, actual: " + err.Error())
		} else if largeObj.S != loadedLargeObj.S {
			t.Fatalf("ReadCache expected %+v actual %+v\n", largeObj.S, loadedLargeObj.S)
		}
	}

	{
		// Write a smaller object to the same file, to make sure it properly truncates, and doesn't leave old text lying around (resulting in malformed json)
		smallObj := Obj{S: `foo`}
		WriteCache(CacheFormatJSON.Encoder, tmpDir, fileName, smallObj)
		loadedSmallObj := Obj{}
		if err := ReadCache(CacheFormatJSON.Decoder, tmpDir, fileName, time.Hour, &loadedSmallObj); err != nil {
			t.Fatalf("GetJSONObjFromFile small error expected nil, actual: " + err.Error())
		} else if smallObj.S != loadedSmallObj.S {
			t.Fatalf("GetJSONObjFromFile expected %+v actual %+v\n", smallObj.S, loadedSmallObj.S)
		}
	}
}

func TestWriteCacheCBOR(t *testing.T) {
	type Obj struct {
		S string `json:"s"`
	}

	fileName := "TestWriteCacheCBOR"
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, fileName)
	defer os.Remove(filePath)

	{
		largeObj := Obj{
			S: `
    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

    Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst.
`,
		}

		WriteCache(CacheFormatCBOR.Encoder, tmpDir, fileName, largeObj)
		loadedLargeObj := Obj{}
		if err := ReadCache(CacheFormatCBOR.Decoder, tmpDir, fileName, time.Hour, &loadedLargeObj); err != nil {
			t.Fatalf("ReadCache large error expected nil, actual: " + err.Error())
		} else if largeObj.S != loadedLargeObj.S {
			t.Fatalf("ReadCache expected %+v actual %+v\n", largeObj.S, loadedLargeObj.S)
		}
	}

	{
		// Write a smaller object to the same file, to make sure it properly truncates, and doesn't leave old text lying around (resulting in malformed cbor)
		smallObj := Obj{S: `foo`}
		WriteCache(CacheFormatCBOR.Encoder, tmpDir, fileName, smallObj)
		loadedSmallObj := Obj{}
		if err := ReadCache(CacheFormatCBOR.Decoder, tmpDir, fileName, time.Hour, &loadedSmallObj); err != nil {
			t.Fatalf("ReadCache small error expected nil, actual: " + err.Error())
		} else if smallObj.S != loadedSmallObj.S {
			t.Fatalf("ReadCache expected %+v actual %+v\n", smallObj.S, loadedSmallObj.S)
		}
	}
}

func TestDefaultCacheFormatIsomorphic(t *testing.T) {
	// Test whether DefaultCacheFormat is isomorphic.
	// That is, whether serializing and deserializing produces the original object.
	// This might seem silly, but encoding/gob is not isomorphic and fails this test.
	// We requires an isomorphic cache format, or config files will be wrong.

	// Delivery Service is one of TC's most complex objects, so use an array of it to test.
	// It's important to test null pointers, as well as pointers to default values (a common failure).

	dsMatchPtr := []tc.DeliveryServiceMatch{
		tc.DeliveryServiceMatch{
			Type:      tc.DSMatchTypeHostRegex,
			SetNumber: 0,
			Pattern:   "foo",
		},
		tc.DeliveryServiceMatch{
			Type:      tc.DSMatchTypeInvalid,
			SetNumber: 42,
			Pattern:   "",
		},
	}
	dsTypeInvalidPtr := tc.DSTypeInvalid

	ds1 := tc.DeliveryServiceNullable{}
	ds1.Active = util.BoolPtr(false)
	ds1.AnonymousBlockingEnabled = nil
	ds1.CacheURL = util.StrPtr("")
	ds1.CCRDNSTTL = util.IntPtr(0)
	ds1.CDNID = nil
	ds1.CDNName = nil
	ds1.CheckPath = util.StrPtr("foo")
	ds1.DisplayName = util.StrPtr("")
	ds1.DSCP = nil
	ds1.LogsEnabled = util.BoolPtr(true)
	// ds1.MatchList = &dsMatchPtr
	ds1.MissLat = nil
	ds1.MissLong = util.FloatPtr(0)
	ds1.Signed = false
	ds1.Type = &dsTypeInvalidPtr
	// ds1.ExampleURLs = []string{"foo", ""}

	ds2 := tc.DeliveryServiceNullable{}
	ds2.Active = nil
	ds2.AnonymousBlockingEnabled = util.BoolPtr(false)
	ds2.CacheURL = util.StrPtr("")
	ds2.CCRDNSTTL = util.IntPtr(0)
	ds2.CDNID = nil
	ds2.CDNName = nil
	ds2.CheckPath = util.StrPtr("foo")
	ds2.DisplayName = util.StrPtr("")
	ds2.DSCP = nil
	ds2.LogsEnabled = util.BoolPtr(true)
	ds2.MatchList = &dsMatchPtr
	ds2.MissLat = util.FloatPtr(0)
	ds2.MissLong = util.FloatPtr(42)
	ds2.Signed = true
	ds2.Type = nil
	ds2.ExampleURLs = nil

	dses := []tc.DeliveryServiceNullable{ds1, ds2}

	fileName := "TestDefaultCacheFormatIsomorphic"

	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, fileName)
	defer os.Remove(filePath)

	WriteCache(DefaultCacheFormat.Encoder, tmpDir, fileName, dses)

	readDSes := []tc.DeliveryServiceNullable{}
	if err := ReadCache(DefaultCacheFormat.Decoder, tmpDir, fileName, time.Hour, &readDSes); err != nil {
		t.Fatalf("ReadCache error expected nil, actual: " + err.Error())
	}
	if len(readDSes) != 2 {
		t.Fatalf("ReadCache error expected 2 dses, actual: %v", len(readDSes))
	}

	if !reflect.DeepEqual(dses[0], readDSes[0]) {
		dsj, _ := json.MarshalIndent(dses[0], "", " ")
		dsrj, _ := json.MarshalIndent(readDSes[0], "", " ")
		t.Errorf("ReadCache expected %+v actual %+v\n", string(dsj), string(dsrj))
	}

	if !reflect.DeepEqual(dses[1], readDSes[1]) {
		dsj, _ := json.MarshalIndent(dses[1], "", " ")
		dsrj, _ := json.MarshalIndent(readDSes[1], "", " ")
		t.Errorf("ReadCache expected %+v actual %+v\n", string(dsj), string(dsrj))
	}
}
