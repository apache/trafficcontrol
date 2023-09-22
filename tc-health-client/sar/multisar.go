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
	"net"
	"os"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-llog"
)

type HostPort struct {
	Host string
	Port int
}

type SARResult struct {
	Host string
	Port int
	RTT  time.Duration
	Err  error
}

// MultiSAR is like SAR for multiple requests.
// SAR has to listen on a raw IP port without a TCP socket, which is relatively inexpensive for a single request,
// but expensive for large numbers of requests.
// MultiSAR uses a single listener on an ephemeral local port for all SAR requests, significantly reducing
// resource costs.
func MultiSAR(log llog.Log, hosts []HostPort, timeout time.Duration) ([]SARResult, error) {
	log = llog.LibInit(log)

	localAddrStr, err := GetLocalAddr()
	if err != nil {
		return nil, errors.New("getting local address: " + err.Error())
	}

	localAddr := net.ParseIP(localAddrStr)
	if localAddr == nil {
		return nil, errors.New("failed to parse local addr '" + localAddrStr + "' as IP")
	}
	if v4 := localAddr.To4(); v4 != nil {
		localAddr = v4
	}

	ephemeralPortHolder, err := GetAndHoldEphemeralPort(localAddrStr)
	if err != nil {
		return nil, errors.New("failed to listen on ephemeral port: " + err.Error())
	}
	defer ephemeralPortHolder.Close()

	srcPort := ephemeralPortHolder.Port()

	// pre-construct all the packets, so we listen for as little time as possible

	// TODO implement Initial Sequence Number ISN per RFC9293§3.4.1? It might be faster to use the same seq num for all packets
	seqNum := uint32(42)

	type HostPortPacket struct {
		HostPort
		TCPHdr TCPHdr
	}

	packets := []HostPortPacket{}

	hostAddr := map[string]string{}    // map[host]addr - note hosts may be IPs, in which case host and addr will be the same
	addrHosts := map[string][]string{} // map[addr][]host - note multiple FQDNs may have the same IP

	results := []SARResult{}

	for _, host := range hosts {
		makeHostErrResult := func(err error) SARResult {
			return SARResult{
				Host: host.Host,
				Port: host.Port,
				RTT:  0,
				Err:  err,
			}
		}

		remoteAddr := net.ParseIP(host.Host)
		if remoteAddr != nil {
			// host is IP
			hostAddr[host.Host] = host.Host
			addrHosts[host.Host] = append(addrHosts[host.Host], host.Host)
		} else {
			// host isn't an IP, assume FQDN
			addrs, err := net.LookupHost(host.Host)
			if err != nil {
				results = append(results, makeHostErrResult(errors.New("lookup up host '"+host.Host+"': "+err.Error())))
				continue
			}
			if len(addrs) == 0 {
				results = append(results, makeHostErrResult(errors.New("looking up host '"+host.Host+"' succeeded, but no addresses were found.")))
				continue
			}
			remoteAddr = net.ParseIP(addrs[0])
			if remoteAddr == nil {
				results = append(results, makeHostErrResult(errors.New("failed to parse addr '"+host.Host+"' ip '"+addrs[0]+"' as IP")))
				continue
			}

			hostAddr[host.Host] = addrs[0]
			addrHosts[addrs[0]] = append(addrHosts[addrs[0]], host.Host)
		}

		if v4 := remoteAddr.To4(); v4 != nil {
			remoteAddr = v4
		}

		// TODO handle IPv6

		window := 256 * 10
		destPort := host.Port
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
			return nil, errors.New("converting native header to byte: " + err.Error())
		}
		hdrBts.SetChecksum(MakeTCPChecksum(hdrBts, localAddr, remoteAddr))
		packets = append(packets, HostPortPacket{HostPort: host, TCPHdr: hdrBts})
	}

	remoteInf := []sarRemoteInf{}
	// Note we need to iterate over hostAddr, *not* hosts.
	// The hosts has all, but we may have failed to resolve some, in which case they won't be sent and we must not listen for their responses.
	for _, host := range hosts {
		addr, ok := hostAddr[host.Host]
		if !ok {
			continue // if it's not in hostAddr, we already added an error result and we won't send a packet, so don't listen for it
		}
		remoteInf = append(remoteInf, sarRemoteInf{
			Host:   host.Host,
			Addr:   addr,
			Port:   host.Port,
			AckNum: seqNum + 1,
		})
	}

	sarListenerResp := []sarListenerResp{}
	sarListenerErr := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(1)
	doneSending := make(chan struct{}, 1) // we don't want to start the timeout until after we send the last packet. Buffer 1, the main sending thread doesn't block
	go func() {
		sarListenerResp, sarListenerErr = sarListener(log, localAddrStr, localAddr, remoteInf, timeout, srcPort, doneSending)
		ephemeralPortHolder.Close() // this isn't strictly necessary, the defer will close this shortly. But this closes it ASAP
		wg.Done()
	}()

	sendTimes := map[HostPort]time.Time{}

	packetSendStart := time.Now()
	for _, packet := range packets {
		sendTime, err := SendPacket(packet.TCPHdr, packet.HostPort.Host)
		if err != nil {
			return nil, errors.New("sending packet: " + err.Error())
		}
		sendTimes[packet.HostPort] = sendTime
	}

	log.Infof("multisar main thread sent %v packets in %vms\n", len(packets), time.Since(packetSendStart)/time.Millisecond)

	doneSending <- struct{}{} // send the doneSending message to the listener, so it sets the timeout
	wg.Wait()                 // wait for the listener to return and set sarListenerResp and sarListenerErr

	if sarListenerErr != nil {
		return nil, errors.New("listening for ACKs: " + sarListenerErr.Error())
	}

	for _, listenResp := range sarListenerResp {
		hosts := addrHosts[listenResp.Addr] // multiple FQDNs may have the same IP, but at L4 we only care about the behavior of the IP
		for _, host := range hosts {
			sendTime, ok := sendTimes[HostPort{Host: host, Port: listenResp.Port}]
			if !ok {
				log.Errorf("SAR listener got packet that was never sent somehow! Should never happen! Response: %+v\n", listenResp)
				continue
			}
			roundTripTime := time.Duration(0)
			if listenResp.Err == nil {
				// this check shouldn't be necessary, the caller should never look at duration if err!=nil. This just makes it easier to debug if they do
				roundTripTime = listenResp.RespTime.Sub(sendTime)
			}
			results = append(results, SARResult{
				Host: host,
				Port: listenResp.Port,
				RTT:  roundTripTime,
				Err:  listenResp.Err,
			})
		}
	}
	return results, nil
}

type sarRemoteInf struct {
	Host   string
	Addr   string
	Port   int
	AckNum uint32
}

type sarListenerResp struct {
	Addr     string
	Port     int
	RespTime time.Time
	Err      error
}

func sarListener(log llog.Log, localAddrStr string, localAddrIP net.IP, remoteArr []sarRemoteInf, timeout time.Duration, srcPort int, doneSending <-chan struct{}) ([]sarListenerResp, error) {

	localAddr, err := net.ResolveIPAddr("ip4", localAddrStr)
	if err != nil {
		return nil, errors.New("resolving local address: " + err.Error())
	}

	// remotes is used to quickly, progressively match multiple requests
	remotes := map[string]map[int]uint32{} // map[remoteAddr][port]acknum
	for _, remote := range remoteArr {
		if _, ok := remotes[remote.Addr]; !ok {
			remotes[remote.Addr] = map[int]uint32{}
		}
		remotes[remote.Addr][remote.Port] = remote.AckNum
	}

	// listen on the local IP for the SynAck
	// TODO should this ListenPacket, to let Go automatically choose the port for us?
	conn, err := net.ListenIP("ip4:tcp", localAddr)
	if err != nil {
		return nil, errors.New("listening for syn-ack response: " + err.Error())
	}
	defer conn.Close()

	// TODO listen on "sent everything" chan, and don't set timeout until everything was sent,
	// to prevent timing out before some things are even sent

	responses := []sarListenerResp{}

	// the max TCP Header size is 60 bytes (IPv4 is variable up to 60, IPv6 is always 40)
	// This means the buffer here won't be just the header, it may also contain part of the payload.
	// We don't care here, our header parsing will stop after reading the final End Of Option List option.
	// But if this were changed to do something with the body, that would need taken into account.
	buf := make([]byte, 60)
	for {
		// after the main thread is done sending all packets, set the timeout
		select {
		case <-doneSending:
			log.Infof("multisar listener got doneSending, setting timeout '%v'\n", timeout)
			// Note this sets the absolute deadline, for all read packets.
			// If we wanted to set the timeout for each individual packet read,
			// we could call conn.SetDeadline inside the for-loop immediately before each conn.ReadFrom call.
			if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
				return nil, errors.New("setting deadline timeout to listen for syn-ack response: " + err.Error())
			}
		default:
		}

		numBts, readRemoteAddr, err := conn.ReadFrom(buf)
		readTime := time.Now()
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				// after we time out, we need to add everything still in remotes into responses as an error
				for remoteAddr, portMap := range remotes {
					for port, _ := range portMap {
						responses = append(responses, sarListenerResp{
							Addr: remoteAddr,
							Port: port,
							Err:  err,
						})
					}
				}
				return responses, nil
			}

			// wasn't a read deadline error, some other kind of error
			return nil, errors.New("reading response: " + err.Error())
		}
		if numBts < 14 { // 14 because the last data we need to look at, the flags, are at index 13
			// log.Warnf("receiveSAR got malformed packet, too short")
			continue
		}

		tcpHdrRaw := TCPHdr(buf[:numBts])

		if int(tcpHdrRaw.DestPort()) != srcPort {
			// log.Warnf("receiveSAR got mismatched dest port %+v expecting %v, ignoring and continuing to listen\n", tcpHdrRaw.DestPort(), srcPort)
			continue
		}

		remoteAddr := readRemoteAddr.String()
		addrMap, ok := remotes[remoteAddr]
		if !ok {
			// this is normal. Listening to an IP makes us receive all packets for all apps using that local IP.
			// They're duplicated to us, so we're not causing anyone else to lose packets.
			// We just need to ignore them
			continue
		}

		// note we swap src and dest ports. The source for the packet we sent is the dest we recieve, and vice versa.

		destPort := int(tcpHdrRaw.SrcPort()) // source of packet is dest for our request
		ackNum, ok := addrMap[destPort]
		if !ok {
			// log.Warnf("receiveSAR got mismatched source port %+v expecting %v, ignoring and continuing to listen\n", tcpHdrRaw.SrcPort(), dstPort)
			continue
		}

		if tcpHdrRaw.AckNum() != ackNum {
			// log.Warnf("receiveSAR got mismatched ack num %+v expecting %v, ignoring and continuing to listen\n", tcpHdrRaw.AckNum(), ackNum)
			continue
		}

		if tcpHdrRaw.SYN() && tcpHdrRaw.ACK() {
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
		} else if tcpHdrRaw.RST() {
			// We're expecting either a Reset or a Syn-Ack
		} else {
			log.Warnln("receiveSAR got packet that wasn't an rst or a synack, but matched the ports and AckNum. Very strange!")
			continue // packet wasn't a Rest or a Syn-Ack, ignore it
		}

		// We found a match for an addr+port+acknum in our requests, that was either a Reset or Syn-Ack response.

		// So now add the time to the response array, remove the request from remoteMap,
		// and return if remoteMap is empty

		responses = append(responses, sarListenerResp{
			Addr:     remoteAddr,
			Port:     destPort,
			RespTime: readTime,
		})

		delete(addrMap, destPort)
		if len(addrMap) == 0 {
			delete(remotes, remoteAddr)
		}

		// log.Warnf("multisar got packet, still have %v remotes\n", len(remotes))

		// if we read all the packets we were expecting, return
		if len(remotes) == 0 {
			return responses, nil
		}
	}
}
