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
	"strconv"
)

// EphemeralPortHolder serves 2 purposes: it gets an ephemeral TCP port, and it holds
// onto the port so the OS doesn't assign it to any other app.
//
// It listens on :0 thereby getting a socket on an ephemeral port.
// It continues listening but never reading from the socket until Close is called.
type EphemeralPortHolder struct {
	listener net.Listener
	port     int
}

// GetAndHoldEphemeralPort gets an ephemeral port, and listens on it to prevent
// the OS assigning the port to other apps.
// Close must be called on the returned EphemeralPortHolder to stop listening.
func GetAndHoldEphemeralPort(addr string) (*EphemeralPortHolder, error) {
	addr += ":0"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.New("listening: " + err.Error())
	}

	// get the port now, so EphemeralPortHolder.Port() doesn't need to return an error

	listenAddr := listener.Addr().String()
	ipPort := SplitLast(listenAddr, ":")
	if len(ipPort) < 1 {
		return nil, errors.New("malformed addr '" + listenAddr + "', should have been ip:port") // should never happen
	}
	portStr := ipPort[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("malformed addr '" + listenAddr + "' port was not an integer") // should never happen
	}
	if port > 65535 || port < 0 {
		return nil, errors.New("malformed addr '" + listenAddr + "' port was outside bounds") // should never happen
	}

	return &EphemeralPortHolder{listener: listener, port: port}, nil
}

func (ph *EphemeralPortHolder) Close() error {
	return ph.listener.Close()
}

func (ph *EphemeralPortHolder) Addr() net.Addr {
	return ph.listener.Addr()
}

func (ph *EphemeralPortHolder) Port() int {
	return ph.port
}
