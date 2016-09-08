package deliveryservice

import (
	"fmt"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	dsdata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservicedata"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/log"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
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

// TODO rename to just 'New'?
func NewStats() Stats {
	return Stats{DeliveryService: map[enum.DeliveryServiceName]dsdata.Stat{}}
}

func setStaticData(dsStats Stats, dsServers map[string][]string) Stats {
	for ds, stat := range dsStats.DeliveryService {
		stat.Common.CachesConfigured.Value = int64(len(dsServers[string(ds)]))
		dsStats.DeliveryService[ds] = stat // TODO consider changing dsStats.DeliveryService[ds] to a pointer so this kind of thing isn't necessary; possibly more performant, as well
	}
	return dsStats
}

func addAvailableData(dsStats Stats, crStates peer.Crstates, serverCachegroups map[enum.CacheName]enum.CacheGroupName, serverDs map[string][]string, serverTypes map[enum.CacheName]enum.CacheType, statHistory map[enum.CacheName][]cache.Result) (Stats, error) {
	for cache, available := range crStates.Caches {
		cacheGroup, ok := serverCachegroups[enum.CacheName(cache)]
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
				stat.Common.IsAvailable.Value = true
				stat.Common.CachesAvailable.Value++
				cacheGroupStats := stat.CacheGroups[enum.CacheGroupName(cacheGroup)]
				cacheGroupStats.IsAvailable.Value = true
				stat.CacheGroups[enum.CacheGroupName(cacheGroup)] = cacheGroupStats
				stat.Total.IsAvailable.Value = true
				typeStats := stat.Type[cacheType]
				typeStats.IsAvailable.Value = true
				stat.Type[cacheType] = typeStats
			}

			// TODO fix nested ifs
			if results, ok := statHistory[enum.CacheName(cache)]; ok {
				if len(results) < 1 {
					log.Warnf("no results %v %v\n", cache, deliveryService)
				} else {
					result := results[0]
					if result.PrecomputedData.Reporting {
						stat.Common.CachesReporting[enum.CacheName(cache)] = true
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
	CacheGroups map[enum.CacheGroupName]LastKbpsData
	Type        map[enum.CacheType]LastKbpsData
	Total       LastKbpsData
}

func (a StatLastKbps) Copy() StatLastKbps {
	b := StatLastKbps{CacheGroups: map[enum.CacheGroupName]LastKbpsData{}, Type: map[enum.CacheType]LastKbpsData{}, Total: a.Total}
	for k, v := range a.CacheGroups {
		b.CacheGroups[k] = v
	}
	for k, v := range a.Type {
		b.Type[k] = v
	}
	return b
}

func newStatLastKbps() StatLastKbps {
	return StatLastKbps{CacheGroups: map[enum.CacheGroupName]LastKbpsData{}, Type: map[enum.CacheType]LastKbpsData{}}
}

type LastKbpsData struct {
	Kbps  float64
	Bytes int64
	Time  time.Time
}

const BytesPerKbps = 1024

// addKbps adds Kbps fields to the NewStats, based on the previous out_bytes in the oldStats, and the time difference.
//
// Traffic Server only updates its data every N seconds. So, often we get a new Stats with the same OutBytes as the previous one,
// So, we must record the last changed value, and the time it changed. Then, if the new OutBytes is different from the previous,
// we set the (new - old) / lastChangedTime as the KBPS, and update the recorded LastChangedTime and LastChangedValue
//
// This specifically returns the given dsStats and lastKbpsStats on error, so it's safe to do persistentStats, persistentLastKbpsStats, err = addKbps(...)
// TODO handle ATS byte rolling (when the `out_bytes` overflows back to 0)
func addKbps(statHistory map[enum.CacheName][]cache.Result, dsStats Stats, lastKbpsStats StatsLastKbps, dsStatsTime time.Time) (Stats, StatsLastKbps, error) {
	for dsName, stat := range dsStats.DeliveryService {
		lastKbpsStat, lastKbpsStatExists := lastKbpsStats.DeliveryServices[dsName]
		if !lastKbpsStatExists {
			lastKbpsStat = newStatLastKbps()
		}

		for cgName, cacheStats := range stat.CacheGroups {
			lastKbpsData, _ := lastKbpsStat.CacheGroups[cgName]

			if cacheStats.OutBytes.Value == lastKbpsData.Bytes {
				cacheStats.Kbps.Value = lastKbpsData.Kbps
				stat.CacheGroups[cgName] = cacheStats
				continue
			}

			if lastKbpsStatExists && lastKbpsData.Bytes != 0 {
				cacheStats.Kbps.Value = float64(cacheStats.OutBytes.Value-lastKbpsData.Bytes) / dsStatsTime.Sub(lastKbpsData.Time).Seconds()
			}

			if cacheStats.Kbps.Value < 0 {
				cacheStats.Kbps.Value = 0
				log.Errorf("addkbps negative cachegroup cacheStats.Kbps.Value: '%v' '%v' %v - %v / %v\n", dsName, cgName, cacheStats.OutBytes.Value, lastKbpsData.Bytes, dsStatsTime.Sub(lastKbpsData.Time).Seconds())
			}

			lastKbpsStat.CacheGroups[cgName] = LastKbpsData{Time: dsStatsTime, Bytes: cacheStats.OutBytes.Value, Kbps: cacheStats.Kbps.Value}
			stat.CacheGroups[cgName] = cacheStats
		}

		for cacheType, cacheStats := range stat.Type {
			lastKbpsData, _ := lastKbpsStat.Type[cacheType]
			if cacheStats.OutBytes.Value == lastKbpsData.Bytes {
				if cacheStats.OutBytes.Value == lastKbpsData.Bytes {
					if lastKbpsData.Kbps < 0 {
						log.Errorf("addkbps negative cachetype cacheStats.Kbps.Value!\n")
						lastKbpsData.Kbps = 0
					}
					cacheStats.Kbps.Value = lastKbpsData.Kbps
					stat.Type[cacheType] = cacheStats
					continue
				}
				if lastKbpsStatExists && lastKbpsData.Bytes != 0 {
					cacheStats.Kbps.Value = float64(cacheStats.OutBytes.Value-lastKbpsData.Bytes) / dsStatsTime.Sub(lastKbpsData.Time).Seconds()
				}
				if cacheStats.Kbps.Value < 0 {
					log.Errorf("addkbps negative cachetype cacheStats.Kbps.Value.\n")
					cacheStats.Kbps.Value = 0
				}
				lastKbpsStat.Type[cacheType] = LastKbpsData{Time: dsStatsTime, Bytes: cacheStats.OutBytes.Value, Kbps: cacheStats.Kbps.Value}
				stat.Type[cacheType] = cacheStats
			}
		}

		totalChanged := lastKbpsStat.Total.Bytes != stat.Total.OutBytes.Value
		if lastKbpsStatExists && lastKbpsStat.Total.Bytes != 0 && totalChanged {
			stat.Total.Kbps.Value = float64(stat.Total.OutBytes.Value-lastKbpsStat.Total.Bytes) / dsStatsTime.Sub(lastKbpsStat.Total.Time).Seconds() / BytesPerKbps
			if stat.Total.Kbps.Value < 0 {
				stat.Total.Kbps.Value = 0
				log.Errorf("addkbps negative stat.Total.Kbps.Value! Deliveryservice '%v' %v - %v / %v\n", dsName, stat.Total.OutBytes.Value, lastKbpsStat.Total.Bytes, dsStatsTime.Sub(lastKbpsStat.Total.Time).Seconds())
			}
		} else {
			stat.Total.Kbps.Value = lastKbpsStat.Total.Kbps
		}

		if totalChanged {
			lastKbpsStat.Total = LastKbpsData{Time: dsStatsTime, Bytes: stat.Total.OutBytes.Value, Kbps: stat.Total.Kbps.Value}
		}

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

		kbps := float64(outBytes-lastCacheKbpsData.Bytes) / result.Time.Sub(lastCacheKbpsData.Time).Seconds() / BytesPerKbps
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
			httpDsStat.Total = httpDsStat.Total.Sum(resultStat.Total)
			httpDsStat.CacheGroups[cachegroup] = httpDsStat.CacheGroups[cachegroup].Sum(resultStat.CacheGroups[cachegroup])
			httpDsStat.Type[serverType] = httpDsStat.Type[serverType].Sum(resultStat.Type[serverType])
			dsStats.DeliveryService[ds] = httpDsStat // TODO determine if necessary
		}
	}

	kbpsStats, kbpsStatsLastKbps, kbpsErr := addKbps(statHistory, dsStats, lastKbpsStats, now)
	log.Infof("CreateStats took %v\n", time.Since(start))
	return kbpsStats, kbpsStatsLastKbps, kbpsErr
}

type StatName string
type StatOld struct {
	Time  int64  `json:"time"`
	Value string `json:"value"`
	Span  int    `json:"span,omitempty"`  // TODO set? remove?
	Index int    `json:"index,omitempty"` // TODO set? remove?
}
type StatsOld struct {
	DeliveryService map[enum.DeliveryServiceName]map[StatName][]StatOld `json:"deliveryService"`
}

func addStatCacheStats(s *StatsOld, c dsdata.StatCacheStats, deliveryService enum.DeliveryServiceName, prefix string, t int64) *StatsOld {
	s.DeliveryService[deliveryService][StatName(prefix+".out_bytes")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.OutBytes.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".isAvailable")] = []StatOld{StatOld{Time: t, Value: fmt.Sprintf("%t", c.IsAvailable.Value)}}
	s.DeliveryService[deliveryService][StatName(prefix+".status_5xx")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Status5xx.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".status_4xx")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Status4xx.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".status_3xx")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Status3xx.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".status_2xx")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Status2xx.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".in_bytes")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.InBytes.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".kbps")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Kbps.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".tps_5xx")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Tps5xx.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".tps_4xx")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Tps4xx.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".tps_3xx")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Tps3xx.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".tps_2xx")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.Tps2xx.Value))}}
	s.DeliveryService[deliveryService][StatName(prefix+".error-string")] = []StatOld{StatOld{Time: t, Value: c.ErrorString.Value}}
	s.DeliveryService[deliveryService][StatName(prefix+".tps_total")] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.TpsTotal.Value))}}
	return s
}

func addCommonData(s *StatsOld, c *dsdata.StatCommon, deliveryService enum.DeliveryServiceName, t int64) *StatsOld {
	s.DeliveryService[deliveryService]["caches-configured"] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.CachesConfigured.Value))}}
	s.DeliveryService[deliveryService]["caches-reporting"] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(len(c.CachesReporting))}}
	s.DeliveryService[deliveryService]["error-string"] = []StatOld{StatOld{Time: t, Value: c.ErrorString.Value}}
	s.DeliveryService[deliveryService]["status"] = []StatOld{StatOld{Time: t, Value: c.Status.Value}}
	s.DeliveryService[deliveryService]["isHealthy"] = []StatOld{StatOld{Time: t, Value: fmt.Sprintf("%t", c.IsHealthy.Value)}}
	s.DeliveryService[deliveryService]["isAvailable"] = []StatOld{StatOld{Time: t, Value: fmt.Sprintf("%t", c.IsAvailable.Value)}}
	s.DeliveryService[deliveryService]["caches-available"] = []StatOld{StatOld{Time: t, Value: strconv.Itoa(int(c.CachesAvailable.Value))}}
	return s
}

// StatsJSON returns an object formatted as expected to be serialized to JSON and served.
func StatsJSON(dsStats Stats) StatsOld {
	now := time.Now().Unix()
	jsonObj := &StatsOld{DeliveryService: map[enum.DeliveryServiceName]map[StatName][]StatOld{}}

	for deliveryService, stat := range dsStats.DeliveryService {
		jsonObj.DeliveryService[deliveryService] = map[StatName][]StatOld{}
		jsonObj = addCommonData(jsonObj, &stat.Common, deliveryService, now)
		for cacheGroup, cacheGroupStats := range stat.CacheGroups {
			jsonObj = addStatCacheStats(jsonObj, cacheGroupStats, deliveryService, string("location."+cacheGroup), now)
		}
		for cacheType, typeStats := range stat.Type {
			jsonObj = addStatCacheStats(jsonObj, typeStats, deliveryService, "type."+cacheType.String(), now)
		}
		jsonObj = addStatCacheStats(jsonObj, stat.Total, deliveryService, "total", now)
	}
	return *jsonObj
}
