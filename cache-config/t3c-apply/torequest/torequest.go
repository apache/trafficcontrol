package torequest

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
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-apply/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3c-apply/util"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

type UpdateStatus int

const (
	UpdateTropsNotNeeded  UpdateStatus = 0
	UpdateTropsNeeded     UpdateStatus = 1
	UpdateTropsSuccessful UpdateStatus = 2
	UpdateTropsFailed     UpdateStatus = 3
)

const (
	TailDiagsLogRelative = "/var/log/trafficserver/diags.log"
	TailRestartTimeOutMS = 60000
	TailReloadTimeOutMS  = 15000
	tailMatch            = `ET_(TASK|NET)\s\d{1,}`
	tailRestartEnd       = "Traffic Server is fully initialized"
	tailReloadEnd        = "remap.config finished loading"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type TrafficOpsReq struct {
	Cfg     config.Cfg
	Pkgs    map[string]bool // map of packages which are installed, either already installed or newly installed by this run.
	plugins map[string]bool // map of verified plugins

	installedPkgs map[string]struct{} // map of packages which were installed by us.
	changedFiles  []string            // list of config files which were changed

	configFiles        map[string]*ConfigFile
	configFileWarnings map[string][]string

	RestartData
}

type ShouldReloadRestart struct {
	ReloadRestart []FileRestartData
}

type FileRestartData struct {
	Name string
	RestartData
}

type RestartData struct {
	TrafficCtlReload     bool // a traffic_ctl_reload is required
	SysCtlReload         bool // a reload of the sysctl.conf is required
	NtpdRestart          bool // ntpd needs restarting
	TeakdRestart         bool // a restart of teakd is required
	TrafficServerRestart bool // a trafficserver restart is required
	RemapConfigReload    bool // remap.config should be reloaded
	HitchReload          bool // hitch should be reloaded
	VarnishReload        bool // varnish should be reloaded
}

type ConfigFile struct {
	Name              string // file name
	Dir               string // install directory
	Path              string // full path
	Service           string // service assigned to
	CfgBackup         string // location to backup the config at 'Path'
	TropsBackup       string // location to backup the TrafficOps Version
	AuditComplete     bool   // audit is complete
	AuditFailed       bool   // audit failed
	AuditError        string // Error generated when AuditFailed is true
	ChangeApplied     bool   // a change has been applied
	ChangeNeeded      bool   // change required
	PreReqFailed      bool   // failed plugin prerequiste check
	RemapPluginConfig bool   // file is a remap plugin config file
	Body              []byte
	Perm              os.FileMode // default file permissions
	Uid               int         // owner uid, default is 0
	Gid               int         // owner gid, default is 0
	Warnings          []string
}

func (u UpdateStatus) String() string {
	var result string
	switch u {
	case 0:
		result = "UpdateTropsNotNeeded"
	case 1:
		result = "UpdateTropsNeeded"
	case 2:
		result = "UpdateTropsSuccessful"
	case 3:
		result = "UpdateTropsFailed"
	}
	return result
}

// commentsFilter is used to remove comment
// lines from config files while making
// comparisons.
func commentsFilter(body []string) []string {
	var newlines []string

	newlines = make([]string, 0)

	for ii := range body {
		line := body[ii]
		if strings.HasPrefix(line, "#") {
			continue
		}
		newlines = append(newlines, line)
	}

	return newlines
}

// newLineFilter removes carriage returns
// from config files while making comparisons.
func newLineFilter(str string) string {
	str = strings.ReplaceAll(str, "\r\n", "\n")
	return strings.TrimSpace(str)
}

// unencodeFilter translates HTML escape
// sequences while making config file comparisons.
func unencodeFilter(body []string) []string {
	var newlines []string

	newlines = make([]string, 0)
	sp := regexp.MustCompile(`\s+`)
	el := regexp.MustCompile(`^\s+|\s+$`)
	am := regexp.MustCompile(`amp;`)
	lt := regexp.MustCompile(`&lt;`)
	gt := regexp.MustCompile(`&gt;`)

	for ii := range body {
		s := body[ii]
		s = sp.ReplaceAllString(s, " ")
		s = el.ReplaceAllString(s, "")
		s = am.ReplaceAllString(s, "")
		s = lt.ReplaceAllString(s, "<")
		s = gt.ReplaceAllString(s, ">")
		s = strings.TrimSpace(s)
		newlines = append(newlines, s)
	}

	return newlines
}

// DumpConfigFiles is used for debugging
func (r *TrafficOpsReq) DumpConfigFiles() {
	for _, cfg := range r.configFiles {
		log.Infof("Name: %s, Dir: %s, Service: %s\n",
			cfg.Name, cfg.Dir, cfg.Service)
	}
}

// NewTrafficOpsReq returns a new TrafficOpsReq object.
func NewTrafficOpsReq(cfg config.Cfg) *TrafficOpsReq {
	return &TrafficOpsReq{
		Cfg:           cfg,
		Pkgs:          map[string]bool{},
		plugins:       map[string]bool{},
		configFiles:   map[string]*ConfigFile{},
		installedPkgs: map[string]struct{}{},
	}
}

// checkConfigFile checks and audits config files.
// The filesAdding parameter is the list of files about to be added, which is needed for verification in case a file is required and about to be created but doesn't exist yet.
func (r *TrafficOpsReq) checkConfigFile(cfg *ConfigFile, filesAdding []string) error {
	if cfg.Name == "" {
		cfg.AuditFailed = true
		return errors.New("Config file name is empty is empty, skipping further checks.")
	}

	if cfg.Dir == "" {
		cfg.AuditFailed = true
		return errors.New("No location information for " + cfg.Name)
	}
	// return if audit has already been done.
	if cfg.AuditComplete {
		return nil
	}

	if !util.MkDirWithOwner(cfg.Dir, r.Cfg.ReportOnly, &cfg.Uid, &cfg.Gid) {
		cfg.AuditFailed = true
		return errors.New("Unable to create the directory '" + cfg.Dir + " for " + "'" + cfg.Name + "'")
	}

	log.Debugf("======== Start processing config file: %s ========\n", cfg.Name)

	if cfg.Name == "50-ats.rules" {
		err := r.processUdevRules(cfg)
		if err != nil {
			cfg.AuditFailed = true
			return errors.New("unable to process udev rules in '" + cfg.Name + "': " + err.Error())
		}
	}

	if cfg.Name == "remap.config" {
		err := r.processRemapOverrides(cfg)
		if err != nil {
			cfg.AuditFailed = true
			return err
		}
	}

	// perform plugin verification
	if cfg.Name == "remap.config" || cfg.Name == "plugin.config" {
		if err := checkRefs(r.Cfg, cfg.Body, filesAdding); err != nil {
			r.configFileWarnings[cfg.Name] = append(r.configFileWarnings[cfg.Name], "failed to verify '"+cfg.Name+"': "+err.Error())
			cfg.AuditFailed = true
			return errors.New("failed to verify '" + cfg.Name + "': " + err.Error())
		}
		log.Infoln("Successfully verified plugins used by '" + cfg.Name + "'")
	}

	if strings.HasSuffix(cfg.Name, ".cer") {
		err, fatal := checkCert(cfg.Body)
		if err != nil {
			r.configFileWarnings[cfg.Name] = append(r.configFileWarnings[cfg.Name], err.Error())
		}
		r.configFileWarnings[cfg.Name] = append(r.configFileWarnings[cfg.Name], cfg.Warnings...)
		if fatal {
			return errors.New(err.Error() + " for: " + cfg.Name)
		}
	}

	changeNeeded, err := diff(r.Cfg, cfg.Body, cfg.Path, r.Cfg.ReportOnly, cfg.Perm, cfg.Uid, cfg.Gid)

	if err != nil {
		cfg.AuditFailed = true
		return errors.New("getting diff: " + err.Error())
	}
	cfg.ChangeNeeded = changeNeeded
	cfg.AuditComplete = true

	log.Infof("======== End processing config file: %s for service: %s ========\n", cfg.Name, cfg.Service)
	return nil
}

// checkStatusFiles ensures that the cache status files reflect
// the status retrieved from Traffic Ops.
func (r *TrafficOpsReq) checkStatusFiles(svrStatus string) error {
	if svrStatus == "" {
		return errors.New("Returning; did not find status from Traffic Ops!")
	} else {
		log.Debugf("Found %s status from Traffic Ops.\n", svrStatus)
	}
	statusFile := filepath.Join(config.StatusDir, svrStatus)
	fileExists, _ := util.FileExists(statusFile)
	if !fileExists {
		log.Errorf("status file %s does not exist.\n", statusFile)
	}
	statuses, err := getStatuses(r.Cfg)
	if err != nil {
		return fmt.Errorf("could not retrieves a statuses list from Traffic Ops: %s\n", err)
	}

	for f := range statuses {
		otherStatus := filepath.Join(config.StatusDir, statuses[f])
		if otherStatus == statusFile {
			continue
		}
		fileExists, _ := util.FileExists(otherStatus)
		if !r.Cfg.ReportOnly && fileExists {
			log.Errorf("Removing other status file %s that exists\n", otherStatus)
			err = os.Remove(otherStatus)
			if err != nil {
				log.Errorf("Error removing %s: %s\n", otherStatus, err)
			}
		}
	}

	if !r.Cfg.ReportOnly {
		if !util.MkDir(config.StatusDir, r.Cfg.ReportOnly) {
			return fmt.Errorf("unable to create '%s'\n", config.StatusDir)
		}
		fileExists, _ := util.FileExists(statusFile)
		if !fileExists {
			err = util.Touch(statusFile)
			if err != nil {
				return fmt.Errorf("unable to touch %s - %s\n", statusFile, err)
			}
		}
	}
	return nil
}

// processRemapOverrides processes remap overrides found from Traffic Ops.
func (r *TrafficOpsReq) processRemapOverrides(cfg *ConfigFile) error {
	from := ""
	newlines := []string{}
	lineCount := 0
	overrideCount := 0
	overridenCount := 0
	overrides := map[string]int{}
	data := cfg.Body

	if len(data) > 0 {
		lines := strings.Split(string(data), "\n")
		for ii := range lines {
			str := lines[ii]
			fields := strings.Fields(str)
			if str == "" || len(fields) < 2 {
				continue
			}
			lineCount++
			from = fields[1]

			_, ok := overrides[from]
			if ok == true { // check if this line should be overriden
				newstr := "##OVERRIDDEN## " + str
				newlines = append(newlines, newstr)
				overridenCount++
			} else if fields[0] == "##OVERRIDE##" { // check for an override
				from = fields[2]
				newlines = append(newlines, "##OVERRIDE##")
				// remove the ##OVERRIDE## comment along with the trailing space
				newstr := strings.TrimPrefix(str, "##OVERRIDE## ")
				// save the remap 'from field' to overrides.
				overrides[from] = 1
				newlines = append(newlines, newstr)
				overrideCount++
			} else { // no override is necessary
				newlines = append(newlines, str)
			}
		}
	} else {
		return errors.New("The " + cfg.Name + " file is empty, nothing to process.")
	}
	if overrideCount > 0 {
		log.Infof("Overrode %d old remap rule(s) with %d new remap rule(s).\n",
			overridenCount, overrideCount)
		newdata := strings.Join(newlines, "\n")
		// strings.Join doesn't add a newline character to
		// the last element in the array and we need one
		// when the data is written out to a file.
		if !strings.HasSuffix(newdata, "\n") {
			newdata = newdata + "\n"
		}
		body := []byte(newdata)
		cfg.Body = body
	}
	return nil
}

// processUdevRules verifies disk drive device ownership and mode
func (r *TrafficOpsReq) processUdevRules(cfg *ConfigFile) error {
	var udevDevices map[string]string

	data := string(cfg.Body)
	lines := strings.Split(data, "\n")

	udevDevices = make(map[string]string)
	for ii := range lines {
		var owner string
		var device string
		line := lines[ii]
		if strings.HasPrefix(line, "KERNEL==") {
			vals := strings.Split(line, "\"")
			if len(vals) >= 3 {
				device = vals[1]
				owner = vals[3]
				if owner == "root" {
					continue
				}
				userInfo, err := user.Lookup(owner)
				if err != nil {
					log.Errorf("no such user on this system: '%s'\n", owner)
					continue
				} else {
					devPath := "/dev/" + device
					fileExists, fileInfo := util.FileExists(devPath)
					if fileExists {
						udevDevices[device] = devPath
						log.Infof("Found device in 50-ats.rules: %s\n", devPath)
						if statStruct, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
							uid := strconv.Itoa(int(statStruct.Uid))
							if uid != userInfo.Uid {
								log.Errorf("Device %s is owned by uid %s, not %s (%s)\n", devPath, uid, owner, userInfo.Uid)
							} else {
								log.Infof("Ownership for disk device %s, is okay\n", devPath)
							}
						} else {
							log.Errorf("Unable to read device owner info for %s\n", devPath)
						}
					}
				}
			}
		}
	}
	fs, err := os.ReadDir("/proc/fs/ext4")
	if err != nil {
		log.Errorln("unable to read /proc/fs/ext4, cannot audit disks for filesystem usage.")
	} else {
		for _, disk := range fs {
			for k, _ := range udevDevices {
				if strings.HasPrefix(k, disk.Name()) {
					log.Warnf("Device %s has an active partition and filesystem!!!!\n", k)
				}
			}
		}
	}

	return nil
}

// readCfgFile reads a config file and return its contents.
func (r *TrafficOpsReq) readCfgFile(cfg *ConfigFile, dir string) ([]byte, error) {
	var data []byte
	var fullFileName string
	if dir == "" {
		fullFileName = cfg.Path
	} else {
		fullFileName = dir + "/" + cfg.Name
	}

	info, err := os.Stat(fullFileName)
	if err != nil {
		return nil, err
	}
	size := info.Size()

	fd, err := os.Open(fullFileName)
	if err != nil {
		return nil, err
	}
	data = make([]byte, size)
	c, err := fd.Read(data)
	if err != nil || int64(c) != size {
		return nil, errors.New("unable to completely read from '" + cfg.Name + "': " + err.Error())
	}
	fd.Close()

	return data, nil
}

const configFileTempSuffix = `.tmp`

// replaceCfgFile replaces an ATS configuration file with one from Traffic Ops.
func (r *TrafficOpsReq) replaceCfgFile(cfg *ConfigFile) (*FileRestartData, error) {
	if r.Cfg.ReportOnly ||
		(r.Cfg.Files != t3cutil.ApplyFilesFlagAll && r.Cfg.Files != t3cutil.ApplyFilesFlagReval) {
		log.Infof("You elected not to replace %s with the version from Traffic Ops.\n", cfg.Name)
		cfg.ChangeApplied = false
		return &FileRestartData{Name: cfg.Name}, nil
	}

	tmpFileName := cfg.Path + configFileTempSuffix
	log.Infof("Writing temp file '%s' with file mode: '%#o' \n", tmpFileName, cfg.Perm)

	// write a new file, then move to the real location
	// because moving is atomic but writing is not.
	// If we just wrote to the real location and the app or OS or anything crashed,
	// we'd end up with malformed files.

	if _, err := util.WriteFileWithOwner(tmpFileName, cfg.Body, &cfg.Uid, &cfg.Gid, cfg.Perm); err != nil {
		return &FileRestartData{Name: cfg.Name}, errors.New("Failed to write temp config file '" + tmpFileName + "': " + err.Error())
	}

	log.Infof("Copying temp file '%s' to real '%s'\n", tmpFileName, cfg.Path)
	if err := os.Rename(tmpFileName, cfg.Path); err != nil {
		return &FileRestartData{Name: cfg.Name}, errors.New("Failed to move temp '" + tmpFileName + "' to real '" + cfg.Path + "': " + err.Error())
	}
	cfg.ChangeApplied = true
	r.changedFiles = append(r.changedFiles, cfg.Path)

	remapConfigReload := cfg.RemapPluginConfig ||
		cfg.Name == "remap.config" ||
		strings.HasPrefix(cfg.Name, "bg_fetch") ||
		strings.HasPrefix(cfg.Name, "hdr_rw_") ||
		strings.HasPrefix(cfg.Name, "regex_remap_") ||
		strings.HasPrefix(cfg.Name, "set_dscp_") ||
		strings.HasPrefix(cfg.Name, "url_sig_") ||
		strings.HasPrefix(cfg.Name, "uri_signing") ||
		strings.HasSuffix(cfg.Name, ".lua")

	trafficCtlReload := strings.HasSuffix(cfg.Dir, "trafficserver") ||
		remapConfigReload ||
		cfg.Name == "ssl_multicert.config" ||
		cfg.Name == "records.config" ||
		(strings.HasSuffix(cfg.Dir, "ssl") && strings.HasSuffix(cfg.Name, ".cer")) ||
		(strings.HasSuffix(cfg.Dir, "ssl") && strings.HasSuffix(cfg.Name, ".key"))

	trafficServerRestart := cfg.Name == "plugin.config"
	ntpdRestart := cfg.Name == "ntpd.conf"
	sysCtlReload := cfg.Name == "sysctl.conf"
	hitchReload := cfg.Name == "hitch.conf"
	varnishReload := cfg.Name == "default.vcl"

	log.Debugf("Reload state after %s: remap.config: %t reload: %t restart: %t ntpd: %t sysctl: %t", cfg.Name, remapConfigReload, trafficCtlReload, trafficServerRestart, ntpdRestart, sysCtlReload)

	log.Debugf("Setting change applied for '%s'\n", cfg.Name)
	return &FileRestartData{
		Name: cfg.Name,
		RestartData: RestartData{
			TrafficCtlReload:     trafficCtlReload,
			SysCtlReload:         sysCtlReload,
			NtpdRestart:          ntpdRestart,
			TrafficServerRestart: trafficServerRestart,
			RemapConfigReload:    remapConfigReload,
			HitchReload:          hitchReload,
			VarnishReload:        varnishReload,
		},
	}, nil
}

// CheckSystemServices is used to verify that packages installed
// are enabled for startup.
func (r *TrafficOpsReq) CheckSystemServices() error {
	if r.Cfg.ServiceAction != t3cutil.ApplyServiceActionFlagRestart {
		return nil
	}
	result, err := getChkconfig(r.Cfg)
	if err != nil {
		log.Errorln(err)
		return err
	}
	for ii := range result {
		name := result[ii]["name"]
		value := result[ii]["value"]
		arrv := strings.Fields(value)
		level := []string{}
		enabled := false
		for jj := range arrv {
			nv := strings.Split(arrv[jj], ":")
			if len(nv) == 2 && strings.Contains(nv[1], "on") {
				level = append(level, nv[0])
				enabled = true
			}
		}
		if !enabled {
			continue
		}
		if r.Cfg.SvcManagement == config.SystemD {
			out, rc, err := util.ExecCommand("/bin/systemctl", "enable", name)
			if err != nil {
				log.Errorf(string(out))
				return errors.New("Unable to enable service " + name + ": " + err.Error())
			}
			if rc == 0 {
				log.Infof("The %s service has been enabled\n", name)
			}
		} else if r.Cfg.SvcManagement == config.SystemV {
			levelValue := strings.Join(level, "")
			_, rc, err := util.ExecCommand("/bin/chkconfig", "--level", levelValue, name, "on")
			if err != nil {
				return errors.New("Unable to enable service " + name + ": " + err.Error())
			}
			if rc == 0 {
				log.Infof("The %s service has been enabled\n", name)
			}
		} else {
			log.Errorf("Unable to ensure %s service is enabled, SvcMananagement type is %s\n", name, r.Cfg.SvcManagement)
		}
	}
	return nil
}

// IsPackageInstalled returns true/false if the named rpm package is installed.
// the prefix before the version is matched.
func (r *TrafficOpsReq) IsPackageInstalled(name string) bool {
	for k, v := range r.Pkgs {
		if strings.HasPrefix(k, name) {
			log.Infof("Found in cache for '%s'", k)
			return v
		}
	}
	if !r.Cfg.RpmDBOk {
		log.Warnf("RPM DB is corrupted cannot run IsPackageInstalled for '%s' and package metadata is unavailable", name)
		return false
	}
	log.Infof("IsPackageInstalled '%v' not found in cache, querying rpm", name)
	pkgArr, err := util.PackageInfo("pkg-query", name)
	if err != nil {
		log.Errorf(`IsPackageInstalled PackageInfo(pkg-query, %v) failed, caching as not installed and returning false! Error: %v\n`, name, err.Error())
		r.Pkgs[name] = false
		return false
	}
	if len(pkgArr) > 0 {
		pkgAndVersion := pkgArr[0]
		log.Infof("IsPackageInstalled '%v' found in rpm, adding '%v' to cache", name, pkgAndVersion)
		r.Pkgs[pkgAndVersion] = true
		return true
	}
	log.Infof("IsPackageInstalled '%v' not found in rpm, adding '%v'=false to cache", name, name)
	r.Pkgs[name] = false
	return false
}

// GetConfigFile fetchs a 'Configfile' by file name.
func (r *TrafficOpsReq) GetConfigFile(name string) (*ConfigFile, bool) {
	cfg, ok := r.configFiles[name]
	return cfg, ok
}

// GetConfigFileList fetches and parses the multipart config files
// for a cache from traffic ops and loads them into the configFiles map.
func (r *TrafficOpsReq) GetConfigFileList() error {
	var atsUid int = 0
	var atsGid int = 0

	atsUser, err := user.Lookup(config.TrafficServerOwner)
	if err != nil {
		log.Errorf("could not lookup the trafficserver, '%s', owner uid, using uid/gid 0",
			config.TrafficServerOwner)
	} else {
		atsUid, err = strconv.Atoi(atsUser.Uid)
		if err != nil {
			log.Errorf("could not parse the ats UID.")
			atsUid = 0
		}
		atsGid, err = strconv.Atoi(atsUser.Gid)
		if err != nil {
			log.Errorf("could not parse the ats GID.")
			atsUid = 0
		}
	}

	allFiles, err := generate(r.Cfg)
	if err != nil {
		return errors.New("requesting data generating config files: " + err.Error())
	}

	r.configFiles = map[string]*ConfigFile{}
	r.configFileWarnings = map[string][]string{}
	var mode os.FileMode
	for _, file := range allFiles {
		if file.Secure {
			mode = 0600
		} else {
			mode = 0644
		}

		r.configFiles[file.Name] = &ConfigFile{
			Name:     file.Name,
			Path:     filepath.Join(file.Path, file.Name),
			Dir:      file.Path,
			Body:     []byte(file.Text),
			Uid:      atsUid,
			Gid:      atsGid,
			Perm:     mode,
			Warnings: file.Warnings,
		}
		for _, warn := range file.Warnings {
			if warn == "" {
				continue
			}
			r.configFileWarnings[file.Name] = append(r.configFileWarnings[file.Name], warn)
		}
	}

	return nil
}

func (r *TrafficOpsReq) PrintWarnings() {
	log.Infoln("======== Summary of config warnings that may need attention. ========")
	for file, warning := range r.configFileWarnings {
		for _, warning := range warning {
			log.Warnf("%s: %s", file, warning)
		}
	}
	log.Infoln("======== End warning summary ========")
}

// CheckRevalidateState retrieves and returns the revalidate status from Traffic Ops.
func (r *TrafficOpsReq) CheckRevalidateState(sleepOverride bool) (UpdateStatus, error) {
	log.Infoln("Checking revalidate state.")
	if !sleepOverride &&
		(r.Cfg.ReportOnly || r.Cfg.Files != t3cutil.ApplyFilesFlagReval) {
		updateStatus := UpdateTropsNotNeeded
		log.Infof("CheckRevalidateState returning %v\n", updateStatus)
		return updateStatus, nil
	}

	updateStatus := UpdateTropsNotNeeded

	serverStatus, err := getUpdateStatus(r.Cfg)
	if err != nil {
		log.Errorln("getting update status: " + err.Error())
		return UpdateTropsNotNeeded, errors.New("getting update status: " + err.Error())
	}
	log.Infof("my status: %s\n", serverStatus.Status)
	if serverStatus.UseRevalPending == false {
		log.Errorln("Update URL: Instant invalidate is not enabled.  Separated revalidation requires upgrading to Traffic Ops version 2.2 and enabling this feature.")
		return UpdateTropsNotNeeded, nil
	}
	if serverStatus.RevalPending == true {
		log.Errorln("Traffic Ops is signaling that a revalidation is waiting to be applied.")
		updateStatus = UpdateTropsNeeded
		if serverStatus.ParentRevalPending == true {
			if r.Cfg.WaitForParents {
				log.Infoln("Traffic Ops is signaling that my parents need to revalidate, not revalidating.")
				updateStatus = UpdateTropsNotNeeded
			} else {
				log.Infoln("Traffic Ops is signaling that my parents need to revalidate, but wait-for-parents is false, revalidating anyway.")
			}
		}
	} else if serverStatus.RevalPending == false && !r.Cfg.ReportOnly && r.Cfg.Files == t3cutil.ApplyFilesFlagReval {
		log.Errorln("In revalidate mode, but no update needs to be applied. I'm outta here.")
		return UpdateTropsNotNeeded, nil
	} else {
		log.Errorln("Traffic Ops is signaling that no revalidations are waiting to be applied.")
		return UpdateTropsNotNeeded, nil
	}

	err = r.checkStatusFiles(serverStatus.Status)
	if err != nil {
		log.Errorln(errors.New("checking status files: " + err.Error()))
	} else {
		log.Infoln("CheckRevalidateState checkStatusFiles returned nil error")
	}

	log.Infof("CheckRevalidateState returning %v\n", updateStatus)
	return updateStatus, nil
}

// CheckSyncDSState retrieves and returns the DS Update status from Traffic Ops.
// The metaData is this run's metadata. It must not be nil, and this function may add to it.
func (r *TrafficOpsReq) CheckSyncDSState(metaData *t3cutil.ApplyMetaData, cfg config.Cfg) (UpdateStatus, error) {
	updateStatus := UpdateTropsNotNeeded
	randDispSec := time.Duration(0)
	log.Debugln("Checking syncds state.")
	//	if r.Cfg.RunMode == t3cutil.ModeSyncDS || r.Cfg.RunMode == t3cutil.ModeBadAss || r.Cfg.RunMode == t3cutil.ModeReport {
	if r.Cfg.Files != t3cutil.ApplyFilesFlagReval {
		serverStatus, err := getUpdateStatus(r.Cfg)
		if err != nil {
			log.Errorln("getting '" + r.Cfg.CacheHostName + "' update status: " + err.Error())
			return updateStatus, err
		}

		if serverStatus.UpdatePending {
			updateStatus = UpdateTropsNeeded
			log.Errorln("Traffic Ops is signaling that an update is waiting to be applied")

			if serverStatus.ParentPending && r.Cfg.WaitForParents {
				log.Errorln("Traffic Ops is signaling that my parents need an update.")
				// TODO should reval really not sleep?
				if !r.Cfg.ReportOnly && r.Cfg.Files != t3cutil.ApplyFilesFlagReval {
					log.Infof("sleeping for %ds to see if the update my parents need is cleared.", randDispSec/time.Second)
					serverStatus, err = getUpdateStatus(r.Cfg)
					if err != nil {
						return updateStatus, err
					}
					if serverStatus.ParentPending || serverStatus.ParentRevalPending {
						log.Errorln("My parents still need an update, bailing.")
						return UpdateTropsNotNeeded, nil
					} else {
						log.Debugln("The update on my parents cleared; continuing.")
					}
				}
			} else {
				log.Debugf("Processing with update: Traffic Ops server status %+v config wait-for-parents %+v", serverStatus, r.Cfg.WaitForParents)
			}
		} else if !r.Cfg.IgnoreUpdateFlag {
			log.Errorln("no queued update needs to be applied.  Running revalidation before exiting.")
			r.RevalidateWhileSleeping(metaData, cfg)
			return UpdateTropsNotNeeded, nil
		} else {
			log.Errorln("Traffic Ops is signaling that no update is waiting to be applied.")
		}

		// check local status files.
		err = r.checkStatusFiles(serverStatus.Status)
		if err != nil {
			log.Errorln(err)
		}
	}
	return updateStatus, nil
}

// CheckReloadRestart determines the final reload/restart state after all config files are processed.
func (r *TrafficOpsReq) CheckReloadRestart(data []FileRestartData) RestartData {
	rd := RestartData{}
	for _, changedFile := range data {
		rd.TrafficCtlReload = rd.TrafficCtlReload || changedFile.TrafficCtlReload
		rd.SysCtlReload = rd.SysCtlReload || changedFile.SysCtlReload
		rd.NtpdRestart = rd.NtpdRestart || changedFile.NtpdRestart
		rd.TeakdRestart = rd.TeakdRestart || changedFile.TeakdRestart
		rd.TrafficServerRestart = rd.TrafficServerRestart || changedFile.TrafficServerRestart
		rd.RemapConfigReload = rd.RemapConfigReload || changedFile.RemapConfigReload
		rd.HitchReload = rd.HitchReload || changedFile.HitchReload
		rd.VarnishReload = rd.VarnishReload || changedFile.VarnishReload
	}
	return rd
}

// ProcessConfigFiles processes all config files retrieved from Traffic Ops.
func (r *TrafficOpsReq) ProcessConfigFiles(metaData *t3cutil.ApplyMetaData) (UpdateStatus, error) {
	var updateStatus UpdateStatus = UpdateTropsNotNeeded
	var auditErrors []string

	log.Infoln(" ======== Start processing config files ========")

	filesAdding := []string{} // list of file names being added, needed for verification.
	for fileName := range r.configFiles {
		filesAdding = append(filesAdding, fileName)
	}

	for _, cfg := range r.configFiles {
		// add service metadata
		if strings.Contains(cfg.Path, "/opt/trafficserver/") || strings.Contains(cfg.Dir, "udev") {
			cfg.Service = "trafficserver"
			if !r.Cfg.InstallPackages && !r.IsPackageInstalled("trafficserver") {
				log.Errorln("Not installing packages, but trafficserver isn't installed. Continuing.")
			}
		} else if strings.Contains(cfg.Path, "/opt/ort") && strings.Contains(cfg.Name, "12M_facts") {
			cfg.Service = "puppet"
		} else if strings.Contains(cfg.Path, "cron") || strings.Contains(cfg.Name, "sysctl.conf") || strings.Contains(cfg.Name, "50-ats.rules") || strings.Contains(cfg.Name, "cron") {
			cfg.Service = "system"
		} else if strings.Contains(cfg.Path, "ntp.conf") {
			cfg.Service = "ntpd"
		} else {
			cfg.Service = "unknown"
		}

		log.Debugf("About to process config file: %s, service: %s\n", cfg.Path, cfg.Service)

		err := r.checkConfigFile(cfg, filesAdding)
		if err != nil {
			log.Errorln(err)
			r.configFiles[cfg.Name].AuditError = err.Error()
		}
	}

	changesRequired := 0
	shouldRestartReload := ShouldReloadRestart{[]FileRestartData{}}

	for _, cfg := range r.configFiles {
		metaData.OwnedFilePaths = append(metaData.OwnedFilePaths, cfg.Path) // all config files are added to OwnedFiles, even if they aren't changed on disk.

		if cfg.ChangeNeeded &&
			!cfg.ChangeApplied &&
			cfg.AuditComplete &&
			!cfg.PreReqFailed &&
			!cfg.AuditFailed {

			changesRequired++
			if cfg.Name == "plugin.config" && r.configFiles["remap.config"].PreReqFailed {
				updateStatus = UpdateTropsFailed
				log.Errorln("plugin.config changed however, prereqs failed for remap.config so I am skipping updates for plugin.config")
				continue
			} else if cfg.Name == "remap.config" && r.configFiles["plugin.config"].PreReqFailed {
				updateStatus = UpdateTropsFailed
				log.Errorln("remap.config changed however, prereqs failed for plugin.config so I am skipping updates for remap.config")
				continue
			} else if cfg.Name == "ip_allow.config" && !r.Cfg.UpdateIPAllow {
				log.Warnln("ip_allow.config changed, not updating! Run with --mode=badass or --syncds-updates-ipallow=true to update!")
				continue
			} else {
				log.Debugf("All Prereqs passed for replacing %s on disk with that in Traffic Ops.\n", cfg.Name)
				reData, err := r.replaceCfgFile(cfg)
				if err != nil {
					log.Errorf("failed to replace the config file, '%s',  on disk with data in Traffic Ops.\n", cfg.Name)
				}
				shouldRestartReload.ReloadRestart = append(shouldRestartReload.ReloadRestart, *reData)
			}
		} else if cfg.AuditFailed {
			auditErrors = append(auditErrors, cfg.AuditError)
			log.Warnf("audit failed for config file: %v Error: %s", cfg.Name, cfg.AuditError)
			updateStatus = UpdateTropsFailed
		}
	}

	if updateStatus == UpdateTropsFailed {
		return UpdateTropsFailed, errors.New(strings.Join(auditErrors, "\n"))
	}

	r.RestartData = r.CheckReloadRestart(shouldRestartReload.ReloadRestart)

	if 0 < len(r.changedFiles) {
		log.Infof("Final state: remap.config: %t reload: %t restart: %t ntpd: %t sysctl: %t", r.RemapConfigReload, r.TrafficCtlReload, r.TrafficServerRestart, r.NtpdRestart, r.SysCtlReload)
	}

	if updateStatus != UpdateTropsFailed && changesRequired > 0 {
		return UpdateTropsNeeded, nil
	}

	return updateStatus, nil
}

// ProcessPackages retrieves a list of required RPM's from Traffic Ops
// and determines which need to be installed or removed on the cache.
func (r *TrafficOpsReq) ProcessPackages() error {
	log.Infoln("Calling ProcessPackages")
	// get the package list for this cache from Traffic Ops.
	pkgs, err := getPackages(r.Cfg)
	if err != nil {
		return errors.New("getting packages: " + err.Error())
	}
	log.Infof("ProcessPackages got %+v\n", pkgs)

	var install []string   // install package list.
	var uninstall []string // uninstall package list
	// loop through the package list to build an install and uninstall list.
	for ii := range pkgs {
		var instpkg string // installed package
		var reqpkg string  // required package
		log.Infof("Processing package %s-%s\n", pkgs[ii].Name, pkgs[ii].Version)
		// check to see if any package by name is installed.
		arr, err := util.PackageInfo("pkg-query", pkgs[ii].Name)
		if err != nil {
			return errors.New("PackgeInfo pkg-query: " + err.Error())
		}
		// go needs the ternary operator :)
		if len(arr) == 1 {
			instpkg = arr[0]
		} else {
			instpkg = ""
		}
		// check if the full package version is installed
		fullPackage := pkgs[ii].Name + "-" + pkgs[ii].Version

		if r.Cfg.InstallPackages {
			if instpkg == fullPackage {
				log.Infof("%s Currently installed and not marked for removal\n", reqpkg)
				r.Pkgs[fullPackage] = true
				continue
			} else if instpkg != "" { // the installed package needs upgrading.
				log.Infof("%s Currently installed and marked for removal\n", instpkg)
				uninstall = append(uninstall, instpkg)
				// the required package needs installing.
				log.Infof("%s is Not installed and is marked for installation.\n", fullPackage)
				install = append(install, fullPackage)
				// get a list of packages that depend on this one and mark dependencies
				// for deletion.
				arr, err = util.PackageInfo("pkg-requires", instpkg)
				if err != nil {
					return errors.New("PackgeInfo pkg-requires: " + err.Error())
				}
				if len(arr) > 0 {
					for jj := range arr {
						log.Infof("%s is Currently installed and depends on %s and needs to be removed.", arr[jj], instpkg)
						uninstall = append(uninstall, arr[jj])
					}
				}
			} else {
				// the required package needs installing.
				log.Infof("%s is Not installed and is marked for installation.\n", fullPackage)
				log.Errorf("%s is Not installed and is marked for installation.\n", fullPackage)
				install = append(install, fullPackage)
			}
		} else {
			// Only check if packages exist and complain if they are wrong.
			if instpkg == fullPackage {
				log.Infof("%s Currently installed.\n", reqpkg)
				r.Pkgs[fullPackage] = true
				continue
			} else if instpkg != "" { // the installed package needs upgrading.
				log.Errorf("%s Wrong version currently installed.\n", instpkg)
				r.Pkgs[instpkg] = true
			} else {
				// the required package needs installing.
				log.Errorf("%s is Not installed.\n", fullPackage)
			}
		}
	}

	log.Debugf("number of packages requiring installation: %d\n", len(install))
	if r.Cfg.ReportOnly {
		log.Errorf("number of packages requiring installation: %d\n", len(install))
	}
	log.Debugf("number of packages requiring removal: %d\n", len(uninstall))
	if r.Cfg.ReportOnly {
		log.Errorf("number of packages requiring removal: %d\n", len(uninstall))
	}

	if r.Cfg.InstallPackages {
		log.Debugf("number of packages requiring installation: %d\n", len(install))
		if r.Cfg.ReportOnly {
			log.Errorf("number of packages requiring installation: %d\n", len(install))
		}
		log.Debugf("number of packages requiring removal: %d\n", len(uninstall))
		if r.Cfg.ReportOnly {
			log.Errorf("number of packages requiring removal: %d\n", len(uninstall))
		}

		if len(install) > 0 {
			for ii := range install {
				result, err := util.PackageAction("info", install[ii])
				if err != nil || result != true {
					return errors.New("Package " + install[ii] + " is not available to install: " + err.Error())
				}
			}
			log.Infoln("All packages available.. proceding..")

			// uninstall packages marked for removal
			if len(install) > 0 && r.Cfg.InstallPackages {
				for jj := range uninstall {
					log.Infof("Uninstalling %s\n", uninstall[jj])
					r, err := util.PackageAction("remove", uninstall[jj])
					if err != nil {
						return errors.New("Unable to uninstall " + uninstall[jj] + " : " + err.Error())
					} else if r == true {
						log.Infof("Package %s was uninstalled\n", uninstall[jj])
					}
				}

				// install the required packages
				for jj := range install {
					pkg := install[jj]
					log.Infof("Installing %s\n", pkg)
					result, err := util.PackageAction("install", pkg)
					if err != nil {
						return errors.New("Unable to install " + pkg + " : " + err.Error())
					} else if result == true {
						r.Pkgs[pkg] = true
						r.installedPkgs[pkg] = struct{}{}
						log.Infof("Package %s was installed\n", pkg)
					}
				}
			}
		}
		if r.Cfg.ReportOnly && len(install) > 0 {
			for ii := range install {
				log.Errorf("\nIn Report mode and %s needs installation.\n", install[ii])
				return errors.New("In Report mode and packages need installation")
			}
		}
	}
	return nil
}

func pkgMetaDataToMap(pmd []string) map[string]bool {
	pkgMap := map[string]bool{}
	for _, pkg := range pmd {
		pkgMap[pkg] = true
	}
	return pkgMap
}

func pkgMatch(pkgMetaData []string, pk string) bool {
	for _, pkg := range pkgMetaData {
		if strings.Contains(pk, pkg) {
			return true
		}
	}
	return false

}

// ProcessPackagesWithMetaData will attempt to get installed package data from
// t3c-apply-metadata.json and log the results.
func (r *TrafficOpsReq) ProcessPackagesWithMetaData(packageMetaData []string) error {
	pkgs, err := getPackages(r.Cfg)
	pkgMdataMap := pkgMetaDataToMap(packageMetaData)
	if err != nil {
		return fmt.Errorf("getting packages: %w", err)
	}
	for _, pkg := range pkgs {
		fullPackage := pkg.Name + "-" + pkg.Version
		if pkgMdataMap[fullPackage] {
			log.Infof("package %s is assumed to be installed according to metadata file", fullPackage)
			r.Pkgs[fullPackage] = true
		} else if pkgMatch(packageMetaData, pkg.Name) {
			log.Infof("package %s is assumed to be installed according to metadata, but doesn't match traffic ops pkg", fullPackage)
			r.Pkgs[fullPackage] = true
		} else {
			log.Infof("package %s does not appear to be installed.", pkg.Name+"-"+pkg.Version)
		}
	}
	return nil
}

func (r *TrafficOpsReq) RevalidateWhileSleeping(metaData *t3cutil.ApplyMetaData, cfg config.Cfg) (UpdateStatus, error) {
	updateStatus, err := r.CheckRevalidateState(true)
	if err != nil {
		return updateStatus, err
	}
	if updateStatus != 0 {
		r.Cfg.Files = t3cutil.ApplyFilesFlagReval
		// TODO verify? This is for revalidating after a syncds, so we probably do want to wait for parents here, and users probably don't for the main syncds run. But, this feels surprising.
		// The better solution is to gut the RevalidateWhileSleeping stuff, once TO can handle more load
		r.Cfg.WaitForParents = true

		err = r.GetConfigFileList()
		if err != nil {
			return updateStatus, err
		}

		updateStatus, err := r.ProcessConfigFiles(metaData)
		if err != nil {
			t3cutil.WriteActionLog(t3cutil.ActionLogActionUpdateFilesReval, t3cutil.ActionLogStatusFailure, metaData)
			return updateStatus, err
		} else {
			t3cutil.WriteActionLog(t3cutil.ActionLogActionUpdateFilesReval, t3cutil.ActionLogStatusSuccess, metaData)
		}

		if err := r.StartServices(&updateStatus, metaData, cfg); err != nil {
			return updateStatus, errors.New("failed to start services: " + err.Error())
		}

		if err := r.UpdateTrafficOps(&updateStatus); err != nil {
			log.Errorf("failed to update Traffic Ops: %s\n", err.Error())
		}

		r.TrafficCtlReload = false
	}

	return updateStatus, nil
}

// StartServices reloads, restarts, or starts ATS as necessary,
// according to the changed config files and run mode.
// Returns nil on success or any error.
func (r *TrafficOpsReq) StartServices(syncdsUpdate *UpdateStatus, metaData *t3cutil.ApplyMetaData, cfg config.Cfg) error {
	serviceNeeds := t3cutil.ServiceNeedsNothing
	if r.Cfg.ServiceAction == t3cutil.ApplyServiceActionFlagRestart {
		serviceNeeds = t3cutil.ServiceNeedsRestart
	} else {
		err := error(nil)
		if serviceNeeds, err = checkReload(r.changedFiles); err != nil {
			return errors.New("determining if service needs restarted - not reloading or restarting! : " + err.Error())
		}
	}

	log.Infof("t3c-check-reload returned '%+v'\n", serviceNeeds)

	// We have our own internal knowledge of files that have been modified as well
	// If check-reload does not know about these and we do, then we should initiate
	// a reload as well
	if serviceNeeds != t3cutil.ServiceNeedsRestart && serviceNeeds != t3cutil.ServiceNeedsReload {
		if r.TrafficCtlReload || r.RemapConfigReload || r.VarnishReload {
			log.Infof("ATS config files unchanged, we updated files via t3c-apply, ATS needs reload")
			serviceNeeds = t3cutil.ServiceNeedsReload
		}
	}
	packageName := "trafficserver"
	if cfg.CacheType == "varnish" {
		packageName = "varnish"
	}

	if (serviceNeeds == t3cutil.ServiceNeedsRestart || serviceNeeds == t3cutil.ServiceNeedsReload) && !r.IsPackageInstalled(packageName) {
		// TODO try to reload/restart anyway? To allow non-RPM installs?
		return errors.New(packageName + " needs " + serviceNeeds.String() + " but is not installed.")
	}

	svcStatus, _, err := util.GetServiceStatus(packageName)
	if err != nil {
		return errors.New("getting trafficserver service status: " + err.Error())
	}

	if r.Cfg.ReportOnly {
		if serviceNeeds == t3cutil.ServiceNeedsRestart {
			log.Errorln("ATS configuration has changed.  The new config will be picked up the next time ATS is started.")
		} else if serviceNeeds == t3cutil.ServiceNeedsReload {
			log.Errorln("ATS configuration has changed. 'traffic_ctl config reload' needs to be run")
		}
		return nil
	} else if r.Cfg.ServiceAction == t3cutil.ApplyServiceActionFlagRestart {
		startStr := "restart"
		if svcStatus != util.SvcRunning {
			startStr = "start"
		}
		if _, err := util.ServiceStart(packageName, startStr); err != nil {
			t3cutil.WriteActionLog(t3cutil.ActionLogActionATSRestart, t3cutil.ActionLogStatusFailure, metaData)
			return errors.New("failed to restart trafficserver")
		}
		t3cutil.WriteActionLog(t3cutil.ActionLogActionATSRestart, t3cutil.ActionLogStatusSuccess, metaData)
		log.Infoln("trafficserver has been " + startStr + "ed")

		if !r.Cfg.NoConfirmServiceAction {
			log.Infoln("confirming ATS restart succeeded")
			if err := doTail(r.Cfg, TailDiagsLogRelative, ".*", tailRestartEnd, TailRestartTimeOutMS); err != nil {
				log.Errorln("error running tail")
			}
		} else {
			log.Infoln("skipping ATS restart success confirmation")
		}
		if *syncdsUpdate == UpdateTropsNeeded {
			*syncdsUpdate = UpdateTropsSuccessful
		}
		return nil // we restarted, so no need to reload
	} else if r.Cfg.ServiceAction == t3cutil.ApplyServiceActionFlagReload {
		if serviceNeeds == t3cutil.ServiceNeedsRestart {
			if *syncdsUpdate == UpdateTropsNeeded {
				*syncdsUpdate = UpdateTropsSuccessful
			}
			log.Errorln("ATS configuration has changed.  The new config will be picked up the next time ATS is started.")
		} else if serviceNeeds == t3cutil.ServiceNeedsReload {
			log.Infoln("ATS configuration has changed, Running 'traffic_ctl config reload' now.")
			reloadCommand := config.TSHome + config.TrafficCtl
			reloadArgs := []string{"config", "reload"}
			if cfg.CacheType == "varnish" {
				reloadCommand = "/usr/sbin/varnishreload"
				reloadArgs = []string{}
			}
			if _, _, err := util.ExecCommand(reloadCommand, reloadArgs...); err != nil {
				t3cutil.WriteActionLog(t3cutil.ActionLogActionATSReload, t3cutil.ActionLogStatusFailure, metaData)

				if *syncdsUpdate == UpdateTropsNeeded {
					*syncdsUpdate = UpdateTropsFailed
				}
				return errors.New("ATS configuration has changed and 'traffic_ctl config reload' failed, check ATS logs: " + err.Error())
			}
			t3cutil.WriteActionLog(t3cutil.ActionLogActionATSReload, t3cutil.ActionLogStatusSuccess, metaData)

			if *syncdsUpdate == UpdateTropsNeeded {
				*syncdsUpdate = UpdateTropsSuccessful
			}
			log.Infoln("ATS 'traffic_ctl config reload' was successful")

			if !r.Cfg.NoConfirmServiceAction {
				log.Infoln("confirming ATS reload succeeded")
				if err := doTail(r.Cfg, TailDiagsLogRelative, tailMatch, tailReloadEnd, TailReloadTimeOutMS); err != nil {
					log.Errorln("error running tail: ", err)
				}
			} else {
				log.Infoln("skipping ATS reload success confirmation")
			}
		}
		if *syncdsUpdate == UpdateTropsNeeded {
			*syncdsUpdate = UpdateTropsSuccessful
		}
		return nil
	}
	return nil
}

func (r *TrafficOpsReq) ShowUpdateStatus(flagType []string, start time.Time, curSetting, newSetting bool) {
	for _, flag := range flagType {
		log.Infof("%s flag currently set to %v, setting to %v took %v", flag, curSetting, newSetting, time.Since(start).Round(time.Millisecond))
	}
}

func (r *TrafficOpsReq) UpdateTrafficOps(syncdsUpdate *UpdateStatus) error {
	var performUpdate bool

	serverStatus, err := getUpdateStatus(r.Cfg)
	if err != nil {
		return errors.New("failed to update Traffic Ops: " + err.Error())
	}

	if *syncdsUpdate == UpdateTropsNotNeeded && (serverStatus.UpdatePending == true || serverStatus.RevalPending == true) {
		performUpdate = true
		log.Errorln("Traffic Ops is signaling that an update is ready to be applied but, none was found! Clearing update state in Traffic Ops anyway.")
	} else if *syncdsUpdate == UpdateTropsNotNeeded {
		log.Errorln("Traffic Ops does not require an update at this time")
		return nil
	} else if *syncdsUpdate == UpdateTropsFailed {
		log.Errorln("Traffic Ops requires an update but, applying the update locally failed.  Traffic Ops is not being updated.")
		return nil
	} else if *syncdsUpdate == UpdateTropsSuccessful {
		performUpdate = true
		log.Errorln("Traffic Ops requires an update and it was applied successfully.  Clearing update state in Traffic Ops.")
	}

	if !performUpdate {
		return nil
	}
	if r.Cfg.ReportOnly {
		log.Errorln("In Report mode and Traffic Ops needs updated you should probably do that manually.")
		return nil
	}

	// TODO: The boolean flags/representation can be removed after ATC (v7.0+)
	if !r.Cfg.ReportOnly && !r.Cfg.NoUnsetUpdateFlag {
		start := time.Now()
		apply := []string{}
		var b bool
		if r.Cfg.Files == t3cutil.ApplyFilesFlagAll {
			b = false
			apply = append(apply, "update")
			log.Infof("Update flag currently set to %v, setting to %v", serverStatus.UpdatePending, b)
			err = sendUpdate(r.Cfg, serverStatus.ConfigUpdateTime, nil, &b, nil)
		} else if r.Cfg.Files == t3cutil.ApplyFilesFlagReval {
			b = false
			apply = append(apply, t3cutil.ApplyFilesFlagReval.String())
			log.Infof("Reval flag currently set to %v, setting to %v", serverStatus.RevalPending, b)
			err = sendUpdate(r.Cfg, nil, serverStatus.RevalidateUpdateTime, nil, &b)
		}
		if err != nil {
			return errors.New("Traffic Ops Update failed: " + err.Error())
		}
		log.Infoln("Traffic Ops has been updated.")
		r.ShowUpdateStatus(apply, start, serverStatus.UpdatePending, b)
	}
	return nil
}
