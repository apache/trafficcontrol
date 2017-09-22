package web

import (
	"crypto/tls"
	"errors"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"net"
	"net/http"
	"time"
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
				log.Errorf("ConnState callback: idle conn is not a InterceptConn: '%T'\n", conn)
			} else {
				// MUST be zeroed when the conn moves to Idle, because the Active callback happens _after_ some/all bytes have been read
				iconn.bytesRead = 0
				iconn.bytesWritten = 0
			}
			connMap.Pop(conn.RemoteAddr().String())
		case http.StateActive:
			if iconn, ok := conn.(*InterceptConn); !ok {
				log.Errorf("ConnState callback: active conn is not a InterceptConn: '%T'\n", conn)
			} else {
				connMap.Push(iconn)
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

// InterceptListenTLS is like InterceptListen but for serving HTTPS.
func InterceptListenTLS(net string, laddr string, certs []tls.Certificate) (net.Listener, *ConnMap, func(net.Conn, http.ConnState), error) {
	interceptListener, connMap, connState, err := InterceptListen(net, laddr)
	if err != nil {
		return nil, nil, nil, errors.New("creating InterceptListen: " + err.Error())
	}
	config := &tls.Config{NextProtos: []string{"h2"}}
	config.Certificates = certs
	tlsListener := tls.NewListener(interceptListener, config)
	return tlsListener, connMap, connState, nil
}

func (l *InterceptListener) Accept() (net.Conn, error) {
	c, err := l.realListener.Accept()
	if err != nil {
		log.Errorf("Accept err: %v\n", err) // TODO stats?
		return c, err
	}
	interceptConn := &InterceptConn{realConn: c}
	l.connMap.Push(interceptConn)
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
