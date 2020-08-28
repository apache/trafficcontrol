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

func TestMakeRecordsDotConfig(t *testing.T) {
	profileName := "myProfile"
	toolName := "myToolName"
	toURL := "https://myto.example.net"
	paramData := map[string]string{
		"param0":                    "val0",
		"param1":                    "val1",
		"param2":                    "val2",
		"test-hostname-replacement": "fooSTRING __HOSTNAME__",
	}

	txt := MakeRecordsDotConfig(profileName, paramData, toolName, toURL)

	testComment(t, txt, profileName, toolName, toURL)

	if !strings.Contains(txt, "param0 val0") {
		t.Errorf("expected config to contain paramData 'param0 val0', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "param1 val1") {
		t.Errorf("expected config to contain paramData 'param1 val1', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "param2 val2") {
		t.Errorf("expected config to contain paramData 'param2 val2', actual: '%v'", txt)
	}
	if !strings.Contains(txt, "test-hostname-replacement fooSTRING __FULL_HOSTNAME__") {
		t.Errorf("expected config to replace 'STRING __HOSTNAME__' with 'STRING __FULL_HOSTNAME__', actual: '%v'", txt)
	}
}

func TestReplaceLineSuffixes(t *testing.T) {
	{
		input := `
foo STRING __HOSTNAME__
bar
baz
`
		expected := `
foo STRING __FULL_HOSTNAME__
bar
baz
`
		actual := replaceLineSuffixes(input, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
		if expected != actual {
			t.Errorf("Expected '%v' Actual '%v'", expected, actual)
		}
	}
	{
		input := `STRING __HOSTNAME__`
		expected := `STRING __FULL_HOSTNAME__`
		actual := replaceLineSuffixes(input, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
		if expected != actual {
			t.Errorf("Expected '%v' Actual '%v'", expected, actual)
		}
	}
	{
		input := `
STRING __HOSTNAME__
`
		expected := `
STRING __FULL_HOSTNAME__
`
		actual := replaceLineSuffixes(input, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
		if expected != actual {
			t.Errorf("Expected '%v' Actual '%v'", expected, actual)
		}
	}
	{
		input := `
  
STRING __HOSTNAME__
`
		expected := `
  
STRING __FULL_HOSTNAME__
`
		actual := replaceLineSuffixes(input, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
		if expected != actual {
			t.Errorf("Expected '%v' Actual '%v'", expected, actual)
		}
	}
	{
		input := `
STRING __HOSTNAME__
  STRING __HOSTNAME__
`
		expected := `
STRING __FULL_HOSTNAME__
  STRING __FULL_HOSTNAME__
`
		actual := replaceLineSuffixes(input, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
		if expected != actual {
			t.Errorf("Expected '%v' Actual '%v'", expected, actual)
		}
	}
	{
		input := `
`
		expected := `
`
		actual := replaceLineSuffixes(input, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
		if expected != actual {
			t.Errorf("Expected '%v' Actual '%v'", expected, actual)
		}
	}
	{
		input := ``
		expected := ``
		actual := replaceLineSuffixes(input, "STRING __HOSTNAME__", "STRING __FULL_HOSTNAME__")
		if expected != actual {
			t.Errorf("Expected '%v' Actual '%v'", expected, actual)
		}
	}
}
