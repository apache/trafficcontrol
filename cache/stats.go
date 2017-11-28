package cache

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

type statHandler struct {
	interfaceName string
	stats         Stats
	statRules     RemapRulesStats
}

// NewStatHandler returns an HTTP handler
func NewStatHandler(interfaceName string, remapRules []RemapRule, stats Stats, statRules RemapRulesStats) http.Handler {
	return statHandler{interfaceName: interfaceName, stats: stats, statRules: statRules}
}

func NewStatHandlerFunc(interfaceName string, remapRules []RemapRule, stats Stats, statRules RemapRulesStats) http.HandlerFunc {
	handler := NewStatHandler(interfaceName, remapRules, stats, statRules)
	f := func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
	return f
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

	Connections() uint64
	IncConnections()
	DecConnections()

	CacheHits() uint64
	AddCacheHit()
	CacheMisses() uint64
	AddCacheMiss()

	CacheSize() uint64
	CacheCapacity() uint64
}

func NewStats(remapRules []RemapRule, cache Cache, cacheCapacityBytes uint64) Stats {
	connections := uint64(0)
	cacheHits := uint64(0)
	cacheMisses := uint64(0)
	return &stats{
		system:             NewStatsSystem(),
		remap:              NewStatsRemaps(remapRules),
		connections:        &connections,
		cacheHits:          &cacheHits,
		cacheMisses:        &cacheMisses,
		cache:              cache,
		cacheCapacityBytes: cacheCapacityBytes,
	}
}

type stats struct {
	system             StatsSystem
	remap              StatsRemaps
	connections        *uint64
	cacheHits          *uint64
	cacheMisses        *uint64
	cache              Cache
	cacheCapacityBytes uint64
}

func (s stats) Connections() uint64   { return atomic.LoadUint64(s.connections) }
func (s stats) IncConnections()       { atomic.AddUint64(s.connections, 1) }
func (s stats) DecConnections()       { atomic.AddUint64(s.connections, ^uint64(0)) }
func (s stats) CacheHits() uint64     { return atomic.LoadUint64(s.cacheHits) }
func (s stats) AddCacheHit()          { atomic.AddUint64(s.cacheHits, 1) }
func (s stats) CacheMisses() uint64   { return atomic.LoadUint64(s.cacheMisses) }
func (s stats) AddCacheMiss()         { atomic.AddUint64(s.cacheMisses, 1) }
func (s *stats) System() StatsSystem  { return StatsSystem(s.system) }
func (s *stats) Remap() StatsRemaps   { return s.remap }
func (s stats) CacheSize() uint64     { return s.cache.Size() }
func (s stats) CacheCapacity() uint64 { return s.cacheCapacityBytes }

type StatsRemaps interface {
	Stats(fqdn string) (StatsRemap, bool)
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

	CacheHits() uint64
	AddCacheHit()
	CacheMisses() uint64
	AddCacheMiss()
}

func getFromFQDN(r RemapRule) string {
	path := r.From
	schemeEnd := `://`
	if i := strings.Index(path, schemeEnd); i != -1 {
		path = path[i+len(schemeEnd):]
	}
	pathStart := `/`
	if i := strings.Index(path, pathStart); i != -1 {
		path = path[:i]
	}
	return path
}

func NewStatsRemaps(remapRules []RemapRule) StatsRemaps {
	m := make(map[string]StatsRemap, len(remapRules))
	for _, rule := range remapRules {
		m[getFromFQDN(rule)] = NewStatsRemap() // must pre-allocate, for threadsafety, so users are never changing the map itself, only the value pointed to.
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
	for rule := range s {
		rules = append(rules, rule)
	}
	return rules
}

func NewStatsRemap() StatsRemap {
	return &statsRemap{}
}

type statsRemap struct {
	inBytes     uint64
	outBytes    uint64
	status2xx   uint64
	status3xx   uint64
	status4xx   uint64
	status5xx   uint64
	cacheHits   uint64
	cacheMisses uint64
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

func (r *statsRemap) CacheHits() uint64 { return atomic.LoadUint64(&r.cacheHits) }
func (r *statsRemap) AddCacheHit()      { atomic.AddUint64(&r.cacheHits, 1) }

func (r *statsRemap) CacheMisses() uint64 { return atomic.LoadUint64(&r.cacheMisses) }
func (r *statsRemap) AddCacheMiss()       { atomic.AddUint64(&r.cacheMisses, 1) }

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

// type StatsATSJSON struct {
// 	Server string            `json:"server"`
// 	Remap  map[string]uint64 `json:"remap"`
// }

type StatsSystemJSON struct {
	InterfaceName        string `json:"inf.name"`
	InterfaceSpeed       int64  `json:"inf.speed"`
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
	ATS    map[string]interface{} `json:"ats"`
	System StatsSystemJSON        `json:"system"`
}

func loadFileAndLogGrep(filename string, grepStr string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Errorf("reading system stat file %v: %v\n", filename, err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		l = strings.TrimLeftFunc(l, unicode.IsSpace)
		if strings.HasPrefix(l, grepStr) {
			return l
		}
	}
	log.Errorf("reading system stat file %v looking for %v: not found\n", filename, grepStr)
	return ""
}

func loadFileAndLog(filename string) string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("reading system stat file %v: %v\n", filename, err)
		return ""
	}
	return strings.TrimSpace(string(f))
}

func loadFileAndLogInt(filename string) int64 {
	s := loadFileAndLog(filename)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Errorf("parsing system stat file %v: %v\n", filename, err)
	}
	return i
}

func (h statHandler) LoadSystemStats() StatsSystemJSON {
	s := StatsSystemJSON{}
	s.InterfaceName = h.interfaceName
	s.InterfaceSpeed = loadFileAndLogInt(fmt.Sprintf("/sys/class/net/%v/speed", h.interfaceName))
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

func (h statHandler) LoadRemapStats() map[string]interface{} {
	statsRemaps := h.stats.Remap()
	rules := statsRemaps.Rules()
	jsonStats := make(map[string]interface{}, len(rules)*8) // remap has 8 members: in, out, 2xx, 3xx, 4xx, 5xx, hits, misses
	jsonStats["server"] = "6.2.1"
	for _, rule := range rules {
		ruleName := rule
		statsRemap, ok := statsRemaps.Stats(ruleName)
		if !ok {
			continue // TODO warn?
		}
		jsonStats["plugin.remap_stats."+ruleName+".in_bytes"] = statsRemap.InBytes()
		jsonStats["plugin.remap_stats."+ruleName+".out_bytes"] = statsRemap.OutBytes()
		jsonStats["plugin.remap_stats."+ruleName+".status_2xx"] = statsRemap.Status2xx()
		jsonStats["plugin.remap_stats."+ruleName+".status_3xx"] = statsRemap.Status3xx()
		jsonStats["plugin.remap_stats."+ruleName+".status_4xx"] = statsRemap.Status4xx()
		jsonStats["plugin.remap_stats."+ruleName+".status_5xx"] = statsRemap.Status5xx()
		jsonStats["plugin.remap_stats."+ruleName+".cache_hits"] = statsRemap.CacheHits()
		jsonStats["plugin.remap_stats."+ruleName+".cache_misses"] = statsRemap.CacheMisses()
	}

	jsonStats["proxy.process.http.current_client_connections"] = h.stats.Connections()
	jsonStats["proxy.process.http.cache_hits"] = h.stats.CacheHits()
	jsonStats["proxy.process.http.cache_misses"] = h.stats.CacheMisses()
	jsonStats["proxy.process.http.cache_capacity_bytes"] = h.stats.CacheCapacity()
	jsonStats["proxy.process.http.cache_size_bytes"] = h.stats.CacheSize()

	return jsonStats
}

func (h statHandler) Allowed(ip net.IP) bool {
	// TODO remove duplication
	for _, network := range h.statRules.Deny {
		if network.Contains(ip) {
			log.Debugf("deny contains ip\n")
			return false
		}
	}
	if len(h.statRules.Allow) == 0 {
		log.Debugf("Allowed len 0\n")
		return true
	}
	for _, network := range h.statRules.Allow {
		if network.Contains(ip) {
			log.Debugf("allow contains ip\n")
			return true
		}
	}
	return false
}

func (h statHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO access log? Stats byte count?
	ip, err := GetIP(r)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		log.Errorln("statHandler ServeHTTP failed to get IP: " + ip.String())
		return
	}
	if !h.Allowed(ip) {
		code := http.StatusForbidden
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		log.Debugln("statHandler.ServeHTTP IP " + ip.String() + " FORBIDDEN") // TODO event?
		return
	}

	// TODO gzip
	system := h.LoadSystemStats() // TODO goroutine on a timer?
	ats := map[string]interface{}{"server": "6.2.1"}
	if r.URL.Query().Get("application") != "system" {
		ats = h.LoadRemapStats()
	}
	stats := StatsJSON{System: system, ATS: ats}

	bytes, err := json.Marshal(stats)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
