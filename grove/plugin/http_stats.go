package plugin

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
	"unicode"

	"github.com/apache/incubator-trafficcontrol/grove/remapdata"
	"github.com/apache/incubator-trafficcontrol/grove/stat"
	"github.com/apache/incubator-trafficcontrol/grove/web"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
)

func init() {
	AddPlugin(10000, Funcs{onRequest: stats})
}

const StatsEndpoint = "/_astats"

func stats(icfg interface{}, d OnRequestData) bool {
	if !strings.HasPrefix(d.R.URL.Path, StatsEndpoint) {
		log.Debugf("plugin onrequest http_stats returning, not in path '" + d.R.URL.Path + "'\n")
		return false
	}

	log.Debugf("plugin onrequest http_stats calling\n")

	w := d.W
	req := d.R

	// TODO access log? Stats byte count?
	ip, err := web.GetIP(req)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		log.Errorln("statHandler ServeHTTP failed to get IP: " + ip.String())
		return true
	}
	if !Allowed(d.StatRules, ip) {
		code := http.StatusForbidden
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
		log.Debugln("statHandler.ServeHTTP IP " + ip.String() + " FORBIDDEN") // TODO event?
		return true
	}

	// TODO gzip
	system := LoadSystemStats(d.Stats, d.InterfaceName) // TODO goroutine on a timer?
	ats := map[string]interface{}{"server": "6.2.1"}
	if req.URL.Query().Get("application") != "system" {
		ats = LoadRemapStats(d.Stats, d.HTTPConns, d.HTTPSConns)
	}
	stats := stat.StatsJSON{System: system, ATS: ats}

	bytes, err := json.Marshal(stats)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte(http.StatusText(code)))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	return true
}

func LoadSystemStats(stats stat.Stats, interfaceName string) stat.StatsSystemJSON {
	s := stat.StatsSystemJSON{}
	s.InterfaceName = interfaceName
	s.InterfaceSpeed = loadFileAndLogInt(fmt.Sprintf("/sys/class/net/%v/speed", interfaceName))
	s.ProcNetDev = loadFileAndLogGrep("/proc/net/dev", interfaceName)
	s.ProcLoadAvg = loadFileAndLog("/proc/loadavg")
	s.ConfigReloadRequests = stats.System().ConfigReloadRequests()
	s.LastReloadRequest = stats.System().LastReloadRequest().Unix()
	s.ConfigReloads = stats.System().ConfigReloads()
	s.LastReload = stats.System().LastReload().Unix()
	s.AstatsLoad = stats.System().AstatsLoad().Unix()
	s.Something = "here" // emulate existing ATS Astats behavior
	return s
}

func LoadRemapStats(stats stat.Stats, httpConns *web.ConnMap, httpsConns *web.ConnMap) map[string]interface{} {
	statsRemaps := stats.Remap()
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

	jsonStats["proxy.process.http.current_client_connections"] = httpConns.Len() + httpsConns.Len()
	jsonStats["proxy.process.http.cache_hits"] = stats.CacheHits()
	jsonStats["proxy.process.http.cache_misses"] = stats.CacheMisses()
	jsonStats["proxy.process.http.cache_capacity_bytes"] = stats.CacheCapacity()
	jsonStats["proxy.process.http.cache_size_bytes"] = stats.CacheSize()

	return jsonStats
}

func Allowed(statRules remapdata.RemapRulesStats, ip net.IP) bool {
	// TODO remove duplication
	for _, network := range statRules.Deny {
		if network.Contains(ip) {
			log.Debugf("deny contains ip\n")
			return false
		}
	}
	if len(statRules.Allow) == 0 {
		log.Debugf("Allowed len 0\n")
		return true
	}
	for _, network := range statRules.Allow {
		if network.Contains(ip) {
			log.Debugf("allow contains ip\n")
			return true
		}
	}
	return false
}

func loadFileAndLog(filename string) string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("reading system stat file %v: %v\n", filename, err)
		return ""
	}
	return strings.TrimSpace(string(f))
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

func loadFileAndLogInt(filename string) int64 {
	s := loadFileAndLog(filename)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Errorf("parsing system stat file %v: %v\n", filename, err)
	}
	return i
}
