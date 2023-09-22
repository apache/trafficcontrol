// Package towrap wraps two versions of Traffic Ops clients to give up-to-date
// information, possibly using legacy API versions.
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
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	legacyClient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/publicsuffix"
)

const localHostIP = "127.0.0.1"

// ErrNilSession is the error returned by operations performed on a nil session.
var ErrNilSession = errors.New("nil session")

// ByteTime is a structure for associating a set of raw data with some CDN
// Snapshot statistics, and a certain time.
type ByteTime struct {
	bytes []byte
	time  time.Time
	stats *tc.CRConfigStats
}

// ByteMapCache is a thread-access-safe map of cache server hostnames to
// ByteTime structures.
type ByteMapCache struct {
	cache *map[string]ByteTime
	m     *sync.RWMutex
}

// NewByteMapCache constructs a new, empty ByteMapCache.
func NewByteMapCache() ByteMapCache {
	return ByteMapCache{m: &sync.RWMutex{}, cache: &map[string]ByteTime{}}
}

// Set sets the entry given by 'key' to a new ByteTime structure with the given
// raw data ('newBytes') and the given statistics ('stats') at the current time.
func (c ByteMapCache) Set(key string, newBytes []byte, stats *tc.CRConfigStats) {
	c.m.Lock()
	defer c.m.Unlock()
	(*c.cache)[key] = ByteTime{bytes: newBytes, stats: stats, time: time.Now()}
}

// Get retrieves the raw data, associated time, and statistics of the entry
// given by 'key'.
func (c ByteMapCache) Get(key string) ([]byte, time.Time, *tc.CRConfigStats) {
	c.m.RLock()
	defer c.m.RUnlock()
	if byteTime, ok := (*c.cache)[key]; !ok {
		return nil, time.Time{}, nil
	} else {
		return byteTime.bytes, byteTime.time, byteTime.stats
	}
}

func (s TrafficOpsSessionThreadsafe) BackupFileExists() bool {
	if _, err := os.Stat(s.CRConfigBackupFile); !os.IsNotExist(err) {
		if _, err = os.Stat(s.TMConfigBackupFile); !os.IsNotExist(err) {
			return true
		}
	}
	return false
}

// CRConfigStat represents a set of statistics from a CDN Snapshot requested at
// a particular time.
type CRConfigStat struct {
	// Err contains any error that may have occurred when obtaining the
	// statistics.
	Err error `json:"error"`
	// ReqAddr is the network address from which the statistics were requested.
	ReqAddr string `json:"request_address"`
	// ReqTime is the time at which the request for statistics was made.
	ReqTime time.Time `json:"request_time"`
	// Stats contains the actual statistics.
	Stats tc.CRConfigStats `json:"stats"`
}

// CopyCRConfigStat makes a deep copy of a slice of CRConfigStats.
func CopyCRConfigStat(old []CRConfigStat) []CRConfigStat {
	newStats := make([]CRConfigStat, len(old))
	copy(newStats, old)
	return newStats
}

// CRConfigHistoryThreadsafe stores history in a circular buffer.
type CRConfigHistoryThreadsafe struct {
	hist   *[]CRConfigStat
	m      *sync.RWMutex
	limit  *uint64
	length *uint64
	pos    *uint64
}

// NewCRConfigHistoryThreadsafe constructs a new, empty
// CRConfigHistoryThreadsafe - this is the ONLY way to safely create a
// CRConfigHistoryThreadsafe, using the zero value of the structure will cause
// all operations to encounter segmentation faults, and there is no way to
// preempt this.
//
// 'limit' indicates the size of the circular buffer - effectively the number of
// entries it will be capable of storing.
func NewCRConfigHistoryThreadsafe(limit uint64) CRConfigHistoryThreadsafe {
	hist := make([]CRConfigStat, limit, limit)
	length := uint64(0)
	pos := uint64(0)
	return CRConfigHistoryThreadsafe{hist: &hist, m: &sync.RWMutex{}, limit: &limit, length: &length, pos: &pos}
}

// Add adds the given stat to the history. Does not add new additions with the
// same remote address and CRConfig Date as the previous.
func (h CRConfigHistoryThreadsafe) Add(i *CRConfigStat) {
	h.m.Lock()
	defer h.m.Unlock()

	if *h.length != 0 {
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
	if *h.length < *h.limit {
		*h.length++
	}
}

// Get retrieves the stored history of CRConfigStat entries.
func (h CRConfigHistoryThreadsafe) Get() []CRConfigStat {
	h.m.RLock()
	defer h.m.RUnlock()
	if *h.length < *h.limit {
		return CopyCRConfigStat((*h.hist)[:*h.length])
	}
	newStats := make([]CRConfigStat, *h.limit)
	copy(newStats, (*h.hist)[*h.pos:])
	copy(newStats[*h.length-*h.pos:], (*h.hist)[:*h.pos])
	return newStats
}

// Len gives the number of currently stored items in the buffer.
//
// An uninitialized buffer has zero length.
func (h CRConfigHistoryThreadsafe) Len() uint64 {
	if h.length == nil {
		return 0
	}
	return *h.length
}

// TrafficOpsSessionThreadsafe provides access to the Traffic Ops client safe
// for multiple goroutines. This fulfills the ITrafficOpsSession interface.
type TrafficOpsSessionThreadsafe struct {
	session            **client.Session // pointer-to-pointer, because we're given a pointer from the Traffic Ops package, and we don't want to copy it.
	legacySession      **legacyClient.Session
	m                  *sync.Mutex
	lastCRConfig       ByteMapCache
	crConfigHist       CRConfigHistoryThreadsafe
	CRConfigBackupFile string
	TMConfigBackupFile string
}

// NewTrafficOpsSessionThreadsafe returns a new threadsafe
// TrafficOpsSessionThreadsafe wrapping the given `Session`.
func NewTrafficOpsSessionThreadsafe(s *client.Session, ls *legacyClient.Session, histLimit uint64, cfg config.Config) TrafficOpsSessionThreadsafe {
	return TrafficOpsSessionThreadsafe{
		CRConfigBackupFile: cfg.CRConfigBackupFile,
		crConfigHist:       NewCRConfigHistoryThreadsafe(histLimit),
		lastCRConfig:       NewByteMapCache(),
		m:                  &sync.Mutex{},
		session:            &s,
		legacySession:      &ls,
		TMConfigBackupFile: cfg.TMConfigBackupFile,
	}
}

// Initialized tells whether or not the TrafficOpsSessionThreadsafe has been
// properly initialized with non-nil sessions.
func (s TrafficOpsSessionThreadsafe) Initialized() bool {
	return s.session != nil && *s.session != nil && s.legacySession != nil && *s.legacySession != nil
}

// Update updates the TrafficOpsSessionThreadsafe's connection information with
// the provided information. It's safe for calling by multiple goroutines, being
// aware that they will race.
func (s *TrafficOpsSessionThreadsafe) Update(
	url string,
	username string,
	password string,
	insecure bool,
	userAgent string,
	useCache bool,
	timeout time.Duration,
) error {
	if s == nil {
		return errors.New("cannot update nil session")
	}
	s.m.Lock()
	defer s.m.Unlock()

	// always set unauthenticated sessions first which can eventually authenticate themselves when attempting requests
	if err := s.setSession(url, username, password, insecure, userAgent, useCache, timeout); err != nil {
		return err
	}
	if err := s.setLegacySession(url, username, password, insecure, userAgent, useCache, timeout); err != nil {
		return err
	}

	session, _, err := client.LoginWithAgent(url, username, password, insecure, userAgent, useCache, timeout)
	if err != nil {
		log.Errorf("logging in using up-to-date client: %v", err)
		legacySession, _, err := legacyClient.LoginWithAgent(url, username, password, insecure, userAgent, useCache, timeout)
		if err != nil || legacySession == nil {
			err = fmt.Errorf("logging in using legacy client: %v", err)
			return err
		}
		*s.legacySession = legacySession
	} else {
		*s.session = session
	}

	return nil
}

// setSession sets the session for the up-to-date client without logging in.
func (s *TrafficOpsSessionThreadsafe) setSession(url, username, password string, insecure bool, userAgent string, useCache bool, timeout time.Duration) error {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return err
	}
	to := client.NewSession(username, password, url, userAgent, &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
		Jar: jar,
	}, useCache)
	*s.session = to
	return nil
}

// setSession sets the session for the legacy client without logging in.
func (s *TrafficOpsSessionThreadsafe) setLegacySession(url, username, password string, insecure bool, userAgent string, useCache bool, timeout time.Duration) error {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return err
	}
	to := legacyClient.NewSession(username, password, url, userAgent, &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
		Jar: jar,
	}, useCache)
	*s.legacySession = to
	return nil
}

// getThreadsafeSession is used internally to get a copy of the session pointer,
// or nil if it doesn't exist. This should not be used outside
// TrafficOpsSessionThreadsafe, and never stored, because part of the purpose of
// rafficOpsSessionThreadsafe is to store a pointer to the Session pointer, so
// it can be updated by one goroutine and immediately used by another. This
// should only be called immediately before using the session, since someone
// else may update it concurrently.
func (s TrafficOpsSessionThreadsafe) get() *client.Session {
	s.m.Lock()
	defer s.m.Unlock()
	if s.session == nil || *s.session == nil {
		return nil
	}
	return *s.session
}

func (s TrafficOpsSessionThreadsafe) getLegacy() *legacyClient.Session {
	s.m.Lock()
	defer s.m.Unlock()
	if s.legacySession == nil || *s.legacySession == nil {
		return nil
	}
	return *s.legacySession
}

// CRConfigHistory gets all of the stored, historical data about CRConfig
// Snapshots' Stats sections.
func (s TrafficOpsSessionThreadsafe) CRConfigHistory() []CRConfigStat {
	return s.crConfigHist.Get()
}

// CRConfigValid checks if the passed tc.CRConfig structure is valid, and
// ensures that it is from the same CDN as the last CRConfig Snapshot, as well
// as that it is newer than the last CRConfig Snapshot.
func (s *TrafficOpsSessionThreadsafe) CRConfigValid(crc *tc.CRConfig, cdn string) error {
	if crc == nil {
		return errors.New("CRConfig is nil")
	}
	if crc.Stats.CDNName == nil {
		return errors.New("CRConfig.Stats.CDN missing")
	}
	if crc.Stats.DateUnixSeconds == nil {
		return errors.New("CRConfig.Stats.Date missing")
	}

	// Note this intentionally takes intended CDN, rather than trusting
	// crc.Stats
	lastCrc, lastCrcTime, lastCrcStats := s.lastCRConfig.Get(cdn)

	if lastCrc == nil {
		return nil
	}
	if lastCrcStats.DateUnixSeconds == nil {
		log.Warnln("TrafficOpsSessionThreadsafe.CRConfigValid returning no error, but last CRConfig Date was missing!")
		return nil
	}
	if lastCrcStats.CDNName == nil {
		log.Warnln("TrafficOpsSessionThreadsafe.CRConfigValid returning no error, but last CRConfig CDN was missing!")
		return nil
	}
	if *lastCrcStats.CDNName != *crc.Stats.CDNName {
		return errors.New("CRConfig.Stats.CDN " + *crc.Stats.CDNName + " different than last received CRConfig.Stats.CDNName " + *lastCrcStats.CDNName + " received at " + lastCrcTime.Format(time.RFC3339Nano))
	}
	if *lastCrcStats.DateUnixSeconds > *crc.Stats.DateUnixSeconds {
		return errors.New("CRConfig.Stats.Date " + strconv.FormatInt(*crc.Stats.DateUnixSeconds, 10) + " older than last received CRConfig.Stats.Date " + strconv.FormatInt(*lastCrcStats.DateUnixSeconds, 10) + " received at " + lastCrcTime.Format(time.RFC3339Nano))
	}
	return nil
}

// CRConfigRaw returns the CRConfig from the Traffic Ops. This is safe for
// multiple goroutines.
func (s TrafficOpsSessionThreadsafe) CRConfigRaw(cdn string) ([]byte, error) {

	var remoteAddr string
	var err error
	var crConfig *tc.CRConfig
	var configBytes []byte
	json := jsoniter.ConfigFastest

	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}
	response, reqInf, err := ss.GetCRConfig(cdn, client.RequestOptions{})
	if reqInf.RemoteAddr != nil {
		remoteAddr = reqInf.RemoteAddr.String()
	}
	if err != nil {
		log.Warnln("getting CRConfig from Traffic Ops using up-to-date client: " + err.Error() + ". Retrying with legacy client")
		ls := s.getLegacy()
		if ls == nil {
			return nil, ErrNilSession
		}
		response, reqInf, err := ls.GetCRConfig(cdn, legacyClient.RequestOptions{})
		if reqInf.RemoteAddr != nil {
			remoteAddr = reqInf.RemoteAddr.String()
		}
		if err != nil {
			log.Errorln("getting CRConfig from Traffic Ops using legacy client: " + err.Error() + ". Checking for backup")
		}
		configBytes, err = json.Marshal(response.Response)
	} else {
		crConfig = &response.Response
		configBytes, err = json.Marshal(crConfig)
	}
	if err != nil {
		crConfig = nil
		log.Warnln("failed to marshal CRConfig using up-to-date client: " + err.Error())
	}

	if err == nil {
		log.Infoln("successfully got CRConfig from Traffic Ops. Writing to backup file")
		if wErr := ioutil.WriteFile(s.CRConfigBackupFile, configBytes, 0644); wErr != nil {
			log.Errorf("failed to write CRConfig backup file: %v", wErr)
		}
	} else {
		if s.BackupFileExists() {
			log.Errorln("using backup file for CRConfig snapshot due to error fetching CRConfig snapshot from Traffic Ops: " + err.Error())
			configBytes, err = ioutil.ReadFile(s.CRConfigBackupFile)
			if err != nil {
				return nil, fmt.Errorf("reading CRConfig backup file: %v", err)
			}
			remoteAddr = localHostIP
			err = nil
		} else {
			return nil, fmt.Errorf("failed to get CRConfig from Traffic Ops (%v), and there is no backup file", err)
		}
	}

	hist := &CRConfigStat{
		Err:     err,
		ReqAddr: remoteAddr,
		ReqTime: time.Now(),
		Stats:   tc.CRConfigStats{},
	}
	defer s.crConfigHist.Add(hist)

	if crConfig == nil {
		if err = json.Unmarshal(configBytes, crConfig); err != nil {
			err = errors.New("invalid JSON: " + err.Error())
			hist.Err = err
			return configBytes, err
		}
	}
	hist.Stats = crConfig.Stats

	if err = s.CRConfigValid(crConfig, cdn); err != nil {
		err = errors.New("invalid CRConfig: " + err.Error())
		hist.Err = err
		return configBytes, err
	}

	s.lastCRConfig.Set(cdn, configBytes, &crConfig.Stats)
	return configBytes, nil
}

// LastCRConfig returns the last CRConfig requested from CRConfigRaw, and the
// time it was returned. This is designed to be used in conjunction with a
// poller which regularly calls CRConfigRaw. If no last CRConfig exists, because
// CRConfigRaw has never been called successfully, this calls CRConfigRaw once
// to try to get the CRConfig from Traffic Ops.
func (s TrafficOpsSessionThreadsafe) LastCRConfig(cdn string) ([]byte, time.Time, error) {
	crConfig, crConfigTime, _ := s.lastCRConfig.Get(cdn)
	if len(crConfig) == 0 {
		b, err := s.CRConfigRaw(cdn)
		return b, time.Now(), err
	}
	return crConfig, crConfigTime, nil
}

func (s TrafficOpsSessionThreadsafe) fetchTMConfig(cdn string) (*tc.TrafficMonitorConfig, error) {
	ss := s.get()
	if ss == nil {
		return nil, ErrNilSession
	}

	m, _, e := ss.GetTrafficMonitorConfig(cdn, client.NewRequestOptions())
	return &m.Response, e
}

func (s TrafficOpsSessionThreadsafe) fetchLegacyTMConfig(cdn string) (*tc.TrafficMonitorConfig, error) {
	ss := s.getLegacy()
	var m tc.TrafficMonitorConfig
	if ss == nil {
		return nil, ErrNilSession
	}

	r, _, e := ss.GetTrafficMonitorConfig(cdn, legacyClient.RequestOptions{})
	if e != nil {
		return nil, e
	}
	m = r.Response
	return &m, e
}

// trafficMonitorConfigMapRaw returns the Traffic Monitor config map from the
// Traffic Ops, directly from the monitoring endpoint. This is not usually
// what is needed, rather monitoring needs the snapshotted CRConfig data, which
// is filled in by `LegacyTrafficMonitorConfigMap`. This is safe for multiple
// goroutines.
func (s TrafficOpsSessionThreadsafe) trafficMonitorConfigMapRaw(cdn string) (*tc.TrafficMonitorConfigMap, error) {
	var config *tc.TrafficMonitorConfig
	var configMap *tc.TrafficMonitorConfigMap
	var err error

	config, err = s.fetchTMConfig(cdn)
	if err != nil {
		log.Warnln("getting Traffic Monitor config from Traffic Ops using up-to-date client: " + err.Error() + ". Retrying with legacy client")
		config, err = s.fetchLegacyTMConfig(cdn)
		if err != nil {
			log.Errorln("getting Traffic Monitor config from Traffic Ops using legacy client: " + err.Error())
		}
	}

	if err == nil {
		log.Infoln("successfully got Traffic Monitor config from Traffic Ops")
		if config == nil {
			return nil, fmt.Errorf("nil Traffic Monitor config after successful fetch")
		}
		configMap, err = tc.TrafficMonitorTransformToMap(config)
	}

	if err != nil {
		// Default error case, no backup file exists
		if !s.BackupFileExists() {
			return nil, err
		}
		log.Errorln("using backup file for monitoring config snapshot due to invalid monitoring config snapshot from Traffic Ops: " + err.Error())

		b, err := ioutil.ReadFile(s.TMConfigBackupFile)
		if err != nil {
			return nil, errors.New("reading TMConfigBackupFile: " + err.Error())
		}

		json := jsoniter.ConfigFastest
		var tmConfig tc.TrafficMonitorConfig
		if err := json.Unmarshal(b, &tmConfig); err != nil {
			return nil, errors.New("unmarshalling backup file monitoring.json: " + err.Error())
		}
		return tc.TrafficMonitorTransformToMap(&tmConfig)
	}

	json := jsoniter.ConfigFastest
	data, err := json.Marshal(*config)
	if err == nil {
		if wErr := ioutil.WriteFile(s.TMConfigBackupFile, data, 0644); wErr != nil {
			log.Errorf("failed to write TM config backup file: %v", wErr)
		}
	}

	return configMap, err
}

// TrafficMonitorConfigMap returns the Traffic Monitor config map from the
// Traffic Ops. This is safe for multiple goroutines.
func (s TrafficOpsSessionThreadsafe) TrafficMonitorConfigMap(cdn string) (*tc.TrafficMonitorConfigMap, error) {
	mc, err := s.trafficMonitorConfigMapRaw(cdn)
	if err != nil {
		return nil, fmt.Errorf("getting monitor config map: %v", err)
	}
	return mc, nil
}

func (s TrafficOpsSessionThreadsafe) fetchServerByHostname(hostName string) (tc.ServerV50, error) {
	ss := s.get()
	if ss == nil {
		return tc.ServerV50{}, ErrNilSession
	}

	params := url.Values{}
	params.Set("hostName", hostName)
	resp, _, err := ss.GetServers(client.RequestOptions{QueryParameters: params})
	if err != nil {
		return tc.ServerV50{}, fmt.Errorf("fetching server by hostname '%s': %v", hostName, err)
	}

	respLen := len(resp.Response)
	if respLen < 1 {
		return tc.ServerV50{}, fmt.Errorf("no server '%s' found in Traffic Ops", hostName)
	}

	var server tc.ServerV50
	var num int
	found := false
	for i, srv := range resp.Response {
		num = i
		if srv.CDNID > -1 && srv.HostName == hostName {
			server = srv
			found = true
			break
		}
	}
	if !found {
		return tc.ServerV50{}, fmt.Errorf("either no server '%s' found in Traffic Ops, or none by that hostName had non-nil CDN", hostName)
	}

	if respLen > 1 {
		log.Warnf("Getting monitor server by hostname '%s' returned %d servers - selecting #%d", hostName, respLen, num)
	}

	return server, nil
}

func (s TrafficOpsSessionThreadsafe) fetchLegacyServerByHostname(hostName string) (tc.ServerV50, error) {
	ss := s.getLegacy()
	if ss == nil {
		return tc.ServerV50{}, ErrNilSession
	}

	params := url.Values{}
	params.Set("hostName", hostName)
	resp, _, err := ss.GetServers(legacyClient.RequestOptions{QueryParameters: params})
	if err != nil {
		return tc.ServerV50{}, fmt.Errorf("fetching server by hostname '%s': %v", hostName, err)
	}

	respLen := len(resp.Response)
	if respLen < 1 {
		return tc.ServerV50{}, fmt.Errorf("no server '%s' found in Traffic Ops", hostName)
	}

	var server tc.ServerV40
	var num int
	found := false
	for i, srv := range resp.Response {
		num = i
		if srv.CDNName != nil && srv.HostName != nil && *srv.HostName == hostName {
			server = srv
			found = true
			break
		}
	}
	if !found {
		return tc.ServerV50{}, fmt.Errorf("either no server '%s' found in Traffic Ops, or none by that hostName had non-nil CDN", hostName)
	}

	if respLen > 1 {
		log.Warnf("Getting monitor server by hostname '%s' returned %d servers - selecting #%d", hostName, respLen, num)
	}

	if len(server.ProfileNames) == 0 {
		return tc.ServerV50{}, fmt.Errorf("server with hostname '%s' has no profile", hostName)
	}
	newServer := server.Upgrade()
	if err != nil {
		return newServer, fmt.Errorf("coercing legacy server to new format: %v", err)
	}
	return newServer, nil
}

// MonitorCDN returns the name of the CDN of a Traffic Monitor with the given
// hostName.
func (s TrafficOpsSessionThreadsafe) MonitorCDN(hostName string) (string, error) {
	var server tc.ServerV50
	var err error

	server, err = s.fetchServerByHostname(hostName)
	if err != nil {
		log.Warnln("getting server by hostname '" + hostName + "' using up-to-date client: " + err.Error() + ". Retrying with legacy client")
		server, err = s.fetchLegacyServerByHostname(hostName)
	}

	if err != nil {
		return "", fmt.Errorf("getting monitor CDN: %v", err)
	}

	// nil-dereference checks done already in each 'fetch' method; they'll just
	// return an error in that case
	return server.CDN, nil
}
