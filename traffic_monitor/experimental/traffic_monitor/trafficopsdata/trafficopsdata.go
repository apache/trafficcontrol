package trafficopsdata

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	"regexp"
	"strings"
	"sync"
)

// DsRegexes maps Delivery Service Regular Expressions to delivery services.
// For performance, we categorize Regular Expressions into 3 categories:
// 1. Direct string matches, with no regular expression matching characters
// 2. .*\.foo\..* expressions, where foo is a direct string match with no regular expression matching characters
// 3. Everything else
// This allows us to do a cheap match on 1 and 2, and only regex match the uncommon case.
// TODO performance tests, whether Go compiled *Regexp is relevantly slower than `strings.Contains` for direct and .foo. matches
type Regexes struct {
	DirectMatches                      map[string]enum.DeliveryServiceName
	DotStartSlashDotFooSlashDotDotStar map[string]enum.DeliveryServiceName
	RegexMatch                         map[*regexp.Regexp]enum.DeliveryServiceName
}

// DeliveryService returns the delivery service which matches the given fqdn, or false.
func (d Regexes) DeliveryService(fqdn string) (enum.DeliveryServiceName, bool) {
	if ds, ok := d.DirectMatches[fqdn]; ok {
		return ds, true
	}
	for matchStr, ds := range d.DotStartSlashDotFooSlashDotDotStar {
		if strings.Contains(fqdn, "."+matchStr+".") {
			return ds, true
		}
	}
	for regex, ds := range d.RegexMatch {
		if regex.MatchString(fqdn) {
			return ds, true
		}
	}
	return "", false
}

func NewRegexes() Regexes {
	return Regexes{DirectMatches: map[string]enum.DeliveryServiceName{}, DotStartSlashDotFooSlashDotDotStar: map[string]enum.DeliveryServiceName{}, RegexMatch: map[*regexp.Regexp]enum.DeliveryServiceName{}}
}

type TOData struct {
	DeliveryServiceServers map[enum.DeliveryServiceName][]enum.CacheName
	ServerDeliveryServices map[enum.CacheName][]enum.DeliveryServiceName
	ServerTypes            map[enum.CacheName]enum.CacheType
	DeliveryServiceTypes   map[enum.DeliveryServiceName]enum.DSType
	DeliveryServiceRegexes Regexes
	ServerCachegroups      map[enum.CacheName]enum.CacheGroupName
}

func New() *TOData {
	return &TOData{
		DeliveryServiceServers: map[enum.DeliveryServiceName][]enum.CacheName{},
		ServerDeliveryServices: map[enum.CacheName][]enum.DeliveryServiceName{},
		ServerTypes:            map[enum.CacheName]enum.CacheType{},
		DeliveryServiceTypes:   map[enum.DeliveryServiceName]enum.DSType{},
		DeliveryServiceRegexes: NewRegexes(),
		ServerCachegroups:      map[enum.CacheName]enum.CacheGroupName{},
	}
}

// This could be made lock-free, if the performance was necessary
type TODataThreadsafe struct {
	toData *TOData
	m      *sync.RWMutex
}

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
	ContentServers map[enum.CacheName]struct {
		DeliveryServices map[enum.DeliveryServiceName][]string `json:"deliveryServices"`
		CacheGroup       string                                `json:"cacheGroup"`
		Type             string                                `json:"type"`
	} `json:"contentServers"`
	DeliveryServices map[enum.DeliveryServiceName]struct {
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
func getDeliveryServiceServers(crc CRConfig) (map[enum.DeliveryServiceName][]enum.CacheName, map[enum.CacheName][]enum.DeliveryServiceName, error) {
	dsServers := map[enum.DeliveryServiceName][]enum.CacheName{}
	serverDses := map[enum.CacheName][]enum.DeliveryServiceName{}

	for serverName, serverData := range crc.ContentServers {
		for deliveryServiceName, _ := range serverData.DeliveryServices {
			dsServers[deliveryServiceName] = append(dsServers[deliveryServiceName], serverName)
			serverDses[serverName] = append(serverDses[serverName], deliveryServiceName)
		}
	}
	return dsServers, serverDses, nil
}

// getDeliveryServiceRegexes gets the regexes of each delivery service, for the given CDN, from Traffic Ops.
// Returns a map[deliveryService][]regex.
func getDeliveryServiceRegexes(crc CRConfig) (Regexes, error) {
	dsRegexes := map[enum.DeliveryServiceName][]string{}

	for dsName, dsData := range crc.DeliveryServices {
		if len(dsData.Matchsets) < 1 {
			return Regexes{}, fmt.Errorf("CRConfig missing regex for '%s'", dsName)
		}
		for _, matchset := range dsData.Matchsets {
			if len(matchset.MatchList) < 1 {
				return Regexes{}, fmt.Errorf("CRConfig missing Regex for '%s'", dsName)
			}
			dsRegexes[dsName] = append(dsRegexes[dsName], matchset.MatchList[0].Regex)
		}
	}

	return createRegexes(dsRegexes)
}

// TODO precompute, move to TOData; call when we get new delivery services, instead of every time we create new stats
func createRegexes(dsToRegex map[enum.DeliveryServiceName][]string) (Regexes, error) {
	dsRegexes := Regexes{
		DirectMatches:                      map[string]enum.DeliveryServiceName{},
		DotStartSlashDotFooSlashDotDotStar: map[string]enum.DeliveryServiceName{},
		RegexMatch:                         map[*regexp.Regexp]enum.DeliveryServiceName{},
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
func getServerCachegroups(crc CRConfig) (map[enum.CacheName]enum.CacheGroupName, error) {
	serverCachegroups := map[enum.CacheName]enum.CacheGroupName{}

	for server, serverData := range crc.ContentServers {
		serverCachegroups[server] = enum.CacheGroupName(serverData.CacheGroup)
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
		serverTypes[server] = t
	}
	return serverTypes, nil
}

func getDeliveryServiceTypes(crc CRConfig) (map[enum.DeliveryServiceName]enum.DSType, error) {
	dsTypes := map[enum.DeliveryServiceName]enum.DSType{}

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
