package varnishcfg

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
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestConfigureCustomVCL(t *testing.T) {
	vb := NewVCLBuilder(&t3cutil.ConfigData{
		ServerParams: []tc.ParameterV50{
			{ConfigFile: "default.vcl", Name: "import", Value: "std"},
			{ConfigFile: "default.vcl", Name: "vcl_recv", Value: "set req.url = std.querysort(req.url);"},
			{
				ConfigFile: "default.vcl",
				Name:       "vcl_deliver",
				Value:      "if (req.status >= 400 && req.status <= 500) {\n\tset req.status = 404;\n}",
			},
		},
	})

	vclFile := newVCLFile(defaultVCLVersion)
	vb.configureCustomVCL(&vclFile)

	expectedVCLFile := newVCLFile(defaultVCLVersion)
	expectedVCLFile.imports = append(expectedVCLFile.imports, "std")
	expectedVCLFile.subroutines["vcl_recv"] = []string{
		"set req.url = std.querysort(req.url);",
	}
	expectedVCLFile.subroutines["vcl_deliver"] = []string{
		"if (req.status >= 400 && req.status <= 500) {",
		"	set req.status = 404;",
		"}",
	}

	if !reflect.DeepEqual(vclFile, expectedVCLFile) {
		t.Errorf("got %v want %v", vclFile, expectedVCLFile)
	}
}
