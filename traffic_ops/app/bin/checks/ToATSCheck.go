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

/* ToATSCheck.go
   This app collects cache disk usage (CDU) and cache hit ratio (CHR) metrics
   via astats from each cache node and submits the results to TOAPI.
*/

package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
	"github.com/romana/rlog"
)

// Traffic Ops connection params
const AllowInsecureConnections = false
const UserAgent = "go/tc-ats-monitor"
const UseClientCache = false
const TrafficOpsRequestTimeout = time.Second * time.Duration(10)

var confForce *bool
var httpClient = &http.Client{Timeout: 5 * time.Second}

type Config struct {
	URL    string `json:"to_url"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

type Server struct {
	id        int
	cdn       string
	name      string
	fqdn      string
	iface     string
	ip4       string
	ip6       string
	httpPort  string
	httpsPort string
	status    string
	file      string
}

type HttpStats struct {
	ATS struct {
		HitFresh            int64 `json:"proxy.process.http.transaction_counts.hit_fresh"`
		HitFreshProcess     int64 `json:"proxy.process.http.transaction_counts.hit_fresh.process"`
		HitRevalidated      int64 `json:"proxy.process.http.transaction_counts.hit_revalidated"`
		MissCold            int64 `json:"proxy.process.http.transaction_counts.miss_cold"`
		MissNotCacheable    int64 `json:"proxy.process.http.transaction_counts.miss_not_cacheable"`
		MissChanged         int64 `json:"proxy.process.http.transaction_counts.miss_changed"`
		MissClientNoCache   int64 `json:"proxy.process.http.transaction_counts.miss_client_no_cache"`
		ErrAborts           int64 `json:"proxy.process.http.transaction_counts.errors.aborts"`
		ErrPossibleAborts   int64 `json:"proxy.process.http.transaction_counts.errors.possible_aborts"`
		ErrConnectFailed    int64 `json:"proxy.process.http.transaction_counts.errors.connect_failed"`
		ErrPreAcceptHangups int64 `json:"proxy.process.http.transaction_counts.errors.pre_accept_hangups"`
		ErrOther            int64 `json:"proxy.process.http.transaction_counts.errors.other"`
		ErrUnclassified     int64 `json:"proxy.process.http.transaction_counts.other.unclassified"`
	} `json:"ats"`
	Time int64
}

type DiskStats struct {
	ATS struct {
		BytesUsed    int64 `json:"proxy.process.cache.bytes_used"`
		BytesTotal   int64 `json:"proxy.process.cache.bytes_total"`
		RCBytesUsed  int64 `json:"proxy.process.cache.ram_cache.bytes_used"`
		BytesUsed1   int64 `json:"proxy.process.cache.volume_1.bytes_used"`
		BytesTotal1  int64 `json:"proxy.process.cache.volume_1.bytes_total"`
		RCBytesUsed1 int64 `json:"proxy.process.cache.volume_1.ram_cache.bytes_used"`
		BytesUsed2   int64 `json:"proxy.process.cache.volume_2.bytes_used"`
		BytesTotal2  int64 `json:"proxy.process.cache.volume_2.bytes_total"`
		RCBytesUsed2 int64 `json:"proxy.process.cache.volume_2.ram_cache.bytes_used"`
	} `json:"ats"`
	Time int64
}

func NewServer(id int, name string, status string) Server {
	server := Server{}
	server.id = id
	server.name = name
	server.status = status
	return server
}

func LoadConfig(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, err
}

func readFile(f string) (HttpStats, error) {
	var stats HttpStats
	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		return stats, err
	}

	stats = HttpStats{}
	err = binary.Read(file, binary.BigEndian, &stats)
	if err != nil {
		return stats, err
	}
	return stats, nil
}

func writeFile(f string, stats HttpStats) error {
	t := time.Now()
	stats.Time = t.Unix()
	file, err := os.Create(f)
	defer file.Close()
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.BigEndian, stats)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) getDiskStats() (DiskStats, error) {
	rlog.Debug("starting getHttpStats()")
	var stats DiskStats
	stats = DiskStats{}
	var ipaddr string
	if s.ip4 != "" {
		ipaddr = s.ip4
	} else if s.ip6 != "" {
		ipaddr = s.ip6
	} else {
		err := errors.New("Server has neither IPv4 nor IPv6 address")
		if err != nil {
			return stats, err
		}
	}
	url := "http://" + ipaddr + ":" + s.httpPort + "/_astats?application=bytes_used;bytes_total&inf.name=" + s.iface
	rlog.Debugf("fetching: %s", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		return stats, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return stats, err
	}
	err = json.Unmarshal(body, &stats)
	stats.Time = time.Now().Unix()
	if err != nil {
		return stats, err
	}
	return stats, err
}

func (s *Server) getHttpStats() (HttpStats, error) {
	rlog.Debug("starting getHttpStats()")
	var stats HttpStats
	stats = HttpStats{}
	var ipaddr string
	if s.ip4 != "" {
		ipaddr = s.ip4
	} else if s.ip6 != "" {
		ipaddr = s.ip6
	} else {
		err := errors.New("Server has neither IPv4 nor IPv6 address")
		if err != nil {
			return stats, err
		}
	}
	url := "http://" + ipaddr + ":" + s.httpPort + "/_astats?application=proxy.process.http.transaction_counts&inf.name=" + s.iface
	rlog.Debugf("fetching: %s", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		return stats, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return stats, err
	}
	err = json.Unmarshal(body, &stats)
	stats.Time = time.Now().Unix()
	if err != nil {
		return stats, err
	}
	return stats, err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getCDU(d DiskStats) (usage int) {
	// TO-DO: get per-volume stats? For now, total usage
	rlog.Infof("disk_total=%d disk_used=%d", d.ATS.BytesUsed, d.ATS.BytesTotal)
	b := float64(d.ATS.BytesUsed)
	t := float64(d.ATS.BytesTotal)
	u := b / t * 100
	rlog.Infof("usage_perc=%f", u)
	usage = int(u)
	return
}

func getRatio(c HttpStats, p HttpStats) (ratio int, seconds int) {
	ratio = -1
	// HitFreshProcess seems to be a redundant metric, so not using it
	// for calculations
	hits := c.ATS.HitFresh + c.ATS.HitRevalidated
	miss := c.ATS.MissCold + c.ATS.MissNotCacheable + c.ATS.MissChanged + c.ATS.MissClientNoCache
	errs := c.ATS.ErrAborts + c.ATS.ErrPossibleAborts + c.ATS.ErrConnectFailed + c.ATS.ErrPreAcceptHangups + c.ATS.ErrOther + c.ATS.ErrUnclassified
	prevHits := p.ATS.HitFresh + p.ATS.HitRevalidated
	prevMiss := p.ATS.MissCold + p.ATS.MissNotCacheable + p.ATS.MissChanged + p.ATS.MissClientNoCache
	prevErrs := p.ATS.ErrAborts + p.ATS.ErrPossibleAborts + p.ATS.ErrConnectFailed + p.ATS.ErrPreAcceptHangups + p.ATS.ErrOther + p.ATS.ErrUnclassified
	h := float64(hits - prevHits)
	m := float64(miss - prevMiss)
	e := float64(errs - prevErrs)
	t := h + m + e
	seconds = int(c.Time - p.Time)
	rlog.Infof("hits=%f miss=%f err=%f total=%f seconds=%d", h, m, e, t, seconds)
	if t != 0 {
		r := (h / t) * 100
		rlog.Infof("hit_ratio=%f", r)
		ratio = int(r)
		return
	}
	return
}

func main() {
	var cpath_new string
	var seconds int

	jobStart := time.Now()

	// define default config file path
	cpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		rlog.Error("Config error:", err)
		os.Exit(1)
	}
	cpath_new = strings.Replace(cpath, "/bin/checks", "/conf/check-config.json", 1)

	// command-line flags
	confPtr := flag.String("conf", cpath_new, "Config file path")
	confName := flag.String("name", "undef", "Check name 'CDU' (cache disk usage) or 'CHR' (cache hit ratio)")
	confInclude := flag.String("host", "undef", "Specific host or regex to include (optional)")
	confSyslog := flag.Bool("syslog", false, "Log check results to syslog")
	confCdn := flag.String("cdn", "all", "Check specific CDN by name")
	confExclude := flag.String("exclude", "undef", "Hostname regex to exclude")
	confReset := flag.Bool("reset", false, "Reset check values in TO to 'blank' state")
	confQuiet := flag.Bool("q", false, "Do not send updates to TO")
	confForce = flag.Bool("f", false, "Force a failure result")
	confPath := flag.String("path", "/var/tmp/tc-checks", "Path to store persistent data files")
	flag.Parse()

	// configure syslog logger
	if *confSyslog == true {
		logwriter, err := syslog.New(syslog.LOG_INFO, os.Args[0])
		if err == nil {
			log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
			log.SetOutput(logwriter)
		}
	}

	if *confName != "CHR" && *confName != "CDU" {
		rlog.Error("Check name should be 'CHR' or 'CDU'")
		os.Exit(1)
	}

	// create persistant data file path, if not present
	os.MkdirAll(*confPath, os.ModePerm)

	// load config json
	config, err := LoadConfig(*confPtr)
	if err != nil {
		rlog.Error("Error loading config:", err)
		os.Exit(1)
	}

	// connect to TO
	session, _, err := toclient.LoginWithAgent(
		config.URL,
		config.User,
		config.Passwd,
		AllowInsecureConnections,
		UserAgent,
		UseClientCache,
		TrafficOpsRequestTimeout)
	if err != nil {
		rlog.Criticalf("An error occurred while logging in: %v\n", err)
		os.Exit(1)
	}

	// Make TO API call for server details
	var servers []tc.Server
	servers, _, err = session.GetServers(nil)
	if err != nil {
		rlog.Criticalf("An error occurred while getting servers: %v\n", err)
		os.Exit(1)
	}

	for _, server := range servers {
		re, err := regexp.Compile("^(MID|EDGE).*")
		if err != nil {
			rlog.Error("supplied exclusion regex does not compile:", err)
			os.Exit(1)
		}
		if re.Match([]byte(server.Type)) {
			serverStart := time.Now()
			if *confInclude != "undef" {
				re_inc, err := regexp.Compile(*confInclude)
				if err != nil {
					rlog.Error("supplied exclusion regex does not compile:", err)
					os.Exit(1)
				}
				if !re_inc.MatchString(server.HostName) {
					rlog.Debugf("%s does not match the provided include regex, skipping", server.HostName)
					continue
				}
			}
			if *confCdn != "all" && *confCdn != server.CDNName {
				rlog.Debugf("%s is not assinged to the specified CDN '%s', skipping", server.HostName, *confCdn)
				continue
			}
			if *confExclude != "undef" {
				re, err := regexp.Compile(*confExclude)
				if err != nil {
					rlog.Error("supplied exclusion regex does not compile:", err)
					os.Exit(1)
				}
				if re.MatchString(server.HostName) {
					rlog.Debugf("%s matches the provided exclude regex, skipping", server.HostName)
					continue
				}
			}
			s := NewServer(server.ID, server.HostName, server.Status)
			defaulStatusValue := -1
			var statusData tc.ServercheckRequestNullable
			statusData.ID = &s.id
			statusData.Name = confName
			statusData.HostName = &s.name
			statusData.Value = &defaulStatusValue
			s.fqdn = s.name + "." + server.DomainName
			s.iface = server.InterfaceName
			s.ip4 = strings.Split(server.IPAddress, "/")[0]
			s.ip6 = strings.Split(server.IP6Address, "/")[0]
			s.file = *confPath + "/" + s.fqdn + "_chr.dat"
			s.cdn = server.CDNName
			rlog.Infof("Next server=%s status=%s", s.fqdn, s.status)

			if *confSyslog {
				log.Printf("Next server=%s status=%s", s.fqdn, s.status)
			}
			if s.status == "REPORTED" && *confReset != true {
				if *confName == "CHR" {
					var check_against_prev bool
					var preStats HttpStats
					var curStats HttpStats
					if fileExists(s.file) {
						rlog.Debugf("Data file exists, using for compare: %s", s.file)
						check_against_prev = true
					} else {
						rlog.Infof("Data file doesn't exist; initializing: %s", s.file)
						check_against_prev = false
					}
					curStats, err = s.getHttpStats()
					if err != nil {
						rlog.Errorf("Error: %s", err)
						continue
					}
					if check_against_prev == true {
						preStats, err = readFile(s.file)
						rlog.Debugf("previous: %v", preStats)
						rlog.Debugf("current: %v", curStats)
						*statusData.Value, seconds = getRatio(curStats, preStats)
						err = writeFile(s.file, curStats)
					} else {
						err = writeFile(s.file, curStats)
						rlog.Info("persistent data initialized; going to next server")
						continue
					}

					if err != nil {
						rlog.Errorf("Error: %s", err)
					}
				} else if *confName == "CDU" {
					var curStats DiskStats
					curStats, err = s.getDiskStats()
					if err != nil {
						rlog.Errorf("Error: %s", err)
						continue
					}
					*statusData.Value = getCDU(curStats)
				}
			}

			serverElapsed := time.Since(serverStart)
			if *confName == "CHR" {
				rlog.Infof("Finished checking server=%s check=%s result=%d cdn=%s seconds=%d elapsed=%s", s.fqdn, *confName, *statusData.Value, s.cdn, seconds, serverElapsed)
				if *confSyslog {
					log.Printf("Finished checking server=%s check=%s result=%d cdn=%s seconds=%d elapsed=%s", s.fqdn, *confName, *statusData.Value, s.cdn, seconds, serverElapsed)
				}
			} else if *confName == "CDU" {
				rlog.Infof("Finished checking server=%s check=%s result=%d cdn=%s elapsed=%s", s.fqdn, *confName, *statusData.Value, s.cdn, serverElapsed)
				if *confSyslog {
					log.Printf("Finished checking check=%s server=%s result=%d cdn=%s elapsed=%s", s.fqdn, *confName, *statusData.Value, s.cdn, serverElapsed)
				}
			}

			if *confQuiet == false {
				rlog.Debug("Sending update to TO")
				_, _, err := session.InsertServerCheckStatus(statusData)
				if err != nil {
					rlog.Error("Error updating server check status with TO:", err)
				}
				// we seem to be reusing TO cons, and updating super-fast, so...
				// lets slow things down a bit
				time.Sleep(100 * time.Millisecond)
			} else {
				rlog.Debug("Skipping update to TO")
			}
		}
	}
	jobElapsed := time.Since(jobStart)
	rlog.Info("Job complete", jobElapsed)
	if *confSyslog {
		log.Print("Job complete totaltime=", jobElapsed)
	}
	os.Exit(0)
}
