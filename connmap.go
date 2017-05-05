package grove

import (
	"fmt"
	"net"
	"sync"
)

type ConnMap struct {
	conns map[string]net.Conn
	m     sync.Mutex
}

func NewConnMap() *ConnMap {
	return &ConnMap{conns: map[string]net.Conn{}}
}

func (cm ConnMap) Pop(remoteAddr string) (net.Conn, bool) {
	fmt.Printf("ConnMap popping '%v'\n", remoteAddr)
	cm.m.Lock()
	defer cm.m.Unlock()
	c, ok := cm.conns[remoteAddr]
	delete(cm.conns, remoteAddr)
	return c, ok
}

func (cm ConnMap) Push(conn net.Conn) {
	fmt.Printf("ConnMap pushing '%v'\n", conn.RemoteAddr().String())
	cm.m.Lock()
	defer cm.m.Unlock()
	cm.conns[conn.RemoteAddr().String()] = conn
}
