package deliveryservice

import (
	"fmt"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/log"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	dsdata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservicedata"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/http_server"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	"net/url"
	"strconv"
	"time"
)

// TODO remove 'ds' and 'stat' from names

// TODO remove DeliveryService and set type to the map directly, or add other members
type Stats struct {
	DeliveryService map[enum.DeliveryServiceName]dsdata.Stat `json:"deliveryService"`
}

func (a Stats) Copy() Stats {
	b := NewStats()
	for k, v := range a.DeliveryService {
		b.DeliveryService[k] = v.Copy()
	}
	return b
}

func (a Stats) Get(name enum.DeliveryServiceName) (dsdata.StatReadonly, bool) {
	ds, ok := a.DeliveryService[name]
	return ds, ok
}

// TODO rename to just 'New'?
func NewStats() Stats {
	return Stats{DeliveryService: map[enum.DeliveryServiceName]dsdata.Stat{}}
}

func setStaticData(dsStats Stats, dsServers map[enum.DeliveryServiceName][]enum.CacheName) Stats {
	for ds, stat := range dsStats.DeliveryService {
		stat.CommonStats.CachesConfiguredNum.Value = int64(len(dsServers[ds]))
		dsStats.DeliveryService[ds] = stat // TODO consider changing dsStats.DeliveryService[ds] to a pointer so this kind of thing isn't necessary; possibly more performant, as well
	}
	return dsStats
}

func addAvailableData(dsStats Stats, crStates peer.Crstates, serverCachegroups map[enum.CacheName]enum.CacheGroupName, serverDs map[enum.CacheName][]enum.DeliveryServiceName, serverTypes map[enum.CacheName]enum.CacheType, statHistory map[enum.CacheName][]cache.Result) (Stats, error) {
	for cache, available := range crStates.Caches {
		cacheGroup, ok := serverCachegroups[cache]
		if !ok {
			log.Warnf("CreateStats not adding availability data for '%s': not found in Cachegroups\n", cache)
			continue
		}
		deliveryServices, ok := serverDs[cache]
		if !ok {
			log.Warnf("CreateStats not adding availability data for '%s': not found in DeliveryServices\n", cache)
			continue
		}
		cacheType, ok := serverTypes[enum.CacheName(cache)]
		if !ok {
			log.Warnf("CreateStats not adding availability data for '%s': not found in Server Types\n", cache)
			continue
		}

		for _, deliveryService := range deliveryServices {
			if deliveryService == "" {
				log.Errorf("EMPTY addAvailableData DS") // various bugs in other functions can cause this - this will help identify and debug them.
				continue
			}

			stat, ok := dsStats.DeliveryService[enum.DeliveryServiceName(deliveryService)]
			if !ok {
				log.Warnf("CreateStats not adding availability data for '%s': not found in Stats\n", cache)
				continue // TODO log warning? Error?
			}

			if available.IsAvailable {
				// c.IsAvailable.Value
				stat.CommonStats.IsAvailable.Value = true
				stat.CommonStats.CachesAvailableNum.Value++
				cacheGroupStats := stat.CacheGroups[enum.CacheGroupName(cacheGroup)]
				cacheGroupStats.IsAvailable.Value = true
				stat.CacheGroups[enum.CacheGroupName(cacheGroup)] = cacheGroupStats
				stat.TotalStats.IsAvailable.Value = true
				typeStats := stat.Types[cacheType]
				typeStats.IsAvailable.Value = true
				stat.Types[cacheType] = typeStats
			}

			// TODO fix nested ifs
			if results, ok := statHistory[enum.CacheName(cache)]; ok {
				if len(results) < 1 {
					log.Warnf("no results %v %v\n", cache, deliveryService)
				} else {
					result := results[0]
					if result.PrecomputedData.Reporting {
						stat.CommonStats.CachesReporting[enum.CacheName(cache)] = true
					} else {
						log.Debugf("no reporting %v %v\n", cache, deliveryService)
					}
				}
			} else {
				log.Debugf("no result for %v %v\n", cache, deliveryService)
			}

			dsStats.DeliveryService[enum.DeliveryServiceName(deliveryService)] = stat // TODO Necessary? Remove?
		}
	}
	return dsStats, nil
}

type StatsLastKbps struct {
	DeliveryServices map[enum.DeliveryServiceName]StatLastKbps
	Caches           map[enum.CacheName]LastKbpsData
}

func NewStatsLastKbps() StatsLastKbps {
	return StatsLastKbps{DeliveryServices: map[enum.DeliveryServiceName]StatLastKbps{}, Caches: map[enum.CacheName]LastKbpsData{}}
}

func (a StatsLastKbps) Copy() StatsLastKbps {
	b := NewStatsLastKbps()
	for k, v := range a.DeliveryServices {
		b.DeliveryServices[k] = v.Copy()
	}
	for k, v := range a.Caches {
		b.Caches[k] = v
	}
	return b
}

// TODO figure a way to associate this type with StatHTTP, with which its members correspond.
type StatLastKbps struct {
	Caches      map[enum.CacheName]LastKbpsData
	CacheGroups map[enum.CacheGroupName]LastKbpsData
	Type        map[enum.CacheType]LastKbpsData
	Total       LastKbpsData
}

func (a StatLastKbps) Copy() StatLastKbps {
	b := StatLastKbps{
		CacheGroups: map[enum.CacheGroupName]LastKbpsData{},
		Type:        map[enum.CacheType]LastKbpsData{},
		Caches:      map[enum.CacheName]LastKbpsData{},
		Total:       a.Total,
	}
	for k, v := range a.CacheGroups {
		b.CacheGroups[k] = v
	}
	for k, v := range a.Type {
		b.Type[k] = v
	}
	for k, v := range a.Caches {
		b.Caches[k] = v
	}
	return b
}

func newStatLastKbps() StatLastKbps {
	return StatLastKbps{
		CacheGroups: map[enum.CacheGroupName]LastKbpsData{},
		Type:        map[enum.CacheType]LastKbpsData{},
		Caches:      map[enum.CacheName]LastKbpsData{},
	}
}

type LastKbpsData struct {
	Kbps  float64
	Bytes int64
	Time  time.Time
}

const BytesPerKilobit = 125

// addKbps adds Kbps fields to the NewStats, based on the previous out_bytes in the oldStats, and the time difference.
//
// Traffic Server only updates its data every N seconds. So, often we get a new Stats with the same OutBytes as the previous one,
// So, we must record the last changed value, and the time it changed. Then, if the new OutBytes is different from the previous,
// we set the (new - old) / lastChangedTime as the KBPS, and update the recorded LastChangedTime and LastChangedValue
//
// This specifically returns the given dsStats and lastKbpsStats on error, so it's safe to do persistentStats, persistentLastKbpsStats, err = addKbps(...)
// TODO handle ATS byte rolling (when the `out_bytes` overflows back to 0)
// TODO break this function up, it's too big.
func addKbps(statHistory map[enum.CacheName][]cache.Result, dsStats Stats, lastKbpsStats StatsLastKbps, dsStatsTime time.Time, serverCachegroups map[enum.CacheName]enum.CacheGroupName, serverTypes map[enum.CacheName]enum.CacheType) (Stats, StatsLastKbps, error) {
	for dsName, stat := range dsStats.DeliveryService {
		lastKbpsStat, lastKbpsStatExists := lastKbpsStats.DeliveryServices[dsName]
		if !lastKbpsStatExists {
			lastKbpsStat = newStatLastKbps()
		}

		for cacheName, cacheStats := range stat.Caches {
			lastCacheStat, lastCacheStatExists := lastKbpsStat.Caches[cacheName]

			if cacheStats.OutBytes.Value == lastCacheStat.Bytes {
				cacheStats.Kbps.Value = lastCacheStat.Kbps
				stat.Caches[cacheName] = cacheStats
				continue
			}

			if lastCacheStatExists && lastCacheStat.Bytes != 0 {
				lastCacheStat.Kbps = float64(cacheStats.OutBytes.Value-lastCacheStat.Bytes) / BytesPerKilobit / stat.CachesTimeReceived[cacheName].Sub(lastCacheStat.Time).Seconds()
			}
			lastCacheStat.Bytes = cacheStats.OutBytes.Value
			lastCacheStat.Time = stat.CachesTimeReceived[cacheName]
			lastKbpsStat.Caches[cacheName] = lastCacheStat

			cacheStats.Kbps.Value = lastCacheStat.Kbps // TODO determine if necessary
			stat.Caches[cacheName] = cacheStats
		}

		// TODO don't add kbps for caches which didn't respond to their last request
		cacheGroups := map[enum.CacheGroupName]LastKbpsData{}
		cacheTypes := map[enum.CacheType]LastKbpsData{}
		total := LastKbpsData{}
		for cacheName, cacheStats := range lastKbpsStat.Caches {
			if !stat.CommonStats.CachesReporting[cacheName] {
				continue
			}

			cacheGroup, ok := serverCachegroups[cacheName]
			if !ok {
				log.Errorf("addkbps cache %v not in cachegroups\n", cacheName)
			} else {
				c := cacheGroups[cacheGroup]
				c.Kbps += cacheStats.Kbps
				cacheGroups[cacheGroup] = c
			}

			cacheType, ok := serverTypes[cacheName]
			if !ok {
				log.Errorf("addkbps cache %v not in types\n", cacheName)
			} else {
				c := cacheTypes[cacheType]
				c.Kbps += cacheStats.Kbps
				cacheTypes[cacheType] = c
			}

			total.Kbps += cacheStats.Kbps
		}

		for cacheGroup, lastKbpsData := range cacheGroups {
			g := stat.CacheGroups[cacheGroup]
			g.Kbps.Value = lastKbpsData.Kbps
			stat.CacheGroups[cacheGroup] = g
		}

		for cacheType, lastKbpsData := range cacheTypes {
			t := stat.Types[cacheType]
			t.Kbps.Value = lastKbpsData.Kbps
			stat.Types[cacheType] = t
		}

		lastKbpsStat.CacheGroups = cacheGroups
		lastKbpsStat.Type = cacheTypes
		lastKbpsStat.Total = total

		stat.TotalStats.Kbps.Value = total.Kbps

		lastKbpsStats.DeliveryServices[dsName] = lastKbpsStat
		dsStats.DeliveryService[dsName] = stat
	}

	for cacheName, results := range statHistory {
		var result *cache.Result
		for _, r := range results {
			// result.Errors can include stat errors where OutBytes was set correctly, so we look for the first non-zero OutBytes rather than the first errorless result
			// TODO add error classes to PrecomputedData, to distinguish stat errors from HTTP errors?
			if r.PrecomputedData.OutBytes == 0 {
				continue
			}
			result = &r
			break
		}

		if result == nil {
			log.Warnf("addkbps cache %v has no results\n", cacheName)
			continue
		}

		outBytes := result.PrecomputedData.OutBytes

		lastCacheKbpsData, ok := lastKbpsStats.Caches[cacheName]
		if !ok {
			// this means this is the first result for this cache - this is a normal condition
			lastKbpsStats.Caches[cacheName] = LastKbpsData{Time: dsStatsTime, Bytes: outBytes, Kbps: 0}
			continue
		}

		if lastCacheKbpsData.Bytes == outBytes {
			// this means this ATS hasn't updated its byte count yet - this is a normal condition
			continue // don't try to kbps, and importantly don't change the time of the last change, if Traffic Server hasn't updated
		}

		if outBytes == 0 {
			log.Errorf("addkbps %v outbytes zero\n", cacheName)
			continue
		}

		kbps := float64(outBytes-lastCacheKbpsData.Bytes) / BytesPerKilobit / result.Time.Sub(lastCacheKbpsData.Time).Seconds()
		if lastCacheKbpsData.Bytes == 0 {
			kbps = 0
			log.Errorf("addkbps cache %v lastCacheKbpsData.Bytes zero\n", cacheName)
		}
		if kbps < 0 {
			log.Errorf("addkbps negative cache kbps: cache %v kbps %v outBytes %v lastCacheKbpsData.Bytes %v dsStatsTime %v lastCacheKbpsData.Time %v\n", cacheName, kbps, outBytes, lastCacheKbpsData.Bytes, dsStatsTime, lastCacheKbpsData.Time) // this is almost certainly a code bug. The only case this would ever be a data issue, would be if Traffic Server returned fewer bytes than previously.
			kbps = 0
		}

		lastKbpsStats.Caches[cacheName] = LastKbpsData{Time: result.Time, Bytes: outBytes, Kbps: kbps}
	}

	return dsStats, lastKbpsStats, nil
}

func CreateStats(statHistory map[enum.CacheName][]cache.Result, toData todata.TOData, crStates peer.Crstates, lastKbpsStats StatsLastKbps, now time.Time) (Stats, StatsLastKbps, error) {
	start := time.Now()
	dsStats := NewStats()
	for deliveryService, _ := range toData.DeliveryServiceServers {
		if deliveryService == "" {
			log.Errorf("EMPTY CreateStats deliveryService")
			continue
		}
		dsStats.DeliveryService[enum.DeliveryServiceName(deliveryService)] = *dsdata.NewStat()
	}
	dsStats = setStaticData(dsStats, toData.DeliveryServiceServers)
	var err error
	dsStats, err = addAvailableData(dsStats, crStates, toData.ServerCachegroups, toData.ServerDeliveryServices, toData.ServerTypes, statHistory) // TODO move after stat summarisation
	if err != nil {
		return dsStats, lastKbpsStats, fmt.Errorf("Error getting Cache availability data: %v", err)
	}

	for server, history := range statHistory {
		if len(history) < 1 {
			continue // TODO warn?
		}
		cachegroup, ok := toData.ServerCachegroups[server]
		if !ok {
			log.Warnf("server %s has no cachegroup, skipping\n", server)
			continue
		}
		serverType, ok := toData.ServerTypes[enum.CacheName(server)]
		if !ok {
			log.Warnf("server %s not in CRConfig, skipping\n", server)
			continue
		}
		result := history[len(history)-1]

		// TODO check result.PrecomputedData.Errors
		for ds, resultStat := range result.PrecomputedData.DeliveryServiceStats {
			if ds == "" {
				log.Errorf("EMPTY precomputed delivery service")
				continue
			}

			if _, ok := dsStats.DeliveryService[ds]; !ok {
				dsStats.DeliveryService[ds] = resultStat
				continue
			}
			httpDsStat := dsStats.DeliveryService[ds]
			httpDsStat.TotalStats = httpDsStat.TotalStats.Sum(resultStat.TotalStats)
			httpDsStat.CacheGroups[cachegroup] = httpDsStat.CacheGroups[cachegroup].Sum(resultStat.CacheGroups[cachegroup])
			httpDsStat.Types[serverType] = httpDsStat.Types[serverType].Sum(resultStat.Types[serverType])
			httpDsStat.Caches[server] = httpDsStat.Caches[server].Sum(resultStat.Caches[server])
			httpDsStat.CachesTimeReceived[server] = resultStat.CachesTimeReceived[server]
			httpDsStat.CommonStats = dsStats.DeliveryService[ds].CommonStats
			dsStats.DeliveryService[ds] = httpDsStat // TODO determine if necessary
		}
	}

	kbpsStats, kbpsStatsLastKbps, kbpsErr := addKbps(statHistory, dsStats, lastKbpsStats, now, toData.ServerCachegroups, toData.ServerTypes)
	log.Infof("CreateStats took %v\n", time.Since(start))
	return kbpsStats, kbpsStatsLastKbps, kbpsErr
}

func addStatCacheStats(s *dsdata.StatsOld, c dsdata.StatCacheStats, deliveryService enum.DeliveryServiceName, prefix string, t int64, filter dsdata.Filter) *dsdata.StatsOld {
	add := func(name, val string) {
		if filter.UseStat(name) {
			s.DeliveryService[deliveryService][dsdata.StatName(prefix+name)] = []dsdata.StatOld{dsdata.StatOld{Time: t, Value: val}}
		}
	}
	add("out_bytes", strconv.Itoa(int(c.OutBytes.Value)))
	add("isAvailable", fmt.Sprintf("%t", c.IsAvailable.Value))
	add("status_5xx", strconv.Itoa(int(c.Status5xx.Value)))
	add("status_4xx", strconv.Itoa(int(c.Status4xx.Value)))
	add("status_3xx", strconv.Itoa(int(c.Status3xx.Value)))
	add("status_2xx", strconv.Itoa(int(c.Status2xx.Value)))
	add("in_bytes", strconv.Itoa(int(c.InBytes.Value)))
	add("kbps", strconv.Itoa(int(c.Kbps.Value)))
	add("tps_5xx", strconv.Itoa(int(c.Tps5xx.Value)))
	add("tps_4xx", strconv.Itoa(int(c.Tps4xx.Value)))
	add("tps_3xx", strconv.Itoa(int(c.Tps3xx.Value)))
	add("tps_2xx", strconv.Itoa(int(c.Tps2xx.Value)))
	add("error", c.ErrorString.Value)
	add("tps_total", strconv.Itoa(int(c.TpsTotal.Value)))
	return s
}

func addCommonData(s *dsdata.StatsOld, c *dsdata.StatCommon, deliveryService enum.DeliveryServiceName, t int64, filter dsdata.Filter) *dsdata.StatsOld {
	add := func(name, val string) {
		if filter.UseStat(name) {
			s.DeliveryService[deliveryService][dsdata.StatName(name)] = []dsdata.StatOld{dsdata.StatOld{Time: t, Value: val}}
		}
	}
	add("caches-configured", strconv.Itoa(int(c.CachesConfiguredNum.Value)))
	add("caches-reporting", strconv.Itoa(len(c.CachesReporting)))
	add("error-string", strconv.Itoa(len(c.CachesReporting)))
	add("status", c.StatusStr.Value)
	add("isHealthy", fmt.Sprintf("%t", c.IsHealthy.Value))
	add("isAvailable", fmt.Sprintf("%t", c.IsAvailable.Value))
	add("caches-available", strconv.Itoa(int(c.CachesAvailableNum.Value)))
	return s
}

// StatsJSON returns an object formatted as expected to be serialized to JSON and served.
func (dsStats Stats) JSON(filter dsdata.Filter, params url.Values) dsdata.StatsOld {
	now := time.Now().Unix()
	jsonObj := &dsdata.StatsOld{
		DeliveryService: map[enum.DeliveryServiceName]map[dsdata.StatName][]dsdata.StatOld{},
		QueryParams:     http_server.ParametersStr(params),
		DateStr:         http_server.DateStr(time.Now()),
	}

	for deliveryService, stat := range dsStats.DeliveryService {
		if !filter.UseDeliveryService(deliveryService) {
			continue
		}
		jsonObj.DeliveryService[deliveryService] = map[dsdata.StatName][]dsdata.StatOld{}
		jsonObj = addCommonData(jsonObj, &stat.CommonStats, deliveryService, now, filter)
		for cacheGroup, cacheGroupStats := range stat.CacheGroups {
			jsonObj = addStatCacheStats(jsonObj, cacheGroupStats, deliveryService, "location."+string(cacheGroup)+".", now, filter)
		}
		for cacheType, typeStats := range stat.Types {
			jsonObj = addStatCacheStats(jsonObj, typeStats, deliveryService, "type."+cacheType.String()+".", now, filter)
		}
		jsonObj = addStatCacheStats(jsonObj, stat.TotalStats, deliveryService, "total.", now, filter)
	}
	return *jsonObj
}
