package plugin

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/apache/trafficcontrol/v8/grove/stat"
	"github.com/apache/trafficcontrol/v8/grove/web"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
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
	if !d.StatRules.Allowed(ip) {
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
	s.Version = stats.System().Version()
	return s
}

func LoadRemapStats(stats stat.Stats, httpConns *web.ConnMap, httpsConns *web.ConnMap) map[string]interface{} {
	statsRemaps := stats.Remap()
	rules := statsRemaps.Rules()
	jsonStats := make(map[string]interface{}, len(rules)*8) // remap has 8 members: in, out, 2xx, 3xx, 4xx, 5xx, hits, misses
	jsonStats["server"] = "6.2.1"                           // emulate a good ATS version
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
