package web

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

type InterceptListener struct {
	realListener net.Listener
	connMap      *ConnMap
}

func getConnStateCallback(connMap *ConnMap) func(net.Conn, http.ConnState) {
	return func(conn net.Conn, state http.ConnState) {
		switch state {
		case http.StateClosed:
			fallthrough
		case http.StateIdle:
			if iconn, ok := conn.(*InterceptConn); !ok {
				log.Infof("ConnState callback: idle conn is not a InterceptConn: '%T'\n", conn)
			} else {
				// MUST be zeroed when the conn moves to Idle, because the Active callback happens _after_ some/all bytes have been read
				iconn.bytesRead = 0
				iconn.bytesWritten = 0
			}
			connMap.Remove(conn.RemoteAddr().String())
		case http.StateActive:
			if iconn, ok := conn.(*InterceptConn); !ok {
				log.Infof("ConnState callback: active conn is not a InterceptConn: '%T'\n", conn)
			} else {
				connMap.Add(iconn)
			}
		}
	}
}

// InterceptListen creates and returns a net.Listener via net.Listen, which is wrapped with an intercepter, which counts Conn read and write bytes. If you want a `grove.NewCacheHandler` to be able to count in and out bytes per remap rule in the stats interface, it must be served with a listener created via InterceptListen or InterceptListenTLS.
func InterceptListen(network, laddr string) (net.Listener, *ConnMap, func(net.Conn, http.ConnState), error) {
	l, err := net.Listen(network, laddr)
	if err != nil {
		return l, nil, nil, err
	}
	connMap := NewConnMap()
	return &InterceptListener{realListener: l, connMap: connMap}, connMap, getConnStateCallback(connMap), nil
}

// InterceptListenTLS is like InterceptListen but for serving HTTPS. It returns the tls.Config, which must be set on the http.Server using this listener for HTTP/2 to be set up.
func InterceptListenTLS(network string, laddr string, certs []tls.Certificate, h2Disabled bool) (net.Listener, *ConnMap, func(net.Conn, http.ConnState), *tls.Config, error) {
	config := &tls.Config{}
	// HTTP2 is enabled if config.DisableHTTP2 is false
	if !h2Disabled {
		config.NextProtos = []string{"h2"}
	}
	config.Certificates = certs
	config.BuildNameToCertificate()
	l, err := net.Listen(network, laddr)
	if err != nil {
		return l, nil, nil, nil, err
	}
	connMap := NewConnMap()

	interceptListener := &InterceptListener{realListener: l, connMap: connMap}
	tlsListener := tls.NewListener(interceptListener, config)
	return tlsListener, connMap, getConnStateCallback(connMap), config, nil
}

func (l *InterceptListener) Accept() (net.Conn, error) {
	c, err := l.realListener.Accept()
	if err != nil {
		log.Errorf("Accept err: %v\n", err) // TODO stats?
		return c, err
	}
	interceptConn := &InterceptConn{realConn: c}
	l.connMap.Add(interceptConn)
	return interceptConn, nil
}

func (l *InterceptListener) Close() error {
	return l.realListener.Close()
}

func (l *InterceptListener) Addr() net.Addr {
	return l.realListener.Addr()
}

type InterceptConn struct {
	realConn     net.Conn
	bytesRead    int
	bytesWritten int
}

func (c *InterceptConn) BytesRead() int {
	return c.bytesRead
}

func (c *InterceptConn) BytesWritten() int {
	return c.bytesWritten
}

func (c *InterceptConn) Read(b []byte) (n int, err error) {
	n, err = c.realConn.Read(b)
	c.bytesRead += n
	return
}
func (c *InterceptConn) Write(b []byte) (n int, err error) {
	n, err = c.realConn.Write(b)
	c.bytesWritten += n
	return
}
func (c *InterceptConn) Close() error {
	return c.realConn.Close()
}
func (c *InterceptConn) LocalAddr() net.Addr {
	return c.realConn.LocalAddr()
}
func (c *InterceptConn) RemoteAddr() net.Addr {
	return c.realConn.RemoteAddr()
}
func (c *InterceptConn) SetDeadline(t time.Time) error {
	return c.realConn.SetDeadline(t)
}
func (c *InterceptConn) SetReadDeadline(t time.Time) error {
	return c.realConn.SetReadDeadline(t)
}
func (c *InterceptConn) SetWriteDeadline(t time.Time) error {
	return c.realConn.SetWriteDeadline(t)
}

func (c *InterceptConn) Real() net.Conn {
	return c.realConn
}

// SetDSCP attempts to set the DSCP flag on the TCP connection of the InterceptConn (TOS in IPv4, Traffic Class in IPv6). Returns any error. Note this is not guaranteed to succeed on all operating systems. In particular, these functions are known to fail on OS X, and may or may not succeed on various Linux kernels and distributions. Callers should expect failure, and it is recommended the connection and request continue to be handled.
func (c *InterceptConn) SetDSCP(dscp int) error {
	if dscp == 0 {
		return nil
	}
	if c == nil {
		return errors.New("Conn is nil")
	}
	realConn := c.Real()
	if realConn == nil {
		return errors.New("real Conn is nil")
	}
	ipv4Err := ipv4.NewConn(realConn).SetTOS(dscp)
	ipv6Err := ipv6.NewConn(realConn).SetTrafficClass(dscp)
	if ipv4Err != nil || ipv6Err != nil {
		errStr := ""
		if ipv4Err != nil {
			errStr = "setting IPv4 TOS: " + ipv4Err.Error()
		}
		if ipv6Err != nil {
			if ipv4Err != nil {
				errStr += "; "
			}
			errStr += "setting IPv6 TrafficClass: " + ipv6Err.Error()
		}
		return errors.New(errStr)
	}
	return nil
}
