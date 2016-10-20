package manager

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/cache"
	ds "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	dsdata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservicedata"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/srvhttp"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
)

type JSONEvents struct {
	Events []Event `json:"events"`
}

type CacheState struct {
	Value bool `json:"value"`
}

type APIPeerStates struct {
	srvhttp.CommonAPIData
	Peers map[enum.TrafficMonitorName]map[enum.CacheName][]CacheState `json:"peers"`
}

// TODO make fields nullable, so error fields can be omitted, letting API callers still get updates for unerrored fields
type CacheStatus struct {
	Type                  *string  `json:"type,omitempty"`
	Status                *string  `json:"status,omitempty"`
	LoadAverage           *float64 `json:"load_average,omitempty"`
	QueryTimeMilliseconds *int64   `json:"query_time_ms,omitempty"`
	BandwidthKbps         *float64 `json:"bandwidth_kbps,omitempty"`
	ConnectionCount       *int64   `json:"connection_count,omitempty"`
}

// CacheStatFilter fulfills the cache.Filter interface, for filtering stats. See the `NewCacheStatFilter` documentation for details on which query parameters are used to filter.
type CacheStatFilter struct {
	historyCount int
	statsToUse   map[string]struct{}
	wildcard     bool
	cacheType    enum.CacheType
	hosts        map[enum.CacheName]struct{}
	cacheTypes   map[enum.CacheName]enum.CacheType
}

func (f *CacheStatFilter) UseCache(name enum.CacheName) bool {
	if _, inHosts := f.hosts[name]; len(f.hosts) != 0 && !inHosts {
		return false
	}
	if f.cacheType != enum.CacheTypeInvalid && f.cacheTypes[name] != f.cacheType {
		return false
	}
	return true
}

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
func NewCacheStatFilter(params url.Values, cacheTypes map[enum.CacheName]enum.CacheType) (cache.Filter, error) {
	validParams := map[string]struct{}{"hc": struct{}{}, "stats": struct{}{}, "wildcard": struct{}{}, "type": struct{}{}, "hosts": struct{}{}}
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

// DSStatFilter fulfills the cache.Filter interface, for filtering stats. See the `NewDSStatFilter` documentation for details on which query parameters are used to filter.
type DSStatFilter struct {
	historyCount     int
	statsToUse       map[string]struct{}
	wildcard         bool
	dsType           enum.DSType
	deliveryServices map[enum.DeliveryServiceName]struct{}
	dsTypes          map[enum.DeliveryServiceName]enum.DSType
}

func (f *DSStatFilter) UseDeliveryService(name enum.DeliveryServiceName) bool {
	if _, inDSes := f.deliveryServices[name]; len(f.deliveryServices) != 0 && !inDSes {
		return false
	}
	if f.dsType != enum.DSTypeInvalid && f.dsTypes[name] != f.dsType {
		return false
	}
	return true
}

func (f *DSStatFilter) UseStat(statName string) bool {
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

func (f *DSStatFilter) WithinStatHistoryMax(n int) bool {
	if f.historyCount == 0 {
		return true
	}
	if n <= f.historyCount {
		return true
	}
	return false
}

// NewDSStatFilter takes the HTTP query parameters and creates a cache.Filter, filtering according to the query parameters passed.
// Query parameters used are `hc`, `stats`, `wildcard`, `type`, and `deliveryservices`.
// If `hc` is 0, all history is returned. If `hc` is empty, 1 history is returned.
// If `stats` is empty, all stats are returned.
// If `wildcard` is empty, `stats` is considered exact.
// If `type` is empty, all types are returned.
func NewDSStatFilter(params url.Values, dsTypes map[enum.DeliveryServiceName]enum.DSType) (dsdata.Filter, error) {
	validParams := map[string]struct{}{"hc": struct{}{}, "stats": struct{}{}, "wildcard": struct{}{}, "type": struct{}{}, "deliveryservices": struct{}{}}
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

	dsType := enum.DSTypeInvalid
	if paramType, exists := params["type"]; exists && len(paramType) > 0 {
		dsType = enum.DSTypeFromString(paramType[0])
		if dsType == enum.DSTypeInvalid {
			return nil, fmt.Errorf("invalid query parameter type '%v' - valid types are: {http, dns}", paramType[0])
		}
	}

	deliveryServices := map[enum.DeliveryServiceName]struct{}{}
	// TODO rename 'hosts' to 'names' for consistency
	if paramNames, exists := params["deliveryservices"]; exists && len(paramNames) > 0 {
		commaNames := strings.Split(paramNames[0], ",")
		for _, name := range commaNames {
			deliveryServices[enum.DeliveryServiceName(name)] = struct{}{}
		}
	}
	// parameters without values are considered names, e.g. `?my-cache-0` or `?my-delivery-service`
	for maybeName, val := range params {
		if len(val) == 0 || (len(val) == 1 && val[0] == "") {
			deliveryServices[enum.DeliveryServiceName(maybeName)] = struct{}{}
		}
	}

	return &DSStatFilter{
		historyCount:     historyCount,
		statsToUse:       statsToUse,
		wildcard:         wildcard,
		dsType:           dsType,
		deliveryServices: deliveryServices,
		dsTypes:          dsTypes,
	}, nil
}

// PeerStateFilter fulfills the cache.Filter interface, for filtering stats. See the `NewPeerStateFilter` documentation for details on which query parameters are used to filter.
type PeerStateFilter struct {
	historyCount int
	cachesToUse  map[enum.CacheName]struct{}
	peersToUse   map[enum.TrafficMonitorName]struct{}
	wildcard     bool
	cacheType    enum.CacheType
	cacheTypes   map[enum.CacheName]enum.CacheType
}

func (f *PeerStateFilter) UsePeer(name enum.TrafficMonitorName) bool {
	if _, inPeers := f.peersToUse[name]; len(f.peersToUse) != 0 && !inPeers {
		return false
	}
	return true
}

func (f *PeerStateFilter) UseCache(name enum.CacheName) bool {
	if f.cacheType != enum.CacheTypeInvalid && f.cacheTypes[name] != f.cacheType {
		return false
	}

	if len(f.cachesToUse) == 0 {
		return true
	}

	if !f.wildcard {
		_, ok := f.cachesToUse[name]
		return ok
	}
	for cacheToUse := range f.cachesToUse {
		if strings.Contains(string(name), string(cacheToUse)) {
			return true
		}
	}
	return false
}

func (f *PeerStateFilter) WithinStatHistoryMax(n int) bool {
	if f.historyCount == 0 {
		return true
	}
	if n <= f.historyCount {
		return true
	}
	return false
}

// NewPeerStateFilter takes the HTTP query parameters and creates a cache.Filter, filtering according to the query parameters passed.
// Query parameters used are `hc`, `stats`, `wildcard`, `typep`, and `hosts`. The `stats` param filters caches. The `hosts` param filters peer Traffic Monitors. The `type` param filters cache types (edge, mid).
// If `hc` is 0, all history is returned. If `hc` is empty, 1 history is returned.
// If `stats` is empty, all stats are returned.
// If `wildcard` is empty, `stats` is considered exact.
// If `type` is empty, all cache types are returned.
func NewPeerStateFilter(params url.Values, cacheTypes map[enum.CacheName]enum.CacheType) (*PeerStateFilter, error) {
	// TODO change legacy `stats` and `hosts` to `caches` and `monitors` (or `peers`).
	validParams := map[string]struct{}{"hc": struct{}{}, "stats": struct{}{}, "wildcard": struct{}{}, "type": struct{}{}, "peers": struct{}{}}
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

	cachesToUse := map[enum.CacheName]struct{}{}
	// TODO rename 'stats' to 'caches'
	if paramStats, exists := params["stats"]; exists && len(paramStats) > 0 {
		commaStats := strings.Split(paramStats[0], ",")
		for _, stat := range commaStats {
			cachesToUse[enum.CacheName(stat)] = struct{}{}
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

	peersToUse := map[enum.TrafficMonitorName]struct{}{}
	if paramNames, exists := params["peers"]; exists && len(paramNames) > 0 {
		commaNames := strings.Split(paramNames[0], ",")
		for _, name := range commaNames {
			peersToUse[enum.TrafficMonitorName(name)] = struct{}{}
		}
	}
	// parameters without values are considered names, e.g. `?my-cache-0` or `?my-delivery-service`
	for maybeName, val := range params {
		if len(val) == 0 || (len(val) == 1 && val[0] == "") {
			peersToUse[enum.TrafficMonitorName(maybeName)] = struct{}{}
		}
	}

	return &PeerStateFilter{
		historyCount: historyCount,
		cachesToUse:  cachesToUse,
		wildcard:     wildcard,
		cacheType:    cacheType,
		peersToUse:   peersToUse,
		cacheTypes:   cacheTypes,
	}, nil
}

// DataRequest takes an `srvhttp.DataRequest`, and the monitored data objects, and returns the appropriate response, and the status code.
func DataRequest(
	req srvhttp.DataRequest,
	opsConfig OpsConfigThreadsafe,
	toSession towrap.ITrafficOpsSession,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	statHistory StatHistoryThreadsafe,
	dsStats DSStatsThreadsafe,
	events EventsThreadsafe,
	staticAppData StaticAppData,
	healthPollInterval time.Duration,
	lastHealthDurations DurationMapThreadsafe,
	fetchCount UintThreadsafe,
	healthIteration UintThreadsafe,
	errorCount UintThreadsafe,
	toData todata.TODataThreadsafe,
	localCacheStatus CacheAvailableStatusThreadsafe,
	lastStats LastStatsThreadsafe,
	unpolledCaches UnpolledCachesThreadsafe,
) (body []byte, responseCode int) {

	// handleErr takes an error, and the request type it came from, and logs. It is ok to call with a nil error, in which case this is a no-op.
	handleErr := func(err error) {
		if err == nil {
			return
		}
		errorCount.Inc()
		log.Errorf("Request Error: %v\n", fmt.Errorf(req.Type.String()+": %v", err))
	}

	// commonReturn takes the body, err, and the data request Type which has been processed. It logs and deals with any error, and returns the appropriate bytes and response code for the `srvhttp`.
	commonReturn := func(body []byte, err error) ([]byte, int) {
		if err == nil {
			return body, http.StatusOK
		}
		handleErr(err)
		return nil, http.StatusInternalServerError
	}

	if unpolledCaches.Any() {
		handleErr(fmt.Errorf("service still starting, some caches unpolled"))
		return []byte("Service Unavailable"), http.StatusServiceUnavailable
	}

	var err error
	switch req.Type {
	case srvhttp.TRConfig:
		cdnName := opsConfig.Get().CdnName
		if toSession == nil {
			return commonReturn(nil, fmt.Errorf("Unable to connect to Traffic Ops"))
		}
		if cdnName == "" {
			return commonReturn(nil, fmt.Errorf("No CDN Configured"))
		}
		return commonReturn(body, err)
	case srvhttp.TRStateDerived:
		body, err = peer.CrstatesMarshall(combinedStates.Get())
		return commonReturn(body, err)
	case srvhttp.TRStateSelf:
		body, err = peer.CrstatesMarshall(localStates.Get())
		return commonReturn(body, err)
	case srvhttp.CacheStats:
		filter, err := NewCacheStatFilter(req.Parameters, toData.Get().ServerTypes)
		if err != nil {
			handleErr(err)
			return []byte(err.Error()), http.StatusBadRequest
		}
		body, err = cache.StatsMarshall(statHistory.Get(), filter, req.Parameters)
		return commonReturn(body, err)
	case srvhttp.DSStats:
		filter, err := NewDSStatFilter(req.Parameters, toData.Get().DeliveryServiceTypes)
		if err != nil {
			handleErr(err)
			return []byte(err.Error()), http.StatusBadRequest
		}
		body, err = json.Marshal(dsStats.Get().JSON(filter, req.Parameters)) // TODO marshall beforehand, for performance? (test to see how often requests are made)
		return commonReturn(body, err)
	case srvhttp.EventLog:
		body, err = json.Marshal(JSONEvents{Events: events.Get()})
		return commonReturn(body, err)
	case srvhttp.PeerStates:
		filter, err := NewPeerStateFilter(req.Parameters, toData.Get().ServerTypes)
		if err != nil {
			handleErr(err)
			return []byte(err.Error()), http.StatusBadRequest
		}

		body, err = json.Marshal(createAPIPeerStates(peerStates.Get(), filter, req.Parameters))
		return commonReturn(body, err)
	case srvhttp.StatSummary:
		return nil, http.StatusNotImplemented
	case srvhttp.Stats:
		body, err = getStats(staticAppData, healthPollInterval, lastHealthDurations.Get(), fetchCount.Get(), healthIteration.Get(), errorCount.Get())
		return commonReturn(body, err)
	case srvhttp.ConfigDoc:
		opsConfigCopy := opsConfig.Get()
		// if the password is blank, leave it blank, so callers can see it's missing.
		if opsConfigCopy.Password != "" {
			opsConfigCopy.Password = "*****"
		}
		body, err = json.Marshal(opsConfigCopy)
		return commonReturn(body, err)
	case srvhttp.APICacheCount: // TODO determine if this should use peerStates
		return []byte(strconv.Itoa(len(localStates.Get().Caches))), http.StatusOK
	case srvhttp.APICacheAvailableCount:
		return []byte(strconv.Itoa(cacheAvailableCount(localStates.Get().Caches))), http.StatusOK
	case srvhttp.APICacheDownCount:
		return []byte(strconv.Itoa(cacheDownCount(localStates.Get().Caches))), http.StatusOK
	case srvhttp.APIVersion:
		s := "traffic_monitor-" + staticAppData.Version + "."
		if len(staticAppData.GitRevision) > 6 {
			s += staticAppData.GitRevision[:6]
		} else {
			s += staticAppData.GitRevision
		}
		return []byte(s), http.StatusOK
	case srvhttp.APITrafficOpsURI:
		return []byte(opsConfig.Get().Url), http.StatusOK
	case srvhttp.APICacheStates:
		body, err = json.Marshal(createCacheStatuses(toData.Get().ServerTypes, statHistory.Get(),
			lastHealthDurations.Get(), localStates.Get().Caches, lastStats.Get(), localCacheStatus))
		return commonReturn(body, err)
	case srvhttp.APIBandwidthKbps:
		serverTypes := toData.Get().ServerTypes
		kbpsStats := lastStats.Get()
		sum := float64(0.0)
		for cache, data := range kbpsStats.Caches {
			if serverTypes[cache] != enum.CacheTypeEdge {
				continue
			}
			sum += data.Bytes.PerSec / ds.BytesPerKilobit
		}
		return []byte(fmt.Sprintf("%f", sum)), http.StatusOK
	case srvhttp.APIBandwidthCapacityKbps:
		statHistory := statHistory.Get()
		cap := int64(0)
		for _, results := range statHistory {
			if len(results) == 0 {
				continue
			}
			cap += results[0].MaxKbps
		}
		return []byte(fmt.Sprintf("%d", cap)), http.StatusOK
	default:
		return commonReturn(nil, fmt.Errorf("Unknown Request Type"))
	}
}

func createCacheStatuses(
	cacheTypes map[enum.CacheName]enum.CacheType,
	statHistory map[enum.CacheName][]cache.Result,
	lastHealthDurations map[enum.CacheName]time.Duration,
	cacheStates map[enum.CacheName]peer.IsAvailable,
	lastStats ds.LastStats,
	localCacheStatusThreadsafe CacheAvailableStatusThreadsafe,
) map[enum.CacheName]CacheStatus {
	conns := createCacheConnections(statHistory)
	statii := map[enum.CacheName]CacheStatus{}
	localCacheStatus := localCacheStatusThreadsafe.Get()

	for cacheName, cacheType := range cacheTypes {
		cacheStatHistory, ok := statHistory[cacheName]
		if !ok {
			log.Warnf("createCacheStatuses stat history missing cache %s\n", cacheName)
			continue
		}

		if len(cacheStatHistory) < 1 {
			log.Warnf("createCacheStatuses stat history empty for cache %s\n", cacheName)
			continue
		}

		log.Debugf("createCacheStatuses NOT empty for cache %s\n", cacheName)

		var loadAverage *float64
		procLoadAvg := cacheStatHistory[0].Astats.System.ProcLoadavg
		if procLoadAvg != "" {
			firstSpace := strings.IndexRune(procLoadAvg, ' ')
			if firstSpace == -1 {
				log.Warnf("WARNING unexpected proc.loadavg '%s' for cache %s\n", procLoadAvg, cacheName)
			} else {
				loadAverageVal, err := strconv.ParseFloat(procLoadAvg[:firstSpace], 64)
				if err != nil {
					log.Warnf("proc.loadavg doesn't contain a float prefix '%s' for cache %s\n", procLoadAvg, cacheName)
				} else {
					loadAverage = &loadAverageVal
				}
			}
		}

		var queryTime *int64
		queryTimeVal, ok := lastHealthDurations[cacheName]
		if !ok {
			log.Warnf("cache not in last health durations cache %s\n", cacheName)
		} else {
			queryTimeInt := int64(queryTimeVal / time.Millisecond)
			queryTime = &queryTimeInt
		}

		var kbps *float64
		lastStat, ok := lastStats.Caches[enum.CacheName(cacheName)]
		if !ok {
			log.Warnf("cache not in last kbps cache %s\n", cacheName)
		} else {
			kbpsVal := lastStat.Bytes.PerSec / float64(ds.BytesPerKilobit)
			kbps = &kbpsVal
		}

		var connections *int64
		connectionsVal, ok := conns[enum.CacheName(cacheName)]
		if !ok {
			log.Warnf("cache not in connections %s\n", cacheName)
		} else {
			connections = &connectionsVal
		}

		var status *string
		statusVal, ok := localCacheStatus[enum.CacheName(cacheName)]
		if !ok {
			log.Warnf("cache not in statuses %s\n", cacheName)
		} else {
			statusString := statusVal.Status + " - "
			if localCacheStatus[enum.CacheName(cacheName)].Available {
				statusString += "available"
			} else {
				statusString += "unavailable"
			}
			status = &statusString
		}

		cacheTypeStr := string(cacheType)
		statii[enum.CacheName(cacheName)] = CacheStatus{Type: &cacheTypeStr, LoadAverage: loadAverage, QueryTimeMilliseconds: queryTime, BandwidthKbps: kbps, ConnectionCount: connections, Status: status}
	}
	return statii
}

func createCacheConnections(statHistory map[enum.CacheName][]cache.Result) map[enum.CacheName]int64 {
	conns := map[enum.CacheName]int64{}
	for server, history := range statHistory {
		for _, result := range history {
			val, ok := result.Astats.Ats["proxy.process.http.current_client_connections"]
			if !ok {
				continue
			}

			v, ok := val.(float64)
			if !ok {
				continue
			}

			conns[server] = int64(v)
			break
		}
	}
	return conns
}

func cacheDownCount(caches map[enum.CacheName]peer.IsAvailable) int {
	count := 0
	for _, available := range caches {
		if !available.IsAvailable {
			count++
		}
	}
	return count
}

func cacheAvailableCount(caches map[enum.CacheName]peer.IsAvailable) int {
	return len(caches) - cacheDownCount(caches)
}

func createAPIPeerStates(peerStates map[enum.TrafficMonitorName]peer.Crstates, filter *PeerStateFilter, params url.Values) APIPeerStates {
	apiPeerStates := APIPeerStates{
		CommonAPIData: srvhttp.GetCommonAPIData(params, time.Now()),
		Peers:         map[enum.TrafficMonitorName]map[enum.CacheName][]CacheState{},
	}

	for peer, state := range peerStates {
		if !filter.UsePeer(peer) {
			continue
		}
		if _, ok := apiPeerStates.Peers[peer]; !ok {
			apiPeerStates.Peers[peer] = map[enum.CacheName][]CacheState{}
		}
		peerState := apiPeerStates.Peers[peer]
		for cache, available := range state.Caches {
			if !filter.UseCache(cache) {
				continue
			}
			peerState[cache] = []CacheState{CacheState{Value: available.IsAvailable}}
		}
		apiPeerStates.Peers[peer] = peerState
	}
	return apiPeerStates
}

type Stats struct {
	MaxMemoryMB         uint64 `json:"Max Memory (MB)"`
	GitRevision         string `json:"git-revision"`
	ErrorCount          uint64 `json:"Error Count"`
	Uptime              uint64 `json:"uptime"`
	FreeMemoryMB        uint64 `json:"Free Memory (MB)"`
	TotalMemoryMB       uint64 `json:"Total Memory (MB)"`
	Version             string `json:"version"`
	DeployDir           string `json:"deploy-dir"`
	FetchCount          uint64 `json:"Fetch Count"`
	QueryIntervalDelta  int    `json:"Query Interval Delta"`
	IterationCount      uint64 `json:"Iteration Count"`
	Name                string `json:"name"`
	BuildTimestamp      string `json:"buildTimestamp"`
	QueryIntervalTarget int    `json:"Query Interval Target"`
	QueryIntervalActual int    `json:"Query Interval Actual"`
	SlowestCache        string `json:"Slowest Cache"`
	LastQueryInterval   int    `json:"Last Query Interval"`
	Microthreads        int    `json:"Goroutines"`
	LastGC              string `json:"Last Garbage Collection"`
	MemAllocBytes       uint64 `json:"Memory Bytes Allocated"`
	MemTotalBytes       uint64 `json:"Total Bytes Allocated"`
	MemSysBytes         uint64 `json:"System Bytes Allocated"`
}

func getLongestPoll(lastHealthTimes map[enum.CacheName]time.Duration) (enum.CacheName, time.Duration) {
	var longestCache enum.CacheName
	var longestTime time.Duration
	for cache, time := range lastHealthTimes {
		if time > longestTime {
			longestTime = time
			longestCache = cache
		}
	}
	return longestCache, longestTime
}

func getStats(staticAppData StaticAppData, pollingInterval time.Duration, lastHealthTimes map[enum.CacheName]time.Duration, fetchCount uint64, healthIteration uint64, errorCount uint64) ([]byte, error) {
	longestPollCache, longestPollTime := getLongestPoll(lastHealthTimes)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var s Stats
	s.MaxMemoryMB = memStats.TotalAlloc / (1024 * 1024)
	s.GitRevision = staticAppData.GitRevision
	s.ErrorCount = errorCount
	s.Uptime = uint64(time.Since(staticAppData.StartTime) / time.Second)
	s.FreeMemoryMB = staticAppData.FreeMemoryMB
	s.TotalMemoryMB = memStats.Alloc / (1024 * 1024) // TODO rename to "used memory" if/when nothing is using the JSON entry
	s.Version = staticAppData.Version
	s.DeployDir = staticAppData.WorkingDir
	s.FetchCount = fetchCount
	s.SlowestCache = string(longestPollCache)
	s.IterationCount = healthIteration
	s.Name = staticAppData.Name
	s.BuildTimestamp = staticAppData.BuildTimestamp
	s.QueryIntervalTarget = int(pollingInterval / time.Millisecond)
	s.QueryIntervalActual = int(longestPollTime / time.Millisecond)
	s.QueryIntervalDelta = s.QueryIntervalActual - s.QueryIntervalTarget
	s.LastQueryInterval = int(math.Max(float64(s.QueryIntervalActual), float64(s.QueryIntervalTarget)))
	s.Microthreads = runtime.NumGoroutine()
	s.LastGC = time.Unix(0, int64(memStats.LastGC)).String()
	s.MemAllocBytes = memStats.Alloc
	s.MemTotalBytes = memStats.TotalAlloc
	s.MemSysBytes = memStats.Sys

	return json.Marshal(s)
}
