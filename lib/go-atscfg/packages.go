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
	"encoding/json"
	"sort"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// PackagesFileName is an unused constant of unknown purpose.
//
// Deprecated: Since the 'package' ConfigFile value is a "dummy" value used to
// indicate packages that should be installed through yum, there isn't a need
// for any definition of a file name, and the Parameter ConfigFile value is
// already exported as PackagesParamConfigFile.
const PackagesFileName = `packages`

// PackagesParamConfigFile is the ConfigFile value of Parameters that define
// system packages to be installed on a cache server.
const PackagesParamConfigFile = `package`

// ContentTypePackages is a MIME type of unknown meaning and purpose.
//
// Deprecated: Since the 'package' ConfigFile value is a "dummy" value used to
// indicate packages that should be installed through yum, there isn't a need
// for any definition of a content type for a file that doesn't exist, and this
// value is never used for anything anyway. The contents of the file as output
// by tc3-generate are actually encoded as JSON, so at best this is inaccurate.
const ContentTypePackages = ContentTypeTextASCII

// LineCommentPackages is used only to convey the idea that since "package
// Parameters" don't define file contents they don't have comments and therefore
// no string signifies the beginning of a comment for this non-existent grammar.
//
// Deprecated: This constant expresses a concept that by its own definition has
// no meaning.
const LineCommentPackages = ""

// PackagesOpts contains settings to configure generation options.
type PackagesOpts struct {
}

// MakePackages returns the 'packages' ATS config file endpoint.
// This is a JSON object, and should be served with an 'application/json' Content-Type.
func MakePackages(
	serverParams []tc.ParameterV5,
	opts *PackagesOpts,
) (Cfg, error) {
	if opts == nil {
		opts = &PackagesOpts{}
	}
	warnings := []string{}

	params := paramsToMultiMap(filterParams(serverParams, PackagesParamConfigFile, "", "", ""))

	pkgs := []pkg{}
	for name, versions := range params {
		for _, version := range versions {
			pkgs = append(pkgs, pkg{Name: name, Version: version})
		}
	}
	sort.Sort(packages(pkgs))
	bts, err := json.Marshal(&pkgs)
	if err != nil {
		// should never happen
		return Cfg{}, makeErr(warnings, "marshalling chkconfig NameVersions: "+err.Error())
	}

	return Cfg{
		Text:        string(bts),
		ContentType: ContentTypePackages,
		LineComment: LineCommentPackages,
		Warnings:    warnings,
	}, nil
}

type pkg struct {
	Name    string
	Version string
}

type packages []pkg

func (ps packages) Len() int { return len(ps) }
func (ps packages) Less(i, j int) bool {
	if ps[i].Name != ps[j].Name {
		return ps[i].Name < ps[j].Name
	}
	return ps[i].Version < ps[j].Version
}
func (ps packages) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }
