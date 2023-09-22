// package sar implements a syn-ack-rst health ping.
// It sends a TCP SYN, waits for an ACK, then immediately sends an RST to kill the connection.
// The primary purpose of this is as a health check, to verify the remote host is reachable, and able and willing to respond.
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
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-llog"
)

// SAR calls SARAddr if host is an IP address, or SARHost if host is an FQDN.
func SAR(log llog.Log, host string, port int, timeout time.Duration) (time.Duration, error) {
	log = llog.LibInit(log)
	if ip := net.ParseIP(host); ip != nil {
		return SARAddr(log, host, port, timeout)
	}
	return SARHost(log, host, port, timeout)
}

// SARAddr sends a syn-ack-reset to the given addr.
// The addr must be an IP.
// On success, the round-trip time to receive the Ack is returned, with a nil error.
// The addr is a string representation of an IPv4 or IPv6 address.
// TODO add optional local addr param
func SARAddr(log llog.Log, addr string, port int, timeout time.Duration) (time.Duration, error) {
	log = llog.LibInit(log)
	remoteAddr := net.ParseIP(addr)
	if remoteAddr == nil {
		return 0, errors.New("failed to parse addr '" + addr + "' as IP")
	}
	if v4 := remoteAddr.To4(); v4 != nil {
		remoteAddr = v4
	}

	localAddrStr, err := GetLocalAddr()
	if err != nil {
		return 0, errors.New("getting local address: " + err.Error())
	}

	localAddr := net.ParseIP(localAddrStr)
	if localAddr == nil {
		return 0, errors.New("failed to parse local addr '" + localAddrStr + "' as IP")
	}
	if v4 := localAddr.To4(); v4 != nil {
		localAddr = v4
	}

	ephemeralPortHolder, err := GetAndHoldEphemeralPort(localAddrStr)
	if err != nil {
		return 0, errors.New("failed to listen on ephemeral port: " + err.Error())
	}
	defer ephemeralPortHolder.Close()

	srcPort := ephemeralPortHolder.Port()

	// TODO implement Initial Sequence Number ISN per RFC9293§3.4.1
	seqNum := uint32(42)
	window := 256 * 10
	destPort := port
	dataOffset := 5 // because we have no options?

	native := TCPHdrNative{
		SrcPort:    uint16(srcPort),
		DestPort:   uint16(destPort),
		SeqNum:     seqNum,
		DataOffset: uint8(dataOffset), // 4 bits
		SYN:        true,
		Window:     uint16(window),
	}
	hdrBts, err := TCPHdrFromNative(native)
	if err != nil {
		return 0, errors.New("converting native header to byte: " + err.Error())
	}
	// the hdrBts packet is still missing the checksum, since we couldn't compute that until we had the marshalled bytes of the packet

	checksum := MakeTCPChecksum(hdrBts, localAddr, remoteAddr)

	hdrBts.SetChecksum(checksum)

	// wait for the recieving SAR goroutine to finish
	wg := sync.WaitGroup{}
	wg.Add(1)

	receiveTime := time.Time{}
	receiveSARErr := error(nil)
	// have to receive in a goroutine, because listening doesn't return, and we need to start listening before we send
	// (or else the response could get to our machine and be discarded between the send and listen)
	go func() {
		remoteAddrStr := addr
		localAddrStr := localAddr.String()
		receiveTime, receiveSARErr = receiveSAR(log, localAddrStr, remoteAddrStr, localAddr, remoteAddr, timeout, srcPort, destPort, int(seqNum+1))
		ephemeralPortHolder.Close() // this isn't strictly necessary, the defer will close this shortly. But this closes it ASAP
		wg.Done()
	}()

	sendTime, err := SendPacket(hdrBts, remoteAddr.String())
	if err != nil {
		return 0, errors.New("sending packet: " + err.Error())
	}

	wg.Wait()

	if receiveSARErr != nil {
		return 0, errors.New("receiving SAR: " + receiveSARErr.Error())
	}

	return receiveTime.Sub(sendTime), nil
}

// SARHost sends a syn-ack-reset to the given hostname.
// It looks up the hostname, and calls SendAddr on the first returned address.
// See SendAddr.
func SARHost(log llog.Log, host string, port int, timeout time.Duration) (time.Duration, error) {
	log = llog.LibInit(log)
	addrs, err := net.LookupHost(host)
	if err != nil {
		return 0, errors.New("lookup up host: " + err.Error())
	} else if len(addrs) == 0 {
		return 0, errors.New("looking up host succeeded, but no addresses were found.")
	}
	return SARAddr(log, addrs[0], port, timeout)
}

// GetLocalAddr gets a local IP, which this package will set as the TCP packet
// source, and then listen on this address to get the syn-ack response.
func GetLocalAddr() (string, error) {
	_, localAddr, err := FindNetInterfaceAddr()
	if err != nil {
		return "", errors.New("finding interface and address: " + err.Error())
	}
	laddr := strings.Split(localAddr.String(), "/")[0] // remove any CIDR
	return laddr, nil
}

// FindNetInterfaceAddr selects a local network interface to use.
// It picks the first interface that isn't loopback, is up, and has addresses.
// It then picks the first address in that interface.
// Returns the selected interface name, the selected address, and any error.
func FindNetInterfaceAddr() (string, net.Addr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", nil, errors.New("getting interfaces: " + err.Error())
	}
	if len(interfaces) == 0 {
		return "", nil, errors.New("no interfaces found")
	}

	errs := []string{}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		} else if len(addrs) == 0 {
			errs = append(errs, "interface '"+iface.Name+"' had no addresses")
			continue
		}
		return iface.Name, addrs[0], nil
	}
	return "", nil, errors.New("no interface found. Errors: " + strings.Join(errs, ", "))
}

// SendPacket sends the given packet to the given address and port.
//
// Note this doesn't take the local IP:port, which will be the source, which the destination will reply to,
// Because those must already be included in the packet.
//
// Returns the time immediately before the packet was sent, and any error.
func SendPacket(packet TCPHdr, addr string) (time.Time, error) {
	// TODO handle IPv6

	// Note this doesn't include the port. When dialing tcp, net.Dial addr must include the port.
	// But when dialing ip, the addr must not include the port.
	// Rather, the port will be determined from the dest port in the packet.
	conn, err := net.Dial("ip4:tcp", addr)
	if err != nil {
		return time.Time{}, errors.New("dialing address: " + err.Error())
	}
	defer conn.Close()

	sendTime := time.Now()
	numBtsSent, err := conn.Write(packet)
	if err != nil {
		return time.Time{}, errors.New("writing: " + err.Error())
	}
	if numBtsSent != len(packet) {
		return time.Time{}, fmt.Errorf("tried to write %v bytes, but only %v sent but no error", len(packet), numBtsSent)
	}
	return sendTime, nil
}

// receiveSAR receives the Syn-Ack in response to the Syn, with the given timeout.
// Upon receiving a SynAck, it sends a RST to help the remote discard the connection quickly, and stops caring.
// The srcPort, dstPort, and ackNum are used to match the SYNACK packet.
// Returns the time that the SynAck was received.
func receiveSAR(log llog.Log, localAddrStr string, remoteAddr string, localAddrIP net.IP, remoteAddrIP net.IP, timeout time.Duration, srcPort int, dstPort int, ackNum int) (time.Time, error) {
	localAddr, err := net.ResolveIPAddr("ip4", localAddrStr)
	if err != nil {
		return time.Time{}, errors.New("resolving local address: " + err.Error())
	}
	// listen on the local IP for the SynAck
	// TODO should this ListenPacket, to let Go automatically choose the port for us?
	conn, err := net.ListenIP("ip4:tcp", localAddr)
	if err != nil {
		return time.Time{}, errors.New("listening for syn-ack response: " + err.Error())
	}
	defer conn.Close()

	// Note this sets the absolute deadline, for all read packets.
	// If we wanted to set the timeout for each individual packet read,
	// we could call conn.SetDeadline inside the for-loop immediately before each conn.ReadFrom call.
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return time.Time{}, errors.New("setting deadline timeout to listen for syn-ack response: " + err.Error())
	}

	for {
		// the max TCP Header size is 60 bytes (IPv4 is variable up to 60, IPv6 is always 40)
		// This means the buffer here won't be just the header, it may also contain part of the payload.
		// We don't care here, our header parsing will stop after reading the final End Of Option List option.
		// But if this were changed to do something with the body, that would need taken into account.
		buf := make([]byte, 60)

		numBts, readRemoteAddr, err := conn.ReadFrom(buf)
		readTime := time.Now()
		if err != nil {
			return time.Time{}, errors.New("reading response: " + err.Error())
		}

		if readRemoteAddr.String() != remoteAddr {
			// this is normal. Listening to an IP makes us receive all packets for all apps using that local IP.
			// They're duplicated to us, so we're not causing anyone else to lose packets.
			// We just need to ignore them
			continue
		}

		tcpHdrRaw := TCPHdr(buf[:numBts])

		// note we swap src and dest ports. The source for the packet we sent is the dest we recieve, and vice versa.

		if int(tcpHdrRaw.SrcPort()) != dstPort {
			// log.Warnf("receiveSAR got mismatched source port %+v expecting %v, ignoring and continuing to listen\n", tcpHdrRaw.SrcPort(), dstPort)
			continue
		}
		if int(tcpHdrRaw.DestPort()) != srcPort {
			// log.Warnf("receiveSAR got mismatched dest port %+v expecting %v, ignoring and continuing to listen\n", tcpHdrRaw.DestPort(), srcPort)
			continue
		}
		if int(tcpHdrRaw.AckNum()) != ackNum {
			// log.Warnf("receiveSAR got mismatched ack num %+v expecting %v, ignoring and continuing to listen\n", tcpHdrRaw.AckNum(), ackNum)
			continue
		}

		tcpHdr, err := TCPHdrToNative(tcpHdrRaw)
		if err != nil {
			return time.Time{}, errors.New("decoding tcp header: " + err.Error())
		}

		if tcpHdr.RST {
			return readTime, nil // TODO should this return an error? We were expecting a syn-ack, not a reset
		}
		if tcpHdr.SYN && tcpHdr.ACK {

			// Note the sendRST below is commented because it isn't necessary:
			// the Linux kernel, even though we're listening on the port, will respond with an RST.
			//
			// Presumably because the kernel knows an ACK for something it doesn't know it sent
			// is malformed, and also because the AckNum doesn't match anything.
			//
			// I'm leaving the below sendRST commented to make it clear that an RST should be sent
			// to help the remote host free resources quickly;
			// but we don't need to, the kernel is doing it for us.

			// if err := sendRST(log, remoteAddr, localAddrIP, remoteAddrIP, int(tcpHdr.AckNum), srcPort, dstPort); err != nil {
			// 	log.Errorln("receiveSAR sendRST error: " + err.Error())
			// 	// don't fail. The RST is just to help the remote host close the connection gracefully.
			// 	// We still succeeded with sending a syn and getting a syn-ack, which is all we need on our end
			// }

			return readTime, nil
		}

		log.Warnln("receiveSAR got packet that wasn't an rst or a synack, but matched the ports and AckNum. Very strange!")
	}
}

// sendRST sends a tcp reset. This is designed to be called immediately after getting the syn-ack
// for the syn-ack-rst health check, to help the remote host close the connection quickly and gracefully.
//
// The seqNum must be the sequence number of the original syn packet plus 1.
func sendRST(log llog.Log, remoteAddrStr string, localAddr net.IP, remoteAddr net.IP, seqNum int, srcPort int, dstPort int) error {
	window := 43690 // TODO change

	dataOffset := 5 // because we have no options?
	native := TCPHdrNative{
		SrcPort:    uint16(srcPort),
		DestPort:   uint16(dstPort),
		SeqNum:     uint32(seqNum),
		DataOffset: uint8(dataOffset), // 4 bits
		RST:        true,
		Window:     uint16(window),
	}
	hdrBts, err := TCPHdrFromNative(native)
	if err != nil {
		return errors.New("converting native header to byte: " + err.Error())
	}
	// the hdrBts packet is still missing the checksum, since we couldn't compute that until we had the marshalled bytes of the packet

	checksum := MakeTCPChecksum(hdrBts, localAddr, remoteAddr)

	hdrBts.SetChecksum(checksum)

	if _, err := SendPacket(hdrBts, remoteAddrStr); err != nil {
		return errors.New("sending packet: " + err.Error())
	}
	return nil
}
