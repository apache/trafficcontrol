package trafficopswrapper

import (
	traffic_ops "github.com/Comcast/traffic_control/traffic_ops/client"
	"sync"
)

type ITrafficOpsSession interface {
	CRConfigRaw(cdn string) ([]byte, error)
	TrafficMonitorConfigMap(cdn string) (*traffic_ops.TrafficMonitorConfigMap, error)
}

type TrafficOpsSessionThreadsafe struct {
	session *traffic_ops.Session
	m       *sync.Mutex
}

func NewTrafficOpsSessionThreadsafe(s *traffic_ops.Session) TrafficOpsSessionThreadsafe {
	return TrafficOpsSessionThreadsafe{s, &sync.Mutex{}}
}

func (s TrafficOpsSessionThreadsafe) CRConfigRaw(cdn string) ([]byte, error) {
	s.m.Lock()
	b, e := s.session.CRConfigRaw(cdn)
	s.m.Unlock()
	return b, e
}

func (s TrafficOpsSessionThreadsafe) TrafficMonitorConfigMap(cdn string) (*traffic_ops.TrafficMonitorConfigMap, error) {
	s.m.Lock()
	d, e := s.session.TrafficMonitorConfigMap(cdn)
	s.m.Unlock()
	return d, e
}
