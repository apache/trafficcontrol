// Copyright 2015-present Basho Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package riak

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	backoff "github.com/basho/backoff"
	proto "github.com/golang/protobuf/proto"
)

// Connection errors
var (
	ErrCannotRead  = errors.New("Cannot read from a non-active or closed connection")
	ErrCannotWrite = errors.New("Cannot write to a non-active or closed connection")
)

// AuthOptions object contains the authentication credentials and tls config
type AuthOptions struct {
	User      string
	Password  string
	TlsConfig *tls.Config
}

type connectionOptions struct {
	remoteAddress       *net.TCPAddr
	connectTimeout      time.Duration
	requestTimeout      time.Duration
	authOptions         *AuthOptions
	tempNetErrorRetries uint16
}

const (
	connCreated state = iota
	connTlsStarting
	connActive
	connInactive
)

type connection struct {
	addr                *net.TCPAddr
	conn                net.Conn
	connectTimeout      time.Duration
	requestTimeout      time.Duration
	tempNetErrorRetries uint16
	authOptions         *AuthOptions
	sizeBuf             []byte
	dataBuf             []byte
	active              bool
	inFlight            bool
	lastUsed            time.Time
	stateData
}

func newConnection(options *connectionOptions) (*connection, error) {
	if options == nil {
		return nil, ErrOptionsRequired
	}
	if options.remoteAddress == nil {
		return nil, ErrAddressRequired
	}
	if options.connectTimeout == 0 {
		options.connectTimeout = defaultConnectTimeout
	}
	if options.requestTimeout == 0 {
		options.requestTimeout = defaultRequestTimeout
	}
	if options.tempNetErrorRetries == 0 {
		options.tempNetErrorRetries = defaultTempNetErrorRetries
	}
	c := &connection{
		addr:                options.remoteAddress,
		connectTimeout:      options.connectTimeout,
		requestTimeout:      options.requestTimeout,
		tempNetErrorRetries: options.tempNetErrorRetries,
		authOptions:         options.authOptions,
		sizeBuf:             make([]byte, 4),
		dataBuf:             make([]byte, defaultInitBuffer),
		inFlight:            false,
		lastUsed:            time.Now(),
	}
	c.initStateData("connCreated", "connTlsStarting", "connActive", "connInactive")
	c.setState(connCreated)
	return c, nil
}

func (c *connection) connect() (err error) {
	dialer := &net.Dialer{
		Timeout:   c.connectTimeout,
		KeepAlive: time.Second * 30,
	}
	c.conn, err = dialer.Dial("tcp", c.addr.String()) // NB: SetNoDelay() is true by default for TCP connections
	if err != nil {
		logError("[Connection]", "error when dialing %s: '%s'", c.addr.String(), err.Error())
		c.close()
	} else {
		logDebug("[Connection]", "connected to: %s", c.addr)
		if err = c.startTls(); err != nil {
			c.close()
			c.setState(connInactive)
			return
		}
		c.setState(connActive)
	}
	return
}

func (c *connection) startTls() error {
	if c.authOptions == nil {
		return nil
	}
	if c.authOptions.TlsConfig == nil {
		return ErrAuthMissingConfig
	}
	c.setState(connTlsStarting)
	startTlsCmd := &startTlsCommand{}
	if err := c.execute(startTlsCmd); err != nil {
		return err
	}
	var tlsConn *tls.Conn
	if tlsConn = tls.Client(c.conn, c.authOptions.TlsConfig); tlsConn == nil {
		return ErrAuthTLSUpgradeFailed
	}
	if err := tlsConn.Handshake(); err != nil {
		return err
	}
	c.conn = tlsConn
	authCmd := &authCommand{
		user:     c.authOptions.User,
		password: c.authOptions.Password,
	}
	return c.execute(authCmd)
}

func (c *connection) available() bool {
	return (c.conn != nil && c.isStateLessThan(connInactive))
}

func (c *connection) close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

func (c *connection) setInFlight(inFlightVal bool) {
	c.inFlight = inFlightVal
}

func (c *connection) execute(cmd Command) (err error) {
	if c.inFlight == true {
		err = fmt.Errorf("[Connection] attempted to run '%s' command on in-use connection", cmd.Name())
		return
	}

	if lc, ok := cmd.(listingCommand); ok {
		allowListing := lc.getAllowListing()
		if !allowListing {
			err = ErrListingDisabled
			cmd.onError(err)
			return
		}
	}

	c.setInFlight(true)
	defer c.setInFlight(false)
	c.lastUsed = time.Now()

	var message []byte
	message, err = getRiakMessage(cmd)
	if err != nil {
		return
	}

	// Use the *greater* of the connection's request timeout
	// or the Command's timeout
	timeout := c.requestTimeout
	if tc, ok := cmd.(timeoutCommand); ok {
		tc := tc.getTimeout()
		if tc > timeout {
			timeout = tc
		}
	}

	if err = c.write(message, timeout); err != nil {
		return
	}

	var response []byte
	var decoded proto.Message
	for {
		response, err = c.read(timeout) // NB: response *will* have entire pb message
		if err != nil {
			cmd.onError(err)
			return
		}

		// Maybe translate RpbErrorResp into golang error
		if err = maybeRiakError(response); err != nil {
			cmd.onError(err)
			return
		}

		if decoded, err = decodeRiakMessage(cmd, response); err != nil {
			cmd.onError(err)
			return
		}

		err = cmd.onSuccess(decoded)
		if err != nil {
			cmd.onError(err)
			return
		}

		if sc, ok := cmd.(streamingCommand); ok {
			// Streaming Commands indicate done
			if sc.isDone() {
				return
			}
		} else {
			// non-streaming command, done at this point
			return
		}
	}
}

func (c *connection) setReadDeadline(t time.Duration) {
	c.conn.SetReadDeadline(time.Now().Add(t))
}

// NB: This will read one full pb message from Riak, or error in doing so
func (c *connection) read(timeout time.Duration) ([]byte, error) {
	if !c.available() {
		return nil, ErrCannotRead
	}

	var err error
	var count int
	var messageLength uint32
	var rt time.Duration = timeout // rt = 'read timeout'
	b := &backoff.Backoff{
		Min:    rt,
		Jitter: true,
	}
	try := uint16(0)

	for {
		c.setReadDeadline(rt)
		if count, err = io.ReadFull(c.conn, c.sizeBuf); err == nil && count == 4 {
			messageLength = binary.BigEndian.Uint32(c.sizeBuf)
			if messageLength > uint32(cap(c.dataBuf)) {
				logDebug("[Connection]", "allocating larger dataBuf of size %d", messageLength)
				c.dataBuf = make([]byte, messageLength)
			} else {
				c.dataBuf = c.dataBuf[0:messageLength]
			}
			// FUTURE: large object warning / error
			// TODO: FUTURE this deadline should subtract the duration taken by the first
			// ReadFull call. Currently it's could wait up to 2X the read timout value
			c.setReadDeadline(rt)
			count, err = io.ReadFull(c.conn, c.dataBuf)
		} else {
			if err == nil && count != 4 {
				err = newClientError(fmt.Sprintf("[Connection] expected to read 4 bytes, only read: %d", count), nil)
			}
		}

		if err == nil && count != int(messageLength) {
			err = newClientError(fmt.Sprintf("[Connection] message length: %d, only read: %d", messageLength, count), nil)
		}

		if err == nil {
			return c.dataBuf, nil
		}

		if try < c.tempNetErrorRetries && isTemporaryNetError(err) {
			rt = b.Duration()
			try++
			logDebug("[Connection]", "temporary error, re-try %v, new read timeout: %v", try, rt)
		} else {
			c.setState(connInactive)
			return nil, err
		}
	}
}

func (c *connection) write(data []byte, timeout time.Duration) error {
	if !c.available() {
		return ErrCannotWrite
	}
	c.conn.SetWriteDeadline(time.Now().Add(timeout))
	count, err := c.conn.Write(data)
	if err != nil {
		c.setState(connInactive)
		return err
	}
	if count != len(data) {
		return newClientError(fmt.Sprintf("[Connection] data length: %d, only wrote: %d", len(data), count), nil)
	}
	return nil
}
