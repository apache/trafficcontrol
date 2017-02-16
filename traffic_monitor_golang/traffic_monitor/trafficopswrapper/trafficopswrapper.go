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
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/crconfig"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// ITrafficOpsSession provides an interface to the Traffic Ops client, so it may be wrapped or mocked.
type ITrafficOpsSession interface {
	CRConfigRaw(cdn string) ([]byte, error)
	LastCRConfig(cdn string) ([]byte, time.Time, error)
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

type ByteTime struct {
	bytes []byte
	time  time.Time
}

type ByteMapCache struct {
	cache *map[string]ByteTime
	m     *sync.RWMutex
}

func NewByteMapCache() ByteMapCache {
	return ByteMapCache{m: &sync.RWMutex{}, cache: &map[string]ByteTime{}}
}

func (c ByteMapCache) Set(key string, newBytes []byte) {
	c.m.Lock()
	defer c.m.Unlock()
	(*c.cache)[key] = ByteTime{bytes: newBytes, time: time.Now()}
}

func (c ByteMapCache) Get(key string) ([]byte, time.Time) {
	c.m.RLock()
	defer c.m.RUnlock()
	if byteTime, ok := (*c.cache)[key]; !ok {
		return nil, time.Time{}
	} else {
		return byteTime.bytes, byteTime.time
	}
}

// TrafficOpsSessionThreadsafe provides access to the Traffic Ops client safe for multiple goroutines. This fulfills the ITrafficOpsSession interface.
type TrafficOpsSessionThreadsafe struct {
	session      **to.Session // pointer-to-pointer, because we're given a pointer from the Traffic Ops package, and we don't want to copy it.
	m            *sync.Mutex
	lastCRConfig ByteMapCache
}

// NewTrafficOpsSessionThreadsafe returns a new threadsafe TrafficOpsSessionThreadsafe wrapping the given `Session`.
func NewTrafficOpsSessionThreadsafe(s *to.Session) TrafficOpsSessionThreadsafe {
	return TrafficOpsSessionThreadsafe{session: &s, m: &sync.Mutex{}, lastCRConfig: NewByteMapCache()}
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
	b, _, err := ss.GetCRConfig(cdn)
	if err == nil {
		s.lastCRConfig.Set(cdn, b)
	}
	return b, err
}

// LastCRConfig returns the last CRConfig requested from CRConfigRaw, and the time it was returned. This is designed to be used in conjunction with a poller which regularly calls CRConfigRaw. If no last CRConfig exists, because CRConfigRaw has never been called successfully, this calls CRConfigRaw once to try to get the CRConfig from Traffic Ops.
func (s TrafficOpsSessionThreadsafe) LastCRConfig(cdn string) ([]byte, time.Time, error) {
	crConfig, crConfigTime := s.lastCRConfig.Get(cdn)
	if crConfig == nil {
		b, err := s.CRConfigRaw(cdn)
		return b, time.Now(), err
	}
	return crConfig, crConfigTime, nil
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
		s := to.TrafficServer{}
		if srv.Profile != nil {
			s.Profile = *srv.Profile
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing Profile field\n", name)
		}
		if srv.Ip != nil {
			s.IP = *srv.Ip
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing IP field\n", name)
		}
		if srv.Status != nil {
			s.Status = string(*srv.Status)
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing Status field\n", name)
		}
		if srv.CacheGroup != nil {
			s.CacheGroup = *srv.CacheGroup
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing CacheGroup field\n", name)
		}
		if srv.Ip6 != nil {
			s.IP6 = *srv.Ip6
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing IP6 field\n", name)
		}
		if srv.Port != nil {
			s.Port = *srv.Port
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing Port field\n", name)
		}
		s.HostName = name
		if srv.Fqdn != nil {
			s.FQDN = *srv.Fqdn
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing FQDN field\n", name)
		}
		if srv.InterfaceName != nil {
			s.InterfaceName = *srv.InterfaceName
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing InterfaceName field\n", name)
		}
		if srv.ServerType != nil {
			s.Type = *srv.ServerType
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing Type field\n", name)
		}
		if srv.HashId != nil {
			s.HashID = *srv.HashId
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing HashId field\n", name)
		}
		mc.TrafficServer[name] = s
	}

	// Dump the "live" monitoring.json monitors, and populate with the "snapshotted" CRConfig
	mc.TrafficMonitor = map[string]to.TrafficMonitor{}
	for name, mon := range crConfig.Monitors {
		// monitorProfile = *mon.Profile
		m := to.TrafficMonitor{}
		if mon.Port != nil {
			m.Port = *mon.Port
		} else {
			log.Warnf("Creating monitor config: CRConfig monitor %s missing Port field\n", name)
		}
		if mon.IP6 != nil {
			m.IP6 = *mon.IP6
		} else {
			log.Warnf("Creating monitor config: CRConfig monitor %s missing IP6 field\n", name)
		}
		if mon.IP != nil {
			m.IP = *mon.IP
		} else {
			log.Warnf("Creating monitor config: CRConfig monitor %s missing IP field\n", name)
		}
		m.HostName = name
		if mon.FQDN != nil {
			m.FQDN = *mon.FQDN
		} else {
			log.Warnf("Creating monitor config: CRConfig monitor %s missing FQDN field\n", name)
		}
		if mon.Profile != nil {
			m.Profile = *mon.Profile
		} else {
			log.Warnf("Creating monitor config: CRConfig monitor %s missing Profile field\n", name)
		}
		if mon.Location != nil {
			m.Location = *mon.Location
		} else {
			log.Warnf("Creating monitor config: CRConfig monitor %s missing Location field\n", name)
		}
		if mon.Status != nil {
			m.Status = string(*mon.Status)
		} else {
			log.Warnf("Creating monitor config: CRConfig monitor %s missing Status field\n", name)
		}
		mc.TrafficMonitor[name] = m
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
