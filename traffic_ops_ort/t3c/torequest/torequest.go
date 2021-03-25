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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/t3c/config"
	"github.com/apache/trafficcontrol/traffic_ops_ort/t3c/util"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type UpdateStatus int

const (
	UpdateTropsNotNeeded  UpdateStatus = 0
	UpdateTropsNeeded     UpdateStatus = 1
	UpdateTropsSuccessful UpdateStatus = 2
	UpdateTropsFailed     UpdateStatus = 3
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type TrafficOpsReq struct {
	Cfg                  config.Cfg
	pkgs                 map[string]bool // map of installed packages
	plugins              map[string]bool // map of verified plugins
	configFiles          map[string]*ConfigFile
	baseBackupDir        string
	TrafficCtlReload     bool   // a traffic_ctl_reload is required
	SysCtlReload         bool   // a reload of the sysctl.conf is required
	NtpdRestart          bool   // ntpd needs restarting
	TeakdRestart         bool   // a restart of teakd is required
	TrafficServerRestart bool   // a trafficserver restart is required
	RemapConfigReload    bool   // remap.config should be reloaded
	unixTimeStr          string // unix time string at program startup.
}

type ConfigFile struct {
	Header            textproto.MIMEHeader
	Name              string // file name
	Dir               string // install directory
	Path              string // full path
	Service           string // service assigned to
	CfgBackup         string // location to backup the config at 'Path'
	TropsBackup       string // location to backup the TrafficOps Version
	AuditComplete     bool   // audit is complete
	AuditFailed       bool   // audit failed
	ChangeApplied     bool   // a change has been applied
	ChangeNeeded      bool   // change required
	PreReqFailed      bool   // failed plugin prerequiste check
	RemapPluginConfig bool   // file is a remap plugin config file
	Body              []byte
	Perm              os.FileMode // default file permissions
	Uid               int         // owner uid, default is 0
	Gid               int         // owner gid, default is 0
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
		fmt.Printf("Name: %s, Dir: %s, Service: %s\n",
			cfg.Name, cfg.Dir, cfg.Service)
	}
}

// NewTrafficOpsReq returns a new TrafficOpsReq object.
func NewTrafficOpsReq(cfg config.Cfg) *TrafficOpsReq {
	unixTimeString := strconv.FormatInt(time.Now().Unix(), 10)

	return &TrafficOpsReq{
		Cfg:           cfg,
		pkgs:          make(map[string]bool),
		plugins:       make(map[string]bool),
		configFiles:   make(map[string]*ConfigFile),
		baseBackupDir: config.TmpBase + "/" + unixTimeString,
		unixTimeStr:   unixTimeString,
	}
}

// atsTcExec is a wrapper to run an atstccfg command.
func (r *TrafficOpsReq) atsTcExec(cmdstr string) ([]byte, error) {
	log.Debugf("cmdstr: %s\n", cmdstr)
	result, err := r.atsTcExecCommand(cmdstr, -1, -1)
	return result, err
}

// atsTcExecCommand is used to run the atstccfg command.
func (r *TrafficOpsReq) atsTcExecCommand(cmdstr string, queueState int, revalState int) ([]byte, error) {
	// adjust log locations used for atstccfg
	// cannot use stdout as this will cause json parsing errors.
	errorLocation := r.Cfg.LogLocationErr
	if errorLocation == "stdout" {
		errorLocation = "stderr"
		log.Infoln("atstccfg error logging has been re-directed to 'stderr'")
	}
	infoLocation := r.Cfg.LogLocationInfo
	if infoLocation == "stdout" {
		infoLocation = "stderr"
		log.Infoln("atstccfg info logging has been re-directed to 'stderr'")
	}
	warningLocation := r.Cfg.LogLocationWarn
	if warningLocation == "stdout" {
		warningLocation = "stderr"
		log.Infoln("atstccfg warning logging has been re-directed to 'stderr'")
	}

	args := []string{
		"--dir=" + config.TSConfigDir,
		"--traffic-ops-timeout-milliseconds=" + strconv.FormatInt(int64(r.Cfg.TOTimeoutMS), 10),
		"--traffic-ops-disable-proxy=" + strconv.FormatBool(r.Cfg.ReverseProxyDisable),
		"--traffic-ops-user=" + r.Cfg.TOUser,
		"--traffic-ops-password=" + r.Cfg.TOPass,
		"--traffic-ops-url=" + r.Cfg.TOURL,
		"--cache-host-name=" + r.Cfg.CacheHostName,
		"--log-location-error=" + errorLocation,
		"--log-location-info=" + infoLocation,
		"--log-location-warning=" + warningLocation,
	}

	if r.Cfg.TOInsecure == true {
		args = append(args, "--traffic-ops-insecure")
	}

	if r.Cfg.DNSLocalBind {
		args = append(args, "--dns-local-bind")
	}

	switch cmdstr {
	case "chkconfig":
		args = append(args, "--get-data=chkconfig")
	case "packages":
		args = append(args, "--get-data=packages")
	case "statuses":
		args = append(args, "--get-data=statuses")
	case "system-info":
		args = append(args, "--get-data=system-info")
	case "update-status":
		args = append(args, "--get-data=update-status")
	case "send-update":
		var queueStatus string = "false"
		var revalStatus string = "false"
		if queueState > 0 {
			queueStatus = "true"
		}
		if revalState > 0 {
			revalStatus = "true"
		}
		args = append(args, "--set-queue-status="+queueStatus)
		args = append(args, "--set-reval-status="+revalStatus)
	case "get-config-files":
		if r.Cfg.RunMode == config.Revalidate {
			args = append(args, "--revalidate-only")
		}
	default:
		return nil, errors.New("invalid command '" + cmdstr + "'")
	}

	cmd := exec.Command(config.AtsTcConfig, args...)
	var outbuf bytes.Buffer
	var errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil {
		return nil, errors.New("Error from atstccfg: " + err.Error() + ": " + errbuf.String())
	}

	return outbuf.Bytes(), nil
}

// backUpFile makes a backup of a config file.
func (r *TrafficOpsReq) backUpFile(cfg *ConfigFile) error {

	// init backup directories
	configBkupDir := filepath.Join(r.baseBackupDir, cfg.Service, "/config_bkp")
	cfg.CfgBackup = filepath.Join(configBkupDir, cfg.Name)
	tropsBkupDir := filepath.Join(r.baseBackupDir, cfg.Service, "/config_trops")
	cfg.TropsBackup = filepath.Join(tropsBkupDir, cfg.Name)

	fileExists, _ := util.FileExists(cfg.Path)
	if fileExists {
		// create config file backup directory
		err := os.MkdirAll(configBkupDir, 0755)
		if err != nil {
			return errors.New("Unable to create config backup directory '" + configBkupDir + "': " + err.Error())
		}

		// backup the config file
		data, err := r.readCfgFile(cfg, "")
		if err != nil {
			return errors.New("Unable to read the config file '" + cfg.Path + "': " + err.Error())
		}
		_, err = r.writeCfgFile(cfg, configBkupDir, data)
		if err != nil {
			return errors.New("Failed to write the config file: " + err.Error())
		}

	} else {
		log.Debugf("Config file: %s doesn't exist. No need to back up.\n", cfg.Name)
	}

	// backup the traffic ops file.
	err := os.MkdirAll(tropsBkupDir, 0755)
	if err != nil {
		return errors.New("Unable to create Trops backup directory '" + tropsBkupDir + "': " + err.Error())
	}
	_, err = r.writeCfgFile(cfg, tropsBkupDir, cfg.Body)
	if err != nil {
		return errors.New("Failed to write the config file from traffic ops: " + err.Error())
	}

	// backup the current config file.
	return nil
}

// checkConfigFile checks and audits config files.
func (r *TrafficOpsReq) checkConfigFile(cfg *ConfigFile) error {
	if cfg.Name == "" {
		cfg.AuditFailed = true
		return errors.New("Config file name is empty is empty, skipping further checks.")
	}

	if cfg.Dir == "" {
		return errors.New("No location information for " + cfg.Name)
	}
	// return if audit has already been done.
	if cfg.AuditComplete == true {
		return nil
	}

	if !util.MkDir(cfg.Dir, r.Cfg) {
		return errors.New("Unable to create the directory '" + cfg.Dir + " for " + "'" + cfg.Name + "'")
	}

	log.Debugf("======== Start processing config file: %s ========\n", cfg.Name)

	if cfg.Name == "remap.config" {
		err := r.processRemapOverrides(cfg)
		if err != nil {
			return err
		}
	}

	// perform plugin verification
	if cfg.Name == "remap.config" || cfg.Name == "plugin.config" {
		err := r.verifyPlugins(cfg)
		if err != nil {
			return err
		}
	}

	// apply traffic ops data filters in preparation for comparison
	// to data on disk.
	tropsData := strings.Split(string(cfg.Body), "\n")
	tropsData = unencodeFilter(tropsData)
	tropsData = commentsFilter(tropsData)

	var diskData []string
	fileExists, _ := util.FileExists(cfg.Path)
	if fileExists {
		data, err := r.readCfgFile(cfg, "")
		if err != nil {
			return errors.New("reading from '" + cfg.Path + "' failed: " + err.Error())
		}
		diskData = strings.Split(string(data), "\n")
	} else { // file doesn't exist on, it's new from Traffic Ops.
		cfg.AuditComplete = true
		cfg.ChangeNeeded = true
		log.Infof("No such file on disk, '%s'\n", cfg.Path)
	}

	// apply disk file data filters in preparation for comparison
	// to data from traffic ops
	diskData = unencodeFilter(diskData)
	diskData = commentsFilter(diskData)

	// apply final new line filters disk and traffic ops data for comparison
	disk := strings.Join(diskData, "\n")
	disk = newLineFilter(disk)

	trops := strings.Join(tropsData, "\n")
	trops = newLineFilter(trops)

	if disk != trops {
		cfg.ChangeNeeded = true
		log.Infof("change needed to %s\n", cfg.Name)
		err := r.backUpFile(cfg)
		if err != nil {
			return errors.New("unable to back up '" + cfg.Name + "': " + err.Error())
		}
	} else {
		cfg.ChangeNeeded = false
		log.Infof("All lines match TrOps for config file: %s\n", cfg.Name)
	}

	if cfg.Name == "50-ats.rules" {
		err := r.processUdevRules(cfg)
		if err != nil {
			return errors.New("unable to process udev rules in '" + cfg.Name + "': " + err.Error())
		}
	}

	log.Infof("======== End processing config file: %s for service: %s ========\n", cfg.Name, cfg.Service)
	cfg.AuditComplete = true

	return nil
}

// checkPlugin verifies ATS plugin requirements are satisfied.
func (r *TrafficOpsReq) checkPlugin(plugin string) error {
	// already verified
	if r.plugins[plugin] == true {
		return nil
	}
	pluginFile := filepath.Join(config.TSHome, "/libexec/trafficserver/", plugin)
	pkgs, err := util.PackageInfo("pkg-provides", pluginFile)
	if err != nil {
		return errors.New("unable to verify plugin " + pluginFile + ": " + err.Error())
	}
	if len(pkgs) == 0 || pkgs == nil { // no package is installed that provides the plugin.
		return errors.New(plugin + ": Package for plugin: " + plugin + ", is not installed.")
	}
	_, ok := r.pkgs[pkgs[0]]
	if !ok {
		return errors.New(plugin + ": Package for plugin: " + plugin + ", is not installed.")
	}
	return nil
}

// checkStatusFiles insures that the cache status files reflect
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
	statuses, err := r.getStatuses()
	if err != nil {
		return fmt.Errorf("could not retrieves a statuses list from Traffic Ops: %s\n", err)
	}

	for f := range statuses {
		otherStatus := filepath.Join(config.StatusDir, statuses[f])
		if otherStatus == statusFile {
			continue
		}
		fileExists, _ := util.FileExists(otherStatus)
		if r.Cfg.RunMode != config.Report && fileExists {
			log.Errorf("Removing other status file %s that exists\n", otherStatus)
			err = os.Remove(otherStatus)
			if err != nil {
				log.Errorf("Error removing %s: %s\n", otherStatus, err)
			}
		}
	}

	if r.Cfg.RunMode != config.Report {
		if !util.MkDir(config.StatusDir, r.Cfg) {
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

// getStatuses fetches a list of cache statuses from Traffic ops.
func (r *TrafficOpsReq) getStatuses() ([]string, error) {
	var statuses []tc.StatusNullable
	sl := []string{}
	out, err := r.atsTcExec("statuses")
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	if err = json.Unmarshal(out, &statuses); err != nil {
		log.Errorln(err)
		return nil, err
	} else {
		for val := range statuses {
			if statuses[val].Name != nil {
				sl = append(sl, *statuses[val].Name)
			}
		}
	}

	return sl, nil
}

// getUpdateStatus retrieves the update statuse for a cache from Traffic Ops.
func (r *TrafficOpsReq) getUpdateStatus() (*tc.ServerUpdateStatus, error) {
	var status tc.ServerUpdateStatus
	out, err := r.atsTcExec("update-status")
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	if err = json.Unmarshal(out, &status); err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugf("ServerUpdateStatus: %#v\n", status)
	return &status, nil
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
	fs, err := ioutil.ReadDir("/proc/fs/ext4")
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

// replaceCfgFile replaces an ATS configuration file with one from Traffic Ops.
func (r *TrafficOpsReq) replaceCfgFile(cfg *ConfigFile) error {
	if r.Cfg.RunMode == config.BadAss || r.Cfg.RunMode == config.SyncDS || r.Cfg.RunMode == config.Revalidate {

		log.Infof("Copying '%s' to '%s'\n", cfg.TropsBackup, cfg.Path)
		data, err := util.ReadFile(cfg.TropsBackup)
		if err != nil {
			return errors.New("Unable to read the config file '" + cfg.TropsBackup + "': " + err.Error())
		}
		_, err = r.writeCfgFile(cfg, "", data)
		if err != nil {
			return errors.New("Failed to write the new config file: " + err.Error())
		}
		cfg.ChangeApplied = true

		r.TrafficCtlReload = r.TrafficCtlReload ||
			strings.HasSuffix(cfg.Dir, "trafficserver") ||
			cfg.RemapPluginConfig ||
			cfg.Name == "remap.config" ||
			cfg.Name == "ssl_multicert.config" ||
			strings.HasPrefix(cfg.Name, "url_sig_") ||
			strings.HasPrefix(cfg.Name, "uri_signing") ||
			strings.HasPrefix(cfg.Name, "hdr_rw_") ||
			(strings.HasSuffix(cfg.Dir, "ssl") && strings.HasSuffix(cfg.Name, ".cer")) ||
			(strings.HasSuffix(cfg.Dir, "ssl") && strings.HasSuffix(cfg.Name, ".key"))

		r.TrafficServerRestart = cfg.Name == "plugin.config"
		r.RemapConfigReload = cfg.RemapPluginConfig || cfg.Name == "remap.config"
		r.NtpdRestart = cfg.Name == "ntpd.conf"
		r.SysCtlReload = cfg.Name == "sysctl.conf"

		log.Debugf("Setting change applied for '%s'\n", cfg.Name)

	} else {
		log.Infof("You elected not to replace %s with the version from Traffic Ops.\n", cfg.Name)
		cfg.ChangeApplied = false
	}
	return nil
}

func (r *TrafficOpsReq) sleepTimer(serverStatus *tc.ServerUpdateStatus) {
	randDispSec := time.Duration(0)
	revalClockSec := time.Duration(0)

	if r.Cfg.Dispersion > 0 {
		randDispSec = util.RandomDuration(r.Cfg.Dispersion) / time.Second
	}
	if r.Cfg.RevalWaitTime > 0 {
		revalClockSec = r.Cfg.RevalWaitTime / time.Second
	}

	if serverStatus.UseRevalPending && r.Cfg.RunMode != config.BadAss {
		log.Infoln("Performing a revalidation check before sleeping...")
		_, err := r.RevalidateWhileSleeping()
		if err != nil {
			log.Errorf("Revalidation check completed with error: %s\n", err)
		} else {
			log.Infoln("Revalidation check complete.")
		}
	}
	if randDispSec < revalClockSec || serverStatus.UseRevalPending == false || r.Cfg.RunMode == config.BadAss {
		log.Infof("Sleeping for %d seconds: ", randDispSec)
	} else {
		log.Infof("%d seconds until next revalidation check.\n", revalClockSec)
		log.Infof("%d seconds remaining in dispersion sleep period\n", randDispSec)
		log.Infof("Sleeping for %d seconds: ", revalClockSec)
	}

	for randDispSec > 0 {
		fmt.Printf(".")
		time.Sleep(time.Second)
		revalClockSec--
		if revalClockSec < 1 && r.Cfg.RunMode != config.BadAss && serverStatus.UseRevalPending {
			fmt.Printf("\n")
			log.Infoln("Interrupting dispersion sleep period for revalidation check.")
			_, err := r.RevalidateWhileSleeping()
			if r.Cfg.RevalWaitTime > 0 {
				revalClockSec = r.Cfg.RevalWaitTime / time.Second
			}
			if err != nil {
				log.Errorf("Revalidation check completed with error: %s\n", err)
			} else {
				log.Infoln("Revalidation check complete.")
			}
			if revalClockSec < randDispSec {
				log.Infof("Revalidation check complete. %d seconds until next revalidation check.", revalClockSec)
				log.Infof("%d seconds remaining in dispersion sleep period\n", randDispSec)
				log.Infof("Sleeping for %d seconds: ", revalClockSec)
			} else {
				log.Infof("Revalidation check complete. %d seconds remaining in dispersion sleep period.\n", randDispSec)
				log.Infof("Sleeping for %d seconds: ", randDispSec)
			}
		}
		randDispSec--
	}
	fmt.Printf("\n")
}

// writeCfgFile writes the 'data' from Traffic Ops to an ATS config file.
func (r *TrafficOpsReq) writeCfgFile(cfg *ConfigFile, dir string, data []byte) (int, error) {
	var fullFileName string

	if dir == "" {
		fullFileName = cfg.Path
	} else {
		fullFileName = dir + "/" + cfg.Name
	}

	fd, err := os.OpenFile(fullFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, cfg.Perm)
	if err != nil {
		return 0, errors.New("unable to open '" + fullFileName + "' for writing: " + err.Error())
	}

	c, err := fd.Write(data)
	if err != nil {
		return 0, errors.New("error writing to '" + fullFileName + "': " + err.Error())
	}
	fd.Close()

	if cfg.Service == "trafficserver" {
		err = os.Chown(fullFileName, cfg.Uid, cfg.Gid)
	} else {
		err = os.Chown(fullFileName, 0, 0)
	}
	if err != nil {
		return 0, errors.New("error changing ownership on '" + fullFileName + "': " + err.Error())
	}

	return c, nil
}

// verifyPlugins is used to verify that the plugin found
// in plugin.config or from a remap.config is installed.
func (r *TrafficOpsReq) verifyPlugins(cfg *ConfigFile) error {

	log.Debugf("Checking plugins for %s\n", cfg.Name)

	str := string(cfg.Body)
	lines := strings.Split(str, "\n")
	if cfg.Name == "plugin.config" {
		for ii := range lines {
			line := strings.TrimSpace(lines[ii])
			if strings.HasPrefix(line, "#") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) > 0 {
				plugin := fields[0]
				// already verified
				if r.plugins[plugin] == true {
					continue
				}
				err := r.checkPlugin(plugin)
				if err != nil {
					cfg.PreReqFailed = true
					return err
				} else {
					r.plugins[plugin] = true
					log.Infof("Plugin %s has been verified.\n", plugin)
				}
			}
		}
	} else if cfg.Name == "remap.config" && len(lines) > 0 {
		for ii := range lines {
			line := lines[ii]
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			plugins := strings.Split(line, "@plugin=")
			if len(plugins) > 0 {
				for jj := range plugins {
					var plugin string = ""
					var plugin_config string = ""
					if jj > 0 {
						params := strings.Split(plugins[jj], "@pparam=")
						param_length := len(params)
						if param_length > 0 {
							switch param_length {
							case 1:
								plugin = strings.TrimSpace(params[0])
								plugin_config = ""
							default:
								plugin = strings.TrimSpace(params[0])
								plugin_config = strings.TrimSpace(params[1])
							}
						}
						if strings.HasSuffix(plugin, ".so") {
							// already verified
							if r.plugins[plugin] == true {
								continue
							}
							err := r.checkPlugin(plugin)
							if err != nil {
								cfg.PreReqFailed = true
								return err
							} else {
								r.plugins[plugin] = true
								log.Infof("Plugin %s has been verified.\n", plugin)
							}
						}
						if plugin_config != "" {
							if strings.HasPrefix(plugin_config, "proxy.config") || strings.HasPrefix(plugin_config, "-") {
								continue
							} else {
								plugin_config = filepath.Base(plugin_config)
								cfg, ok := r.configFiles[plugin_config]
								if ok {
									cfg.RemapPluginConfig = true
									log.Debugf("%s is a remap plugin config file\n", plugin_config)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// CheckSystemServices is used to verify that packages installed
// are enabled for startup.
func (r *TrafficOpsReq) CheckSystemServices() error {
	if r.Cfg.RunMode == config.BadAss {
		out, err := r.atsTcExec("chkconfig")
		if err != nil {
			log.Errorln(err)
			return err
		} else {
			var result []map[string]string
			if err = json.Unmarshal(out, &result); err != nil {
				return err
			}
			for ii := range result {
				name := result[ii]["name"]
				value := result[ii]["value"]
				arrv := strings.Fields(value)
				var level []string
				var enabled bool = false
				for jj := range arrv {
					nv := strings.Split(arrv[jj], ":")
					if len(nv) == 2 && strings.Contains(nv[1], "on") {
						level = append(level, nv[0])
						enabled = true
					}
				}
				if enabled == true {
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
						log.Errorf("Unable to insure %s service is enabled, SvcMananagement type is %s\n", name, r.Cfg.SvcManagement)
					}
				}
			}
		}
	}
	return nil
}

// IsPackageInstalled returns true/false if the named rpm package is installed.
// the prefix before the version is matched.
func (r *TrafficOpsReq) IsPackageInstalled(name string) bool {
	for k, v := range r.pkgs {
		if strings.HasPrefix(k, name) {
			return v
		}
	}
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

	fBytes, err := r.atsTcExec("get-config-files")
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(fBytes)
	var hdr string
	for {
		hdr, err = buf.ReadString('\n')
		if err != nil {
			return err
		} else {
			if strings.HasPrefix(hdr, "Content-Type: multipart/mixed") {
				break
			}
		}
	}
	// split on the ": " to trim off the leading space
	har := strings.Split(hdr, ": ")
	if len(har) > 1 && (har[0] != rfc.ContentType || !strings.HasPrefix(har[1], rfc.ContentTypeMultiPartMixed)) {
		return errors.New("invalid config file data received from Traffic Ops, not in multipart format.")
	}
	msg := &mail.Message{
		Header: map[string][]string{har[0]: {har[1]}},
		Body:   bytes.NewReader(buf.Bytes()),
	}
	mediaType, params, err := mime.ParseMediaType(msg.Header.Get(rfc.ContentType))
	if err != nil {
		return err
	}
	if strings.HasPrefix(mediaType, rfc.ContentTypeMultiPartMixed) {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		// if the configFiles map is not empty, create a new map
		// and the old map will be garbage collected.
		if len(r.configFiles) > 0 {
			log.Infoln("intializing a new configFiles map")
			r.configFiles = make(map[string]*ConfigFile)
		}
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			dataBytes, err := ioutil.ReadAll(p)
			if err != nil {
				return err
			}
			path := p.Header.Get("Path")
			if path != "" {
				cf := &ConfigFile{
					Header: p.Header,
					Name:   filepath.Base(path),
					Path:   path,
					Dir:    filepath.Dir(path),
					Body:   dataBytes,
					Uid:    atsUid,
					Gid:    atsGid,
					Perm:   0644,
				}
				r.configFiles[cf.Name] = cf
			}
		}
	}
	return nil
}

// GetHeaderComment looks up the tm.toolname parameter from traffic ops.
func (r *TrafficOpsReq) GetHeaderComment() string {
	var toolName string = ""
	out, err := r.atsTcExec("system-info")
	if err != nil {
		log.Errorln(err)
	} else {
		var result map[string]interface{}
		if err = json.Unmarshal(out, &result); err != nil {
			log.Errorln(err)
		} else {
			tn := result["tm.toolname"]
			if tn, ok := tn.(string); ok {
				toolName = tn
				log.Infof("Found tm.toolname: %v\n", toolName)
			} else {
				log.Errorln("Did not find tm.toolname!")
			}
		}
	}
	return toolName
}

// CheckRevalidateState retrieves and returns the revalidate status from Traffic Ops.
func (r *TrafficOpsReq) CheckRevalidateState(sleepOverride bool) (UpdateStatus, error) {
	updateStatus := UpdateTropsNotNeeded
	log.Infoln("Checking revalidate state.")

	if r.Cfg.RunMode == config.Revalidate || sleepOverride {
		serverStatus, err := r.getUpdateStatus()
		log.Infof("my status: %s\n", serverStatus.Status)
		if err != nil {
			log.Errorln(err)
			return updateStatus, err
		} else {
			if serverStatus.UseRevalPending == false {
				log.Errorln("Update URL: Instant invalidate is not enabled.  Separated revalidation requires upgrading to Traffic Ops version 2.2 and enabling this feature.")
				return UpdateTropsNotNeeded, nil
			}
			if serverStatus.RevalPending == true {
				log.Errorln("Traffic Ops is signaling that a revalidation is waiting to be applied.")
				updateStatus = UpdateTropsNeeded
				if serverStatus.ParentRevalPending == true {
					log.Errorln("Traffic Ops is signaling that my parents need to revalidate.")
					// no update needed until my parents are updated.
					updateStatus = UpdateTropsNotNeeded
				}
			} else if serverStatus.RevalPending == false && r.Cfg.RunMode == config.Revalidate {
				log.Errorln("In revalidate mode, but no update needs to be applied. I'm outta here.")
				return UpdateTropsNotNeeded, nil
			} else {
				log.Errorln("Traffic Ops is signaling that no revalidations are waiting to be applied.")
				return UpdateTropsNotNeeded, nil
			}
		}

		err = r.checkStatusFiles(serverStatus.Status)
		if err != nil {
			log.Errorln(err)
		}
	}

	return updateStatus, nil
}

// CheckSYncDSState retrieves and returns the DS Update status from Traffic Ops.
func (r *TrafficOpsReq) CheckSyncDSState() (UpdateStatus, error) {
	updateStatus := UpdateTropsNotNeeded
	randDispSec := time.Duration(0)
	if r.Cfg.Dispersion > 0 {
		randDispSec = util.RandomDuration(r.Cfg.Dispersion)
	}
	log.Debugln("Checking syncds state.")
	if r.Cfg.RunMode == config.SyncDS || r.Cfg.RunMode == config.BadAss || r.Cfg.RunMode == config.Report {
		serverStatus, err := r.getUpdateStatus()
		if err != nil {
			log.Errorln(err)
			return updateStatus, err
		}

		if serverStatus.UpdatePending {
			if r.Cfg.Dispersion > 0 {
				log.Infof("Sleeping for %ds (dispersion) before proceeding with updates.\n\n", (randDispSec / time.Second))
				r.sleepTimer(serverStatus)
			}
			updateStatus = UpdateTropsNeeded
			log.Errorln("Traffic Ops is signaling that an update is waiting to be applied")

			if serverStatus.ParentPending && r.Cfg.WaitForParents && !serverStatus.UseRevalPending {
				log.Errorln("Traffic Ops is signaling that my parents need an update.")
				if r.Cfg.RunMode == config.SyncDS {
					log.Infof("In syncds mode, sleeping for %ds to see if the update my parents need is cleared.", randDispSec/time.Second)
					r.sleepTimer(serverStatus)
					serverStatus, err = r.getUpdateStatus()
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
				log.Debugf("Traffic Ops is signaling that my parents do not need an update, or wait_for_parents is false.")
			}
		} else if r.Cfg.RunMode == config.SyncDS {
			log.Errorln("In syncds mode, but no syncds update needs to be applied.  Running revalidation before exiting.")
			r.RevalidateWhileSleeping()
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

// ProcessConfigFiles processes all config files retrieved from Traffic Ops.
func (r *TrafficOpsReq) ProcessConfigFiles() (UpdateStatus, error) {
	var updateStatus UpdateStatus = UpdateTropsNotNeeded

	log.Infoln(" ======== Start processing config files ========")

	for _, cfg := range r.configFiles {
		// add service metadata
		if strings.Contains(cfg.Path, "/opt/trafficserver/") || strings.Contains(cfg.Dir, "udev") {
			cfg.Service = "trafficserver"
			if r.Cfg.RunMode == config.SyncDS && !r.IsPackageInstalled("trafficserver") {
				return UpdateTropsFailed, errors.New("In syncds mode, but trafficserver isn't installed. Bailing.")
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

		log.Debugf("In %s mode, I'm about to process config file: %s, service: %s\n", r.Cfg.RunMode, cfg.Path, cfg.Service)

		err := r.checkConfigFile(cfg)
		if err != nil {
			log.Errorln(err)
		}
	}

	changesRequired := 0

	for _, cfg := range r.configFiles {
		if cfg.ChangeNeeded &&
			!cfg.ChangeApplied &&
			cfg.AuditComplete &&
			!cfg.PreReqFailed &&
			!cfg.AuditFailed {

			changesRequired++
			if cfg.Name == "plugin.config" && r.configFiles["remap.config"].PreReqFailed == true {
				updateStatus = UpdateTropsFailed
				log.Errorln("plugin.config changed however, prereqs failed for remap.config so I am skipping updates for plugin.config")
				continue
			} else if cfg.Name == "remap.config" && r.configFiles["plugin.config"].PreReqFailed == true {
				updateStatus = UpdateTropsFailed
				log.Errorln("remap.config changed however, prereqs failed for plugin.config so I am skipping updates for remap.config")
				continue
			} else {
				log.Debugf("All Prereqs passed for replacing %s on disk with that in Traffic Ops.\n", cfg.Name)
				err := r.replaceCfgFile(cfg)
				if err != nil {
					log.Errorf("failed to replace the config file, '%s',  on disk with data in Traffic Ops.\n", cfg.Name)
				}
			}
		}
	}

	if updateStatus != UpdateTropsFailed && changesRequired > 0 {
		return UpdateTropsNeeded, nil
	}

	return updateStatus, nil
}

// ProcessPackages retrievies a list of required RPM's from Traffic Ops
// and determines which need to be installed or removed on the cache.
func (r *TrafficOpsReq) ProcessPackages() error {
	var pkgs []Package
	var install []string   // install package list.
	var uninstall []string // uninstall package list

	// get the package list for this cache from Traffic Ops.
	out, err := r.atsTcExec("packages")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(out, &pkgs); err != nil {
		return err
	}

	// loop through the package list to build an install and uninstall list.
	for ii := range pkgs {
		var instpkg string // installed package
		var reqpkg string  // required package
		log.Infof("Processing package %s-%s\n", pkgs[ii].Name, pkgs[ii].Version)
		// check to see if any package by name is installed.
		arr, err := util.PackageInfo("pkg-query", pkgs[ii].Name)
		if err != nil {
			return err
		}
		// go needs the ternary operator :)
		if len(arr) == 1 {
			instpkg = arr[0]
		} else {
			instpkg = ""
		}
		// check if the full package version is installed
		fullPackage := pkgs[ii].Name + "-" + pkgs[ii].Version

		if instpkg == fullPackage {
			log.Infof("%s Currently installed and not marked for removal\n", reqpkg)
			r.pkgs[fullPackage] = true
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
				return err
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
	}
	log.Debugf("number of packages requiring installation: %d\n", len(install))
	if r.Cfg.RunMode == config.Report {
		log.Errorf("number of packages requiring installation: %d\n", len(install))
	}
	log.Debugf("number of packages requiring removal: %d\n", len(uninstall))
	if r.Cfg.RunMode == config.Report {
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
		if len(install) > 0 && r.Cfg.RunMode == config.BadAss {
			for jj := range uninstall {
				log.Infof("Uninstalling %s\n", install[jj])
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
					r.pkgs[pkg] = true
					log.Infof("Package %s was installed\n", pkg)
				}
			}
		}
	}
	if r.Cfg.RunMode == config.Report && len(install) > 0 {
		for ii := range install {
			log.Errorf("\nIn Report mode and %s needs installation.\n", install[ii])
			return errors.New("In Report mode and packages need installation")
		}
	}
	return nil
}

func (r *TrafficOpsReq) RevalidateWhileSleeping() (UpdateStatus, error) {
	updateStatus, err := r.CheckRevalidateState(true)
	if err != nil {
		return updateStatus, err
	}
	if updateStatus != 0 {
		r.Cfg.RunMode = config.Revalidate
		err = r.GetConfigFileList()
		if err != nil {
			return updateStatus, err
		}

		updateStatus, err := r.ProcessConfigFiles()
		if err != nil {
			return updateStatus, err
		}

		result := r.StartServices(&updateStatus)
		if !result {
			return updateStatus, errors.New("failed to start services.")
		}

		// update Traffic Ops
		_, err = r.UpdateTrafficOps(&updateStatus)
		if err != nil {
			log.Errorf("failed to update Traffic Ops: %s\n", err.Error())
		}

		r.TrafficCtlReload = false
	}

	return updateStatus, nil
}

func (r *TrafficOpsReq) StartServices(syncdsUpdate *UpdateStatus) bool {
	startSuccess := true

	// start ATS
	if r.IsPackageInstalled("trafficserver") {
		svcStatus, _, err := util.GetServiceStatus("trafficserver")
		if err != nil {
			log.Errorf("error getting 'trafficserver' run status: %s", err)
			startSuccess = false
		} else if r.Cfg.RunMode == config.BadAss {
			if svcStatus == util.SvcRunning {
				running, err := util.ServiceStart("trafficserver", "restart")
				if err != nil {
					log.Errorf("failed to restart trafficserver.")
					startSuccess = false
				} else if running {
					log.Infof("trafficserver has been restarted.")
					if r.TrafficCtlReload {
						log.Infoln("trafficserver was just started, no need to run 'traffic_ctl config reload'")
						r.TrafficCtlReload = false
					}
				}
			} else {
				running, err := util.ServiceStart("trafficserver", "start")
				if err != nil {
					startSuccess = false
					log.Errorf("trafficserver failed to start, running 'traffic_ctl config reload' will also fail: %s\n", err.Error())
				} else if running {
					log.Infoln("trafficserver was successfully started.")
					r.TrafficServerRestart = false
					if r.TrafficCtlReload {
						log.Infoln("trafficserver was just started, no need to run 'traffic_ctl config reload'")
						r.TrafficCtlReload = false
					}
				}
			}
		}

		if svcStatus == util.SvcRunning && r.TrafficCtlReload && !r.TrafficServerRestart {
			switch r.Cfg.RunMode {
			case config.Report:
				log.Errorln("ATS configuration has changed.  'traffic_ctl config reload' needs to be run")
				break
			case config.BadAss:
				fallthrough
			case config.SyncDS:
				fallthrough
			case config.Revalidate:
				log.Infoln("ATS configuration has changed, Running 'traffic_ctl config reload' now.")
				_, _, err := util.ExecCommand(config.TSHome+config.TrafficCtl, "config", "reload")
				if err != nil {
					if *syncdsUpdate == UpdateTropsNeeded {
						*syncdsUpdate = UpdateTropsFailed
						log.Errorf("ATS configuration has change and 'traffic_ctl config reload' has failed, check ATS logs: %s", err.Error())
						startSuccess = false
					}
				} else {
					if *syncdsUpdate == UpdateTropsNeeded {
						log.Infoln("ATS 'traffic_ctl config reload' was successful")
						*syncdsUpdate = UpdateTropsSuccessful
					}
				}
			default:
				log.Errorln("ATS configuration has changed.  'traffic_ctl config reload was not run.")
			}
		} else if r.TrafficCtlReload && (!startSuccess || svcStatus != util.SvcRunning) {
			log.Errorln("ATS configuration has changed.  The new config will be picked up the next time ATS is started.")
			if *syncdsUpdate == UpdateTropsNeeded {
				*syncdsUpdate = UpdateTropsSuccessful
				log.Errorln("'traffic_ctl config reload' was not run but, Traffic Ops is being updated anyway.")
			}
		}
	} else {
		log.Errorln("trafficserver is not installed.")
	}

	return startSuccess
}

func (r *TrafficOpsReq) UpdateTrafficOps(syncdsUpdate *UpdateStatus) (bool, error) {
	var updateResult bool

	serverStatus, err := r.getUpdateStatus()
	if err != nil {
		return false, errors.New("failed to update Traffic Ops: " + err.Error())
	}

	if *syncdsUpdate == UpdateTropsNotNeeded && (serverStatus.UpdatePending == true || serverStatus.RevalPending == true) {
		updateResult = true
		log.Errorln("Traffic Ops is signaling that an update is ready to be applied but, none was found! Clearing update state in Traffic Ops anyway.")
	} else if *syncdsUpdate == UpdateTropsNotNeeded {
		log.Errorln("Traffic Ops does not require an update at this time")
		return true, nil
	} else if *syncdsUpdate == UpdateTropsFailed {
		log.Errorln("Traffic Ops requires an update but, applying the update locally failed.  Traffic Ops is not being updated.")
		return true, nil
	} else if *syncdsUpdate == UpdateTropsSuccessful {
		updateResult = true
		log.Errorln("Traffic Ops requires an update and it was applied successfully.  Clearing update state in Traffic Ops.")
	}

	if updateResult {
		switch r.Cfg.RunMode {
		case config.Report:
			log.Errorln("In Report mode and Traffic Ops needs updated you should probably do that manually.")
			break
		case config.BadAss:
			fallthrough
		case config.SyncDS:
			if serverStatus.RevalPending {
				_, err = r.atsTcExecCommand("send-update", 0, 1)
			} else {
				_, err = r.atsTcExecCommand("send-update", 0, 0)
			}
		case config.Revalidate:
			_, err = r.atsTcExecCommand("send-update", 1, 0)
		}
		if err != nil {
			return false, errors.New("Traffic Ops Update failed: " + err.Error())
		} else {
			log.Errorln("Traffic Ops has been updated.")
		}
	}
	return true, nil
}
