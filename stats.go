package grove

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type statHandler struct {
	interfaceName string
	stats         Stats
}

// NewStatHandler returns an HTTP handler
func NewStatHandler(interfaceName string, remapRules []string) (http.Handler, Stats) {
	stats := NewStats(remapRules)
	return statHandler{interfaceName: interfaceName, stats: stats}, stats
}

func NewStatHandlerFunc(interfaceName string, remapRules []string) (http.HandlerFunc, Stats) {
	handler, rules := NewStatHandler(interfaceName, remapRules)
	f := func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
	return f, rules
}

type StatsSystem interface {
	AddConfigReloadRequests()
	SetLastReloadRequest(time.Time)
	AddConfigReload()
	SetLastReload(time.Time)
	SetAstatsLoad(time.Time)

	ConfigReloadRequests() uint64
	LastReloadRequest() time.Time
	ConfigReloads() uint64
	LastReload() time.Time
	AstatsLoad() time.Time
}

type Stats interface {
	System() StatsSystem
	Remap() StatsRemaps
}

func NewStats(remapRules []string) Stats {
	return &stats{system: NewStatsSystem(), remap: NewStatsRemaps(remapRules)}
}

type stats struct {
	system StatsSystem
	remap  StatsRemaps
}

func (s *stats) System() StatsSystem {
	return StatsSystem(s.system)
}

func (s *stats) Remap() StatsRemaps { return s.remap }

type StatsRemaps interface {
	Stats(remapRule string) (StatsRemap, bool)
	Rules() []string
}

type StatsRemap interface {
	InBytes() uint64
	AddInBytes(uint64)
	OutBytes() uint64
	AddOutBytes(uint64)
	Status2xx() uint64
	AddStatus2xx(uint64)
	Status3xx() uint64
	AddStatus3xx(uint64)
	Status4xx() uint64
	AddStatus4xx(uint64)
	Status5xx() uint64
	AddStatus5xx(uint64)
}

func NewStatsRemaps(remapRules []string) StatsRemaps {
	m := make(map[string]StatsRemap, len(remapRules))
	for _, rule := range remapRules {
		m[rule] = NewStatsRemap() // must pre-allocate, for threadsafety, so users are never changing the map itself, only the value pointed to.
	}
	return statsRemaps(m)
}

type statsRemaps map[string]StatsRemap

func (s statsRemaps) Stats(rule string) (StatsRemap, bool) {
	r, ok := s[rule]
	return r, ok
}

func (s statsRemaps) Rules() []string {
	rules := make([]string, len(s))
	for rule, _ := range s {
		rules = append(rules, rule)
	}
	return rules
}

func NewStatsRemap() StatsRemap {
	return &statsRemap{}
}

type statsRemap struct {
	inBytes   uint64
	outBytes  uint64
	status2xx uint64
	status3xx uint64
	status4xx uint64
	status5xx uint64
}

func (r *statsRemap) InBytes() uint64       { return atomic.LoadUint64(&r.inBytes) }
func (r *statsRemap) AddInBytes(v uint64)   { atomic.AddUint64(&r.inBytes, v) }
func (r *statsRemap) OutBytes() uint64      { return atomic.LoadUint64(&r.outBytes) }
func (r *statsRemap) AddOutBytes(v uint64)  { atomic.AddUint64(&r.outBytes, v) }
func (r *statsRemap) Status2xx() uint64     { return atomic.LoadUint64(&r.status2xx) }
func (r *statsRemap) AddStatus2xx(v uint64) { atomic.AddUint64(&r.status2xx, v) }
func (r *statsRemap) Status3xx() uint64     { return atomic.LoadUint64(&r.status3xx) }
func (r *statsRemap) AddStatus3xx(v uint64) { atomic.AddUint64(&r.status3xx, v) }
func (r *statsRemap) Status4xx() uint64     { return atomic.LoadUint64(&r.status4xx) }
func (r *statsRemap) AddStatus4xx(v uint64) { atomic.AddUint64(&r.status4xx, v) }
func (r *statsRemap) Status5xx() uint64     { return atomic.LoadUint64(&r.status5xx) }
func (r *statsRemap) AddStatus5xx(v uint64) { atomic.AddUint64(&r.status5xx, v) }

func NewStatsSystem() StatsSystem {
	return &statsSystem{}
}

type statsSystem struct {
	configReloadRequests      uint64
	lastReloadRequestUnixNano int64
	configReloads             uint64
	lastReloadUnixNano        int64
	astatsLoadUnixNano        int64
}

func (s *statsSystem) ConfigReloadRequests() uint64 {
	return atomic.LoadUint64(&s.configReloadRequests)
}
func (s *statsSystem) AddConfigReloadRequests() {
	atomic.AddUint64(&s.configReloadRequests, 1)
}
func (s *statsSystem) LastReloadRequest() time.Time {
	return time.Unix(0, atomic.LoadInt64(&s.lastReloadRequestUnixNano))
}
func (s *statsSystem) SetLastReloadRequest(t time.Time) {
	atomic.StoreInt64(&s.lastReloadRequestUnixNano, t.UnixNano())
}
func (s *statsSystem) ConfigReloads() uint64 {
	return atomic.LoadUint64(&s.configReloads)
}
func (s *statsSystem) AddConfigReload() {
	atomic.AddUint64(&s.configReloads, 1)
}
func (s *statsSystem) LastReload() time.Time {
	return time.Unix(0, atomic.LoadInt64(&s.lastReloadUnixNano))
}
func (s *statsSystem) SetLastReload(t time.Time) {
	atomic.StoreInt64(&s.lastReloadUnixNano, t.UnixNano())
}
func (s *statsSystem) AstatsLoad() time.Time {
	return time.Unix(0, atomic.LoadInt64(&s.astatsLoadUnixNano))
}
func (s *statsSystem) SetAstatsLoad(t time.Time) {
	atomic.StoreInt64(&s.astatsLoadUnixNano, t.UnixNano())
}

const ATSVersion = "5.3.2" // of course, we're not really ATS. We're terrible liars.

type StatsATSJSON struct {
	Server string            `json:"server"`
	Remap  map[string]uint64 `json:"remap"`
}

type StatsSystemJSON struct {
	InterfaceName        string `json:"inf.name"`
	InterfaceSpeed       string `json:"inf.speed"`
	ProcNetDev           string `json:"proc.net.dev"`
	ProcLoadAvg          string `json:"proc.loadavg"`
	ConfigReloadRequests uint64 `json:"configReloadRequests"`
	LastReloadRequest    int64  `json:"lastReloadRequest"`
	ConfigReloads        uint64 `json:"configReloads"`
	LastReload           int64  `json:"lastReload"`
	AstatsLoad           int64  `json:"astatsLoad"`
	Something            string `json:"something"`
}

type StatsJSON struct {
	ATS    StatsATSJSON    `json:"ats"`
	System StatsSystemJSON `json:"system"`
}

func loadFileAndLogGrep(filename string, grepStr string) string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error reading system stat file %v: %v\n", filename, err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		if i := strings.Index(l, grepStr); i == -1 {
			return l
		}
	}
	fmt.Printf("Error reading system stat file %v looking for %v: not found\n", filename, grepStr)
	return ""
}

func loadFileAndLog(filename string) string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading system stat file %v: %v\n", filename, err)
		return ""
	}
	return string(f)
}

func (h statHandler) LoadSystemStats() StatsSystemJSON {
	s := StatsSystemJSON{}
	s.InterfaceName = h.interfaceName
	s.InterfaceSpeed = loadFileAndLog(fmt.Sprintf("/sys/class/net/%v/speed", h.interfaceName))
	s.ProcNetDev = loadFileAndLogGrep("/proc/net/dev", h.interfaceName)
	s.ProcLoadAvg = loadFileAndLog("/proc/loadavg")
	s.ConfigReloadRequests = h.stats.System().ConfigReloadRequests()
	s.LastReloadRequest = h.stats.System().LastReloadRequest().Unix()
	s.ConfigReloads = h.stats.System().ConfigReloads()
	s.LastReload = h.stats.System().LastReload().Unix()
	s.AstatsLoad = h.stats.System().AstatsLoad().Unix()
	s.Something = "here" // emulate existing ATS Astats behavior
	return s
}

func (h statHandler) LoadRemapStats() map[string]uint64 {
	statsRemaps := h.stats.Remap()
	rules := statsRemaps.Rules()
	jsonStats := make(map[string]uint64, len(rules)*6) // remap has 6 members: in, out, 2xx, 3xx, 4xx, 5xx
	for _, rule := range rules {
		statsRemap, ok := statsRemaps.Stats(rule)
		if !ok {
			continue // TODO warn?
		}
		jsonStats[fmt.Sprintf("plugin.remap_stats.%s.in_bytes", rule)] = statsRemap.InBytes()
		jsonStats[fmt.Sprintf("plugin.remap_stats.%s.out_bytes", rule)] = statsRemap.OutBytes()
		jsonStats[fmt.Sprintf("plugin.remap_stats.%s.status_2xx", rule)] = statsRemap.Status2xx()
		jsonStats[fmt.Sprintf("plugin.remap_stats.%s.status_3xx", rule)] = statsRemap.Status3xx()
		jsonStats[fmt.Sprintf("plugin.remap_stats.%s.status_4xx", rule)] = statsRemap.Status4xx()
		jsonStats[fmt.Sprintf("plugin.remap_stats.%s.status_5xx", rule)] = statsRemap.Status5xx()
	}
	return jsonStats
}

func (h statHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	system := h.LoadSystemStats() // TODO goroutine on a timer?
	remap := h.LoadRemapStats()
	stats := StatsJSON{System: system, ATS: StatsATSJSON{Server: ATSVersion, Remap: remap}}
	bytes, err := json.Marshal(stats)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
