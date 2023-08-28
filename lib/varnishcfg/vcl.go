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

import "fmt"

const defaultVCLVersion = "4.1"

// vclFile contains all VCL components
type vclFile struct {
	version     string
	imports     []string
	acls        map[string][]string
	backends    map[string]backend
	subroutines map[string][]string
}

func newVCLFile(version string) vclFile {
	return vclFile{
		version:     version,
		imports:     make([]string, 0),
		acls:        make(map[string][]string),
		backends:    make(map[string]backend),
		subroutines: make(map[string][]string),
	}
}

func (v vclFile) String() string {
	txt := fmt.Sprintf("vcl %s;\n", v.version)
	for _, i := range v.imports {
		txt += fmt.Sprintf("import %s;\n", i)
	}

	for name, backend := range v.backends {
		txt += fmt.Sprintf("backend %s {\n", name)
		txt += fmt.Sprint(backend)
		txt += fmt.Sprint("}\n")
	}
	// varnishd will fail if there are no backends defined
	if len(v.backends) == 0 {
		txt += fmt.Sprint("backend default none;\n")
	}

	for name, acl := range v.acls {
		txt += fmt.Sprintf("acl %s {\n", name)
		for _, entry := range acl {
			txt += fmt.Sprintf("\t%s;\n", entry)
		}
		txt += fmt.Sprint("}\n")
	}

	// has to be before other subroutines for variables initialization
	if _, ok := v.subroutines["vcl_init"]; ok {
		txt += fmt.Sprint("sub vcl_init {\n")
		for _, entry := range v.subroutines["vcl_init"] {
			txt += fmt.Sprintf("\t%s\n", entry)
		}
		txt += fmt.Sprint("}\n")
	}

	for name, subroutine := range v.subroutines {
		if name == "vcl_init" {
			continue
		}
		txt += fmt.Sprintf("sub %s {\n", name)
		for _, entry := range subroutine {
			txt += fmt.Sprintf("\t%s\n", entry)
		}
		txt += fmt.Sprint("}\n")
	}

	return txt
}

type backend struct {
	host string
	port int
}

func (b backend) String() string {
	txt := fmt.Sprintf("\t.host = \"%s\";\n", b.host)
	txt += fmt.Sprintf("\t.port = \"%d\";\n", b.port)
	return txt
}
