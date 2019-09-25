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

func TestMakeATSDotRules(t *testing.T) {
	profileName := "myProfile"
	toolName := "myToolName"
	toURL := "https://myto.example.net"
	paramData := map[string]string{
		"Drive_Prefix":      "/dev/sd",
		"Drive_Letters":     "a,b,c,d,e",
		"RAM_Drive_Prefix":  "/dev/ra",
		"RAM_Drive_Letters": "f,g,h",
	}

	txt := MakeATSDotRules(profileName, paramData, toolName, toURL)

	testComment(t, txt, profileName, toolName, toURL)

	if count := strings.Count(txt, "\n"); count != 9 { // one line for each drive letter, plus 1 comment
		t.Errorf("expected one line for each drive letter plus a comment, actual: '%v' count %v", txt, count)
	}

	if !strings.Contains(txt, "sda") {
		t.Errorf("expected sda for drive letter, actual: '%v'", txt)
	}
	if !strings.Contains(txt, "rah") {
		t.Errorf("expected sda for drive letter, actual: '%v'", txt)
	}
}
