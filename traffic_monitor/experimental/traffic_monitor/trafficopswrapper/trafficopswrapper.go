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
	"fmt"
	"sync"

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
	Parameters(profileName string) ([]to.Parameter, error)
	DeliveryServices() ([]to.DeliveryService, error)
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

// TrafficMonitorConfigMap returns the Traffic Monitor config map from the Traffic Ops. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) TrafficMonitorConfigMap(cdn string) (*to.TrafficMonitorConfigMap, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil, ErrNilSession
	}
	d, e := (*s.session).TrafficMonitorConfigMap(cdn)
	return d, e
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
