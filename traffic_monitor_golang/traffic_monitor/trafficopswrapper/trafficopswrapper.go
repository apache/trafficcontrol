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

func (s TrafficOpsSessionThreadsafe) URL() (string, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return "", ErrNilSession
	}
	url := (*s.session).URL
	return url, nil
}

func (s TrafficOpsSessionThreadsafe) User() (string, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return "", ErrNilSession
	}
	user := (*s.session).UserName
	return user, nil
}

// TrafficOpsSessionThreadsafe provides access to the Traffic Ops client safe for multiple goroutines. This fulfills the ITrafficOpsSession interface.
type TrafficOpsSessionThreadsafe struct {
	session **to.Session // pointer-to-pointer, because we're given a pointer from the Traffic Ops package, and we don't want to copy it.
	m       *sync.Mutex
}

// NewTrafficOpsSessionThreadsafe returns a new threadsafe TrafficOpsSessionThreadsafe wrapping the given `Session`.
func NewTrafficOpsSessionThreadsafe(s *to.Session) TrafficOpsSessionThreadsafe {
	return TrafficOpsSessionThreadsafe{&s, &sync.Mutex{}}
}

// CRConfigRaw returns the CRConfig from the Traffic Ops. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) CRConfigRaw(cdn string) ([]byte, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil, ErrNilSession
	}
	b, _, e := (*s.session).GetCRConfig(cdn)
	return b, e
}

// TrafficMonitorConfigMapRaw returns the Traffic Monitor config map from the Traffic Ops, directly from the monitoring.json endpoint. This is not usually what is needed, rather monitoring needs the snapshotted CRConfig data, which is filled in by `TrafficMonitorConfigMap`. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) TrafficMonitorConfigMapRaw(cdn string) (*to.TrafficMonitorConfigMap, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil, ErrNilSession
	}
	d, e := (*s.session).TrafficMonitorConfigMap(cdn)
	return d, e
}

// TrafficMonitorConfigMap returns the Traffic Monitor config map from the Traffic Ops. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) TrafficMonitorConfigMap(cdn string) (*to.TrafficMonitorConfigMap, error) {
	mc, err := s.TrafficMonitorConfigMapRaw(cdn)
	if err != nil {
		fmt.Printf("DEBUG4 TrafficMonitorConfigMap err: %v\n", err)
		return nil, fmt.Errorf("getting monitor config map: %v", err)
	}

	crcData, err := s.CRConfigRaw(cdn)
	if err != nil {
		fmt.Printf("DEBUG4 CRConfigRaw err: %v\n", err)
		return nil, fmt.Errorf("getting CRConfig: %v", err)
	}

	crConfig := crconfig.CRConfig{}
	if err := json.Unmarshal(crcData, &crConfig); err != nil {
		fmt.Printf("DEBUG4 CRConfig Unmarshal err: %v\n", err)
		return nil, fmt.Errorf("Error unmarshalling CRConfig JSON: %v", err)
		return nil, err
	}

	mc, err = CreateMonitorConfig(crConfig, mc)
	if err != nil {
		fmt.Printf("DEBUG4 CreateMonitorConfig err: %v\n", err)
		return nil, fmt.Errorf("Error creating Traffic Monitor Config: %v", err)
	}

	// mcMap, err := to.TrafficMonitorTransformToMap(mc)
	// if err != nil {
	// 	fmt.Printf("DEBUG4 TrafficMonitorTransformToMap err: %v\n", err)
	// 	return nil, fmt.Errorf("Error transforming Traffic Monitor Config to Map: %v", err)
	// }

	// debug

	// if bytes, err := json.Marshal(mcMap); err != nil {
	// 	fmt.Printf("DEBUG4 error marshalling map: %v\n", err)
	// } else {
	// 	fmt.Printf("DEBUG4 New Map: %v\n\n", string(bytes))
	// }

	return mc, nil
}

func CreateMonitorConfig(crConfig crconfig.CRConfig, mc *to.TrafficMonitorConfigMap) (*to.TrafficMonitorConfigMap, error) {
	// mc := to.TrafficMonitorConfig{}

	// cgs, err := s.CacheGroups()
	// if err != nil {
	// 	return nil, fmt.Errorf("Error getting CacheGroups: %v", err)
	// }

	// allProfiles, err := s.Profiles()
	// if err != nil {
	// 	return nil, fmt.Errorf("Error getting Profiles: %v", err)
	// }

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

	// for _, cg := range cgs {
	// 	mc.CacheGroups = append(mc.CacheGroups, to.TMCacheGroup{
	// 		Name: cg.Name,
	// 		Coordinates: to.Coordinates{
	// 			Latitude:  cg.Latitude,
	// 			Longitude: cg.Longitude,
	// 		},
	// 	})
	// }

	// monitorProfile := ""

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

	// monitorParams, err := s.Parameters(monitorProfile)
	// if err != nil {
	// 	return nil, fmt.Errorf("Error getting profile %v parameters: %v", monitorProfile, err)
	// }
	// mc.Config = map[string]interface{}{}
	// for _, param := range monitorParams {
	// 	if numParam, err := strconv.ParseFloat(param.Value, 64); err == nil {
	// 		mc.Config[param.Name] = numParam
	// 	} else {
	// 		mc.Config[param.Name] = param.Value
	// 	}
	// }

	// Dump the "live" monitoring.json DeliveryServices, and populate with the "snapshotted" CRConfig
	mc.DeliveryService = map[string]to.TMDeliveryService{}
	for name, _ := range crConfig.DeliveryServices {
		mc.DeliveryService[name] = to.TMDeliveryService{
			XMLID:              name,
			TotalTPSThreshold:  0,          // TODO verify
			Status:             "Reported", // TODO verify
			TotalKbpsThreshold: 0,          // TODO verify
		}
	}

	// mc.Profiles = []to.TMProfile{}
	// for _, prof := range allProfiles {
	// 	if strings.HasPrefix(prof.Name, "EDGE") || strings.HasPrefix(prof.Name, "TEAK") {
	// 		mc.Profiles = append(mc.Profiles, to.TMProfile{
	// 			Name: prof.Name,
	// 			Type: "EDGE",
	// 		})
	// 	} else if strings.HasPrefix(prof.Name, "MID") {
	// 		mc.Profiles = append(mc.Profiles, to.TMProfile{
	// 			Name: prof.Name,
	// 			Type: "MID",
	// 		})
	// 	}
	// }
	// for profI, prof := range mc.Profiles {
	// 	// MID2_TOP_v5.3.2-757
	// 	prof.Parameters.Thresholds = map[string]to.HealthThreshold{}
	// 	// TODO lock
	// 	params, err := s.Parameters(prof.Name)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("Error getting profile %v parameters: %v", prof, err)
	// 	}
	// 	fmt.Printf("DEBUG5 profile %v len(params) %v\n", prof.Name, len(params))
	// 	for _, param := range params {
	// 		if param.Name == "health.connection.timeout" {
	// 			i, err := strconv.Atoi(param.Value)
	// 			if err != nil {
	// 				return nil, fmt.Errorf("Error getting profile %v parameter %v: %v", prof, param.Name, err)
	// 			}
	// 			prof.Parameters.HealthConnectionTimeout = i
	// 		} else if param.Name == "health.polling.url" {
	// 			prof.Parameters.HealthPollingURL = param.Value
	// 		} else if param.Name == "history.count" {
	// 			i, err := strconv.Atoi(param.Value)
	// 			if err != nil {
	// 				return nil, fmt.Errorf("Error getting profile %v parameter %v: %v", prof, param.Name, err)
	// 			}
	// 			prof.Parameters.HistoryCount = i
	// 		} else if strings.HasPrefix(param.Name, "health.threshold.") {
	// 			stat := param.Name[len("health.threshold."):]
	// 			thresh, err := to.StrToThreshold(param.Value)
	// 			if err != nil {
	// 				return nil, fmt.Errorf("Error getting profile %v parameter %v: %v", prof, param.Name, err)
	// 			}
	// 			prof.Parameters.Thresholds[stat] = thresh
	// 		}
	// 	}
	// 	mc.Profiles[profI] = prof
	// }
	return mc, nil
}

// Set sets the internal Traffic Ops session. This is safe for multiple goroutines, being aware they will race.
func (s TrafficOpsSessionThreadsafe) Set(session *to.Session) {
	s.m.Lock()
	defer s.m.Unlock()
	*s.session = session
}

func (s TrafficOpsSessionThreadsafe) Servers() ([]to.Server, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil, ErrNilSession
	}
	return (*s.session).Servers()
}

func (s TrafficOpsSessionThreadsafe) Profiles() ([]to.Profile, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil, ErrNilSession
	}
	return (*s.session).Profiles()
}

func (s TrafficOpsSessionThreadsafe) Parameters(profileName string) ([]to.Parameter, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil, ErrNilSession
	}
	return (*s.session).Parameters(profileName)
}

func (s TrafficOpsSessionThreadsafe) DeliveryServices() ([]to.DeliveryService, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil, ErrNilSession
	}
	return (*s.session).DeliveryServices()
}

func (s TrafficOpsSessionThreadsafe) CacheGroups() ([]to.CacheGroup, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil, ErrNilSession
	}
	return (*s.session).CacheGroups()
}
