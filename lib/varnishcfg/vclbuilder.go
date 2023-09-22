// Package varnishcfg manages generating configuration files
// for Varnish cache and Hitch proxy using data from Traffic Ops APIs.
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
	"fmt"
	"strings"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
)

// VCLBuilder builds the default VCL file using TO data.
type VCLBuilder struct {
	toData *t3cutil.ConfigData
	// opts
}

// NewVCLBuilder returns a new VCLBuilder object.
func NewVCLBuilder(toData *t3cutil.ConfigData) VCLBuilder {
	return VCLBuilder{
		toData: toData,
	}
}

// BuildVCLFile builds the default VCL file.
func (vb *VCLBuilder) BuildVCLFile() (string, []string, error) {
	warnings := make([]string, 0)
	v := newVCLFile(defaultVCLVersion)

	// access control should be added first to ensure no request processed if it is not allowed
	aclWarnings, err := vb.configureAccessControl(&v)
	warnings = append(warnings, aclWarnings...)
	if err != nil {
		return "", nil, fmt.Errorf("(warnings: %s) %w", strings.Join(warnings, ", "), err)
	}

	atsMajorVersion := uint(9)

	parents, dataWarns, err := atscfg.MakeParentDotConfigData(
		vb.toData.DeliveryServices,
		vb.toData.Server,
		vb.toData.Servers,
		vb.toData.Topologies,
		vb.toData.ServerParams,
		vb.toData.ParentConfigParams,
		vb.toData.ServerCapabilities,
		vb.toData.DSRequiredCapabilities,
		vb.toData.CacheGroups,
		vb.toData.DeliveryServiceServers,
		vb.toData.CDN,
		&atscfg.ParentConfigOpts{},
		atsMajorVersion,
	)
	warnings = append(warnings, dataWarns...)
	if err != nil {
		return "", nil, fmt.Errorf("(warnings: %s) %w", strings.Join(warnings, ", "), err)
	}
	profileDSes, dsWarnings := atscfg.GetProfileDSes(
		vb.toData.Server,
		vb.toData.Servers,
		vb.toData.DeliveryServices,
		vb.toData.DeliveryServiceServers,
	)
	warnings = append(warnings, dsWarnings...)
	cacheWarnings := vb.configureUncacheableDSes(&v, profileDSes)
	warnings = append(warnings, cacheWarnings...)

	dirWarnings, err := vb.configureDirectors(&v, parents)
	warnings = append(warnings, dirWarnings...)

	return fmt.Sprint(v), warnings, err
}
