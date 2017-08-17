package grove

import (
	"net"
	"sync"
	// "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

type ConnMap struct {
	conns map[string]net.Conn
	m     sync.Mutex
}

func NewConnMap() *ConnMap {
	return &ConnMap{conns: map[string]net.Conn{}}
}

func (cm *ConnMap) Pop(remoteAddr string) (net.Conn, bool) {
	// log.Debugf("ConnMap popping '%v'\n", remoteAddr)
	cm.m.Lock()
	defer cm.m.Unlock()
	c, ok := cm.conns[remoteAddr]
	delete(cm.conns, remoteAddr)
	return c, ok
}

func (cm *ConnMap) Push(conn net.Conn) {
	// log.Debugf("ConnMap pushing '%v'\n", conn.RemoteAddr().String())
	cm.m.Lock()
	defer cm.m.Unlock()
	cm.conns[conn.RemoteAddr().String()] = conn
}
