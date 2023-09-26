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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

// IPAllowConfigFileName is both the file name of ip_allow.config files and the
// ConfigFile value of Parameters that can affect its contents.
const IPAllowConfigFileName = `ip_allow.config`

// ContentTypeIPAllowDotConfig is the MIME type of the contents of an
// ip_allow.config file.
const ContentTypeIPAllowDotConfig = ContentTypeTextASCII

// LineCommentIPAllowDotConfig is the string that signifies the start of a line
// comment in the grammar of an ip_allow.config file.
const LineCommentIPAllowDotConfig = LineCommentHash

// ParamPurgeAllowIP is the Name of Parameters which can specify an IP address
// to be allowed ALL methods on the cache server (not just PURGE).
//
// TODO: replace instances of "purge_allow_ip" with this currently-unused
// constant.
// TODO: unexport? This is only used in this package.
const ParamPurgeAllowIP = "purge_allow_ip"

// ParamCoalesceMaskLenV4 is the Name of a Parameter which can be used to set
// the mask length for IPv4 network specifications used to coalesce groups of
// addresses in an ip_allow.config file.
//
// TODO: unexport? This is only used in this package.
const ParamCoalesceMaskLenV4 = "coalesce_masklen_v4"

// ParamCoalesceNumberV4 is the Name of a Parameter that can be used to set the
// minimum number of IPv4 addresses in an ip_allow.config file that must be
// covered by a mask in order to allow coalescing them.
//
// TODO: unexport? This is only used in this package.
const ParamCoalesceNumberV4 = "coalesce_number_v4"

// ParamCoalesceMaskLenV6 is the Name of a Parameter which can be used to set
// the mask length for IPv6 network specifications used to coalesce groups of
// addresses in an ip_allow.config file.
//
// TODO: unexport? This is only used in this package.
const ParamCoalesceMaskLenV6 = "coalesce_masklen_v6"

// ParamCoalesceNumberV6 is the Name of a Parameter that can be used to set the
// minimum number of IPv6 addresses in an ip_allow.config file that must be
// covered by a mask in order to allow coalescing them.
//
// TODO: unexport? This is only used in this package.
const ParamCoalesceNumberV6 = "coalesce_number_v6"

// DefaultCoalesceMaskLenV4 is the value that will be used for the mask length
// for IPv4 network specifications used to coalesce groups of addresses in an
// ip_allow.config file, if a value is not given by a Parameter.
//
// TODO: unexport? This is only used in this package.
const DefaultCoalesceMaskLenV4 = 24

// DefaultCoalesceNumberV4 is the value that will be used for the minimum number
// of IPv4 addresses in an ip_allow.config file that must be covered by a mask
// in order to allow coalescing them, if a value is not given by a Parameter.
//
// TODO: unexport? This is only used in this package.
const DefaultCoalesceNumberV4 = 5

// DefaultCoalesceMaskLenV6 is the value that will be used for the mask length
// for IPv6 network specifications used to coalesce groups of addresses in an
// ip_allow.config file, if a value is not given by a Parameter.
//
// TODO: unexport? This is only used in this package.
const DefaultCoalesceMaskLenV6 = 48

// DefaultCoalesceNumberV6 is the value that will be used for the minimum number
// of IPv6 addresses in an ip_allow.config file that must be covered by a mask
// in order to allow coalescing them, if a value is not given by a Parameter.
//
// TODO: unexport? This is only used in this package.
const DefaultCoalesceNumberV6 = 5

// IPAllowDotConfigOpts contains settings to configure generation options.
type IPAllowDotConfigOpts struct {
	// HdrComment is the header comment to include at the beginning of the file.
	// This should be the text desired, without comment syntax (like # or //). The file's comment syntax will be added.
	// To omit the header comment, pass the empty string.
	HdrComment string
}

// MakeIPAllowDotConfig creates the ip_allow.config ATS config file.
// The childServers is a list of servers which are children for this Mid-tier server. This should be empty for Edge servers.
// More specifically, it should be the list of edges whose cachegroup's parent_cachegroup or secondary_parent_cachegroup is the cachegroup of this Mid server.
func MakeIPAllowDotConfig(
	serverParams []tc.ParameterV5,
	server *Server,
	servers []Server,
	cacheGroups []tc.CacheGroupNullableV5,
	topologies []tc.TopologyV5,
	opt *IPAllowDotConfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &IPAllowDotConfigOpts{}
	}
	warnings := []string{}

	if server.CacheGroup == "" {
		return Cfg{}, makeErr(warnings, "this server missing Cachegroup")
	}
	if server.HostName == "" {
		return Cfg{}, makeErr(warnings, "this server missing HostName")
	}

	params := paramsToMultiMap(filterParams(serverParams, IPAllowConfigFileName, "", "", ""))

	ipAllowDat := []ipAllowData{}

	// default for coalesce_ipv4 = 24, 5 and for ipv6 48, 5; override with the parameters in the server profile.
	coalesceMaskLenV4 := DefaultCoalesceMaskLenV4
	coalesceNumberV4 := DefaultCoalesceNumberV4
	coalesceMaskLenV6 := DefaultCoalesceMaskLenV6
	coalesceNumberV6 := DefaultCoalesceNumberV6

	for name, vals := range params {
		for _, val := range vals {
			switch name {
			case ParamPurgeAllowIP:
				for _, ip := range strings.Split(val, ",") {
					ipAllowDat = append(ipAllowDat, allowAll(strings.TrimSpace(ip)))
				}
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
		ipAllowDat = append([]ipAllowData{allowAll(`127.0.0.1`)}, ipAllowDat...)
		ipAllowDat = append([]ipAllowData{allowAll(`::1`)}, ipAllowDat...)
		ipAllowDat = append(ipAllowDat, allowAllButPushPurgeDelete(`0.0.0.0-255.255.255.255`))
		ipAllowDat = append(ipAllowDat, allowAllButPushPurgeDelete(`::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff`))
	} else {

		ips := []*net.IPNet{}
		ip6s := []*net.IPNet{}

		cgMap := map[string]tc.CacheGroupNullableV5{}
		for _, cg := range cacheGroups {
			if cg.Name == nil {
				return Cfg{}, makeErr(warnings, "got cachegroup with nil name!")
			}
			cgMap[*cg.Name] = cg
		}

		if server.CacheGroup == "" {
			return Cfg{}, makeErr(warnings, "server had nil Cachegroup!")
		}

		serverCG, ok := cgMap[server.CacheGroup]
		if !ok {
			return Cfg{}, makeErr(warnings, "server cachegroup not in cachegroups!")
		}

		childCGNames := getTopologyDirectChildren(tc.CacheGroupName(server.CacheGroup), topologies)

		childCGs := map[string]tc.CacheGroupNullableV5{}
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
			if childServer.CacheGroup == "" {
				warnings = append(warnings, "Servers had server with nil Cachegroup, skipping!")
				continue
			} else if childServer.HostName == "" {
				warnings = append(warnings, "Servers had server with nil HostName, skipping!")
				continue
			}

			// We need to add IPs to the allow of
			// - all children of this server
			// - all monitors, if this server is a Mid
			//
			_, isChild := childCGs[childServer.CacheGroup]
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
							warnings = append(warnings, "server '"+server.HostName+"' IP '"+svAddr.Address+" is not an IP address or CIDR - skipping!")
						} else if ip == nil {
							// not a CIDR or IP - error out
							warnings = append(warnings, "server '"+server.HostName+"' IP '"+svAddr.Address+" failed to parse as IP or CIDR - skipping!")
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
			ipAllowDat = append(ipAllowDat, allowAllButPushPurge(util.RangeStr(cidr)))
		}
		for _, cidr := range cidr6s {
			ipAllowDat = append(ipAllowDat, allowAllButPushPurge(util.RangeStr(cidr)))
		}

		// allow RFC 1918 server space - TODO JvD: parameterize
		ipAllowDat = append(ipAllowDat, allowAllButPushPurge(`10.0.0.0-10.255.255.255`))
		ipAllowDat = append(ipAllowDat, allowAllButPushPurge(`172.16.0.0-172.31.255.255`))
		ipAllowDat = append(ipAllowDat, allowAllButPushPurge(`192.168.0.0-192.168.255.255`))

		// order matters, so sort before adding the denys
		sort.Sort(ipAllowDatas(ipAllowDat))

		// start by allowing everything to localhost, including PURGE and PUSH
		ipAllowDat = append([]ipAllowData{allowAll(`127.0.0.1`)}, ipAllowDat...)
		ipAllowDat = append([]ipAllowData{allowAll(`::1`)}, ipAllowDat...)

		// end with a deny
		ipAllowDat = append(ipAllowDat, denyAll(`0.0.0.0-255.255.255.255`))
		ipAllowDat = append(ipAllowDat, denyAll(`::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff`))
	}

	text := makeHdrComment(opt.HdrComment)
	for _, al := range ipAllowDat {
		text += `src_ip=` + al.Src + ` action=` + al.Action + ` method=` + al.Method + "\n"
	}

	return Cfg{
		Text:        text,
		ContentType: ContentTypeHostingDotConfig,
		LineComment: LineCommentHostingDotConfig,
		Warnings:    warnings,
	}, nil
}

type ipAllowData struct {
	Src    string
	Action string
	Method string
}

type ipAllowDatas []ipAllowData

func (is ipAllowDatas) Len() int      { return len(is) }
func (is ipAllowDatas) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is ipAllowDatas) Less(i, j int) bool {
	if is[i].Src != is[j].Src {
		return is[i].Src < is[j].Src
	}
	if is[i].Action != is[j].Action {
		return is[i].Action < is[j].Action
	}
	return is[i].Method < is[j].Method
}

type ipAllowServer struct {
	IPAddress  string
	IP6Address string
}

type serversSortByName []Server

func (ss serversSortByName) Len() int      { return len(ss) }
func (ss serversSortByName) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }
func (ss serversSortByName) Less(i, j int) bool {
	if ss[j].HostName == "" {
		return false
	} else if ss[i].HostName == "" {
		return true
	}
	return ss[i].HostName < ss[j].HostName
}

// ActionAllow specifies that IPs matching the line's pattern should be allowed
// to use the explicitly listed HTTP request methods.
const ActionAllow = "ip_allow"

// ActionDeny specifies that IPs matching the line's pattern should NOT be
// allowed to use the explicitly listed HTTP request methods.
const ActionDeny = "ip_deny"

// These are the ATS HTTP request methods that may be granted using ipallow
// configuration file rules.
const (
	// "ALL" has the special meaning "all HTTP Methods".
	MethodAll = "ALL"
	// MethodPush is a non-standard HTTP request method understood by Apache
	// Traffic Server used to directly insert content for the requested URI into
	// the cache.
	MethodPush = "PUSH"
	// MethodPurge is a non-standard HTTP request method understood by Apache
	// Traffic Server used to directly remove the requested URI from the cache.
	MethodPurge  = "PURGE"
	MethodDelete = "DELETE"
)

// MethodSeparator is the string used as an infix operator in ipallow.config ATS
// configuration files to express a union of allowed HTTP request methods.
const MethodSeparator = `|`

// allowAllButPushPurge is a helper func to build a ipAllowData for the given range string immediately allowing all Methods except Push and Purge.
func allowAllButPushPurge(rangeStr string) ipAllowData {
	// Note denying methods implicitly and immediately allows all other methods!
	// So Deny PUSH|PURGE will make all other methods
	// immediately allowed, regardless of any later deny rules!
	methodPushPurge := strings.Join([]string{MethodPush, MethodPurge}, MethodSeparator)
	return ipAllowData{
		Src:    rangeStr,
		Action: ActionDeny,
		Method: methodPushPurge,
	}
}

// allowAllButPushPurgeDelete is a helper func to build a ipAllowData for the given range string immediately allowing all Methods except PUSH, PURGE, and DELETE.
func allowAllButPushPurgeDelete(rangeStr string) ipAllowData {
	// Note denying methods implicitly and immediately allows all other methods!
	// So Deny PUSH|PURGE will make all other methods
	// immediately allowed, regardless of any later deny rules!
	methodPushPurgeDelete := strings.Join([]string{MethodPush, MethodPurge, MethodDelete}, MethodSeparator)
	return ipAllowData{
		Src:    rangeStr,
		Action: ActionDeny,
		Method: methodPushPurgeDelete,
	}
}

// allowAll is a helper func to build a ipAllowData for the given range string immediately allowing all Methods, including Push and Purge.
func allowAll(rangeStr string) ipAllowData {
	return ipAllowData{
		Src:    rangeStr,
		Action: ActionAllow,
		Method: MethodAll,
	}
}

// denyAll is a helper func to build a ipAllowData for the given range string immediately denying all Methods.
func denyAll(rangeStr string) ipAllowData {
	return ipAllowData{
		Src:    rangeStr,
		Action: ActionDeny,
		Method: MethodAll,
	}
}
