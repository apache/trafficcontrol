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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

const IPAllowConfigFileName = `ip_allow.config`
const ContentTypeIPAllowDotConfig = ContentTypeTextASCII
const LineCommentIPAllowDotConfig = LineCommentHash

type IPAllowData struct {
	Src    string
	Action string
	Method string
}

type IPAllowDatas []IPAllowData

func (is IPAllowDatas) Len() int      { return len(is) }
func (is IPAllowDatas) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is IPAllowDatas) Less(i, j int) bool {
	if is[i].Src != is[j].Src {
		return is[i].Src < is[j].Src
	}
	if is[i].Action != is[j].Action {
		return is[i].Action < is[j].Action
	}
	return is[i].Method < is[j].Method
}

const ParamPurgeAllowIP = "purge_allow_ip"
const ParamCoalesceMaskLenV4 = "coalesce_masklen_v4"
const ParamCoalesceNumberV4 = "coalesce_number_v4"
const ParamCoalesceMaskLenV6 = "coalesce_masklen_v6"
const ParamCoalesceNumberV6 = "coalesce_number_v6"

type IPAllowServer struct {
	IPAddress  string
	IP6Address string
}

const DefaultCoalesceMaskLenV4 = 24
const DefaultCoalesceNumberV4 = 5
const DefaultCoalesceMaskLenV6 = 48
const DefaultCoalesceNumberV6 = 5

type ServersSortByName []tc.ServerNullable

func (ss ServersSortByName) Len() int      { return len(ss) }
func (ss ServersSortByName) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }
func (ss ServersSortByName) Less(i, j int) bool {
	if ss[j].HostName == nil {
		return false
	} else if ss[i].HostName == nil {
		return true
	}
	return *ss[i].HostName < *ss[j].HostName
}

// MakeIPAllowDotConfig creates the ip_allow.config ATS config file.
// The childServers is a list of servers which are children for this Mid-tier server. This should be empty for Edge servers.
// More specifically, it should be the list of edges whose cachegroup's parent_cachegroup or secondary_parent_cachegroup is the cachegroup of this Mid server.
func MakeIPAllowDotConfig(
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	params map[string][]string, // map[name]value - config file should always be ip_allow.config
	server *tc.ServerNullable,
	servers []tc.ServerNullable,
	cacheGroups []tc.CacheGroupNullable,
) string {
	if server.HostName == nil {
		return "ERROR: server missing hostname"
	}

	ipAllowData := []IPAllowData{}
	const ActionAllow = "ip_allow"
	const ActionDeny = "ip_deny"
	const MethodAll = "ALL"

	// localhost is trusted.
	ipAllowData = append(ipAllowData, IPAllowData{
		Src:    `127.0.0.1`,
		Action: ActionAllow,
		Method: MethodAll,
	})
	ipAllowData = append(ipAllowData, IPAllowData{
		Src:    `::1`,
		Action: ActionAllow,
		Method: MethodAll,
	})

	// default for coalesce_ipv4 = 24, 5 and for ipv6 48, 5; override with the parameters in the server profile.
	coalesceMaskLenV4 := DefaultCoalesceMaskLenV4
	coalesceNumberV4 := DefaultCoalesceNumberV4
	coalesceMaskLenV6 := DefaultCoalesceMaskLenV6
	coalesceNumberV6 := DefaultCoalesceNumberV6

	for name, vals := range params {
		for _, val := range vals {
			switch name {
			case "purge_allow_ip":
				ipAllowData = append(ipAllowData, IPAllowData{
					Src:    val,
					Action: ActionAllow,
					Method: MethodAll,
				})
			case ParamCoalesceMaskLenV4:
				if vi, err := strconv.Atoi(val); err != nil {
					log.Warnln("MakeIPAllowDotConfig got param '" + name + "' val '" + val + "' not a number, ignoring!")
				} else if coalesceMaskLenV4 != DefaultCoalesceMaskLenV4 {
					log.Warnln("MakeIPAllowDotConfig got multiple param '" + name + "' - ignoring  val '" + val + "'!")
				} else {
					coalesceMaskLenV4 = vi
				}
			case ParamCoalesceNumberV4:
				if vi, err := strconv.Atoi(val); err != nil {
					log.Warnln("MakeIPAllowDotConfig got param '" + name + "' val '" + val + "' not a number, ignoring!")
				} else if coalesceNumberV4 != DefaultCoalesceNumberV4 {
					log.Warnln("MakeIPAllowDotConfig got multiple param '" + name + "' - ignoring  val '" + val + "'!")
				} else {
					coalesceNumberV4 = vi
				}
			case ParamCoalesceMaskLenV6:
				if vi, err := strconv.Atoi(val); err != nil {
					log.Warnln("MakeIPAllowDotConfig got param '" + name + "' val '" + val + "' not a number, ignoring!")
				} else if coalesceMaskLenV6 != DefaultCoalesceMaskLenV6 {
					log.Warnln("MakeIPAllowDotConfig got multiple param '" + name + "' - ignoring  val '" + val + "'!")
				} else {
					coalesceMaskLenV6 = vi
				}
			case ParamCoalesceNumberV6:
				if vi, err := strconv.Atoi(val); err != nil {
					log.Warnln("MakeIPAllowDotConfig got param '" + name + "' val '" + val + "' not a number, ignoring!")
				} else if coalesceNumberV6 != DefaultCoalesceNumberV6 {
					log.Warnln("MakeIPAllowDotConfig got multiple param '" + name + "' - ignoring  val '" + val + "'!")
				} else {
					coalesceNumberV6 = vi
				}
			}
		}
	}

	// for edges deny "PUSH|PURGE|DELETE", allow everything else to everyone.
	isMid := strings.HasPrefix(server.Type, tc.MidTypePrefix)
	if !isMid {
		ipAllowData = append(ipAllowData, IPAllowData{
			Src:    `0.0.0.0-255.255.255.255`,
			Action: ActionDeny,
			Method: `PUSH|PURGE|DELETE`,
		})
		ipAllowData = append(ipAllowData, IPAllowData{
			Src:    `::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff`,
			Action: ActionDeny,
			Method: `PUSH|PURGE|DELETE`,
		})
	} else {

		ips := []*net.IPNet{}
		ip6s := []*net.IPNet{}

		cgMap := map[string]tc.CacheGroupNullable{}
		for _, cg := range cacheGroups {
			if cg.Name == nil {
				return "ERROR: got cachegroup with nil name!'"
			}
			cgMap[*cg.Name] = cg
		}

		if server.Cachegroup == nil {
			return "ERROR: server had nil Cachegroup!"
		}

		serverCG, ok := cgMap[*server.Cachegroup]
		if !ok {
			return "ERROR: Server cachegroup not in cachegroups!"
		}

		childCGs := map[string]tc.CacheGroupNullable{}
		for cgName, cg := range cgMap {
			if (cg.ParentName != nil && *cg.ParentName == *serverCG.Name) || (cg.SecondaryParentName != nil && *cg.SecondaryParentName == *serverCG.Name) {
				childCGs[cgName] = cg
			}
		}

		// sort servers, to guarantee things like IP coalescing are deterministic
		sort.Sort(ServersSortByName(servers))
		for _, childServer := range servers {
			if childServer.Cachegroup == nil {
				log.Errorln("Servers had server with nil Cachegroup, skipping!")
				continue
			} else if childServer.HostName == nil {
				log.Errorln("Servers had server with nil HostName, skipping!")
				continue
			}

			// We need to add IPs to the allow of
			// - all children of this server
			// - all monitors, if this server is a Mid
			//
			// TODO: handle Topologies. Mids currently block everything but Edges
			//       We should decide how to handle that in a post-edge-mid world.
			//       That probably means adding all child and monitor IPs, and blocking everything else,
			//       for all non-first-tier caches.
			//
			_, isChild := childCGs[*childServer.Cachegroup]
			if !isChild && (!strings.HasPrefix(server.Type, tc.MidTypePrefix) || (string(childServer.Type) != tc.MonitorTypeName)) {
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
							log.Errorln("MakeIPAllowDotConfig server '" + *server.HostName + "' IP '" + svAddr.Address + " is not an IP address or CIDR - skipping!")
						} else if ip == nil {
							// not a CIDR or IP - error out
							log.Errorln("MakeIPAllowDotConfig server '" + *server.HostName + "' IP '" + svAddr.Address + " failed to parse as IP or CIDR - skipping!")
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
			ipAllowData = append(ipAllowData, IPAllowData{
				Src:    util.RangeStr(cidr),
				Action: ActionAllow,
				Method: MethodAll,
			})
		}
		for _, cidr := range cidr6s {
			ipAllowData = append(ipAllowData, IPAllowData{
				Src:    util.RangeStr(cidr),
				Action: ActionAllow,
				Method: MethodAll,
			})
		}

		// allow RFC 1918 server space - TODO JvD: parameterize
		ipAllowData = append(ipAllowData, IPAllowData{
			Src:    `10.0.0.0-10.255.255.255`,
			Action: ActionAllow,
			Method: MethodAll,
		})
		ipAllowData = append(ipAllowData, IPAllowData{
			Src:    `172.16.0.0-172.31.255.255`,
			Action: ActionAllow,
			Method: MethodAll,
		})
		ipAllowData = append(ipAllowData, IPAllowData{
			Src:    `192.168.0.0-192.168.255.255`,
			Action: ActionAllow,
			Method: MethodAll,
		})

		// order matters, so sort before adding the denys
		sort.Sort(IPAllowDatas(ipAllowData))

		// end with a deny
		ipAllowData = append(ipAllowData, IPAllowData{
			Src:    `0.0.0.0-255.255.255.255`,
			Action: ActionDeny,
			Method: MethodAll,
		})
		ipAllowData = append(ipAllowData, IPAllowData{
			Src:    `::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff`,
			Action: ActionDeny,
			Method: MethodAll,
		})
	}

	text := GenericHeaderComment(*server.HostName, toToolName, toURL)
	for _, al := range ipAllowData {
		text += `src_ip=` + al.Src + ` action=` + al.Action + ` method=` + al.Method + "\n"
	}
	return text
}
