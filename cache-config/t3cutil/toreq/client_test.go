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
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

func TestRequestInfoStr(t *testing.T) {

	type Expected struct {
		Desc     string
		Input    toclientlib.ReqInf
		Path     string
		Contains []string
	}
	expecteds := []Expected{
		{"all_nils", toclientlib.ReqInf{RemoteAddr: nil, StatusCode: 0, RespHeaders: nil}, "", []string{"code=0"}},
		{"nil_addr", toclientlib.ReqInf{RemoteAddr: nil, StatusCode: 200, RespHeaders: http.Header{}}, "", []string{"code=200"}},
		{"zero_code", toclientlib.ReqInf{RemoteAddr: makeIP("192.0.2.1"), StatusCode: 0, RespHeaders: http.Header{}}, "", []string{"code=0"}},
		{"nil_header", toclientlib.ReqInf{RemoteAddr: makeIP("192.0.2.1"), StatusCode: 200, RespHeaders: nil}, "", []string{"code=200"}},
		{"ip", toclientlib.ReqInf{RemoteAddr: makeIP("192.0.2.1"), StatusCode: 200, RespHeaders: http.Header{}}, "", []string{"ip=192.0.2.1"}},
		{"date", toclientlib.ReqInf{RemoteAddr: makeIP("192.0.2.1"), StatusCode: 200, RespHeaders: http.Header{"Date": {"Mon, 29 Nov 2021 10:11:12 GMT"}}}, "", []string{`date="Mon, 29 Nov 2021 10:11:12 GMT"`}},
		{"age", toclientlib.ReqInf{RemoteAddr: makeIP("192.0.2.1"), StatusCode: 200, RespHeaders: http.Header{"Age": {"4242"}}}, "", []string{`age=42`}},
		{"path", toclientlib.ReqInf{RemoteAddr: makeIP("192.0.2.1"), StatusCode: 200, RespHeaders: http.Header{"Age": {"4242"}}}, "my-endpoint", []string{`path=my-endpoint`}},
		{"path_age_code_ip_date", toclientlib.ReqInf{RemoteAddr: makeIP("192.0.2.9"), StatusCode: 206, RespHeaders: http.Header{"Age": {"99"}, "Date": {`some malformed date`}}}, "my-everything-endpoint", []string{`path=my-everything-endpoint`, `age=99`, `code=206`, `ip=192.0.2.9`, `date="some malformed date"`}},
	}

	for _, expected := range expecteds {
		t.Run(expected.Desc, func(t *testing.T) {
			output := RequestInfoStr(expected.Input, expected.Path)
			for _, expContain := range expected.Contains {
				if !strings.Contains(output, expContain) {
					t.Errorf("expected input '%v' path '%v' to contain '%v' actual '%v'", expected.Input, expected.Path, expContain, output)
				}
			}
		})
	}
}

func makeIP(ip string) *net.IPAddr { return &net.IPAddr{IP: net.ParseIP(ip), Zone: ""} }
