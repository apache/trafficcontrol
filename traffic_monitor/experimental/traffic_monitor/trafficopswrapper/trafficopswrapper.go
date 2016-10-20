package trafficopswrapper

import (
	"fmt"
	"sync"

	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

type ITrafficOpsSession interface {
	CRConfigRaw(cdn string) ([]byte, error)
	TrafficMonitorConfigMap(cdn string) (*to.TrafficMonitorConfigMap, error)
	Set(session *to.Session)
}

type TrafficOpsSessionThreadsafe struct {
	session **to.Session // pointer-to-pointer, because we're given a pointer from the Traffic Ops package, and we don't want to copy it.
	m       *sync.Mutex
}

func NewTrafficOpsSessionThreadsafe(s *to.Session) TrafficOpsSessionThreadsafe {
	return TrafficOpsSessionThreadsafe{&s, &sync.Mutex{}}
}

func (s TrafficOpsSessionThreadsafe) CRConfigRaw(cdn string) ([]byte, error) {
	s.m.Lock()
	if s.session == nil || *s.session == nil {
		return nil, fmt.Errorf("nil session")
	}
	b, _, e := (*s.session).GetCRConfig(cdn)
	s.m.Unlock()
	return b, e
}

func (s TrafficOpsSessionThreadsafe) TrafficMonitorConfigMap(cdn string) (*to.TrafficMonitorConfigMap, error) {
	s.m.Lock()
	if s.session == nil || *s.session == nil {
		return nil, fmt.Errorf("nil session")
	}
	d, e := (*s.session).TrafficMonitorConfigMap(cdn)
	s.m.Unlock()
	return d, e
}

func (s TrafficOpsSessionThreadsafe) Set(session *to.Session) {
	s.m.Lock()
	*s.session = session
	s.m.Unlock()
}
