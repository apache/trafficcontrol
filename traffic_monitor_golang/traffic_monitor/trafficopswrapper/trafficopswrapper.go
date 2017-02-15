package trafficopswrapper

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
	"encoding/json"
	"fmt"
	"sync"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/crconfig"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// ITrafficOpsSession provides an interface to the Traffic Ops client, so it may be wrapped or mocked.
type ITrafficOpsSession interface {
	CRConfigRaw(cdn string) ([]byte, error)
	TrafficMonitorConfigMap(cdn string) (*to.TrafficMonitorConfigMap, error)
	Set(session *to.Session)
	URL() (string, error)
	User() (string, error)
	Servers() ([]to.Server, error)
	Profiles() ([]to.Profile, error)
	Parameters(profileName string) ([]to.Parameter, error)
	DeliveryServices() ([]to.DeliveryService, error)
	CacheGroups() ([]to.CacheGroup, error)
}

var ErrNilSession = fmt.Errorf("nil session")

// TrafficOpsSessionThreadsafe provides access to the Traffic Ops client safe for multiple goroutines. This fulfills the ITrafficOpsSession interface.
type TrafficOpsSessionThreadsafe struct {
	session **to.Session // pointer-to-pointer, because we're given a pointer from the Traffic Ops package, and we don't want to copy it.
	m       *sync.Mutex
}

// NewTrafficOpsSessionThreadsafe returns a new threadsafe TrafficOpsSessionThreadsafe wrapping the given `Session`.
func NewTrafficOpsSessionThreadsafe(s *to.Session) TrafficOpsSessionThreadsafe {
	return TrafficOpsSessionThreadsafe{&s, &sync.Mutex{}}
}

// Set sets the internal Traffic Ops session. This is safe for multiple goroutines, being aware they will race.
func (s TrafficOpsSessionThreadsafe) Set(session *to.Session) {
	s.m.Lock()
	defer s.m.Unlock()
	*s.session = session
}

// getThreadsafeSession is used internally to get a copy of the session pointer, or nil if it doesn't exist. This should not be used outside TrafficOpsSessionThreadsafe, and never stored, because part of the purpose of TrafficOpsSessionThreadsafe is to store a pointer to the Session pointer, so it can be updated by one goroutine and immediately used by another. This should only be called immediately before using the session, since someone else may update it concurrently.
func (s TrafficOpsSessionThreadsafe) get() *to.Session {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil
	}
	return *s.session
}

func (s TrafficOpsSessionThreadsafe) URL() (string, error) {
	ss := s.get()
	if ss == nil {
		return "", ErrNilSession
	}
	return ss.URL, nil
}

func (s TrafficOpsSessionThreadsafe) User() (string, error) {
	ss := s.get()
	if ss == nil {
		return "", ErrNilSession
	}
	return ss.UserName, nil
}

// CRConfigRaw returns the CRConfig from the Traffic Ops. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) CRConfigRaw(cdn string) ([]byte, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	b, _, e := ss.GetCRConfig(cdn)
	return b, e
}

// TrafficMonitorConfigMapRaw returns the Traffic Monitor config map from the Traffic Ops, directly from the monitoring.json endpoint. This is not usually what is needed, rather monitoring needs the snapshotted CRConfig data, which is filled in by `TrafficMonitorConfigMap`. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) trafficMonitorConfigMapRaw(cdn string) (*to.TrafficMonitorConfigMap, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	return ss.TrafficMonitorConfigMap(cdn)
}

// TrafficMonitorConfigMap returns the Traffic Monitor config map from the Traffic Ops. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) TrafficMonitorConfigMap(cdn string) (*to.TrafficMonitorConfigMap, error) {
	mc, err := s.trafficMonitorConfigMapRaw(cdn)
	if err != nil {
		return nil, fmt.Errorf("getting monitor config map: %v", err)
	}

	crcData, err := s.CRConfigRaw(cdn)
	if err != nil {
		return nil, fmt.Errorf("getting CRConfig: %v", err)
	}

	crConfig := crconfig.CRConfig{}
	if err := json.Unmarshal(crcData, &crConfig); err != nil {
		return nil, fmt.Errorf("unmarshalling CRConfig JSON: %v", err)
	}

	mc, err = CreateMonitorConfig(crConfig, mc)
	if err != nil {
		return nil, fmt.Errorf("creating Traffic Monitor Config: %v", err)
	}

	return mc, nil
}

func CreateMonitorConfig(crConfig crconfig.CRConfig, mc *to.TrafficMonitorConfigMap) (*to.TrafficMonitorConfigMap, error) {
	// Dump the "live" monitoring.json servers, and populate with the "snapshotted" CRConfig
	mc.TrafficServer = map[string]to.TrafficServer{}
	for name, srv := range crConfig.ContentServers {
		mc.TrafficServer[name] = to.TrafficServer{
			Profile:       *srv.Profile,
			IP:            *srv.Ip,
			Status:        string(*srv.Status),
			CacheGroup:    *srv.CacheGroup,
			IP6:           *srv.Ip6,
			Port:          *srv.Port,
			HostName:      name,
			FQDN:          *srv.Fqdn,
			InterfaceName: *srv.InterfaceName,
			Type:          *srv.ServerType,
			HashID:        *srv.HashId,
		}
	}

	// Dump the "live" monitoring.json monitors, and populate with the "snapshotted" CRConfig
	mc.TrafficMonitor = map[string]to.TrafficMonitor{}
	for name, mon := range crConfig.Monitors {
		// monitorProfile = *mon.Profile
		mc.TrafficMonitor[name] = to.TrafficMonitor{
			Port:     *mon.Port,
			IP6:      *mon.IP6,
			IP:       *mon.IP,
			HostName: name,
			FQDN:     *mon.FQDN,
			Profile:  *mon.Profile,
			Location: *mon.Location,
			Status:   string(*mon.Status),
		}
	}

	// Dump the "live" monitoring.json DeliveryServices, and populate with the "snapshotted" CRConfig
	// But keep using the monitoring.json thresholds, because they're not in the CRConfig.
	rawDeliveryServices := mc.DeliveryService
	mc.DeliveryService = map[string]to.TMDeliveryService{}
	for name, _ := range crConfig.DeliveryServices {
		if rawDS, ok := rawDeliveryServices[name]; ok {
			// use the raw DS if it exists, because the CRConfig doesn't have thresholds or statuses
			mc.DeliveryService[name] = rawDS
		} else {
			mc.DeliveryService[name] = to.TMDeliveryService{
				XMLID:              name,
				TotalTPSThreshold:  0,
				Status:             "REPORTED",
				TotalKbpsThreshold: 0,
			}
		}
	}
	return mc, nil
}

func (s TrafficOpsSessionThreadsafe) Servers() ([]to.Server, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	return ss.Servers()
}

func (s TrafficOpsSessionThreadsafe) Profiles() ([]to.Profile, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	return ss.Profiles()
}

func (s TrafficOpsSessionThreadsafe) Parameters(profileName string) ([]to.Parameter, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	return ss.Parameters(profileName)
}

func (s TrafficOpsSessionThreadsafe) DeliveryServices() ([]to.DeliveryService, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	return ss.DeliveryServices()
}

func (s TrafficOpsSessionThreadsafe) CacheGroups() ([]to.CacheGroup, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	return ss.CacheGroups()
}
