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
	"os"
	"path/filepath"
	"testing"
	"time"
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

		WriteCacheJSON(tmpDir, fileName, largeObj)
		loadedLargeObj := Obj{}
		if err := GetJSONObjFromFile(tmpDir, fileName, time.Hour, &loadedLargeObj); err != nil {
			t.Fatalf("GetJSONObjFromFile large error expected nil, actual: " + err.Error())
		} else if largeObj.S != loadedLargeObj.S {
			t.Fatalf("GetJSONObjFromFile expected %+v actual %+v\n", largeObj.S, loadedLargeObj.S)
		}
	}

	{
		// Write a smaller object to the same file, to make sure it properly truncates, and doesn't leave old text lying around (resulting in malformed json)
		smallObj := Obj{S: `foo`}
		WriteCacheJSON(tmpDir, fileName, smallObj)
		loadedSmallObj := Obj{}
		if err := GetJSONObjFromFile(tmpDir, fileName, time.Hour, &loadedSmallObj); err != nil {
			t.Fatalf("GetJSONObjFromFile small error expected nil, actual: " + err.Error())
		} else if smallObj.S != loadedSmallObj.S {
			t.Fatalf("GetJSONObjFromFile expected %+v actual %+v\n", smallObj.S, loadedSmallObj.S)
		}
	}
}
