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
	"fmt"
)

// TCPHdrNative is a convenience struct containing deserialized TCP Header data.
// This is more convenient but less efficient than TCPHdr.
type TCPHdrNative struct {
	SrcPort    uint16
	DestPort   uint16
	SeqNum     uint32
	AckNum     uint32
	DataOffset uint8 // 4 bits
	Reserved   uint8 // 4 bits
	CWR        bool
	ECE        bool
	URG        bool
	ACK        bool
	PSH        bool
	RST        bool
	SYN        bool
	FIN        bool
	Window     uint16
	Checksum   uint16 // Kernel will set this if it's 0
	Urgent     uint16
	Options    []TCPHdrNativeOption
}

type TCPHdrNativeOption struct {
	Kind int
	Len  int
	Data []byte
}

// TCPHdrToNative creates a TCPHdrNative from a TCPHdr.
func TCPHdrToNative(hdr TCPHdr) (TCPHdrNative, error) {
	native := TCPHdrNative{}
	const minTCPHdrSize = 20
	if len(hdr) < minTCPHdrSize {
		return TCPHdrNative{}, fmt.Errorf("malformed header, minimum TCP header size is %v but hdr was %v bytes", minTCPHdrSize, len(hdr))
	}

	native.SrcPort = hdr.SrcPort()
	native.DestPort = hdr.DestPort()
	native.SeqNum = hdr.SeqNum()
	native.AckNum = hdr.AckNum()
	native.DataOffset = hdr.DataOffset()
	native.Reserved = hdr.Reserved()
	native.CWR = hdr.CWR()
	native.ECE = hdr.ECE()
	native.URG = hdr.URG()
	native.ACK = hdr.ACK()
	native.PSH = hdr.PSH()
	native.RST = hdr.RST()
	native.SYN = hdr.SYN()
	native.FIN = hdr.FIN()
	native.Window = hdr.Window()
	native.Checksum = hdr.Checksum()
	native.Urgent = hdr.Urgent()
	native.Options = []TCPHdrNativeOption{}

	if native.DataOffset < 5 {
		return native, nil // no options
	}

	prevOptionsSize := 0
	optionNum := 0
	for {
		// TODO add bounds checking, for malformed options

		option := TCPHdrNativeOption{}
		option.Kind = int(hdr.OptionKind(prevOptionsSize, 0))
		if option.Kind == TCPHdrOptionEndOfOptionList {
			// EoOL is length 1, no length octet and no data octets
			native.Options = append(native.Options, option)
			prevOptionsSize += 1
			optionNum++
			break
		}

		if option.Kind == TCPHdrOptionNoOperation {
			// NoOp is length 1, no length octet and no data octets
			native.Options = append(native.Options, option)
			prevOptionsSize += 1
			optionNum++
			continue
		}

		option.Len = int(hdr.OptionLen(prevOptionsSize, optionNum))
		option.Data = hdr.OptionData(prevOptionsSize, option.Len, optionNum)
		native.Options = append(native.Options, option)
		prevOptionsSize += option.Len
		optionNum++
	}
	return native, nil
}

// TCPHdrFromNative creates a TCPHdr bytes from a TCPHdrNative.
// The created TCPHdr is ready to send over the wire.
func TCPHdrFromNative(native TCPHdrNative) (TCPHdr, error) {
	// TODO pre-calculate and allocate options size
	hdr := TCPHdr(make([]byte, 20, 20)) // allocate the mandatory 20 bytes up front

	binary.BigEndian.PutUint16(hdr[0:2], uint16(native.SrcPort))
	binary.BigEndian.PutUint16(hdr[2:4], uint16(native.DestPort))
	binary.BigEndian.PutUint32(hdr[4:8], uint32(native.SeqNum))
	binary.BigEndian.PutUint32(hdr[8:12], uint32(native.AckNum))

	// TODO write a single uint32 for offset+reserved+control+window? Should be faster

	hdr[12] = (native.DataOffset << 4) | (native.Reserved) // single byte so no endian byte order conversion. offset is 4 bits and reserved is 4 bits

	// TODO this is terribly inefficient. Go has no way to convert bool to int without a conditional.
	//      The TCPHdr is designed to be efficient, and TCPHdrNative is designed to be convenient if less efficient,
	//      But it shouldn't be *that* inefficient. Maybe we should just use uint8?

	boolToByte := func(b bool) byte {
		if b {
			return 1
		}
		return 0
	}

	cwr := boolToByte(native.CWR)
	ece := boolToByte(native.ECE)
	urg := boolToByte(native.URG)
	ack := boolToByte(native.ACK)
	psh := boolToByte(native.PSH)
	rst := boolToByte(native.RST)
	syn := boolToByte(native.SYN)
	fin := boolToByte(native.FIN)

	hdr[13] = cwr | (ece << 1) | (urg << 2) | (ack << 3) | (psh << 4) | (rst << 5) | (syn << 6) | (fin << 7)

	hdr[13] = (cwr << 7) | (ece << 6) | (urg << 5) | (ack << 4) | (psh << 3) | (rst << 2) | (syn << 1) | (fin)
	binary.BigEndian.PutUint16(hdr[14:16], native.Window)
	binary.BigEndian.PutUint16(hdr[16:18], native.Checksum)
	binary.BigEndian.PutUint16(hdr[18:20], native.Urgent)

	for _, opt := range native.Options {
		hdr = append(hdr, uint8(opt.Kind)) // uint8, no endian byte order conversion
		if opt.Len == 0 {
			continue
		}
		hdr = append(hdr, uint8(opt.Len))
		hdr = append(hdr, opt.Data[:opt.Len-2]...) // -2 because Len includes the 1-octet kind and 1-octet length fields.
	}
	return hdr, nil
}
