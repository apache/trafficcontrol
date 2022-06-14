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

/* ToPingCheck.go
   Used for checking ILO ping,  MTU test, 10G (IPv4), and 10G6 (IPv6) pings.
*/

package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
)

// Traffic Ops connection params
const AllowInsecureConnections = false
const UserAgent = "go/tc-ping-monitor"
const UseClientCache = false
const TrafficOpsRequestTimeout = time.Second * time.Duration(10)

var confForce *bool

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
	ilo       string
	ip4       string
	ip6       string
	mtu       int
	failcount int
	status    string
	failaddr  string
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

func (s *Server) ping(name string) bool {
	// I haven't found a good way to do this natively with go yet, which makes me sad
	s.failcount = 0
	size := 30 //default
	ok := true
	var addr string
	switch name {
	case "IPv4", "10G":
		log.Infof("IPv4")
		if s.ip4 != "" {
			ok, addr = ping4(size, s.ip4)
			if ok == false {
				s.failaddr = addr
			}
		}
	case "IPv6", "10G6":
		log.Infof("IPv6")
		if s.ip6 != "" {
			ok, addr = ping6(size, s.ip6)
			if ok == false {
				s.failaddr = addr
			}
		}
	case "ILO":
		log.Infof("ILO")
		if s.ilo != "" {
			match, err := regexp.MatchString(":", s.ilo)
			if err != nil {
				log.Errorf("Match error:", err)
				os.Exit(1)
			}
			if match {
				ok, addr = ping6(size, s.ilo)
			} else {
				ok, addr = ping4(size, s.ilo)
				if ok == false {
					s.failaddr = addr
				}
			}
		}
	case "MTU":
		log.Infof("MTU")
		ok4 := true
		if s.ip4 != "" {
			// subtract protocol headers from MTU to get payload size
			size = s.mtu - 28
			ok4, addr = ping4(size, s.ip4)
		} else {
			log.Warnf("no IPv4 address detected (skipping)")
		}
		if ok4 == false {
			s.failaddr = addr
		}
		ok6 := true
		if s.ip6 != "" {
			// subtract protocol headers from MTU to get payload size
			size = s.mtu - 48
			ok6, addr = ping6(size, s.ip6)
		} else {
			log.Warnf("no IPv6 address detected (skipping)")
		}
		if ok6 == false {
			if len(s.failaddr) > 0 {
				s.failaddr = s.failaddr + "," + addr
			} else {
				s.failaddr = addr
			}
		}
		if !(ok4 && ok6) {
			ok = false
		}
	}
	return ok
}

func ping4(size int, addr string) (bool, string) {
	log.Infof("size: ", strconv.Itoa(size))
	out, err := exec.Command("/bin/ping", "-M", "do", "-s", strconv.Itoa(size), "-c", "2", addr).Output()
	if err != nil {
		log.Warnf("ping failed for %s: %s", addr, err.Error())
		return false, addr
	}
	log.Debugf("ping output:\n%v", out)
	return true, addr
}

func ping6(size int, addr string) (bool, string) {
	log.Infof("size: ", strconv.Itoa(size))
	out, err := exec.Command("/bin/ping6", "-M", "do", "-s", strconv.Itoa(size), "-c", "2", addr).Output()
	if err != nil {
		log.Warnf("ping failed for %s: %s", addr, err.Error())
		return false, addr
	}
	log.Debugf("ping output:\n%v", out)
	return true, addr
}

func main() {
	var cpath_new string
	var ok bool

	jobStart := time.Now()

	// define default config file path
	cpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Errorf("Config error:", err)
		os.Exit(1)
	}
	cpath_new = strings.Replace(cpath, "/bin/checks", "/conf/check-config.json", 1)

	// command-line flags
	confPtr := flag.String("conf", cpath_new, "Config file path")
	confName := flag.String("name", "undef", "'10G|IPv4', '10G6|IPv6', 'ILO', 'MTU'")
	confInclude := flag.String("host", "undef", "Specific host or regex to include (optional)")
	confCdn := flag.String("cdn", "all", "Check specific CDN by name")
	confExclude := flag.String("exclude", "undef", "Hostname regex to exclude")
	confReset := flag.Bool("reset", false, "Reset check values in TO to 'blank' state")
	confQuiet := flag.Bool("q", false, "Do not send updates to TO")
	confForce = flag.Bool("f", false, "Force a failure result")
	flag.Parse()

	reName, err := regexp.Compile("^(10G|10G6|IPv4|IPv6|ILO|MTU)$")
	if err != nil {
		log.Errorf("supplied exclusion regex does not compile:", err)
		os.Exit(1)
	}
	if !(reName.Match([]byte(*confName))) {
		log.Errorf("Check name must be one of the following:")
		log.Errorf("'10G' (legacy) or 'IPv4' (new) for IPv4 interface check")
		log.Errorf("'10G6' (legacy) or 'IPv6' (new) for IPv6 interface check")
		log.Errorf("'ILO' out-of-band mgmt interface check")
		log.Errorf("'MTU' uses the MTU value for the server in TO to check MTU (checks both v4 and v6, if available)")
		os.Exit(1)
	}

	// load config json
	config, err := LoadConfig(*confPtr)
	if err != nil {
		log.Errorf("Error loading config:", err)
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
		log.Errorf("An error occurred while logging in: %v\n", err)
		os.Exit(1)
	}

	// Make TO API call for server details
	var servers tc.ServersV3Response
	servers, _, err = session.GetServersWithHdr(nil, nil)
	if err != nil {
		log.Errorf("An error occurred while getting servers: %v\n", err)
		os.Exit(1)
	}

	for _, server := range servers.Response {
		re, err := regexp.Compile("^(MID|EDGE).*")
		if err != nil {
			log.Errorf("supplied exclusion regex does not compile:", err)
			os.Exit(1)
		}
		if re.Match([]byte(server.Type)) {
			serverStart := time.Now()
			if *confInclude != "undef" {
				re_inc, err := regexp.Compile(*confInclude)
				if err != nil {
					log.Errorf("supplied exclusion regex does not compile:", err)
					os.Exit(1)
				}
				if !re_inc.MatchString(*server.HostName) {
					log.Debugf("%s does not match the provided include regex, skipping", server.HostName)
					continue
				}
			}
			if *confCdn != "all" && *confCdn != *server.CDNName {
				log.Debugf("%s is not assinged to the specified CDN '%s', skipping", server.HostName, *confCdn)
				continue
			}
			if *confExclude != "undef" {
				re, err := regexp.Compile(*confExclude)
				if err != nil {
					log.Errorf("supplied exclusion regex does not compile:", err)
					os.Exit(1)
				}
				if re.MatchString(*server.HostName) {
					log.Debugf("%s matches the provided exclude regex, skipping", server.HostName)
					continue
				}
			}
			s := NewServer(*server.ID, *server.HostName, *server.Status, -1)
			defaulStatusValue := -1
			var statusData tc.ServercheckRequestNullable
			statusData.ID = &s.id
			statusData.Name = confName
			statusData.HostName = &s.name
			statusData.Value = &defaulStatusValue
			s.fqdn = s.name + "." + *server.DomainName
			for _, interf := range server.Interfaces {
				for _, addr := range interf.IPAddresses {
					if s.ip4 == "" && strings.Count(addr.Address, ":") == 0 {
						s.ip4 = strings.Split(addr.Address, "/")[0]
					}
					if s.ip6 == "" && strings.Count(addr.Address, ":") > 0 {
						s.ip6 = strings.Split(addr.Address, "/")[0]
					}
				}
			}
			if server.ILOIPAddress == nil {
				s.ilo = ""
			} else {
				if strings.Contains(*server.ILOIPAddress, "/") {
					s.ilo = strings.Split(*server.ILOIPAddress, "/")[0]
				} else {
					s.ilo = *server.ILOIPAddress
				}
			}
			s.cdn = *server.CDNName
			s.mtu = int(*server.Interfaces[0].MTU)
			log.Infof("Next server=%s status=%s", s.fqdn, s.status)
			if (s.status == "REPORTED" || s.status == "ADMIN_DOWN") && *confReset != true {
				ok = s.ping(*confName)
				log.Infof("ok: %v", ok)
				if ok == false {
					s.failcount = 1
				}
			}

			// send status update to TO
			if s.failcount == -1 {
				// server not checked
				*statusData.Value = -1
			} else if s.failcount > 0 {
				// server had failures
				log.Infof("result=failure server=%s status=%s check=%s addr=%s", s.fqdn, s.status, *confName, s.failaddr)
				*statusData.Value = 0
			} else {
				// server looks OK
				log.Infof("result=success server=%s status=%s", s.fqdn, s.status)
				*statusData.Value = 1
			}
			serverElapsed := time.Since(serverStart)
			log.Infof("Finished checking server=%s result=%d cdn=%s elapsed=%s", s.fqdn, *statusData.Value, s.cdn, serverElapsed)
			if *confQuiet == false {
				log.Debugf("Sending update to TO")
				_, _, err := session.InsertServerCheckStatus(statusData)
				if err != nil {
					log.Errorf("Error updating server check status with TO:", err)
				}
			} else {
				log.Debugf("Skipping update to TO")
			}
		}
	}
	jobElapsed := time.Since(jobStart)
	log.Infof("Job complete totaltime=%s", jobElapsed)
	os.Exit(0)
}