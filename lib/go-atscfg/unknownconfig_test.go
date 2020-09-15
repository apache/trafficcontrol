package atscfg

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
	"strings"
	"testing"
)

func TestMakeUnknownConfig(t *testing.T) {
	profileName := "myProfile"
	toolName := "myToolName"
	toURL := "https://myto.example.net"
	paramData := map[string]string{
		"param0": "val0",
		"param1": "val1",
		"param2": "val2",
	}

	txt := MakeUnknownConfig(profileName, paramData, toolName, toURL)

	testComment(t, txt, profileName, toolName, toURL)

	if !strings.Contains(txt, "val0") {
		t.Errorf("expected config to contain paramData value 'val0', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "val1") {
		t.Errorf("expected config to contain paramData value 'val1', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "val2") {
		t.Errorf("expected config to contain paramData value 'val2', actual: '%v'", txt)
	}
	if strings.Contains(txt, "param0") {
		t.Errorf("expected config to NOT contain paramData name 'param0', actual: '%v'", txt)
	}

	paramData["header"] = "none"

	txt = MakeUnknownConfig(profileName, paramData, toolName, toURL)

	firstLine := strings.TrimSpace(strings.SplitN(txt, "\n", 2)[0]) // SplitN always returns at least 1 element, no need to check len before indexing
	if strings.HasPrefix(firstLine, "#") {
		t.Errorf("expected config with 'header=none' to NOT contain header line, actual: '%v'", txt)
	}

	paramData["header"] = "foobar"

	txt = MakeUnknownConfig(profileName, paramData, toolName, toURL)

	firstLine = strings.TrimSpace(strings.SplitN(txt, "\n", 2)[0]) // SplitN always returns at least 1 element, no need to check len before indexing
	if firstLine != "foobar" {
		t.Errorf("expected config with 'header=foobar' to contain header 'foobar', actual: '%v'", txt)
	}
}
