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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestMakeURLSigConfig(t *testing.T) {
	profileName := "myProfile"
	toolName := "myToolName"
	toURL := "https://myto.example.net"
	paramData := map[string]string{
		"key1": "foo",
		"key2": "bar",
		//todo
	}
	urlSigKeys := tc.URLSigKeys{}

	txt := MakeURLSigConfig(profileName, urlSigKeys, paramData, toolName, toURL)

	txt = strings.Replace(txt, " ", "", -1)

	if !strings.Contains(txt, "key1=foo") {
		t.Errorf("expected param key key1=foo, actual '%v'", txt)
	}
	if !strings.Contains(txt, "key2=bar") {
		t.Errorf("expected param key key1=foo, actual '%v'", txt)
	}

	urlSigKeys["urlsigkeys-1"] = "urlsigkeys-val-1"
	urlSigKeys["urlsigkeys-2"] = "urlsigkeys-val-2"

	txt = MakeURLSigConfig(profileName, urlSigKeys, paramData, toolName, toURL)

	txt = strings.Replace(txt, " ", "", -1)

	if !strings.Contains(txt, "urlsigkeys-1=urlsigkeys-val-1") {
		t.Errorf("expected param key key1=foo, actual '%v'", txt)
	}
	if !strings.Contains(txt, "urlsigkeys-2=urlsigkeys-val-2") {
		t.Errorf("expected param key key1=foo, actual '%v'", txt)
	}
	if strings.Contains(txt, "key1=foo") {
		t.Errorf("expected config to NOT contain param data keys 'key1=foo' if urlsig keys exist, actual '%v'", txt)
	}
	if strings.Contains(txt, "key2=bar") {
		t.Errorf("expected config to NOT contain param data keys 'key2=bar' if urlsig keys exist, actual '%v'", txt)
	}
}
