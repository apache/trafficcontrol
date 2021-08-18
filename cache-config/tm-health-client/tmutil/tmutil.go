package tmutil

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/cache-config/tm-health-client/config"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/datareq"
	"github.com/apache/trafficcontrol/traffic_monitor/tmclient"
	"gopkg.in/yaml.v2"
)

const (
	serverRequest  = "https://tp.cdn.comcast.net/api/3.0/servers?type=RASCAL"
	TrafficCtl     = "traffic_ctl"
	ParentsFile    = "parent.config"
	StrategiesFile = "strategies.yaml"
)

type ParentAvailable interface {
	available() bool
}

// the necessary fields of a trafficserver parents config file needed to
// read the file and keep track of it's modification time.
type ParentConfigFile struct {
	Filename       string
	LastModifyTime int64
}

// the necessary data required to keep track of trafficserver config
// files, lists of parents a trafficserver instance uses, and directory
// locations used for configuration and trafficserver executables.
type ParentInfo struct {
	ParentDotConfig        ParentConfigFile
	StrategiesDotYaml      ParentConfigFile
	TrafficServerBinDir    string
	TrafficServerConfigDir string
	Parents                map[string]ParentStatus
	Cfg                    config.Cfg
}

// when reading the 'strategies.yaml', these fields are used to help
// parse out fail_over objects.
type FailOver struct {
	MaxSimpleRetries      int      `yaml:"max_simple_retries,omitempty"`
	MaxUnavailableRetries int      `yaml:"max_unavailable_retries,omitempty"`
	RingMode              string   `yaml:"ring_mode,omitempty"`
	ResponseCodes         []int    `yaml:"response_codes,omitempty"`
	MarkDownCodes         []int    `yaml:"markdown_codes,omitempty"`
	HealthCheck           []string `yaml:"health_check,omitempty"`
}

// the trafficserver 'HostStatus' fields that are necessary to interface
// with the trafficserver 'traffic_ctl' command.
type ParentStatus struct {
	Fqdn         string
	ActiveReason bool
	LocalReason  bool
	ManualReason bool
}

// used to get the overall parent availablity from the
// HostStatus markdown reasons.  all markdown reasons
// must be true for a parent to be considered available.
func (p ParentStatus) available() bool {
	if !p.ActiveReason {
		return false
	} else if !p.LocalReason {
		return false
	} else if !p.ManualReason {
		return false
	}
	return true
}

// used to log that a parent's status is either UP or
// DOWN based upon the HostStatus reason codes.  to
// be considered UP, all reason codes must be 'true'.
func (p ParentStatus) Status() string {
	if !p.ActiveReason {
		return "DOWN"
	} else if !p.LocalReason {
		return "DOWN"
	} else if !p.ManualReason {
		return "DOWN"
	}
	return "UP"
}

type StatusReason int

// these are the HostStatus reason codes used withing
// trafficserver.
const (
	ACTIVE StatusReason = iota
	LOCAL
	MANUAL
)

// used for logging a parent's HostStatus reason code
// setting.
func (s StatusReason) String() string {
	switch s {
	case ACTIVE:
		return "ACTIVE"
	case LOCAL:
		return "LOCAL"
	case MANUAL:
		return "MANUAL"
	}
	return "UNDEFINED"
}

// the fields used from 'strategies.yaml' that describe
// a parent.
type Host struct {
	HostName  string     `yaml:"host"`
	Protocols []Protocol `yaml:"protocol"`
}

// the protocol object in 'strategies.yaml' that help to
// describe a parent.
type Protocol struct {
	Scheme           string  `yaml:"scheme"`
	Port             int     `yaml:"port"`
	Health_check_url string  `yaml:"health_check_url,omitempty"`
	Weight           float64 `yaml:"weight,omitempty"`
}

// a trafficserver strategy object from 'strategies.yaml'.
type Strategy struct {
	Strategy        string   `yaml:"strategy"`
	Policy          string   `yaml:"policy"`
	HashKey         string   `yaml:"hash_key,omitempty"`
	GoDirect        bool     `yaml:"go_direct,omitempty"`
	ParentIsProxy   bool     `yaml:"parent_is_proxy,omitempty"`
	CachePeerResult bool     `yaml:"cache_peer_result,omitempty"`
	Scheme          string   `yaml:"scheme"`
	FailOvers       FailOver `yaml:"failover,omitempty"`
}

// the top level array defintions in a trafficserver 'strategies.yaml'
// configuration file.
type Strategies struct {
	Strategy []Strategy    `yaml:"strategies"`
	Hosts    []Host        `yaml:"hosts"`
	Groups   []interface{} `yaml:"groups"`
}

// used at startup to load a trafficservers list of parents from
// it's 'parent.config', 'strategies.yaml' and current parent
// status from trafficservers HostStatus subsystem.
func NewParentInfo(cfg config.Cfg) (*ParentInfo, error) {

	parentConfig := filepath.Join(cfg.TrafficServerConfigDir, ParentsFile)
	modTime, err := getFileModificationTime(parentConfig)
	if err != nil {
		return nil, errors.New("error reading " + ParentsFile + ": " + err.Error())
	}
	parents := ParentConfigFile{
		Filename:       parentConfig,
		LastModifyTime: modTime,
	}

	stratyaml := filepath.Join(cfg.TrafficServerConfigDir, StrategiesFile)
	modTime, err = getFileModificationTime(stratyaml)
	if err != nil {
		return nil, errors.New("error reading " + StrategiesFile + ": " + err.Error())
	}

	strategies := ParentConfigFile{
		Filename:       filepath.Join(cfg.TrafficServerConfigDir, StrategiesFile),
		LastModifyTime: modTime,
	}

	parentInfo := ParentInfo{
		ParentDotConfig:        parents,
		StrategiesDotYaml:      strategies,
		TrafficServerBinDir:    cfg.TrafficServerBinDir,
		TrafficServerConfigDir: cfg.TrafficServerConfigDir,
	}

	// initialize the trafficserver parents map.
	parentStatus := make(map[string]ParentStatus)

	// read the 'parent.config'.
	if err := parentInfo.readParentConfig(parentStatus); err != nil {
		return nil, errors.New("loading " + ParentsFile + " file: " + err.Error())
	}

	// read the strategies.yaml.
	if err := parentInfo.readStrategies(parentStatus); err != nil {
		return nil, errors.New("loading parent " + StrategiesFile + " file: " + err.Error())
	}

	// collect the trafficserver parent status from the HostStatus subsystem.
	if err := parentInfo.readHostStatus(parentStatus); err != nil {
		return nil, errors.New("reading trafficserver host status: " + err.Error())
	}

	log.Infof("startup loaded %d parent records\n", len(parentStatus))

	parentInfo.Parents = parentStatus
	parentInfo.Cfg = cfg

	return &parentInfo, nil
}

// Queries a traffic monitor that is monitoring the trafficserver instance running on a host to
// obtain the availability, health, of a parent used by trafficserver.
func (c *ParentInfo) GetCacheStatuses() (map[tc.CacheName]datareq.CacheStatus, error) {

	tmHostName, err := c.findATrafficMonitor()
	if err != nil {
		return nil, errors.New("finding a trafficmonitor: " + err.Error())
	}
	tmc := tmclient.New("http://"+tmHostName, config.GetRequestTimeout())

	return tmc.CacheStatuses()
}

// The main polling function that keeps the parents list current if
// with any changes to the trafficserver 'parent.config' or 'strategies.yaml'.
// Also, it keeps parent status current with the the trafficserver HostStatus
// subsystem.  Finally, on each poll cycle a trafficmonitor is queried to check
// that all parents used by this trafficserver are available for use based upon
// the trafficmonitors idea from it's health protocol.  Parents are marked up or
// down in the trafficserver subsystem based upon that hosts current status and
// the status that trafficmonitor health protocol has determined for a parent.
func (c *ParentInfo) PollAndUpdateCacheStatus() {
	pollingInterval := config.GetTMPollingInterval()
	log.Infoln("polling started")

	for {
		if err := c.UpdateParentInfo(); err != nil {
			log.Errorf("could not load new ATS parent info: %s\n", err.Error())
		} else {
			log.Debugf("updated parent info, total number of parents: %d\n", len(c.Parents))
		}

		// read traffic manager cache statuses.
		caches, err := c.GetCacheStatuses()
		if err != nil {
			log.Errorln(err.Error())
		}

		for k, v := range caches {
			hostName := string(k)
			cs, ok := c.Parents[hostName]
			if ok {
				tmAvailable := *v.CombinedAvailable
				if cs.available() != tmAvailable {
					if !c.Cfg.EnableActiveMarkdowns {
						if !tmAvailable {
							log.Infof("TM reports that %s is not available and should be marked DOWN but, mark downs are disabled by configuration", hostName)
						} else {
							log.Infof("TM reports that %s is available and should be marked UP but, mark up is disabled by configuration", hostName)
						}
					} else {
						if err = c.markParent(cs.Fqdn, *v.Status, tmAvailable); err != nil {
							log.Errorln(err.Error())
						}
					}
				}
			}
		}

		time.Sleep(pollingInterval)
	}
}

// Used by the polling function to update the parents list from
// changes to 'parent.config' and 'strategies.yaml'.  The parents
// availability is also updated to reflect the current state from
// the trafficserver HostStatus subsystem.
func (c *ParentInfo) UpdateParentInfo() error {
	ptime, err := getFileModificationTime(c.ParentDotConfig.Filename)
	if err != nil {
		return errors.New("error reading " + ParentsFile + ": " + err.Error())
	}
	stime, err := getFileModificationTime(c.StrategiesDotYaml.Filename)
	if err != nil {
		return errors.New("error reading " + StrategiesFile + ": " + err.Error())
	}
	if c.ParentDotConfig.LastModifyTime < ptime {
		// read the 'parent.config'.
		if err := c.readParentConfig(c.Parents); err != nil {
			return errors.New("updating " + ParentsFile + " file: " + err.Error())
		} else {
			log.Infof("updated parents from new %s, total parents: %d\n", ParentsFile, len(c.Parents))
		}
	}

	if c.StrategiesDotYaml.LastModifyTime < stime {
		// read the 'strategies.yaml'.
		if err := c.readStrategies(c.Parents); err != nil {
			return errors.New("updating parent " + StrategiesFile + " file: " + err.Error())
		} else {
			log.Infof("updated parents from new %s total parents: %d\n", StrategiesFile, len(c.Parents))
		}
	}

	// collect the trafficserver current host status.
	if err := c.readHostStatus(c.Parents); err != nil {
		return errors.New("reading latest trafficserver host status: " + err.Error())
	}

	return nil
}

// choose an available trafficmonitor, returns an error if
// there are none.
func (c *ParentInfo) findATrafficMonitor() (string, error) {
	var tmHostname string
	lth := len(c.Cfg.TrafficMonitors)
	if lth == 0 {
		return "", errors.New("there are no available traffic monitors")
	}

	// build an array of available traffic monitors.
	tms := make([]string, 0)
	for k, v := range c.Cfg.TrafficMonitors {
		if v == true {
			log.Debugf("traffic monitor %s is available\n", k)
			tms = append(tms, k)
		}
	}

	// choose one at random.
	lth = len(tms)
	if lth > 0 {
		rand.Seed(time.Now().UnixNano())
		r := (rand.Intn(lth))
		tmHostname = tms[r]
	} else {
		return "", errors.New("there are no available traffic monitors")
	}

	log.Infof("polling: %s\n", tmHostname)

	return tmHostname, nil
}

// get the file modification times for a trafficserver configuration
// file.
func getFileModificationTime(fn string) (int64, error) {
	f, err := os.Open(fn)
	if err != nil {
		return 0, errors.New("opening " + fn + ": " + err.Error())
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return 0, errors.New("unable to get file status for " + fn + ": " + err.Error())
	}
	return finfo.ModTime().UnixNano(), nil
}

// parse out the hostname of a parent listed in parents.config
// or 'strategies.yaml'. the hostname can be an IP address.
func parseFqdn(fqdn string) string {
	var hostName string
	if ip := net.ParseIP(fqdn); ip == nil {
		// not an IP, get the hostname
		flds := strings.Split(fqdn, ".")
		hostName = flds[0]
	} else { // use the IP addr
		hostName = fqdn
	}
	return hostName
}

// used to mark a parent as up or down in the trafficserver HostStatus
// subsystem.
func (c *ParentInfo) markParent(fqdn string, cacheStatus string, available bool) error {
	hostName := parseFqdn(fqdn)
	tc := filepath.Join(c.TrafficServerBinDir, TrafficCtl)
	reason := c.Cfg.ReasonCode
	var status string
	if available {
		status = "up"
	} else {
		status = "down"
	}

	cmd := exec.Command(tc, "host", status, "--reason", reason, fqdn)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New("marking " + fqdn + " " + status + ": " + TrafficCtl + " error: " + err.Error())
	}
	pv, ok := c.Parents[hostName]
	if ok {
		switch reason {
		case "active":
			pv.ActiveReason = available
		case "local":
			pv.LocalReason = available
		}
	}
	c.Parents[hostName] = pv

	if !available {
		log.Infof("marked parent %s DOWN, cache status was: %s\n", hostName, cacheStatus)
	} else {
		log.Infof("marked parent %s UP, cache status was: %s\n", hostName, cacheStatus)
	}
	return nil
}

// reads the current parent statuses from the trafficserver HostStatus
// subsystem.
func (c *ParentInfo) readHostStatus(parentStatus map[string]ParentStatus) error {
	tc := filepath.Join(c.TrafficServerBinDir, TrafficCtl)

	cmd := exec.Command(tc, "metric", "match", "host_status")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(TrafficCtl + " error: " + err.Error())
	}

	if len((stdout.Bytes())) > 0 {
		var activeReason bool
		var localReason bool
		var manualReason bool
		var hostName string
		var fqdn string
		scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "proxy.process.host_status.") {
				fields := strings.Split(line, " ")
				if len(fields) > 0 {
					fqdnField := strings.Split(fields[0], "proxy.process.host_status.")
					if len(fqdnField) > 0 {
						fqdn = fqdnField[1]
					}
					statField := strings.Split(fields[1], ",")
					if len(statField) == 5 {
						if strings.HasPrefix(statField[1], "ACTIVE:UP") {
							activeReason = true
						} else if strings.HasPrefix(statField[1], "ACTIVE:DOWN") {
							activeReason = false
						}
						if strings.HasPrefix(statField[2], "LOCAL:UP") {
							localReason = true
						} else if strings.HasPrefix(statField[2], "LOCAL:DOWN") {
							localReason = false
						}
						if strings.HasPrefix(statField[3], "MANUAL:UP") {
							manualReason = true
						} else if strings.HasPrefix(statField[3], "MANUAL:DOWN") {
							manualReason = false
						}
					}
					pstat := ParentStatus{
						Fqdn:         fqdn,
						ActiveReason: activeReason,
						LocalReason:  localReason,
						ManualReason: manualReason,
					}
					hostName = parseFqdn(fqdn)
					pv, ok := parentStatus[hostName]
					// create the ParentStatus struct and add it to the
					// Parents map only if an entry in the map does not
					// already exist.
					if !ok {
						parentStatus[hostName] = pstat
						log.Infof("added Host '%s' from ATS Host Status to the parents map\n", hostName)
					} else {
						available := pstat.available()
						if pv.available() != available {
							log.Infof("host status for '%s' has changed to %s\n", hostName, pstat.Status())
							parentStatus[hostName] = pstat
						}
					}
				}
			}
		}
		log.Debugf("processed trafficserver host status results, total parents: %d\n", len(parentStatus))
	}
	return nil
}

// load parents list from the Trafficserver 'parent.config' file.
func (c *ParentInfo) readParentConfig(parentStatus map[string]ParentStatus) error {
	fn := c.ParentDotConfig.Filename

	_, err := os.Stat(fn)
	if err != nil {
		log.Warnf("skipping 'parents': %s\n", err.Error())
		return nil
	}

	log.Debugf("loading %s\n", fn)

	f, err := os.Open(fn)

	if err != nil {
		return errors.New("failed to open + " + fn + " :" + err.Error())
	}
	defer f.Close()

	finfo, err := os.Stat(fn)
	if err != nil {
		return errors.New("failed to Stat + " + fn + " :" + err.Error())
	}
	c.ParentDotConfig.LastModifyTime = finfo.ModTime().UnixNano()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		sbytes := scanner.Bytes()
		if sbytes[0] == 35 { // skip comment lines, 35 is a '#'.
			continue
		}
		// search for the parent list.
		if i := strings.Index(string(sbytes), "parent="); i > 0 {
			var plist []string
			res := bytes.Split(sbytes, []byte("\""))
			// 'parent.config' parent separators are ';' or ','.
			plist = strings.Split(strings.TrimSpace(string(res[1])), ";")
			if len(plist) == 1 {
				plist = strings.Split(strings.TrimSpace(string(res[1])), ",")
			}
			// parse the parent list to get each hostName and it's associated
			// port.
			if len(plist) > 1 {
				for _, v := range plist {
					parent := strings.Split(v, ":")
					if len(parent) == 2 {
						fqdn := parent[0]
						hostName := parseFqdn(fqdn)
						_, ok := parentStatus[hostName]
						// create the ParentStatus struct and add it to the
						// Parents map only if an entry in the map does not
						// already exist.
						if !ok {
							pstat := ParentStatus{
								Fqdn:         strings.TrimSpace(fqdn),
								ActiveReason: true,
								LocalReason:  true,
								ManualReason: true,
							}
							parentStatus[hostName] = pstat
							log.Debugf("added Host '%s' from %s to the parents map\n", hostName, fn)
						}
					}
				}
			}
		}
	}
	return nil
}

// load the parent hosts from 'strategies.yaml'.
func (c *ParentInfo) readStrategies(parentStatus map[string]ParentStatus) error {
	var includes []string
	fn := c.StrategiesDotYaml.Filename

	_, err := os.Stat(fn)
	if err != nil {
		log.Warnf("skipping 'strategies': %s\n", err.Error())
		return nil
	}

	log.Debugf("loading %s\n", fn)

	// open the strategies file for scanning.
	f, err := os.Open(fn)
	if err != nil {
		return errors.New("failed to open + " + fn + " :" + err.Error())
	}
	defer f.Close()

	finfo, err := os.Stat(fn)
	if err != nil {
		return errors.New("failed to Stat + " + fn + " :" + err.Error())
	}
	c.StrategiesDotYaml.LastModifyTime = finfo.ModTime().UnixNano()

	scanner := bufio.NewScanner(f)

	// search for any yaml files that should be included in the
	// yaml stream.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#include") {
			fields := strings.Split(line, " ")
			if len(fields) >= 2 {
				includeFile := filepath.Join(c.TrafficServerConfigDir, fields[1])
				includes = append(includes, includeFile)
			}
		}
	}

	includes = append(includes, fn)

	var yamlContent string

	// load all included and 'strategies yaml' files to
	// the yamlContent.
	for _, includeFile := range includes {
		log.Debugf("loading %s\n", includeFile)
		content, err := ioutil.ReadFile(includeFile)
		if err != nil {
			return errors.New(err.Error())
		}

		yamlContent = yamlContent + string(content)
	}

	strategies := Strategies{}

	if err := yaml.Unmarshal([]byte(yamlContent), &strategies); err != nil {
		return errors.New("failed to unmarshall " + fn + ": " + err.Error())
	}

	for _, host := range strategies.Hosts {
		fqdn := host.HostName
		hostName := parseFqdn(fqdn)
		// create the ParentStatus struct and add it to the
		// Parents map only if an entry in the map does not
		// already exist.
		_, ok := parentStatus[hostName]
		if !ok {
			pstat := ParentStatus{
				Fqdn:         strings.TrimSpace(fqdn),
				ActiveReason: true,
				LocalReason:  true,
				ManualReason: true,
			}
			parentStatus[hostName] = pstat
			log.Debugf("added Host '%s' from %s to the parents map\n", hostName, fn)
		}
	}
	return nil
}
