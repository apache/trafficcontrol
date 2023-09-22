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
	"net"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func (v VCLBuilder) configureAccessControl(vclFile *vclFile) ([]string, error) {
	warnings := make([]string, 0)
	allowAllIPs := []string{
		`"127.0.0.1"`,
		`"::1"`,
	}

	purgeIPs := atscfg.GetPurgeIPs(v.toData.ServerParams)
	for _, ip := range purgeIPs {
		if parts := strings.Split(ip, `/`); len(parts) == 2 {
			allowAllIPs = append(allowAllIPs, fmt.Sprintf(`"%s"/%s`, parts[0], parts[1]))
			continue
		}
		allowAllIPs = append(allowAllIPs, fmt.Sprintf(`"%s"`, ip))
	}

	if v.toData.Server.Type == tc.CacheTypeEdge.String() {
		configureAccessControlForEdge(vclFile.acls, vclFile.subroutines, allowAllIPs)
		return warnings, nil
	}
	if v.toData.Server.Type == tc.CacheTypeMid.String() {

		coalesceMaskLenV4, coalesceNumberV4, coalesceMaskLenV6, coalesceNumberV6, ws := atscfg.GetCoalesceMaskAndNumber(v.toData.ServerParams)

		warnings = append(warnings, ws...)
		cidrs, cidr6s, ws, err := atscfg.GetAllowedCIDRsForMid(
			v.toData.Server,
			v.toData.Servers,
			v.toData.CacheGroups,
			v.toData.Topologies,
			coalesceNumberV4,
			coalesceMaskLenV4,
			coalesceNumberV6,
			coalesceMaskLenV6,
		)
		warnings = append(warnings, ws...)
		if err != nil {
			return warnings, err
		}
		allowAllButPushPurge := cidrsToVarnishCIDRs(append(cidrs, cidr6s...))
		allowAllButPushPurge = append(allowAllButPushPurge, `"10.0.0.0"/8`, `"172.16.0.0"/12`, `"192.168.0.0"/16`)

		configureAccessControlForMid(vclFile.acls, vclFile.subroutines, allowAllIPs, allowAllButPushPurge)
	}
	return warnings, nil
}

// cidrsToVarnishCIDRs converts CIDRs from the format IP/mask to "IP"/mask
func cidrsToVarnishCIDRs(cidrs []*net.IPNet) []string {
	varnishCIDRs := make([]string, 0)
	for _, cidr := range cidrs {
		cidrParts := strings.Split(cidr.String(), "/")
		varnishCIDR := fmt.Sprintf(`"%s"/%s`, cidrParts[0], cidrParts[1])
		varnishCIDRs = append(varnishCIDRs, varnishCIDR)
	}
	return varnishCIDRs
}

func configureAccessControlForEdge(acls, subroutines map[string][]string, allowAllIPs []string) {
	acls["allow_all"] = append(acls["allow_all"], allowAllIPs...)

	subroutines["vcl_recv"] = append(subroutines["vcl_recv"], []string{
		`if ((req.method == "PUSH" || req.method == "PURGE" || req.method == "DELETE") && !client.ip ~ allow_all) {`,
		`	return (synth(405));`,
		`}`,
		`if (req.method == "PURGE") {`,
		`	return (purge);`,
		`}`,
	}...)
}

func configureAccessControlForMid(acls, subroutines map[string][]string, allowAllIPs, allowAllButPushPurge []string) {
	acls["allow_all"] = append(acls["allow_all"], allowAllIPs...)
	acls["allow_all_but_push_purge"] = append(acls["allow_all_but_push_purge"], allowAllButPushPurge...)

	subroutines["vcl_recv"] = append(subroutines["vcl_recv"], []string{
		// push and purge are not allowed except for allow_all acl
		`if ((req.method == "PUSH" || req.method == "PURGE") && client.ip ~ allow_all_but_push_purge) {`,
		`	return (synth(405));`,
		`}`,
		// mid cache only accepts requests from from allowed IPs
		`if (!client.ip ~ allow_all_but_push_purge && !client.ip ~ allow_all) {`,
		`	return (synth(405));`,
		`}`,
		`if (req.method == "PURGE") {`,
		`	return (purge);`,
		`}`,
	}...)
}
