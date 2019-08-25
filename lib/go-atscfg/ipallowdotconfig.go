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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

const IPAllowConfigFileName = `ip_allow.config`

type IPAllowData struct {
	Src    string
	Action string
	Method string
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

// MakeIPAllowDotConfig creates the ip_allow.config ATS config file.
// The childServers is a list of servers which are children for this Mid-tier server. This should be empty for Edge servers.
// More specifically, it should be the list of edges whose cachegroup's parent_cachegroup or secondary_parent_cachegroup is the cachegroup of this Mid server.
func MakeIPAllowDotConfig(
	serverName tc.CacheName,
	serverType tc.CacheType,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	params map[string][]string, // map[name]value - config file should always be ip_allow.config
	childServers map[tc.CacheName]IPAllowServer,
) string {
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
				} else if coalesceNumberV6 != DefaultCoalesceMaskLenV6 {
					log.Warnln("MakeIPAllowDotConfig got multiple param '" + name + "' - ignoring  val '" + val + "'!")
				} else {
					coalesceNumberV6 = vi
				}
			}
		}
	}

	// for edges deny "PUSH|PURGE|DELETE", allow everything else to everyone.
	isMid := strings.HasPrefix(string(serverType), tc.MidTypePrefix)
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
		for serverName, server := range childServers {

			if ip := net.ParseIP(server.IPAddress).To4(); ip != nil {
				// got an IP - convert it to a CIDR and add it to the list
				ips = append(ips, util.IPToCIDR(ip))
			} else {
				// not an IP, try a CIDR
				if ip, cidr, err := net.ParseCIDR(server.IPAddress); err != nil {
					// not a CIDR or IP - error out
					log.Errorln("MakeIPAllowDotConfig server '" + string(serverName) + "' IP '" + server.IPAddress + " is not an IPv4 address or CIDR - skipping!")
				} else {
					// got a valid CIDR - now make sure it's v4
					ip = ip.To4()
					if ip == nil {
						// valid CIDR, but not v4
						log.Errorln("MakeIPAllowDotConfig server '" + string(serverName) + "' IP '" + server.IPAddress + " is a CIDR, but not v4 - skipping!")
					} else {
						// got a valid IPv4 CIDR - add it to the list
						ips = append(ips, cidr)
					}
				}
			}

			if server.IP6Address != "" {
				ip6 := net.ParseIP(server.IP6Address)
				if ip6 != nil && ip6.To4() == nil {
					// got a valid IPv6 - add it to the list
					ip6s = append(ip6s, util.IPToCIDR(ip6))
				} else {
					// not a v6 IP, try a CIDR
					if ip, cidr, err := net.ParseCIDR(server.IP6Address); err != nil {
						// not a CIDR or IP - error out
						log.Errorln("MakeIPAllowDotConfig server '" + string(serverName) + "' IP6 '" + server.IP6Address + " is not an IPv6 address or CIDR - skipping!")
					} else {
						// got a valid CIDR - now make sure it's v6
						ip = ip.To4()
						if ip != nil {
							// valid CIDR, but not v6
							log.Errorln("MakeIPAllowDotConfig server '" + string(serverName) + "' IP6 '" + server.IPAddress + " is a CIDR, but not v6 - skipping!")
						} else {
							// got a valid IPv6 CIDR - add it to the list
							ip6s = append(ip6s, cidr)
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

	text := GenericHeaderComment(string(serverName), toToolName, toURL)
	for _, al := range ipAllowData {
		text += `src_ip=` + al.Src + ` action=` + al.Action + ` method=` + al.Method + "\n"
	}
	return text
}
