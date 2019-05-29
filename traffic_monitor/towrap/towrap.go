package towrap

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
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/client"
)

// ITrafficOpsSession provides an interface to the Traffic Ops client, so it may be wrapped or mocked.
type ITrafficOpsSession interface {
	CRConfigRaw(cdn string) ([]byte, error)
	LastCRConfig(cdn string) ([]byte, time.Time, error)
	TrafficMonitorConfigMap(cdn string) (*tc.TrafficMonitorConfigMap, error)
	Set(session *client.Session)
	URL() (string, error)
	User() (string, error)
	Servers() ([]tc.Server, error)
	Profiles() ([]tc.Profile, error)
	Parameters(profileName string) ([]tc.Parameter, error)
	DeliveryServices() ([]tc.DeliveryService, error)
	CacheGroups() ([]tc.CacheGroupNullable, error)
	CRConfigHistory() []CRConfigStat
}

var ErrNilSession = fmt.Errorf("nil session")

// TODO rename CRConfigCacheObj
type ByteTime struct {
	bytes []byte
	time  time.Time
	stats *tc.CRConfigStats
}

type ByteMapCache struct {
	cache *map[string]ByteTime
	m     *sync.RWMutex
}

func NewByteMapCache() ByteMapCache {
	return ByteMapCache{m: &sync.RWMutex{}, cache: &map[string]ByteTime{}}
}

func (c ByteMapCache) Set(key string, newBytes []byte, stats *tc.CRConfigStats) {
	c.m.Lock()
	defer c.m.Unlock()
	(*c.cache)[key] = ByteTime{bytes: newBytes, stats: stats, time: time.Now()}
}

func (c ByteMapCache) Get(key string) ([]byte, time.Time, *tc.CRConfigStats) {
	c.m.RLock()
	defer c.m.RUnlock()
	if byteTime, ok := (*c.cache)[key]; !ok {
		return nil, time.Time{}, nil
	} else {
		return byteTime.bytes, byteTime.time, byteTime.stats
	}
}

// CRConfigHistoryThreadsafe stores history in a circular buffer.
type CRConfigHistoryThreadsafe struct {
	hist  *[]CRConfigStat
	m     *sync.RWMutex
	limit *uint64
	len   *uint64
	pos   *uint64
}

func NewCRConfigHistoryThreadsafe(limit uint64) CRConfigHistoryThreadsafe {
	hist := make([]CRConfigStat, limit, limit)
	len := uint64(0)
	pos := uint64(0)
	return CRConfigHistoryThreadsafe{hist: &hist, m: &sync.RWMutex{}, limit: &limit, len: &len, pos: &pos}
}

// Add adds the given stat to the history. Does not add new additions with the same remote address and CRConfig Date as the previous.
func (h CRConfigHistoryThreadsafe) Add(i *CRConfigStat) {
	h.m.Lock()
	defer h.m.Unlock()

	if *h.len != 0 {
		last := (*h.hist)[(*h.pos-1)%*h.limit]
		datesEqual := (i.Stats.DateUnixSeconds == nil && last.Stats.DateUnixSeconds == nil) || (i.Stats.DateUnixSeconds != nil && last.Stats.DateUnixSeconds != nil && *i.Stats.DateUnixSeconds == *last.Stats.DateUnixSeconds)
		cdnsEqual := (i.Stats.CDNName == nil && last.Stats.CDNName == nil) || (i.Stats.CDNName != nil && last.Stats.CDNName != nil && *i.Stats.CDNName == *last.Stats.CDNName)
		reqAddrsEqual := i.ReqAddr == last.ReqAddr
		if reqAddrsEqual && datesEqual && cdnsEqual {
			return
		}
	}

	(*h.hist)[*h.pos] = *i
	*h.pos = (*h.pos + 1) % *h.limit
	if *h.len < *h.limit {
		*h.len++
	}
}

func (h CRConfigHistoryThreadsafe) Get() []CRConfigStat {
	h.m.RLock()
	defer h.m.RUnlock()
	if *h.len < *h.limit {
		return CopyCRConfigStat((*h.hist)[:*h.len])
	}
	new := make([]CRConfigStat, *h.limit)
	copy(new, (*h.hist)[*h.pos:])
	copy(new[*h.len-*h.pos:], (*h.hist)[:*h.pos])
	return new
}

func CopyCRConfigStat(old []CRConfigStat) []CRConfigStat {
	new := make([]CRConfigStat, len(old))
	copy(new, old)
	return new
}

type CRConfigStat struct {
	ReqTime time.Time        `json:"request_time"`
	ReqAddr string           `json:"request_address"`
	Stats   tc.CRConfigStats `json:"stats"`
	Err     error            `json:"error"`
}

// TrafficOpsSessionThreadsafe provides access to the Traffic Ops client safe for multiple goroutines. This fulfills the ITrafficOpsSession interface.
type TrafficOpsSessionThreadsafe struct {
	session      **client.Session // pointer-to-pointer, because we're given a pointer from the Traffic Ops package, and we don't want to copy it.
	m            *sync.Mutex
	lastCRConfig ByteMapCache
	crConfigHist CRConfigHistoryThreadsafe
}

// NewTrafficOpsSessionThreadsafe returns a new threadsafe TrafficOpsSessionThreadsafe wrapping the given `Session`.
func NewTrafficOpsSessionThreadsafe(s *client.Session, crConfigHistoryLimit uint64) TrafficOpsSessionThreadsafe {
	return TrafficOpsSessionThreadsafe{session: &s, m: &sync.Mutex{}, lastCRConfig: NewByteMapCache(), crConfigHist: NewCRConfigHistoryThreadsafe(crConfigHistoryLimit)}
}

// Set sets the internal Traffic Ops session. This is safe for multiple goroutines, being aware they will race.
func (s TrafficOpsSessionThreadsafe) Set(session *client.Session) {
	s.m.Lock()
	defer s.m.Unlock()
	*s.session = session
}

// getThreadsafeSession is used internally to get a copy of the session pointer, or nil if it doesn't exist. This should not be used outside TrafficOpsSessionThreadsafe, and never stored, because part of the purpose of TrafficOpsSessionThreadsafe is to store a pointer to the Session pointer, so it can be updated by one goroutine and immediately used by another. This should only be called immediately before using the session, since someone else may update it concurrently.
func (s TrafficOpsSessionThreadsafe) get() *client.Session {
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

func (s TrafficOpsSessionThreadsafe) CRConfigHistory() []CRConfigStat {
	return s.crConfigHist.Get()
}

func (s *TrafficOpsSessionThreadsafe) CRConfigValid(crc *tc.CRConfig, cdn string) error {
	// Note this intentionally takes intended CDN, rather than trusting crc.Stats
	lastCrc, lastCrcTime, lastCrcStats := s.lastCRConfig.Get(cdn)
	if lastCrc == nil {
		return nil
	}
	if lastCrcStats.DateUnixSeconds == nil {
		log.Warnln("TrafficOpsSessionThreadsafe.CRConfigValid returning no error, but last CRConfig Date was missing!")
		return nil
	}
	if *lastCrcStats.CDNName != *crc.Stats.CDNName {
		return errors.New("CRConfig.Stats.CDN " + *crc.Stats.CDNName + " different than last received CRConfig.Stats.CDNName " + *lastCrcStats.CDNName + " received at " + lastCrcTime.Format(time.RFC3339Nano))
	}
	if crc.Stats.DateUnixSeconds == nil {
		return errors.New("CRConfig.Stats.Date missing")
	}
	if *lastCrcStats.DateUnixSeconds > *crc.Stats.DateUnixSeconds {
		return errors.New("CRConfig.Stats.Date " + strconv.FormatInt(*crc.Stats.DateUnixSeconds, 10) + " older than last received CRConfig.Stats.Date " + strconv.FormatInt(*lastCrcStats.DateUnixSeconds, 10) + " received at " + lastCrcTime.Format(time.RFC3339Nano))
	}
	return nil
}

// CRConfigRaw returns the CRConfig from the Traffic Ops. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) CRConfigRaw(cdn string) ([]byte, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	b, reqInf, err := ss.GetCRConfig(cdn)

	hist := &CRConfigStat{time.Now(), reqInf.RemoteAddr.String(), tc.CRConfigStats{}, err}
	defer s.crConfigHist.Add(hist)

	if err != nil {
		return b, err
	}

	crc := &tc.CRConfig{}
	if err = json.Unmarshal(b, crc); err != nil {
		err = errors.New("invalid JSON: " + err.Error())
		hist.Err = err
		return b, err
	}
	hist.Stats = crc.Stats

	if err = s.CRConfigValid(crc, cdn); err != nil {
		err = errors.New("invalid CRConfig: " + err.Error())
		hist.Err = err
		return b, err
	}

	s.lastCRConfig.Set(cdn, b, &crc.Stats)
	return b, nil
}

// LastCRConfig returns the last CRConfig requested from CRConfigRaw, and the time it was returned. This is designed to be used in conjunction with a poller which regularly calls CRConfigRaw. If no last CRConfig exists, because CRConfigRaw has never been called successfully, this calls CRConfigRaw once to try to get the CRConfig from Traffic Ops.
func (s TrafficOpsSessionThreadsafe) LastCRConfig(cdn string) ([]byte, time.Time, error) {
	crConfig, crConfigTime, _ := s.lastCRConfig.Get(cdn)
	if crConfig == nil {
		b, err := s.CRConfigRaw(cdn)
		return b, time.Now(), err
	}
	return crConfig, crConfigTime, nil
}

// TrafficMonitorConfigMapRaw returns the Traffic Monitor config map from the Traffic Ops, directly from the monitoring.json endpoint. This is not usually what is needed, rather monitoring needs the snapshotted CRConfig data, which is filled in by `TrafficMonitorConfigMap`. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) trafficMonitorConfigMapRaw(cdn string) (*tc.TrafficMonitorConfigMap, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	configMap, _, error := ss.GetTrafficMonitorConfigMap(cdn)
	return configMap, error
}

// TrafficMonitorConfigMap returns the Traffic Monitor config map from the Traffic Ops. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) TrafficMonitorConfigMap(cdn string) (*tc.TrafficMonitorConfigMap, error) {
	mc, err := s.trafficMonitorConfigMapRaw(cdn)
	if err != nil {
		return nil, fmt.Errorf("getting monitor config map: %v", err)
	}

	crcData, err := s.CRConfigRaw(cdn)
	if err != nil {
		return nil, fmt.Errorf("getting CRConfig: %v", err)
	}

	crConfig := tc.CRConfig{}
	if err := json.Unmarshal(crcData, &crConfig); err != nil {
		return nil, fmt.Errorf("unmarshalling CRConfig JSON : %v", err)
	}

	mc, err = CreateMonitorConfig(crConfig, mc)
	if err != nil {
		return nil, fmt.Errorf("creating Traffic Monitor Config: %v", err)
	}

	return mc, nil
}

func CreateMonitorConfig(crConfig tc.CRConfig, mc *tc.TrafficMonitorConfigMap) (*tc.TrafficMonitorConfigMap, error) {
	// Dump the "live" monitoring.json servers, and populate with the "snapshotted" CRConfig
	mc.TrafficServer = map[string]tc.TrafficServer{}
	for name, srv := range crConfig.ContentServers {
		s := tc.TrafficServer{}
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
		if srv.ServerStatus != nil {
			s.ServerStatus = string(*srv.ServerStatus)
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
		if srv.HttpsPort != nil {
			s.HTTPSPort = *srv.HttpsPort
		} else {
			log.Warnf("Creating monitor config: CRConfig server %s missing HttpsPort field\n", name)
		}
		mc.TrafficServer[name] = s
	}

	// Dump the "live" monitoring.json monitors, and populate with the "snapshotted" CRConfig
	mc.TrafficMonitor = map[string]tc.TrafficMonitor{}
	for name, mon := range crConfig.Monitors {
		// monitorProfile = *mon.Profile
		m := tc.TrafficMonitor{}
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
		if mon.ServerStatus != nil {
			m.ServerStatus = string(*mon.ServerStatus)
		} else {
			log.Warnf("Creating monitor config: CRConfig monitor %s missing ServerStatus field\n", name)
		}
		mc.TrafficMonitor[name] = m
	}

	// Dump the "live" monitoring.json DeliveryServices, and populate with the "snapshotted" CRConfig
	// But keep using the monitoring.json thresholds, because they're not in the CRConfig.
	rawDeliveryServices := mc.DeliveryService
	mc.DeliveryService = map[string]tc.TMDeliveryService{}
	for name, _ := range crConfig.DeliveryServices {
		if rawDS, ok := rawDeliveryServices[name]; ok {
			// use the raw DS if it exists, because the CRConfig doesn't have thresholds or statuses
			mc.DeliveryService[name] = rawDS
		} else {
			mc.DeliveryService[name] = tc.TMDeliveryService{
				XMLID:              name,
				TotalTPSThreshold:  0,
				ServerStatus:       "REPORTED",
				TotalKbpsThreshold: 0,
			}
		}
	}
	return mc, nil
}

func (s TrafficOpsSessionThreadsafe) Servers() ([]tc.Server, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	servers, _, error := ss.GetServers()
	return servers, error
}

func (s TrafficOpsSessionThreadsafe) Profiles() ([]tc.Profile, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	profiles, _, error := ss.GetProfiles()
	return profiles, error
}

func (s TrafficOpsSessionThreadsafe) Parameters(profileName string) ([]tc.Parameter, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	parameters, _, error := ss.GetParametersByProfileName(profileName)
	return parameters, error
}

func (s TrafficOpsSessionThreadsafe) DeliveryServices() ([]tc.DeliveryService, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	deliveryServices, _, error := ss.GetDeliveryServices()
	return deliveryServices, error
}

func (s TrafficOpsSessionThreadsafe) CacheGroups() ([]tc.CacheGroupNullable, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	cacheGroups, _, error := ss.GetCacheGroupsNullable()
	return cacheGroups, error
}
