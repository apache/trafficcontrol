package config

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
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-apply"

var TSHome string = "/opt/trafficserver"

const DefaultTSConfigDir = "/opt/trafficserver/etc/trafficserver"

const (
	StatusDir          = "/var/lib/trafficcontrol-cache-config/status"
	Chkconfig          = "/sbin/chkconfig"
	Service            = "/sbin/service"
	SystemCtl          = "/bin/systemctl"
	TmpBase            = "/tmp/trafficcontrol-cache-config"
	TrafficCtl         = "/bin/traffic_ctl"
	TrafficServerOwner = "ats"
)

type SvcManagement int

const (
	Unknown SvcManagement = 0
	SystemD SvcManagement = 1
	SystemV SvcManagement = 2 // legacy System V Init.
)

func (s SvcManagement) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case SystemD:
		return "SystemD"
	case SystemV:
		return "SystemV"
	}
	return "Unknown"
}

type Cfg struct {
	LogLocationDebug    string
	LogLocationErr      string
	LogLocationInfo     string
	LogLocationWarn     string
	CacheHostName       string
	SvcManagement       SvcManagement
	Retries             int
	ReverseProxyDisable bool
	RpmDBOk             bool
	SkipOSCheck         bool
	UseStrategies       t3cutil.UseStrategiesFlag
	TOInsecure          bool
	TOTimeoutMS         time.Duration
	TOUser              string
	TOPass              string
	TOURL               string
	DNSLocalBind        bool
	WaitForParents      bool
	YumOptions          string
	// UseGit is whether to create and maintain a git repo of config changes.
	// Note this only applies to the ATS config directory inferred or set via the flag.
	//      It does not do anything for config files generated outside that location.
	UseGit                      UseGitFlag
	NoCache                     bool
	SyncDSUpdatesIPAllow        bool
	OmitViaStringRelease        bool
	NoOutgoingIP                bool
	DisableParentConfigComments bool
	DefaultClientEnableH2       *bool
	DefaultClientTLSVersions    *string
	// MaxMindLocation is a URL string for a download location for a maxmind database
	// for use with either HeaderRewrite or Maxmind_ACL plugins
	MaxMindLocation string
	TsHome          string
	TsConfigDir     string

	ServiceAction          t3cutil.ApplyServiceActionFlag
	NoConfirmServiceAction bool

	ReportOnly        bool
	GoDirect          string
	Files             t3cutil.ApplyFilesFlag
	InstallPackages   bool
	IgnoreUpdateFlag  bool
	NoUnsetUpdateFlag bool
	UpdateIPAllow     bool
	Version           string
	GitRevision       string
	LocalATSVersion   string
	CacheType         string
}

func (cfg Cfg) AppVersion() string { return t3cutil.VersionStr(AppName, cfg.Version, cfg.GitRevision) }

type UseGitFlag string

const (
	UseGitAuto    = "auto"
	UseGitYes     = "yes"
	UseGitNo      = "no"
	UseGitInvalid = ""
)

func StrToUseGitFlag(str string) UseGitFlag {
	str = strings.ToLower(strings.TrimSpace(str))
	switch str {
	case UseGitAuto:
		fallthrough
	case UseGitYes:
		fallthrough
	case UseGitNo:
		return UseGitFlag(str)
	default:
		return UseGitInvalid
	}
}

type WaitForParentsFlag string

const WaitForParentsDefault = WaitForParentsReval
const TimeAndDateLayout = "Jan 2, 2006 15:04 MST"

const (
	WaitForParentsTrue    = "true"
	WaitForParentsFalse   = "false"
	WaitForParentsReval   = "reval"
	WaitForParentsInvalid = ""
)

func StrToWaitForParentsFlag(str string) WaitForParentsFlag {
	str = strings.ToLower(strings.TrimSpace(str))
	switch str {
	case WaitForParentsTrue:
		fallthrough
	case WaitForParentsFalse:
		fallthrough
	case WaitForParentsReval:
		return WaitForParentsFlag(str)
	default:
		return WaitForParentsInvalid
	}
}

func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationErr) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationDebug) }
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // event logging is not used.

func directoryExists(dir string) (bool, os.FileInfo) {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	}
	return info.IsDir(), info
}

// derives the ATS Installation directory from
// the rpm config file list.
func GetTSPackageHome() string {
	var dir []string
	var output bytes.Buffer
	var tsHome string = ""
	var files []string

	cmd := exec.Command("/bin/rpm", "-q", "-c", "trafficserver")
	cmd.Stdout = &output
	err := cmd.Run()
	// on error or if the trafficserver rpm is not installed indicated
	// by a return code of '1', return an empty string.
	if err != nil || cmd.ProcessState.ExitCode() == 1 {
		return ""
	} else {
		responseStr := string(output.Bytes())
		files = strings.Split(responseStr, "\n")
	}

	if files != nil && err == nil { // trafficserver is installed, derive TSHome
		for ii := range files {
			line := strings.TrimSpace(files[ii])
			if strings.Contains(line, "etc/trafficserver") {
				dir = strings.Split(line, "/etc")
				break
			}
		}
		if dir != nil && len(dir) >= 0 {
			if ok, _ := directoryExists(dir[0]); ok == true {
				tsHome = dir[0]
			}
		}
	}
	return tsHome
}

const (
	rpmDBBdb             = "bdb"
	rpmDBSquLite         = "sqlite"
	rpmDBUnknown         = "unknown"
	rpmDBVerifyCmd       = "/usr/lib/rpm/rpmdb_verify"
	sqliteRpmDbVerifyCmd = "/bin/sqlite3"
	rpmDir               = "/var/lib/rpm"
	sqliteRpmDB          = "rpmdb.sqlite"
)

// getRpmDBBackend uses "%_db_backend" macro to get the database type
func getRpmDBBackend() (string, error) {
	var outBuf bytes.Buffer
	cmd := exec.Command("/bin/rpm", "-E", "%_db_backend")
	cmd.Stdout = &outBuf
	err := cmd.Run()
	if err != nil {
		return rpmDBUnknown, err
	}
	return strings.TrimSpace(outBuf.String()), nil
}

// isSqliteInstalled looks to see if the sqlite3 executable
// is installed which is needed to do the db verify
func isSqliteInstalled() bool {
	sqliteUtil := isCommandAvailable("/bin/sqlite3")
	return sqliteUtil
}

// verifies the rpm database files. if there is any database corruption
// it will return false
func verifyRpmDB(rpmDir string) bool {
	exclude := regexp.MustCompile(`(^\.|^__)`)
	dbFiles, err := os.ReadDir(rpmDir)
	if err != nil {
		return false
	}
	for _, file := range dbFiles {
		if exclude.Match([]byte(file.Name())) {
			continue
		}
		cmd := exec.Command(rpmDBVerifyCmd, rpmDir+"/"+file.Name())
		err := cmd.Run()
		if err != nil || cmd.ProcessState.ExitCode() > 0 {
			return false
		}
	}
	return true
}

// verifySqliteRpmDB runs PRAGMA quick_check
// requires /bin/sqlite3
func verifySqliteRpmDB(sqliteDB string) bool {
	args := []string{sqliteDB, `PRAGMA quick_check`}
	cmd := exec.Command(sqliteRpmDbVerifyCmd, args...)
	err := cmd.Run()
	if err != nil || cmd.ProcessState.ExitCode() > 0 {
		return false
	}
	return true
}

func GetCfg(appVersion string, gitRevision string) (Cfg, error) {
	var err error
	toInfoLog := []string{}

	version := getopt.BoolLong("version", 'E', "Print version information and exit.")
	cacheHostNamePtr := getopt.StringLong("cache-host-name", 'H', "", "Host name of the cache to generate config for. Must be the server host name in Traffic Ops, not a URL, and not the FQDN")
	retriesPtr := getopt.IntLong("num-retries", 'r', 3, "[number] retry connection to Traffic Ops URL [number] times, default is 3")
	reverseProxyDisablePtr := getopt.BoolLong("reverse-proxy-disable", 'p', "[false | true] bypass the reverse proxy even if one has been configured default is false")
	skipOSCheckPtr := getopt.BoolLong("skip-os-check", 'C', "[false | true] skip os check, default is false")
	toInsecurePtr := getopt.BoolLong("traffic-ops-insecure", 'I', "[true | false] ignore certificate errors from Traffic Ops")
	toTimeoutMSPtr := getopt.IntLong("traffic-ops-timeout-milliseconds", 't', 30000, "Timeout in milli-seconds for Traffic Ops requests, default is 30000")
	toURLPtr := getopt.StringLong("traffic-ops-url", 'u', "", "Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with the environment variable TO_URL")
	toUserPtr := getopt.StringLong("traffic-ops-user", 'U', "", "Traffic Ops username. Required. May also be set with the environment variable TO_USER")
	toPassPtr := getopt.StringLong("traffic-ops-password", 'P', "", "Traffic Ops password. Required. May also be set with the environment variable TO_PASS")
	tsHomePtr := getopt.StringLong("trafficserver-home", 'R', "", "Trafficserver Package directory. May also be set with the environment variable TS_HOME")
	dnsLocalBindPtr := getopt.BoolLong("dns-local-bind", 'b', "[true | false] whether to use the server's Service Addresses to set the ATS DNS local bind address")
	help := getopt.BoolLong("help", 'h', "Print usage information and exit")
	useGitStr := getopt.StringLong("git", 'g', "auto", "Create and use a git repo in the config directory. Options are yes, no, and auto. If yes, create and use. If auto, use if it exist. Default is auto.")
	noCachePtr := getopt.BoolLong("no-cache", 'n', "Whether to not use a cache and make conditional requests to Traffic Ops")
	syncdsUpdatesIPAllowPtr := getopt.BoolLong("syncds-updates-ipallow", 'S', "Whether syncds mode will update ipallow. This exists because ATS had a bug where reloading after changing ipallow would block everything. Default is false.")
	omitViaStringReleasePtr := getopt.BoolLong("omit-via-string-release", 'e', "Whether to set the records.config via header to the ATS release from the RPM. Default true.")
	noOutgoingIP := getopt.BoolLong("no-outgoing-ip", 'i', "Whether to not set the records.config outgoing IP to the server's addresses in Traffic Ops. Default is false.")
	disableParentConfigCommentsPtr := getopt.BoolLong("disable-parent-config-comments", 'c', "Whether to disable verbose parent.config comments. Default false.")
	defaultEnableH2 := getopt.BoolLong("default-client-enable-h2", '2', "Whether to enable HTTP/2 on Delivery Services by default, if they have no explicit Parameter. This is irrelevant if ATS records.config is not serving H2. If omitted, H2 is disabled.")
	defaultClientTLSVersions := getopt.StringLong("default-client-tls-versions", 'V', "", "Comma-delimited list of default TLS versions for Delivery Services with no Parameter, e.g. --default-tls-versions='1.1,1.2,1.3'. If omitted, all versions are enabled.")
	maxmindLocationPtr := getopt.StringLong("maxmind-location", 'M', "", "URL of a maxmind gzipped database file, to be installed into the trafficserver etc directory.")
	verbosePtr := getopt.CounterLong("verbose", 'v', `Log verbosity. Logging is output to stderr. By default, errors are logged. To log warnings, pass '-v'. To log info, pass '-vv'. To omit error logging, see '-s'`)
	cache := getopt.StringLong("cache", 'T', "ats", "Cache server type. Generate configuration files for specific cache server type, e.g. 'ats', 'varnish'.")
	const silentFlagName = "silent"
	silentPtr := getopt.BoolLong(silentFlagName, 's', `Silent. Errors are not logged, and the 'verbose' flag is ignored. If a fatal error occurs, the return code will be non-zero but no text will be output to stderr`)

	const goDirectFlagName = "go-direct"
	goDirectPtr := getopt.StringLong(goDirectFlagName, 'G', "old", "[true|false|old] default will set go_direct to old and it will be based on opposite of parent_is_proxy directivefalse, you can also set go_direct true, or false.")

	const waitForParentsFlagName = "wait-for-parents"
	waitForParentsPtr := getopt.BoolLong(waitForParentsFlagName, 'W', "[true | false] do not update if parent_pending = 1 in the update json. Default is false")

	const serviceActionFlagName = "service-action"
	const defaultServiceAction = t3cutil.ApplyServiceActionFlagReload
	serviceActionPtr := getopt.EnumLong(serviceActionFlagName, 'a', []string{string(t3cutil.ApplyServiceActionFlagReload), string(t3cutil.ApplyServiceActionFlagRestart), string(t3cutil.ApplyServiceActionFlagNone), ""}, "", "action to perform on Traffic Server and other system services. Only reloads if necessary, but always restarts. Default is 'reload'")

	const noConfirmServiceActionFlagName = "no-confirm-service-action"
	noConfirmServiceAction := getopt.BoolLong(noConfirmServiceActionFlagName, 0, "Whether to skip waiting and confirming the service action succeeded (reload or restart) via t3c-tail. Default is false.")

	const reportOnlyFlagName = "report-only"
	reportOnlyPtr := getopt.BoolLong(reportOnlyFlagName, 'o', "Log information about necessary files and actions, but take no action. Default is false")

	const filesFlagName = "files"
	const defaultFiles = t3cutil.ApplyFilesFlagAll
	filesPtr := getopt.EnumLong(filesFlagName, 'f', []string{string(t3cutil.ApplyFilesFlagAll), string(t3cutil.ApplyFilesFlagReval), ""}, "", "[all | reval] Which files to generate. If reval, the Traffic Ops server reval_pending flag is used instead of the upd_pending flag. Default is 'all'")

	const installPackagesFlagName = "install-packages"
	installPackagesPtr := getopt.BoolLong(installPackagesFlagName, 'k', "Whether to install necessary packages. Default is false.")

	const ignoreUpdateFlagName = "ignore-update-flag"
	ignoreUpdateFlagPtr := getopt.BoolLong(ignoreUpdateFlagName, 'F', "Whether to ignore the upd_pending or reval_pending flag in Traffic Ops, and always generate and apply files. If true, the flag is still unset in Traffic Ops after files are applied. Default is false.")
	noUnsetUpdateFlagPtr := getopt.BoolLong("no-unset-update-flag", 'd', "Whether to not unset the update flag in Traffic Ops after applying files. This option makes it possible to generate test or debug configuration from a production Traffic Ops without un-setting queue or reval flags. Default is false.")

	const updateIPAllowFlagName = "update-ipallow"
	updateIPAllowPtr := getopt.BoolLong(updateIPAllowFlagName, 'A', "Whether ipallow file will be updated if necessary. This exists because ATS had a bug where reloading after changing ipallow would block everything. Default is false.")

	const useStrategiesFlagName = "use-strategies"
	const defaultUseStrategies = t3cutil.UseStrategiesFlagFalse
	useStrategiesPtr := getopt.EnumLong(useStrategiesFlagName, 0, []string{string(t3cutil.UseStrategiesFlagTrue), string(t3cutil.UseStrategiesFlagCore), string(t3cutil.UseStrategiesFlagCore), ""}, "", "[true | core| false] whether to generate config using strategies.yaml instead of parent.config. If true use the parent_select plugin, if 'core' use ATS core strategies, if false use parent.config.")

	const useLocalATSVersionFlagName = "local-ats-version"
	useLocalATSVersionPtr := getopt.BoolLong(useLocalATSVersionFlagName, 0, "[true | false] whether to use the local installed ATS version for config generation. If false, attempt to use the Server Package Parameter and fall back to ATS 5. If true and the local ATS version cannot be found, an error will be logged and the version set to ATS 5. Default is false")

	const runModeFlagName = "run-mode"
	runModePtr := getopt.StringLong(runModeFlagName, 'm', "", `[badass | report | revalidate | syncds] run mode. Optional, convenience flag which sets other flags for common usage scenarios.
syncds     keeps the defaults:
                --report-only=false
                --files=all
                --install-packages=false
                --service-action=reload
                --ignore-update-flag=false
                --update-ipallow=false
                --no-unset-update-flag=false
revalidate sets --files=reval
                --wait-for-parents=true
badass     sets --install-packages=true
                --service-action=restart
                --ignore-update-flag=true
                --update-ipallow=true
report     sets --report-only=true
                --no-unset-update-flag=true
                --silent

Note the 'syncds' settings are all the flag defaults. Hence, if no mode is set, the default is effectively 'syncds'.

If any of the related flags are also set, they override the mode's default behavior.`)

	getopt.Parse()

	// The mode is never exposed outside this function to prevent accidentally changing behavior based on it,
	// so we want to log what flags the mode set here, to aid debugging.
	// But we can't do that until the loggers are initialized.
	modeLogStrs := []string{}
	fatalLogStrs := []string{}
	if getopt.IsSet(runModeFlagName) {
		runMode := t3cutil.StrToMode(*runModePtr)
		if runMode == t3cutil.ModeInvalid {
			fatalLogStrs = append(fatalLogStrs, *runModePtr+" is an invalid mode.")
		}
		modeLogStrs = append(modeLogStrs, "t3c-apply is running in "+runMode.String()+" mode")
		switch runMode {
		case t3cutil.ModeSyncDS:
			// syncds flags are all the defaults, no need to change anything
		case t3cutil.ModeRevalidate:
			if !getopt.IsSet(filesFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+filesFlagName+"="+t3cutil.ApplyFilesFlagReval.String())
				*filesPtr = t3cutil.ApplyFilesFlagReval.String()
			}
			if !getopt.IsSet(waitForParentsFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+waitForParentsFlagName+"="+"true")
				*waitForParentsPtr = true
			}
		case t3cutil.ModeBadAss:
			if !getopt.IsSet(serviceActionFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+serviceActionFlagName+"="+t3cutil.ApplyServiceActionFlagRestart.String())
				*serviceActionPtr = t3cutil.ApplyServiceActionFlagRestart.String()
			}
			if !getopt.IsSet(installPackagesFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+installPackagesFlagName+"="+"true")
				*installPackagesPtr = true
			}
			if !getopt.IsSet(ignoreUpdateFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+ignoreUpdateFlagName+"="+"true")
				*ignoreUpdateFlagPtr = true
			}
			if !getopt.IsSet(updateIPAllowFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+updateIPAllowFlagName+"="+"true")
				*updateIPAllowPtr = true
			}
		case t3cutil.ModeReport:
			if !getopt.IsSet(reportOnlyFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+reportOnlyFlagName+"="+"true")
				*reportOnlyPtr = true
			}
			if !getopt.IsSet(ignoreUpdateFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+ignoreUpdateFlagName+"="+"true")
				*ignoreUpdateFlagPtr = true
			}
			if !getopt.IsSet(silentFlagName) {
				modeLogStrs = append(modeLogStrs, runMode.String()+" setting --"+silentFlagName+"="+"true")
				*silentPtr = true
			}
		}
	}

	switch *goDirectPtr {
	case "false", "true", "old":
		if !getopt.IsSet(goDirectFlagName) {
			modeLogStrs = append(modeLogStrs, goDirectFlagName+" not set using default 'old'")
		}
	default:
		modeLogStrs = append(modeLogStrs, *goDirectPtr+" is not a valid go-direct option setting default 'old'")
		*goDirectPtr = "old"
	}

	if *serviceActionPtr == "" {
		*serviceActionPtr = defaultServiceAction.String()
	}
	if *filesPtr == "" {
		*filesPtr = defaultFiles.String()
	}

	if !getopt.IsSet(useStrategiesFlagName) {
		*useStrategiesPtr = defaultUseStrategies.String()
	}

	logLocationError := log.LogLocationStderr
	logLocationWarn := log.LogLocationNull
	logLocationInfo := log.LogLocationNull
	logLocationDebug := log.LogLocationNull
	if *silentPtr {
		logLocationError = log.LogLocationNull
	} else {
		if *verbosePtr >= 1 {
			logLocationWarn = log.LogLocationStderr
		}
		if *verbosePtr >= 2 {
			logLocationInfo = log.LogLocationStderr
			logLocationDebug = log.LogLocationStderr // t3c only has 3 verbosity options: none (-s), error (default or --verbose=0), warning (-v), and info (-vv). Any code calling log.Debug is treated as Info.
		}
	}

	if *verbosePtr > 2 {
		fatalLogStrs = append(fatalLogStrs, "Too many verbose options. The maximum log verbosity level is 2 (-vv or --verbose=2) for errors (0), warnings (1), and info (2)")
	}

	var cacheHostName string
	if len(*cacheHostNamePtr) > 0 {
		cacheHostName = *cacheHostNamePtr
	} else {
		cacheHostName, err = os.Hostname()
		if err != nil {
			fatalLogStrs = append(fatalLogStrs, "Could not get the hostname from the O.S., please supply a hostname: "+err.Error())
		}
		// strings.Split always returns a slice with at least 1 element, so we don't need a len check
		cacheHostName = strings.Split(cacheHostName, ".")[0]
	}

	useGit := StrToUseGitFlag(*useGitStr)

	if useGit == UseGitInvalid {
		fatalLogStrs = append(fatalLogStrs, "Invalid git flag '"+*useGitStr+"'. Valid options are yes, no, auto.")
	}

	retries := *retriesPtr
	reverseProxyDisable := *reverseProxyDisablePtr
	skipOsCheck := *skipOSCheckPtr
	useStrategies := t3cutil.UseStrategiesFlag(*useStrategiesPtr)
	toInsecure := *toInsecurePtr
	toTimeoutMS := time.Millisecond * time.Duration(*toTimeoutMSPtr)
	toURL := *toURLPtr
	toUser := *toUserPtr
	toPass := *toPassPtr
	dnsLocalBind := *dnsLocalBindPtr
	maxmindLocation := *maxmindLocationPtr

	if *version {
		cfg := &Cfg{Version: appVersion, GitRevision: gitRevision}
		fmt.Println(cfg.AppVersion())
		os.Exit(0)
	} else if *help {
		Usage()
		return Cfg{}, nil
	}

	urlSourceStr := "argument" // for error messages
	if toURL == "" {
		urlSourceStr = "environment variable"
		toURL = os.Getenv("TO_URL")
	} else {
		os.Setenv("TO_URL", toURL)
	}
	if toUser == "" {
		toUser = os.Getenv("TO_USER")
	} else {
		os.Setenv("TO_USER", toUser)
	}
	if *toPassPtr == "" {
		toPass = os.Getenv("TO_PASS")
	} else {
		os.Setenv("TO_PASS", toPass)
	}

	rpmDBisOk := true
	rpmDBType, err := getRpmDBBackend()
	if err != nil {
		toInfoLog = append(toInfoLog, fmt.Sprintf("error getting db type: %s", err.Error()))
		rpmDBType = rpmDBUnknown
	}

	if rpmDBType == rpmDBSquLite {
		sqliteUtil := isSqliteInstalled()
		toInfoLog = append(toInfoLog, fmt.Sprintf("RPM database is %s", rpmDBSquLite))
		if sqliteUtil {
			rpmDBisOk = verifySqliteRpmDB(rpmDir + "/" + sqliteRpmDB)
			toInfoLog = append(toInfoLog, fmt.Sprintf("RPM database is ok: %t", rpmDBisOk))
		} else {
			toInfoLog = append(toInfoLog, "/bin/sqlite3 not available, RPM database not checked")
		}
	} else if rpmDBType == rpmDBBdb {
		toInfoLog = append(toInfoLog, fmt.Sprintf("RPM database is %s", rpmDBBdb))
		rpmDBisOk = verifyRpmDB(rpmDir)
		toInfoLog = append(toInfoLog, fmt.Sprintf("RPM database is ok: %t", rpmDBisOk))
	} else {
		toInfoLog = append(toInfoLog, fmt.Sprintf("RPM DB type is %s DB check will be skipped", rpmDBUnknown))
	}

	if *installPackagesPtr && !rpmDBisOk {
		if t3cutil.StrToMode(*runModePtr) == t3cutil.ModeBadAss {
			fatalLogStrs = append(fatalLogStrs, "RPM database check failed unable to install packages cannot continue in badass mode")
		} else {
			fatalLogStrs = append(fatalLogStrs, "RPM database check failed unable to install packages cannot continue")
		}
	}

	// set TSHome
	var tsHome = ""
	if *tsHomePtr != "" {
		tsHome = *tsHomePtr
		toInfoLog = append(toInfoLog, fmt.Sprintf("set TSHome from command line: '%s'", TSHome))
	}
	if *tsHomePtr == "" { // evironment or rpm check.
		tsHome = os.Getenv("TS_HOME") // check for the environment variable.
		if tsHome != "" {
			toInfoLog = append(toInfoLog, fmt.Sprintf("set TSHome from TS_HOME environment variable '%s'\n", TSHome))
		} else if rpmDBisOk { // check using the config file listing from the rpm package if rpmdb is ok.
			tsHome = GetTSPackageHome()
			if tsHome != "" {
				toInfoLog = append(toInfoLog, fmt.Sprintf("set TSHome from the RPM config file  list '%s'\n", TSHome))
			}
		} else if tsHome == "" {
			toInfoLog = append(toInfoLog, fmt.Sprintf("no override for TSHome was found, using the configured default: '%s'\n", TSHome))
		}
	}

	tsConfigDir := DefaultTSConfigDir

	if tsHome != "" {
		TSHome = tsHome
		tsConfigDir = tsHome + "/etc/trafficserver"
		if cache != nil && *cache == "varnish" {
			tsConfigDir = tsHome + "/etc/varnish"
		}
		toInfoLog = append(toInfoLog, fmt.Sprintf("TSHome: %s, TSConfigDir: %s\n", TSHome, tsConfigDir))
	}

	atsVersionStr := ""
	if *useLocalATSVersionPtr {
		atsVersionStr, err = GetATSVersionStr(tsHome)
		if err != nil {
			fatalLogStrs = append(fatalLogStrs, "getting local ATS version: "+err.Error())
		}
	}
	toInfoLog = append(toInfoLog, fmt.Sprintf("ATSVersionStr: '%s'\n", atsVersionStr))

	usageStr := "basic usage: t3c-apply --traffic-ops-url=myurl --traffic-ops-user=myuser --traffic-ops-password=mypass --cache-host-name=my-cache"
	if strings.TrimSpace(toURL) == "" {
		fatalLogStrs = append(fatalLogStrs, "Missing required argument --traffic-ops-url or TO_URL environment variable. "+usageStr)
	}
	if strings.TrimSpace(toUser) == "" {
		fatalLogStrs = append(fatalLogStrs, "Missing required argument --traffic-ops-user or TO_USER environment variable. "+usageStr)
	}
	if strings.TrimSpace(toPass) == "" {
		fatalLogStrs = append(fatalLogStrs, "Missing required argument --traffic-ops-password or TO_PASS environment variable. "+usageStr)
	}
	if strings.TrimSpace(cacheHostName) == "" {
		fatalLogStrs = append(fatalLogStrs, "Missing required argument --cache-host-name. "+usageStr)
	}

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return Cfg{}, errors.New("parsing Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	} else if err = validateURL(toURLParsed); err != nil {
		return Cfg{}, errors.New("invalid Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	}

	svcManagement := getOSSvcManagement()
	yumOptions := os.Getenv("YUM_OPTIONS")

	cfg := Cfg{
		LogLocationDebug:            logLocationDebug,
		LogLocationErr:              logLocationError,
		LogLocationInfo:             logLocationInfo,
		LogLocationWarn:             logLocationWarn,
		CacheHostName:               cacheHostName,
		SvcManagement:               svcManagement,
		Retries:                     retries,
		RpmDBOk:                     rpmDBisOk,
		ReverseProxyDisable:         reverseProxyDisable,
		SkipOSCheck:                 skipOsCheck,
		UseStrategies:               useStrategies,
		TOInsecure:                  toInsecure,
		TOTimeoutMS:                 toTimeoutMS,
		TOUser:                      toUser,
		TOPass:                      toPass,
		TOURL:                       toURL,
		DNSLocalBind:                dnsLocalBind,
		WaitForParents:              *waitForParentsPtr,
		YumOptions:                  yumOptions,
		UseGit:                      useGit,
		NoCache:                     *noCachePtr,
		SyncDSUpdatesIPAllow:        *syncdsUpdatesIPAllowPtr,
		UpdateIPAllow:               *updateIPAllowPtr,
		OmitViaStringRelease:        *omitViaStringReleasePtr,
		NoOutgoingIP:                *noOutgoingIP,
		DisableParentConfigComments: *disableParentConfigCommentsPtr,
		DefaultClientEnableH2:       defaultEnableH2,
		DefaultClientTLSVersions:    defaultClientTLSVersions,
		MaxMindLocation:             maxmindLocation,
		TsHome:                      TSHome,
		TsConfigDir:                 tsConfigDir,
		GoDirect:                    *goDirectPtr,
		ServiceAction:               t3cutil.ApplyServiceActionFlag(*serviceActionPtr),
		NoConfirmServiceAction:      *noConfirmServiceAction,
		ReportOnly:                  *reportOnlyPtr,
		Files:                       t3cutil.ApplyFilesFlag(*filesPtr),
		InstallPackages:             *installPackagesPtr,
		IgnoreUpdateFlag:            *ignoreUpdateFlagPtr,
		NoUnsetUpdateFlag:           *noUnsetUpdateFlagPtr,
		Version:                     appVersion,
		GitRevision:                 gitRevision,
		LocalATSVersion:             atsVersionStr,
		CacheType:                   *cache,
	}

	if err = log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("Initializing loggers: " + err.Error() + "\n")
	}

	if len(fatalLogStrs) > 0 {
		for _, str := range fatalLogStrs {
			str = strings.TrimSpace(str)
			log.Errorln(str)
		}
		return Cfg{}, errors.New("fatal error has occurred")
	}
	for _, str := range modeLogStrs {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		log.Infoln(str)
	}
	for _, msg := range toInfoLog {
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		log.Infoln(msg)
	}

	printConfig(cfg)

	return cfg, nil
}

func validateURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme expected 'http' or 'https', actual '" + u.Scheme + "'")
	}
	if strings.TrimSpace(u.Host) == "" {
		return errors.New("no host")
	}
	return nil
}

func isCommandAvailable(name string) bool {
	command := name + " --version"
	cmd := exec.Command("/bin/sh", "-c", command)
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func getOSSvcManagement() SvcManagement {
	var _svcManager SvcManagement

	if isCommandAvailable(SystemCtl) {
		_svcManager = SystemD
	} else if isCommandAvailable(Service) {
		_svcManager = SystemV
	}
	if !isCommandAvailable(Chkconfig) {
		return Unknown
	}

	// we have what we need
	return _svcManager
}

func GetATSVersionStr(tsHome string) (string, error) {
	tsPath := tsHome
	tsPath = filepath.Join(tsPath, "bin")
	tsPath = filepath.Join(tsPath, "traffic_server")

	stdOut, stdErr, code := t3cutil.Do(`sh`, `-c`, `set -o pipefail && `+tsPath+` --version | head -1 | awk '{print $3}'`)
	if code != 0 {
		return "", fmt.Errorf("traffic_server --version returned code %v stderr '%v' stdout '%v'", code, string(stdErr), string(stdOut))
	}
	atsVersion := strings.TrimSpace(string(stdOut))
	if atsVersion == "" {
		return "", fmt.Errorf("traffic_server --version returned nothing, code %v stderr '%v' stdout '%v'", code, string(stdErr), string(stdOut))
	}
	return atsVersion, nil
}

func printConfig(cfg Cfg) {
	// TODO add new flags
	log.Debugf("LogLocationDebug: %s\n", cfg.LogLocationDebug)
	log.Debugf("LogLocationErr: %s\n", cfg.LogLocationErr)
	log.Debugf("LogLocationInfo: %s\n", cfg.LogLocationInfo)
	log.Debugf("LogLocationWarn: %s\n", cfg.LogLocationWarn)
	log.Debugf("CacheHostName: %s\n", cfg.CacheHostName)
	log.Debugf("SvcManagement: %s\n", cfg.SvcManagement)
	log.Debugf("GoDirect: %s\n", cfg.GoDirect)
	log.Debugf("Retries: %d\n", cfg.Retries)
	log.Debugf("ReverseProxyDisable: %t\n", cfg.ReverseProxyDisable)
	log.Debugf("SkipOSCheck: %t\n", cfg.SkipOSCheck)
	log.Debugf("TOInsecure: %t\n", cfg.TOInsecure)
	log.Debugf("TOTimeoutMS: %d\n", cfg.TOTimeoutMS)
	log.Debugf("TOUser: %s\n", cfg.TOUser)
	log.Debugf("TOPass: Pass len: '%d'\n", len(cfg.TOPass))
	log.Debugf("TOURL: %s\n", cfg.TOURL)
	log.Debugf("TSHome: %s\n", TSHome)
	log.Debugf("LocalATSVersion: %s\n", cfg.LocalATSVersion)
	log.Debugf("WaitForParents: %v\n", cfg.WaitForParents)
	log.Debugf("ServiceAction: %v\n", cfg.ServiceAction)
	log.Debugf("NoConfirmServiceAction: %v\n", cfg.NoConfirmServiceAction)
	log.Debugf("YumOptions: %s\n", cfg.YumOptions)
	log.Debugf("MaxmindLocation: %s\n", cfg.MaxMindLocation)
}

func Usage() {
	getopt.PrintUsage(os.Stdout)
}
