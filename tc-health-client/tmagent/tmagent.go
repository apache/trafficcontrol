package tmagent

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
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/tc-health-client/config"
	"github.com/apache/trafficcontrol/tc-health-client/util"
	"github.com/apache/trafficcontrol/traffic_monitor/tmclient"
	"gopkg.in/yaml.v2"
)

const (
	TrafficCtl     = "traffic_ctl"
	ParentsFile    = "parent.config"
	StrategiesFile = "strategies.yaml"
)

type ParentAvailable interface {
	available(reasonCode string) bool
}

// the necessary data required to keep track of trafficserver config
// files, lists of parents a trafficserver instance uses, and directory
// locations used for configuration and trafficserver executables.
type ParentInfo struct {
	ParentDotConfig        util.ConfigFile
	StrategiesDotYaml      util.ConfigFile
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
	Fqdn                 string
	ActiveReason         bool
	LocalReason          bool
	ManualReason         bool
	LastTmPoll           int64
	UnavailablePollCount int
}

// used to get the overall parent availablity from the
// HostStatus markdown reasons.  all markdown reasons
// must be true for a parent to be considered available.
func (p ParentStatus) available(reasonCode string) bool {
	rc := false

	switch reasonCode {
	case "active":
		rc = p.ActiveReason
	case "local":
		rc = p.LocalReason
	case "manual":
		rc = p.ManualReason
	}
	return rc
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
	modTime, err := util.GetFileModificationTime(parentConfig)
	if err != nil {
		return nil, errors.New("error reading " + ParentsFile + ": " + err.Error())
	}
	parents := util.ConfigFile{
		Filename:       parentConfig,
		LastModifyTime: modTime,
	}

	stratyaml := filepath.Join(cfg.TrafficServerConfigDir, StrategiesFile)
	modTime, err = util.GetFileModificationTime(stratyaml)
	if err != nil {
		return nil, errors.New("error reading " + StrategiesFile + ": " + err.Error())
	}

	strategies := util.ConfigFile{
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
		return nil, fmt.Errorf("reading trafficserver host status: %w", err)
	}

	log.Infof("startup loaded %d parent records\n", len(parentStatus))

	parentInfo.Parents = parentStatus
	parentInfo.Cfg = cfg

	return &parentInfo, nil
}

// Queries a traffic monitor that is monitoring the trafficserver instance running on a host to
// obtain the availability, health, of a parent used by trafficserver.
func (c *ParentInfo) GetCacheStatuses() (tc.CRStates, error) {

	tmHostName, err := c.findATrafficMonitor()
	if err != nil {
		return tc.CRStates{}, errors.New("finding a trafficmonitor: " + err.Error())
	}
	tmc := tmclient.New("http://"+tmHostName, config.GetRequestTimeout())

	// Use a proxy to query TM if the ProxyURL is set
	if c.Cfg.ParsedProxyURL != nil {
		tmc.Transport = &http.Transport{Proxy: http.ProxyURL(c.Cfg.ParsedProxyURL)}
	}

	return tmc.CRStates(true)
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
	cycleCount := 0
	log.Infoln("polling started")

	for {
		pollingInterval := config.GetTMPollingInterval()
		// check for config file updates
		newCfg := config.Cfg{
			HealthClientConfigFile: c.Cfg.HealthClientConfigFile,
		}
		isNew, err := config.LoadConfig(&newCfg)
		if err != nil {
			log.Errorf("error reading changed config file %s: %s\n", c.Cfg.HealthClientConfigFile.Filename, err.Error())
		}
		if isNew {
			if err = config.ReadCredentials(&newCfg, false); err != nil {
				log.Errorln("could not load credentials for config updates, keeping the old config")
			} else {
				if err = config.GetTrafficMonitors(&newCfg); err != nil {
					log.Errorln("could not update the list of trafficmonitors, keeping the old config")
				} else {
					config.UpdateConfig(&c.Cfg, &newCfg)
					log.Infoln("the configuration has been successfully updated")
				}
			}
		} else { // check for updates to the credentials file
			if c.Cfg.CredentialFile.Filename != "" {
				modTime, err := util.GetFileModificationTime(c.Cfg.CredentialFile.Filename)
				if err != nil {
					log.Errorf("could not stat the credential file %s", c.Cfg.CredentialFile.Filename)
				}
				if modTime > c.Cfg.CredentialFile.LastModifyTime {
					log.Infoln("the credential file has changed, loading new credentials")
					if err = config.ReadCredentials(&c.Cfg, true); err != nil {
						log.Errorf("could not load credentials from the updated credential file: %s", c.Cfg.CredentialFile.Filename)
					}
				}
			}
		}

		// check for parent and strategies config file updates, and trafficserver
		// host status changes.  If an error is encountered reading data the current
		// parents lists and hoststatus remains unchanged.
		if err := c.UpdateParentInfo(); err != nil {
			log.Errorf("could not load new ATS parent info: %s\n", err.Error())
		} else {
			log.Debugf("updated parent info, total number of parents: %d\n", len(c.Parents))
		}

		// read traffic manager cache statuses.
		_c, err := c.GetCacheStatuses()

		// get the current poll time
		now := time.Now().Unix()

		caches := _c.Caches
		if err != nil {
			log.Errorf("error in TrafficMonitor polling: %s\n", err.Error())
			if err = config.GetTrafficMonitors(&c.Cfg); err != nil {
				log.Errorln("could not update the list of trafficmonitors, keeping the old config")
			} else {
				log.Infoln("updated TrafficMonitor statuses from TrafficOps")
			}
			time.Sleep(pollingInterval)
			continue
		}

		for k, v := range caches {
			hostName := string(k)
			cs, ok := c.Parents[hostName]
			if ok {
				// update the polling time
				cs.LastTmPoll = now
				c.Parents[hostName] = cs
				tmAvailable := v.IsAvailable
				if cs.available(c.Cfg.ReasonCode) != tmAvailable {
					// do not mark down if the configuration disables mark downs.
					if !c.Cfg.EnableActiveMarkdowns && !tmAvailable {
						log.Infof("TM reports that %s is not available and should be marked DOWN but, mark downs are disabled by configuration", hostName)
					} else {
						// See issue #6448, the status field used in api/cache-status is not
						// available in the publish/CrStates endpoint.  For now, will not
						// use it.
						//if err = c.markParent(cs.Fqdn, *v.Status, tmAvailable); err != nil {
						if err = c.markParent(cs.Fqdn, tmAvailable); err != nil {
							log.Errorln(err.Error())
						}
					}
				}
				// if the host is available clear the unavailable poll count if not 0.
				if cs.available(c.Cfg.ReasonCode) && tmAvailable {
					if cs.UnavailablePollCount > 0 {
						log.Debugf("resetting the UnavailablePollCount for %s from %d to 0",
							hostName, cs.UnavailablePollCount)
						cs.UnavailablePollCount = 0
						c.Parents[hostName] = cs
					}
				}
			}
		}

		// periodically update the TrafficMonitor list and statuses
		if cycleCount == c.Cfg.TmUpdateCycles {
			cycleCount = 0
			if err = config.GetTrafficMonitors(&c.Cfg); err != nil {
				log.Errorln("could not update the list of trafficmonitors, keeping the old config")
			} else {
				log.Infoln("updated TrafficMonitor statuses from TrafficOps")
			}
		} else {
			cycleCount++
		}

		time.Sleep(pollingInterval)
	}
}

// Used by the polling function to update the parents list from
// changes to 'parent.config' and 'strategies.yaml'.  The parents
// availability is also updated to reflect the current state from
// the trafficserver HostStatus subsystem.
func (c *ParentInfo) UpdateParentInfo() error {
	ptime, err := util.GetFileModificationTime(c.ParentDotConfig.Filename)
	if err != nil {
		return errors.New("error reading " + ParentsFile + ": " + err.Error())
	}
	stime, err := util.GetFileModificationTime(c.StrategiesDotYaml.Filename)
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
		return errors.New("trafficserver may not be running: " + err.Error())
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

	log.Debugf("polling: %s\n", tmHostname)

	return tmHostname, nil
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

func (c *ParentInfo) execTrafficCtl(fqdn string, available bool) error {
	reason := c.Cfg.ReasonCode
	tc := filepath.Join(c.TrafficServerBinDir, TrafficCtl)

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

	return nil
}

// used to mark a parent as up or down in the trafficserver HostStatus
// subsystem.
//
// TODO see issue #6448, add cacheStatus back when available in CrStates
//func (c *ParentInfo) markParent(fqdn string, cacheStatus string, available bool) error {
func (c *ParentInfo) markParent(fqdn string, available bool) error {
	var hostAvailable bool
	var err error
	hostName := parseFqdn(fqdn)

	log.Debugf("fqdn: %s, available: %v", fqdn, available)

	pv, ok := c.Parents[hostName]
	if ok {
		activeReason := pv.ActiveReason
		localReason := pv.LocalReason
		unavailablePollCount := pv.UnavailablePollCount

		log.Debugf("hostName: %s, UnavailablePollCount: %d, available: %v", hostName, unavailablePollCount, available)

		if !available { // unavailable
			unavailablePollCount += 1
			if unavailablePollCount < c.Cfg.UnavailablePollThreshold {
				log.Infof("TM indicates %s is unavailable but the UnavailablePollThreshold has not been reached", hostName)
				hostAvailable = true
			} else {
				// marking the host down
				err = c.execTrafficCtl(fqdn, available)
				if err != nil {
					log.Errorln(err.Error())
				}
				if err == nil {
					// TODO see issue 6448, add cacheStatus back when available in CrStates
					// log.Infof("marked parent %s DOWN, cache status was: %s\n", hostName, cacheStatus)
					hostAvailable = false
					log.Infof("marked parent %s DOWN", hostName)
				}
			}
		} else { // available
			// marking the host up
			err = c.execTrafficCtl(fqdn, available)
			if err == nil {
				hostAvailable = true
				// reset the unavilable poll count
				unavailablePollCount = 0
				// TODO see issue #6448, add cacheStatus back when available in CrStates
				//log.Infof("marked parent %s UP, cache status was: %s\n", hostName, cacheStatus)
				log.Infof("marked parent %s UP", hostName)
			} else {
				hostAvailable = false
			}
		}

		// update parent info
		if err == nil {
			reason := c.Cfg.ReasonCode
			switch reason {
			case "active":
				activeReason = hostAvailable
			case "local":
				localReason = hostAvailable
			}
			// save updates
			pv.ActiveReason = activeReason
			pv.LocalReason = localReason
			pv.UnavailablePollCount = unavailablePollCount
			c.Parents[hostName] = pv
			log.Debugf("Updated parent status: %v", pv)
		}
	}
	return err
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
		return fmt.Errorf("%s error: %s", TrafficCtl, stderr.String())
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
						Fqdn:                 fqdn,
						ActiveReason:         activeReason,
						LocalReason:          localReason,
						ManualReason:         manualReason,
						LastTmPoll:           0,
						UnavailablePollCount: 0,
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
						available := pstat.available(c.Cfg.ReasonCode)
						if pv.available(c.Cfg.ReasonCode) != available {
							log.Infof("host status for '%s' has changed to %s\n", hostName, pstat.Status())
							pstat.LastTmPoll = pv.LastTmPoll
							pstat.UnavailablePollCount = pv.UnavailablePollCount
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
								Fqdn:                 strings.TrimSpace(fqdn),
								ActiveReason:         true,
								LocalReason:          true,
								ManualReason:         true,
								LastTmPoll:           0,
								UnavailablePollCount: 0,
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
				Fqdn:                 strings.TrimSpace(fqdn),
				ActiveReason:         true,
				LocalReason:          true,
				ManualReason:         true,
				LastTmPoll:           0,
				UnavailablePollCount: 0,
			}
			parentStatus[hostName] = pstat
			log.Debugf("added Host '%s' from %s to the parents map\n", hostName, fn)
		}
	}
	return nil
}
