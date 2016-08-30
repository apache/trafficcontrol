package deliveryservice

import (
	"fmt"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	dsdata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservicedata"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	"strconv"
	"strings"
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
	for ds, istat := range dsStats.DeliveryService {
		istat.CommonData().CachesConfigured.Value = int64(len(dsServers[string(ds)]))
	}
	return dsStats
}

func addAvailableData(dsStats Stats, crStates peer.Crstates, serverCachegroups map[enum.CacheName]enum.CacheGroupName, serverDs map[string][]string, serverTypes map[enum.CacheName]enum.CacheType, statHistory map[enum.CacheName][]cache.Result) (Stats, error) {
	for cache, available := range crStates.Caches {
		cacheGroup, ok := serverCachegroups[enum.CacheName(cache)]
		if !ok {
			fmt.Printf("WARNING: CreateStats not adding availability data for '%s': not found in Cachegroups\n", cache)
			continue
		}
		deliveryServices, ok := serverDs[cache]
		if !ok {
			fmt.Printf("WARNING: CreateStats not adding availability data for '%s': not found in DeliveryServices\n", cache)
			continue
		}
		cacheType, ok := serverTypes[enum.CacheName(cache)]
		if !ok {
			fmt.Printf("WARNING: CreateStats not adding availability data for '%s': not found in Server Types\n", cache)
			continue
		}

		for _, deliveryService := range deliveryServices {
			iStat, ok := dsStats.DeliveryService[enum.DeliveryServiceName(deliveryService)]
			if !ok || iStat == nil {
				fmt.Printf("WARNING: CreateStats not adding availability data for '%s': not found in Stats\n", cache)
				continue // TODO log warning? Error?
			}

			if available.IsAvailable {
				// c.IsAvailable.Value
				iStat.CommonData().IsAvailable.Value = true
				iStat.CommonData().CachesAvailable.Value++
				if stat, ok := iStat.(*dsdata.StatHTTP); ok {
					cacheGroupStats := stat.CacheGroups[enum.CacheGroupName(cacheGroup)]
					cacheGroupStats.IsAvailable.Value = true
					stat.CacheGroups[enum.CacheGroupName(cacheGroup)] = cacheGroupStats
					stat.Total.IsAvailable.Value = true
					typeStats := stat.Type[cacheType]
					typeStats.IsAvailable.Value = true
					stat.Type[cacheType] = typeStats
				} else if _, ok := iStat.(*dsdata.StatDNS); ok {
				} else {
					return dsStats, fmt.Errorf("Unknown stat type for Delivery Service '%s': %v", deliveryService, iStat)
				}
			}

			// TODO fix nested ifs
			if results, ok := statHistory[enum.CacheName(cache)]; ok {
				if len(results) < 1 {
					fmt.Printf("WARNING no results %v %v\n", cache, deliveryService)
				} else {
					result := results[0]
					if result.PrecomputedData.Reporting {
						iStat.CommonData().CachesReporting[enum.CacheName(cache)] = true
					} else {
						fmt.Printf("DEBUG no reporting %v %v\n", cache, deliveryService)
					}
				}
			} else {
				fmt.Printf("DEBUG no result for %v %v\n", cache, deliveryService)
			}

			dsStats.DeliveryService[enum.DeliveryServiceName(deliveryService)] = iStat // TODO Necessary? Remove?
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
func addKbps(statHistory map[enum.CacheName][]cache.Result, dsStats Stats, lastKbpsStats StatsLastKbps, dsStatsTime time.Time) (Stats, StatsLastKbps, error) {
	for dsName, iStat := range dsStats.DeliveryService {
		if _, ok := iStat.(*dsdata.StatDNS); ok {
			continue
		}
		if _, ok := iStat.(*dsdata.StatHTTP); !ok {
			fmt.Printf("WARNING: addKbps got unknown stat type %T\n", iStat)
			continue
		}
		stat := iStat.(*dsdata.StatHTTP)

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

			if lastKbpsStatExists {
				cacheStats.Kbps.Value = float64(cacheStats.OutBytes.Value-lastKbpsData.Bytes) / dsStatsTime.Sub(lastKbpsData.Time).Seconds()
			}

			lastKbpsStat.CacheGroups[cgName] = LastKbpsData{Time: dsStatsTime, Bytes: cacheStats.OutBytes.Value, Kbps: cacheStats.Kbps.Value}
			stat.CacheGroups[cgName] = cacheStats
		}

		for cacheType, cacheStats := range stat.Type {
			lastKbpsData, _ := lastKbpsStat.Type[cacheType]
			if cacheStats.OutBytes.Value == lastKbpsData.Bytes {
				if cacheStats.OutBytes.Value == lastKbpsData.Bytes {
					cacheStats.Kbps.Value = lastKbpsData.Kbps
					stat.Type[cacheType] = cacheStats
					continue
				}
				if lastKbpsStatExists {
					cacheStats.Kbps.Value = float64(cacheStats.OutBytes.Value-lastKbpsData.Bytes) / dsStatsTime.Sub(lastKbpsData.Time).Seconds()
				}
				lastKbpsStat.Type[cacheType] = LastKbpsData{Time: dsStatsTime, Bytes: cacheStats.OutBytes.Value, Kbps: cacheStats.Kbps.Value}
				stat.Type[cacheType] = cacheStats
			}
		}
		if lastKbpsStatExists {
			stat.Total.Kbps.Value = float64(stat.Total.OutBytes.Value-lastKbpsStat.Total.Bytes) / dsStatsTime.Sub(lastKbpsStat.Total.Time).Seconds()
		} else {
			stat.Total.Kbps.Value = lastKbpsStat.Total.Kbps
		}
		lastKbpsStat.Total = LastKbpsData{Time: dsStatsTime, Bytes: stat.Total.OutBytes.Value, Kbps: stat.Total.Kbps.Value}

		lastKbpsStats.DeliveryServices[dsName] = lastKbpsStat
	}

	for cacheName, results := range statHistory { // map[enum.CacheName]int64
		if len(results) < 1 {
			continue // TODO warn?
		}
		result := results[0]
		outBytes := result.PrecomputedData.OutBytes

		lastCacheKbpsData, ok := lastKbpsStats.Caches[cacheName]
		if !ok {
			lastKbpsStats.Caches[cacheName] = LastKbpsData{Time: dsStatsTime, Bytes: outBytes, Kbps: 0}
			continue
		}

		if lastCacheKbpsData.Bytes == outBytes {
			continue // don't try to kbps, and importantly don't change the time of the last change, if Traffic Server hasn't updated
		}

		if outBytes == 0 {
			fmt.Printf("ERROR adding kbps %v outbytes zero\n", cacheName)
			continue
		}

		kbps := float64(outBytes-lastCacheKbpsData.Bytes) / result.Time.Sub(lastCacheKbpsData.Time).Seconds() / BytesPerKbps
		if lastCacheKbpsData.Bytes == 0 {
			kbps = 0
			fmt.Printf("ERROR adding kbps %v lastCacheKbpsData.Bytes zero\n", cacheName)
		}
		if kbps < 0 {
			kbps = 0
			// TODO figure out what to do. Print error. Explode. Definitely don't set kbps negative.
			fmt.Printf("ERROR negative kbps: %v kbps %v outBytes %v lastCacheKbpsData.Bytes %v dsStatsTime %v lastCacheKbpsData.Time %v\n", cacheName, kbps, outBytes, lastCacheKbpsData.Bytes, dsStatsTime, lastCacheKbpsData.Time)
		}

		lastKbpsStats.Caches[cacheName] = LastKbpsData{Time: result.Time, Bytes: outBytes, Kbps: kbps}
	}

	return dsStats, lastKbpsStats, nil
}

func CreateStats(statHistory map[enum.CacheName][]cache.Result, toData todata.TOData, crStates peer.Crstates, lastKbpsStats StatsLastKbps, now time.Time) (Stats, StatsLastKbps, error) {
	start := time.Now()
	dsStats := NewStats()
	for deliveryService, _ := range toData.DeliveryServiceServers {
		dsType, ok := toData.DeliveryServiceTypes[deliveryService]
		if !ok {
			return Stats{}, lastKbpsStats, fmt.Errorf("deliveryservice %s missing type", deliveryService)
		}
		if dsType == enum.DSTypeHTTP {
			dsStats.DeliveryService[enum.DeliveryServiceName(deliveryService)] = dsdata.NewStatHTTP()
		} else if dsType == enum.DSTypeDNS {
			dsStats.DeliveryService[enum.DeliveryServiceName(deliveryService)] = dsdata.NewStatDNS()
		} else {
			return Stats{}, lastKbpsStats, fmt.Errorf("unknown type for '%s': %v", deliveryService, dsType)
		}
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
			fmt.Printf("WARNING server %s has no cachegroup, skipping\n", server)
			continue
		}
		serverType, ok := toData.ServerTypes[enum.CacheName(server)]
		if !ok {
			fmt.Printf("WARNING server %s not in CRConfig, skipping\n", server)
			continue
		}
		result := history[len(history)-1]

		for ds, stat := range result.PrecomputedData.DeliveryServiceStats {
			switch stat.(type) {
			case *dsdata.StatHTTP:
				resultHttpStat := stat.(*dsdata.StatHTTP)
				if _, ok := dsStats.DeliveryService[ds]; !ok {
					dsStats.DeliveryService[ds] = resultHttpStat
					continue
				}
				httpDsStat, ok := dsStats.DeliveryService[ds].(*dsdata.StatHTTP)
				if !ok {
					fmt.Printf("WARNING precomputed stat does not match dsStats\n")
					continue
				}
				httpDsStat.Total = httpDsStat.Total.Sum(resultHttpStat.Total)
				httpDsStat.CacheGroups[cachegroup] = httpDsStat.CacheGroups[cachegroup].Sum(resultHttpStat.CacheGroups[cachegroup])
				httpDsStat.Type[serverType] = httpDsStat.Type[serverType].Sum(resultHttpStat.Type[serverType])
				dsStats.DeliveryService[ds] = httpDsStat // TODO determine if necessary
			case *dsdata.StatDNS:
				resultDnsStat := stat.(*dsdata.StatDNS)
				dsStatsDnsStat, ok := dsStats.DeliveryService[ds].(*dsdata.StatDNS)
				if !ok {
					fmt.Printf("WARNING precomputed DNS stat does not match dsStats\n")
					continue
				}
				dsStatsDnsStat.Sum(resultDnsStat)
				dsStats.DeliveryService[ds] = dsStatsDnsStat // TODO determine if necessary
			}
		}
	}

	kbpsStats, kbpsStatsLastKbps, kbpsErr := addKbps(statHistory, dsStats, lastKbpsStats, now)
	fmt.Printf("CreateStats took %v\n", time.Since(start))
	return kbpsStats, kbpsStatsLastKbps, kbpsErr
}

// processStat and its subsidiary functions act as a State Machine, flowing the stat thru states for each "." component of the stat name
// TODO fix this being crazy slow. THIS IS THE BOTTLENECK
func processStat(dsStats *Stats, dsRegexes todata.Regexes, dsTypes map[string]enum.DSType, cachegroup enum.CacheGroupName, server enum.CacheName, serverType enum.CacheType, stat string, value interface{}) (enum.DeliveryServiceName, dsdata.Stat, error) {
	parts := strings.Split(stat, ".")
	if len(parts) < 1 {
		return "", nil, fmt.Errorf("stat has no initial part")
	}

	switch parts[0] {
	case "plugin":
		return processStatPlugin(dsStats, dsRegexes, dsTypes, cachegroup, server, serverType, stat, parts[1:], value)
	case "proxy":
		return "", nil, dsdata.ErrNotProcessedStat
	default:
		return "", nil, fmt.Errorf("stat has unknown initial part '%s'", parts[0])
	}
}

func processStatPlugin(dsStats *Stats, dsRegexes todata.Regexes, dsTypes map[string]enum.DSType, cachegroup enum.CacheGroupName, server enum.CacheName, serverType enum.CacheType, stat string, statParts []string, value interface{}) (enum.DeliveryServiceName, dsdata.Stat, error) {
	if len(statParts) < 1 {
		return "", nil, fmt.Errorf("stat has no plugin part")
	}
	switch statParts[0] {
	case "remap_stats":
		return processStatPluginRemapStats(dsStats, dsRegexes, dsTypes, cachegroup, server, serverType, stat, statParts[1:], value)
	default:
		return "", nil, fmt.Errorf("stat has unknown plugin part '%s'", statParts[0])
	}
}

func processStatPluginRemapStats(dsStats *Stats, dsRegexes todata.Regexes, dsTypes map[string]enum.DSType, cachegroup enum.CacheGroupName, server enum.CacheName, serverType enum.CacheType, stat string, statParts []string, value interface{}) (enum.DeliveryServiceName, dsdata.Stat, error) {
	if len(statParts) < 2 {
		return "", nil, fmt.Errorf("stat has no remap_stats deliveryservice and name parts")
	}

	fqdn := strings.Join(statParts[:len(statParts)-1], ".")
	statName := statParts[len(statParts)-1]
	ds, ok := dsRegexes.DeliveryService(fqdn)
	if !ok {
		return ds, nil, fmt.Errorf("%s matched no delivery service", fqdn)
	}

	if _, ok := dsTypes[string(ds)]; !ok {
		return ds, nil, fmt.Errorf("delivery service %s not found in types map", ds)
	}

	addedStat, err := addStat(dsStats.DeliveryService[ds], statName, value, string(ds), server, serverType, cachegroup, dsTypes)
	if err != nil {
		return ds, nil, err
	}
	return ds, addedStat, nil
}

func addStat(iStat dsdata.Stat, name string, val interface{}, ds string, server enum.CacheName, serverType enum.CacheType, cachegroup enum.CacheGroupName, dsTypes map[string]enum.DSType) (dsdata.Stat, error) {
	if iStat == nil {
		return iStat, fmt.Errorf("addStat given nil stat for %s", ds)
	}

	var common *dsdata.StatCommon
	common = iStat.CommonData()

	if name == "error_string" {
		valStr, ok := val.(string)
		if !ok {
			return iStat, fmt.Errorf("stat '%s' value expected string actual '%v' type %T", name, val, val)
		}
		common.ErrorString.Value += valStr + "; " // TODO figure out what the delimiter should be
	}
	common.Status.Value = "REPORTED" // TODO fix?

	if stat, ok := iStat.(*dsdata.StatHTTP); ok {
		newCachegroupStat, err := addCacheStat(stat.CacheGroups[cachegroup], name, val)
		if err != nil {
			return stat, err
		}
		stat.CacheGroups[enum.CacheGroupName(cachegroup)] = newCachegroupStat

		newTypeStat, err := addCacheStat(stat.Type[serverType], name, val)
		if err != nil {
			return stat, err
		}
		stat.Type[serverType] = newTypeStat

		newTotal, err := addCacheStat(stat.Total, name, val)
		if err != nil {
			return stat, err
		}
		stat.Total = newTotal

		return stat, nil
	}
	if stat, ok := iStat.(*dsdata.StatDNS); ok {
		// TODO handle DNS DS stats
		return stat, nil
	}
	return iStat, fmt.Errorf("delivery service %s type is invalid", iStat)
}

// addCacheStat adds the given stat to the existing stat. Note this adds, it doesn't overwrite. Numbers are summed, strings are concatenated.
// TODO make this less duplicate code somehow.
func addCacheStat(stat dsdata.StatCacheStats, name string, val interface{}) (dsdata.StatCacheStats, error) {
	switch name {
	case "status_2xx":
		v, ok := val.(float64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status2xx.Value += int64(v)
	case "status_3xx":
		v, ok := val.(float64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status3xx.Value += int64(v)
	case "status_4xx":
		v, ok := val.(float64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status4xx.Value += int64(v)
	case "status_5xx":
		v, ok := val.(float64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status5xx.Value += int64(v)
	case "out_bytes":
		v, ok := val.(float64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.OutBytes.Value += int64(v)
	case "is_available":
		fmt.Println("DEBUGa got is_available")
		v, ok := val.(bool)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected bool actual '%v' type %T", name, val, val)
		}
		if v {
			stat.IsAvailable.Value = true
		}
	case "in_bytes":
		v, ok := val.(float64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.InBytes.Value += v
	case "tps_2xx":
		v, ok := val.(int64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps2xx.Value += v
	case "tps_3xx":
		v, ok := val.(int64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps3xx.Value += v
	case "tps_4xx":
		v, ok := val.(int64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps4xx.Value += v
	case "tps_5xx":
		v, ok := val.(int64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps5xx.Value += v
	case "error_string":
		v, ok := val.(string)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected string actual '%v' type %T", name, val, val)
		}
		stat.ErrorString.Value += v + ", "
	case "tps_total":
		v, ok := val.(int64)
		if !ok {
			return stat, fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.TpsTotal.Value += v
	case "status_unknown":
		return stat, dsdata.ErrNotProcessedStat
	default:
		return stat, fmt.Errorf("unknown stat '%s'", name)
	}
	return stat, nil
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

	for deliveryService, dsStat := range dsStats.DeliveryService {
		jsonObj.DeliveryService[deliveryService] = map[StatName][]StatOld{}
		jsonObj = addCommonData(jsonObj, dsStat.CommonData(), deliveryService, now)
		if stat, ok := dsStat.(*dsdata.StatHTTP); ok {
			for cacheGroup, cacheGroupStats := range stat.CacheGroups {
				jsonObj = addStatCacheStats(jsonObj, cacheGroupStats, deliveryService, string("location."+cacheGroup), now)
			}
			for cacheType, typeStats := range stat.Type {
				jsonObj = addStatCacheStats(jsonObj, typeStats, deliveryService, "type."+cacheType.String(), now)
			}
			jsonObj = addStatCacheStats(jsonObj, stat.Total, deliveryService, "total", now)
		}
	}
	return *jsonObj
}
