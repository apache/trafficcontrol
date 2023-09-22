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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/tc-health-client/config"
	"github.com/apache/trafficcontrol/v8/tc-health-client/util"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"

	"gopkg.in/yaml.v3"
)

const (
	TrafficCtl     = "traffic_ctl"
	ParentsFile    = "parent.config"
	StrategiesFile = "strategies.yaml"
)

// this global is used to auto select the
// proper ATS traffic_ctl command to use
// when querying host status. for ATS
// version 10 and greater this will remain
// at 0.  For ATS version 9, this will be
// auto updated to 1
var traffic_ctl_index = 0

// ParentInfo contains the necessary data required to keep track of trafficserver config
// files, lists of parents a trafficserver instance uses, and directory
// locations used for configuration and trafficserver executables.
//
// All AtomicPtr members MUST NOT be modified.
// All AtomicPtr members are accessed concurrently by multiple goroutines,
// and their objects must be atomically replaced by a single writer, and never modified.
type ParentInfo struct {
	ParentDotConfig        util.ConfigFile
	StrategiesDotYaml      util.ConfigFile
	TrafficServerBinDir    string
	TrafficServerConfigDir string

	Cfg    *util.AtomicPtr[config.Cfg]
	TOData *util.AtomicPtr[TOData]

	TrafficMonitorHealth *util.AtomicPtr[TrafficMonitorHealth]
	ParentHealthL4       *util.AtomicPtr[ParentHealth]
	ParentHealthL7       *util.AtomicPtr[ParentHealth]
	ParentServiceHealth  *util.AtomicPtr[ParentServiceHealth]
	ParentHealthLog      io.WriteCloser

	MarkdownMethods map[config.HealthMethod]struct{}
	HealthMethods   map[config.HealthMethod]struct{}

	parents util.SyncMap[string, ParentStatus]

	// ParentHostFQDNs maps hostnames to fqdns, of parents.
	// This is necessary because Traffic Monitor's CRStates API only maps hostnames to health,
	// which is mostly okay for CDN caches; but we can't key on hostname anywhere else, because hostnames for non-cache parents aren't unique or meaningful.
	// So, we key on FQDN everywhere, and this mapping is used to map TM hostnames to their FQDNs.
	//
	// Don't use this for anything that includes non-cache (origin) servers.
	// Hostnames for any servers but caches are not unique and cannot be used as keys.
	ParentHostFQDNs util.SyncMap[string, string]
}

// TOData is the Traffic Ops data needed by various services.
type TOData struct {
	UserAgent string
	TOClient  *toclient.Session
	Monitors  map[string]struct{} `json:"trafficmonitors,omitempty"` // set[fqdn]
	Caches    map[string]struct{} `json:"caches,omitempty"`          // set[fqdn]
}

// Clone copies the TOData into a new object.
// Note the TOClient, which is safe for use by multiple goroutines, is not copied.
func (td *TOData) Clone() *TOData {
	newTD := NewTOData(td.UserAgent)
	for fqdn, _ := range td.Monitors {
		newTD.Monitors[fqdn] = struct{}{}
	}
	for fqdn, _ := range td.Caches {
		newTD.Caches[fqdn] = struct{}{}
	}
	newTD.TOClient = td.TOClient
	return newTD
}

func NewTOData(userAgent string) *TOData {
	return &TOData{
		UserAgent: userAgent,
		Monitors:  map[string]struct{}{},
		Caches:    map[string]struct{}{},
	}
}

// CombineParentHealth takes a ParentHealth directly queried from the parent service (ATS or the Origin)
// and the ParentServiceHealth directly queried from the parent's HealthService,
// combines them,
// and returns a ParentServiceHealth with each direct parent's RecursiveParentHealth containing
// both the direct ParentHealth and the ParentServiceHealth.
//
// If any parent exists in the direct ParentHealth but not ParentServiceHealth (which may be common,
// if cache parents aren't running the Health Service; and will always be true for origin parents),
// the direct parent will be inserted into the ParentServiceHealth.
//
// If any parent exists in the ParentServiceHealth but not ParentHealth, that means the actual parent
// service (ATS or the Origin) was unreachable, even though the Parent Health Service was reachable.
// In which case, we don't want to ever consider that parent healthy, so it will not be included.
func CombineParentHealth(parentHealthL4 *ParentHealth, parentHealthL7 *ParentHealth, parentServiceHealth *ParentServiceHealth) *ParentServiceHealth {
	// have to copy, we aren't allowed to modify the ParentHealth or ParentServiceHealth
	// (because there are many concurrent readers).

	combined := NewParentServiceHealth()
	for parentFQDN, parentDirectHealthL4 := range parentHealthL4.ParentHealthPollResults {
		parentCombined := RecursiveParentHealth{}

		parentHealthL4Copy := parentDirectHealthL4
		parentCombined.ParentHealthL4 = &parentHealthL4Copy

		// generally both L4 and L7 health should always exist,
		// but this might occur if the service just started, or just got a new parent,
		// and it hasn't had time to poll both yet
		if parentHealthL7, parentHealthL7Exists := parentHealthL7.ParentHealthPollResults[parentFQDN]; parentHealthL7Exists {
			parentHealthL7Copy := parentHealthL7
			parentCombined.ParentHealthL7 = &parentHealthL7Copy
		}

		if parentServiceHealth, ok := parentServiceHealth.ParentServiceHealthPollResults[parentFQDN]; ok {
			if parentServiceHealth.ParentServiceHealth == nil {
				log.Errorln("combining parent health: parent '" + parentFQDN + "' was in parent service health with a null recursive.ParentServiceHealth! Should never happen, the health should only ever contain polled service objects")
				continue
			}
			parentServiceHealthCopy := *parentServiceHealth.ParentServiceHealth
			parentCombined.ParentServiceHealth = &parentServiceHealthCopy
		}
		// TODO warn/err if parent is in parentHealth but not parentServiceHealth? And the reverse?
		combined.ParentServiceHealthPollResults[parentFQDN] = parentCombined
	}
	return combined
}

// LoadParentStatus returns parent's status and whether the parent existed.
// It is safe for multiple goroutines.
func (pi *ParentInfo) LoadParentStatus(fqdn string) (ParentStatus, bool) {
	return pi.parents.Load(fqdn)
}

// StoreParentStatus Sets a parent's status.
// It is safe for multiple goroutines.
func (pi *ParentInfo) StoreParentStatus(fqdn string, st ParentStatus) {
	pi.parents.Store(fqdn, st)
}

// GetParents returns a list of parent FQDNs.
// The map of parents is not locked during execution, so if the map data is changed concurrently,
// whether or not changes are returned will be nondeterministic.
// It is safe for multiple goroutines.
func (pi *ParentInfo) GetParents() []string {
	parentFQDNs := []string{}
	pi.parents.Range(func(parentFQDN string, parent ParentStatus) bool {
		parentFQDNs = append(parentFQDNs, parentFQDN)
		return true
	})
	return parentFQDNs
}

// GetCacheParents returns a list of parent FQDNs which are caches (as opposed to origins).
// See GetParents for concurrency and other details.
func (pi *ParentInfo) GetCacheParents() []string {
	// TODO make this smarter.
	parentFQDNs := []string{}
	toData := pi.TOData.Get()
	caches := toData.Caches
	numParents := 0
	numCaches := 0
	pi.parents.Range(func(parentFQDN string, parent ParentStatus) bool {
		numParents++
		if _, ok := caches[parentFQDN]; !ok {
			// log.Debugf("GetCacheParents parent '%v' not in caches, skipping\n", parent.Fqdn)
			return true // skip parents that aren't caches (i.e. origins)
		}
		numCaches++
		parentFQDNs = append(parentFQDNs, parentFQDN)
		return true
	})
	return parentFQDNs
}

// LoadOrStoreParentStatus adds the given status if fqdn does not exist,
// and returns the existing value if it does exist.
// Returns true if the value was loaded, false if it was stored.
//
// This behaves identical to sync.Map.LoadOrStore.
//
// It is safe for multiple goroutines.
func (pi *ParentInfo) LoadOrStoreParentStatus(fdqn string, status ParentStatus) (ParentStatus, bool) {
	return pi.parents.LoadOrStore(fdqn, status)
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
func NewParentInfo(cfgPtr *util.AtomicPtr[config.Cfg]) (*ParentInfo, error) {
	cfg := cfgPtr.Get()

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

	// read the 'parent.config'.
	if err := parentInfo.readParentConfig(); err != nil {
		return nil, errors.New("loading " + ParentsFile + " file: " + err.Error())
	}

	// read the strategies.yaml.
	if err := parentInfo.readStrategies(cfg.MonitorStrategiesPeers); err != nil {
		return nil, errors.New("loading parent " + StrategiesFile + " file: " + err.Error())
	}

	// collect the trafficserver parent status from the HostStatus subsystem.
	if err := parentInfo.readHostStatus(cfg); err != nil {
		return nil, fmt.Errorf("reading trafficserver host status: %w", err)
	}

	// TODO remove old parents no longer in parent.config, strategies.yaml, or traffic_ctl host status

	// log.Infof("startup loaded %d parent records\n", len(parentStatus))
	// TODO track how many elements are in the map?
	log.Infof("startup loaded %v parent records\n", len(parentInfo.GetParents()))

	parentInfo.Cfg = cfgPtr
	parentInfo.ParentHealthL4 = NewParentHealthPtr()
	parentInfo.ParentHealthL7 = NewParentHealthPtr()
	parentInfo.TrafficMonitorHealth = util.NewAtomicPtr(NewTrafficMonitorHealth())
	parentInfo.ParentServiceHealth = NewParentServiceHealthPtr()
	parentInfo.TOData = util.NewAtomicPtr(NewTOData(cfg.UserAgent))

	parentInfo.HealthMethods = map[config.HealthMethod]struct{}{}
	for _, hm := range *cfg.HealthMethods {
		parentInfo.HealthMethods[hm] = struct{}{}
	}

	parentInfo.MarkdownMethods = map[config.HealthMethod]struct{}{}
	for _, hm := range *cfg.MarkdownMethods {
		parentInfo.MarkdownMethods[hm] = struct{}{}
	}

	return &parentInfo, nil
}

// Used by the polling function to update the parents list from
// changes to 'parent.config' and 'strategies.yaml'.  The parents
// availability is also updated to reflect the current state from
// the trafficserver HostStatus subsystem.
func (pi *ParentInfo) UpdateParentInfo(cfg *config.Cfg) error {
	ptime, err := util.GetFileModificationTime(pi.ParentDotConfig.Filename)
	if err != nil {
		return errors.New("error reading " + ParentsFile + ": " + err.Error())
	}
	stime, err := util.GetFileModificationTime(pi.StrategiesDotYaml.Filename)
	if err != nil {
		return errors.New("error reading " + StrategiesFile + ": " + err.Error())
	}
	if pi.ParentDotConfig.LastModifyTime < ptime {
		// read the 'parent.config'.
		if err := pi.readParentConfig(); err != nil {
			return errors.New("updating " + ParentsFile + " file: " + err.Error())
		} else {
			// log.Infof("updated parents from new %s, total parents: %d\n", ParentsFile, len(pi.Parents))
			// TODO track map len
			log.Infof("tm-agent total_parents=%v event=\"updated parents from new parent.config\"\n", len(pi.GetParents()))
		}
	}

	if pi.StrategiesDotYaml.LastModifyTime < stime {
		// read the 'strategies.yaml'.
		if err := pi.readStrategies(cfg.MonitorStrategiesPeers); err != nil {
			return errors.New("updating parent " + StrategiesFile + " file: " + err.Error())
		} else {
			// log.Infof("updated parents from new %s total parents: %d\n", StrategiesFile, len(pi.Parents))
			// TODO track map len
			log.Infof("tm-agent total_parents=%v event=\"updated parents from new strategies.yaml\"\n", len(pi.GetParents()))
		}
	}

	// collect the trafficserver current host status.
	if err := pi.readHostStatus(cfg); err != nil {
		return errors.New("trafficserver may not be running: " + err.Error())
	}

	return nil
}

func (pi *ParentInfo) WritePollState() error {
	cfg := pi.Cfg.Get()
	data, err := json.MarshalIndent(pi, "", "\t")
	if err != nil {
		return fmt.Errorf("marshaling configuration state: %s\n", err.Error())
	} else {
		err = os.WriteFile(cfg.PollStateJSONLog, data, 0644)
		if err != nil {
			return fmt.Errorf("writing configuration state: %s\n", err.Error())
		}
	}
	return nil
}

// findATrafficMonitor chooses an available trafficmonitor,
// and returns an error if there are none.
func (pi *ParentInfo) findATrafficMonitor() (string, error) {
	toData := pi.TOData.Get()

	var tmHostname string
	lth := len(toData.Monitors)
	if lth == 0 {
		return "", errors.New("there are no available traffic monitors")
	}

	// build an array of available traffic monitors.
	tms := []string{}
	for fqdn, _ := range toData.Monitors {
		tms = append(tms, fqdn)
	}

	// choose one at random.
	lth = len(tms)
	if lth > 0 {
		// TODO make deterministic. Hash of hostname?
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

// readParentConfig loads the parents list from the Trafficserver 'parent.config' file.
func (pi *ParentInfo) readParentConfig() error {
	fn := pi.ParentDotConfig.Filename

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
	pi.ParentDotConfig.LastModifyTime = finfo.ModTime().UnixNano()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		sbytes := scanner.Bytes()
		sbytes = bytes.TrimSpace(sbytes)
		if len(sbytes) == 0 {
			continue // skip blank lines
		}
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

						{
							parentHostName := parseFqdn(fqdn)
							pi.ParentHostFQDNs.Store(parentHostName, fqdn)
						}

						// create the ParentStatus struct and add it to the
						// Parents map only if an entry in the map does not
						// already exist.
						_, loaded := pi.LoadOrStoreParentStatus(fqdn, ParentStatus{
							Fqdn:                 strings.TrimSpace(fqdn),
							ActiveReason:         true,
							LocalReason:          true,
							ManualReason:         true,
							LastTmPoll:           0,
							UnavailablePollCount: 0,
						})
						if !loaded {
							log.Debugf("added Host '%s' from %s to the parents map\n", fqdn, fn)
						}
					}
				}
			}
		}
	}
	return nil
}

// load the parent hosts from 'strategies.yaml'.
func (pi *ParentInfo) readStrategies(monitorPeers bool) error {
	var includes []string
	fn := pi.StrategiesDotYaml.Filename

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
	pi.StrategiesDotYaml.LastModifyTime = finfo.ModTime().UnixNano()

	scanner := bufio.NewScanner(f)

	// search for any yaml files that should be included in the
	// yaml stream.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#include") {
			fields := strings.Split(line, " ")
			if len(fields) >= 2 {
				includeFile := filepath.Join(pi.TrafficServerConfigDir, fields[1])
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

	// If we are to not monitor peers, this will set the hosts to non-peer hosts only
	if !monitorPeers {
		type YAMLHosts struct {
			Hosts []yaml.Node `yaml:"hosts"`
		}

		// Empty the hosts since we want to rebuild
		strategies.Hosts = []Host{}

		var yamlPeers YAMLHosts

		if err := yaml.Unmarshal([]byte(yamlContent), &yamlPeers); err != nil {
			return errors.New("failed to unmarshall " + fn + ": " + err.Error())
		}

		for _, host := range yamlPeers.Hosts {
			if host.Anchor[0:4] != "peer" {
				var hostObj Host
				if err := host.Decode(&hostObj); err != nil {
					return errors.New("Failed to unmarshall non-peer object. " + fn + ": " + err.Error())
				}

				strategies.Hosts = append(strategies.Hosts, hostObj)
			}
		}
	}

	for _, host := range strategies.Hosts {
		fqdn := host.HostName

		{
			parentHostName := parseFqdn(fqdn)
			pi.ParentHostFQDNs.Store(parentHostName, fqdn)
		}

		_, loaded := pi.LoadOrStoreParentStatus(fqdn, ParentStatus{
			Fqdn:                 strings.TrimSpace(fqdn),
			ActiveReason:         true,
			LocalReason:          true,
			ManualReason:         true,
			LastTmPoll:           0,
			UnavailablePollCount: 0,
		})
		if !loaded {
			log.Debugf("added Host '%s' from %s to the parents map\n", fqdn, fn)
		}
	}
	return nil
}

// GetTOData fetches data needed from Traffic Ops and refreshes it in the ParentInfo object.
//
// Note it takes a Config, even though ParentInfo has a AtomicPtr[Config], because callers typically
// need to atomically operate on a single Config.
// If a caller doesn't already have a Config, simply call pi.GetTOData(pi.Cfg.Get()).
func (pi *ParentInfo) GetTOData(cfg *config.Cfg) error {
	// TODO can we use the t3c cache here?
	toData := pi.TOData.Get().Clone()

	if toData.TOClient == nil {
		session, _, err := toclient.LoginWithAgent(cfg.TOUrl, cfg.TOUser, cfg.TOPass, true, toData.UserAgent, false, cfg.TORequestTimeout)
		if err != nil {
			return fmt.Errorf("could not establish a TrafficOps session: %w", err)
		} else {
			toData.TOClient = session
		}
	}

	srvs, reqInf, err := toData.TOClient.GetServers(toclient.NewRequestOptions())
	if err != nil {
		// next time we'll login again and get a new session.
		toData.TOClient = nil
		pi.TOData.Set(toData)
		return errors.New("error fetching Trafficmonitor server list: " + err.Error())
	} else if reqInf.StatusCode >= 300 || reqInf.StatusCode < 200 {
		// Provide logging around a potential issue
		return fmt.Errorf("Traffic Ops returned a non 2xx status code. Expected 2xx, got %v", reqInf.StatusCode)
	}

	toData.Monitors = map[string]struct{}{}
	toData.Caches = map[string]struct{}{}
	for _, sv := range srvs.Response {
		log.Debugf("GetTOData server '%v' type '%v'\n", *sv.HostName, sv.Type)
		if sv.HostName == nil {
			log.Errorf("Traffic Ops returned server with nil hostname, skipping!")
		} else if sv.CDNName == nil || sv.DomainName == nil || sv.Status == nil {
			log.Errorf("Traffic Ops returned server '" + *sv.HostName + "' with nil required fields, skipping!")
		}
		if *sv.CDNName != cfg.CDNName {
			continue
		}
		if sv.Type == tc.MonitorTypeName && tc.CacheStatus(*sv.Status) == tc.CacheStatusOnline {
			fqdn := *sv.HostName + "." + *sv.DomainName
			toData.Monitors[fqdn] = struct{}{}
			continue
		}
		if tc.CacheType(sv.Type) == tc.CacheTypeEdge || tc.CacheType(sv.Type) == tc.CacheTypeMid {
			fqdn := *sv.HostName + "." + *sv.DomainName
			toData.Caches[fqdn] = struct{}{}
			continue
		}
	}

	pi.TOData.Set(toData)

	return nil
}

// readHostStatus reads the current parent statuses from the trafficserver HostStatus subsystem.
func (pi *ParentInfo) readHostStatus(cfg *config.Cfg) error {
	tc := filepath.Join(pi.TrafficServerBinDir, TrafficCtl)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// auto select traffic_ctl command for ATS version 9 or 10 and later
	for i := traffic_ctl_index; i <= 1; i++ {
		var err error
		switch i {
		case 0: // ATS version 10 and later
			cmd := exec.Command(tc, "host", "status")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
		case 1: // ATS version 9
			cmd := exec.Command(tc, "metric", "match", "host_status")
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
		}
		if err == nil {
			break
		}
		if err != nil && i == 0 {
			log.Infof("%s command used is not for ATS version 10 or later, downgrading to ATS version 9\n", TrafficCtl)
			traffic_ctl_index = 1
			continue
		}
		if err != nil {
			return fmt.Errorf("%s error: %s", TrafficCtl, stderr.String())
		}
	}

	if len((stdout.Bytes())) > 0 {
		var activeReason bool
		var localReason bool
		var manualReason bool
		var fqdn string
		scanner := bufio.NewScanner(bytes.NewReader(stdout.Bytes()))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			fields := strings.Split(line, " ")
			/*
			 * For ATS Version 9, the host status uses internal stats and prefixes
			 * the fqdn field from the output of the traffic_ctl host status and metric
			 * match commands with "proxy.process.host_status".  Going forward starting
			 * with ATS Version 10, internal stats are no-longer used and the fqdn field
			 * is no-longer prefixed with the "proxy.process.host_status" string.
			 */
			if len(fields) == 2 {
				// check for ATS version 9 output.
				fqdnField := strings.Split(fields[0], "proxy.process.host_status.")
				if len(fqdnField) == 2 { // ATS version 9
					fqdn = fqdnField[1]
				} else { // ATS version 10 and greater
					fqdn = fqdnField[0]
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
					MarkUpPollCount:      0,
				}
				log.Debugf("processed host status record: %v\n", pstat)

				{
					parentHostName := parseFqdn(fqdn)
					pi.ParentHostFQDNs.Store(parentHostName, fqdn)
				}

				// create the ParentStatus struct and add it to the
				// Parents map only if an entry in the map does not
				// already exist.
				pv, loaded := pi.LoadOrStoreParentStatus(fqdn, pstat)
				if !loaded {
					log.Infof("added Host '%s' from ATS Host Status to the parents map\n", fqdn)
				} else {
					available := pstat.available(cfg.ReasonCode)
					if pv.available(cfg.ReasonCode) != available {
						log.Infof("host status for '%s' has changed to %s\n", fqdn, pstat.Status())
						pstat.LastTmPoll = pv.LastTmPoll
						pstat.UnavailablePollCount = pv.UnavailablePollCount
						pstat.MarkUpPollCount = pv.MarkUpPollCount
						pi.StoreParentStatus(fqdn, pstat)
					}
				}
			}
		}
		// log.Debugf("processed trafficserver host status results, total parents: %d\n", len(parentStatus))
		// TODO count parentStatus len?
		log.Debugf("tm-agent total_parents=%v\n", len(pi.GetParents()))
	}
	return nil
}
