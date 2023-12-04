package todata

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
	"regexp"
	"strings"
	"sync"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/towrap"

	jsoniter "github.com/json-iterator/go"
)

// Regexes maps Delivery Service Regular Expressions to delivery services.
// For performance, we categorize Regular Expressions into 3 categories:
// 1. Direct string matches, with no regular expression matching characters
// 2. .*\.foo\..* expressions, where foo is a direct string match with no regular expression matching characters
// 3. Everything else
// This allows us to do a cheap match on 1 and 2, and only regex match the uncommon case.
type Regexes struct {
	DirectMatches                      map[string]tc.DeliveryServiceName
	DotStartSlashDotFooSlashDotDotStar map[string]tc.DeliveryServiceName
	RegexMatch                         map[*regexp.Regexp]tc.DeliveryServiceName
}

// DeliveryService returns the delivery service which matches the given fqdn, or false.
func (d Regexes) DeliveryService(domain, subdomain, subsubdomain string) (tc.DeliveryServiceName, bool) {
	if ds, ok := d.DotStartSlashDotFooSlashDotDotStar[subdomain]; ok {
		return ds, true
	}
	fqdn := fmt.Sprintf("%s.%s.%s", subsubdomain, subdomain, domain)
	if ds, ok := d.DirectMatches[fqdn]; ok {
		return ds, true
	}
	for regex, ds := range d.RegexMatch {
		if regex.MatchString(fqdn) {
			return ds, true
		}
	}
	return "", false
}

// NewRegexes constructs a new Regexes object, initializing internal pointer members.
func NewRegexes() Regexes {
	return Regexes{
		DirectMatches:                      map[string]tc.DeliveryServiceName{},
		DotStartSlashDotFooSlashDotDotStar: map[string]tc.DeliveryServiceName{},
		RegexMatch:                         map[*regexp.Regexp]tc.DeliveryServiceName{},
	}
}

// TOData holds CDN data fetched from Traffic Ops.
type TOData struct {
	DeliveryServiceRegexes Regexes
	DeliveryServiceServers map[tc.DeliveryServiceName][]tc.CacheName
	DeliveryServiceTypes   map[tc.DeliveryServiceName]tc.DSTypeCategory
	ServerCachegroups      map[tc.CacheName]tc.CacheGroupName
	ServerDeliveryServices map[tc.CacheName][]tc.DeliveryServiceName
	ServerTypes            map[tc.CacheName]tc.CacheType
}

// New returns a new empty TOData object, initializing pointer members.
func New() *TOData {
	return &TOData{
		DeliveryServiceServers: map[tc.DeliveryServiceName][]tc.CacheName{},
		ServerDeliveryServices: map[tc.CacheName][]tc.DeliveryServiceName{},
		ServerTypes:            map[tc.CacheName]tc.CacheType{},
		DeliveryServiceTypes:   map[tc.DeliveryServiceName]tc.DSTypeCategory{},
		DeliveryServiceRegexes: NewRegexes(),
		ServerCachegroups:      map[tc.CacheName]tc.CacheGroupName{},
	}
}

// TODataThreadsafe provides safe access for multiple goroutine writers and one goroutine reader, to the encapsulated TOData object.
// This could be made lock-free, if the performance was necessary
type TODataThreadsafe struct {
	toData *TOData
	m      *sync.RWMutex
}

// NewThreadsafe returns a new TOData object, wrapped to be safe for multiple goroutine readers and a single writer.
func NewThreadsafe() TODataThreadsafe {
	return TODataThreadsafe{m: &sync.RWMutex{}, toData: New()}
}

// Get returns the current TOData. Callers MUST NOT modify returned data. Mutation IS NOT threadsafe
// If callers need to modify, a new GetMutable() should be added which copies.
func (d TODataThreadsafe) Get() TOData {
	d.m.RLock()
	defer d.m.RUnlock()
	return *d.toData
}

func (d TODataThreadsafe) set(newTOData TOData) {
	d.m.Lock()
	*d.toData = newTOData
	d.m.Unlock()
}

// CRConfig is the CrConfig data needed by TOData. Note this is not all data in the CRConfig.
// TODO change strings to type?
type CRConfig struct {
	ContentServers map[tc.CacheName]struct {
		DeliveryServices map[tc.DeliveryServiceName][]string `json:"deliveryServices"`
		CacheGroup       string                              `json:"cacheGroup"`
		Type             string                              `json:"type"`
	} `json:"contentServers"`
	DeliveryServices map[tc.DeliveryServiceName]struct {
		Topology  tc.TopologyName `json:"topology"`
		Matchsets []struct {
			Protocol  string `json:"protocol"`
			MatchList []struct {
				Regex string `json:"regex"`
				Type  string `json:"match-type"`
			} `json:"matchlist"`
		} `json:"matchsets"`
	} `json:"deliveryServices"`
	Topologies map[tc.TopologyName]struct {
		Nodes []string `json:"nodes"`
	}
}

// Update updates the TOData data with the last fetched CDN
func (d TODataThreadsafe) Update(to towrap.TrafficOpsSessionThreadsafe, cdn string, mc tc.TrafficMonitorConfigMap) error {
	crConfigBytes, _, err := to.LastCRConfig(cdn)
	if err != nil {
		return fmt.Errorf("getting last CRConfig: %v", err)
	}

	newTOData := TOData{}

	var crConfig CRConfig
	json := jsoniter.ConfigFastest
	err = json.Unmarshal(crConfigBytes, &crConfig)
	if err != nil {
		return fmt.Errorf("unmarshalling CRconfig: %v", err)
	}

	// TODO: remove the fallback behavior on the CRConfig in ATC 8.0 (https://github.com/apache/trafficcontrol/issues/6627)
	newTOData.DeliveryServiceServers, newTOData.ServerDeliveryServices = getDeliveryServiceServers(crConfig, mc)

	// TODO: remove the fallback behavior on the CRConfig in ATC 8.0 (https://github.com/apache/trafficcontrol/issues/6627)
	newTOData.DeliveryServiceTypes, err = getDeliveryServiceTypes(crConfig, mc)
	if err != nil {
		return fmt.Errorf("getting delivery service types from Traffic Ops: %v", err)
	}

	// TODO: remove the fallback behavior on the CRConfig in ATC 8.0 (https://github.com/apache/trafficcontrol/issues/6627)
	newTOData.DeliveryServiceRegexes, err = getDeliveryServiceRegexes(crConfig, mc)
	if err != nil {
		return fmt.Errorf("getting delivery service regexes from Traffic Ops: %v", err)
	}

	newTOData.ServerCachegroups = getServerCachegroups(mc)

	newTOData.ServerTypes, err = getServerTypes(mc)
	if err != nil {
		return fmt.Errorf("getting server types from monitoring config: %v", err)
	}

	d.set(newTOData)
	return nil
}

// getDeliveryServiceServers gets the servers on each delivery services, for the given CDN, from Traffic Ops.
func getDeliveryServiceServers(crc CRConfig, mc tc.TrafficMonitorConfigMap) (map[tc.DeliveryServiceName][]tc.CacheName, map[tc.CacheName][]tc.DeliveryServiceName) {
	dsServers := map[tc.DeliveryServiceName][]tc.CacheName{}
	serverDses := map[tc.CacheName][]tc.DeliveryServiceName{}

	topologyCacheGroupDses := map[string][]tc.DeliveryServiceName{}

	if canUseMonitorConfig(mc) {
		for deliveryServiceName, deliveryService := range mc.DeliveryService {
			if deliveryService.Topology == "" {
				continue
			}
			dsName := tc.DeliveryServiceName(deliveryServiceName)
			for _, cacheGroup := range mc.Topology[deliveryService.Topology].Nodes {
				topologyCacheGroupDses[cacheGroup] = append(topologyCacheGroupDses[cacheGroup], dsName)
			}
		}

		for serverName, serverData := range mc.TrafficServer {
			srvName := tc.CacheName(serverName)
			if cacheGroupDses, inTopology := topologyCacheGroupDses[serverData.CacheGroup]; inTopology {
				for _, deliveryServiceName := range cacheGroupDses {
					dsServers[deliveryServiceName] = append(dsServers[deliveryServiceName], srvName)
				}
				serverDses[srvName] = append(serverDses[srvName], cacheGroupDses...)
			}
			for _, serverDS := range serverData.DeliveryServices {
				dsName := tc.DeliveryServiceName(serverDS.XmlId)
				dsServers[dsName] = append(dsServers[dsName], srvName)
				serverDses[srvName] = append(serverDses[srvName], dsName)
			}

		}
		return dsServers, serverDses
	}

	for deliveryServiceName, deliveryService := range crc.DeliveryServices {
		if deliveryService.Topology == "" {
			continue
		}
		for _, cacheGroup := range crc.Topologies[deliveryService.Topology].Nodes {
			topologyCacheGroupDses[cacheGroup] = append(topologyCacheGroupDses[cacheGroup], deliveryServiceName)
		}
	}

	for serverName, serverData := range crc.ContentServers {
		if cacheGroupDses, inTopology := topologyCacheGroupDses[serverData.CacheGroup]; inTopology {
			for _, deliveryServiceName := range cacheGroupDses {
				dsServers[deliveryServiceName] = append(dsServers[deliveryServiceName], serverName)
			}
			serverDses[serverName] = append(serverDses[serverName], cacheGroupDses...)
		}
		for deliveryServiceName := range serverData.DeliveryServices {
			dsServers[deliveryServiceName] = append(dsServers[deliveryServiceName], serverName)
			serverDses[serverName] = append(serverDses[serverName], deliveryServiceName)
		}

	}
	return dsServers, serverDses
}

// getDeliveryServiceRegexes gets the regexes of each delivery service, for the given CDN, from Traffic Ops.
// Returns a map[deliveryService][]regex.
func getDeliveryServiceRegexes(crc CRConfig, mc tc.TrafficMonitorConfigMap) (Regexes, error) {
	dsRegexes := map[tc.DeliveryServiceName][]string{}

	if canUseMonitorConfig(mc) {
		for dsName, dsData := range mc.DeliveryService {
			if len(dsData.HostRegexes) < 1 {
				log.Warnln("TMConfig missing regex for delivery service " + dsName)
				continue
			}
			dsRegexes[tc.DeliveryServiceName(dsName)] = dsData.HostRegexes
		}
		return createRegexes(dsRegexes)
	}

	for dsName, dsData := range crc.DeliveryServices {
		for _, matchset := range dsData.Matchsets {
			if len(matchset.MatchList) < 1 {
				log.Warnln("CRConfig missing regex for delivery service '" + string(dsName) + "' matchset protocol '" + matchset.Protocol + "'")
				continue
			}
			if matchset.MatchList[0].Type != "HOST" {
				continue
			}
			dsRegexes[dsName] = append(dsRegexes[dsName], matchset.MatchList[0].Regex)
		}
		if len(dsRegexes[dsName]) == 0 {
			return Regexes{}, fmt.Errorf("CRConfig missing regex for '%s'", dsName)
		}
	}

	return createRegexes(dsRegexes)
}

// TODO precompute, move to TOData; call when we get new delivery services, instead of every time we create new stats
func createRegexes(dsToRegex map[tc.DeliveryServiceName][]string) (Regexes, error) {
	dsRegexes := Regexes{
		DirectMatches:                      map[string]tc.DeliveryServiceName{},
		DotStartSlashDotFooSlashDotDotStar: map[string]tc.DeliveryServiceName{},
		RegexMatch:                         map[*regexp.Regexp]tc.DeliveryServiceName{},
	}

	for ds, regexStrs := range dsToRegex {
		for _, regexStr := range regexStrs {
			prefix := `.*\.`
			suffix := `\..*`
			if strings.HasPrefix(regexStr, prefix) && strings.HasSuffix(regexStr, suffix) {
				matchStr := regexStr[len(prefix) : len(regexStr)-len(suffix)]
				if otherDs, ok := dsRegexes.DotStartSlashDotFooSlashDotDotStar[matchStr]; ok {
					return dsRegexes, fmt.Errorf("duplicate regex %s (%s) in %s and %s", regexStr, matchStr, ds, otherDs)
				}
				dsRegexes.DotStartSlashDotFooSlashDotDotStar[matchStr] = ds
				continue
			}
			if !strings.ContainsAny(regexStr, `[]^\:{}()|?+*,=%@<>!'`) {
				if otherDs, ok := dsRegexes.DirectMatches[regexStr]; ok {
					return dsRegexes, fmt.Errorf("duplicate Regex %s in %s and %s", regexStr, ds, otherDs)
				}
				dsRegexes.DirectMatches[regexStr] = ds
				continue
			}
			// TODO warn? regex matches are unusual
			r, err := regexp.Compile(regexStr)
			if err != nil {
				return dsRegexes, fmt.Errorf("regex %s failed to compile: %v", regexStr, err)
			}
			dsRegexes.RegexMatch[r] = ds
		}
	}
	return dsRegexes, nil
}

// getServerCachegroups gets the cachegroup of each ATS Edge+Mid Cache server, for the given CDN, from Traffic Ops.
// Returns a map[server]cachegroup.
func getServerCachegroups(mc tc.TrafficMonitorConfigMap) map[tc.CacheName]tc.CacheGroupName {
	serverCachegroups := map[tc.CacheName]tc.CacheGroupName{}

	for server, serverData := range mc.TrafficServer {
		serverCachegroups[tc.CacheName(server)] = tc.CacheGroupName(serverData.CacheGroup)
	}
	return serverCachegroups
}

// getServerTypes gets the cache type of each ATS Edge+Mid Cache server, for the given CDN, from Traffic Ops.
func getServerTypes(mc tc.TrafficMonitorConfigMap) (map[tc.CacheName]tc.CacheType, error) {
	serverTypes := map[tc.CacheName]tc.CacheType{}

	for server, serverData := range mc.TrafficServer {
		t := tc.CacheTypeFromString(serverData.Type)
		if t == tc.CacheTypeInvalid {
			return nil, fmt.Errorf("getServerTypes monitoring config unknown type for '%s': '%s'", server, serverData.Type)
		}
		serverTypes[tc.CacheName(server)] = t
	}
	return serverTypes, nil
}

// canUseMonitorConfig returns true if we can prefer monitor config data to crconfig data.
func canUseMonitorConfig(mc tc.TrafficMonitorConfigMap) bool {
	for _, dsData := range mc.DeliveryService {
		return dsData.Type != ""
	}
	return false
}

func getDeliveryServiceTypes(crc CRConfig, mc tc.TrafficMonitorConfigMap) (map[tc.DeliveryServiceName]tc.DSTypeCategory, error) {
	dsTypes := map[tc.DeliveryServiceName]tc.DSTypeCategory{}
	if canUseMonitorConfig(mc) {
		for dsName, dsData := range mc.DeliveryService {
			dsType := tc.DSTypeCategoryFromString(dsData.Type)
			if dsType == tc.DSTypeCategoryInvalid {
				log.Warnln("monitoring config invalid DS type category for delivery service '" + string(dsName) + "' matchset protocol '" + dsData.Type + "'; skipping")
				continue
			}
			dsTypes[tc.DeliveryServiceName(dsName)] = dsType
		}
		return dsTypes, nil
	}

	for dsName, dsData := range crc.DeliveryServices {
		if len(dsData.Matchsets) < 1 {
			return nil, fmt.Errorf("CRConfig missing protocol for '%s'", dsName)
		}
		dsTypeStr := dsData.Matchsets[0].Protocol
		dsType := tc.DSTypeCategoryFromString(dsTypeStr)
		if dsType == tc.DSTypeCategoryInvalid {
			log.Warnln("CRConfig invalid matchset protocol for delivery service '" + string(dsName) + "' matchset protocol '" + dsTypeStr + "'; skipping")
			continue
		}
		dsTypes[dsName] = dsType
	}
	return dsTypes, nil
}
