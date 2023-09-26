package util

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
	"bytes"
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

// BitsPerByte is the number of bits in a byte.
const BitsPerByte = 8

// CoalesceIPs combines ips into CIDRs, by combining overlapping networks into networks of size coalesceMaskLen, if there are at least coalesceNumber IPs in the larger mask.
func CoalesceIPs(ips []net.IP, coalesceNumber int, coalesceMaskLen int) []*net.IPNet {
	if len(ips) == 0 {
		return nil
	}

	maskIP := ips[0].To4()
	isV4 := maskIP != nil
	if maskIP == nil {
		maskIP = ips[0]
	}

	mask := net.CIDRMask(coalesceMaskLen, len(maskIP)*BitsPerByte)

	type IPNetSources struct {
		Net     *net.IPNet
		Sources []net.IP
	}

	nets := []IPNetSources{}

iploop:
	for _, ip := range ips {
		ipIsV4 := ip.To4() != nil
		if isV4 != ipIsV4 {
			log.Errorln("CoalesceIPs got both V4 and V6 IPs, ignoring IP '" + ip.String() + "'")
			continue
		}
		for i, net := range nets {
			if net.Net.Contains(ip) {
				nets[i].Sources = append(nets[i].Sources, ip)
				continue iploop
			}
		}

		ipnet := &net.IPNet{IP: ip.Mask(mask), Mask: mask}
		nets = append(nets, IPNetSources{ipnet, []net.IP{ip}})
	}

	finalNets := []*net.IPNet{}
	for _, ipnet := range nets {
		if len(ipnet.Sources) >= coalesceNumber {
			finalNets = append(finalNets, ipnet.Net)
			continue
		}

		for _, ip := range ipnet.Sources {
			finalNets = append(finalNets, IPToCIDR(ip))
		}
	}

	return finalNets
}

// CoalesceCIDRs coalesces cidrs into a smaller set of CIDRs, by combining overlapping networks into networks of size coalesceMaskLen, if there are at least coalesceNumber cidrs in the larger mask.
func CoalesceCIDRs(cidrs []*net.IPNet, coalesceNumber int, coalesceMaskLen int) []*net.IPNet {
	if len(cidrs) == 0 {
		return nil
	}

	maskIP := cidrs[0].IP.To4()
	isV4 := maskIP != nil
	if maskIP == nil {
		maskIP = cidrs[0].IP
	}

	mask := net.CIDRMask(coalesceMaskLen, len(maskIP)*BitsPerByte)

	type IPNetSources struct {
		Net     *net.IPNet
		Sources []*net.IPNet
	}

	nets := []IPNetSources{}

iploop:
	for _, cidr := range cidrs {
		ipIsV4 := cidr.IP.To4() != nil
		if isV4 != ipIsV4 {
			log.Errorln("CoalesceIPs got both V4 and V6 IPs, ignoring CIDR '" + cidr.String() + "'")
			continue
		}

		for i, net := range nets {
			if CIDRIsSubset(cidr, net.Net) {
				nets[i].Sources = append(nets[i].Sources, cidr)
				continue iploop
			}
			if CIDRIsSubset(net.Net, cidr) {
				// if the existing net is a subset of the new cidr, replace the existing net with our larger cidr
				nets[i].Net = cidr
				nets[i].Sources = append(nets[i].Sources, cidr)
				continue iploop
			}
		}

		// use the larger of the coalesceMaskLen and this cidr's mask
		largerMask := mask
		if bytes.Compare(cidr.Mask, mask) < 1 {

			// Note this means cidr.Mask is numerically smaller, but that actually means it's masking more things.
			// Note bytes.Compare is defined to be lexographical, and we need a bit-wise comparison, but that's actually the same.
			largerMask = cidr.Mask
		}

		ipnet := &net.IPNet{IP: cidr.IP.Mask(largerMask), Mask: largerMask}
		nets = append(nets, IPNetSources{ipnet, []*net.IPNet{cidr}})
	}

	finalNets := []*net.IPNet{}
	for _, ipnet := range nets {
		if len(ipnet.Sources) >= coalesceNumber {
			finalNets = append(finalNets, ipnet.Net)
			continue
		}

		for _, cidr := range ipnet.Sources {
			finalNets = append(finalNets, cidr)
		}
	}

	return finalNets
}

// CIDRIsSubset returns whether na is a subset (possibly improper) of nb.
func CIDRIsSubset(na *net.IPNet, nb *net.IPNet) bool {
	return nb.Contains(FirstIP(na)) && nb.Contains(LastIP(na))
}

// FirstIP returns the first IP in the CIDR.
// For example, The CIDR 192.0.2.0/24 returns 192.0.2.0.
func FirstIP(ipn *net.IPNet) net.IP {
	return ipn.IP.Mask(ipn.Mask)
}

// LastIP returns the last IP in the CIDR.
// For example, The CIDR 192.0.2.0/24 returns 192.0.2.255.
func LastIP(ipn *net.IPNet) net.IP {
	inverseMask := make([]byte, len(ipn.Mask), len(ipn.Mask))
	for i, b := range ipn.Mask {
		inverseMask[i] = b ^ 0xFF
	}

	maxIPBts := make([]byte, len(ipn.IP), len(ipn.IP))

	for i, b := range ipn.IP {
		maxIPBts[i] = b | inverseMask[i]
	}

	maxIP := net.IP(maxIPBts)
	return maxIP
}

// RangeStr returns the hyphenated range of IPs.
// For example, The CIDR 192.0.2.0/24 returns "192.0.2.0-192.0.2.255".
func RangeStr(ipn *net.IPNet) string {
	firstIP := FirstIP(ipn)
	lastIP := LastIP(ipn)
	if firstIP.Equal(lastIP) {
		return firstIP.String()
	}
	return firstIP.String() + "-" + lastIP.String()
}

// IPToCIDR returns the CIDR containing just the given IP. For IPv6, this means /128, for IPv4, /32.
func IPToCIDR(ip net.IP) *net.IPNet {
	fullMask := net.IPMask([]byte{255, 255, 255, 255})
	if isV4 := ip.To4() != nil; !isV4 {
		fullMask = net.IPMask([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255})
	}
	return &net.IPNet{IP: ip, Mask: fullMask}
}

// IP4ToNum converts the passed string to a 32-bit unsigned integer where each
// byte that makes up the number is one of the bytes of the IPv4 address.
//
// The address is encoded with each byte left-to-right making up the
// most-to-least significant bytes in the resulting number. If the passed
// string cannot be parsed as an IPv4 address in standard notation, an error is
// returned.
func IP4ToNum(ip string) (uint32, error) {
	parts := strings.Split(ip, `.`)
	if len(parts) != 4 {
		return 0, errors.New("malformed IPv4")
	}
	intParts := []uint32{}
	for _, part := range parts {
		i, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return 0, errors.New("malformed IPv4")
		}
		intParts = append(intParts, uint32(i))
	}

	num := intParts[3]
	num += intParts[2] << 8
	num += intParts[1] << 16
	num += intParts[0] << 24

	return num, nil
}

// IP4InRange checks if the given string IP address falls within the specified
// hyphen-delimited range.
//
// The range should be of the form "start-end" e.g. "192.0.2.0-192.0.2.255".
// If either the input IP address or either end of this range fail to parse as
// IP addresses - or if the range is malformed - an error is returned.
//
// The behavior of this utility is undefined if the start of the IP range does
// not encode via IP4ToNum to a lower number than the end of the range.
func IP4InRange(ip, ipRange string) (bool, error) {
	ab := strings.Split(ipRange, `-`)
	if len(ab) != 2 {
		if len(ab) == 1 { // no range check for equality
			return ip == ipRange, nil
		}
		return false, errors.New("malformed range")
	}
	ipNum, err := IP4ToNum(ip)
	if err != nil {
		return false, errors.New("malformed ip")
	}
	aNum, err := IP4ToNum(ab[0])
	if err != nil {
		return false, errors.New("malformed range (first part)")
	}
	bNum, err := IP4ToNum(ab[1])
	if err != nil {
		return false, errors.New("malformed range (second part)")
	}
	return ipNum >= aNum && ipNum <= bNum, nil
}
