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

/* ToFQDNCheck.go
   This app verifies that forward DNS (A/AAAA) records match what is
   recorded in TODB for each server. Optionally, it will validate
   reverse (PTR) records as well to ensure that they agree with
   forward DNS.
*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"log/syslog"
	"net"
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
const UserAgent = "go/tc-fqdn-monitor"
const UseClientCache = false
const TrafficOpsRequestTimeout = time.Second * time.Duration(10)

var (
	confForce  *bool
	confPtr    *bool
	confSyslog *bool
)

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
	ip4       string
	ip6       string
	httpPort  string
	httpsPort string
	failcount int
	status    string
}

func NewServer(id int, name string, status string, f int) Server {
	server := Server{}
	server.id = id
	server.name = name
	server.status = status
	server.failcount = f
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

func contains(s []net.IP, i net.IP) bool {
	for _, a := range s {
		if a.String() == i.String() {
			return true
		}
	}
	return false
}

func (s *Server) check_dns(server tc.Server) string {
	var dnsrr []net.IP
	var tcrr []net.IP
	s.failcount = 0
	rlog.Infof("checking A/AAAA records for %v", s.fqdn)
	if *confSyslog {
		log.Printf("checking forward records for %v", s.fqdn)
	}
	s.cdn = server.CDNName
	s.ip4 = strings.Split(server.IPAddress, "/")[0]
	s.ip6 = strings.Split(server.IP6Address, "/")[0]
	if s.ip4 != "" {
		tcrr = append(tcrr, net.ParseIP(s.ip4))
	}
	if s.ip6 != "" {
		tcrr = append(tcrr, net.ParseIP(s.ip6))
	}
	rlog.Debugf("Addrs for %s in TO: %v", s.name, tcrr)
	//tcrr = append(tcrr, net.ParseIP("2001:db8:a:b::1"))
	//dnsrr = append(dnsrr, net.ParseIP("2001:DB8:a:b:0:0:0:1"))
	iprecords, _ := net.LookupIP(s.fqdn)
	for _, ip := range iprecords {
		dnsrr = append(dnsrr, ip)
	}
	if *confForce == true {
		rlog.Info("Force failure option specified")
		dnsrr = nil
	}
	rlog.Debugf("Addrs for %s in DNS: %v", s.name, dnsrr)
	if len(tcrr) != len(dnsrr) {
		msg := "TC and DNS have different number of records"
		s.failcount = 1
		return msg
	}
	for _, addr := range tcrr {
		rlog.Debugf("checking if %v is in DNS... ", addr)
		if !contains(dnsrr, addr) {
			rlog.Debug("no")
			msg := "expected A or AAAA record not found in DNS: " + addr.String()
			s.failcount = 1
			return msg
		}
	}
	if *confPtr == true {
		rlog.Infof("checking PTR records for %v", tcrr)
		if *confSyslog {
			log.Printf("checking PTR records for %v", tcrr)
		}
		var res *net.Resolver
		res = &net.Resolver{
			PreferGo: true,
		}
		for _, ip_addr := range tcrr {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel() // make sure all paths cancel the context to avoid context leak
			ptr_records, err := res.LookupAddr(ctx, ip_addr.String())
			if err != nil {
				msg := err.Error()
				s.failcount = 1
				return msg
			}
			rlog.Debugf("%s PTR: %v", ip_addr, ptr_records)
			if len(ptr_records) > 1 {
				msg := "too many PTR records found for " + ip_addr.String()
				s.failcount = 1
				return msg
			}
			ptr := strings.Trim(ptr_records[0], ".")
			ptr = strings.ToLower(ptr)
			fqdn := strings.ToLower(s.fqdn)
			if ptr != fqdn {
				msg := "unexpected PTR found in DNS: " + ptr
				s.failcount = 1
				return msg
			}
		}
	}
	return "ok"
}

func main() {
	var cpath_new string

	jobStart := time.Now()

	// define default config file path
	cpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		rlog.Error("Config error:", err)
		os.Exit(1)
	}
	cpath_new = strings.Replace(cpath, "/bin/checks", "/conf/check-config.json", 1)

	// command-line flags
	confConf := flag.String("conf", cpath_new, "Config file path")
	confName := flag.String("name", "FQDN", "Check name to pass to TO, e.g. 'FQDN'")
	confInclude := flag.String("host", "undef", "Specific host or regex to include (optional)")
	confSyslog = flag.Bool("syslog", false, "Log check results to syslog")
	confCdn := flag.String("cdn", "all", "Check specific CDN by name")
	confExclude := flag.String("exclude", "undef", "Hostname regex to exclude")
	confReset := flag.Bool("reset", false, "Reset check values in TO to 'blank' state")
	confQuiet := flag.Bool("q", false, "Do not send updates to TO")
	confForce = flag.Bool("f", false, "Force a failure result")
	confPtr = flag.Bool("ptr", false, "Validate DNS PTR record(s)")
	flag.Parse()

	// configure syslog logger
	if *confSyslog == true {
		logwriter, err := syslog.New(syslog.LOG_INFO, os.Args[0])
		if err == nil {
			log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
			log.SetOutput(logwriter)
		}
	}

	if *confName == "undef" {
		rlog.Error("Must specify check name for update to send to TO")
		os.Exit(1)
	}

	// load config json
	config, err := LoadConfig(*confConf)
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
	servers, _, err = session.GetServers()
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
			s := NewServer(server.ID, server.HostName, server.Status, -1)
			defaulStatusValue := -1
			var statusData tc.ServercheckRequestNullable
			var msg string
			statusData.ID = &s.id
			statusData.Name = confName
			statusData.HostName = &s.name
			statusData.Value = &defaulStatusValue
			s.fqdn = s.name + "." + server.DomainName
			rlog.Infof("Next server=%s status=%s", s.fqdn, s.status)
			if *confSyslog {
				log.Printf("Next server=%s status=%s", s.fqdn, s.status)
			}
			if (s.status == "REPORTED" || s.status == "ADMIN_DOWN") && *confReset != true {
				msg = s.check_dns(server)
			}

			// send status update to TO
			if s.failcount == -1 {
				// server not checked
				*statusData.Value = -1
			} else if s.failcount > 0 {
				// server had failures
				rlog.Infof("result=failure server=%s status=%s error=%s", s.fqdn, s.status, msg)
				if *confSyslog {
					log.Printf("result=failure server=%s status=%s error=%s", s.fqdn, s.status, msg)
				}
				*statusData.Value = 0
			} else {
				// server looks OK
				rlog.Infof("result=success server=%s status=%s", s.fqdn, s.status)
				if *confSyslog {
					log.Printf("result=success server=%s status=%s", s.fqdn, s.status)
				}
				*statusData.Value = 1
			}
			serverElapsed := time.Since(serverStart)
			rlog.Infof("Finished checking server=%s result=%d cdn=%s elapsed=%s", s.fqdn, *statusData.Value, s.cdn, serverElapsed)
			if *confSyslog {
				log.Printf("Finished checking server=%s result=%d cdn=%s elapsed=%s", s.fqdn, *statusData.Value, s.cdn, serverElapsed)
			}
			if *confQuiet == false {
				rlog.Debug("Sending update to TO")
				_, _, err := session.InsertServerCheckStatus(statusData)
				if err != nil {
					rlog.Error("Error updating server check status with TO:", err)
				}
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
