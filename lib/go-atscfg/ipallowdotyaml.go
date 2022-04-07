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
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

const IPAllowYamlFileName = `ip_allow.yaml`
const ContentTypeIPAllowDotYAML = ContentTypeYAML
const LineCommentIPAllowDotYAML = LineCommentHash

// const ParamPurgeAllowIP = "purge_allow_ip"
// const ParamCoalesceMaskLenV4 = "coalesce_masklen_v4"
// const ParamCoalesceNumberV4 = "coalesce_number_v4"
// const ParamCoalesceMaskLenV6 = "coalesce_masklen_v6"
// const ParamCoalesceNumberV6 = "coalesce_number_v6"

// const DefaultCoalesceMaskLenV4 = 24
// const DefaultCoalesceNumberV4 = 5
// const DefaultCoalesceMaskLenV6 = 48
// const DefaultCoalesceNumberV6 = 5

// AStatsDotConfigOpts contains settings to configure generation options.
type IPAllowDotYAMLOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeIPAllowDotYAML creates the ip_allow.yaml ATS 9+ config file.
func MakeIPAllowDotYAML(
	serverParams []tc.Parameter,
	server *Server,
	servers []Server,
	cacheGroups []tc.CacheGroupNullable,
	topologies []tc.Topology,
	opt *IPAllowDotYAMLOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &IPAllowDotYAMLOpts{}
	}
	warnings := []string{}

	if server.Cachegroup == nil {
		return Cfg{}, makeErr(warnings, "this server missing Cachegroup")
	}
	if server.HostName == nil {
		return Cfg{}, makeErr(warnings, "this server missing HostName")
	}

	params := paramsToMultiMap(filterParams(serverParams, IPAllowConfigFileName, "", "", ""))

	ipAllowDat := []ipAllowYAMLData{}

	// localhost is trusted.
	ipAllowDat = append([]ipAllowYAMLData{yamlAllowAll(`127.0.0.1`)}, ipAllowDat...)
	ipAllowDat = append([]ipAllowYAMLData{yamlAllowAll(`::1`)}, ipAllowDat...)

	// default for coalesce_ipv4 = 24, 5 and for ipv6 48, 5; override with the parameters in the server profile.
	coalesceMaskLenV4 := DefaultCoalesceMaskLenV4
	coalesceNumberV4 := DefaultCoalesceNumberV4
	coalesceMaskLenV6 := DefaultCoalesceMaskLenV6
	coalesceNumberV6 := DefaultCoalesceNumberV6

	for name, vals := range params {
		for _, val := range vals {
			switch name {
			case ParamPurgeAllowIP:
				ipAllowDat = append(ipAllowDat, yamlAllowAll(val))
			case ParamCoalesceMaskLenV4:
				if vi, err := strconv.Atoi(val); err != nil {
					warnings = append(warnings, "got param '"+name+"' val '"+val+"' not a number, ignoring!")
				} else if coalesceMaskLenV4 != DefaultCoalesceMaskLenV4 {
					warnings = append(warnings, "got multiple param '"+name+"' - ignoring  val '"+val+"'!")
				} else {
					coalesceMaskLenV4 = vi
				}
			case ParamCoalesceNumberV4:
				if vi, err := strconv.Atoi(val); err != nil {
					warnings = append(warnings, "got param '"+name+"' val '"+val+"' not a number, ignoring!")
				} else if coalesceNumberV4 != DefaultCoalesceNumberV4 {
					warnings = append(warnings, "got multiple param '"+name+"' - ignoring  val '"+val+"'!")
				} else {
					coalesceNumberV4 = vi
				}
			case ParamCoalesceMaskLenV6:
				if vi, err := strconv.Atoi(val); err != nil {
					warnings = append(warnings, "got param '"+name+"' val '"+val+"' not a number, ignoring!")
				} else if coalesceMaskLenV6 != DefaultCoalesceMaskLenV6 {
					warnings = append(warnings, "got multiple param '"+name+"' - ignoring  val '"+val+"'!")
				} else {
					coalesceMaskLenV6 = vi
				}
			case ParamCoalesceNumberV6:
				if vi, err := strconv.Atoi(val); err != nil {
					warnings = append(warnings, "got param '"+name+"' val '"+val+"' not a number, ignoring!")
				} else if coalesceNumberV6 != DefaultCoalesceNumberV6 {
					warnings = append(warnings, "got multiple param '"+name+"' - ignoring  val '"+val+"'!")
				} else {
					coalesceNumberV6 = vi
				}
			}
		}
	}

	// for edges deny "PUSH|PURGE|DELETE", allow everything else to everyone.
	isMid := strings.HasPrefix(server.Type, tc.MidTypePrefix)
	if !isMid {
		ipAllowDat = append(ipAllowDat, yamlAllowAllButPushPurgeDelete(`0.0.0.0/0`))
		ipAllowDat = append(ipAllowDat, yamlAllowAllButPushPurgeDelete(`::/0`))
	} else {

		ips := []*net.IPNet{}
		ip6s := []*net.IPNet{}

		cgMap := map[string]tc.CacheGroupNullable{}
		for _, cg := range cacheGroups {
			if cg.Name == nil {
				return Cfg{}, makeErr(warnings, "got cachegroup with nil name!")
			}
			cgMap[*cg.Name] = cg
		}

		if server.Cachegroup == nil {
			return Cfg{}, makeErr(warnings, "server had nil Cachegroup!")
		}

		serverCG, ok := cgMap[*server.Cachegroup]
		if !ok {
			return Cfg{}, makeErr(warnings, "server cachegroup not in cachegroups!")
		}

		childCGNames := getTopologyDirectChildren(tc.CacheGroupName(*server.Cachegroup), topologies)

		childCGs := map[string]tc.CacheGroupNullable{}
		for cgName, _ := range childCGNames {
			childCGs[string(cgName)] = cgMap[string(cgName)]
		}

		for cgName, cg := range cgMap {
			if (cg.ParentName != nil && *cg.ParentName == *serverCG.Name) || (cg.SecondaryParentName != nil && *cg.SecondaryParentName == *serverCG.Name) {
				childCGs[cgName] = cg
			}
		}

		// sort servers, to guarantee things like IP coalescing are deterministic
		sort.Sort(serversSortByName(servers))
		for _, childServer := range servers {
			if childServer.Cachegroup == nil {
				warnings = append(warnings, "Servers had server with nil Cachegroup, skipping!")
				continue
			} else if childServer.HostName == nil {
				warnings = append(warnings, "Servers had server with nil HostName, skipping!")
				continue
			}

			// We need to add IPs to the allow of
			// - all children of this server
			// - all monitors, if this server is a Mid
			//
			_, isChild := childCGs[*childServer.Cachegroup]
			if !isChild && !strings.HasPrefix(server.Type, tc.MidTypePrefix) && string(childServer.Type) != tc.MonitorTypeName {
				continue
			}

			for _, svInterface := range childServer.Interfaces {
				for _, svAddr := range svInterface.IPAddresses {
					if ip := net.ParseIP(svAddr.Address); ip != nil {
						// got an IP - convert it to a CIDR and add it to the list
						if ip4 := ip.To4(); ip4 != nil {
							ips = append(ips, util.IPToCIDR(ip4))
						} else {
							ip6s = append(ip6s, util.IPToCIDR(ip))
						}
					} else {
						// not an IP, try a CIDR
						if ip, cidr, err := net.ParseCIDR(svAddr.Address); err != nil {
							// not a CIDR or IP - error out
							warnings = append(warnings, "server '"+*server.HostName+"' IP '"+svAddr.Address+" is not an IP address or CIDR - skipping!")
						} else if ip == nil {
							// not a CIDR or IP - error out
							warnings = append(warnings, "server '"+*server.HostName+"' IP '"+svAddr.Address+" failed to parse as IP or CIDR - skipping!")
						} else {
							// got a valid CIDR - add it to the list
							if ip4 := ip.To4(); ip4 != nil {
								ips = append(ips, cidr)
							} else {
								ip6s = append(ip6s, cidr)
							}
						}
					}
				}
			}
		}

		cidrs := util.CoalesceCIDRs(ips, coalesceNumberV4, coalesceMaskLenV4)
		cidr6s := util.CoalesceCIDRs(ip6s, coalesceNumberV6, coalesceMaskLenV6)

		for _, cidr := range cidrs {
			ipAllowDat = append(ipAllowDat, yamlAllowAllButPushPurge(cidr.String()))
		}
		for _, cidr := range cidr6s {
			ipAllowDat = append(ipAllowDat, yamlAllowAllButPushPurge(cidr.String()))
		}

		// allow RFC 1918 server space - TODO JvD: parameterize
		ipAllowDat = append(ipAllowDat, yamlAllowAllButPushPurge(`10.0.0.0/8`))
		ipAllowDat = append(ipAllowDat, yamlAllowAllButPushPurge(`172.16.0.0/12`))
		ipAllowDat = append(ipAllowDat, yamlAllowAllButPushPurge(`192.168.0.0/16`))

		// order matters, so sort before adding the denys
		sort.Sort(ipAllowYAMLDatas(ipAllowDat))

		// start with a deny for PUSH and PURGE - TODO CDL: parameterize
		// but leave purge open through localhost
		// Edges already deny PUSH and PURGE

		// start by allowing everything to localhost, including PURGE and PUSH
		ipAllowDat = append([]ipAllowYAMLData{yamlAllowAll(`127.0.0.1`)}, ipAllowDat...)
		ipAllowDat = append([]ipAllowYAMLData{yamlAllowAll(`::1`)}, ipAllowDat...)

		// end with a deny
		ipAllowDat = append(ipAllowDat, yamlDenyAll(`0.0.0.0/0`))
		ipAllowDat = append(ipAllowDat, yamlDenyAll(`::/0`))
	}

	text := makeHdrComment(opt.HdrComment)
	text += `
ip_allow:`
	for _, al := range ipAllowDat {
		text += `
  - apply: in
    ip_addrs: ` + al.Src + `
    action: ` + al.Action + `
    methods:`
		for _, method := range al.Methods {
			text += `
      - ` + method
		}
	}
	text += "\n"

	return Cfg{
		Text:        text,
		ContentType: ContentTypeHostingDotConfig,
		LineComment: LineCommentHostingDotConfig,
		Warnings:    warnings,
	}, nil
}

type ipAllowYAMLData struct {
	Src     string
	Action  string
	Methods []string
}

type ipAllowYAMLDatas []ipAllowYAMLData

func (is ipAllowYAMLDatas) Len() int      { return len(is) }
func (is ipAllowYAMLDatas) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is ipAllowYAMLDatas) Less(i, j int) bool {
	if is[i].Src != is[j].Src {
		return is[i].Src < is[j].Src
	}
	if is[i].Action != is[j].Action {
		return is[i].Action < is[j].Action
	}
	if len(is[i].Methods) < len(is[j].Methods) {
		return true
	}
	for mi := 0; mi < len(is[i].Methods); mi++ {
		if is[i].Methods[mi] != is[j].Methods[mi] {
			return is[i].Methods[mi] < is[j].Methods[mi]
		}
	}
	return false
}

const YAMLActionAllow = "allow"
const YAMLActionDeny = "deny"
const YAMLMethodAll = "ALL"

// yamlAllowAllButPushPurge is a helper func to build a ipAllowYAMLData for the given range string immediately allowing all Methods except Push and Purge.
func yamlAllowAllButPushPurge(rangeStr string) ipAllowYAMLData {
	// Note denying methods implicitly and immediately allows all other methods!
	// So Deny PUSH|PURGE will make all other methods
	// immediately allowed, regardless of any later deny rules!
	methodPushPurge := []string{MethodPush, MethodPurge}
	return ipAllowYAMLData{
		Src:     rangeStr,
		Action:  YAMLActionDeny,
		Methods: methodPushPurge,
	}
}

// yamlAllowAllButPushPurgeDelete is a helper func to build a ipAllowYAMLData for the given range string immediately allowing all Methods except PUSH, PURGE, and DELETE.
func yamlAllowAllButPushPurgeDelete(rangeStr string) ipAllowYAMLData {
	// Note denying methods implicitly and immediately allows all other methods!
	// So Deny PUSH|PURGE will make all other methods
	// immediately allowed, regardless of any later deny rules!
	methodPushPurgeDelete := []string{MethodPush, MethodPurge, MethodDelete}
	return ipAllowYAMLData{
		Src:     rangeStr,
		Action:  YAMLActionDeny,
		Methods: methodPushPurgeDelete,
	}
}

// yamlAllowAll is a helper func to build a ipAllowYAMLData for the given range string immediately allowing all Methods, including Push and Purge.
func yamlAllowAll(rangeStr string) ipAllowYAMLData {
	return ipAllowYAMLData{
		Src:     rangeStr,
		Action:  YAMLActionAllow,
		Methods: []string{YAMLMethodAll},
	}
}

// yamlDenyAll is a helper func to build a ipAllowYAMLData for the given range string immediately denying all Methods.
func yamlDenyAll(rangeStr string) ipAllowYAMLData {
	return ipAllowYAMLData{
		Src:     rangeStr,
		Action:  YAMLActionDeny,
		Methods: []string{YAMLMethodAll},
	}
}
