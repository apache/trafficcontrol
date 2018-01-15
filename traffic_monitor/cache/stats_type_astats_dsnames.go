package cache

// stats_type_astats_dsnames is a Stat format similar to the default Astats, but with Fully Qualified Domain Names replaced with Delivery Service names (xml_id)
// That is, stat names are of the form: `"plugin.remap_stats.delivery-service-name.stat-name"`

import(
	"strconv"
	"strings"
	"fmt"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/dsdata"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/todata"
)

func init() {
 	AddStatsType("astats-dsnames", astatsParse, astatsdsnamesPrecompute)
}

func astatsdsnamesPrecompute(cache tc.CacheName, toData todata.TOData, rawStats map[string]interface{}, system AstatsSystem) PrecomputedData {
	stats := map[tc.DeliveryServiceName]dsdata.Stat{}
	precomputed := PrecomputedData{}
	var err error
	if precomputed.OutBytes, err = astatsdsnamesOutBytes(system.ProcNetDev, system.InfName); err != nil {
		precomputed.OutBytes = 0
		log.Errorf("astatsdsnamesPrecompute %s handle precomputing outbytes '%v'\n", cache, err)
	}

	kbpsInMbps := int64(1000)
	precomputed.MaxKbps = int64(system.InfSpeed) * kbpsInMbps

	for stat, value := range rawStats {
		var err error
		stats, err = astatsdsnamesProcessStat(cache, stats, toData, stat, value)
		if err != nil && err != dsdata.ErrNotProcessedStat {
			log.Infof("precomputing cache %v stat %v value %v error %v", cache, stat, value, err)
			precomputed.Errors = append(precomputed.Errors, err)
		}
	}
	precomputed.DeliveryServiceStats = stats
	return precomputed
}

// astatsdsnamesProcessStat and its subsidiary functions act as a State Machine, flowing the stat thru states for each "." component of the stat name
func astatsdsnamesProcessStat(server tc.CacheName, stats map[tc.DeliveryServiceName]dsdata.Stat, toData todata.TOData, stat string, value interface{}) (map[tc.DeliveryServiceName]dsdata.Stat, error) {
	parts := strings.Split(stat, ".")
	if len(parts) < 1 {
		return stats, fmt.Errorf("stat has no initial part")
	}

	switch parts[0] {
	case "plugin":
		return astatsdsnamesProcessStatPlugin(server, stats, toData, stat, parts[1:], value)
	case "proxy":
		return stats, dsdata.ErrNotProcessedStat
	case "server":
		return stats, dsdata.ErrNotProcessedStat
	default:
		return stats, fmt.Errorf("stat '%s' has unknown initial part '%s'", stat, parts[0])
	}
}

func astatsdsnamesProcessStatPlugin(server tc.CacheName, stats map[tc.DeliveryServiceName]dsdata.Stat, toData todata.TOData, stat string, statParts []string, value interface{}) (map[tc.DeliveryServiceName]dsdata.Stat, error) {
	if len(statParts) < 1 {
		return stats, fmt.Errorf("stat has no plugin part")
	}
	switch statParts[0] {
	case "remap_stats":
		return astatsdsnamesProcessStatPluginRemapStats(server, stats, toData, stat, statParts[1:], value)
	default:
		return stats, fmt.Errorf("stat has unknown plugin part '%s'", statParts[0])
	}
}

func astatsdsnamesProcessStatPluginRemapStats(server tc.CacheName, stats map[tc.DeliveryServiceName]dsdata.Stat, toData todata.TOData, stat string, statParts []string, value interface{}) (map[tc.DeliveryServiceName]dsdata.Stat, error) {
	if len(statParts) < 3 {
		return stats, fmt.Errorf("stat has no remap_stats deliveryservice and name parts")
	}

	ds := tc.DeliveryServiceName(statParts[0])
	statName := statParts[len(statParts)-1]

	if _, ok := toData.DeliveryServiceTypes[ds]; !ok {
		return stats, fmt.Errorf("no delivery service match for name '%v' stat '%v'\n", ds, statName)
	}

	dsStat, ok := stats[ds]
	if !ok {
		newStat := dsdata.NewStat()
		dsStat = *newStat
	}

	if err := astatsdstypesAddCacheStat(&dsStat.TotalStats, statName, value); err != nil {
		return stats, err
	}

	cachegroup, ok := toData.ServerCachegroups[server]
	if !ok {
		return stats, fmt.Errorf("server missing from TOData.ServerCachegroups")
	}
	dsStat.CacheGroups[cachegroup] = dsStat.TotalStats

	cacheType, ok := toData.ServerTypes[server]
	if !ok {
		return stats, fmt.Errorf("server missing from TOData.ServerTypes")
	}
	dsStat.Types[cacheType] = dsStat.TotalStats

	dsStat.Caches[server] = dsStat.TotalStats

	stats[ds] = dsStat
	return stats, nil
}

// astatsdsnamesOutBytes takes the proc.net.dev string, and the interface name, and returns the bytes field
// NOTE this is superficially duplicated from astatsOutBytes, but they are conceptually different, because the `astats` format changing should not necessarily affect the `astats-dstypes` format. The MUST be kept separate, and code between them MUST NOT be de-duplicated.
func astatsdsnamesOutBytes(procNetDev, iface string) (int64, error) {
	if procNetDev == "" {
		return 0, fmt.Errorf("procNetDev empty")
	}
	if iface == "" {
		return 0, fmt.Errorf("iface empty")
	}
	ifacePos := strings.Index(procNetDev, iface)
	if ifacePos == -1 {
		return 0, fmt.Errorf("interface '%s' not found in proc.net.dev '%s'", iface, procNetDev)
	}

	procNetDevIfaceBytes := procNetDev[ifacePos+len(iface)+1:]
	procNetDevIfaceBytesArr := strings.Fields(procNetDevIfaceBytes) // TODO test
	if len(procNetDevIfaceBytesArr) < 10 {
		return 0, fmt.Errorf("proc.net.dev iface '%v' unknown format '%s'", iface, procNetDev)
	}
	procNetDevIfaceBytes = procNetDevIfaceBytesArr[8]

	return strconv.ParseInt(procNetDevIfaceBytes, 10, 64)
}

// astatsdstypesAddCacheStat adds the given stat to the existing stat. Note this adds, it doesn't overwrite. Numbers are summed, strings are concatenated.
func astatsdstypesAddCacheStat(stat *dsdata.StatCacheStats, name string, val interface{}) error {
	// TODO make this less duplicate code somehow.
	// NOTE this is superficially duplicated from astatsAddCacheStat, but they are conceptually different, because the `astats` format changing should not necessarily affect the `astats-dstypes` format. The MUST be kept separate, and code between them MUST NOT be de-duplicated.
	switch name {
	case "status_2xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status2xx.Value += int64(v)
	case "status_3xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status3xx.Value += int64(v)
	case "status_4xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status4xx.Value += int64(v)
	case "status_5xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status5xx.Value += int64(v)
	case "out_bytes":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.OutBytes.Value += int64(v)
	case "is_available":
		v, ok := val.(bool)
		if !ok {
			return fmt.Errorf("stat '%s' value expected bool actual '%v' type %T", name, val, val)
		}
		if v {
			stat.IsAvailable.Value = true
		}
	case "in_bytes":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.InBytes.Value += v
	case "tps_2xx":
		v, ok := val.(int64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps2xx.Value += float64(v)
	case "tps_3xx":
		v, ok := val.(int64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps3xx.Value += float64(v)
	case "tps_4xx":
		v, ok := val.(int64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps4xx.Value += float64(v)
	case "tps_5xx":
		v, ok := val.(int64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps5xx.Value += float64(v)
	case "error_string":
		v, ok := val.(string)
		if !ok {
			return fmt.Errorf("stat '%s' value expected string actual '%v' type %T", name, val, val)
		}
		stat.ErrorString.Value += v + ", "
	case "tps_total":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.TpsTotal.Value += v
	case "status_unknown":
		return dsdata.ErrNotProcessedStat
	default:
		return fmt.Errorf("unknown stat '%s'", name)
	}
	return nil
}
