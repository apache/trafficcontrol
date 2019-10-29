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

func TestMakeSetDSCPDotConfig(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"
	dscpNumStr := "42"

	txt := MakeSetDSCPDotConfig(cdnName, toToolName, toURL, dscpNumStr)

	if !strings.Contains(txt, string(cdnName)) {
		t.Errorf("expected: cdnName '" + string(cdnName) + "', actual: missing")
	}
	if !strings.Contains(txt, toToolName) {
		t.Errorf("expected: toToolName '" + toToolName + "', actual: missing")
	}
	if !strings.Contains(txt, toURL) {
		t.Errorf("expected: toURL '" + toURL + "', actual: missing")
	}
	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: missing")
	}

	if !strings.Contains(txt, dscpNumStr) {
		t.Errorf("expected: dscp number '"+dscpNumStr+"' in config, actual '%v'", txt)
	}
}

func TestMakeSetDSCPDotConfigNonNumber(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"
	dscpNumStr := "42a"

	txt := MakeSetDSCPDotConfig(cdnName, toToolName, toURL, dscpNumStr)

	if !strings.Contains(strings.ToLower(txt), "error") {
		t.Errorf("expected: error from non-number dscp, actual '%v'", txt)
	}
}
