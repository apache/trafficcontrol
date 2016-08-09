package trafficopsdata

// TODO move to its own package?

import (
	"encoding/json"
	"fmt"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	towrap "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	"sync"
)

type TOData struct {
	DeliveryServiceServers map[string][]string
	ServerDeliveryServices map[string]string
	ServerTypes            map[enum.CacheName]enum.CacheType
	DeliveryServiceTypes   map[string]enum.DSType
	DeliveryServiceRegexes map[string][]string
	ServerCachegroups      map[enum.CacheName]enum.CacheGroupName
}

func New() *TOData {
	return &TOData{
		DeliveryServiceServers: map[string][]string{},
		ServerDeliveryServices: map[string]string{},
		ServerTypes:            map[enum.CacheName]enum.CacheType{},
		DeliveryServiceTypes:   map[string]enum.DSType{},
		DeliveryServiceRegexes: map[string][]string{},
		ServerCachegroups:      map[enum.CacheName]enum.CacheGroupName{},
	}
}

// This could be made lock-free, if the performance was necessary
type TODataThreadsafe struct {
	toData *TOData
	m      *sync.Mutex
}

func NewThreadsafe() TODataThreadsafe {
	return TODataThreadsafe{m: &sync.Mutex{}, toData: New()}
}

func (d TODataThreadsafe) Get() TOData {
	d.m.Lock()
	defer func() {
		d.m.Unlock()
	}()
	return *d.toData
}

func (d TODataThreadsafe) set(newTOData TOData) {
	d.m.Lock()
	*d.toData = newTOData
	d.m.Unlock()
}

// CRConfig is the CrConfig data needed by TOData. Note this is not all data in the CRConfig.
type CRConfig struct {
	ContentServers map[string]struct {
		DeliveryServices map[string][]string `json:"deliveryServices"`
		CacheGroup       string              `json:"cacheGroup"`
		Type             string              `json:"type"`
	} `json:"contentServers"`
	DeliveryServices map[string]struct {
		Matchsets []struct {
			Protocol  string `json:"protocol"`
			MatchList []struct {
				Regex string `json:"regex"`
			} `json:"matchlist"`
		} `json:"matchsets"`
	} `json:"deliveryServices"`
}

// Fetch gets the CRConfig from Traffic Ops, creates the TOData maps, and atomically sets the TOData.
// TODO since the session is threadsafe, each TOData get func below could be put in a goroutine, if performance mattered
// TODO change called funcs to take CRConfigRaw instead of toSession
func (d TODataThreadsafe) Fetch(to towrap.ITrafficOpsSession, cdn string) error {
	newTOData := TOData{}

	crConfigBytes, err := to.CRConfigRaw(cdn)
	if err != nil {
		return fmt.Errorf("Error getting CRconfig from Traffic Ops: %v", err)
	}
	var crConfig CRConfig
	if err := json.Unmarshal(crConfigBytes, &crConfig); err != nil {
		return fmt.Errorf("Error unmarshalling CRconfig: %v", err)
	}

	newTOData.DeliveryServiceServers, newTOData.ServerDeliveryServices, err = getDeliveryServiceServers(crConfig)
	if err != nil {
		return err
	}

	newTOData.DeliveryServiceTypes, err = getDeliveryServiceTypes(crConfig)
	if err != nil {
		return fmt.Errorf("Error getting delivery service types from Traffic Ops: %v\n", err)
	}

	newTOData.DeliveryServiceRegexes, err = getDeliveryServiceRegexes(crConfig)
	if err != nil {
		return fmt.Errorf("Error getting delivery service regexes from Traffic Ops: %v\n", err)
	}

	newTOData.ServerCachegroups, err = getServerCachegroups(crConfig)
	if err != nil {
		return fmt.Errorf("Error getting server cachegroups from Traffic Ops: %v\n", err)
	}

	newTOData.ServerTypes, err = getServerTypes(crConfig)
	if err != nil {
		return fmt.Errorf("Error getting server types from Traffic Ops: %v\n", err)
	}

	d.set(newTOData)
	return nil
}

// getDeliveryServiceServers gets the servers on each delivery services, for the given CDN, from Traffic Ops.
// Returns a map[deliveryService][]server, and a map[server]deliveryService
func getDeliveryServiceServers(crc CRConfig) (map[string][]string, map[string]string, error) {
	dsServers := map[string][]string{}
	serverDs := map[string]string{}

	for serverName, serverData := range crc.ContentServers {
		for deliveryServiceName, _ := range serverData.DeliveryServices {
			dsServers[deliveryServiceName] = append(dsServers[deliveryServiceName], serverName)
			serverDs[serverName] = deliveryServiceName
		}
	}
	return dsServers, serverDs, nil
}

// getDeliveryServiceRegexes gets the regexes of each delivery service, for the given CDN, from Traffic Ops.
// Returns a map[deliveryService][]regex.
func getDeliveryServiceRegexes(crc CRConfig) (map[string][]string, error) {
	dsRegexes := map[string][]string{}

	for dsName, dsData := range crc.DeliveryServices {
		if len(dsData.Matchsets) < 1 {
			return nil, fmt.Errorf("CRConfig missing regex for '%s'", dsName)
		}
		for _, matchset := range dsData.Matchsets {
			if len(matchset.MatchList) < 1 {
				return nil, fmt.Errorf("CRConfig missing Regex for '%s'", dsName)
			}
			dsRegexes[dsName] = append(dsRegexes[dsName], matchset.MatchList[0].Regex)
		}
	}
	return dsRegexes, nil
}

// getServerCachegroups gets the cachegroup of each ATS Edge+Mid Cache server, for the given CDN, from Traffic Ops.
// Returns a map[server]cachegroup.
func getServerCachegroups(crc CRConfig) (map[enum.CacheName]enum.CacheGroupName, error) {
	serverCachegroups := map[enum.CacheName]enum.CacheGroupName{}

	for server, serverData := range crc.ContentServers {
		serverCachegroups[enum.CacheName(server)] = enum.CacheGroupName(serverData.CacheGroup)
	}
	return serverCachegroups, nil
}

// getServerTypes gets the cache type of each ATS Edge+Mid Cache server, for the given CDN, from Traffic Ops.
func getServerTypes(crc CRConfig) (map[enum.CacheName]enum.CacheType, error) {
	serverTypes := map[enum.CacheName]enum.CacheType{}

	for server, serverData := range crc.ContentServers {
		t := enum.CacheTypeFromString(serverData.Type)
		if t == enum.CacheTypeInvalid {
			return nil, fmt.Errorf("getServerTypes CRConfig unknown type for '%s': '%s'", server, serverData.Type)
		}
		serverTypes[enum.CacheName(server)] = t
	}
	return serverTypes, nil
}

func getDeliveryServiceTypes(crc CRConfig) (map[string]enum.DSType, error) {
	dsTypes := map[string]enum.DSType{}

	for dsName, dsData := range crc.DeliveryServices {
		if len(dsData.Matchsets) < 1 {
			return nil, fmt.Errorf("CRConfig missing protocol for '%s'", dsName)
		}
		dsTypeStr := dsData.Matchsets[0].Protocol
		dsType := enum.DSTypeFromString(dsTypeStr)
		if dsType == enum.DSTypeInvalid {
			return nil, fmt.Errorf("CRConfig unknowng protocol for '%s': '%s'", dsName, dsTypeStr)
		}
		dsTypes[dsName] = dsType
	}
	return dsTypes, nil
}
