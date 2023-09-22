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
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-generate"

const ExitCodeSuccess = 0
const ExitCodeErrGeneric = 1
const ExitCodeNotFound = 104
const ExitCodeBadRequest = 100

var ErrNotFound = errors.New("not found")
var ErrBadRequest = errors.New("bad request")

type Cfg struct {
	ListPlugins        bool
	LogLocationErr     string
	LogLocationInfo    string
	LogLocationDebug   string
	LogLocationWarn    string
	RevalOnly          bool
	Dir                string
	UseStrategies      t3cutil.UseStrategiesFlag
	GoDirect           string
	ViaRelease         bool
	SetDNSLocalBind    bool
	NoOutgoingIP       bool
	ATSMajorVersion    uint
	ParentComments     bool
	DefaultEnableH2    bool
	DefaultTLSVersions []atscfg.TLSVersion
	Version            string
	GitRevision        string
	Cache              string
}

func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationErr) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationDebug) }
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // app doesn't use the event logger.

func (cfg Cfg) AppVersion() string { return t3cutil.VersionStr(AppName, cfg.Version, cfg.GitRevision) }

// GetCfg gets the application configuration, from arguments and environment variables.
func GetCfg(appVersion string, gitRevision string) (Cfg, error) {
	version := getopt.BoolLong("version", 'V', "Print version information and exit.")
	listPlugins := getopt.BoolLong("list-plugins", 'l', "Print the list of plugins.")
	help := getopt.BoolLong("help", 'h', "Print usage information and exit")
	revalOnly := getopt.BoolLong("revalidate-only", 'y', "Whether to exclude files not named 'regex_revalidate.config'")
	dir := getopt.StringLong("dir", 'D', "", "ATS config directory, used for config files without location parameters or with relative paths. May be blank. If blank and any required config file location parameter is missing or relative, will error.")
	viaRelease := getopt.BoolLong("via-string-release", 'r', "Whether to use the Release value from the RPM package as a replacement for the ATS version specified in the build that is returned in the Via and Server headers from ATS.")
	dnsLocalBind := getopt.BoolLong("dns-local-bind", 'b', "Whether to use the server's Service Addresses to set the ATS DNS local bind address.")
	disableParentConfigComments := getopt.BoolLong("disable-parent-config-comments", 'c', "Disable adding a comments to parent.config individual lines")
	defaultEnableH2 := getopt.BoolLong("default-client-enable-h2", '2', "Whether to enable HTTP/2 on Delivery Services by default, if they have no explicit Parameter. This is irrelevant if ATS records.config is not serving H2. If omitted, H2 is disabled.")
	defaultTLSVersionsStr := getopt.StringLong("default-client-tls-versions", 'T', "", "Comma-delimited list of default TLS versions for Delivery Services with no Parameter, e.g. '--default-tls-versions=1.1,1.2,1.3'. If omitted, all versions are enabled.")
	noOutgoingIP := getopt.BoolLong("no-outgoing-ip", 'i', "Whether to not set the records.config outgoing IP to the server's addresses in Traffic Ops. Default is false.")
	atsVersion := getopt.StringLong("ats-version", 'a', "", "The ATS version, e.g. 9.1.2-42.abc123.el7.x86_64. If omitted, generation will attempt to get the ATS version from the Server Parameters, and fall back to lib/go-atscfg.DefaultATSVersion")
	verbosePtr := getopt.CounterLong("verbose", 'v', `Log verbosity. Logging is output to stderr. By default, errors are logged. To log warnings, pass '-v'. To log info, pass '-vv'. To omit error logging, see '-s'`)
	silentPtr := getopt.BoolLong("silent", 's', `Silent. Errors are not logged, and the 'verbose' flag is ignored. If a fatal error occurs, the return code will be non-zero but no text will be output to stderr`)
	cache := getopt.StringLong("cache", 'C', "ats", "Cache server type. Generate configuration files for specific cache server type, e.g. 'ats', 'varnish'.")

	const useStrategiesFlagName = "use-strategies"
	const defaultUseStrategies = t3cutil.UseStrategiesFlagFalse
	useStrategiesPtr := getopt.EnumLong(useStrategiesFlagName, 0, []string{string(t3cutil.UseStrategiesFlagTrue), string(t3cutil.UseStrategiesFlagCore), string(t3cutil.UseStrategiesFlagFalse), string(t3cutil.UseStrategiesFlagCore), ""}, "", "[true | core| false] whether to generate config using strategies.yaml instead of parent.config. If true use the parent_select plugin, if 'core' use ATS core strategies, if false use parent.config.")

	const goDirectFlagName = "go-direct"
	goDirectPtr := getopt.StringLong(goDirectFlagName, 'G', "false", "[true|false|old] default will set go_direct to false, you can set go_direct true, or old will be based on opposite of parent_is_proxy directive.")

	getopt.Parse()

	if *version {
		cfg := &Cfg{Version: appVersion, GitRevision: gitRevision}
		fmt.Println(cfg.AppVersion())
		os.Exit(0)
	} else if *help {
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	} else if *listPlugins {
		return Cfg{ListPlugins: true}, nil
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
		return Cfg{}, errors.New("Too many verbose options. The maximum log verbosity level is 2 (-vv or --verbose=2) for errors (0), warnings (1), and info (2)")
	}

	// The flag takes the full version, for forward-compatibility in case we need it in the future,
	// but we only need the major version at the moment.
	atsMajorVersion := uint(0)
	if *atsVersion != "" {
		err := error(nil)
		atsMajorVersion, err = atscfg.GetATSMajorVersionFromATSVersion(*atsVersion)
		if err != nil {
			return Cfg{}, errors.New("parsing ATS version '" + *atsVersion + "': " + err.Error())
		}
	}

	defaultTLSVersions := atscfg.DefaultDefaultTLSVersions

	*defaultTLSVersionsStr = strings.TrimSpace(*defaultTLSVersionsStr)
	if len(*defaultTLSVersionsStr) > 0 {
		defaultTLSVersionsStrs := strings.Split(*defaultTLSVersionsStr, ",")

		defaultTLSVersions = []atscfg.TLSVersion{}
		for _, tlsVersionStr := range defaultTLSVersionsStrs {
			tlsVersion := atscfg.StringToTLSVersion(tlsVersionStr)
			if tlsVersion == atscfg.TLSVersionInvalid {
				return Cfg{}, errors.New("unknown TLS Version '" + tlsVersionStr + "' in '" + *defaultTLSVersionsStr + "'")
			}
			defaultTLSVersions = append(defaultTLSVersions, tlsVersion)
		}
	}

	if !getopt.IsSet(useStrategiesFlagName) {
		*useStrategiesPtr = defaultUseStrategies.String()
	}

	switch *goDirectPtr {
	case "false", "true", "old":
	default:
		return Cfg{}, errors.New(goDirectFlagName + " should be false, true, or old")
	}

	cfg := Cfg{
		LogLocationErr:     logLocationError,
		LogLocationWarn:    logLocationWarn,
		LogLocationInfo:    logLocationInfo,
		LogLocationDebug:   logLocationDebug,
		ListPlugins:        *listPlugins,
		RevalOnly:          *revalOnly,
		Dir:                *dir,
		ViaRelease:         *viaRelease,
		SetDNSLocalBind:    *dnsLocalBind,
		NoOutgoingIP:       *noOutgoingIP,
		ATSMajorVersion:    atsMajorVersion,
		ParentComments:     !(*disableParentConfigComments),
		DefaultEnableH2:    *defaultEnableH2,
		DefaultTLSVersions: defaultTLSVersions,
		Version:            appVersion,
		GitRevision:        gitRevision,
		UseStrategies:      t3cutil.UseStrategiesFlag(*useStrategiesPtr),
		GoDirect:           *goDirectPtr,
		Cache:              *cache,
	}
	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("Initializing loggers: " + err.Error() + "\n")
	}
	return cfg, nil
}

func ValidateURL(u *url.URL) error {
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
