package datareq

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/cache"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
)

// CacheStatFilter fulfills the cache.Filter interface, for filtering stats. See the `NewCacheStatFilter` documentation for details on which query parameters are used to filter.
type CacheStatFilter struct {
	historyCount int
	statsToUse   map[string]struct{}
	wildcard     bool
	cacheType    enum.CacheType
	hosts        map[enum.CacheName]struct{}
	cacheTypes   map[enum.CacheName]enum.CacheType
}

// UseCache returns whether the given cache is in the filter.
func (f *CacheStatFilter) UseCache(name enum.CacheName) bool {
	if _, inHosts := f.hosts[name]; len(f.hosts) != 0 && !inHosts {
		return false
	}
	if f.cacheType != enum.CacheTypeInvalid && f.cacheTypes[name] != f.cacheType {
		return false
	}
	return true
}

// UseStat returns whether the given stat is in the filter.
func (f *CacheStatFilter) UseStat(statName string) bool {
	if len(f.statsToUse) == 0 {
		return true
	}
	if !f.wildcard {
		_, ok := f.statsToUse[statName]
		return ok
	}
	for statToUse := range f.statsToUse {
		if strings.Contains(statName, statToUse) {
			return true
		}
	}
	return false
}

// WithinStatHistoryMax returns whether the given history index is less than the max history of this filter.
func (f *CacheStatFilter) WithinStatHistoryMax(n int) bool {
	if f.historyCount == 0 {
		return true
	}
	if n <= f.historyCount {
		return true
	}
	return false
}

// NewCacheStatFilter takes the HTTP query parameters and creates a CacheStatFilter which fulfills the `cache.Filter` interface, filtering according to the query parameters passed.
// Query parameters used are `hc`, `stats`, `wildcard`, `type`, and `hosts`.
// If `hc` is 0, all history is returned. If `hc` is empty, 1 history is returned.
// If `stats` is empty, all stats are returned.
// If `wildcard` is empty, `stats` is considered exact.
// If `type` is empty, all cache types are returned.
func NewCacheStatFilter(path string, params url.Values, cacheTypes map[enum.CacheName]enum.CacheType) (cache.Filter, error) {
	validParams := map[string]struct{}{
		"hc":       struct{}{},
		"stats":    struct{}{},
		"wildcard": struct{}{},
		"type":     struct{}{},
		"hosts":    struct{}{},
		"cache":    struct{}{},
	}
	if len(params) > len(validParams) {
		return nil, fmt.Errorf("invalid query parameters")
	}
	for param := range params {
		if _, ok := validParams[param]; !ok {
			return nil, fmt.Errorf("invalid query parameter '%v'", param)
		}
	}

	historyCount := 1
	if paramHc, exists := params["hc"]; exists && len(paramHc) > 0 {
		v, err := strconv.Atoi(paramHc[0])
		if err == nil {
			historyCount = v
		}
	}

	statsToUse := map[string]struct{}{}
	if paramStats, exists := params["stats"]; exists && len(paramStats) > 0 {
		commaStats := strings.Split(paramStats[0], ",")
		for _, stat := range commaStats {
			statsToUse[stat] = struct{}{}
		}
	}

	wildcard := false
	if paramWildcard, exists := params["wildcard"]; exists && len(paramWildcard) > 0 {
		wildcard, _ = strconv.ParseBool(paramWildcard[0]) // ignore errors, error => false
	}

	cacheType := enum.CacheTypeInvalid
	if paramType, exists := params["type"]; exists && len(paramType) > 0 {
		cacheType = enum.CacheTypeFromString(paramType[0])
		if cacheType == enum.CacheTypeInvalid {
			return nil, fmt.Errorf("invalid query parameter type '%v' - valid types are: {edge, mid}", paramType[0])
		}
	}

	hosts := map[enum.CacheName]struct{}{}
	if paramHosts, exists := params["hosts"]; exists && len(paramHosts) > 0 {
		commaHosts := strings.Split(paramHosts[0], ",")
		for _, host := range commaHosts {
			hosts[enum.CacheName(host)] = struct{}{}
		}
	}
	if paramHosts, exists := params["cache"]; exists && len(paramHosts) > 0 {
		commaHosts := strings.Split(paramHosts[0], ",")
		for _, host := range commaHosts {
			hosts[enum.CacheName(host)] = struct{}{}
		}
	}

	pathArgument := getPathArgument(path)
	if pathArgument != "" {
		hosts[enum.CacheName(pathArgument)] = struct{}{}
	}

	// parameters without values are considered hosts, e.g. `?my-cache-0`
	for maybeHost, val := range params {
		if len(val) == 0 || (len(val) == 1 && val[0] == "") {
			hosts[enum.CacheName(maybeHost)] = struct{}{}
		}
	}

	return &CacheStatFilter{
		historyCount: historyCount,
		statsToUse:   statsToUse,
		wildcard:     wildcard,
		cacheType:    cacheType,
		hosts:        hosts,
		cacheTypes:   cacheTypes,
	}, nil
}
