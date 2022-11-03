package sar

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
	"encoding/binary"
)

// TCPHdr is a TCP header
//
// TCPHdr should always be constructed with at least 20 bytes, the minimum TCP header size.
// Users constructing a TCPHdr and giving it to a trusting caller must ensure it is at least 20 bytes and valid.
// Users receiving an untrusted TCPHdr must check it is at least 20 bytes.
// This type and its functions do not check size or validity, for performance reasons, and therefore may
// panic if the TCPHdr is too small or otherwise malformed.
//
// Note TCPHdr does no validation of any kind, and is completely unaware of the contents of the packet,
// including the semantics of any TCP Options.
//
// For a more convenient but slower implementation, use TCPHdrDecoded (TODO implement)
//
// See Valid.
type TCPHdr []byte

func (th TCPHdr) SrcPort() uint16 { return binary.BigEndian.Uint16(th[:2]) }

func (th TCPHdr) DestPort() uint16 { return binary.BigEndian.Uint16(th[2:4]) }

func (th TCPHdr) SeqNum() uint32 {
	return binary.BigEndian.Uint32(th[4:8])
}

func (th TCPHdr) AckNum() uint32 {
	return binary.BigEndian.Uint32(th[8:12])
}

// DataOffset is the DOffset TCP field. See RFC9293§3.1.
//
// Note the Data Offset field is 4 bits; unfortunately, uint8 is the smallest Go type.
// Therefore, the most significant 4 bits will be unused and are not part of the field.
func (th TCPHdr) DataOffset() uint8 {
	do := th[12]         // don't need an Endian convertion, it's just a single byte
	do = do & 0b00001111 // We could just not do this mask. Callers shouldn't be looking at the higher bits anyway. Is this ever a security issue?
	return do
}

// Reserved is the reserved TCP field. See RFC9293§3.1.
//
// Note the Reserved field is 4 bits; unfortunately, uint8 is the smallest Go type.
// Therefore, the most significant 4 bits will be unused and are not part of the field.
func (th TCPHdr) Reserved() uint8 {
	do := th[12] // don't need an Endian convertion, it's just a single byte
	do = do >> 4
	return do
}

// Control is the control flag TCP fields: CWR, ECE, URG, ACK, PSH, RST, SYN, FIN.
//
// This provided for convenience and performance, when bit-shifting is faster or more convenient
// than using booleans. Each control flag is also available via individual boolean funcs.
func (th TCPHdr) Control() uint8 {
	return th[13]
}

// CWR is the Congestion Window Reduced flag, part of the ECN Explicit Congestion Notification.
// See RFC3168§6.1, RFC9293§3.1.
func (th TCPHdr) CWR() bool {
	return th[13]&0b10000000 != 0
}

// ECE is the ECN-Echo flag, part of the ECN Explicit Congestion Notification.
// See RFC3168§6.1, RFC9293§3.1.
func (th TCPHdr) ECE() bool {
	return th[13]&0b01000000 != 0
}

// URG is the Urgent flag. See RFC9293§3.1.
func (th TCPHdr) URG() bool {
	return th[13]&0b00100000 != 0
}

// ACK is the Acknowledgement flag. See RFC9293§3.1.
func (th TCPHdr) ACK() bool {
	return th[13]&0b00010000 != 0
}

// PSH is the Push flag. See RFC9293§3.1, RFC9293§3.9.1
func (th TCPHdr) PSH() bool {
	return th[13]&0b00001000 != 0
}

// RST is the connection Reset flag. See RFC9293§3.1.
func (th TCPHdr) RST() bool {
	return th[13]&0b00000100 != 0
}

// SYN is the synchronize sequence numbers flag. See RFC9293§3.1.
func (th TCPHdr) SYN() bool {
	return th[13]&0b00000010 != 0
}

// FIN is the finish connection flag. See RFC9293§3.1.
func (th TCPHdr) FIN() bool {
	return th[13]&0b00000001 != 0
}

func (th TCPHdr) Window() uint16 {
	return binary.BigEndian.Uint16(th[14:16])
}

func (th TCPHdr) Checksum() uint16 {
	return binary.BigEndian.Uint16(th[16:18])
}

func (th TCPHdr) SetChecksum(cs uint16) {
	binary.BigEndian.PutUint16(th[16:18], cs)
}

func (th TCPHdr) Urgent() uint16 {
	return binary.BigEndian.Uint16(th[18:20])
}

// OptionKind returns the kind of the nth option. See RFC9293§3.2.
//
// The prevSize is the sum of the sizes of all options less than n.
// This is necessary to find the next option.
//
// Note the End of Option List the last option is of kind 0. So to iterate over all options, continuously request Option(prevSize, n) (tracking the sum of the previous sizes)
// until the returned option Kind is 0.
//
// Note per RFC9293§3.2 all options except kind 0 (End of Option) and kind 1 (No Operation) have lengths. Callers must detect this and not query OptionLength for Kind 0 or 1.
//
// See RFC9293§3.2 for more details.
func (th TCPHdr) OptionKind(prevSize int, n int) uint8 {
	return th[prevSize] // single byte, no Endian conversion
}

// OptionLen returns the kind of the nth option. See RFC9293§3.2.
//
// The prevSize is the sum of the sizes of all options less than n.
// This is necessary to find the next option.
//
// Note per RFC9293§3.2 all options except kind 0 (End of Option) and kind 1 (No Operation) have lengths. Callers must detect this and not query OptionLen for Kind 0 or 1.
//
// See RFC9293§3.2 for more details.
func (th TCPHdr) OptionLen(prevSize int, n int) uint8 {
	return th[prevSize+1] // single byte, no Endian conversion
}

// OptionData returns the data of the nth option. See RFC9293§3.2.
//
// The prevSize is the sum of the sizes of all options less than n.
// This is necessary to find the next option.
//
// Note per RFC9293§3.2 all options except kind 0 (End of Option) and kind 1 (No Operation) have lengths. Callers must detect this and not query OptionLength or OptionData for Kind 0 or 1.
// Note the Option Length includes the Kind and Length fields, each 1 octet. Therefore, options of length 2 have no data. Options of length <=2 must not call OptionData.
//
// The optionLen is the Length field. Note the Option Length includes the Kind and Length. The optionLen must not be the length of the data, it must be the Option Length field.
// Therefore, the returned data will be 2 octets less than optionLen.
//
// Note the OptionLength includes padding. Therefore, the data returned by this func will include padding, in addition to the semantic option data. This func is unaware of the semantics
// of any options.
//
// Note the OptionsData bytes are returned in network byte order (Big Endian), not machine byte order (frequently Little Endian).
// Because TCPHdr doesn't know about options, it lacks the context to know the sizes and locations of integers inside particular options,
// and therefore cannot convert byte order. Callers must convert byte order as necessary. See the binary package.
//
// See RFC9293§3.2 for more details.
func (th TCPHdr) OptionData(prevSize int, optionLen int, n int) []byte {
	return th[prevSize+2 : optionLen-2]
}

const TCPHdrOptionEndOfOptionList = 0
const TCPHdrOptionNoOperation = 1

// ProtocolNumberTCP is the protocol number for TCP.
// See IANA Protocol Numbers, RFC791.
const ProtocolNumberTCP = 6

// MakeTCPChecksum creates the TCP packet checksum. See RFC9293§3.1
// TODO handle IPv6
func MakeTCPChecksum(data []byte, sourceIP []byte, destIP []byte) uint16 {
	// TODO validate sourceIP and destIP are either 4 or 16 len?
	// TODO handle IPv6

	// TCP Checksum includes an IP "psuedo-header" (to protect against mis-routed packets). See RFC9293§3.1.
	pseudoHdr := []byte{
		sourceIP[0], sourceIP[1], sourceIP[2], sourceIP[3],
		destIP[0], destIP[1], destIP[2], destIP[3],
		0, // defined to be zero, see RFC9293§3.1.
		ProtocolNumberTCP,
		0, // tcp segment length, will calculate below
		0, // tcp segment length, will calculate below
	}
	binary.BigEndian.PutUint16(pseudoHdr[10:12], uint16(len(data))) // tcp segment length
	// TODO pad pseudoHdr?

	btsToSum := append(pseudoHdr, data...)

	lenBtsToSum := len(btsToSum)
	nextWord := uint16(0)
	sum := uint32(0)
	for i := 0; i+1 < lenBtsToSum; i += 2 {
		nextWord = uint16(btsToSum[i])<<8 | uint16(btsToSum[i+1])
		sum += uint32(nextWord)
	}
	if lenBtsToSum%2 != 0 { // is there an odd byte?
		sum += uint32(btsToSum[len(btsToSum)-1])
	}

	// add carry (if any), which itself may have a carry so add that too (if any)
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)

	// bitwise complement
	return uint16(^sum)
}
