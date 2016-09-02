package cache

import (
	"encoding/json"
	"fmt"
	dsdata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservicedata"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/log"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	"io"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	ResultChannel chan Result
	Notify        int
	ToData        *todata.TODataThreadsafe
	PeerStates    *peer.CRStatesPeersThreadsafe
}

// NewHandler does NOT precomputes stat data before calling ResultChannel, and Result.Precomputed will be nil
func NewHandler() Handler {
	return Handler{ResultChannel: make(chan Result)}
}

// NewPrecomputeHandler precomputes stat data and populates result.Precomputed before passing to ResultChannel.
func NewPrecomputeHandler(toData todata.TODataThreadsafe, peerStates peer.CRStatesPeersThreadsafe) Handler {
	return Handler{ResultChannel: make(chan Result), ToData: &toData, PeerStates: &peerStates}
}

func (h Handler) Precompute() bool {
	return h.ToData != nil && h.PeerStates != nil
}

type PrecomputedData struct {
	DeliveryServiceStats map[enum.DeliveryServiceName]dsdata.Stat
	OutBytes             int64
	Errors               []error
	Reporting            bool
}

type Result struct {
	Id        string
	Available bool
	Errors    []error
	Astats    Astats
	Time      time.Time
	Vitals    Vitals
	PrecomputedData
	PollID       uint64
	PollFinished chan<- uint64
}

type Vitals struct {
	LoadAvg    float64
	BytesOut   int64
	BytesIn    int64
	KbpsOut    int64
	MaxKbpsOut int64
}

type Stat struct {
	Time  int64       `json:"time"`
	Value interface{} `json:"value"`
}

type Stats struct {
	Caches map[string]map[string][]Stat `json:"caches"`
}

const (
	NOTIFY_NEVER = iota
	NOTIFY_CHANGE
	NOTIFY_ALWAYS
)

func StatsMarshall(statHistory map[enum.CacheName][]Result, historyCount int) ([]byte, error) {
	var stats Stats

	stats.Caches = map[string]map[string][]Stat{}

	count := 1

	for id, history := range statHistory {
		for _, result := range history {
			for stat, value := range result.Astats.Ats {
				s := Stat{
					Time:  result.Time.UnixNano() / 1000000,
					Value: value,
				}

				_, exists := stats.Caches[string(id)]

				if !exists {
					stats.Caches[string(id)] = map[string][]Stat{}
				}

				stats.Caches[string(id)][stat] = append(stats.Caches[string(id)][stat], s)
			}

			if historyCount > 0 && count == historyCount {
				break
			}

			count++
		}
	}

	return json.Marshal(stats)
}

func (handler Handler) Handle(id string, r io.Reader, err error, pollId uint64, pollFinished chan<- uint64) {
	log.Debugf("poll %v %v handle start\n", pollId, time.Now())
	result := Result{
		Id:           id,
		Available:    false,
		Errors:       []error{},
		Time:         time.Now(), // TODO change this to be computed the instant we get the result back, to minimise inaccuracy
		PollID:       pollId,
		PollFinished: pollFinished,
	}

	if err != nil {
		log.Errorf("%v handler given error '%v'\n", id, err) // error here, in case the thing that called Handle didn't error
		result.Errors = append(result.Errors, err)
		handler.ResultChannel <- result
		return
	}

	if r == nil {
		log.Errorf("%v handle reader nil\n", id)
		result.Errors = append(result.Errors, fmt.Errorf("handler got nil reader"))
		handler.ResultChannel <- result
		return
	}

	result.PrecomputedData.Reporting = true

	if err := json.NewDecoder(r).Decode(&result.Astats); err != nil {
		log.Errorf("%s procnetdev decode error '%v'\n", id, err)
		result.Errors = append(result.Errors, err)
		handler.ResultChannel <- result
		return
	}

	if result.Astats.System.ProcNetDev == "" {
		log.Warnf("addkbps %s procnetdev empty\n", id)
	}

	log.Debugf("poll %v %v handle decode end\n", pollId, time.Now())

	if err != nil {
		result.Errors = append(result.Errors, err)
		log.Errorf("addkbps handle %s error '%v'\n", id, err)
	} else {
		result.Available = true
	}

	if handler.Precompute() {
		log.Debugf("poll %v %v handle precompute start\n", pollId, time.Now())
		result = handler.precompute(result)
		log.Debugf("poll %v %v handle precompute end\n", pollId, time.Now())
	}
	log.Debugf("poll %v %v handle write start\n", pollId, time.Now())
	handler.ResultChannel <- result
	log.Debugf("poll %v %v handle end\n", pollId, time.Now())
}

// outBytes takes the proc.net.dev string, and the interface name, and returns the bytes field
func outBytes(procNetDev, iface string) (int64, error) {
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
	procNetDevIfaceBytes = strings.TrimLeft(procNetDevIfaceBytes, " ")
	spacePos := strings.Index(procNetDevIfaceBytes, " ")
	if spacePos != -1 {
		procNetDevIfaceBytes = procNetDevIfaceBytes[:spacePos]
	}
	return strconv.ParseInt(procNetDevIfaceBytes, 10, 64)
}

// precompute does the calculations which are possible with only this one cache result.
func (handler Handler) precompute(result Result) Result {
	todata := handler.ToData.Get()
	stats := map[enum.DeliveryServiceName]dsdata.Stat{}

	var err error
	if result.PrecomputedData.OutBytes, err = outBytes(result.Astats.System.ProcNetDev, result.Astats.System.InfName); err != nil {
		result.PrecomputedData.OutBytes = 0
		log.Errorf("addkbps %s handle precomputing outbytes '%v'\n", result.Id, err)
	}

	for stat, value := range result.Astats.Ats {
		var err error
		stats, err = processStat(result.Id, stats, todata, stat, value)
		if err != nil && err != dsdata.ErrNotProcessedStat {
			log.Errorf("precomputing cache %v stat %v value %v error %v", result.Id, stat, value, err)
			result.PrecomputedData.Errors = append(result.PrecomputedData.Errors, err)
		}
	}
	result.PrecomputedData.DeliveryServiceStats = stats
	return result
}

// processStat and its subsidiary functions act as a State Machine, flowing the stat thru states for each "." component of the stat name
// TODO fix this being crazy slow. THIS IS THE BOTTLENECK
func processStat(server string, stats map[enum.DeliveryServiceName]dsdata.Stat, toData todata.TOData, stat string, value interface{}) (map[enum.DeliveryServiceName]dsdata.Stat, error) {
	parts := strings.Split(stat, ".")
	if len(parts) < 1 {
		return stats, fmt.Errorf("stat has no initial part")
	}

	switch parts[0] {
	case "plugin":
		return processStatPlugin(server, stats, toData, stat, parts[1:], value)
	case "proxy":
		return stats, dsdata.ErrNotProcessedStat
	case "server":
		return stats, dsdata.ErrNotProcessedStat
	default:
		return stats, fmt.Errorf("stat '%s' has unknown initial part '%s'", stat, parts[0])
	}
}

func processStatPlugin(server string, stats map[enum.DeliveryServiceName]dsdata.Stat, toData todata.TOData, stat string, statParts []string, value interface{}) (map[enum.DeliveryServiceName]dsdata.Stat, error) {
	if len(statParts) < 1 {
		return stats, fmt.Errorf("stat has no plugin part")
	}
	switch statParts[0] {
	case "remap_stats":
		return processStatPluginRemapStats(server, stats, toData, stat, statParts[1:], value)
	default:
		return stats, fmt.Errorf("stat has unknown plugin part '%s'", statParts[0])
	}
}

func processStatPluginRemapStats(server string, stats map[enum.DeliveryServiceName]dsdata.Stat, toData todata.TOData, stat string, statParts []string, value interface{}) (map[enum.DeliveryServiceName]dsdata.Stat, error) {
	if len(statParts) < 2 {
		return stats, fmt.Errorf("stat has no remap_stats deliveryservice and name parts")
	}

	fqdn := strings.Join(statParts[:len(statParts)-1], ".")

	ds, ok := toData.DeliveryServiceRegexes.DeliveryService(fqdn)
	if !ok {
		return stats, fmt.Errorf("ERROR no delivery service match for fqdn '%v' stat '%v'\n", fqdn, strings.Join(statParts, "."))
	}
	if ds == "" {
		return stats, fmt.Errorf("ERROR EMPTY delivery service fqdn %v stat %v\n", fqdn, strings.Join(statParts, "."))
	}

	statName := statParts[len(statParts)-1]

	dsStat, ok := stats[ds]
	if !ok {
		newStat := dsdata.NewStat()
		dsStat = *newStat
	}

	if err := addCacheStat(&dsStat.Total, statName, value); err != nil {
		return stats, err
	}

	cachegroup, ok := toData.ServerCachegroups[enum.CacheName(server)]
	if !ok {
		return stats, fmt.Errorf("server missing from TOData.ServerCachegroups") // TODO check logs, make sure this isn't normal
	}
	dsStat.CacheGroups[cachegroup] = dsStat.Total

	cacheType, ok := toData.ServerTypes[enum.CacheName(server)]
	if !ok {
		return stats, fmt.Errorf("server missing from TOData.ServerTypes")
	}
	dsStat.Type[cacheType] = dsStat.Total
	stats[ds] = dsStat
	return stats, nil
}

// addCacheStat adds the given stat to the existing stat. Note this adds, it doesn't overwrite. Numbers are summed, strings are concatenated.
// TODO make this less duplicate code somehow.
func addCacheStat(stat *dsdata.StatCacheStats, name string, val interface{}) error {
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
		log.Debugln("got is_available")
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
		stat.Tps2xx.Value += v
	case "tps_3xx":
		v, ok := val.(int64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps3xx.Value += v
	case "tps_4xx":
		v, ok := val.(int64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps4xx.Value += v
	case "tps_5xx":
		v, ok := val.(int64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Tps5xx.Value += v
	case "error_string":
		v, ok := val.(string)
		if !ok {
			return fmt.Errorf("stat '%s' value expected string actual '%v' type %T", name, val, val)
		}
		stat.ErrorString.Value += v + ", "
	case "tps_total":
		v, ok := val.(int64)
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
